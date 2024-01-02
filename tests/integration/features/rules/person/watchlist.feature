Feature: Watchlist Rule

    Background:
        * url baseURLCompliance

	    * def WatchlistRuleOffer = "TEST_OFFER" + uuid() 
        * def TestRulePartner = "TestWatchlistRulePartner" + uuid()

		* def problemCodeDateOfBirthNoInputtedOrEnriched = "DATE_OF_BIRTH_NOT_INPUTTED_OR_ENRICHED"
		* def problemCodePersonFoundOnWatchlist = "PERSON_FOUND_ON_WATCHLIST"

		* def offer = { offer_type : '#(WatchlistRuleOffer)', product: 'maquininha'}

        Given path '/offers'
		And header Content-Type = 'application/json'
		And request offer
		When method POST

		* def ruleSetConfig = { watchlist: {  want_pep_tag: true, wanted_sources: ["OFAC_CAPTA","OFAC"] }}

        Given url mockURL
        And path '/temis-config/cadastral-validation-configs'
		And header Content-Type = 'application/json'
		And def catalogIndividual = call CreateSingleLevelCatalog { offer_type: '#(WatchlistRuleOffer)',  role_type: '#(RoleType.Customer)', person_type: '#(ProfileType.Individual)', account_flag: false, rules_config: #(ruleSetConfig)}
		And request catalogIndividual
		When method POST	
		
		Given url mockURL
        And path '/temis-config/cadastral-validation-configs'
		And header Content-Type = 'application/json'
		And def catalogCompany = call CreateSingleLevelCatalog { offer_type: '#(WatchlistRuleOffer)', role_type: '#(RoleType.Customer)', person_type: '#(ProfileType.Company)', account_flag: false, rules_config: #(ruleSetConfig)}
		And request catalogCompany
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

    Scenario: Reject by invalid date of birth

		* def documentNumber =  CPFGenerator()
		* def profile = { profile_id: '#(profileID)', offer_type: '#(WatchlistRuleOffer)', role_type: '#(RoleType.Customer)', profile_type: '#(ProfileType.Individual)', document_number: '#(documentNumber)' }
		* profile.individual = {}
		* profile.individual.date_of_birth = null
		
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
				result: '#(RuleResult.Analysing)',
				rule_set_result: [
					{
						set: '#(RuleSet.Watchlist)',
						name: '#(RuleName.Watchlist)',
						result: '#(RuleResult.Analysing)',
						metadata: 'Date of birth is not present in profile and was not enriched',
						problems: [
							{
								code: '#(problemCodeDateOfBirthNoInputtedOrEnriched)'
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
		And retry until response.result == RuleResult.Analysing
        When method GET
		And match response contains deep expectedResponse

		Given url mockURL
		And path "/subscribe/state-events/" + profileID
		And retry until response.length == 2
		When method GET
		Then assert responseStatus == 200
		And match response contains deep expectedCreatedStateEvent
		And match response contains deep expectedChangedStateEvent

	Scenario: Should find wanted sources for profile individual

		* def documentNumber =  CPFGenerator()
		* def profile = { profile_id: '#(profileID)', offer_type: '#(WatchlistRuleOffer)', role_type: '#(RoleType.Customer)', profile_type: '#(ProfileType.Individual)', document_number: '#(documentNumber)', name: 'ANA PAULA DA SILVA OLIVEIRA'}
		* profile.individual = {}
		* profile.individual.date_of_birth = '2021-03-29T10:38:20.744555-03:00'
		
		Given url mockURL
		And path '/v1/temis/profile/' + profileID
		And request profile
		When method POST
		Then assert responseStatus == 200

		* def watchlistResponse = {document_number: '#(documentNumber)', sources:["OFAC_CAPTA", "OFAC"]}

		Given url mockURL
		And path '/watchlist'
		And request watchlistResponse
		When method POST
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
		* def ruleResults = response.rule_set_result
		Then assert responseStatus == 200
		Then assert response.result == RuleResult.Analysing
		Then assert ruleResults[0].set == RuleSet.Watchlist
		Then assert ruleResults[0].name == RuleName.Watchlist
		Then assert ruleResults[0].result == RuleResult.Analysing
		Then assert ruleResults[0].problems[0].code == problemCodePersonFoundOnWatchlist
		Then match ruleResults[0].problems[0].detail contains ["OFAC_CAPTA", "OFAC"]		
		Then match ruleResults[0].metadata[0].sources contains ["OFAC_CAPTA", "OFAC"]
		Then match ruleResults[0].metadata[0] == 
		"""
			{
        		other: "##string",
        		entries: "#[] #string",
        		sources: "#[] #string",
        		watch: "##string",
        		link: "#string",
        		name: "#string",        
        		title: "#string"                        
    		}
		"""

		Given url mockURL
		And path "/subscribe/state-events/" + profileID
		And retry until response.length == 2
		When method GET
		Then assert responseStatus == 200
		And match response contains deep expectedCreatedStateEvent
		And match response contains deep expectedChangedStateEvent

	Scenario: Should find wanted sources for profile individual and reject directly when status parameter is provided

		* def customOfferType = "TEST_CUSTOM_RESULT_MATCH_WATCHLIST_FOUND" 
		* def customOffer = { offer_type : '#(customOfferType)', product: 'maquininha' }

		Given path '/offers'
		And header Content-Type = 'application/json'
		And request offer
		When method POST

		* def ruleSetConfig = { watchlist: {  want_pep_tag: true, wanted_sources: ["OFAC_CAPTA","OFAC"],  has_match_in_watch_list_status: '#(RuleResult.Rejected)'}}

		Given url mockURL
        And path '/temis-config/cadastral-validation-configs'
		And header Content-Type = 'application/json'
		And def catalogIndividual = call CreateSingleLevelCatalog { offer_type: '#(customOfferType)',  role_type: '#(RoleType.Customer)', person_type: '#(ProfileType.Individual)', account_flag: false, rules_config: #(ruleSetConfig)}
		And request catalogIndividual
		When method POST	

		* def documentNumber =  CPFGenerator()
		* def profile = { profile_id: '#(profileID)', offer_type: '#(customOfferType)', role_type: '#(RoleType.Customer)', profile_type: '#(ProfileType.Individual)', document_number: '#(documentNumber)', name: 'ANA PAULA DA SILVA OLIVEIRA'}
		* profile.individual = {}
		* profile.individual.date_of_birth = '2021-03-29T10:38:20.744555-03:00'
		
		Given url mockURL
		And path '/v1/temis/profile/' + profileID
		And request profile
		When method POST
		Then assert responseStatus == 200

		* def watchlistResponse = {document_number: '#(documentNumber)', sources:["Ofac_CAPTA", "Ofac"]}

		Given url mockURL
		And path '/watchlist'
		And request watchlistResponse
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
		Then assert ruleResults[0].set == RuleSet.Watchlist
		Then assert ruleResults[0].name == RuleName.Watchlist
		Then assert ruleResults[0].result == RuleResult.Rejected
		Then assert ruleResults[0].problems[0].code == problemCodePersonFoundOnWatchlist
		Then match ruleResults[0].problems[0].detail contains ["OFAC_CAPTA", "OFAC"]		
		Then match ruleResults[0].metadata[0].sources contains ["Ofac_CAPTA", "Ofac"]
		Then match ruleResults[0].metadata[0] == 
		"""
			{
        		other: "##string",
        		entries: "#[] #string",
        		sources: "#[] #string",
        		watch: "##string",
        		link: "#string",
        		name: "#string",        
        		title: "#string"                        
    		}
		"""	

		Given url mockURL
		And path "/subscribe/state-events/" + profileID
		And retry until response.length == 2
		When method GET
		Then assert responseStatus == 200
		And match response contains deep expectedCreatedStateEvent
		And match response contains deep expectedChangedStateEvent

	Scenario: Should not validate if company has not LegalName

		* def documentNumber =  CNPJGenerator()
		* def profile = { profile_id: '#(profileID)', offer_type: '#(WatchlistRuleOffer)', role_type: '#(RoleType.Customer)', profile_type: '#(ProfileType.Company)', document_number: '#(documentNumber)' }
		
		Given url mockURL
		And path '/v1/temis/profile/' + profileID
		And request profile
		When method POST
		Then assert responseStatus == 200

		* def event = { profile_id: '#(profileID)', entity_id: '#(profileID)', entity_type: 'PROFILE', event_type: 'PROFILE_CHANGED', parent_type: 'PROFILE', update_date: '2023-08-03T19:29:06.63556Z'  }
		* string json = event
    	* def result = RegistrationEventsPublisher.publish(json)
		* eval sleep(3000)

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

	Scenario: Should not approve profile company in watchlist

		* def documentNumber =  CNPJGenerator()
		* def companyName = 'WatchlistCompany' + uuid()
		* def profile = { profile_id: '#(profileID)', offer_type: '#(WatchlistRuleOffer)', role_type: '#(RoleType.Customer)', profile_type: '#(ProfileType.Company)', document_number: '#(documentNumber)', company: { legal_name: '#(companyName)'} }
		
		Given url mockURL
		And path '/v1/temis/profile/' + profileID
		And request profile
		When method POST
		Then assert responseStatus == 200

		* def watchlistResponse = {document_number: '#(documentNumber)', sources:["OFAC_CAPTA", "OFAC"]}
		* def watchlistResponses = []
		* def void = watchlistResponses.add(watchlistResponse)

		Given url mockURL
		And path '/watchlist/company/' + companyName
		And request watchlistResponses
		When method POST
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
		* def ruleResults = response.rule_set_result
		Then assert responseStatus == 200
		Then assert response.result == RuleResult.Analysing
		Then assert ruleResults[0].set == RuleSet.Watchlist
		Then assert ruleResults[0].name == RuleName.Watchlist
		Then assert ruleResults[0].result == RuleResult.Analysing
		Then assert ruleResults[0].problems[0].code == problemCodePersonFoundOnWatchlist
		Then match ruleResults[0].problems[0].detail contains ["OFAC_CAPTA", "OFAC"]		
		Then match ruleResults[0].metadata[0].sources contains ["OFAC_CAPTA", "OFAC"]

		Given url mockURL
		And path "/subscribe/state-events/" + profileID
		And retry until response.length == 2
		When method GET
		Then assert responseStatus == 200
		And match response contains deep expectedCreatedStateEvent
		And match response contains deep expectedChangedStateEvent
				
	Scenario: Should not approve profile company with legalname type "name/name"

		* def documentNumber =  CNPJGenerator()
		* def companyID = uuid()
		* def companyName = 'WatchlistCompany/' + companyID
		* def profile = { profile_id: '#(profileID)', offer_type: '#(WatchlistRuleOffer)', role_type: '#(RoleType.Customer)', profile_type: '#(ProfileType.Company)', document_number: '#(documentNumber)', company: { legal_name: '#(companyName)'} }

		Given url mockURL
		And path '/v1/temis/profile/' + profileID
		And request profile
		When method POST
		Then assert responseStatus == 200

		* def watchlistResponse = {document_number: '#(documentNumber)', sources:["OFAC_CAPTA", "OFAC"]}
		* def watchlistResponses = []
		* def void = watchlistResponses.add(watchlistResponse)

		Given url mockURL
		And path '/watchlist/company/' + 'WatchlistCompany' + companyID
		And request watchlistResponses
		When method POST
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
		* def ruleResults = response.rule_set_result
		Then assert responseStatus == 200
		Then assert response.result == RuleResult.Analysing
		Then assert ruleResults[0].set == RuleSet.Watchlist
		Then assert ruleResults[0].name == RuleName.Watchlist
		Then assert ruleResults[0].result == RuleResult.Analysing
		Then assert ruleResults[0].problems[0].code == problemCodePersonFoundOnWatchlist
		Then match ruleResults[0].problems[0].detail contains ["OFAC_CAPTA", "OFAC"]
		Then match ruleResults[0].metadata[0].sources contains ["OFAC_CAPTA", "OFAC"]

		Given url mockURL
		And path "/subscribe/state-events/" + profileID
		And retry until response.length == 2
		When method GET
		Then assert responseStatus == 200
		And match response contains deep expectedCreatedStateEvent
		And match response contains deep expectedChangedStateEvent

	Scenario: Should not find occurrences in watchlist for profile company

		* def documentNumber =  CNPJGenerator()
		* def companyID = uuid()
		* def companyName = 'WatchlistCompany/' + companyID
		* def profile = { profile_id: '#(profileID)', offer_type: '#(WatchlistRuleOffer)',  role_type: '#(RoleType.Customer)', profile_type: '#(ProfileType.Company)', document_number: '#(documentNumber)', company: { legal_name: '#(companyName)'} }
		
		Given url mockURL
		And path '/v1/temis/profile/' + profileID
		And request profile
		When method POST
		Then assert responseStatus == 200

		* def watchlistResponse = []

		Given url mockURL
		And path '/watchlist/company/' + 'WatchlistCompany' + companyID
		And request watchlistResponse
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
		And assert response.result == RuleResult.Approved

		Given url mockURL
		And path "/subscribe/state-events/" + profileID
		And retry until response.length == 2
		When method GET
		Then assert responseStatus == 200
		And match response contains deep expectedCreatedStateEvent
		And match response contains deep expectedChangedStateEvent
