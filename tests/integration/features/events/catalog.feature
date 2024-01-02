Feature: CATALOG API

  	Background:
    * url baseURLCompliance
    * configure retry = { count: 20, interval: 1000 }	
    
    * def catalogOffer = 'CATALOG_OFFER_' + CPFGenerator()
    * def offer = { offer_type : '#(catalogOffer)', product: 'maquininha' }

    Given path '/offers'
    And header Content-Type = 'application/json'
    And request offer
    When method POST
    And assert responseStatus == 201    

    Scenario: Should use the catalog with partner
        
        * def ruleSetConfig = { manual_block: {} }

        * def catalogWithoutPartner = call CreateSingleLevelCatalog { offer_type: '#(catalogOffer)', role_type: '#(RoleType.Customer)', person_type: '#(ProfileType.Individual)', account_flag: false, rules_config: #(ruleSetConfig)}
        
        Given url mockURL
        And path '/temis-config/cadastral-validation-configs'
        And header Content-Type = 'application/json'
        And request catalogWithoutPartner
        When method POST
        Then assert responseStatus == 200

        * def ruleSetConfig = { watchlist: {} }
        
        * def partnerID = uuid()
        * def catalogWithPartner = call CreateSingleLevelCatalog { offer_type: '#(catalogOffer)',  role_type: '#(RoleType.Customer)', person_type: '#(ProfileType.Individual)', partner_id: '#(partnerID)', account_flag: false, rules_config: #(ruleSetConfig)}

        Given url mockURL
        And path '/temis-config/cadastral-validation-configs'
        And header Content-Type = 'application/json'
        And request catalogWithPartner
        When method POST
        Then assert responseStatus == 200
        
        * def profileID = uuid()
        * def documentNumber =  CPFGenerator()
		* def profile = { partner_id: '#(partnerID)', profile_id: '#(profileID)', offer_type: '#(catalogOffer)',  role_type: '#(RoleType.Customer)', profile_type: '#(ProfileType.Individual)', document_number: '#(documentNumber)' }

		Given url mockURL
		And path '/v1/temis/profile/' + profileID
		And request profile
		When method POST
		Then assert responseStatus == 200

        * def event = { profile_id: '#(profileID)', entity_id: '#(profileID)', entity_type: 'PROFILE', event_type: 'PROFILE_CREATED', parent_type: 'PROFILE', update_date: '2023-08-03T19:29:06.63556Z'  }
		* string json = event
        
    	* def result = RegistrationEventsPublisher.publish(json)

        Given url baseURLCompliance
		And path '/state/' + profileID
        And header Content-Type = 'application/json'
        And params ({ only_pending: false })
        And retry until response.rule_set_result.length == 1	
        When method GET
		Then assert responseStatus == 200
		And match response contains deep { entity_id: '#(profileID)', rule_set_result: [ { set: 'WATCHLIST' } ] }

    Scenario: Should use the catalog without partner
        
        * def ruleSetConfig = { manual_block: {} }

        * def catalogWithoutPartner = call CreateSingleLevelCatalog { offer_type: '#(catalogOffer)',  role_type: '#(RoleType.Customer)', person_type: '#(ProfileType.Individual)', account_flag: false, rules_config: #(ruleSetConfig)}
        
        Given url mockURL
        And path '/temis-config/cadastral-validation-configs'
        And header Content-Type = 'application/json'
        And request catalogWithoutPartner
        When method POST
        Then assert responseStatus == 200

        * def ruleSetConfig = { watchlist: {} }
        
        * def anyPartnerID = uuid()
        * def catalogWithPartner = call CreateSingleLevelCatalog { offer_type: '#(catalogOffer)',  role_type: '#(RoleType.Customer)', person_type: '#(ProfileType.Individual)', partner_id: '#(anyPartnerID)', account_flag: false, rules_config: #(ruleSetConfig)}
        
        Given url mockURL
        And path '/temis-config/cadastral-validation-configs'
        And header Content-Type = 'application/json'
        And request catalogWithPartner
        When method POST
        Then assert responseStatus == 200
        
        * def profileID = uuid()
        * def partnerID = uuid()
        * def documentNumber =  CPFGenerator()
		* def profile = { partner_id: '#(partnerID)', profile_id: '#(profileID)', offer_type: '#(catalogOffer)',  role_type: '#(RoleType.Customer)', profile_type: '#(ProfileType.Individual)', document_number: '#(documentNumber)' }

		Given url mockURL
		And path '/v1/temis/profile/' + profileID
		And request profile
		When method POST
		Then assert responseStatus == 200

        * def event = { profile_id: '#(profileID)', entity_id: '#(profileID)', entity_type: 'PROFILE', event_type: 'PROFILE_CREATED', parent_type: 'PROFILE', update_date: '2023-08-03T19:29:06.63556Z'  }
		* string json = event
    	
        * def result = RegistrationEventsPublisher.publish(json)

        Given url baseURLCompliance
		And path '/state/' + profileID
        And header Content-Type = 'application/json'
        And params ({ only_pending: false })
        And retry until response.rule_set_result.length == 1			
        When method GET
		Then assert responseStatus == 200
		And match response contains deep { entity_id: '#(profileID)', rule_set_result: [ { set: 'MANUAL_BLOCK' } ] }

    Scenario: Should use the catalog with correct partner
        
        * def ruleSetConfig = { manual_block: {} }

        * def anyPartnerID = uuid()
        * def catalogWithAnyPartner = call CreateSingleLevelCatalog { offer_type: '#(catalogOffer)',  role_type: '#(RoleType.Customer)', person_type: '#(ProfileType.Individual)', partner_id: '#(anyPartnerID)', account_flag: false, rules_config: #(ruleSetConfig)}
        
        Given url mockURL
        And path '/temis-config/cadastral-validation-configs'
        And header Content-Type = 'application/json'
        And request catalogWithAnyPartner
        When method POST
        Then assert responseStatus == 200

        * def ruleSetConfig = { watchlist: {} }
        
        * def correctPartnerID = uuid()
        * def catalogWithCorrectPartner = call CreateSingleLevelCatalog { offer_type: '#(catalogOffer)',  role_type: '#(RoleType.Customer)', person_type: '#(ProfileType.Individual)', partner_id: '#(correctPartnerID)', account_flag: false, rules_config: #(ruleSetConfig)}
        
        Given url mockURL
        And path '/temis-config/cadastral-validation-configs'
        And header Content-Type = 'application/json'
        And request catalogWithCorrectPartner
        When method POST
        Then assert responseStatus == 200
        
        * def profileID = uuid()
        * def documentNumber =  CPFGenerator()
		* def profile = { partner_id: '#(correctPartnerID)', profile_id: '#(profileID)',  offer_type: '#(catalogOffer)', role_type: '#(RoleType.Customer)', profile_type: '#(ProfileType.Individual)', document_number: '#(documentNumber)' }

		Given url mockURL
		And path '/v1/temis/profile/' + profileID
		And request profile
		When method POST
		Then assert responseStatus == 200

        * def event = { profile_id: '#(profileID)', entity_id: '#(profileID)', entity_type: 'PROFILE', event_type: 'PROFILE_CREATED', parent_type: 'PROFILE', update_date: '2023-08-03T19:29:06.63556Z'  }
		* string json = event
    	
        * def result = RegistrationEventsPublisher.publish(json)

        Given url baseURLCompliance
		And path '/state/' + profileID
        And header Content-Type = 'application/json'
        And params ({ only_pending: false })
        And retry until response.rule_set_result.length == 1		
        When method GET
		Then assert responseStatus == 200
		And match response contains deep { entity_id: '#(profileID)', rule_set_result: [ { set: 'WATCHLIST' } ] }

