// Checks the built site (run after `next build`) without a browser:
// - every internal link and image points at a built page or a public file
// - every page has a canonical URL, and an og:url for itself
// - every page in the production sitemap still exists, or redirects somewhere
// External links are checked by the nightly Cypress run instead.
import fs from 'fs/promises'
import path from 'path'

import { XMLParser } from 'fast-xml-parser'

const APP_DIR = '.next/server/app'
const PUBLIC_DIR = 'public'
const PROD_SITEMAP_URL = 'https://nitric.io/docs/sitemap-0.xml'
const SITE_ORIGIN = 'https://nitric.io'

const routesManifest = JSON.parse(
  await fs.readFile('.next/routes-manifest.json', 'utf8'),
)
const basePath = routesManifest.basePath
const redirectPatterns = routesManifest.redirects
  .filter((redirect) => !redirect.internal)
  .map((redirect) => new RegExp(redirect.regex))

async function listHtmlFiles(dir) {
  const files = []
  for (const entry of await fs.readdir(dir, { withFileTypes: true })) {
    const entryPath = path.join(dir, entry.name)
    if (entry.isDirectory()) {
      files.push(...(await listHtmlFiles(entryPath)))
    } else if (entry.name.endsWith('.html')) {
      files.push(entryPath)
    }
  }
  return files
}

function pagePathForFile(file) {
  const route = path
    .relative(APP_DIR, file)
    .replace(/\.html$/, '')
    .replace(/(^|\/)index$/, '')
  return route ? `${basePath}/${route}` : basePath
}

async function publicFileExists(urlPath) {
  const file = path.join(PUBLIC_DIR, urlPath.slice(basePath.length))
  return fs
    .stat(file)
    .then((stat) => stat.isFile())
    .catch(() => false)
}

function decodeEntities(value) {
  return value.replace(/&amp;/g, '&').replace(/&quot;/g, '"')
}

// Resolves an href/src to a path on this site, or null when it points elsewhere
function siteUrlPath(reference, pagePath) {
  if (reference.startsWith('#') || reference.startsWith('mailto:')) {
    return null
  }

  const url = new URL(reference, `${SITE_ORIGIN}${pagePath}`)
  if (url.origin !== SITE_ORIGIN || !url.pathname.startsWith(basePath)) {
    return null
  }

  // Next's optimized images and bundles are build artifacts, not authored links
  if (url.pathname.startsWith(`${basePath}/_next/`)) {
    return null
  }

  return url.pathname
}

function attribute(html, pattern) {
  return html.match(pattern)?.[1]
}

const files = (await listHtmlFiles(APP_DIR)).filter(
  (file) => !file.includes('_not-found'),
)
const pages = new Map()
for (const file of files) {
  pages.set(pagePathForFile(file), await fs.readFile(file, 'utf8'))
}

const failures = []

for (const [pagePath, html] of pages) {
  for (const [, reference] of html.matchAll(
    /<(?:a|img)\s[^>]*?\b(?:href|src)="([^"]*)"/g,
  )) {
    const target = siteUrlPath(decodeEntities(reference), pagePath)
    if (
      target === null ||
      pages.has(target) ||
      (await publicFileExists(target))
    ) {
      continue
    }
    failures.push(`${pagePath}: broken link to ${reference}`)
  }

  // The canonical URL may point elsewhere, e.g. guides republished from another site
  if (!attribute(html, /<link rel="canonical" href="([^"]*)"/)) {
    failures.push(`${pagePath}: missing canonical URL`)
  }

  const ogUrl = attribute(html, /<meta property="og:url" content="([^"]*)"/)
  if (!ogUrl) {
    failures.push(`${pagePath}: missing og:url`)
  } else if (new URL(ogUrl).pathname !== pagePath) {
    failures.push(`${pagePath}: og:url is ${ogUrl}`)
  }
}

const sitemap = new XMLParser().parse(
  await fetch(PROD_SITEMAP_URL).then((res) => res.text()),
)
const prodPaths = sitemap.urlset.url.map((entry) =>
  new URL(entry.loc).pathname.replace('/_index', ''),
)
for (const prodPath of prodPaths) {
  const redirected = redirectPatterns.some((pattern) => pattern.test(prodPath))
  if (!pages.has(prodPath) && !redirected) {
    failures.push(`${prodPath}: in the production sitemap but no longer exists`)
  }
}

console.log(
  `Checked ${pages.size} pages and ${prodPaths.length} production sitemap paths`,
)
if (failures.length > 0) {
  console.error(failures.join('\n'))
  console.error(`${failures.length} problem(s) found`)
  process.exit(1)
}
