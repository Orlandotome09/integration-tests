Feature: Contract Rule Incomplete

    Background:
        * url baseURLCompliance

	    * def PreconditionContractRuleOffer = "TEST_OFFER" + uuid()
		* def PreconditionContractRuleOfferNotApproved = "TEST_OFFER" + uuid() 

		* def ruleSetPrecondition = "PRECONDITION_CONTRACT_SET"
		* def ruleNameProfileApproved = "PROFILE_APPROVED"

        * def problemCodeStateNotFound = "STATE_NOT_FOUND"
        * def problemCodeProfileNotApproved = "PROFILE_NOT_APPROVED"
	
        * def offer = { offer_type : '#(PreconditionContractRuleOffer)', product: 'maquininha'}

        Given path '/offers'
		And header Content-Type = 'application/json'
		And request offer
		When method POST

        * def ruleSetConfig = { manual_block: {} }

        Given url mockURL
        And path '/temis-config/cadastral-validation-configs'
		And header Content-Type = 'application/json'
		And def catalog = call CreateSingleLevelCatalog { offer_type: '#(PreconditionContractRuleOffer)', role_type: '#(RoleType.Customer)', person_type: '#(ProfileType.Individual)', account_flag: false, rrules_config: #(ruleSetConfig)}
		And request catalog
		When method POST

       
        * def offerNotApproved = { offer_type : '#(PreconditionContractRuleOfferNotApproved)' , product: 'maquininha'}

        Given path '/offers'
		And header Content-Type = 'application/json'
		And request offerNotApproved
		When method POST

        * def ruleSetConfig = { bureau: {} }

        Given url mockURL
        And path '/temis-config/cadastral-validation-configs'
		And header Content-Type = 'application/json'
		And def catalog = call CreateSingleLevelCatalog { offer_type: '#(PreconditionContractRuleOfferNotApproved)', role_type: '#(RoleType.Customer)', person_type: '#(ProfileType.Individual)', account_flag: false, rules_config: #(ruleSetConfig) }
		And request catalog
		When method POST    

	Scenario: Should not find profile state
        		
		* def contractID = uuid()
		* def profileID = uuid()
		* def contract = { contract_id: '#(contractID)', profile_id: '#(profileID)', document_id: ''}

		Given url mockURL
		And path '/v1/temis/contract/' + contractID
		And request contract
		When method POST
    
		* def event = { contract_id: '#(contractID)', entity_id: '#(contractID)', entity_type: 'CONTRACT', event_type: 'CONTRACT_CREATED', parent_type: 'CONTRACT', update_date: '2023-08-03T19:29:06.63556Z'  }
		* string json = event
    	
		* def result = RegistrationEventsPublisher.publish(json)

		* def expectedResponse =
		"""
			{
			  result: '#(RuleResult.Approved)'
			}
		"""

		Given url baseURLCompliance
		And path '/state/' + contractID
        And header Content-Type = 'application/json'
        And params ({ only_pending: false })
		And retry until response.result == RuleResult.Approved 		
        When method GET
		Then assert responseStatus == 200
		Then match response contains deep expectedResponse

    Scenario: Should find profile not approved
        
        * def profileID = uuid()
		* def documentNumber = DocumentNormalizer(CPFGenerator())
		* def profile = { profile_id: '#(profileID)', offer_type: '#(PreconditionContractRuleOfferNotApproved)', role_type: '#(RoleType.Customer)', profile_type: '#(ProfileType.Individual)', document_number: '#(documentNumber)' }

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
        
        * def contractID = uuid()
		* def contract = { contract_id: '#(contractID)', profile_id: '#(profileID)', document_id: ''}

		Given url mockURL
		And path '/v1/temis/contract/' + contractID
		And request contract
		When method POST

		* def event = { contract_id: '#(contractID)', entity_id: '#(contractID)', entity_type: 'CONTRACT', event_type: 'CONTRACT_CREATED', parent_type: 'CONTRACT', update_date: '2023-08-03T19:29:06.63556Z'  }
		* string json = event

    	* def result = RegistrationEventsPublisher.publish(json)

		* def expectedResponse =
		"""
			{
			  result: '#(RuleResult.Approved)'
			}
		"""

		Given url baseURLCompliance
		And path '/state/' + contractID
        And header Content-Type = 'application/json'
        And params ({ only_pending: false })
		And retry until response.result == RuleResult.Approved		
        When method GET
		Then assert responseStatus == 200
		And match response contains deep expectedResponse

    Scenario: Should find profile approved

		* def profileID = uuid()
		* def documentNumber = DocumentNormalizer(CPFGenerator())
		* def profile = { profile_id: '#(profileID)', offer_type: '#(PreconditionContractRuleOffer)',  role_type: '#(RoleType.Customer)', profile_type: '#(ProfileType.Individual)', document_number: '#(documentNumber)' }

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

        * def contractID = uuid()
		* def contract = { contract_id: '#(contractID)', profile_id: '#(profileID)', document_id: ''}

		Given url mockURL
		And path '/v1/temis/contract/' + contractID
		And request contract
		When method POST

		* def event = { contract_id: '#(contractID)', entity_id: '#(contractID)', entity_type: 'CONTRACT', event_type: 'CONTRACT_CREATED', parent_type: 'CONTRACT', update_date: '2023-08-03T19:29:06.63556Z'  }
		* string json = event
    	
		* def result = RegistrationEventsPublisher.publish(json)

		* def expectedResponse =
		"""
			{
				result: '#(RuleResult.Approved)'
			}
		"""

		Given url baseURLCompliance
		And path '/state/' + contractID
        And header Content-Type = 'application/json'
        And params ({ only_pending: false })
		And retry until response.result == RuleResult.Approved
        When method GET
		Then assert responseStatus == 200
		And match response contains deep expectedResponse
