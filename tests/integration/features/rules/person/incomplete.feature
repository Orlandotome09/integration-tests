Feature: Incomplete Rule

	Background:
		* url baseURLCompliance
		
		* def IncompleteRuleOffer = "TEST_OFFER" + uuid() 
		* def IncompleteRuleOfferCompany = "TEST_OFFER" + uuid() 

		* def ProblemCodeDateOfBirthRequired = "DATE_OF_BIRTH_REQUIRED"
		* def ProblemCodePhoneRequired = "PHONE_REQUIRED"
		* def ProblemCodeEmailRequired = "EMAIL_REQUIRED"
		* def ProblemCodePepRequired = "PEP_REQUIRED"
		* def ProblemCodeLastNameRequired = "LAST_NAME_REQUIRED"
		* def ProblemCodeFileNotFoundIdentification = "FILE_NOT_FOUND_IDENTIFICATION"
		* def ProblemCodeDocumentNotFoundIdentification = "DOCUMENT_NOT_FOUND_IDENTIFICATION"
		* def ProblemCodeDocumentNotFoundRegistrationForm = "DOCUMENT_NOT_FOUND_REGISTRATION_FORM"
		* def ProblemCodeAddressNotFound = "ADDRESS_NOT_FOUND"
		* def ProblemCodeDocumentNotFoundStatuteSocial = "DOCUMENT_NOT_FOUND_STATUTE_SOCIAL"
		* def ProblemCodeDocumentNotFoundCorporateDocument = "DOCUMENT_NOT_FOUND_CORPORATE_DOCUMENT"

		* def offer = { offer_type : '#(IncompleteRuleOffer)', product: 'maquininha'}

		Given path '/offers'
		And header Content-Type = 'application/json'
		And request offer
		When method POST

		* def func = read('incompleteRuleSetConfig.js')
		* def ruleSetConfig = call func {}

		Given url mockURL
        And path '/temis-config/cadastral-validation-configs'
		And header Content-Type = 'application/json'
		And def catalog = call CreateSingleLevelCatalog { offer_type: '#(IncompleteRuleOffer)', role_type: '#(RoleType.Customer)', person_type: '#(ProfileType.Individual)', account_flag: false, rules_config: #(ruleSetConfig)}
		And request catalog
		When method POST

		* def offer = { offer_type : '#(IncompleteRuleOfferCompany)', product: 'maquininha'}

		Given path '/offers'
		And header Content-Type = 'application/json'
		And request offer
		When method POST

		* def createIncompleteRuleSetConfigCompany = read('incompleteRuleSetConfigCompany.js')
		* def ruleSetConfig = call createIncompleteRuleSetConfigCompany {}

		Given url mockURL
        And path '/temis-config/cadastral-validation-configs'
		And header Content-Type = 'application/json'
		And def catalog = call CreateSingleLevelCatalog { offer_type: '#(IncompleteRuleOfferCompany)', role_type: '#(RoleType.Merchant)', person_type: '#(ProfileType.Company)', account_flag: false, rules_config: #(ruleSetConfig)}
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
		
	Scenario: Should find incompleted profile

		* def profile = { profile_id: '#(profileID)', offer_type: '#(IncompleteRuleOffer)',  role_type: '#(RoleType.Customer)', profile_type: '#(ProfileType.Individual)' }

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
		And retry until response.result == RuleResult.Incomplete
		When method GET
		* def ruleResults = response.rule_set_result
		Then assert ruleResults[0].result == RuleResult.Incomplete
		Then assert ruleResults[0].set == RuleSet.Incomplete
		Then assert ruleResults[0].name == RuleName.RequiredFieldsNotFound
		Then assert ruleResults[0].problems[0].code == ProblemCodeDateOfBirthRequired
		Then assert ruleResults[0].problems[1].code == ProblemCodePhoneRequired
		Then assert ruleResults[0].problems[2].code == ProblemCodeEmailRequired
		Then assert ruleResults[0].problems[3].code == ProblemCodePepRequired
		Then assert ruleResults[0].problems[4].code == ProblemCodeLastNameRequired
		Then assert ruleResults[1].result == RuleResult.Incomplete
		Then assert ruleResults[1].set == RuleSet.Incomplete
		Then assert ruleResults[1].name == RuleName.AddressNotFound
		Then assert ruleResults[1].problems[0].code == ProblemCodeAddressNotFound
		Then assert ruleResults[2].result == RuleResult.Incomplete
		Then assert ruleResults[2].set == RuleSet.Incomplete
		Then assert ruleResults[2].name == RuleName.DocumentNotFound
		Then assert ruleResults[2].problems[0].code == ProblemCodeDocumentNotFoundIdentification
		Then assert ruleResults[2].problems[1].code == ProblemCodeDocumentNotFoundRegistrationForm

		Given url mockURL
		And path "/subscribe/state-events/" + profileID
		And retry until response.length == 2
		When method GET
		Then assert responseStatus == 200
		And match response contains deep expectedCreatedStateEvent
		And match response contains deep expectedChangedStateEvent

	Scenario: Should find completed profile

		* def profile = call CreateProfileIndividual { profile_id: '#(profileID)', offer_type: '#(IncompleteRuleOffer)', role_type: '#(RoleType.Customer)'}
		* profile.individual.pep = false
		* profile.expiration_date = "2024-01-01T00:00:00Z"

		Given url mockURL
		And path '/v1/temis/profile/' + profileID
		And request profile
		When method POST
		Then assert responseStatus == 200

		* def documentID1 = uuid()
		* def documentIdentification = { document_id: '#(documentID1)', entity_id: '#(profileID)', type: 'IDENTIFICATION', sub_type: 'RG'}
		* documentIdentification.expiration_date = "2024-02-02"
		* documentIdentification.emission_date = "2021-02-02T00:00:00Z"
		* def documentID2 = uuid()
		* def documentRegistrationForm = { document_id: '#(documentID2)', entity_id: '#(profileID)', type: 'REGISTRATION_FORM'}
		* documentRegistrationForm.expiration_date = "2024-03-03"
		* documentRegistrationForm.emission_date = "2021-03-03T00:00:00Z"
		* def documents = []
		* def void = documents.add(documentIdentification)
		* def void = documents.add(documentRegistrationForm)
		
		Given url mockURL
		And path '/v1/temis/documents'
		And params ({ entity_id: profileID })
		And request documents
		When method POST

		* def documentFileID1 = uuid()
		* def documentFiles1 = [{ document_file_id: '#(documentFileID1)', document_id: '#(documentID1)'}]

		Given url mockURL
		And path '/v1/temis/document/' + documentID1 + '/files'
		And request documentFiles1
		When method POST

		* def documentFileID2 = uuid()
		* def documentFiles2 = [{ document_file_id: '#(documentFileID2)', document_id: '#(documentID2)'}]

		Given url mockURL
		And path '/v1/temis/document/' + documentID2 + '/files'
		And request documentFiles2
		When method POST

		* def addresses = [{ type: 'RESIDENTIAL', street: 'Rua Roberto Menotti'}]

		Given url mockURL
		And path '/v1/temis/addresses'
		And params ({ profile_id: profileID })
		And request addresses
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
		Then assert response.result == RuleResult.Approved
		Then assert ruleResults[0].set == RuleSet.Incomplete
		Then assert ruleResults[0].name == RuleName.RequiredFieldsNotFound
		Then assert ruleResults[0].result == RuleResult.Approved
		Then assert ruleResults[0].pending == false
		Then assert ruleResults[0].metadata == null

		Then assert ruleResults[1].set == RuleSet.Incomplete
		Then assert ruleResults[1].name == RuleName.AddressNotFound
		Then assert ruleResults[1].result == RuleResult.Approved
		Then assert ruleResults[1].pending == false
		Then assert ruleResults[1].metadata == null

		Then assert ruleResults[2].set == RuleSet.Incomplete
		Then assert ruleResults[2].name == RuleName.DocumentNotFound
		Then assert ruleResults[2].result == RuleResult.Approved
		Then assert ruleResults[2].pending == false
		Then assert ruleResults[2].metadata == null

		Given url mockURL
		And path "/subscribe/state-events/" + profileID
		And retry until response.length == 2
		When method GET
		Then assert responseStatus == 200
		And match response contains deep expectedCreatedStateEvent
		And match response contains deep expectedChangedStateEvent

		Given url baseURLCompliance
		And path '/profile/' + profileID
		And header Content-Type = 'application/json'
		When method GET
		Then assert responseStatus == 200
		And assert response.person.documents[0].expiration_date == documentIdentification.expiration_date
		And assert response.person.documents[0].emission_date == documentIdentification.emission_date
		And assert response.person.documents[1].expiration_date == documentRegistrationForm.expiration_date
		And assert response.person.documents[1].emission_date == documentRegistrationForm.emission_date

	Scenario: Should approve publicly traded company

		* def documentNumber =  DocumentNormalizer(CNPJGenerator())
		* def partnerID = uuid()
		* def parentID = uuid()
		* def params =
        """
            {
                profile_id: '#(profileID)',
                partner_id: '#(partnerID)',
                parent_id: '#(parentID)',
                offer_type: '#(IncompleteRuleOfferCompany)',
                role_type: '#(RoleType.Merchant)',
                document_number: '#(documentNumber)',
                callback_url: '/url'
            }
        """
		* def profileCompany = call CreateProfileCompany params
		* profileCompany.company.legal_nature = "2046"

		Given url mockURL
		And path '/v1/temis/profile/' + profileID
		And request profileCompany
		When method POST
		Then assert responseStatus == 200

		* def documentID1 = uuid()
		* def documentCorporate = { document_id: '#(documentID1)', entity_id: '#(profileID)', type: 'CORPORATE_DOCUMENT'}
		* def documentID2 = uuid()
		* def documentBusiness = { document_id: '#(documentID2)', entity_id: '#(profileID)', type: 'APPOINTMENT_DOCUMENT', sub_type: "MINUTES_OF_ELECTION"}
		* def documentID3 = uuid()
		* def documentFinancial = { document_id: '#(documentID3)', entity_id: '#(profileID)', type: 'CONSTITUTION_DOCUMENT', sub_type: "STATUTE_SOCIAL"}
		* def documents = []
		* def void = documents.add(documentCorporate)
		* def void = documents.add(documentBusiness)
		* def void = documents.add(documentFinancial)
		
		Given url mockURL
		And path '/v1/temis/documents'
		And params ({ entity_id: profileID })
		And request documents
		When method POST

		* def documentFileID1 = uuid()
		* def documentFiles1 = [{ document_file_id: '#(documentFileID1)', document_id: '#(documentID1)'}]

		Given url mockURL
		And path '/v1/temis/document/' + documentID1 + '/files'
		And request documentFiles1
		When method POST

		* def documentFileID2 = uuid()
		* def documentFiles2 = [{ document_file_id: '#(documentFileID2)', document_id: '#(documentID2)'}]

		Given url mockURL
		And path '/v1/temis/document/' + documentID2 + '/files'
		And request documentFiles2
		When method POST

		* def documentFileID3 = uuid()
		* def documentFiles3 = [{ document_file_id: '#(documentFileID3)', document_id: '#(documentID3)'}]

		Given url mockURL
		And path '/v1/temis/document/' + documentID3 + '/files'
		And request documentFiles3
		When method POST

		* def addresses = [{ type: 'COMERCIAL', street: 'Rua Roberto Menotti'}]

		Given url mockURL
		And path '/v1/temis/addresses'
		And params ({ profile_id: profileID })
		And request addresses
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
		Then assert response.result == RuleResult.Approved
		Then assert ruleResults[0].set == RuleSet.Incomplete
		Then assert ruleResults[0].name == RuleName.AddressNotFound
		Then assert ruleResults[0].result == RuleResult.Approved
		Then assert ruleResults[0].pending == false
		Then assert ruleResults[0].metadata == null

		Then assert ruleResults[1].set == RuleSet.Incomplete
		Then assert ruleResults[1].name == RuleName.DocumentNotFound
		Then assert ruleResults[1].result == RuleResult.Approved
		Then assert ruleResults[1].pending == false
		Then assert ruleResults[1].metadata == null

		Given url mockURL
		And path "/subscribe/state-events/" + profileID
		And retry until response.length == 2
		When method GET
		Then assert responseStatus == 200
		And match response contains deep expectedCreatedStateEvent
		And match response contains deep expectedChangedStateEvent

	Scenario: Should not approve publicly traded company when it does not have specific documents
	
		* def documentNumber =  DocumentNormalizer(CNPJGenerator())
		* def partnerID = uuid()
		* def parentID = uuid()
		* def params =
        """
            {
                profile_id: '#(profileID)',
                partner_id: '#(partnerID)',
                parent_id: '#(parentID)',
                offer_type: '#(IncompleteRuleOfferCompany)',
                role_type: '#(RoleType.Merchant)',
                document_number: '#(documentNumber)',
                callback_url: '/url'
            }
        """
		* def profileCompany = call CreateProfileCompany params
		* profileCompany.company.legal_nature = "2046"

		Given url mockURL
		And path '/v1/temis/profile/' + profileID
		And request profileCompany
		When method POST
		Then assert responseStatus == 200

		* def documentID1 = uuid()
		* def documentCorporate = { document_id: '#(documentID1)', entity_id: '#(profileID)', type: 'CORPORATE_DOCUMENT'}
		* def documentID2 = uuid()
		* def documentBusiness = { document_id: '#(documentID2)', entity_id: '#(profileID)', type: 'APPOINTMENT_DOCUMENT', sub_type: "MINUTES_OF_ELECTION"}
		* def documents = []
		* def void = documents.add(documentCorporate)
		* def void = documents.add(documentBusiness)
		
		Given url mockURL
		And path '/v1/temis/documents'
		And params ({ entity_id: profileID })
		And request documents
		When method POST

		* def documentFileID1 = uuid()
		* def documentFiles1 = [{ document_file_id: '#(documentFileID1)', document_id: '#(documentID1)'}]

		Given url mockURL
		And path '/v1/temis/document/' + documentID1 + '/files'
		And request documentFiles1
		When method POST

		* def documentFileID2 = uuid()
		* def documentFiles2 = [{ document_file_id: '#(documentFileID2)', document_id: '#(documentID2)'}]

		Given url mockURL
		And path '/v1/temis/document/' + documentID2 + '/files'
		And request documentFiles2
		When method POST

		* def addresses = [{ type: 'COMERCIAL', street: 'Rua Roberto Menotti'}]

		Given url mockURL
		And path '/v1/temis/addresses'
		And params ({ profile_id: profileID })
		And request addresses
		When method POST

		* def event = { profile_id: '#(profileID)', entity_id: '#(profileID)', entity_type: 'PROFILE', event_type: 'PROFILE_CHANGED', parent_type: 'PROFILE', update_date: '2023-08-03T19:29:06.63556Z'  }
		* string json = event
    	
		* def result = RegistrationEventsPublisher.publish(json)

		Given url baseURLCompliance
		And path '/state/' + profileID
		And header Content-Type = 'application/json'
		And params ({ only_pending: false })
		And retry until response.result == RuleResult.Incomplete
		When method GET 
		* def ruleResults = response.rule_set_result
		Then assert response.result == RuleResult.Incomplete
		Then assert ruleResults[0].set == RuleSet.Incomplete
		Then assert ruleResults[0].name == RuleName.AddressNotFound
		Then assert ruleResults[0].result == RuleResult.Approved
		Then assert ruleResults[0].pending == false
		Then assert ruleResults[0].metadata == null

		Then assert ruleResults[1].set == RuleSet.Incomplete
		Then assert ruleResults[1].name == RuleName.DocumentNotFound
		Then assert ruleResults[1].result == RuleResult.Incomplete
		Then assert ruleResults[1].problems[0].code == ProblemCodeDocumentNotFoundStatuteSocial

		Given url mockURL
		And path "/subscribe/state-events/" + profileID
		And retry until response.length == 2
		When method GET
		Then assert responseStatus == 200
		And match response contains deep expectedCreatedStateEvent
		And match response contains deep expectedChangedStateEvent

	Scenario: Should not approve a publicly traded company when it does not have specific documents or the journey
		
		* def documentNumber =  DocumentNormalizer(CNPJGenerator())
		* def partnerID = uuid()
		* def parentID = uuid()
		* def params =
        """
            {
                profile_id: '#(profileID)',
                partner_id: '#(partnerID)',
                parent_id: '#(parentID)',
                offer_type: '#(IncompleteRuleOfferCompany)',
                role_type: '#(RoleType.Merchant)',
                document_number: '#(documentNumber)',
                callback_url: '/url'
            }
        """
		* def profileCompany = call CreateProfileCompany params
		* profileCompany.company.legal_nature = "2046"

		Given url mockURL
		And path '/v1/temis/profile/' + profileID
		And request profileCompany
		When method POST
		Then assert responseStatus == 200

		* def documentID1 = uuid()
		* def documentBusiness = { document_id: '#(documentID1)', entity_id: '#(profileID)', type: 'APPOINTMENT_DOCUMENT', sub_type: "MINUTES_OF_ELECTION"}
		* def documentID2 = uuid()
		* def documentFinancial = { document_id: '#(documentID2)', entity_id: '#(profileID)', type: 'CONSTITUTION_DOCUMENT', sub_type: "STATUTE_SOCIAL"}
		* def documents = []
		* def void = documents.add(documentBusiness)
		* def void = documents.add(documentFinancial)
		
		Given url mockURL
		And path '/v1/temis/documents'
		And params ({ entity_id: profileID })
		And request documents
		When method POST

		* def documentFileID1 = uuid()
		* def documentFiles1 = [{ document_file_id: '#(documentFileID1)', document_id: '#(documentID1)'}]

		Given url mockURL
		And path '/v1/temis/document/' + documentID1 + '/files'
		And request documentFiles1
		When method POST

		* def documentFileID2 = uuid()
		* def documentFiles2 = [{ document_file_id: '#(documentFileID2)', document_id: '#(documentID2)'}]

		Given url mockURL
		And path '/v1/temis/document/' + documentID2 + '/files'
		And request documentFiles2
		When method POST

		* def addresses = [{ type: 'COMERCIAL', street: 'Rua Roberto Menotti'}]

		Given url mockURL
		And path '/v1/temis/addresses'
		And params ({ profile_id: profileID })
		And request addresses
		When method POST

		* def event = { profile_id: '#(profileID)', entity_id: '#(profileID)', entity_type: 'PROFILE', event_type: 'PROFILE_CHANGED', parent_type: 'PROFILE', update_date: '2023-08-03T19:29:06.63556Z'  }
		* string json = event
    	
		* def result = RegistrationEventsPublisher.publish(json)

		Given url baseURLCompliance
		And path '/state/' + profileID
		And header Content-Type = 'application/json'
		And params ({ only_pending: false })
		And retry until response.result == RuleResult.Incomplete
		When method GET 
		* def ruleResults = response.rule_set_result
		Then assert response.result == RuleResult.Incomplete
		Then assert ruleResults[0].set == RuleSet.Incomplete
		Then assert ruleResults[0].name == RuleName.AddressNotFound
		Then assert ruleResults[0].result == RuleResult.Approved
		Then assert ruleResults[0].pending == false
		Then assert ruleResults[0].metadata == null

		Then assert ruleResults[1].set == RuleSet.Incomplete
		Then assert ruleResults[1].name == RuleName.DocumentNotFound
		Then assert ruleResults[1].result == RuleResult.Incomplete
		Then assert ruleResults[1].problems[0].code == ProblemCodeDocumentNotFoundCorporateDocument

		Given url mockURL
		And path "/subscribe/state-events/" + profileID
		And retry until response.length == 2
		When method GET
		Then assert responseStatus == 200
		And match response contains deep expectedCreatedStateEvent
		And match response contains deep expectedChangedStateEvent

	Scenario: Should approve a company that is not publicly traded when it does not have specific documents but has journey documents
		
		* def documentNumber =  DocumentNormalizer(CNPJGenerator())
		* def partnerID = uuid()
		* def parentID = uuid()
		* def params =
        """
            {
                profile_id: '#(profileID)',
                partner_id: '#(partnerID)',
                parent_id: '#(parentID)',
                offer_type: '#(IncompleteRuleOfferCompany)',
                role_type: '#(RoleType.Merchant)',
                document_number: '#(documentNumber)',
                callback_url: '/url'
            }
        """
		* def profileCompany = call CreateProfileCompany params
		# SA Code for privately held company = 2054
		* profileCompany.company.legal_nature = "2054" 

		Given url mockURL
		And path '/v1/temis/profile/' + profileID
		And request profileCompany
		When method POST
		Then assert responseStatus == 200

		* def documentID1 = uuid()
		* def documentCorporate = { document_id: '#(documentID1)', entity_id: '#(profileID)', type: 'CORPORATE_DOCUMENT'}
		* def documents = []
		* def void = documents.add(documentCorporate)
				
		Given url mockURL
		And path '/v1/temis/documents'
		And params ({ entity_id: profileID })
		And request documents
		When method POST

		* def documentFileID1 = uuid()
		* def documentFiles1 = [{ document_file_id: '#(documentFileID1)', document_id: '#(documentID1)'}]

		Given url mockURL
		And path '/v1/temis/document/' + documentID1 + '/files'
		And request documentFiles1
		When method POST

		* def addresses = [{ type: 'COMERCIAL', street: 'Rua Roberto Menotti'}]

		Given url mockURL
		And path '/v1/temis/addresses'
		And params ({ profile_id: profileID })
		And request addresses
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
		Then assert response.result == RuleResult.Approved
		Then assert ruleResults[0].set == RuleSet.Incomplete
		Then assert ruleResults[0].name == RuleName.AddressNotFound
		Then assert ruleResults[0].result == RuleResult.Approved
		Then assert ruleResults[0].pending == false
		Then assert ruleResults[0].metadata == null

		Then assert ruleResults[1].set == RuleSet.Incomplete
		Then assert ruleResults[1].name == RuleName.DocumentNotFound
		Then assert ruleResults[1].result == RuleResult.Approved
		Then assert ruleResults[1].pending == false
		Then assert ruleResults[1].metadata == null

		Given url mockURL
		And path "/subscribe/state-events/" + profileID
		And retry until response.length == 2
		When method GET
		Then assert responseStatus == 200
		And match response contains deep expectedCreatedStateEvent
		And match response contains deep expectedChangedStateEvent

	Scenario: Should not approve a company that is not publicly traded when it does not have specific documents or journey documents
		
		* def documentNumber =  DocumentNormalizer(CNPJGenerator())
		* def partnerID = uuid()
		* def parentID = uuid()
		* def params =
        """
            {
                profile_id: '#(profileID)',
                partner_id: '#(partnerID)',
                parent_id: '#(parentID)',
                offer_type: '#(IncompleteRuleOfferCompany)',
                role_type: '#(RoleType.Merchant)',
                document_number: '#(documentNumber)',
                callback_url: '/url'
            }
        """
		* def profileCompany = call CreateProfileCompany params
		# SA Code for privately held company = 2054
		* profileCompany.company.legal_nature = "2054" 

		Given url mockURL
		And path '/v1/temis/profile/' + profileID
		And request profileCompany
		When method POST
		Then assert responseStatus == 200
		
		* def addresses = [{ type: 'COMERCIAL', street: 'Rua Roberto Menotti'}]

		Given url mockURL
		And path '/v1/temis/addresses'
		And params ({ profile_id: profileID })
		And request addresses
		When method POST

		* def event = { profile_id: '#(profileID)', entity_id: '#(profileID)', entity_type: 'PROFILE', event_type: 'PROFILE_CHANGED', parent_type: 'PROFILE', update_date: '2023-08-03T19:29:06.63556Z'  }
		* string json = event
    	
		* def result = RegistrationEventsPublisher.publish(json)

		Given url baseURLCompliance
		And path '/state/' + profileID
		And header Content-Type = 'application/json'
		And params ({ only_pending: false })
		And retry until response.result == RuleResult.Incomplete
		When method GET 
		* def ruleResults = response.rule_set_result
		Then assert response.result == RuleResult.Incomplete
		Then assert ruleResults[0].set == RuleSet.Incomplete
		Then assert ruleResults[0].name == RuleName.AddressNotFound
		Then assert ruleResults[0].result == RuleResult.Approved
		Then assert ruleResults[0].pending == false
		Then assert ruleResults[0].metadata == null

		Then assert ruleResults[1].set == RuleSet.Incomplete
		Then assert ruleResults[1].name == RuleName.DocumentNotFound
		Then assert ruleResults[1].result == RuleResult.Incomplete
		Then assert ruleResults[1].problems[0].code == ProblemCodeDocumentNotFoundCorporateDocument

		Given url mockURL
		And path "/subscribe/state-events/" + profileID
		And retry until response.length == 2
		When method GET
		Then assert responseStatus == 200
		And match response contains deep expectedCreatedStateEvent
		And match response contains deep expectedChangedStateEvent
