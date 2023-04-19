describe('Basic Navigation', () => {
  it('Home Page Nav', () => {
    cy.visit('/')
    cy.contains('OSU Class Data Explorer')

    // select CS160 on class dropdown
    cy.get('input.autocomplete-input').click()
    cy.get('.autocomplete-list .autocomplete-list-item').contains('CS160').click()

    // confirm the info card loaded correctly
    cy.get('h2').contains('CS160')
    cy.contains('Credits: 4')
    cy.contains('Computer Science intro is a class designed to do a whole bunch of stuff')
    // confirm there is a .button-link with href=/class/CS160
    cy.get('.button-link').should('have.attr', 'href', '/class/CS160')

    // Go to class page
    cy.get('p .button-link[href="/class/CS160"]').click()

    // wait for class.html to load
    cy.url().should('include', '/class/CS160')

    ////////////////////////////////////////////////////////////class page
    // confirm the info card loaded correctly
    cy.get('h2').contains('CS160')
    cy.contains('Credits: 4')
    cy.contains('Computer Science intro is a class designed to do a whole bunch of stuff')
    // confirm there is a .button-link with link back home
    cy.get('.button-link').should('have.attr', 'href', '/')
  })
})
