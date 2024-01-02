Feature: Test Legal Representative Rule

    Background:
        * url baseURLCompliance

	    * def LegalRepresentativesRuleOffer = "TEST_OFFER" + uuid()   
		* def ProblemCodeLegalRepresentativeNotApproved = "LEGAL_REPRESENTATIVE_NOT_APPROVED"

		* def offer = { offer_type : '#(LegalRepresentativesRuleOffer)', product: 'maquininha'}

        Given path '/offers'
		And header Content-Type = 'application/json'
		And request offer
		When method POST

        * def ruleSetConfig = { legal_representative: {} }

		Given url mockURL
        And path '/temis-config/cadastral-validation-configs'
        And header Content-Type = 'application/json'
		And def catalog = call CreateSingleLevelCatalog { offer_type: '#(LegalRepresentativesRuleOffer)', role_type: '#(RoleType.Customer)', person_type: '#(ProfileType.Company)', account_flag: false, rules_config: #(ruleSetConfig)}
		And request catalog
		When method POST	

		* def ruleSetConfig = { bureau: {} }

		Given url mockURL
        And path '/temis-config/cadastral-validation-configs'
        And header Content-Type = 'application/json'
		And def catalog = call CreateSingleLevelCatalog { offer_type: '#(LegalRepresentativesRuleOffer)', role_type: '#(RoleType.LegalRepresentative)', person_type: '#(ProfileType.Individual)', account_flag: false, rules_config: #(ruleSetConfig)}
		And request catalog
		When method POST

        * def profileID = uuid()
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


		* def legalRepresentativeID1 = uuid()
		* def legalRepresentativeID2 = uuid()
		* def expectedLegalRepresentative1CreatedStateEvent = 
		"""
		{
			"id": "#string",
			"entity_id": "#(legalRepresentativeID1)",
			"entity_type": "#(EntityType.ComplianceState)",
			"event_type": "#(EventType.State.Created)",
			"update_date":"#string",
			"data": {
				"content": "#(personStateEventContentSchema)",
			}
		}	
		"""

		* def expectedLegalRepresentative1ChangedStateEvent = 
		"""
		{
			"id": "#string",
			"entity_id": "#(legalRepresentativeID1)",
			"entity_type": "#(EntityType.ComplianceState)",
			"event_type": "#(EventType.State.Changed)",
			"update_date":"#string",
			"data": {
				"content": "#(personStateEventContentSchema)",
			}
		}	
		"""

		* def expectedLegalRepresentative2CreatedStateEvent = 
		"""
		{
			"id": "#string",
			"entity_id": "#(legalRepresentativeID2)",
			"entity_type": "#(EntityType.ComplianceState)",
			"event_type": "#(EventType.State.Created)",
			"update_date":"#string",
			"data": {
				"content": "#(personStateEventContentSchema)",
			}
		}	
		"""

		* def expectedLegalRepresentative2ChangedStateEvent = 
		"""
		{
			"id": "#string",
			"entity_id": "#(legalRepresentativeID2)",
			"entity_type": "#(EntityType.ComplianceState)",
			"event_type": "#(EventType.State.Changed)",
			"update_date":"#string",
			"data": {
				"content": "#(personStateEventContentSchema)",
			}
		}	
		"""

    Scenario: Legal representatives is rejeceted by Bureau

		* def documentNumber = DocumentNormalizer(CNPJGenerator())
		* def profile = {profile_id:'#(profileID)', document_number: '#(documentNumber)', offer_type: '#(LegalRepresentativesRuleOffer)', role_type: '#(RoleType.Customer)', profile_type: '#(ProfileType.Company)'}

		Given url mockURL
		And path '/v1/temis/profile/' + profileID
		And request profile
		When method POST
		Then assert responseStatus == 200

		* def legalRepresentatives = 
		"""
		[
			{
				legal_representative_id: '#(legalRepresentativeID1)',
				profile_id: '#(profileID)',
				document_number: '35125595063',
				foreign_tax_residency: true,
				country_of_tax_residency: "USA"
			},
			{
				legal_representative_id: '#(legalRepresentativeID2)',
				profile_id: '#(profileID)',
				document_number: '07937447095',
				foreign_tax_residency: true,
				country_of_tax_residency: "USA"
			}
		]
		"""

		Given url mockURL
		And path '/v1/temis/legal-representatives'
		And params { profile_id : '#(profileID)' }
		And request legalRepresentatives
		When method POST

		* def event = { profile_id: '#(profileID)', entity_id: '#(profileID)', entity_type: 'PROFILE', event_type: 'PROFILE_CHANGED', parent_type: 'PROFILE', update_date: '2023-08-03T19:29:06.63556Z'  }
		* string json = event

    	* def result = RegistrationEventsPublisher.publish(json)

		* def metadata1 = "Legal Representative " + legalRepresentativeID1 + " is not Approved"
		* def metadata2 = "Legal Representative "+legalRepresentativeID2+" is not Approved"
		* def expected = 
		"""
		{
			"entity_id": "#(profileID)",
			"engine_name": "#(Engine.Profile)",
			"rule_set_result": [
				{
				"result": "#(RuleResult.Rejected)",
				"metadata": [
					"#(metadata1)",
					"#(metadata2)"
				],
				"set": "#(RuleSet.LegalRepresentatives)",
				"pending": false,
				"name": "#(RuleName.LegalRepresentativeResult)",
				"problems": [
					{
					"code": "LEGAL_REPRESENTATIVE_NOT_APPROVED",
					"detail": [
						"#(legalRepresentativeID1)",
						"#(legalRepresentativeID2)"
					]
					}
				],
				"step_number" : 0 
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
		Then match response.rule_set_result == expected.rule_set_result

		* def expectedLegalRepresentatives = 
		"""
		[
			{
			  "document_number": "#(legalRepresentatives[0].document_number)",
			  "role_type": "#(EntityType.LegalRepresentative)",
			  "individual": {
				"foreign_tax_residency": "#(legalRepresentatives[0].foreign_tax_residency)",
				"country_of_tax_residency": "#(legalRepresentatives[0].country_of_tax_residency)",
				"phones": [
					{
						"country_code": "",
						"number": "",
						"area_code": "",
						"type": ""
					}
				],
			  },
			  "entity_id": "#string",
			  "entity_type": "#(EntityType.LegalRepresentative)",
			  "profile_id": "#(legalRepresentatives[0].profile_id)",
			  "legal_representative_id": "#(legalRepresentatives[0].legal_representative_id)",
			  "offer_type": "#(LegalRepresentativesRuleOffer)"
			},
			{
			  "document_number": "#(legalRepresentatives[1].document_number)",
			  "role_type": "#(EntityType.LegalRepresentative)",
			  "individual": {
				"foreign_tax_residency": "#(legalRepresentatives[1].foreign_tax_residency)",
				"country_of_tax_residency": "#(legalRepresentatives[1].country_of_tax_residency)",
				"phones": [
					{
						"country_code": "",
						"number": "",
						"area_code": "",
						"type": ""
					}
				],
			  },
			  "entity_id": "#string",
			  "entity_type": "#(EntityType.LegalRepresentative)",
			  "profile_id": "#(legalRepresentatives[1].profile_id)",
			  "legal_representative_id": "#(legalRepresentatives[1].legal_representative_id)",
			  "offer_type": "#(LegalRepresentativesRuleOffer)"
			}
		  ]
		  """

		Given url baseURLCompliance
		And path '/profile/' + profileID
		And header Content-Type = 'application/json'
		When method GET
		Then assert responseStatus == 200
		Then match response.legal_representatives contains deep expectedLegalRepresentatives[0]
		Then match response.legal_representatives contains deep expectedLegalRepresentatives[1]

		Given url mockURL
		And path "/subscribe/state-events/" + profileID
		And retry until response.length == 2
		When method GET
		Then assert responseStatus == 200
		And match response contains deep expectedProfileCreatedStateEvent
		And match response contains deep expectedProfileChangedStateEvent

        Given url mockURL
		And path "/subscribe/state-events/" + legalRepresentativeID1
		And retry until response.length == 2
		When method GET
		Then assert responseStatus == 200
		And match response contains deep expectedLegalRepresentative1CreatedStateEvent
		And match response contains deep expectedLegalRepresentative1ChangedStateEvent

        Given url mockURL
		And path "/subscribe/state-events/" + legalRepresentativeID2
		And retry until response.length == 2
		When method GET
		Then assert responseStatus == 200
		And match response contains deep expectedLegalRepresentative2CreatedStateEvent
		And match response contains deep expectedLegalRepresentative2ChangedStateEvent

	@CreateLegalRepresentativeApprovedSerasa
	Scenario: Legal representatives approved

		* def documentNumber = DocumentNormalizer(CNPJGenerator())
		* def profile = {profile_id:'#(profileID)', document_number: '#(documentNumber)', offer_type: '#(LegalRepresentativesRuleOffer)',  role_type: '#(RoleType.Customer)', profile_type: '#(ProfileType.Company)'}
		* profile.expiration_date = "2024-01-01T00:00:00Z"

		Given url mockURL
		And path '/v1/temis/profile/' + profileID
		And request profile
		When method POST
		Then assert responseStatus == 200

		* def legalRepresentatives = 
		"""
		[
        	{
				legal_representative_id: '#(legalRepresentativeID1)',
				profile_id: "#(profileID)",
				document_number: '27832498048',
				foreign_tax_residency: true,
				country_of_tax_residency: "USA",
				expiration_date: "2024-02-02T00:00:00Z"
			},
			{
				legal_representative_id: '#(legalRepresentativeID2)',
				profile_id: "#(profileID)",
				document_number: '94527672002',
				foreign_tax_residency: true,
				country_of_tax_residency: "USA",
				expiration_date: "2024-03-03T00:00:00Z"
			}
    	]
		"""

		Given url mockURL
		And path '/v1/temis/legal-representatives'
		And params { profile_id : '#(profileID)' }
		And request legalRepresentatives
		When method POST

		* def enrichmentResponse = { situation: 1 }

		Given url mockURL
		And path '/temis-enrichment/individual/' + legalRepresentatives[0].document_number
		And request enrichmentResponse
		When method POST

		Given url mockURL
		And path '/temis-enrichment/individual/' + legalRepresentatives[1].document_number
		And request enrichmentResponse
		When method POST

		* def event = { profile_id: '#(profileID)', entity_id: '#(profileID)', entity_type: 'PROFILE', event_type: 'PROFILE_CHANGED', parent_type: 'PROFILE', update_date: '2023-08-03T19:29:06.63556Z'  }
		* string json = event

    	* def result = RegistrationEventsPublisher.publish(json)
		
		* def expected =
		"""
		{
			"entity_id": "#(profileID)",
			"engine_name": "#(Engine.Profile)",
			"result": "#(RuleResult.Approved)",
			"rule_set_result": [
			  {
				"result": "#(RuleResult.Approved)",
				"metadata": null,
				"set": "#(RuleSet.LegalRepresentatives)",
				"pending": false,
				"name": "#(RuleName.LegalRepresentativeResult)",
				"step_number" : 0 
			  }
			]
		  }
		"""

		Given url baseURLCompliance
		And path '/state/' + profileID
        And header Content-Type = 'application/json'
        And params ({ only_pending: false })
		And retry until response.rule_set_result[0].result == RuleResult.Approved
        When method GET
		Then match response == expected

		* def expectedLegalRepresentatives = 
		"""
		[
			{
			  "document_number": "#(legalRepresentatives[0].document_number)",
			  "role_type": "#(EntityType.LegalRepresentative)",
			  "individual": {
				"foreign_tax_residency": "#(legalRepresentatives[0].foreign_tax_residency)",
				"country_of_tax_residency": "#(legalRepresentatives[0].country_of_tax_residency)",
				"phones": [
					{
						"country_code": "",
						"number": "",
						"area_code": "",
						"type": ""
					}
				],
			  },
			  "entity_id": "#string",
			  "entity_type": "#(EntityType.LegalRepresentative)",
			  "profile_id": "#(legalRepresentatives[0].profile_id)",
			  "legal_representative_id": "#(legalRepresentatives[0].legal_representative_id)",
			  "offer_type": "#(LegalRepresentativesRuleOffer)"
			},
			{
			  "document_number": "#(legalRepresentatives[1].document_number)",
			  "role_type": "#(EntityType.LegalRepresentative)",
			  "individual": {
				"foreign_tax_residency": "#(legalRepresentatives[1].foreign_tax_residency)",
				"country_of_tax_residency": "#(legalRepresentatives[1].country_of_tax_residency)",
				"phones": [
					{
						"country_code": "",
						"number": "",
						"area_code": "",
						"type": ""
					}
				],
			  },
			  "entity_id": "#string",
			  "entity_type": "#(EntityType.LegalRepresentative)",
			  "profile_id": "#(legalRepresentatives[1].profile_id)",
			  "legal_representative_id": "#(legalRepresentatives[1].legal_representative_id)",
			  "offer_type": "#(LegalRepresentativesRuleOffer)"
			}
		  ]
		  """

		Given url baseURLCompliance
		And path '/profile/' + profileID
		And header Content-Type = 'application/json'
		When method GET
		Then assert responseStatus == 200
		Then match response.legal_representatives contains deep expectedLegalRepresentatives[0]
		Then match response.legal_representatives contains deep expectedLegalRepresentatives[1]
		And assert response.expiration_date == "2024-01-01T00:00:00Z"
		And assert response.legal_representatives[0].expiration_date == "2024-02-02T00:00:00Z"
		And assert response.legal_representatives[1].expiration_date == "2024-03-03T00:00:00Z"

		Given url mockURL
		And path "/subscribe/state-events/" + profileID
		And retry until response.length == 2
		When method GET
		Then assert responseStatus == 200
		And match response contains deep expectedProfileCreatedStateEvent
		And match response contains deep expectedProfileChangedStateEvent

        Given url mockURL
		And path "/subscribe/state-events/" + legalRepresentativeID1
		And retry until response.length == 2
		When method GET
		Then assert responseStatus == 200
		And match response contains deep expectedLegalRepresentative1CreatedStateEvent
		And match response contains deep expectedLegalRepresentative1ChangedStateEvent

        Given url mockURL
		And path "/subscribe/state-events/" + legalRepresentativeID2
		And retry until response.length == 2
		When method GET
		Then assert responseStatus == 200
		And match response contains deep expectedLegalRepresentative2CreatedStateEvent
		And match response contains deep expectedLegalRepresentative2ChangedStateEvent