Feature: PROFILE API

  Background:
  * url baseURLCompliance

  Scenario: Get compliance profile with ownership structure content using profile ID

    * def createdOwnerShipApprovedShareholding = call read('../../rules/profile/ownershipStructure/ownership_structure.feature@CreateOwnerShipApprovedShareholding')
    * def profileID = createdOwnerShipApprovedShareholding.profileID
    
    Given path '/profile/' + profileID
    And header Content-Type = 'application/json'
    When method GET
    Then assert responseStatus == 200
    And assert response.ownership_structure.shareholders.length == 4
  
  Scenario: Get compliance profile with enriched ownership structure content using profile ID

    * def createdEnrichedOwnerShipApprovedShareholding = call read('../../rules/profile/ownershipStructure/ownership_structure.feature@CreateEnrichedOwnerShipApprovedShareholding')
    * def profileID = createdEnrichedOwnerShipApprovedShareholding.profileID
    
    Given path '/profile/' + profileID
    And header Content-Type = 'application/json'
    When method GET
    Then assert responseStatus == 200
    And assert response.person.enriched_information.ownership_structure.shareholders.length == 3

  Scenario: Get compliance profile with enriched directors content using profile ID

    * def createdProfile = call read('../../rules/profile/boardOfDirectors/boardOfDirectors.feature@CreateProfileWithEnrichedDirectors')
    * def profileID = createdProfile.profileID
    
    Given path '/profile/' + profileID
    And header Content-Type = 'application/json'
    When method GET
    Then assert responseStatus == 200
    And assert response.person.enriched_information.board_of_directors.length == 2

  Scenario: Get compliance profile with registered directors content using profile ID

    * def createdProfile = call read('../../rules/profile/boardOfDirectors/boardOfDirectors.feature@CreateProfileWithRegisteredDirectors')
    * def profileID = createdProfile.profileID
    
    Given path '/profile/' + profileID
    And header Content-Type = 'application/json'
    When method GET
    Then assert responseStatus == 200
    And assert response.board_of_directors.length == 2

  Scenario: Get profile not found

    * def profileID = uuid()

    Given path '/profile/' + profileID
    And header Content-Type = 'application/json'
    When method GET
    Then assert responseStatus == 404
    And match response == {"error": "profile not found"}

  Scenario: Should create and return profile when it does not exist yet

    * def createdLegalRepresentativeApprovedSerasa = call read('../../rules/profile/legalRepresentative/legal_representatives.feature@CreateLegalRepresentativeApprovedSerasa')
    * def profileID = createdLegalRepresentativeApprovedSerasa.profileID
    
    Given url mockURL
    And path '/postgres/delete/compliance_profiles/profile_id/'+profileID
    And request ""
    When method POST
    Then assert responseStatus == 200

    Given url baseURLCompliance
    And path '/profile/' + profileID
    And header Content-Type = 'application/json'
    When method GET
    Then assert responseStatus == 200
    
  Scenario: Get compliance profile with enriched CAF Provider
    * def createProfileEnrichedCAFProvider = call read('classpath:features/rules/person/caf_analysis.feature@CreateProfileEnrichedCAFProvider')
    * def profileID = createProfileEnrichedCAFProvider.profileID
    * def partnerID = createProfileEnrichedCAFProvider.partnerID
    * def requestID =  createProfileEnrichedCAFProvider.requestID
    * def documentNumber = createProfileEnrichedCAFProvider.documentNumber
    * def offer = createProfileEnrichedCAFProvider.CAFRuleOffer

    * def expectedCAFProviderResponse =
      """
        {
          "person": {
            "document_number": "#(documentNumber)",
            "person_type": "INDIVIDUAL",
            "partner_id": "#(partnerID)",
            "offer_type": "#(offer)",
            "profile_id": "#(profileID)",
            "entity_id": "#(profileID)",
            "entity_type": "PROFILE",
            "role_type": "CUSTOMER",
            "enriched_information": {
                "name": "JOAO DA SILVA",
                "birth_date": "25/12/1980",
                "providers": [
                    {
                    "provider_name": "CAF_ENRICHER",
                    "request_id": "#(requestID)",
                    "status": "APPROVED"
                    }
                ]
            }
          }
        }    
      """
    
    Given url baseURLCompliance
    And path '/profile/' + profileID
    And header Content-Type = 'application/json'
    When method GET
    Then assert responseStatus == 200
    And match response contains deep expectedCAFProviderResponse