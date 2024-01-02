Feature: EXTERNAL API

  Background:
    * url baseURLExternal
    * print "Testing host :",baseURLExternal

  Scenario: Check system availability

    Given  path '/wellness'
    When method GET
    Then assert responseStatus == 200
    And assert response.DBStatus == true
