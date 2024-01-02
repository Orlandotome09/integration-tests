Feature: Bureau Rule

    Background:
        * url baseURLCompliance

	    * def BureauRuleOffer = "TEST_OFFER" + uuid() 
		
		* def problemCodeBureauStatusNotRegular = "BUREAU_STATUS_NOT_REGULAR"
		* def bureauStatusPendingRegularization = "PENDING_REGULARIZATION"
	
		* def offer = { offer_type : '#(BureauRuleOffer)', product: 'maquininha' }

        Given path '/offers'
		And header Content-Type = 'application/json'
		And request offer
		When method POST

        * def ruleSetConfig = { bureau: {} }

        Given url mockURL
        And path '/temis-config/cadastral-validation-configs'
		And header Content-Type = 'application/json'
		And def catalog = call CreateSingleLevelCatalog { offer_type: '#(BureauRuleOffer)', role_type: 'CUSTOMER', person_type: 'INDIVIDUAL', account_flag: false, rules_config: #(ruleSetConfig)}
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

	Scenario: Should REJECT a person NOT FOUND in Bureau

		* def documentNumber =  CPFGenerator()
		* def profile = { profile_id: '#(profileID)', offer_type: '#(BureauRuleOffer)', role_type: 'CUSTOMER', profile_type: 'INDIVIDUAL', document_number: '#(documentNumber)' }

		Given url mockURL
		And path '/v1/temis/profile/' + profileID
		And request profile
		When method POST
		Then assert responseStatus == 200

		* def event = { profile_id: '#(profileID)', entity_id: '#(profileID)', entity_type: 'PROFILE', event_type: 'PROFILE_CHANGED', parent_type: 'PROFILE', update_date: '2023-08-03T19:29:06.63556Z'  }
		* string json = event
    	
		* def result = RegistrationEventsPublisher.publish(json)

		* def expectedResponse =
		"""
			{
				result: '#(RuleResult.Rejected)',
				rule_set_result: [
					{
						set: '#(RuleSet.Bureau)',										
						name: '#(RuleName.CustomerNotFoundInSerasa)',										
						result: '#(RuleResult.Rejected)'										
					},
					{
						set: '#(RuleSet.Bureau)',										
						name: '#(RuleName.CustomerHasProblemsInSerasa)',										
						result: '#(RuleResult.Ignored)'										
					}
				]
			}
		"""

		Given url baseURLCompliance
		And path '/state/' + profileID
        And header Content-Type = 'application/json'
        And params ({ only_pending: false })		
		And retry until response.result == RuleResult.Rejected
        When method GET
		Then assert responseStatus == 200
		And match response contains deep expectedResponse

		Given url mockURL
		And path "/subscribe/state-events/" + profileID
		And retry until response.length == 2
		When method GET
		Then assert responseStatus == 200
		And match response contains deep expectedCreatedStateEvent
		And match response contains deep expectedChangedStateEvent

	Scenario: Given a CUSTOM RESULT is set, should APPROVE a person NOT FOUND in Bureau
		
		* def customOfferType = "TEST_CUSTOM_RESULT_NOT_FOUND_BUREAU_RULE" 
		* def customOffer = { offer_type : '#(customOfferType)', product: 'maquininha' }

		Given url baseURLCompliance
		And path '/offers'
		And header Content-Type = 'application/json'
		And request customOffer
		When method POST
		Then assert responseStatus == 201 || responseStatus == 409

		* def ruleSetConfig = { bureau: {not_found_in_serasa_status: '#(RuleResult.Approved)', not_found_in_serasa_pending: false} }

		Given url mockURL
        And path '/temis-config/cadastral-validation-configs'
		And header Content-Type = 'application/json'
		And def catalog = call CreateSingleLevelCatalog { offer_type: '#(customOfferType)', role_type: 'CUSTOMER', person_type: 'INDIVIDUAL', account_flag: false, rules_config: #(ruleSetConfig)}
		And request catalog
		When method POST

		* def documentNumber =  CPFGenerator()
		* def profile = { profile_id: '#(profileID)', offer_type: '#(customOfferType)', role_type: 'CUSTOMER', profile_type: 'INDIVIDUAL', document_number: '#(documentNumber)' }

		Given url mockURL
		And path '/v1/temis/profile/' + profileID
		And request profile
		When method POST
		Then assert responseStatus == 200

		* def event = { profile_id: '#(profileID)', entity_id: '#(profileID)', entity_type: 'PROFILE', event_type: 'PROFILE_CHANGED', parent_type: 'PROFILE', update_date: '2023-08-03T19:29:06.63556Z'  }
		* string json = event
    	
		* def result = RegistrationEventsPublisher.publish(json)

		* def expectedResponse =
		"""
			{
				result: '#(RuleResult.Approved)',
				rule_set_result: [
					{
						name: '#(RuleName.CustomerNotFoundInSerasa)',										
						result: '#(RuleResult.Approved)',										
						pending: false
					},
					{
						name: '#(RuleName.CustomerHasProblemsInSerasa)',										
						result: '#(RuleResult.Ignored)',
						pending: false										
					}
				]
			}
		"""

		Given url baseURLCompliance
		And path '/state/' + profileID
        And header Content-Type = 'application/json'
        And params ({ only_pending: false })		
		And retry until response.result == RuleResult.Approved
        When method GET
		Then assert responseStatus == 200
		And match response contains deep expectedResponse

		Given url mockURL
		And path "/subscribe/state-events/" + profileID
		And retry until response.length == 2
		When method GET
		Then assert responseStatus == 200
		And match response contains deep expectedCreatedStateEvent
		And match response contains deep expectedChangedStateEvent

	Scenario: Should REJECT a person with CANCELLED status in Bureau
		
		* def documentNumber =  CPFGenerator()
		* def profile = { profile_id: '#(profileID)', offer_type: '#(BureauRuleOffer)', role_type: 'CUSTOMER', profile_type: 'INDIVIDUAL', document_number: '#(documentNumber)' }

		Given url mockURL
		And path '/v1/temis/profile/' + profileID
		And request profile
		When method POST
		Then assert responseStatus == 200

		* def enrichmentResponse = { situation: 4}

		Given url mockURL
		And path '/temis-enrichment/individual/' + documentNumber
		And request enrichmentResponse
		When method POST

		* def event = { profile_id: '#(profileID)', entity_id: '#(profileID)', entity_type: 'PROFILE', event_type: 'PROFILE_CHANGED', parent_type: 'PROFILE', update_date: '2023-08-03T19:29:06.63556Z'  }
		* string json = event
    	
		* def result = RegistrationEventsPublisher.publish(json)

		* def expectedResponse =
		"""
			{
				result: '#(RuleResult.Rejected)',
				rule_set_result: [
					{
						name: '#(RuleName.CustomerNotFoundInSerasa)',										
						result: '#(RuleResult.Approved)'										
					},
					{
						name: '#(RuleName.CustomerHasProblemsInSerasa)',										
						result: '#(RuleResult.Rejected)',										
						problems: [
							{
								code: '#(problemCodeBureauStatusNotRegular)'											
							}
						]										
					}
				]		
			}
		"""

		Given url baseURLCompliance
		And path '/state/' + profileID
        And header Content-Type = 'application/json'
        And params ({ only_pending: false })		
		And retry until response.result == RuleResult.Rejected
        When method GET
		Then assert responseStatus == 200
		And match response contains deep expectedResponse

		Given url mockURL
		And path "/subscribe/state-events/" + profileID
		And retry until response.length == 2
		When method GET
		Then assert responseStatus == 200
		And match response contains deep expectedCreatedStateEvent
		And match response contains deep expectedChangedStateEvent
	
	Scenario: Given a CUSTOM RESULT is set, should APPROVE a person with CANCELLED status in Bureau

		* def customOfferType = "TEST_CUSTOM_HAS_PROBLEMS_RESULT_FOR_BUREAU_RULE" 
		* def customOffer = { offer_type : '#(customOfferType)' , product: 'maquininha'}

		Given url baseURLCompliance
		And path '/offers'
		And header Content-Type = 'application/json'
		And request customOffer
		When method POST
		Then assert responseStatus == 201 || responseStatus == 409

		* def ruleSetConfig = { bureau: {has_problems_in_serasa_status: '#(RuleResult.Approved)', has_problems_in_serasa_pending: false} }

		Given url mockURL
        And path '/temis-config/cadastral-validation-configs'
		And header Content-Type = 'application/json'
		And def catalog = call CreateSingleLevelCatalog { offer_type: '#(customOfferType)', role_type: 'CUSTOMER', person_type: 'INDIVIDUAL', account_flag: false, rules_config: #(ruleSetConfig)}
		And request catalog
		When method POST

		* def documentNumber =  CPFGenerator()
		* def profile = { profile_id: '#(profileID)', offer_type: '#(customOfferType)',  role_type: 'CUSTOMER', profile_type: 'INDIVIDUAL', document_number: '#(documentNumber)' }

		Given url mockURL
		And path '/v1/temis/profile/' + profileID
		And request profile
		When method POST
		Then assert responseStatus == 200

		* def enrichmentResponse = { situation: 4}

		Given url mockURL
		And path '/temis-enrichment/individual/' + documentNumber
		And request enrichmentResponse
		When method POST

		* def event = { profile_id: '#(profileID)', entity_id: '#(profileID)', entity_type: 'PROFILE', event_type: 'PROFILE_CHANGED', parent_type: 'PROFILE', update_date: '2023-08-03T19:29:06.63556Z'  }
		* string json = event
    	
		* def result = RegistrationEventsPublisher.publish(json)

		* def expectedResponse =
		"""
			{
				rule_set_result: [
					{
						name: '#(RuleName.CustomerNotFoundInSerasa)',
						result: '#(RuleResult.Approved)'
					},
					{
						name: '#(RuleName.CustomerHasProblemsInSerasa)',
						result: '#(RuleResult.Approved)',
						pending: false,
						problems: [
							{
								code: '#(problemCodeBureauStatusNotRegular)'
							}
						]
					}
				]
			}
		"""

		Given url baseURLCompliance
		And path '/state/' + profileID
        And header Content-Type = 'application/json'
        And params ({ only_pending: false })		
		And retry until response.result == RuleResult.Approved
        When method GET
		Then assert responseStatus == 200
		And match response contains deep expectedResponse

		Given url mockURL
		And path "/subscribe/state-events/" + profileID
		And retry until response.length == 2
		When method GET
		Then assert responseStatus == 200
		And match response contains deep expectedCreatedStateEvent
		And match response contains deep expectedChangedStateEvent

    Scenario: Should APPROVE a person with REGULAR status in Bureau

		* def documentNumber =  CPFGenerator()
		* def profile = { profile_id: '#(profileID)', offer_type: '#(BureauRuleOffer)',  role_type: 'CUSTOMER', profile_type: 'INDIVIDUAL', document_number: '#(documentNumber)' }

		Given url mockURL
		And path '/v1/temis/profile/' + profileID
		And request profile
		When method POST
		Then assert responseStatus == 200

		* def enrichmentResponse = { situation: 1 }

		Given url mockURL
		And path '/temis-enrichment/individual/' + documentNumber
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
		Then assert ruleResults[0].name == RuleName.CustomerNotFoundInSerasa
		Then assert ruleResults[0].result == RuleResult.Approved
		Then assert ruleResults[1].name == RuleName.CustomerHasProblemsInSerasa
		Then assert ruleResults[1].result == RuleResult.Approved

		Given url mockURL
		And path "/subscribe/state-events/" + profileID
		And retry until response.length == 2
		When method GET
		Then assert responseStatus == 200
		And match response contains deep expectedCreatedStateEvent
		And match response contains deep expectedChangedStateEvent

	Scenario: Should APPROVE a person with PENDING status in Bureau
		
		* def customOfferType = "TEST_APPROVED_STATUSES_BUREAU_RULE" 
		* def customOffer = { offer_type : '#(customOfferType)', product: 'maquininha' }

		Given url baseURLCompliance
		And path '/offers'
		And header Content-Type = 'application/json'
		And request customOffer
		When method POST
		Then assert responseStatus == 201 || responseStatus == 409

		* def ruleSetConfig = { bureau: { approved_statuses: ['#(bureauStatusPendingRegularization)'] } }

		Given url mockURL
        And path '/temis-config/cadastral-validation-configs'
		And header Content-Type = 'application/json'
		And def catalog = call CreateSingleLevelCatalog { offer_type: '#(customOfferType)', role_type: 'CUSTOMER', person_type: 'INDIVIDUAL', account_flag: false, rules_config: #(ruleSetConfig)}
		And request catalog
		When method POST

		* def documentNumber =  CPFGenerator()
		* def profile = { profile_id: '#(profileID)', offer_type: '#(customOfferType)', role_type: 'CUSTOMER', profile_type: 'INDIVIDUAL', document_number: '#(documentNumber)' }

		Given url mockURL
		And path '/v1/temis/profile/' + profileID
		And request profile
		When method POST
		Then assert responseStatus == 200

		* def enrichmentResponse = { situation: 2 }

		Given url mockURL
		And path '/temis-enrichment/individual/' + documentNumber
		And request enrichmentResponse
		When method POST

		* def event = { profile_id: '#(profileID)', entity_id: '#(profileID)', entity_type: 'PROFILE', event_type: 'PROFILE_CHANGED', parent_type: 'PROFILE', update_date: '2023-08-03T19:29:06.63556Z'  }
		* string json = event
    	
		* def result = RegistrationEventsPublisher.publish(json)

		* def expectedResponse =
		"""
			{
				result: '#(RuleResult.Approved)',
				rule_set_result: [
					{
						name: '#(RuleName.CustomerNotFoundInSerasa)',										
						result: '#(RuleResult.Approved)',										
						pending: false
					},
					{
						name: '#(RuleName.CustomerHasProblemsInSerasa)',										
						result: '#(RuleResult.Approved)'										
					}
				]
			}		
		"""

		Given url baseURLCompliance
		And path '/state/' + profileID
        And header Content-Type = 'application/json'
        And params ({ only_pending: false })		
		And retry until response.result == RuleResult.Approved
        When method GET
		Then assert responseStatus == 200
		And match response contains deep expectedResponse

		Given url mockURL
		And path "/subscribe/state-events/" + profileID
		And retry until response.length == 2
		When method GET
		Then assert responseStatus == 200
		And match response contains deep expectedCreatedStateEvent
		And match response contains deep expectedChangedStateEvent

	Scenario: Should REJECT a person DECEASED in Bureau

		* def documentNumber =  CPFGenerator()
		* def profile = { profile_id: '#(profileID)', offer_type: '#(BureauRuleOffer)', role_type: 'CUSTOMER', profile_type: 'INDIVIDUAL', document_number: '#(documentNumber)' }

		Given url mockURL
		And path '/v1/temis/profile/' + profileID
		And request profile
		When method POST
		Then assert responseStatus == 200

		* def enrichmentResponse = { situation: 0}

		Given url mockURL
		And path '/temis-enrichment/individual/' + documentNumber
		And request enrichmentResponse
		When method POST

		* def event = { profile_id: '#(profileID)', entity_id: '#(profileID)', entity_type: 'PROFILE', event_type: 'PROFILE_CHANGED', parent_type: 'PROFILE', update_date: '2023-08-03T19:29:06.63556Z'  }
		* string json = event
    	
		* def result = RegistrationEventsPublisher.publish(json)

		* def expectedResponse =
		"""
			{
				result: '#(RuleResult.Rejected)',
				rule_set_result: [
					{
						name: '#(RuleName.CustomerNotFoundInSerasa)',										
						result: '#(RuleResult.Approved)'										
					},
					{
						name: '#(RuleName.CustomerHasProblemsInSerasa)',										
						result: '#(RuleResult.Rejected)',										
						problems: [
							{
								code: '#(problemCodeBureauStatusNotRegular)'											
							}
						]										
					}
				]		
			}
		"""

		Given url baseURLCompliance
		And path '/state/' + profileID
        And header Content-Type = 'application/json'
        And params ({ only_pending: false })		
		And retry until response.result == RuleResult.Rejected
        When method GET
		Then assert responseStatus == 200
		And match response contains deep expectedResponse

		Given url mockURL
		And path "/subscribe/state-events/" + profileID
		And retry until response.length == 2
		When method GET
		Then assert responseStatus == 200
		And match response contains deep expectedCreatedStateEvent
		And match response contains deep expectedChangedStateEvent