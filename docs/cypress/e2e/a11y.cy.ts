import * as pages from '../fixtures/pages.json'

describe('a11y accessiblity test suite', () => {
  pages.forEach((page) => {
    it(`Should test page ${page} for a11y violations on desktop`, () => {
      cy.viewport('macbook-16')
      cy.visit(page)
      cy.injectAxe()
      cy.wait(100)
      cy.checkA11y(
        undefined,
        {
          includedImpacts: ['critical'],
        },
        (violations) => {
          cy.task(
            'log',
            `${violations.length} accessibility violation${
              violations.length === 1 ? '' : 's'
            } ${violations.length === 1 ? 'was' : 'were'} detected`,
          )
          // pluck specific keys to keep the table readable
          const violationData = violations.map(
            ({ id, impact, description, nodes }) => ({
              id,
              impact,
              description,
              nodes: nodes.length,
            }),
          )

          cy.task('table', violationData)

          // console.error(JSON.stringify(violations));
        },
      )
    })

    it(`Should test page ${page} for a11y violations on mobile`, () => {
      cy.viewport('iphone-x')
      cy.visit(page)
      cy.injectAxe()
      cy.wait(100)
      cy.checkA11y(
        undefined,
        {
          includedImpacts: ['critical'],
        },
        (violations) => {
          cy.task(
            'log',
            `${violations.length} accessibility violation${
              violations.length === 1 ? '' : 's'
            } ${violations.length === 1 ? 'was' : 'were'} detected`,
          )
          // pluck specific keys to keep the table readable
          const violationData = violations.map(
            ({ id, impact, description, nodes }) => ({
              id,
              impact,
              description,
              nodes: nodes.length,
            }),
          )

          cy.task('table', violationData)

          // console.error(JSON.stringify(violations));
        },
      )
    })
  })
})
