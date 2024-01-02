Feature: Blacklist Rule

	Background:
		* url baseURLCompliance

		* def BlacklistRuleOffer = "TEST_OFFER" + uuid() 
		* def TestRulePartner = "TestBlacklistRulePartner"
		
		* def offer = { offer_type : '#(BlacklistRuleOffer)', product: 'maquininha'}

		Given path '/offers'
		And header Content-Type = 'application/json'
		And request offer
		When method POST

		* def ruleSetConfig = { blacklist: {} }

		Given url mockURL
        And path '/temis-config/cadastral-validation-configs'
		And header Content-Type = 'application/json'
		And def catalogIndividual = call CreateSingleLevelCatalog { offer_type: '#(BlacklistRuleOffer)',  role_type: '#(RoleType.Customer)', person_type: '#(ProfileType.Individual)', account_flag: false, rules_config: #(ruleSetConfig) }
		And request catalogIndividual
		When method POST

		Given url mockURL
        And path '/temis-config/cadastral-validation-configs'
		And header Content-Type = 'application/json'
		And def catalogCompany = call CreateSingleLevelCatalog { offer_type: '#(BlacklistRuleOffer)', role_type: '#(RoleType.Customer)', person_type: '#(ProfileType.Company)', account_flag: false, rules_config: #(ruleSetConfig)}
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

	Scenario: Profile individual present in a Blacklist

		* def documentNumber =  CPFGenerator()
		* def profile = { profile_id: '#(profileID)', offer_type: '#(BlacklistRuleOffer)', role_type: '#(RoleType.Customer)', profile_type: '#(ProfileType.Individual)', document_number: '#(documentNumber)' }
		* profile.individual = {}
		* profile.individual.date_of_birth = null

        * def internalList = [{justification: 'Profile present on the black list'}]

		Given url mockURL
		And path '/temis-restrictive-lists/internal-list'
		And param document_number = documentNumber 
		And request internalList
		When method POST
		Then assert responseStatus == 200        

		Given url mockURL
		And path '/v1/temis/profile/' + profileID
		And request profile
		When method POST
		Then assert responseStatus == 200

		* def event = { profile_id: '#(profileID)', entity_id: '#(profileID)', update_date: '2023-08-03T19:29:06.63556Z'  }
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
		Then assert ruleResults[0].set == RuleSet.Blacklist
		Then assert ruleResults[0].name == RuleName.Blacklist
		Then assert ruleResults[0].result == RuleResult.Analysing

		Given url mockURL
		And path "/subscribe/state-events/" + profileID
		And retry until response.length == 2
		When method GET
		Then assert responseStatus == 200
		And match response contains deep expectedCreatedStateEvent
		And match response contains deep expectedChangedStateEvent

	Scenario: Profile individual not present in a Blacklist

		* def documentNumber =  CPFGenerator()
		* def profile = { profile_id: '#(profileID)', offer_type: '#(BlacklistRuleOffer)', role_type: '#(RoleType.Customer)', profile_type: '#(ProfileType.Individual)', document_number: '#(documentNumber)' }
		* profile.individual = {}
		* profile.individual.date_of_birth = null
		* def validationBody = { justifications: [], status: '', partnerId: '123', documentNumber: '#(documentNumber)' }

        * def internalList = []

		Given url mockURL
		And path '/temis-restrictive-lists/internal-list'
		And param document_number = documentNumber 
		And request internalList
		When method POST
		Then assert responseStatus == 200      

		Given url mockURL
		And path '/v1/temis/profile/' + profileID
		And request profile
		When method POST
		Then assert responseStatus == 200

		* def event = { profile_id: '#(profileID)', entity_id: '#(profileID)', update_date: '2023-08-03T19:29:06.63556Z'  }
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

		Given url mockURL
		And path "/subscribe/state-events/" + profileID
		And retry until response.length == 2
		When method GET
		Then assert responseStatus == 200
		And match response contains deep expectedCreatedStateEvent
		And match response contains deep expectedChangedStateEvent

	Scenario: Profile company present in a Blacklist

		* def documentNumber =  CNPJGenerator()
		* def documentNormalized = DocumentNormalizer(documentNumber)
		* def profile = { profile_id: '#(profileID)', offer_type: '#(BlacklistRuleOffer)', role_type: '#(RoleType.Customer)', profile_type: '#(ProfileType.Company)', document_number: '#(documentNormalized)' }

		* def internalList = [{justification: 'Profile present on the black list'}]

		Given url mockURL
		And path '/temis-restrictive-lists/internal-list'
		And param document_number = documentNormalized 
		And request internalList
		When method POST
		Then assert responseStatus == 200     

		Given url mockURL
		And path '/v1/temis/profile/' + profileID
		And request profile
		When method POST
		Then assert responseStatus == 200

		* def event = { profile_id: '#(profileID)', entity_id: '#(profileID)', update_date: '2023-08-03T19:29:06.63556Z'  }
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
		Then assert ruleResults[0].set == RuleSet.Blacklist
		Then assert ruleResults[0].name == RuleName.Blacklist
		Then assert ruleResults[0].result == RuleResult.Analysing
		
		Given url mockURL
		And path "/subscribe/state-events/" + profileID
		And retry until response.length == 2
		When method GET
		Then assert responseStatus == 200
		And match response contains deep expectedCreatedStateEvent
		And match response contains deep expectedChangedStateEvent

	Scenario: Profile company not present in a Blacklist

		* def documentNumber =  CNPJGenerator()
		* def documentNormalized = DocumentNormalizer(documentNumber)
		* def profile = { profile_id: '#(profileID)', offer_type: '#(BlacklistRuleOffer)', role_type: '#(RoleType.Customer)', profile_type: '#(ProfileType.Company)', document_number: '#(documentNormalized)' }

		* def internalList = []

		Given url mockURL
		And path '/temis-restrictive-lists/internal-list'
		And param document_number = documentNormalized 
		And request internalList
		When method POST
		Then assert responseStatus == 200    

		Given url mockURL
		And path '/v1/temis/profile/' + profileID
		And request profile
		When method POST
		Then assert responseStatus == 200

		* def event = { profile_id: '#(profileID)', entity_id: '#(profileID)', update_date: '2023-08-03T19:29:06.63556Z'  }
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

		Given url mockURL
		And path "/subscribe/state-events/" + profileID
		And retry until response.length == 2
		When method GET
		Then assert responseStatus == 200
		And match response contains deep expectedCreatedStateEvent
		And match response contains deep expectedChangedStateEvent