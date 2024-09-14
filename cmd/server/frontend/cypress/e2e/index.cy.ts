describe('Basic Navigation', () => {
  it('should list expected classes and data', () => {
    cy.visit('/')
    cy.contains('OSU Class Data Explorer')

    // select CS160 on class dropdown
    cy.get('.subject>.svelecte>.sv-control').click()
    cy.get('.subject .sv-dropdown-content .sv-item--wrap').contains('CS').click()

    cy.get('.class>.svelecte>.sv-control').click()
    cy.get('.class .sv-dropdown-content .sv-item--wrap').contains('CS160').click()

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

  it('should have all subjects on subject page', function () {
    cy.visit('/subject.html')

    // select CS subject
    cy.get('.sv-control').click()
    cy.get('.sv-dropdown-content .sv-item--wrap').contains('CS').click()
  });
})
