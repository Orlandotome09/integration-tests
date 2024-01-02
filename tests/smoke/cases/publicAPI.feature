Feature: PUBLIC API

  Background:
    * url baseURLPublic
    * print "Testing host :",baseURLPublic

  Scenario: Check system availability

    Given  path '/wellness'
    When method GET
    Then assert responseStatus == 200
    And assert response.DBStatus == true
