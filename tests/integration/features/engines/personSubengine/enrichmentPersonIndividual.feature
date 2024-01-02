Feature: Enrichment at person individual

    Background:
        * url baseURLCompliance

		* def ruleSetOwnershipStructure = "OWNERSHIP_STRUCTURE"
		* def ruleNameShareholders = "SHAREHOLDERS"
		* def ruleSetUnderAge = "IS_UNDER_AGE"
		* def ruleNameUnderAge = "CUSTOMER_IS_UNDER_AGE"
		* def ruleSetLegalRepresentatives = "LEGAL_REPRESENTATIVES"
		* def ruleNameLegalRepresentativesResult = "LEGAL_REPRESENTATIVES_RESULT"

	Scenario: should not enrich the shareholder by enrichment and reject
		
		* def offerType = "TEST_OFFER" + uuid() 
		* def offer = { offer_type : '#(offerType)', product: 'test' }
		
        Given path '/offers'
		And header Content-Type = 'application/json'
		And request offer
		When method POST
		Then assert responseStatus == 201 || responseStatus == 409
		
		* def ruleSetConfigName = "RULE_SET_CONFIG_" + CPFGenerator()
		* def ruleSetConfig = { ownership_structure: {} }

		* def catalogParams = { offer_type: '#(offerType)', role_type: '#(RoleType.Customer)', person_type: '#(ProfileType.Company)', steps: [ { rules_config: #(ruleSetConfig)} ] }
		* def catalog = CreateMultiLevelCatalog(catalogParams)

		Given url mockURL
		And path '/temis-config/cadastral-validation-configs'
		And header Content-Type = 'application/json'
		And request catalog
		When method POST
		Then assert responseStatus == 200

		* def ruleSetConfigName = "RULE_SET_CONFIG_" + CPFGenerator()
		* def ruleSetConfig = { under_age: {} }

		* def catalogParams = { offer_type: '#(offerType)', role_type: '#(RoleType.Shareholder)', person_type: '#(ProfileType.Individual)', steps: [ { rules_config: #(ruleSetConfig) } ], enrich_flag: false}
		* def catalog = CreateMultiLevelCatalog(catalogParams)

		Given url mockURL
		And path '/temis-config/cadastral-validation-configs'
		And header Content-Type = 'application/json'
		And request catalog
		When method POST
		Then assert responseStatus == 200

		* def profileID = uuid()
		* def profileDocumentNumber =  DocumentNormalizer(CNPJGenerator())
		* def profile = { profile_id: '#(profileID)', offer_type: '#(offerType)', role_type: '#(RoleType.Customer)', profile_type: '#(ProfileType.Company)', document_number: '#(profileDocumentNumber)', company: { legal_name: 'CUSTOMER_COMPANY_NAME_LTDA' } }
		* def shareholderDocumentNumber = CPFGenerator()

		Given url mockURL
		And path '/v1/temis/profile/' + profileID
		And request profile
		When method POST
		Then assert responseStatus == 200

		* def ownershipStructure =
		"""
			{
				legal_entity_id: '#(profileDocumentNumber)',
				final_beneficiaries_counted: 1,
				shareholding_sum: 100.0,
				shareholders: [
					{
						parent_legal_entity: '#(profileDocumentNumber)',
						shareholding: 100,
						role: "#(RoleType.Shareholder)",
						type: "#(ProfileType.Individual)",
						name: 'SHAREHOLDER_INDIVIDUAL_NAME_LTDA',
						document_number: '#(shareholderDocumentNumber)',
						nationality: "BRA",
						birth_date: "10/06/2017",
						pep: false
					}
				]
			}
		"""

		Given url mockURL
		And path '/temis-enrichment/ownership-structure/' + profileDocumentNumber
		And request ownershipStructure
		When method POST
		Then assert responseStatus == 200

		* def event = { profile_id: '#(profileID)', entity_id: '#(profileID)', entity_type: 'PROFILE', event_type: 'PROFILE_CHANGED', parent_type: 'PROFILE', update_date: '2023-08-03T19:29:06.63556Z'  }
		* string json = event
    	* def result = RegistrationEventsPublisher.publish(json)
		* eval sleep(2000)		

		* def expectedResponse =
		"""
			{
				result: '#(RuleResult.Rejected)',
				rule_set_result: [
					{
						set: '#(ruleSetOwnershipStructure)',
						name: '#(ruleNameShareholders)',
						result: '#(RuleResult.Rejected)'
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

		* def shareholderID = response.rule_set_result[1].problems[0].detail[0].shareholder_id

		* def expectedResponse =
		"""
			{
				result: '#(RuleResult.Rejected)',
				rule_set_result: [
					{
						set: '#(ruleSetUnderAge)',
						name: '#(ruleNameUnderAge)',
						result: '#(RuleResult.Rejected)'
					}
				]
			}
		"""

		Given url baseURLCompliance
		And path '/state/' + shareholderID
        And header Content-Type = 'application/json'
        And params ({ only_pending: false })		
		And retry until response.result == RuleResult.Rejected
        When method GET
		Then assert responseStatus == 200
		And match response contains deep expectedResponse

	Scenario: should enrich shareholder individual by enrichment and approve
		
		* def offerType = "TEST_OFFER" + uuid() 
		* def offer = { offer_type : '#(offerType)', product: 'test' }
		
        Given path '/offers'
		And header Content-Type = 'application/json'
		And request offer
		When method POST
		Then assert responseStatus == 201 || responseStatus == 409
		
		* def ruleSetConfigName = "RULE_SET_CONFIG_" + CPFGenerator()
		* def ruleSetConfig = { ownership_structure: {} }

        * def catalogParams = { offer_type: '#(offerType)', role_type: '#(RoleType.Customer)', person_type: '#(ProfileType.Company)', steps: [ { rules_config: #(ruleSetConfig) } ] }
		* def catalog = CreateMultiLevelCatalog(catalogParams)

		Given url mockURL
		And path '/temis-config/cadastral-validation-configs'
		And header Content-Type = 'application/json'
		And request catalog
		When method POST
		Then assert responseStatus == 200

		* def ruleSetConfigName = "RULE_SET_CONFIG_" + CPFGenerator()
		* def ruleSetConfig = { under_age: {} }

		* def catalogParams = { offer_type: '#(offerType)', role_type: '#(RoleType.Shareholder)', person_type: '#(ProfileType.Individual)', steps: [ { rules_config: #(ruleSetConfig) } ], enrich_flag: true}
		* def catalog = CreateMultiLevelCatalog(catalogParams)

		Given url mockURL
		And path '/temis-config/cadastral-validation-configs'
		And header Content-Type = 'application/json'
		And request catalog
		When method POST
		Then assert responseStatus == 200

		* def profileID = uuid()
		* def profileDocumentNumber =  DocumentNormalizer(CNPJGenerator())
		* def profile = { profile_id: '#(profileID)', offer_type: '#(offerType)', role_type: '#(RoleType.Customer)', profile_type: '#(ProfileType.Company)', document_number: '#(profileDocumentNumber)', company: { legal_name: 'CUSTOMER_COMPANY_NAME_LTDA' } }

		Given url mockURL
		And path '/v1/temis/profile/' + profileID
		And request profile
		When method POST
		Then assert responseStatus == 200

		* def shareholderDocumentNumber = CPFGenerator()

		* def ownershipStructure =
		"""
			{
				legal_entity_id: '#(profileID)',
				final_beneficiaries_counted: 1,
				shareholding_sum: 100.0,
				shareholders: [
					{
						parent_legal_entity: '#(profileID)',
						shareholding: 100,
						role: "",
						type: "#(ProfileType.Individual)",
						name: 'SHAREHOLDER INDIVIDUAL NAME DA SILVA',
						document_number: '#(shareholderDocumentNumber)',
						nationality: "BRA",
						birth_date: "10/06/2017",
						pep: false
					}
				]
			}
		"""

		Given url mockURL
		And path '/temis-enrichment/ownership-structure/' + profileDocumentNumber
		And request ownershipStructure
		When method POST
		Then assert responseStatus == 200

		* def enrichmentResponse = { birth_date: '06/05/1991', situation: 2, name: 'SHAREHOLDER DA SILVA OLIVEIRA' }

		Given url mockURL
		And path '/temis-enrichment/individual/' + shareholderDocumentNumber
		And request enrichmentResponse
		When method POST

		* def event = { profile_id: '#(profileID)', entity_id: '#(profileID)', entity_type: 'PROFILE', event_type: 'PROFILE_CHANGED', parent_type: 'PROFILE', update_date: '2023-08-03T19:29:06.63556Z'  }
		* string json = event
    	* def result = RegistrationEventsPublisher.publish(json)
		* eval sleep(1000)		

		* def expectedResponse =
		"""
			{
				result: '#(RuleResult.Approved)',
				rule_set_result: [
					{
						set: '#(ruleSetOwnershipStructure)',
						name: '#(ruleNameShareholders)',
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

		* def expectedEnrichment = 
		"""
			{
				"person": {
					"enriched_information": {
						"ownership_structure": {
							"shareholders": [
								{
									"ownership_percent": 100,
									"level": 0,
									"person": {
										"individual": {
											"date_of_birth": "1991-05-06T00:00:00Z",
											"date_of_birth_inputted": "2017-06-10T00:00:00Z",
										},
										"enriched_information": {
											"status": "PENDING_REGULARIZATION",
											"name": "#(enrichmentResponse.name)",
											"birth_date": "#(enrichmentResponse.birth_date)"
										}
									}
								}
							]
						}
					}
				}
			}
		"""

		Given url baseURLCompliance
		And path '/profile/' + profileID
        And header Content-Type = 'application/json'
        When method GET
		Then assert responseStatus == 200
		And match response contains deep expectedEnrichment


	Scenario: should not enrich the shareholder by registration and reject
		
		* def offerType = "TEST_OFFER" + uuid() 
		* def offer = { offer_type : '#(offerType)', product: 'test' }
		
        Given path '/offers'
		And header Content-Type = 'application/json'
		And request offer
		When method POST
		Then assert responseStatus == 201 || responseStatus == 409
		
		* def ruleSetConfigName = "RULE_SET_CONFIG_" + CPFGenerator()
		* def ruleSetConfig = { ownership_structure: {} }

		* def catalogParams = { offer_type: '#(offerType)', role_type: '#(RoleType.Customer)', person_type: '#(ProfileType.Company)', steps: [ { rules_config: #(ruleSetConfig) } ] }
		* def catalog = CreateMultiLevelCatalog(catalogParams)

		Given url mockURL
		And path '/temis-config/cadastral-validation-configs'
		And header Content-Type = 'application/json'
		And request catalog
		When method POST
		Then assert responseStatus == 200

		* def ruleSetConfigName = "RULE_SET_CONFIG_" + CPFGenerator()
		* def ruleSetConfig = { under_age: {} }

		* def catalogParams = { offer_type: '#(offerType)', role_type: '#(RoleType.Shareholder)', person_type: '#(ProfileType.Individual)', steps: [ { rules_config: #(ruleSetConfig) } ], enrich_flag: false}
		* def catalog = CreateMultiLevelCatalog(catalogParams)

		Given url mockURL
		And path '/temis-config/cadastral-validation-configs'
		And header Content-Type = 'application/json'
		And request catalog
		When method POST
		Then assert responseStatus == 200

		* def profileID = uuid()
		* def profileDocumentNumber =  DocumentNormalizer(CNPJGenerator())
		* def profile = { profile_id: '#(profileID)', offer_type: '#(offerType)', role_type: '#(RoleType.Customer)', profile_type: '#(ProfileType.Company)', document_number: '#(profileDocumentNumber)', company: { legal_name: 'CUSTOMER_COMPANY_NAME_LTDA' } }
		* def shareholderDocumentNumber = CPFGenerator()

		Given url mockURL
		And path '/v1/temis/profile/' + profileID
		And request profile
		When method POST
		Then assert responseStatus == 200

		* def ownershipStructure =
		"""
			{
				legal_entity_id: '#(profileDocumentNumber)',
				final_beneficiaries_counted: 1,
				shareholding_sum: 0,
				shareholders: [
					{
						parent_legal_entity: '#(profileDocumentNumber)',
						shareholding: 0,
						role: "#(RoleType.Shareholder)",
						type: "#(ProfileType.Individual)",
						name: 'SHAREHOLDER_INDIVIDUAL_NAME_LTDA',
						document_number: '#(shareholderDocumentNumber)',
						nationality: "BRA",
						birth_date: "10/06/2017",
						pep: false
					}
				]
			}
		"""

		Given url mockURL
		And path '/temis-enrichment/ownership-structure/' + profileDocumentNumber
		And request ownershipStructure
		When method POST
		Then assert responseStatus == 200

		* def shareholderID = uuid()
		* def ownershipStructure =
		"""
			{
				legal_entity_id: '#(profileDocumentNumber)',
				final_beneficiaries_counted: 1,
				shareholding_sum: 100.0,
				shareholders: [
					{
						shareholder_id: '#(shareholderID)',
						parent_legal_entity: '#(profileDocumentNumber)',
						shareholding: 100,
						role: "#(RoleType.Shareholder)",
						type: "#(ProfileType.Individual)",
						name: 'SHAREHOLDER_INDIVIDUAL_NAME_LTDA',
						document_number: '#(shareholderDocumentNumber)',
						nationality: "BRA",
						birth_date: "10/06/2017",
						pep: false
					}
				]
			}
		"""

		Given url mockURL
		And path '/v1/temis/ownership-structure/' + profileID
		And request ownershipStructure
		When method POST
		Then assert responseStatus === 200

		* def event = { profile_id: '#(profileID)', entity_id: '#(profileID)', entity_type: 'PROFILE', event_type: 'PROFILE_CHANGED', parent_type: 'PROFILE', update_date: '2023-08-03T19:29:06.63556Z'  }
		* string json = event
    	* def result = RegistrationEventsPublisher.publish(json)
		* eval sleep(1000)		

		* def expectedResponse =
		"""
			{
				result: '#(RuleResult.Rejected)',
				rule_set_result: [
					{
						set: '#(ruleSetOwnershipStructure)',
						name: '#(ruleNameShareholders)',
						result: '#(RuleResult.Rejected)'
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

		* def shareholderID = response.rule_set_result[1].problems[0].detail[0].shareholder_id

		* def expectedResponse =
		"""
			{
				result: '#(RuleResult.Rejected)',
				rule_set_result: [
					{
						set: '#(ruleSetUnderAge)',
						name: '#(ruleNameUnderAge)',
						result: '#(RuleResult.Rejected)'
					}
				]
			}
		"""

		Given url baseURLCompliance
		And path '/state/' + shareholderID
        And header Content-Type = 'application/json'
        And params ({ only_pending: false })		
		And retry until response.result == RuleResult.Rejected
        When method GET
		Then assert responseStatus == 200
		And match response contains deep expectedResponse

	Scenario: should enrich shareholder individual by registration and approve
		
		* def offerType = "OFFER_" + CPFGenerator()
		* def offer = { offer_type : '#(offerType)', product: 'test' }
		
        Given path '/offers'
		And header Content-Type = 'application/json'
		And request offer
		When method POST
		Then assert responseStatus == 201 || responseStatus == 409
		
		* def ruleSetConfigName = "RULE_SET_CONFIG_" + CPFGenerator()
		* def ruleSetConfig = { ownership_structure: {} }

		* def catalogParams = { offer_type: '#(offerType)', role_type: '#(RoleType.Customer)', person_type: '#(ProfileType.Company)', steps: [ { rules_config: #(ruleSetConfig) } ] }
		* def catalog = CreateMultiLevelCatalog(catalogParams)

		Given url mockURL
		And path '/temis-config/cadastral-validation-configs'
		And header Content-Type = 'application/json'
		And request catalog
		When method POST
		Then assert responseStatus == 200

		* def ruleSetConfigName = "RULE_SET_CONFIG_" + CPFGenerator()
		* def ruleSetConfig = { under_age: {} }

		* def catalogParams = { offer_type: '#(offerType)', role_type: '#(RoleType.Shareholder)', person_type: '#(ProfileType.Individual)', steps: [ { rules_config: #(ruleSetConfig) } ], enrich_flag: true}
		* def catalog = CreateMultiLevelCatalog(catalogParams)

		Given url mockURL
		And path '/temis-config/cadastral-validation-configs'
		And header Content-Type = 'application/json'
		And request catalog
		When method POST
		Then assert responseStatus == 200

		* def profileID = uuid()
		* def profileDocumentNumber =  DocumentNormalizer(CNPJGenerator())
		* def profile = { profile_id: '#(profileID)', offer_type: '#(offerType)', role_type: '#(RoleType.Customer)', profile_type: '#(ProfileType.Company)', document_number: '#(profileDocumentNumber)', company: { legal_name: 'CUSTOMER_COMPANY_NAME_LTDA' } }

		Given url mockURL
		And path '/v1/temis/profile/' + profileID
		And request profile
		When method POST
		Then assert responseStatus == 200

		* def shareholderDocumentNumber = CPFGenerator()

		* def ownershipStructure =
		"""
			{
				legal_entity_id: '#(profileID)',
				final_beneficiaries_counted: 1,
				shareholding_sum: 0,
				shareholders: [
					{
						parent_legal_entity: '#(profileID)',
						shareholding: 0,
						role: "",
						type: "#(ProfileType.Individual)",
						name: 'SHAREHOLDER_INDIVIDUAL_NAME_LTDA',
						document_number: '#(shareholderDocumentNumber)',
						nationality: "BRA",
						birth_date: "10/06/2017",
						pep: false
					}
				]
			}
		"""

		Given url mockURL
		And path '/temis-enrichment/ownership-structure/' + profileDocumentNumber
		And request ownershipStructure
		When method POST
		Then assert responseStatus == 200

		* def shareholderID = uuid()
		* def ownershipStructure =
		"""
			{
				legal_entity_id: '#(profileDocumentNumber)',
				final_beneficiaries_counted: 1,
				shareholding_sum: 100.0,
				shareholders: [
					{
						shareholder_id: '#(shareholderID)',
						parent_legal_entity: '#(profileDocumentNumber)',
						shareholding: 100,
						role: "#(RoleType.Shareholder)",
						type: "#(ProfileType.Individual)",
						name: 'SHAREHOLDER_INDIVIDUAL_NAME_LTDA',
						document_number: '#(shareholderDocumentNumber)',
						nationality: "BRA",
						birth_date: "10/06/2017",
						pep: false
					}
				]
			}
		"""

		Given url mockURL
		And path '/v1/temis/ownership-structure/' + profileID
		And request ownershipStructure
		When method POST
		Then assert responseStatus === 200

		* def enrichmentResponse = { birth_date: '06/05/1991', situation: 2 }

		Given url mockURL
		And path '/temis-enrichment/individual/' + shareholderDocumentNumber
		And request enrichmentResponse
		When method POST

		* def event = { profile_id: '#(profileID)', entity_id: '#(profileID)', entity_type: 'PROFILE', event_type: 'PROFILE_CHANGED', parent_type: 'PROFILE', update_date: '2023-08-03T19:29:06.63556Z'  }
		* string json = event
    	* def result = RegistrationEventsPublisher.publish(json)
		* eval sleep(1000)		

		* def expectedResponse =
		"""
			{
				result: '#(RuleResult.Approved)',
				rule_set_result: [
					{
						set: '#(ruleSetOwnershipStructure)',
						name: '#(ruleNameShareholders)',
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

	Scenario: should not enrich legal representative individual and reject

		* def offerType = "TEST_OFFER" + uuid() 
		* def offer = { offer_type : '#(offerType)', product: 'test' }
		
        Given path '/offers'
		And header Content-Type = 'application/json'
		And request offer
		When method POST
		Then assert responseStatus == 201 || responseStatus == 409
		
		* def ruleSetConfigName = "RULE_SET_CONFIG_" + CPFGenerator()
		* def ruleSetConfig = { legal_representative: {} }

		* def catalogParams = { offer_type: '#(offerType)', role_type: '#(RoleType.Customer)', person_type: '#(ProfileType.Company)', steps: [ { rules_config: #(ruleSetConfig) } ] }
		* def catalog = CreateMultiLevelCatalog(catalogParams)

		Given url mockURL
		And path '/temis-config/cadastral-validation-configs'
		And header Content-Type = 'application/json'
		And request catalog
		When method POST
		Then assert responseStatus == 200

		* def ruleSetConfigName = "RULE_SET_CONFIG_" + CPFGenerator()
		* def ruleSetConfig = { under_age: {} }

		* def catalogParams = { offer_type: '#(offerType)', role_type: 'LEGAL_REPRESENTATIVE', person_type: '#(ProfileType.Individual)', steps: [ { rules_config: #(ruleSetConfig) } ], enrich_flag: false}
		* def catalog = CreateMultiLevelCatalog(catalogParams)

		Given url mockURL
		And path '/temis-config/cadastral-validation-configs'
		And header Content-Type = 'application/json'
		And request catalog
		When method POST
		Then assert responseStatus == 200

		* def profileID = uuid()
		* def profileDocumentNumber =  DocumentNormalizer(CNPJGenerator())
		* def profile = { profile_id: '#(profileID)', offer_type: '#(offerType)', role_type: '#(RoleType.Customer)', profile_type: '#(ProfileType.Company)', document_number: '#(profileDocumentNumber)', company: { legal_name: 'CUSTOMER_COMPANY_NAME_LTDA' } }

		Given url mockURL
		And path '/v1/temis/profile/' + profileID
		And request profile
		When method POST
		Then assert responseStatus == 200

		* def legalRepresentativeId = uuid()
		* def legalRepresentativeDocumentNumber = CPFGenerator()

		* def legalRepresentatives =
		"""
			[
				{
					legal_representative_id: '#(legalRepresentativeId)',
					profile_id: '#(legalRepresentativeId)',
					document_number: '#(legalRepresentativeDocumentNumber)',
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
		* eval sleep(1000)		

		* def expectedResponse =
		"""
			{
				result: '#(RuleResult.Rejected)',
				rule_set_result: [
					{
						set: '#(ruleSetLegalRepresentatives)',
						name: '#(ruleNameLegalRepresentativesResult)',
						result: '#(RuleResult.Rejected)'
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

		* def shareholderID = response.rule_set_result[0].problems[0].detail[0]

		* def expectedResponse =
		"""
			{
				result: '#(RuleResult.Rejected)',
				rule_set_result: [
					{
						set: '#(ruleSetUnderAge)',
						name: '#(ruleNameUnderAge)',
						result: '#(RuleResult.Rejected)'
					}
				]
			}
		"""

		Given url baseURLCompliance
		And path '/state/' + shareholderID
        And header Content-Type = 'application/json'
        And params ({ only_pending: false })		
		And retry until response.result == RuleResult.Rejected
        When method GET
		Then assert responseStatus == 200
		And match response contains deep expectedResponse

	Scenario: should enrich legal representative individual and approve
		
		* def offerType = "TEST_OFFER" + uuid() 
		* def offer = { offer_type : '#(offerType)', product: 'test' }
		
        Given path '/offers'
		And header Content-Type = 'application/json'
		And request offer
		When method POST
		Then assert responseStatus == 201 || responseStatus == 409
		
		* def ruleSetConfigName = "RULE_SET_CONFIG_" + CPFGenerator()
		* def ruleSetConfig = { legal_representative: {} }

		* def catalogParams = { offer_type: '#(offerType)', role_type: '#(RoleType.Customer)', person_type: '#(ProfileType.Company)', steps: [ { rules_config: #(ruleSetConfig) } ] }
		* def catalog = CreateMultiLevelCatalog(catalogParams)

		Given url mockURL
		And path '/temis-config/cadastral-validation-configs'
		And header Content-Type = 'application/json'
		And request catalog
		When method POST
		Then assert responseStatus == 200

		* def ruleSetConfigName = "RULE_SET_CONFIG_" + CPFGenerator()
		* def ruleSetConfig = { under_age: {} }

		* def catalogParams = { offer_type: '#(offerType)', role_type: 'LEGAL_REPRESENTATIVE', person_type: '#(ProfileType.Individual)', steps: [ { rules_config: #(ruleSetConfig) } ], enrich_flag: true}
		* def catalog = CreateMultiLevelCatalog(catalogParams)

		Given url mockURL
		And path '/temis-config/cadastral-validation-configs'
		And header Content-Type = 'application/json'
		And request catalog
		When method POST
		Then assert responseStatus == 200

		* def profileID = uuid()
		* def profileDocumentNumber =  DocumentNormalizer(CNPJGenerator())
		* def profile = { profile_id: '#(profileID)', offer_type: '#(offerType)', role_type: '#(RoleType.Customer)', profile_type: '#(ProfileType.Company)', document_number: '#(profileDocumentNumber)', company: { legal_name: 'CUSTOMER_COMPANY_NAME_LTDA' } }

		Given url mockURL
		And path '/v1/temis/profile/' + profileID
		And request profile
		When method POST
		Then assert responseStatus == 200

		* def legalRepresentativeId = uuid()
		* def legalRepresentativeDocumentNumber = CPFGenerator()

		* def legalRepresentatives =
		"""
			[
				{
					legal_representative_id: '#(legalRepresentativeId)',
					profile_id: '#(legalRepresentativeId)',
					document_number: '#(legalRepresentativeDocumentNumber)',
				}
			]
		"""

		Given url mockURL
		And path '/v1/temis/legal-representatives'
		And params { profile_id : '#(profileID)' }
		And request legalRepresentatives
		When method POST

		* def enrichmentResponse = { birth_date: '06/05/1991', situation: 2 }

		Given url mockURL
		And path '/temis-enrichment/individual/' + legalRepresentativeDocumentNumber
		And request enrichmentResponse
		When method POST

		* def event = { profile_id: '#(profileID)', entity_id: '#(profileID)', entity_type: 'PROFILE', event_type: 'PROFILE_CHANGED', parent_type: 'PROFILE', update_date: '2023-08-03T19:29:06.63556Z'  }
		* string json = event
    	* def result = RegistrationEventsPublisher.publish(json)
		* eval sleep(1000)		

		* def expectedResponse =
		"""
			{
				result: '#(RuleResult.Approved)',
				rule_set_result: [
					{
						set: '#(ruleSetLegalRepresentatives)',
						name: '#(ruleNameLegalRepresentativesResult)',
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