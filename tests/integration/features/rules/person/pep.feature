Feature: Pep Rule

    Background:
        * url baseURLCompliance

	    * def PepRuleOffer = "TEST_OFFER" + uuid()  

		* def offer = { offer_type : '#(PepRuleOffer)', product: 'maquininha'}

        Given path '/offers'
		And header Content-Type = 'application/json'
		And request offer
		When method POST

		* def ruleSetConfig = { pep: {}, watchlist: {  want_pep_tag: true, wanted_sources: ["OFAC_CAPTA","OFAC"] }}    

		Given url mockURL
        And path '/temis-config/cadastral-validation-configs'
		And header Content-Type = 'application/json'
		And def catalog = call CreateSingleLevelCatalog { offer_type: '#(PepRuleOffer)',  role_type: '#(RoleType.Customer)', person_type: '#(ProfileType.Individual)', account_flag: false, rules_config: #(ruleSetConfig)}
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

    Scenario: Should be analysing by pep found in watchlist    

		* def documentNumber =  CPFGenerator()
		* def profile = { profile_id: '#(profileID)', offer_type: '#(PepRuleOffer)',  role_type: '#(RoleType.Customer)', profile_type: '#(ProfileType.Individual)', document_number: '#(documentNumber)' }
		* profile.individual = {}
		* profile.individual.date_of_birth = '2021-03-29T10:38:20.744555-03:00'
		
		Given url mockURL
		And path '/v1/temis/profile/' + profileID
		And request profile
		When method POST
		Then assert responseStatus == 200

		Given url mockURL
		And path '/watchlist'
		And def data = {document_number: '#(documentNumber)', sources:["OFAC_CAPTA", "OFAC", "PEP"]}
		And request data
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
		* def ruleResult = response.rule_set_result[1]
		Then assert ruleResult.set == RuleSet.Pep
		Then assert ruleResult.name == RuleName.Pep
		Then assert ruleResult.result == RuleResult.Analysing
		Then assert ruleResult.pending == true
		And  assert ruleResult.metadata.pep_sources[0] == "WATCHLIST" 

		Given url mockURL
		And path "/subscribe/state-events/" + profileID
		And retry until response.length == 2
		When method GET
		Then assert responseStatus == 200
		And match response contains deep expectedCreatedStateEvent
		And match response contains deep expectedChangedStateEvent

	Scenario: Should be analysing by self declared pep

		* def documentNumber =  CPFGenerator()
		* def profile = { profile_id: '#(profileID)', offer_type: '#(PepRuleOffer)', role_type: '#(RoleType.Customer)', profile_type: '#(ProfileType.Individual)', document_number: '#(documentNumber)' }
		* profile.individual = {}
		* profile.individual.date_of_birth = '2021-03-29T10:38:20.744555-03:00'
		* profile.individual.pep = true

		Given url mockURL
		And path '/v1/temis/profile/' + profileID
		And request profile
		When method POST
		Then assert responseStatus == 200

		Given url mockURL
		And path '/watchlist'
		And def data = {document_number: '#(documentNumber)', sources:["OFAC_CAPTA", "OFAC"]}
		And request data
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
		* def ruleResult = response.rule_set_result[1]
		Then assert ruleResult.set == RuleSet.Pep
		Then assert ruleResult.name == RuleName.Pep
		Then assert ruleResult.result == RuleResult.Analysing
		Then assert ruleResult.pending == true
		And  assert ruleResult.metadata.pep_sources[0] == "SELF_DECLARED" 

		Given url mockURL
		And path "/subscribe/state-events/" + profileID
		And retry until response.length == 2
		When method GET
		Then assert responseStatus == 200
		And match response contains deep expectedCreatedStateEvent
		And match response contains deep expectedChangedStateEvent

	Scenario: Should be analysing when is on PEP COAF list

		* def documentNumber =  CPFGenerator()
		* def profile = { profile_id: '#(profileID)', offer_type: '#(PepRuleOffer)', role_type: '#(RoleType.Customer)', profile_type: '#(ProfileType.Individual)', document_number: '#(documentNumber)' }
		* profile.individual = {}
		* profile.individual.date_of_birth = '2021-03-29T10:38:20.744555-03:00'
		* profile.individual.pep = false

		Given url mockURL
		And path '/v1/temis/profile/' + profileID
		And request profile
		When method POST
		Then assert responseStatus == 200

		* def pepCOAF = 
		"""
		{
			"document_number": "#(documentNumber)",
			"name": "Some name",
			"role": "Senador",
			"institution": "Congreso",
			"start_date":"2000-01-01",
			"end_date":"2050-01-01"
		}	
		"""

		Given url mockURL
		And path "/temis-restrictive-lists/pep/" + documentNumber
		And request pepCOAF
		When method POST
		Then assert responseStatus == 201

		Given url mockURL
		And path "/temis-restrictive-lists/pep/" + documentNumber
		When method GET
		Then assert responseStatus == 200
	
		* def event = { profile_id: '#(profileID)', entity_id: '#(profileID)', entity_type: 'PROFILE', event_type: 'PROFILE_CHANGED', parent_type: 'PROFILE', update_date: '2023-08-03T19:29:06.63556Z'  }
		* string json = event
    	
		* def result = RegistrationEventsPublisher.publish(json)

		Given url baseURLCompliance
		And path '/state/' + profileID
        And header Content-Type = 'application/json'
        And params ({ only_pending: false })		
		And retry until response.result == RuleResult.Analysing
        When method GET
		* def ruleResult = response.rule_set_result[1]
		Then assert ruleResult.set == RuleSet.Pep
		And  assert ruleResult.name == RuleName.Pep
		And  assert ruleResult.result == RuleResult.Analysing
		And  assert ruleResult.pending == true
		And  assert ruleResult.metadata.pep_sources[0] == "COAF" 
		And  assert ruleResult.metadata.pep_information.document_number == documentNumber 

		Given url mockURL
		And path "/subscribe/state-events/" + profileID
		And retry until response.length == 2
		When method GET
		Then assert responseStatus == 200
		And match response contains deep expectedCreatedStateEvent
		And match response contains deep expectedChangedStateEvent

	Scenario: Should approve pep rule

		* def documentNumber =  CPFGenerator()
		* def profile = { profile_id: '#(profileID)', offer_type: '#(PepRuleOffer)',  role_type: '#(RoleType.Customer)', profile_type: '#(ProfileType.Individual)', document_number: '#(documentNumber)' }
		* profile.individual = {}
		* profile.individual.date_of_birth = '2021-03-29T10:38:20.744555-03:00'
		
		Given url mockURL
		And path '/v1/temis/profile/' + profileID
		And request profile
		When method POST
		Then assert responseStatus == 200

		Given url mockURL
		And path '/watchlist'
		And def data = {document_number: '#(documentNumber)', sources:["OFAC_CAPTA", "OFAC"]}
		And request data
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
		* def ruleResult = response.rule_set_result[1]
		Then assert ruleResult.set == RuleSet.Pep
		Then assert ruleResult.name == RuleName.Pep
		Then assert ruleResult.result == RuleResult.Approved
		Then assert ruleResult.pending == false

		Given url mockURL
		And path "/subscribe/state-events/" + profileID
		And retry until response.length == 2
		When method GET
		Then assert responseStatus == 200
		And match response contains deep expectedCreatedStateEvent
		And match response contains deep expectedChangedStateEvent