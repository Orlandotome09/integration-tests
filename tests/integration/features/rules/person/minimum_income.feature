Feature: Test Minimum Income Rule

Background:
	* url baseURLCompliance

	* def MinimumIncomeRuleOfferDefault = "TEST_OFFER" + uuid()  
	* def MinimumIncomeRuleOfferCustom = "TEST_OFFER" + uuid()  

	* def problemCodePersonHasInsufficientMinimumIncome = "PERSON_HAS_INSUFFICIENT_INCOME"

	* def offerDefault = { offer_type : '#(MinimumIncomeRuleOfferDefault)', product : 'maquininha'}
	* def offerCustom = { offer_type : '#(MinimumIncomeRuleOfferCustom)', product : 'maquininha'}

	Given path '/offers'
	And header Content-Type = 'application/json'
	And request offerDefault
	When method POST

	Given path '/offers'
	And header Content-Type = 'application/json'
	And request offerCustom
	When method POST

	* def ruleSetConfigDefault = { minimum_income: {} }
	* def ruleSetConfigCustom = { minimum_income: {minimum_income: 3000.00} }

	Given url mockURL
	And path '/temis-config/cadastral-validation-configs'
	And header Content-Type = 'application/json'
	And def catalog = call CreateSingleLevelCatalog { offer_type: '#(MinimumIncomeRuleOfferDefault)', role_type: '#(RoleType.Customer)', person_type: '#(ProfileType.Individual)', account_flag: false, rules_config: #(ruleSetConfigDefault)}
	And request catalog
	When method POST

	Given url mockURL
	And path '/temis-config/cadastral-validation-configs'
	And header Content-Type = 'application/json'
	And def catalog = call CreateSingleLevelCatalog { offer_type: '#(MinimumIncomeRuleOfferCustom)', role_type: '#(RoleType.Customer)', person_type: '#(ProfileType.Individual)', account_flag: false, rules_config: #(ruleSetConfigCustom)}
	And request catalog
	When method POST

	* def profileID = uuid()
	* def documentNumber = DocumentNormalizer(CPFGenerator())
	* def profile =
	"""
	{
		"profile_id": "#(profileID)",
		"document_number": "#(documentNumber)",
		"offer_type": "#(MinimumIncomeRuleOfferDefault)",
		"role_type":"#(RoleType.Customer)",
		"profile_type": "#(ProfileType.Individual)",
		"individual":{
			"income": 0
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
	
Scenario: Should approve profile individual with minimum income above the required
	* profile.offer_type = MinimumIncomeRuleOfferDefault
	* profile.individual.income = 1001.00

	Given url mockURL
	And path '/v1/temis/profile/' + profileID
	And request profile
	When method POST
	Then assert responseStatus == 200

	* def event = { profile_id: '#(profileID)', entity_id: '#(profileID)', entity_type: 'PROFILE', event_type: 'PROFILE_CHANGED', parent_type: 'PROFILE', update_date: '2023-08-03T19:29:06.63556Z'  }
	* string json = event
	
	* def result = RegistrationEventsPublisher.publish(json)

	* def expectedMinimumIncomeResult = 
	"""
	{
		"step_number" : 0 ,
		"set": "#(RuleSet.MinimumIncome)",
		"name": "#(RuleName.InsufficientIncome)",
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
	And match response.rule_set_result contains deep expectedMinimumIncomeResult

	Given url mockURL
	And path "/subscribe/state-events/" + profileID
	And retry until response.length == 2
	When method GET
	Then assert responseStatus == 200
	And match response contains deep expectedProfileCreatedStateEvent
	And match response contains deep expectedProfileChangedStateEvent

Scenario: Should approve profile individual with minimum income equal the required
	* profile.offer_type = MinimumIncomeRuleOfferDefault
	* profile.individual.income = 1000.00

	Given url mockURL
	And path '/v1/temis/profile/' + profileID
	And request profile
	When method POST
	Then assert responseStatus == 200

	* def event = { profile_id: '#(profileID)', entity_id: '#(profileID)', entity_type: 'PROFILE', event_type: 'PROFILE_CHANGED', parent_type: 'PROFILE', update_date: '2023-08-03T19:29:06.63556Z'  }
	* string json = event
	
	* def result = RegistrationEventsPublisher.publish(json)

	* def expectedMinimumIncomeResult = 
	"""
	{
		"step_number" : 0 ,
		"set": "#(RuleSet.MinimumIncome)",
		"name": "#(RuleName.InsufficientIncome)",
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
	And match response.rule_set_result contains deep expectedMinimumIncomeResult

	Given url mockURL
	And path "/subscribe/state-events/" + profileID
	And retry until response.length == 2
	When method GET
	Then assert responseStatus == 200
	And match response contains deep expectedProfileCreatedStateEvent
	And match response contains deep expectedProfileChangedStateEvent

Scenario: Should analysing profile individual with minimum income below the required
	* profile.offer_type = MinimumIncomeRuleOfferDefault
	* profile.individual.income = 999.00

	Given url mockURL
	And path '/v1/temis/profile/' + profileID
	And request profile
	When method POST
	Then assert responseStatus == 200

	* def event = { profile_id: '#(profileID)', entity_id: '#(profileID)', entity_type: 'PROFILE', event_type: 'PROFILE_CHANGED', parent_type: 'PROFILE', update_date: '2023-08-03T19:29:06.63556Z'  }
	* string json = event
	
	* def result = RegistrationEventsPublisher.publish(json)

	* def expectedMinimumIncomeResult = 
	"""
	{
		"step_number" : 0,
		"set": "#(RuleSet.MinimumIncome)",
		"name": "#(RuleName.InsufficientIncome)",
		"result": "#(RuleResult.Analysing)",
		"pending": true,
		"metadata": "Profile income is below 1000",
		"problems": [
			{
				"code": "PERSON_HAS_INSUFFICIENT_INCOME",
				"detail": {
					"minimum_income_required": 1000,
					"person_income": 999
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
	And match response.rule_set_result contains deep expectedMinimumIncomeResult

	Given url mockURL
	And path "/subscribe/state-events/" + profileID
	And retry until response.length == 2
	When method GET
	Then assert responseStatus == 200
	And match response contains deep expectedProfileCreatedStateEvent
	And match response contains deep expectedProfileChangedStateEvent

Scenario: Should approve profile individual with minimum income above the custom required
	* profile.offer_type = MinimumIncomeRuleOfferCustom
	* profile.individual.income = 3001.00

	Given url mockURL
	And path '/v1/temis/profile/' + profileID
	And request profile
	When method POST
	Then assert responseStatus == 200

	* def event = { profile_id: '#(profileID)', entity_id: '#(profileID)', entity_type: 'PROFILE', event_type: 'PROFILE_CHANGED', parent_type: 'PROFILE', update_date: '2023-08-03T19:29:06.63556Z'  }
	* string json = event
	
	* def result = RegistrationEventsPublisher.publish(json)

	* def expectedMinimumIncomeResult = 
	"""
	{
		"step_number" : 0 ,
		"set": "#(RuleSet.MinimumIncome)",
		"name": "#(RuleName.InsufficientIncome)",
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
	And match response.rule_set_result contains deep expectedMinimumIncomeResult

	Given url mockURL
	And path "/subscribe/state-events/" + profileID
	And retry until response.length == 2
	When method GET
	Then assert responseStatus == 200
	And match response contains deep expectedProfileCreatedStateEvent
	And match response contains deep expectedProfileChangedStateEvent


Scenario: Should analysing profile individual with minimum income below the custom required
	* profile.offer_type = MinimumIncomeRuleOfferCustom
	* profile.individual.income = 2999.00

	Given url mockURL
	And path '/v1/temis/profile/' + profileID
	And request profile
	When method POST
	Then assert responseStatus == 200

	* def event = { profile_id: '#(profileID)', entity_id: '#(profileID)', entity_type: 'PROFILE', event_type: 'PROFILE_CHANGED', parent_type: 'PROFILE', update_date: '2023-08-03T19:29:06.63556Z'  }
	* string json = event
	
	* def result = RegistrationEventsPublisher.publish(json)

	* def expectedMinimumIncomeResult = 
	"""
	{
		"step_number" : 0 ,
		"set": "#(RuleSet.MinimumIncome)",
		"name": "#(RuleName.InsufficientIncome)",
		"result": "#(RuleResult.Analysing)",
		"pending": true,
		"metadata": "Profile income is below 3000",
		"problems": [
			{
				"code": "PERSON_HAS_INSUFFICIENT_INCOME",
				"detail": {
					"minimum_income_required": 3000,
					"person_income": 2999
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
	And match response.rule_set_result contains deep expectedMinimumIncomeResult

	Given url mockURL
	And path "/subscribe/state-events/" + profileID
	And retry until response.length == 2
	When method GET
	Then assert responseStatus == 200
	And match response contains deep expectedProfileCreatedStateEvent
	And match response contains deep expectedProfileChangedStateEvent