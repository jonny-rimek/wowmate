describe('Upload Page', () => {
  it('is able to upload a file', () => {
    cy.visit('/upload')

    cy.fixture('otherside-empty.txt', 'base64').then(content => {
      cy.get('#dropzone').upload(content, 'otherside-empty.txt')
    })
  })
})
