describe('Mythicplus Page', () => {
  context('1080p resolution', () => {
    beforeEach(() => {
      // run these tests as if in a desktop
      // browser with a 720p monitor
      cy.viewport(1920, 1080)
    })

    it('should visit the home page', () => {
      cy.visit('/')
    })

    it('should be able to navigate to m+ page', () => {
      cy.get('[data-cy=mythicplus]').click()
    })
    /*
      - no previous button
      - next button
      - click next
      - previous btn exists
      - click prev
      - prev gone again
      - go other side
      - click first log
     */
  })
})
