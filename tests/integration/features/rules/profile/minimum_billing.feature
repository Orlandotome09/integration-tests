Feature: Test Minimum Billing Rule

	Background:
		* url baseURLCompliance

		* def MinimumBillingRuleOfferDefault = "TEST_OFFER" + uuid()  
		* def MinimumBillingRuleOfferCustom = "TEST_OFFER" + uuid()  

		* def problemCodeCompanyHasInsufficientBilling = "COMPANY_HAS_INSUFFICIENT_BILLING"

		* def offerDefault = { offer_type : '#(MinimumBillingRuleOfferDefault)', product : 'maquininha'}
		* def offerCustom = { offer_type : '#(MinimumBillingRuleOfferCustom)', product : 'maquininha'}

		Given path '/offers'
		And header Content-Type = 'application/json'
		And request offerDefault
		When method POST

		Given path '/offers'
		And header Content-Type = 'application/json'
		And request offerCustom
		When method POST

		* def ruleSetConfigDefault = { minimum_billing: {} }
		* def ruleSetConfigCustom = { minimum_billing: {minimum_billing: 3000.00} }

		Given url mockURL
        And path '/temis-config/cadastral-validation-configs'
		And header Content-Type = 'application/json'
		And def catalog = call CreateSingleLevelCatalog { offer_type: '#(MinimumBillingRuleOfferDefault)', role_type: '#(RoleType.Customer)', person_type: '#(ProfileType.Company)', account_flag: false, rules_config: #(ruleSetConfigDefault)}
		And request catalog
		When method POST

		Given url mockURL
        And path '/temis-config/cadastral-validation-configs'
		And header Content-Type = 'application/json'
		And def catalog = call CreateSingleLevelCatalog { offer_type: '#(MinimumBillingRuleOfferCustom)', role_type: '#(RoleType.Customer)', person_type: '#(ProfileType.Company)', account_flag: false, rules_config: #(ruleSetConfigCustom)}
		And request catalog
		When method POST

		* def profileID = uuid()
		* def documentNumber = DocumentNormalizer(CNPJGenerator())
		* def profile =
		"""
		{
			"profile_id": "#(profileID)",
			"document_number": "#(documentNumber)",
			"offer_type": "#(MinimumBillingRuleOfferDefault)",
			"role_type":"#(RoleType.Customer)",
			"profile_type": "#(ProfileType.Company)",
			"company":{
				"annual_income": 0,
				"annual_income_currency": "BRL"
			}
		}    
		"""

		* def expectedProfileCreatedStateEvent = 
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

		* def expectedProfileChangedStateEvent = 
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
		
	Scenario: Should approve profile company with minimum billing above the required
		* profile.offer_type = MinimumBillingRuleOfferDefault
		* profile.company.annual_income = 1001.00

		Given url mockURL
		And path '/v1/temis/profile/' + profileID
		And request profile
		When method POST
		Then assert responseStatus == 200

		* def event = { profile_id: '#(profileID)', entity_id: '#(profileID)', entity_type: 'PROFILE', event_type: 'PROFILE_CHANGED', parent_type: 'PROFILE', update_date: '2023-08-03T19:29:06.63556Z'  }
		* string json = event
		
		* def result = RegistrationEventsPublisher.publish(json)

		* def expectedMinimumBillingResult = 
		"""
		{
			"step_number" : 0 ,
			"set": "#(RuleSet.MinimumBilling)",
			"name": "#(RuleName.InsufficientBilling)",
			"result": "#(RuleResult.Approved)",
			"metadata": null,
			"pending": false,
		}	
		"""

		Given url baseURLCompliance
		And path '/state/' + profileID
		And header Content-Type = 'application/json'
		And params ({ only_pending: false })
		And retry until response.result ==  RuleResult.Approved
		When method GET
		And match response.rule_set_result contains deep expectedMinimumBillingResult

		Given url mockURL
		And path "/subscribe/state-events/" + profileID
		And retry until response.length == 2
		When method GET
		Then assert responseStatus == 200
		And match response contains deep expectedProfileCreatedStateEvent
		And match response contains deep expectedProfileChangedStateEvent

	Scenario: Should approve profile company with minimum billing equal the required
		* profile.offer_type = MinimumBillingRuleOfferDefault
		* profile.company.annual_income = 1000

		Given url mockURL
		And path '/v1/temis/profile/' + profileID
		And request profile
		When method POST
		Then assert responseStatus == 200

		* def event = { profile_id: '#(profileID)', entity_id: '#(profileID)', entity_type: 'PROFILE', event_type: 'PROFILE_CHANGED', parent_type: 'PROFILE', update_date: '2023-08-03T19:29:06.63556Z'  }
		* string json = event
		
		* def result = RegistrationEventsPublisher.publish(json)

		* def expectedMinimumBillingResult = 
		"""
		{
			"step_number" : 0 ,
			"set": "#(RuleSet.MinimumBilling)",
			"name": "#(RuleName.InsufficientBilling)",
			"result": "#(RuleResult.Approved)",
			"metadata": null,
			"pending": false,
		}	
		"""

		Given url baseURLCompliance
		And path '/state/' + profileID
		And header Content-Type = 'application/json'
		And params ({ only_pending: false })
		And retry until response.result ==  RuleResult.Approved
		When method GET
		And match response.rule_set_result contains deep expectedMinimumBillingResult

		Given url mockURL
		And path "/subscribe/state-events/" + profileID
		And retry until response.length == 2
		When method GET
		Then assert responseStatus == 200
		And match response contains deep expectedProfileCreatedStateEvent
		And match response contains deep expectedProfileChangedStateEvent
	
	Scenario: Should analysing profile company with minimum billing below the required
		* profile.offer_type = MinimumBillingRuleOfferDefault
		* profile.company.annual_income = 999.00

		Given url mockURL
		And path '/v1/temis/profile/' + profileID
		And request profile
		When method POST
		Then assert responseStatus == 200

		* def event = { profile_id: '#(profileID)', entity_id: '#(profileID)', entity_type: 'PROFILE', event_type: 'PROFILE_CHANGED', parent_type: 'PROFILE', update_date: '2023-08-03T19:29:06.63556Z'  }
		* string json = event
		
		* def result = RegistrationEventsPublisher.publish(json)

		* def expectedMinimumBillingResult = 
		"""
		{
			"step_number" : 0 ,
			"set": "#(RuleSet.MinimumBilling)",
			"name": "#(RuleName.InsufficientBilling)",
			"result": "#(RuleResult.Analysing)",
			"pending": true,
			"metadata": "Company billing is below 1000",
			"problems": [
				{
					"code": "COMPANY_HAS_INSUFFICIENT_BILLING",
					"detail": {
						"minimum_billing_required": 1000,
						"company_billing": 999
					}
				}
			]
		}
		"""
		
		Given url baseURLCompliance
		And path '/state/' + profileID
		And header Content-Type = 'application/json'
		And params ({ only_pending: false })
		And retry until response.result ==  RuleResult.Analysing
		When method GET
		And match response.rule_set_result contains deep expectedMinimumBillingResult

		Given url mockURL
		And path "/subscribe/state-events/" + profileID
		And retry until response.length == 2
		When method GET
		Then assert responseStatus == 200
		And match response contains deep expectedProfileCreatedStateEvent
		And match response contains deep expectedProfileChangedStateEvent

	Scenario: Should approve profile company with minimum billing above the custom required
		* profile.offer_type = MinimumBillingRuleOfferCustom
		* profile.company.annual_income = 3001.00

		Given url mockURL
		And path '/v1/temis/profile/' + profileID
		And request profile
		When method POST
		Then assert responseStatus == 200

		* def event = { profile_id: '#(profileID)', entity_id: '#(profileID)', entity_type: 'PROFILE', event_type: 'PROFILE_CHANGED', parent_type: 'PROFILE', update_date: '2023-08-03T19:29:06.63556Z'  }
		* string json = event
		
		* def result = RegistrationEventsPublisher.publish(json)

		* def expectedMinimumBillingResult = 
		"""
		{
			"step_number" : 0 ,
			"set": "#(RuleSet.MinimumBilling)",
			"name": "#(RuleName.InsufficientBilling)",
			"result": "#(RuleResult.Approved)",
			"metadata": null,
			"pending": false,
		}	
		"""

		Given url baseURLCompliance
		And path '/state/' + profileID
		And header Content-Type = 'application/json'
		And params ({ only_pending: false })
		And retry until response.result ==  RuleResult.Approved
		When method GET
		And match response.rule_set_result contains deep expectedMinimumBillingResult

		Given url mockURL
		And path "/subscribe/state-events/" + profileID
		And retry until response.length == 2
		When method GET
		Then assert responseStatus == 200
		And match response contains deep expectedProfileCreatedStateEvent
		And match response contains deep expectedProfileChangedStateEvent


	Scenario: Should analysing profile company with minimum billing below the custom required
		* profile.offer_type = MinimumBillingRuleOfferCustom
		* profile.company.annual_income = 2999.00

		Given url mockURL
		And path '/v1/temis/profile/' + profileID
		And request profile
		When method POST
		Then assert responseStatus == 200

		* def event = { profile_id: '#(profileID)', entity_id: '#(profileID)', entity_type: 'PROFILE', event_type: 'PROFILE_CHANGED', parent_type: 'PROFILE', update_date: '2023-08-03T19:29:06.63556Z'  }
		* string json = event
		
		* def result = RegistrationEventsPublisher.publish(json)

		* def expectedMinimumBillingResult = 
		"""
		{
			"step_number" : 0 ,
			"set": "#(RuleSet.MinimumBilling)",
			"name": "#(RuleName.InsufficientBilling)",
			"result": "#(RuleResult.Analysing)",
			"pending": true,
			"metadata": "Company billing is below 3000",
			"problems": [
				{
					"code": "COMPANY_HAS_INSUFFICIENT_BILLING",
					"detail": {
						"minimum_billing_required": 3000,
						"company_billing": 2999
					}
				}
			]
		}
		"""
		
		Given url baseURLCompliance
		And path '/state/' + profileID
		And header Content-Type = 'application/json'
		And params ({ only_pending: false })
		And retry until response.result ==  RuleResult.Analysing
		When method GET
		And match response.rule_set_result contains deep expectedMinimumBillingResult

		Given url mockURL
		And path "/subscribe/state-events/" + profileID
		And retry until response.length == 2
		When method GET
		Then assert responseStatus == 200
		And match response contains deep expectedProfileCreatedStateEvent
		And match response contains deep expectedProfileChangedStateEvent