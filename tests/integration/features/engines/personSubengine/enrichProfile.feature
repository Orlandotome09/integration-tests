Feature: Enrichment of profile

Background:
    * url baseURLCompliance

    * def offerType = "TEST_OFFER" + uuid() 

    * def offer = { offer_type : '#(offerType)', product: 'maquininha' }

    Given path '/offers'
    And header Content-Type = 'application/json'
    And request offer
    When method POST

    Given url mockURL
    And path '/temis-config/cadastral-validation-configs'
    And header Content-Type = 'application/json'
    And def catalog = call CreateSingleLevelCatalog { offer_type: '#(offerType)', role_type: '#(RoleType.Customer)', person_type: '#(ProfileType.Individual)', account_flag: false, enrich_flag: true}
    And request catalog
    When method POST

    * def profileID = uuid()

Scenario: should enrich profile

    * def profileID = uuid()
    * def profileDocumentNumber =  DocumentNormalizer(CPFGenerator())
    * def profile = { profile_id: '#(profileID)', offer_type: '#(offerType)', role_type: '#(RoleType.Customer)', profile_type: '#(ProfileType.Individual)', document_number: '#(profileDocumentNumber)', "callback_urls": [{"url": "https://","notification_type": "COMPLIANCE_STATUS"}], individual: {first_name: "aaaa"}} 
 
    Given url mockURL
    And path '/v1/temis/profile/' + profileID
    And request profile
    When method POST
    Then assert responseStatus == 200

    * def enrichmentResponse = { situation: 1, name: "SomeName"}

    Given url mockURL
    And path '/temis-enrichment/individual/' + profileDocumentNumber
    And request enrichmentResponse
    When method POST

    * def event = { profile_id: '#(profileID)', entity_id: '#(profileID)', entity_type: 'PROFILE', event_type: 'PROFILE_CHANGED', parent_type: 'PROFILE', update_date: '2023-08-03T19:29:06.63556Z'  }
    * string json = event
    * def result = RegistrationEventsPublisher.publish(json)
    * eval sleep(1000)		

    Given url baseURLCompliance
    And path '/state/' + profileID
    And header Content-Type = 'application/json'
    And params ({ only_pending: false })		
    And retry until response.result == RuleResult.Approved
    When method GET
    Then assert responseStatus == 200