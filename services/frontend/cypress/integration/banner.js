describe('Banner Warning', () => {
  it('is able to upload a file', () => {
    cy.visit('/')
    cy.get('[data-cy=banner]').should('exist')
    cy.get('[data-cy=hide-banner]').should('exist').click()
    cy.get('[data-cy=banner]').should('not.exist')
    cy.reload(true)
    cy.get('[data-cy=banner]').should('not.exist')
  })
})
