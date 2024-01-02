Feature: Test Ownership Structure Rule

	Background:
		* url baseURLCompliance

		* def OwnershipStructureRuleOffer = "TEST_OFFER" + uuid()  
		* def PartnerOwnershipStructure = "PartnerOwnershipStructure"

		* def problemCodeShareholdingNotAchieveMinimunRequired = "SHAREHOLDING_NOT_ACHIEVE_MINIMUM_REQUIRED"

		* def offer = { offer_type : '#(OwnershipStructureRuleOffer)', product : 'maquininha'}

		Given path '/offers'
		And header Content-Type = 'application/json'
		And request offer
		When method POST

		* def ruleSetConfig = { ownership_structure: {} }

		Given url mockURL
        And path '/temis-config/cadastral-validation-configs'
        And header Content-Type = 'application/json'
		And def catalog = call CreateSingleLevelCatalog { offer_type: '#(OwnershipStructureRuleOffer)', role_type: '#(RoleType.Customer)', person_type: '#(ProfileType.Company)', account_flag: false, rules_config: #(ruleSetConfig)}
		And request catalog
		When method POST

		* def ruleSetConfig = {  bureau: {} }

		Given url mockURL
        And path '/temis-config/cadastral-validation-configs'
        And header Content-Type = 'application/json'
		And def catalog = call CreateSingleLevelCatalog { offer_type: '#(OwnershipStructureRuleOffer)', role_type: '#(RoleType.Shareholder)', person_type: '#(ProfileType.Individual)', account_flag: false, rules_config: #(ruleSetConfig)}
		And request catalog
		When method POST

		Given url mockURL
        And path '/temis-config/cadastral-validation-configs'
        And header Content-Type = 'application/json'
		And def catalog = call CreateSingleLevelCatalog { offer_type: '#(OwnershipStructureRuleOffer)', role_type: '#(RoleType.Shareholder)', person_type: '#(ProfileType.Company)', account_flag: false, rules_config: #(ruleSetConfig)}
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

		* def expectedShareholder1CreatedStateEvent = 
		"""
		{
			"id": "#string",
			"entity_id": "#(shareholderID1)",
			"entity_type": "#(EntityType.ComplianceState)",
			"event_type": "#(EventType.State.Created)",
			"update_date":"#string",
			"data": {
				"content": "#(personStateEventContentSchema)",
			}
		}	
		"""

		* def expectedShareholder1ChangedStateEvent = 
		"""
		{
			"id": "#string",
			"entity_id": "#(shareholderID1)",
			"entity_type": "#(EntityType.ComplianceState)",
			"event_type": "#(EventType.State.Changed)",
			"update_date":"#string",
			"data": {
				"content": "#(personStateEventContentSchema)",
			}
		}	
		"""

		* def expectedShareholder2CreatedStateEvent = 
		"""
		{
			"id": "#string",
			"entity_id": "#(shareholderID2)",
			"entity_type": "#(EntityType.ComplianceState)",
			"event_type": "#(EventType.State.Created)",
			"update_date":"#string",
			"data": {
				"content": "#(personStateEventContentSchema)",
			}
		}	
		"""

		* def expectedShareholder2ChangedStateEvent = 
		"""
		{
			"id": "#string",
			"entity_id": "#(shareholderID2)",
			"entity_type": "#(EntityType.ComplianceState)",
			"event_type": "#(EventType.State.Changed)",
			"update_date":"#string",
			"data": {
				"content": "#(personStateEventContentSchema)",
			}
		}	
		"""

		* def expectedShareholder3CreatedStateEvent = 
		"""
		{
			"id": "#string",
			"entity_id": "#(shareholderID3)",
			"entity_type": "#(EntityType.ComplianceState)",
			"event_type": "#(EventType.State.Created)",
			"update_date":"#string",
			"data": {
				"content": "#(personStateEventContentSchema)",
			}
		}	
		"""

		* def expectedShareholder3ChangedStateEvent = 
		"""
		{
			"id": "#string",
			"entity_id": "#(shareholderID3)",
			"entity_type": "#(EntityType.ComplianceState)",
			"event_type": "#(EventType.State.Changed)",
			"update_date":"#string",
			"data": {
				"content": "#(personStateEventContentSchema)",
			}
		}	
		"""
	
		* def expectedShareholder4CreatedStateEvent = 
		"""
		{
			"id": "#string",
			"entity_id": "#(shareholderID4)",
			"entity_type": "#(EntityType.ComplianceState)",
			"event_type": "#(EventType.State.Created)",
			"update_date":"#string",
			"data": {
				"content": "#(personStateEventContentSchema)",
			}
		}	
		"""

		* def expectedShareholder4ChangedStateEvent = 
		"""
		{
			"id": "#string",
			"entity_id": "#(shareholderID4)",
			"entity_type": "#(EntityType.ComplianceState)",
			"event_type": "#(EventType.State.Changed)",
			"update_date":"#string",
			"data": {
				"content": "#(personStateEventContentSchema)",
			}
		}	
		"""

	Scenario: Should not approve Shareholding and ignore Shareholds validation

		* def documentNumber = DocumentNormalizer(CNPJGenerator())
		* def profile = {profile_id:'#(profileID)', document_number: '#(documentNumber)', offer_type: '#(OwnershipStructureRuleOffer)', role_type: '#(RoleType.Customer)', profile_type: '#(ProfileType.Company)'}

		Given url mockURL
		And path '/v1/temis/profile/' + profileID
		And request profile
		When method POST
		Then assert responseStatus == 200

		* def shareholder1 = DocumentNormalizer(CNPJGenerator())
		* def shareholder2 = CPFGenerator()
		* def shareholder3 = CPFGenerator()
		* def func = read('classpath:features/rules/profile/ownershipStructure/data/scenario1/enrichedOwnershipStructureLess95.js')
		* def ownershipStructure = call func {document_number : '#(documentNumber)', shareholder1: '#(shareholder1)', shareholder2: '#(shareholder2)', shareholder3: '#(shareholder3)'}

		Given url mockURL
		And path '/temis-enrichment/ownership-structure/' + documentNumber
		And request ownershipStructure
		When method POST
		Then assert responseStatus == 200

		* def ownershipStructure = { final_beneficiaries_counted: 10,  shareholding_sum: 90}

		Given url mockURL
		And path '/v1/temis/ownership-structure/' + profileID
		And request ownershipStructure
		When method POST
		Then assert responseStatus == 200

		* def event = { profile_id: '#(profileID)', entity_id: '#(profileID)', entity_type: 'PROFILE', event_type: 'PROFILE_CHANGED', parent_type: 'PROFILE', update_date: '2023-08-03T19:29:06.63556Z'  }
		* string json = event
		
		* def result = RegistrationEventsPublisher.publish(json)
		* eval sleep(1000)

		* def expectedShareholdingResult = 
		"""
		{
			"step_number" : 0 ,
			"set": "#(RuleSet.OwnershipStructure)",
			"name": "#(RuleName.Shareholding)",
			"result": "#(RuleResult.Analysing)",
			"metadata": "#string",
			"pending": true,
			"problems": "##[]"
		}	
		"""
		* def expectedShareholdersResult =
		"""
		{
			"step_number" : 0,
			"set": "#(RuleSet.OwnershipStructure)",
			"name": "#(RuleName.Shareholders)",
			"result": "#(RuleResult.Ignored)",
			"metadata": null,
			"pending": false,
		}	
		"""

		Given url baseURLCompliance
		And path '/state/' + profileID
		And header Content-Type = 'application/json'
		And params ({ only_pending: false })
		And retry until response.result ==  RuleResult.Analysing
		When method GET
		And match response.rule_set_result contains deep expectedShareholdingResult
		And match response.rule_set_result contains deep expectedShareholdersResult

		Given url mockURL
		And path "/subscribe/state-events/" + profileID
		And retry until response.length == 2
		When method GET
		Then assert responseStatus == 200
		And match response contains deep expectedProfileCreatedStateEvent
		And match response contains deep expectedProfileChangedStateEvent

	Scenario: Should approve Shareholding, by enriched ownership structure, and not approve Shareholders

		* def documentNumber = DocumentNormalizer(CNPJGenerator())
		* def profile =  { profile_id: '#(profileID)', document_number: '#(documentNumber)', offer_type: '#(OwnershipStructureRuleOffer)', role_type: '#(RoleType.Customer)', profile_type: '#(ProfileType.Company)'}

		Given url mockURL
		And path '/v1/temis/profile/' + profileID
		And request profile
		When method POST
		Then assert responseStatus == 200

		* def shareholder1 = DocumentNormalizer(CNPJGenerator())
		* def shareholder2 = CPFGenerator()
		* def shareholder3 = CPFGenerator()
		* def func = read('classpath:features/rules/profile/ownershipStructure/data/scenario2/enrichedOwnershipStructureMore95.js')
		* def ownershipStructure = call func {document_number : '#(documentNumber)', shareholder1: '#(shareholder1)', shareholder2: '#(shareholder2)', shareholder3: '#(shareholder3)'}

		Given url mockURL
		And path '/temis-enrichment/ownership-structure/' + documentNumber
		And request ownershipStructure
		When method POST
		Then assert responseStatus == 200

		* def enrichmentResponse = { situation: 3 }

		Given url mockURL
		And path '/temis-enrichment/legal-entity/' + shareholder1
		And request enrichmentResponse
		When method POST

		* def enrichmentResponse = { situation: 4 }

		Given url mockURL
		And path '/temis-enrichment/individual/' + shareholder2
		And request enrichmentResponse
		When method POST

		Given url mockURL
		And path '/temis-enrichment/individual/' + shareholder3
		And request enrichmentResponse
		When method POST

		* def event = { profile_id: '#(profileID)', entity_id: '#(profileID)', entity_type: 'PROFILE', event_type: 'PROFILE_CHANGED', parent_type: 'PROFILE', update_date: '2023-08-03T19:29:06.63556Z'  }
		* string json = event
		
		* def result = RegistrationEventsPublisher.publish(json)
		* eval sleep(1000)

		* def expectedShareholdingResult = 
		"""
		{
			"step_number" : 0,
			"set": "#(RuleSet.OwnershipStructure)",
			"name": "#(RuleName.Shareholding)",
			"result": "#(RuleResult.Approved)",
			"pending": false,
			"metadata": "#null"
		}	
		"""
		* def metadata1 = "Shareholder with Document Number "+shareholder1+" is not Approved"
		* def metadata2 = "Shareholder with Document Number "+shareholder2+" is not Approved"
		* def metadata3 = "Shareholder with Document Number "+shareholder3+" is not Approved"
		* def expectedShareholdersResult =
		"""
		{
			"step_number" : 0,
			"set": "#(RuleSet.OwnershipStructure)",
			"name": "#(RuleName.Shareholders)",
			"result": "#(RuleResult.Rejected)",
			"pending": false,
			"metadata": [
				"#(metadata1)",
				"#(metadata2)",
				"#(metadata3)"
			],
			"problems": [
				{
					"code": "SHAREHOLDER_NOT_APPROVED",
					"detail": [
						{
							"document_number": "#(shareholder1)",
							"shareholder_id": "#string"
						},
						{
							"document_number": "#(shareholder2)",
							"shareholder_id": "#string"
						},
						{
							"document_number": "#(shareholder3)",
							"shareholder_id": "#string"
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
		And retry until response.result ==  RuleResult.Rejected
		When method GET
		Then match response.rule_set_result contains deep expectedShareholdingResult
		Then match response.rule_set_result[1].metadata == expectedShareholdersResult.metadata
		Then match response.rule_set_result[1].problems == expectedShareholdersResult.problems

		Given url mockURL
		And path "/subscribe/state-events/" + profileID
		And retry until response.length == 2
		When method GET
		Then assert responseStatus == 200
		And match response contains deep expectedProfileCreatedStateEvent
		And match response contains deep expectedProfileChangedStateEvent

		Given url baseURLCompliance
		And path '/profile/' + profileID
		And header Content-Type = 'application/json'
		When method GET
		Then assert responseStatus == 200

		* def profile = response
        * def shareholderID1 = profile.person.enriched_information.ownership_structure.shareholders[0].shareholder_id
		* def shareholderID2 = profile.person.enriched_information.ownership_structure.shareholders[1].shareholder_id
		* def shareholderID3 = profile.person.enriched_information.ownership_structure.shareholders[2].shareholder_id

		Given url mockURL
		And path "/subscribe/state-events/" + shareholderID1
		And retry until response.length == 2
		When method GET
		Then assert responseStatus == 200
		And match response contains deep expectedShareholder1CreatedStateEvent
		And match response contains deep expectedShareholder1ChangedStateEvent

		Given url mockURL
		And path "/subscribe/state-events/" + shareholderID2
		And retry until response.length == 2
		When method GET
		Then assert responseStatus == 200
		And match response contains deep expectedShareholder2CreatedStateEvent
		And match response contains deep expectedShareholder2ChangedStateEvent

		Given url mockURL
		And path "/subscribe/state-events/" + shareholderID3
		And retry until response.length == 2
		When method GET
		Then assert responseStatus == 200
		And match response contains deep expectedShareholder3CreatedStateEvent
		And match response contains deep expectedShareholder3ChangedStateEvent

	@CreateEnrichedOwnerShipApprovedShareholding
	Scenario: Should approve Shareholding, by enriched ownership structure, and approve Shareholders Rule

		* def documentNumber = DocumentNormalizer(CNPJGenerator())
		* def profile =  { profile_id: '#(profileID)', document_number: '#(documentNumber)', offer_type: '#(OwnershipStructureRuleOffer)', role_type: '#(RoleType.Customer)', profile_type: '#(ProfileType.Company)'}

		Given url mockURL
		And path '/v1/temis/profile/' + profileID
		And request profile
		When method POST
		Then assert responseStatus == 200

		* def shareholder1 = DocumentNormalizer(CNPJGenerator())
		* def shareholder2 = CPFGenerator()
		* def shareholder3 = CPFGenerator()
		* def func = read('classpath:features/rules/profile/ownershipStructure/data/scenario3/enrichedOwnershipStructureMore95.js')
		* def ownershipStructure = call func {document_number : '#(documentNumber)', shareholder1: '#(shareholder1)', shareholder2: '#(shareholder2)', shareholder3: '#(shareholder3)'}

		Given url mockURL
		And path '/temis-enrichment/ownership-structure/' + documentNumber
		And request ownershipStructure
		When method POST
		Then assert responseStatus == 200

		* def enrichResponse = { situation: 2 }

		Given url mockURL
		And path '/temis-enrichment/legal-entity/' + shareholder1
		And request enrichResponse
		When method POST

		* def enrichResponse = { situation: 1 }

		Given url mockURL
		And path '/temis-enrichment/individual/' + shareholder2
		And request enrichResponse
		When method POST

		Given url mockURL
		And path '/temis-enrichment/individual/' + shareholder3
		And request enrichResponse
		When method POST

		* def event = { profile_id: '#(profileID)', entity_id: '#(profileID)', entity_type: 'PROFILE', event_type: 'PROFILE_CHANGED', parent_type: 'PROFILE', update_date: '2023-08-03T19:29:06.63556Z'  }
		* string json = event
		
		* def result = RegistrationEventsPublisher.publish(json)
		* eval sleep(1000)

		* def expectedShareholdingResult =
		"""
		{
			"result": "#(RuleResult.Approved)",
			"metadata": null,
			"set": "#(RuleSet.OwnershipStructure)",
			"pending": false,
			"name": "#(RuleName.Shareholding)",
			"step_number" : 0 
		}
		"""
		* def expectedShareholdersResult =
		"""
		{
			"result": "#(RuleResult.Approved)",
			"metadata": null,
			"set": "#(RuleSet.OwnershipStructure)",
			"pending": false,
			"name": "#(RuleName.Shareholders)",
			"step_number" : 0 
		}	
		"""

		Given url baseURLCompliance
		And path '/state/' + profileID
		And header Content-Type = 'application/json'
		And params ({ only_pending: false })
		And retry until response.result ==  RuleResult.Approved
		When method GET
		Then match response.rule_set_result contains deep expectedShareholdingResult
		Then match response.rule_set_result contains deep expectedShareholdersResult

		Given url mockURL
		And path "/subscribe/state-events/" + profileID
		And retry until response.length == 2
		When method GET
		Then assert responseStatus == 200
		And match response contains deep expectedProfileCreatedStateEvent
		And match response contains deep expectedProfileChangedStateEvent

		Given url baseURLCompliance
		And path '/profile/' + profileID
		And header Content-Type = 'application/json'
		When method GET
		Then assert responseStatus == 200

		* def profile = response
        * def shareholderID1 = profile.person.enriched_information.ownership_structure.shareholders[0].shareholder_id
		* def shareholderID2 = profile.person.enriched_information.ownership_structure.shareholders[1].shareholder_id
		* def shareholderID3 = profile.person.enriched_information.ownership_structure.shareholders[2].shareholder_id

		Given url mockURL
		And path "/subscribe/state-events/" + shareholderID1
		And retry until response.length == 2
		When method GET
		Then assert responseStatus == 200
		And match response contains deep expectedShareholder1CreatedStateEvent
		And match response contains deep expectedShareholder1ChangedStateEvent

		Given url mockURL
		And path "/subscribe/state-events/" + shareholderID2
		And retry until response.length == 2
		When method GET
		Then assert responseStatus == 200
		And match response contains deep expectedShareholder2CreatedStateEvent
		And match response contains deep expectedShareholder2ChangedStateEvent

		Given url mockURL
		And path "/subscribe/state-events/" + shareholderID3
		And retry until response.length == 2
		When method GET
		Then assert responseStatus == 200
		And match response contains deep expectedShareholder3CreatedStateEvent
		And match response contains deep expectedShareholder3ChangedStateEvent

	@CreateOwnerShipApprovedShareholding
	Scenario: Should approve Shareholding, by manually filled ownership structure, and approve Shareholders Rule

		* def documentNumber = DocumentNormalizer(CNPJGenerator())
		* def profile =  { profile_id: '#(profileID)', document_number: '#(documentNumber)', offer_type: '#(OwnershipStructureRuleOffer)', role_type: '#(RoleType.Customer)', profile_type: '#(ProfileType.Company)'}

		Given url mockURL
		And path '/v1/temis/profile/' + profileID
		And request profile
		When method POST
		Then assert responseStatus == 200

		* def func = read('classpath:features/rules/profile/ownershipStructure/data/scenario4/enrichedOwnershipStructureLegalEntity.js')
		* def ownershipStructure = call func { document_number: '#(documentNumber)'}

		Given url mockURL
		And path '/temis-enrichment/ownership-structure/' + documentNumber
		And request ownershipStructure
		When method POST
		Then assert responseStatus == 200

		* def shareholder1 = DocumentNormalizer(CNPJGenerator())
		* def shareholder2 = CPFGenerator()
		* def shareholder3 = CPFGenerator()
		* def shareholder4 = CPFGenerator()
		* def shareholderID1 = uuid()
		* def shareholderID2 = uuid()
		* def shareholderID3 = uuid()
		* def shareholderID4 = uuid()
		* def param = 
		"""
			{
				document_number : '#(documentNumber)', 
				shareholder1: '#(shareholder1)', 
				shareholder2: '#(shareholder2)', 
				shareholder3: '#(shareholder3)', 
				shareholder4: '#(shareholder4)',
				shareholderID1: '#(shareholderID1)',
				shareholderID2: '#(shareholderID2)',
				shareholderID3: '#(shareholderID3)',
				shareholderID4: '#(shareholderID4)'
			}
		"""
		* def func = read('classpath:features/rules/profile/ownershipStructure/data/scenario4/manuallyFilledOwnershipStructure.js')
		* def ownershipStructure = call func param

		Given url mockURL
		And path '/v1/temis/ownership-structure/' + profileID
		And request ownershipStructure
		When method POST
		Then assert responseStatus === 200

		* def enrichResponse = { situation: 2 }

		Given url mockURL
		And path '/temis-enrichment/legal-entity/' + shareholder1
		And request enrichResponse
		When method POST

		* def enrichResponse = { situation: 1 }

		Given url mockURL
		And path '/temis-enrichment/individual/' + shareholder2
		And request enrichResponse
		When method POST

		Given url mockURL
		And path '/temis-enrichment/individual/' + shareholder3
		And request enrichResponse
		When method POST

		Given url mockURL
		And path '/temis-enrichment/individual/' + shareholder4
		And request enrichResponse
		When method POST

		* def event = { profile_id: '#(profileID)', entity_id: '#(profileID)', entity_type: 'PROFILE', event_type: 'PROFILE_CHANGED', parent_type: 'PROFILE', update_date: '2023-08-03T19:29:06.63556Z'  }
		* string json = event
		
		* def result = RegistrationEventsPublisher.publish(json)
		* eval sleep(2000)

		* def expectedShareholdingResult =
		"""
		{
			"result": "#(RuleResult.Approved)",
			"metadata": null,
			"set": "#(RuleSet.OwnershipStructure)",
			"pending": false,
			"name": "#(RuleName.Shareholding)",
			"step_number" : 0 
		}
		"""
		* def expectedShareholdersResult =
		"""
		{
			"result": "#(RuleResult.Approved)",
			"metadata": null,
			"set": "#(RuleSet.OwnershipStructure)",
			"pending": false,
			"name": "#(RuleName.Shareholders)",
			"step_number" : 0 
		}
		"""	

		Given url baseURLCompliance
		And path '/state/' + profileID
		And header Content-Type = 'application/json'
		And params ({ only_pending: false })
		And retry until response.result ==  RuleResult.Approved
		When method GET
		Then match response.rule_set_result contains deep expectedShareholdingResult
		Then match response.rule_set_result contains deep expectedShareholdersResult

		Given url mockURL
		And path "/subscribe/state-events/" + profileID
		And retry until response.length == 2
		When method GET
		Then assert responseStatus == 200
		And match response contains deep expectedProfileCreatedStateEvent
		And match response contains deep expectedProfileChangedStateEvent

		Given url mockURL
		And path "/subscribe/state-events/" + shareholderID1
		And retry until response.length == 2
		When method GET
		Then assert responseStatus == 200
		And match response contains deep expectedShareholder1CreatedStateEvent
		And match response contains deep expectedShareholder1ChangedStateEvent

		Given url mockURL
		And path "/subscribe/state-events/" + shareholderID2
		And retry until response.length == 2
		When method GET
		Then assert responseStatus == 200
		And match response contains deep expectedShareholder2CreatedStateEvent
		And match response contains deep expectedShareholder2ChangedStateEvent

		Given url mockURL
		And path "/subscribe/state-events/" + shareholderID3
		And retry until response.length == 2
		When method GET
		Then assert responseStatus == 200
		And match response contains deep expectedShareholder3CreatedStateEvent
		And match response contains deep expectedShareholder3ChangedStateEvent

		Given url mockURL
		And path "/subscribe/state-events/" + shareholderID4
		And retry until response.length == 2
		When method GET
		Then assert responseStatus == 200
		And match response contains deep expectedShareholder4CreatedStateEvent
		And match response contains deep expectedShareholder4ChangedStateEvent

	Scenario: should not validate enriched shareholders that have same document number, for a given Ownership Structure

		* def documentNumber = DocumentNormalizer(CNPJGenerator())
		* def profile =  { profile_id: '#(profileID)', document_number: '#(documentNumber)', offer_type: '#(OwnershipStructureRuleOffer)', role_type: '#(RoleType.Customer)', profile_type: '#(ProfileType.Company)'}

		Given url mockURL
		And path '/v1/temis/profile/' + profileID
		And request profile
		When method POST
		Then assert responseStatus == 200

		* def shareholder1 = DocumentNormalizer(CNPJGenerator())
		* def shareholder2 = CPFGenerator()
		* def shareholder3 = CPFGenerator()
		* def func = read('classpath:features/rules/profile/ownershipStructure/data/scenario5/enrichedOwnershipStructure.js');
		* def ownershipStructure = call func {document_number : '#(documentNumber)', shareholder1: '#(shareholder1)', shareholder2: '#(shareholder2)', shareholder3: '#(shareholder3)', shareholder4: '#(shareholder2)'}

		Given url mockURL
		And path '/temis-enrichment/ownership-structure/' + documentNumber
		And request ownershipStructure
		When method POST
		Then assert responseStatus == 200

		* def enrichResponse = { situation: 3 }

		Given url mockURL
		And path '/temis-enrichment/legal-entity/' + shareholder1
		And request enrichResponse
		When method POST

		Given url mockURL
		And path '/temis-enrichment/individual/' + shareholder2
		And request enrichResponse
		When method POST

		Given url mockURL
		And path '/temis-enrichment/individual/' + shareholder3
		And request enrichResponse
		When method POST

		* def event = { profile_id: '#(profileID)', entity_id: '#(profileID)', entity_type: 'PROFILE', event_type: 'PROFILE_CHANGED', parent_type: 'PROFILE', update_date: '2023-08-03T19:29:06.63556Z'  }
		* string json = event
		
		* def result = RegistrationEventsPublisher.publish(json)
		* eval sleep(1000)

		Given url baseURLCompliance
		And path '/state/' + profileID
		And header Content-Type = 'application/json'
		And params ({ only_pending: true })
		And retry until response.result ==  RuleResult.Rejected
		When method GET
		Then assert response.result == RuleResult.Rejected

		Given url mockURL
		And path "/subscribe/state-events/" + profileID
		And retry until response.length == 2
		When method GET
		Then assert responseStatus == 200
		And match response contains deep expectedProfileCreatedStateEvent
		And match response contains deep expectedProfileChangedStateEvent

		Given url baseURLCompliance
		And path "/profile/" + profileID
		When method GET
		Then assert responseStatus == 200

		* def profile = response
        * def shareholderID1 = profile.person.enriched_information.ownership_structure.shareholders[0].shareholder_id
		* def shareholderID2 = profile.person.enriched_information.ownership_structure.shareholders[1].shareholder_id
		* def shareholderID3 = profile.person.enriched_information.ownership_structure.shareholders[2].shareholder_id

		Given url mockURL
		And path "/subscribe/state-events/" + shareholderID1
		And retry until response.length == 2
		When method GET
		Then assert responseStatus == 200
		And match response contains deep expectedShareholder1CreatedStateEvent
		And match response contains deep expectedShareholder1ChangedStateEvent

		Given url mockURL
		And path "/subscribe/state-events/" + shareholderID2
		And retry until response.length == 2
		When method GET
		Then assert responseStatus == 200
		And match response contains deep expectedShareholder2CreatedStateEvent
		And match response contains deep expectedShareholder2ChangedStateEvent

		Given url mockURL
		And path "/subscribe/state-events/" + shareholderID3
		And retry until response.length == 2
		When method GET
		Then assert responseStatus == 200
		And match response contains deep expectedShareholder3CreatedStateEvent
		And match response contains deep expectedShareholder3ChangedStateEvent

	Scenario: should not validate manually filled shareholders that have same document number, for a given Ownership Structure

		* def documentNumber = DocumentNormalizer(CNPJGenerator())
		* def profile =  { profile_id: '#(profileID)', document_number: '#(documentNumber)', offer_type: '#(OwnershipStructureRuleOffer)', role_type: '#(RoleType.Customer)', profile_type: '#(ProfileType.Company)'}

		Given url mockURL
		And path '/v1/temis/profile/' + profileID
		And request profile
		When method POST
		Then assert responseStatus == 200

		* def func = read('classpath:features/rules/profile/ownershipStructure/data/scenario6/enrichedOwnershipStructure.js')
		* def ownershipStructure = call func {document_number : '#(documentNumber)' }

		Given url mockURL
		And path '/temis-enrichment/ownership-structure/' + documentNumber
		And request ownershipStructure
		When method POST
		Then assert responseStatus == 200

		* def shareholder1 = DocumentNormalizer(CNPJGenerator())
		* def shareholder2 = CPFGenerator()
		* def shareholder3 = CPFGenerator()
		* def shareholderID1 = uuid()
		* def shareholderID2 = uuid()
		* def shareholderID3 = uuid()
		* def shareholderID4 = uuid()
		* def params = 
		"""
			{
				document_number : '#(documentNumber)', 
				shareholder1: '#(shareholder1)', 
				shareholder2: '#(shareholder2)', 
				shareholder3: '#(shareholder3)', 
				shareholder4: '#(shareholder2)',
				shareholderID1: '#(shareholderID1)',
				shareholderID2: '#(shareholderID2)',
				shareholderID3: '#(shareholderID3)',
				shareholderID4: '#(shareholderID4)'
			}
		"""
		* def func = read('classpath:features/rules/profile/ownershipStructure/data/scenario6/manuallyFilledOwnershipStructure.js')
		* def ownershipStructure = call func params
		
		Given url mockURL
		And path '/v1/temis/ownership-structure/' + profileID
		And request ownershipStructure
		When method POST
		Then assert responseStatus == 200

		* def enrichResponse = { situation: 3 }

		Given url mockURL
		And path '/temis-enrichment/legal-entity/' + shareholder1
		And request enrichResponse
		When method POST

		Given url mockURL
		And path '/temis-enrichment/individual/' + shareholder2
		And request enrichResponse
		When method POST

		Given url mockURL
		And path '/temis-enrichment/individual/' + shareholder3
		And request enrichResponse
		When method POST

		* def event = { profile_id: '#(profileID)', entity_id: '#(profileID)', entity_type: 'PROFILE', event_type: 'PROFILE_CHANGED', parent_type: 'PROFILE', update_date: '2023-08-03T19:29:06.63556Z'  }
		* string json = event
		
		* def result = RegistrationEventsPublisher.publish(json)
		* eval sleep(1500)

		Given url baseURLCompliance
		And path '/state/' + profileID
		And header Content-Type = 'application/json'
		And params ({ only_pending: true })
		And retry until response.result ==  RuleResult.Rejected
		When method GET
		Then assert response.result == RuleResult.Rejected

		Given url mockURL
		And path "/subscribe/state-events/" + profileID
		And retry until response.length == 2
		When method GET
		Then assert responseStatus == 200
		And match response contains deep expectedProfileCreatedStateEvent
		And match response contains deep expectedProfileChangedStateEvent

		Given url mockURL
		And path "/subscribe/state-events/" + shareholderID1
		And retry until response.length == 2
		When method GET
		Then assert responseStatus == 200
		And match response contains deep expectedShareholder1CreatedStateEvent
		And match response contains deep expectedShareholder1ChangedStateEvent

		Given url mockURL
		And path "/subscribe/state-events/" + shareholderID2
		And retry until response.length == 2
		When method GET
		Then assert responseStatus == 200
		And match response contains deep expectedShareholder2CreatedStateEvent
		And match response contains deep expectedShareholder2ChangedStateEvent

		Given url mockURL
		And path "/subscribe/state-events/" + shareholderID3
		And retry until response.length == 2
		When method GET
		Then assert responseStatus == 200
		And match response contains deep expectedShareholder3CreatedStateEvent
		And match response contains deep expectedShareholder3ChangedStateEvent

		Given url mockURL
		And path "/subscribe/state-events/" + shareholderID4
		And retry until response.length == 0
		When method GET
		Then assert responseStatus == 200

	Scenario: should validate shareholders from manually filled ownership structure after overriding shareholding rule

		* def documentNumber = DocumentNormalizer(CNPJGenerator())
		* def profile =  { profile_id: '#(profileID)', document_number: '#(documentNumber)', offer_type: '#(OwnershipStructureRuleOffer)', role_type: '#(RoleType.Customer)', profile_type: '#(ProfileType.Company)'}

		Given url mockURL
		And path '/v1/temis/profile/' + profileID
		And request profile
		When method POST
		Then assert responseStatus == 200

		* def shareholder1 = CPFGenerator()
		* def shareholder2 = CPFGenerator()
		* def shareholder3 = CPFGenerator()
		* def func = read('classpath:features/rules/profile/ownershipStructure/data/enrichedLessThan95.js');
		* def enrichedOwnershipStructure = call func {document_number: '#(documentNumber)', shareholder1: '#(shareholder1)', shareholder2: '#(shareholder2)', shareholder3: '#(shareholder3)'}

		Given url mockURL
		And path '/temis-enrichment/ownership-structure/' + documentNumber
		And request enrichedOwnershipStructure
		When method POST
		Then assert responseStatus == 200

		* def func = read('classpath:features/rules/profile/ownershipStructure/data/filledLessThan95.js')
		* def shareholderID1 = uuid()
		* def shareholderID2 = uuid()
		* def shareholderID3 = uuid()
		* def shareholderID4 = uuid()
		* def params =
		"""
			{
				document_number: '#(documentNumber)', 
				shareholder1: '#(shareholder1)', 
				shareholder2: '#(shareholder2)', 
				shareholder3: '#(shareholder3)',
				shareholderID1:'#(shareholderID1)',
				shareholderID2: '#(shareholderID2)',
				shareholderID3: '#(shareholderID3)',
				shareholderID4: '#(shareholderID4)'}
		"""
		* def filledOwnershipStructure = call func params

		Given url mockURL
		And path '/v1/temis/ownership-structure/' + profileID
		And request filledOwnershipStructure
		When method POST
		Then assert responseStatus == 200

		* def shareholderCompanyDocument = filledOwnershipStructure.shareholders[0].document_number
		* def shareholderIndividualOneDocument = filledOwnershipStructure.shareholders[0].shareholders[0].document_number
		* def shareholderIndividualTwoDocument = filledOwnershipStructure.shareholders[0].shareholders[1].document_number

		Given url mockURL
		And path '/temis-enrichment/legal-entity/' + shareholderCompanyDocument
		And request { situation: 2 }
		When method POST

		Given url mockURL
		And path '/temis-enrichment/individual/' + shareholderIndividualOneDocument
		And request { situation: 1 }
		When method POST

		Given url mockURL
		And path '/temis-enrichment/individual/' + shareholderIndividualTwoDocument
		And request { situation: 1 }
		When method POST

		* def event = { profile_id: '#(profileID)', entity_id: '#(profileID)', entity_type: 'PROFILE', event_type: 'PROFILE_CHANGED', parent_type: 'PROFILE', update_date: '2023-08-03T19:29:06.63556Z'  }
		* string json = event
		
		* def result = RegistrationEventsPublisher.publish(json)
		* eval sleep(1000)

		Given url baseURLCompliance
		And path '/state/' + profileID
		And retry until response.result == RuleResult.Analysing
		When method GET
		Then assert responseStatus == 200
		Then assert response.result == RuleResult.Analysing
		Then assert response.rule_set_result[0].name == "SHAREHOLDING"
		Then assert response.rule_set_result[0].result == RuleResult.Analysing
		Then assert response.rule_set_result[1].name == "SHAREHOLDERS"
		Then assert response.rule_set_result[1].result == RuleResult.Ignored

		* def override = 
		"""
			{
				"entity_id": "#(profileID)",
				"entity_type": "#(Engine.Profile)",
				"rule_set": "OWNERSHIP_STRUCTURE",
				"rule_name": "SHAREHOLDING",
				"result": "#(RuleResult.Approved)",
				"comments": "me",
				"author": "me"
			}
		"""

		Given path '/override'
		And request override
		When method POST
		Then assert responseStatus == 200

		Given url baseURLCompliance
		And path '/state/' + profileID
		And retry until response.result == RuleResult.Approved
		When method GET
		Then assert responseStatus == 200
		Then assert response.result == RuleResult.Approved
		Then assert response.rule_set_result[0].name == "SHAREHOLDING"
		Then assert response.rule_set_result[0].result == RuleResult.Approved
		Then assert response.rule_set_result[1].name == "SHAREHOLDERS"
		Then assert response.rule_set_result[1].result == RuleResult.Approved

		Given url baseURLCompliance
		And path "/profile/" + profileID
		When method GET
		Then assert responseStatus == 200
		Then assert response.ownership_structure.shareholders.length == 3

		Given url mockURL
		And path "/subscribe/state-events/" + profileID
		And retry until response.length == 3
		When method GET
		Then assert responseStatus == 200
		And match response contains deep expectedProfileCreatedStateEvent
		And match response contains deep expectedProfileChangedStateEvent

		Given url mockURL
		And path "/subscribe/state-events/" + shareholderID1
		And retry until response.length == 2
		When method GET
		Then assert responseStatus == 200
		And match response contains deep expectedShareholder1CreatedStateEvent
		And match response contains deep expectedShareholder1ChangedStateEvent

		Given url mockURL
		And path "/subscribe/state-events/" + shareholderID2
		And retry until response.length == 2
		When method GET
		Then assert responseStatus == 200
		And match response contains deep expectedShareholder2CreatedStateEvent
		And match response contains deep expectedShareholder2ChangedStateEvent

		Given url mockURL
		And path "/subscribe/state-events/" + shareholderID3
		And retry until response.length == 2
		When method GET
		Then assert responseStatus == 200
		And match response contains deep expectedShareholder3CreatedStateEvent
		And match response contains deep expectedShareholder3ChangedStateEvent

		Given url mockURL
		And path "/subscribe/state-events/" + shareholderID4
		And retry until response.length == 0
		When method GET
		Then assert responseStatus == 200

	Scenario: should validate shareholders from enriched ownership structure after overriding shareholding rule

		* def documentNumber = DocumentNormalizer(CNPJGenerator())
		* def profile =  { profile_id: '#(profileID)', document_number: '#(documentNumber)', offer_type: '#(OwnershipStructureRuleOffer)', role_type: '#(RoleType.Customer)', profile_type: '#(ProfileType.Company)'}

		Given url mockURL
		And path '/v1/temis/profile/' + profileID
		And request profile
		When method POST
		Then assert responseStatus == 200

		* def shareholder1 = CPFGenerator()
		* def shareholder2 = CPFGenerator()
		* def shareholder3 = CPFGenerator()
		* def func = read('classpath:features/rules/profile/ownershipStructure/data/enrichedLessThan95.js');
		* def enrichedOwnershipStructure = call func {document_number: '#(documentNumber)', shareholder1: '#(shareholder1)', shareholder2: '#(shareholder2)', shareholder3: '#(shareholder3)'}

		Given url mockURL
		And path '/temis-enrichment/ownership-structure/' + documentNumber
		And request enrichedOwnershipStructure
		When method POST
		Then assert responseStatus == 200

		* def shareholderCompanyDocument = enrichedOwnershipStructure.shareholders[0].document_number
		* def shareholderIndividualOneDocument = enrichedOwnershipStructure.shareholders[0].shareholders[0].document_number
		* def shareholderIndividualTwoDocument = enrichedOwnershipStructure.shareholders[0].shareholders[1].document_number

		Given url mockURL
		And path '/temis-enrichment/legal-entity/' + shareholderCompanyDocument
		And request { situation: 2 }
		When method POST

		Given url mockURL
		And path '/temis-enrichment/individual/' + shareholderIndividualOneDocument
		And request { situation: 1 }
		When method POST

		Given url mockURL
		And path '/temis-enrichment/individual/' + shareholderIndividualTwoDocument
		And request { situation: 1 }
		When method POST

		* def event = { profile_id: '#(profileID)', entity_id: '#(profileID)', entity_type: 'PROFILE', event_type: 'PROFILE_CHANGED', parent_type: 'PROFILE', update_date: '2023-08-03T19:29:06.63556Z'  }
		* string json = event
		
		* def result = RegistrationEventsPublisher.publish(json)
		* eval sleep(1000)

		Given url baseURLCompliance
		And path '/state/' + profileID
		And retry until response.result == RuleResult.Analysing
		When method GET
		Then assert responseStatus == 200
		Then assert response.result == RuleResult.Analysing
		Then assert response.rule_set_result[0].name == "SHAREHOLDING"
		Then assert response.rule_set_result[0].result == RuleResult.Analysing
		Then assert response.rule_set_result[1].name == "SHAREHOLDERS"
		Then assert response.rule_set_result[1].result == RuleResult.Ignored

		* def override = 
		"""
			{
				"entity_id": "#(profileID)",
				"entity_type": "#(Engine.Profile)",
				"rule_set": "OWNERSHIP_STRUCTURE",
				"rule_name": "SHAREHOLDING",
				"result": "#(RuleResult.Approved)",
				"comments": "me",
				"author": "me"
			}
		"""

		Given path '/override'
		And request override
		When method POST
		Then assert responseStatus == 200

		Given url baseURLCompliance
		And path '/state/' + profileID
		And retry until response.result == RuleResult.Approved
		When method GET
		Then assert responseStatus == 200
		Then assert response.result == RuleResult.Approved
		Then assert response.rule_set_result[0].name == "SHAREHOLDING"
		Then assert response.rule_set_result[0].result == RuleResult.Approved
		Then assert response.rule_set_result[1].name == "SHAREHOLDERS"
		Then assert response.rule_set_result[1].result == RuleResult.Approved

		Given url baseURLCompliance
		And path "/profile/" + profileID
		When method GET
		Then assert responseStatus == 200
		Then assert response.person.enriched_information.ownership_structure.shareholders.length == 3

		* def profile = response
        * def shareholderID1 = profile.person.enriched_information.ownership_structure.shareholders[0].shareholder_id
		* def shareholderID2 = profile.person.enriched_information.ownership_structure.shareholders[1].shareholder_id
		* def shareholderID3 = profile.person.enriched_information.ownership_structure.shareholders[2].shareholder_id

		Given url mockURL
		And path "/subscribe/state-events/" + profileID
		And retry until response.length == 3
		When method GET
		Then assert responseStatus == 200
		And match response contains deep expectedProfileCreatedStateEvent
		And match response contains deep expectedProfileChangedStateEvent

		Given url mockURL
		And path "/subscribe/state-events/" + shareholderID1
		And retry until response.length == 2
		When method GET
		Then assert responseStatus == 200
		And match response contains deep expectedShareholder1CreatedStateEvent
		And match response contains deep expectedShareholder1ChangedStateEvent

		Given url mockURL
		And path "/subscribe/state-events/" + shareholderID2
		And retry until response.length == 2
		When method GET
		Then assert responseStatus == 200
		And match response contains deep expectedShareholder2CreatedStateEvent
		And match response contains deep expectedShareholder2ChangedStateEvent

		Given url mockURL
		And path "/subscribe/state-events/" + shareholderID3
		And retry until response.length == 2
		When method GET
		Then assert responseStatus == 200
		And match response contains deep expectedShareholder3CreatedStateEvent
		And match response contains deep expectedShareholder3ChangedStateEvent

	Scenario: should not validate shareholders after overriding shareholding rule

		* def documentNumber = DocumentNormalizer(CNPJGenerator())
		* def profile =  { profile_id: '#(profileID)', document_number: '#(documentNumber)', offer_type: '#(OwnershipStructureRuleOffer)', role_type: '#(RoleType.Customer)', profile_type: '#(ProfileType.Company)'}

		Given url mockURL
		And path '/v1/temis/profile/' + profileID
		And request profile
		When method POST
		Then assert responseStatus == 200

		* def event = { profile_id: '#(profileID)', entity_id: '#(profileID)', entity_type: 'PROFILE', event_type: 'PROFILE_CHANGED', parent_type: 'PROFILE', update_date: '2023-08-03T19:29:06.63556Z'  }
		* string json = event
		
		* def result = RegistrationEventsPublisher.publish(json)
		* eval sleep(1000)

		Given url baseURLCompliance
		And path '/state/' + profileID
		And retry until response.result == RuleResult.Analysing
		When method GET
		Then assert responseStatus == 200
		Then assert response.result == RuleResult.Analysing
		Then assert response.rule_set_result[0].name == "SHAREHOLDING"
		Then assert response.rule_set_result[0].result == RuleResult.Analysing
		Then assert response.rule_set_result[1].name == "SHAREHOLDERS"
		Then assert response.rule_set_result[1].result == RuleResult.Ignored

		* def override = 
		"""
			{
				"entity_id": "#(profileID)",
				"entity_type": "#(Engine.Profile)",
				"rule_set": "OWNERSHIP_STRUCTURE",
				"rule_name": "SHAREHOLDING",
				"result": "#(RuleResult.Approved)",
				"comments": "me",
				"author": "me"
			}
		"""

		Given path '/override'
		And request override
		When method POST
		Then assert responseStatus == 200

		Given url baseURLCompliance
		And path '/state/' + profileID
		And retry until response.result == RuleResult.Approved
		When method GET
		Then assert responseStatus == 200
		Then assert response.result == RuleResult.Approved
		Then assert response.rule_set_result[0].name == "SHAREHOLDING"
		Then assert response.rule_set_result[0].result == RuleResult.Approved
		Then assert response.rule_set_result[1].name == "SHAREHOLDERS"
		Then assert response.rule_set_result[1].result == RuleResult.Ignored

		Given url mockURL
		And path "/subscribe/state-events/" + profileID
		And retry until response.length == 3
		When method GET
		Then assert responseStatus == 200
		And match response contains deep expectedProfileCreatedStateEvent
		And match response contains deep expectedProfileChangedStateEvent

	Scenario: Should set Shareholding and Shareholders rules as Ignored if profile LegalNature is one of the following:  2143,2046,1210,3204,3212,3999 

		* def documentNumber = DocumentNormalizer(CNPJGenerator())
		* def profile =  { profile_id: '#(profileID)', document_number: '#(documentNumber)', offer_type: '#(OwnershipStructureRuleOffer)', role_type: '#(RoleType.Customer)', profile_type: '#(ProfileType.Company)', company: {legal_nature: "2143"}}

		Given url mockURL
		And path '/v1/temis/profile/' + profileID
		And request profile
		When method POST
		Then assert responseStatus == 200

		* def func = read('classpath:features/rules/profile/ownershipStructure/data/scenario3/enrichedOwnershipStructureMore95.js')
		* def ownershipStructure = call func {}

		Given url mockURL
		And path '/temis-enrichment/ownership-structure/' + documentNumber
		And request ownershipStructure
		When method POST
		Then assert responseStatus == 200

		* def enrichResponse = { situation: 2 }

		Given url mockURL
		And path '/temis-enrichment/legal-entity/' + '83298913000163'
		And request enrichResponse
		When method POST

		* def enrichResponse = { situation: 1 }

		* def event = { profile_id: '#(profileID)', entity_id: '#(profileID)', entity_type: 'PROFILE', event_type: 'PROFILE_CHANGED', parent_type: 'PROFILE', update_date: '2023-08-03T19:29:06.63556Z'  }
		* string json = event
		
		* def result = RegistrationEventsPublisher.publish(json)
		* eval sleep(1000)

		Given url baseURLCompliance
		And path '/state/' + profileID
		And header Content-Type = 'application/json'
		And params ({ only_pending: false })
		And retry until response.result ==  RuleResult.Approved
		When method GET
		Then match response.rule_set_result contains deep {set: "OWNERSHIP_STRUCTURE", name: "SHAREHOLDING", result: "#(RuleResult.Ignored)" }
		Then match response.rule_set_result contains deep {set: "OWNERSHIP_STRUCTURE", name: "SHAREHOLDERS", result: "#(RuleResult.Ignored)" }

		Given url mockURL
		And path "/subscribe/state-events/" + profileID
		And retry until response.length == 2
		When method GET
		Then assert responseStatus == 200
		And match response contains deep expectedProfileCreatedStateEvent
		And match response contains deep expectedProfileChangedStateEvent