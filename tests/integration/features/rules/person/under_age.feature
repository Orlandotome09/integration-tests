Feature: UnderAge Rule

    Background:
        * url baseURLCompliance

        * def UnderAgeRuleOffer = "TEST_OFFER" + uuid() 
        * def offer = { offer_type : '#(UnderAgeRuleOffer)', product: 'maquininha' }

        Given path '/offers'
        And header Content-Type = 'application/json'
        And request offer
        When method POST

        * def ruleSetConfig = { under_age: {} }

        Given url mockURL
        And path '/temis-config/cadastral-validation-configs'
        And header Content-Type = 'application/json'
        And def catalog = call CreateSingleLevelCatalog { offer_type: '#(UnderAgeRuleOffer)', role_type: '#(RoleType.Customer)', person_type: '#(ProfileType.Individual)', account_flag: false, rules_config: #(ruleSetConfig)}
        And request catalog
        When method POST

        * def UnderAgeRuleOffer2 = "TEST_OFFER_2" + uuid() 
        * def offer = { offer_type : '#(UnderAgeRuleOffer2)', product: 'maquininha' }

        Given path '/offers'
        And header Content-Type = 'application/json'
        And request offer
        When method POST

        * def ruleSetConfig = { under_age: { minimum_age: 30} }

        Given url mockURL
        And path '/temis-config/cadastral-validation-configs'
        And header Content-Type = 'application/json'
        And def catalog = call CreateSingleLevelCatalog { offer_type: '#(UnderAgeRuleOffer2)', role_type: '#(RoleType.Customer)', person_type: '#(ProfileType.Individual)', account_flag: false, rules_config: #(ruleSetConfig)}
        And request catalog
        When method POST

        * def profileID = uuid()
		* def expectedCreatedStateEvent = 
		"""
		{
			"id": "#string",
			"entity_id": "#(profileID)",
			"entity_type": "#(EntityType.ComplianceState)",
			"event_type": "#(EventType.State.Created)",
			"update_date":"#string",
			"data": {
				"content": "#(profileStateEventContentSchema)",
			}
		}	
		"""

		* def expectedChangedStateEvent = 
		"""
		{
			"id": "#string",
			"entity_id": "#(profileID)",
			"entity_type": "#(EntityType.ComplianceState)",
			"event_type": "#(EventType.State.Changed)",
			"update_date":"#string",
			"data": {
				"content": "#(profileStateEventContentSchema)",
			}
		}	
		"""

    Scenario: should approve profile individual over 18

        * def documentNumber =  CPFGenerator()
        * def dateOfBirth = generateDateOfBirthForAge(18)
		* def profile = 
		"""
			{ 
                "profile_id": '#(profileID)', 
                "offer_type": '#(UnderAgeRuleOffer)',  
                "role_type": '#(RoleType.Customer)', 
                "profile_type": '#(ProfileType.Individual)', 
                "document_number": '#(documentNumber)',
                "individual":{
                    "date_of_birth": "#(dateOfBirth)"
                }
            }
		"""
        
        Given url mockURL
		And path '/v1/temis/profile/' + profileID
		And request profile
		When method POST
		Then assert responseStatus == 200

		* def event = { profile_id: '#(profileID)', entity_id: '#(profileID)', entity_type: 'PROFILE', event_type: 'PROFILE_CHANGED', parent_type: 'PROFILE', update_date: '2023-08-03T19:29:06.63556Z'  }
		* string json = event
    	
		* def result = RegistrationEventsPublisher.publish(json)

        Given url baseURLCompliance
		And path '/state/' + profileID
        And header Content-Type = 'application/json'
        And params ({ only_pending: false })		
		And retry until response.result == RuleResult.Approved
        When method GET
		Then assert responseStatus == 200

        Given url mockURL
		And path "/subscribe/state-events/" + profileID
		And retry until response.length == 2
		When method GET
		Then assert responseStatus == 200
		And match response contains deep expectedCreatedStateEvent
		And match response contains deep expectedChangedStateEvent

    Scenario: should reject profile individual under 18

        * def documentNumber =  CPFGenerator()
        * def dateOfBirth = generateDateOfBirthForAge(17)
		* def profile = 
		"""
			{ 
                "profile_id": '#(profileID)', 
                "offer_type": '#(UnderAgeRuleOffer)',  
                "role_type": '#(RoleType.Customer)', 
                "profile_type": '#(ProfileType.Individual)', 
                "document_number": '#(documentNumber)',
                "individual":{
                    "date_of_birth":  "#(dateOfBirth)"
                }
            }
		"""
        
        Given url mockURL
		And path '/v1/temis/profile/' + profileID
		And request profile
		When method POST
		Then assert responseStatus == 200

		* def event = { profile_id: '#(profileID)', entity_id: '#(profileID)', entity_type: 'PROFILE', event_type: 'PROFILE_CHANGED', parent_type: 'PROFILE', update_date: '2023-08-03T19:29:06.63556Z'  }
		* string json = event
    	
		* def result = RegistrationEventsPublisher.publish(json)

        Given url baseURLCompliance
		And path '/state/' + profileID
        And header Content-Type = 'application/json'
        And params ({ only_pending: false })		
		And retry until response.result == RuleResult.Rejected
        When method GET
		Then assert responseStatus == 200

        Given url mockURL
		And path "/subscribe/state-events/" + profileID
		And retry until response.length == 2
		When method GET
		Then assert responseStatus == 200
		And match response contains deep expectedCreatedStateEvent
		And match response contains deep expectedChangedStateEvent


    Scenario: should approve profile individual over 30

        * def documentNumber =  CPFGenerator()
        * def dateOfBirth = generateDateOfBirthForAge(30)
		* def profile = 
		"""
			{ 
                "profile_id": '#(profileID)', 
                "offer_type": '#(UnderAgeRuleOffer2)',  
                "role_type": '#(RoleType.Customer)', 
                "profile_type": '#(ProfileType.Individual)', 
                "document_number": '#(documentNumber)',
                "individual":{
                    "date_of_birth": "#(dateOfBirth)"
                }
            }
		"""
        
        Given url mockURL
		And path '/v1/temis/profile/' + profileID
		And request profile
		When method POST
		Then assert responseStatus == 200

		* def event = { profile_id: '#(profileID)', entity_id: '#(profileID)', entity_type: 'PROFILE', event_type: 'PROFILE_CHANGED', parent_type: 'PROFILE', update_date: '2023-08-03T19:29:06.63556Z'  }
		* string json = event
    	
		* def result = RegistrationEventsPublisher.publish(json)

        Given url baseURLCompliance
		And path '/state/' + profileID
        And header Content-Type = 'application/json'
        And params ({ only_pending: false })		
		And retry until response.result == RuleResult.Approved
        When method GET
		Then assert responseStatus == 200

        Given url mockURL
		And path "/subscribe/state-events/" + profileID
		And retry until response.length == 2
		When method GET
		Then assert responseStatus == 200
		And match response contains deep expectedCreatedStateEvent
		And match response contains deep expectedChangedStateEvent

    
    Scenario: should reject profile individual under 30

        * def documentNumber =  CPFGenerator()
        * def dateOfBirth = generateDateOfBirthForAge(29)
		* def profile = 
		"""
			{ 
                "profile_id": '#(profileID)', 
                "offer_type": '#(UnderAgeRuleOffer2)',  
                "role_type": '#(RoleType.Customer)', 
                "profile_type": '#(ProfileType.Individual)', 
                "document_number": '#(documentNumber)',
                "individual":{
                    "date_of_birth": "#(dateOfBirth)"
                }
            }
		"""
        
        Given url mockURL
		And path '/v1/temis/profile/' + profileID
		And request profile
		When method POST
		Then assert responseStatus == 200

		* def event = { profile_id: '#(profileID)', entity_id: '#(profileID)', entity_type: 'PROFILE', event_type: 'PROFILE_CHANGED', parent_type: 'PROFILE', update_date: '2023-08-03T19:29:06.63556Z'  }
		* string json = event
    	
		* def result = RegistrationEventsPublisher.publish(json)

        Given url baseURLCompliance
		And path '/state/' + profileID
        And header Content-Type = 'application/json'
        And params ({ only_pending: false })		
		And retry until response.result == RuleResult.Rejected
        When method GET
		Then assert responseStatus == 200

        Given url mockURL
		And path "/subscribe/state-events/" + profileID
		And retry until response.length == 2
		When method GET
		Then assert responseStatus == 200
		And match response contains deep expectedCreatedStateEvent
		And match response contains deep expectedChangedStateEvent
    
    