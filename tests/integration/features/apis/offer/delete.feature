Feature: Delete offer

  Background:
  * url baseURLCompliance
  
  * def documentNumber = CPFGenerator()

  Scenario: Delete existing offer
    
    * def offerType = 'deletingOffer'
    * def product = 'maquininha' 

    Given path '/offers'
    And header Content-Type = 'application/json'
    And def offer = {offer_type: '#(offerType)', product: 'maquininha'}
    And request offer
    When method POST
    Then assert responseStatus == 201

    Given path '/offers/' + offerType
    And header Content-Type = 'application/json'
    When method DELETE
    Then assert responseStatus == 200

  Scenario: Delete non existing offer
    
    * def offerType = 'nonExistingOffer'
    * def product = 'maquininha' 

    Given path '/offers/' + offerType
    And header Content-Type = 'application/json'
    When method DELETE
    Then assert responseStatus == 412

 