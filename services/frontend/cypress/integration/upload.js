describe('Upload Page', () => {
  it('is able to upload a file', () => {
    cy.visit('/upload')

    cy.fixture('otherside-empty.txt', 'base64').then(content => {
      cy.get('#dropzone').upload(content, 'otherside-empty.txt')
    })
    cy.fixture('otherside-empty.txt.gz', 'base64').then(content => {
      cy.get('#dropzone').upload(content, 'otherside-empty.txt.gz')
    })
    cy.fixture('otherside-empty.zip', 'base64').then(content => {
      cy.get('#dropzone').upload(content, 'otherside-empty.zip')
    })
  })
})
