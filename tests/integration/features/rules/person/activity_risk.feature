Feature: Activity Risk Rule

    Background:
        * url baseURLCompliance

	    * def ActivityRiskRuleOffer = "TEST_OFFER" + uuid() 

        * def problemCodeEconomicalActivityRiskUndefined = "ECONOMICAL_ACTIVITY_RISK_UNDEFINED"
        * def problemCodeEconomicalActivityRiskHigh = "ECONOMICAL_ACTIVITY_RISK_HIGH"

        * def highRiskActivity = "9900-8/00"
    
		* def offer = { offer_type : '#(ActivityRiskRuleOffer)', product: 'maquininha' }

        Given path '/offers'
		And header Content-Type = 'application/json'
		And request offer
		When method POST

        * def ruleSetConfig = { activity_risk: {} }

        Given url mockURL
        And path '/temis-config/cadastral-validation-configs'
		And header Content-Type = 'application/json'
		And def catalog = call CreateSingleLevelCatalog { offer_type: '#(ActivityRiskRuleOffer)', role_type: '#(RoleType.Customer)', person_type: '#(ProfileType.Company)', account_flag: false, rules_config: #(ruleSetConfig)}
		And request catalog
		When method POST

        Given url mockURL
        And path '/temis-config/cadastral-validation-configs'
		And header Content-Type = 'application/json'
		And def catalog = call CreateSingleLevelCatalog { offer_type: '#(ActivityRiskRuleOffer)',  role_type: '#(RoleType.Customer)', person_type: '#(ProfileType.Individual)', account_flag: false, rules_config: #(ruleSetConfig)}
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

	Scenario: Should not create state due profile is not company (error executing rule)

		* def documentNumber =  CPFGenerator()
		* def profile = 
		"""
			{ profile_id: '#(profileID)', offer_type: '#(ActivityRiskRuleOffer)',  role_type: '#(RoleType.Customer)', profile_type: '#(ProfileType.Individual)', document_number: '#(documentNumber)' }
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
		And retry until response.result == RuleResult.Created
        When method GET
		Then assert responseStatus == 200
		Then assert response.result == RuleResult.Created

		Given url mockURL
		And path "/subscribe/state-events/" + profileID
		And retry until response.length == 1
		When method GET
		Then assert responseStatus == 200
		And match response contains deep expectedCreatedStateEvent

	Scenario: should not find company in bureau
		
		* def documentNumber =  DocumentNormalizer(CNPJGenerator())
		* def profile = { profile_id: '#(profileID)', offer_type: '#(ActivityRiskRuleOffer)',  role_type: '#(RoleType.Customer)', profile_type: '#(ProfileType.Company)', document_number: '#(documentNumber)' }

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
		* def ruleResults = response.rule_set_result
		Then assert responseStatus == 200
		Then assert response.result == RuleResult.Rejected
		Then assert ruleResults[0].name == RuleName.HighRiskActivity
		Then assert ruleResults[0].result == RuleResult.Rejected
		Then assert ruleResults[0].metadata == "company (" + documentNumber + ") not found in bureau"

		Given url mockURL
		And path "/subscribe/state-events/" + profileID
		And retry until response.length == 2
		When method GET
		Then assert responseStatus == 200
		And match response contains deep expectedCreatedStateEvent
		And match response contains deep expectedChangedStateEvent

    Scenario: should find high risk activity    	

		* def documentNumber =  DocumentNormalizer(CNPJGenerator())
		* def profile = { profile_id: '#(profileID)', offer_type: '#(ActivityRiskRuleOffer)',  role_type: '#(RoleType.Customer)', profile_type: '#(ProfileType.Company)', document_number: '#(documentNumber)' }

		Given url mockURL
		And path '/v1/temis/profile/' + profileID
		And request profile
		When method POST
		Then assert responseStatus == 200

		* def enrichmentResponse = { cnae: '#(highRiskActivity)' }

		Given url mockURL
		And path '/temis-enrichment/legal-entity/' + documentNumber
		And request enrichmentResponse
		When method POST

		* def event = { profile_id: '#(profileID)', entity_id: '#(profileID)', entity_type: 'PROFILE', event_type: 'PROFILE_CHANGED', parent_type: 'PROFILE', update_date: '2023-08-03T19:29:06.63556Z'  }
		* string json = event
    	
		* def result = RegistrationEventsPublisher.publish(json)

		Given url baseURLCompliance
		And path '/state/' + profileID
        And header Content-Type = 'application/json'
        And params ({ only_pending: false })		
		And retry until response.result == RuleResult.Analysing
        When method GET
		* def ruleResults = response.rule_set_result
		Then assert response.result == RuleResult.Analysing
		Then assert ruleResults[0].name == RuleName.HighRiskActivity
		Then assert ruleResults[0].result == RuleResult.Analysing
        Then assert ruleResults[0].pending == true
        Then assert ruleResults[0].metadata != null
		Then assert ruleResults[0].problems[0].code == problemCodeEconomicalActivityRiskHigh
		Then assert ruleResults[0].problems[0].detail.activity_code == highRiskActivity

		Given url mockURL
		And path "/subscribe/state-events/" + profileID
		And retry until response.length == 2
		When method GET
		Then assert responseStatus == 200
		And match response contains deep expectedCreatedStateEvent
		And match response contains deep expectedChangedStateEvent

	Scenario: should not find high risk activity
		
		* def documentNumber =  DocumentNormalizer(CNPJGenerator())
		* def profile = { profile_id: '#(profileID)', offer_type: '#(ActivityRiskRuleOffer)',  role_type: '#(RoleType.Customer)', profile_type: '#(ProfileType.Company)', document_number: '#(documentNumber)' }

		Given url mockURL
		And path '/v1/temis/profile/' + profileID
		And request profile
		When method POST
		Then assert responseStatus == 200

		* def enrichmentResponse = { cnae: "11111" }

		Given url mockURL
		And path '/temis-enrichment/legal-entity/' + documentNumber
		And request enrichmentResponse
		When method POST

		* def event = { profile_id: '#(profileID)', entity_id: '#(profileID)', entity_type: 'PROFILE', event_type: 'PROFILE_CHANGED', parent_type: 'PROFILE', update_date: '2023-08-03T19:29:06.63556Z'  }
		* string json = event
    	
		* def result = RegistrationEventsPublisher.publish(json)

		Given url baseURLCompliance
		And path '/state/' + profileID
        And header Content-Type = 'application/json'
        And params ({ only_pending: false })		
		And retry until response.result == RuleResult.Approved
        When method GET
		* def ruleResults = response.rule_set_result
		Then assert responseStatus == 200
		Then assert response.result == RuleResult.Approved
		Then assert ruleResults[0].name == RuleName.HighRiskActivity
		Then assert ruleResults[0].result == RuleResult.Approved

		Given url mockURL
		And path "/subscribe/state-events/" + profileID
		And retry until response.length == 2
		When method GET
		Then assert responseStatus == 200
		And match response contains deep expectedCreatedStateEvent
		And match response contains deep expectedChangedStateEvent

	