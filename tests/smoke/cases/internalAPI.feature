Feature: INTERNAL API

  Background:
    * url baseURLInternal
    * print "Testing host :",baseURLInternal

  Scenario: Check system availability

    Given  path '/wellness'
    When method GET
    Then assert responseStatus == 200
    And assert response.DBStatus == true