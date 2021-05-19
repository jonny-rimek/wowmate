describe('Mythicplus Page - Desktop', () => {
  context('1080p resolution', () => {
    beforeEach(() => {
      cy.viewport(1920, 1080)
    })

    it('is able to paginate', () => {
      cy.visit('/')
      cy.get('[data-cy=mythicplus]').click()
      cy.get('[data-cy=dungeon]').first().should("have.class", "data-cy-active")
      cy.get('[data-cy=prev]').should('not.exist')
      cy.get('[data-cy=next]').should('exist').click()
      cy.get('[data-cy=dungeon]').first().should("have.class", "data-cy-active")
      cy.get('[data-cy=prev]').should('exist').click()
      cy.get('[data-cy=prev]').should('not.exist')
      cy.get('[data-cy=dungeon]').first().should("have.class", "data-cy-active")
    })

    it('can check a logs for a specific dungeon', () => {
      cy.visit('/mythicplus/all')
      cy.get('[data-cy=dungeon]').first().click()
      cy.get('[data-cy=log]').first().click()
    })
  })
})

describe('Mythicplus Page - Mobile', () => {
  // default viewport - iphone 8

  it('is able to click away the mobile menu', () => {
    cy.visit('/')
    cy.get('[data-cy=mobile-menu]').should('not.be.visible')
    cy.get('[data-cy=mobile-menu-btn]').click()
    cy.get('[data-cy=mobile-menu]').should('be.visible')
    cy.get('h2').click()
    cy.get('[data-cy=mobile-menu]').should('not.be.visible')
  })

  it('is able to paginate', () => {
    cy.visit('/')
    cy.get('[data-cy=mobile-menu-btn]').click()
    cy.get('[data-cy=mythicplus-mobile]').click()
    cy.get('[data-cy=dungeon-mobile]').first().click()
    cy.get('[data-cy=prev]').should('not.exist')
    cy.get('[data-cy=next]').should('exist').click({ force: true})
    cy.get('[data-cy=prev]').should('exist').click({ force: true})
    cy.get('[data-cy=prev]').should('not.exist')
  })

  it('can check a logs for a specific dungeon', () => {
    cy.visit('/mythicplus/all')
    cy.get('[data-cy=log]').first().click()
    cy.get('[data-cy=dungeon-name]').should('exist')
  })
})
