Feature: OFFER API

  Background:
  * url baseURLCompliance
  * def documentNumber = CPFGenerator()
  
  Scenario: List offers

    Given path '/offers'
    And header Content-Type = 'application/json'
    And def offer = {offer_type: 'OFFER1', product: 'maquininha'}
    And request offer
    When method POST

    Given path '/offers/' + 'OFFER1'
    And header Content-Type = 'application/json'
    When method GET
    
    * def offer1 = response

    Given path '/offers'
    And header Content-Type = 'application/json'
    And def offer = {offer_type: 'OFFER2', product: 'maquininha'}
    And request offer
    When method POST

    Given path '/offers/' + 'OFFER2'
    And header Content-Type = 'application/json'
    When method GET

    * def offer2 = response 

    Given path '/offers'
    And header Content-Type = 'application/json'
    When method GET
    Then assert responseStatus == 200
    And match response contains offer1
    And match response contains offer2


  Scenario: List no offers

    Given path '/offers'
    And header Content-Type = 'application/json'
    When method GET
    Then assert responseStatus == 200
    And match response contains []