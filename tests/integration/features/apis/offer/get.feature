Feature: Get offer

  Background:
  * url baseURLCompliance
  * def documentNumber = CPFGenerator()

  Scenario: Get non existing offer

    * def offerType = 'nonExistingOffer'

    Given path '/offers/' + offerType
    And header Content-Type = 'application/json'
    When method GET
    Then assert responseStatus == 404
    

 Scenario: get offer

    * def offerType = "TEST_OFFER" + uuid() 
    * def product = 'maquininha'

    Given path '/offers'
    And header Content-Type = 'application/json'
    And def offer = {offer_type: '#(offerType)', product: '#(product)'}
    And request offer
    When method POST
    Then assert responseStatus == 201
    And assert response.offer_type == offerType

    Given path '/offers/' + offerType
    And header Content-Type = 'application/json'
    When method GET
    Then assert responseStatus == 200
    And assert response.offer_type == offerType
  