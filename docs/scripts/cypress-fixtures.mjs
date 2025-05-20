import fs from 'fs/promises'

import { XMLParser } from 'fast-xml-parser'
const parser = new XMLParser()

import { allDocuments } from '../.contentlayer/generated/index.mjs'
import basePages from '../src/assets/sitemap.json' with { type: 'json' }

const FIXTURE_PATH = 'cypress/fixtures/pages.json'

const PROD_PAGES_PATH = 'cypress/fixtures/prod_pages.json'

async function run() {
  try {
    const staticPages = basePages.map((page) => `/docs${page}`)

    const docPages = allDocuments.map((doc) => {
      return !doc.slug ? '/docs' : `/docs/${doc.slug}`
    })

    const paths = [...staticPages, ...docPages].sort((a, b) =>
      a.localeCompare(b),
    )

    console.log(`${paths.length} paths found. Generating fixture`)

    await fs.writeFile(FIXTURE_PATH, JSON.stringify(paths))

    // get current prod pages
    const PROD_BASE_URL = 'https://nitric.io/docs'

    const xml = await fetch(`${PROD_BASE_URL}/sitemap-0.xml`).then((res) =>
      res.text(),
    )

    const jsonObj = parser.parse(xml, false)

    const prodPaths = jsonObj.urlset.url
      .map((p) =>
        p.loc === PROD_BASE_URL
          ? '/docs'
          : `/docs${p.loc
              .substring(PROD_BASE_URL.length, p.loc.length)
              .replace('/_index', '')}`,
      )
      .sort((a, b) => a.localeCompare(b))

    console.log(`${prodPaths.length} prod paths found. Generating fixture`)

    await fs.writeFile(PROD_PAGES_PATH, JSON.stringify(prodPaths))
  } catch (error) {
    console.log(error.message)

    process.exit(1)
  }
}

run()
