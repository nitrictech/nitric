import * as pages from '../fixtures/pages.json'
import * as prodPages from '../fixtures/prod_pages.json'

const CORRECT_CODES = [200]
// site should not return internal redirects
const REDIRECT_CODES = [301, 302, 304, 307, 308]
// other non standard codes, like 999 from linkedin
const OTHER_CODES = [999]

// URLs that we accept 429s for
const ACCEPTED_RATE_LIMITED_URLS = [
  'https://github.com/nitrictech/nitric',
  // Add more URLs here as needed
]

const IGNORED_URLS = [
  'googleads.g.doubleclick.net',
  'youtube.com/api',
  'localhost:49152',
  'localhost:3000',
  'localhost:5000',
  'app.supabase.io',
  'podbay.fm',
  'docs.planetscale.com',
  'turborepo.org',
  'prisma.io',
  'jestjs.io',
  'github.com/nitrictech/go-sdk',
  'https://account.region.auth0.com',
  'https://scoop-docs.vercel.app/',
  'https://vercel.com/new/clone?repository-url=https://github.com/nitrictech/nitric-todo&env=API_BASE_URL',
  'http://localhost:4000',
  'http://localhost:4001',
  'http://localhost:4321',
  'https://www.gutenberg.org/cache/epub/42671/pg42671.txt',
  'https://stackoverflow.com/help/minimal-reproducible-example',
  'https://support.discord.com/hc/en-us/articles/228383668-Intro-to-Webhooks',
  'https://jwt.io',
  'https://portal.azure.com',
]

const rootBaseUrl = Cypress.config('baseUrl')

const isInternalUrl = (url: string) => {
  return (
    url.startsWith(rootBaseUrl) || url.startsWith('./') || url.startsWith('../')
  )
}

const getCleanInternalUrl = (url: string, currentPage: string) => {
  if (url.startsWith(rootBaseUrl)) {
    return url.replace(rootBaseUrl, '')
  }

  // Handle relative paths
  if (url.startsWith('./') || url.startsWith('../')) {
    // Get the directory of the current page
    const currentDir = currentPage.substring(
      0,
      currentPage.lastIndexOf('/') + 1,
    )
    // Resolve the relative path
    const fullPath = new URL(url, `${rootBaseUrl}${currentDir}`).pathname
    return fullPath.replace(rootBaseUrl, '')
  }

  return url
}

const isExternalUrl = (url: string) => {
  return !url.includes('localhost')
}

const isAcceptedRateLimitedUrl = (url: string) => {
  return ACCEPTED_RATE_LIMITED_URLS.some((acceptedUrl) =>
    url.startsWith(acceptedUrl),
  )
}

const req = (
  url: string,
  retryCount = 0,
  followRedirect = false,
  visitedLinks: Record<string, boolean> = {},
): any => {
  return cy
    .request({
      url,
      followRedirect,
      failOnStatusCode: false,
      gzip: false,
    })
    .then((resp) => {
      // Handle rate limiting (429) with exponential backoff
      if (resp.status === 429 && retryCount < 3) {
        const retryAfter = resp.headers['retry-after']
          ? parseInt(
              Array.isArray(resp.headers['retry-after'])
                ? resp.headers['retry-after'][0]
                : resp.headers['retry-after'],
            )
          : null
        const waitTime = retryAfter
          ? retryAfter * 1000
          : Math.min(500 * Math.pow(2, retryCount), 5000)

        cy.log(
          `Rate limited for ${url}, waiting ${waitTime}ms before retry ${retryCount + 1}/3`,
        )
        cy.wait(waitTime)
        return req(url, retryCount + 1, followRedirect, visitedLinks)
      }

      // Handle timeouts with exponential backoff
      if (resp.status === 408 && retryCount < 3) {
        const waitTime = Math.min(200 * Math.pow(2, retryCount), 2000)
        cy.log(
          `Request timeout for ${url}, waiting ${waitTime}ms before retry ${retryCount + 1}/3`,
        )
        cy.wait(waitTime)
        return req(url, retryCount + 1, followRedirect, visitedLinks)
      }

      return resp
    })
}

describe('Broken links test suite', () => {
  const VISITED_SUCCESSFUL_LINKS = {}
  const BATCH_SIZE = 10 // Process links in batches of 10

  pages.forEach((page) => {
    it(`Should visit page ${page} and check all links`, () => {
      cy.viewport('macbook-16')
      cy.visit(page)

      const links = cy.get("a:not([href*='mailto:']),img")

      links
        .filter((_i, link) => {
          const href = link.getAttribute('href')
          const src = link.getAttribute('src')

          return !IGNORED_URLS.some(
            (l) => href?.includes(l) || src?.includes(l),
          )
        })
        .then(($links) => {
          const linkPromises = []
          const linksToCheck = []

          $links.each((_i, link) => {
            const baseUrl =
              link.getAttribute('href') || link.getAttribute('src')
            if (!baseUrl) {
              cy.log('Skipping link with no href/src:', link)
              return
            }

            // Skip if the URL is just a hash fragment
            if (baseUrl.startsWith('#')) {
              cy.log('Skipping hash fragment:', baseUrl)
              return
            }

            const url = baseUrl.split('#')[0]
            if (!url) {
              cy.log('Skipping empty URL from:', baseUrl)
              return
            }

            if (VISITED_SUCCESSFUL_LINKS[url]) {
              cy.log(`Skipping already checked link: ${url}`)
              return
            }

            linksToCheck.push(url)
          })

          // Process links in batches
          for (let i = 0; i < linksToCheck.length; i += BATCH_SIZE) {
            const batch = linksToCheck.slice(i, i + BATCH_SIZE)
            const batchPromises = batch.map((url) => {
              if (!url) {
                cy.log('Skipping empty URL in batch')
                return Promise.resolve()
              }

              if (isInternalUrl(url)) {
                const cleanUrl = getCleanInternalUrl(url, page)
                if (!pages.includes(cleanUrl)) {
                  assert.fail(`${cleanUrl} is not part of the pages fixture`)
                }
                VISITED_SUCCESSFUL_LINKS[url] = true
                return Promise.resolve()
              }

              return req(url, 0, false, VISITED_SUCCESSFUL_LINKS).then(
                (res: Cypress.Response<any>) => {
                  let acceptableCodes = CORRECT_CODES
                  if (
                    REDIRECT_CODES.includes(res.status) &&
                    !isExternalUrl(url)
                  ) {
                    assert.fail(
                      `${url} returned ${res.status} to ${res.headers['location']}`,
                    )
                  } else if (res.status === 429) {
                    // After all retries, if we still get a 429, only mark as successful for accepted URLs
                    if (isAcceptedRateLimitedUrl(url)) {
                      cy.log(
                        `Rate limited for accepted URL ${url} after all retries, marking as successful`,
                      )
                      VISITED_SUCCESSFUL_LINKS[url] = true
                      return
                    } else {
                      assert.fail(
                        `${url} returned 429 (Rate Limited) and is not in the accepted list`,
                      )
                    }
                  } else {
                    acceptableCodes = [
                      ...CORRECT_CODES,
                      ...REDIRECT_CODES,
                      ...OTHER_CODES,
                    ]
                  }

                  if (acceptableCodes.includes(res.status)) {
                    VISITED_SUCCESSFUL_LINKS[url] = true
                  }

                  expect(res.status).oneOf(
                    acceptableCodes,
                    `${url} returned ${res.status}`,
                  )
                },
              )
            })

            linkPromises.push(Promise.all(batchPromises))
          }

          return Promise.all(linkPromises)
        })
    })
  })
})

describe('Current links test suite', () => {
  prodPages.forEach((path) => {
    it(`should visit page ${path} in the current prod sitemap`, function () {
      req(path, 3, true).then((res: Cypress.Response<any>) => {
        expect(res.status).to.be.oneOf(CORRECT_CODES)
      })
    })
  })
})
