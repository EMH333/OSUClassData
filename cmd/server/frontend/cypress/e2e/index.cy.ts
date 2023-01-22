describe('template spec', () => {
  it('passes', () => {
    cy.visit('/')
    cy.contains('OSU Class Data Explorer')

    // select CS160 on class dropdown
    cy.get('input.autocomplete-input').click()
    cy.get('.autocomplete-list .autocomplete-list-item').contains('CS160').click()
    //cy.get('select').select('CS160')

    // confirm there is a h2 that contains CS160
    cy.get('h2').contains('CS160')
  })
})
