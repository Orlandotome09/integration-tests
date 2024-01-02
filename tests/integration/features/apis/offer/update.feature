Feature: OFFER API

  Background:
  * url baseURLCompliance

  * def documentNumber = CPFGenerator()

  Scenario: Update existing offer
    * def offerType = "TEST_OFFER" + uuid() 
    * def product = 'maquininha'
    * def offer = {offer_type: '#(offerType)', product: '#(product)'}

    Given path '/offers'
    And header Content-Type = 'application/json'
    And request offer
    When method POST

    * def createdAt = response.created_at
    * def updatedAt = response.updated_at
    Then assert responseStatus == 201
    And assert response.offer_type == offerType
    And assert response.product == product

    Given path '/offers/' + offerType
    And header Content-Type = 'application/json'
    * def Offer_Update = { product: 'WebPayments'}
    And request Offer_Update
    When method PATCH
    Then assert responseStatus == 200
    And assert response.offer_type == offerType
    And assert response.product == 'WebPayments'

    Given path '/offers/' + offerType
    And header Content-Type = 'application/json'
    When method GET
    Then assert responseStatus == 200
    And assert response.offer_type == offerType
    And assert response.product == 'WebPayments'
    And assert response.created_at == createdAt
    And assert response.updated_at != updatedAt     

  Scenario: Update non existing offer

    * def offerType = 'nonExistingOffer'
    * def product = 'maquininha' 

    Given path '/offers/' + offerType
    And header Content-Type = 'application/json'
    And def offer = { product: '#(product)'}
    And request offer
    When method PATCH
    Then assert responseStatus == 412

  