Feature: Create offer

  Background:
  * url baseURLCompliance
  
  * def documentNumber = CPFGenerator()

    Scenario: should create offer

      * def createdOfferType = "TEST_OFFER" + uuid() 
      * def product = 'maquininha'

      Given path '/offers'
      And header Content-Type = 'application/json'
      And def offer = {offer_type: '#(createdOfferType)', product: '#(product)'}
      And request offer
      When method POST
      Then assert responseStatus == 201
      And assert response.offer_type == createdOfferType

    Scenario: should not create offer, required fields

      * def expectedError = "Type is required, Product is required"

      Given path '/offers'
      And header Content-Type = 'application/json'
      And def offer = {}
      And request offer
      When method POST
      Then assert responseStatus == 400
      Then assert response.error == expectedError 

    Scenario: should not create duplicated offer

      * def duplicatedOfferType = "TEST_OFFER" + uuid() 

      Given path '/offers'
      And header Content-Type = 'application/json'
      And def offer = {offer_type: '#(duplicatedOfferType)', product: 'maquininha'}
      And request offer
      When method POST
      Then assert responseStatus == 201

      Given path '/offers'
      And header Content-Type = 'application/json'
      And def offer = {offer_type: '#(duplicatedOfferType)', product: 'maquininha'}
      And request offer
      When method POST
      Then assert responseStatus == 409

  