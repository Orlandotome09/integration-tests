Feature: Send State Event

    Background:
        * url baseURLCompliance
        * def BureauRuleOffer = "TEST_OFFER" + uuid()
        * def DirectorRuleOffer = "TEST_OFFER" + uuid()


        * def problemCodeBureauStatusNotRegular = "BUREAU_STATUS_NOT_REGULAR"
        * def bureauStatusPendingRegularization = "PENDING_REGULARIZATION"
        * def ruleSetIncompleteContract = "INCOMPLETE_CONTRACT"
        * def ruleSetPreconditionContract = "PRECONDITION_CONTRACT_SET"

		* def ruleNameDocumentNotFound = "DOCUMENT_NOT_FOUND"
		* def ruleNameFileNotFound = "FILE_NOT_FOUND"
        * def ruleNameProfileApproved = "PROFILE_APPROVED"

        * def offer = { offer_type : '#(BureauRuleOffer)', product: 'maquininha' }

        Given path '/offers'
        And header Content-Type = 'application/json'
        And request offer
        When method POST

        * def offerDirector = { offer_type : '#(DirectorRuleOffer)', product: 'maquininha' }

        Given path '/offers'
        And header Content-Type = 'application/json'
        And request offerDirector
        When method POST

        * def ruleSetConfig = { bureau: {} }

        Given url mockURL
        And path '/temis-config/cadastral-validation-configs'
        And header Content-Type = 'application/json'
        And def catalog = call CreateSingleLevelCatalog { offer_type: '#(BureauRuleOffer)', role_type: 'CUSTOMER', person_type: '#(ProfileType.Individual)', account_flag: false, rules_config: #(ruleSetConfig)}
        And request catalog
        When method POST

        * def createdCatalogIndividual = response

        Given url mockURL
        And path '/temis-config/cadastral-validation-configs'
        And header Content-Type = 'application/json'
        And def catalog = call CreateSingleLevelCatalog { offer_type: '#(BureauRuleOffer)', role_type: 'CUSTOMER', person_type: '#(ProfileType.Company)', account_flag: false, rules_config: #(ruleSetConfig)}
        And request catalog
        When method POST

        * def createdCatalogCompany = response

        * def ruleSetConfig = { board_of_directors: {} }

        Given url mockURL
        And path '/temis-config/cadastral-validation-configs'
        And header Content-Type = 'application/json'
        And def catalog = call CreateSingleLevelCatalog { offer_type: '#(DirectorRuleOffer)', role_type: '#(RoleType.Customer)', person_type: '#(ProfileType.Company)', account_flag: false, rules_config: #(ruleSetConfig) }
        And catalog.product_config.enrich_profile_with_bureau_data = true
        And request catalog
        When method POST	

        * def ruleSetConfig = { bureau: {} }

        Given url mockURL
        And path '/temis-config/cadastral-validation-configs'
        And header Content-Type = 'application/json'
        And def catalog = call CreateSingleLevelCatalog { offer_type: '#(DirectorRuleOffer)', role_type: '#(RoleType.Director)', person_type: '#(ProfileType.Individual)', account_flag: false, rules_config: #(ruleSetConfig)}
        And request catalog
        When method POST

        * def createdCatalogDirector = response

        * def profileID = uuid()

    Scenario: should send profile individual state event when approved in bureau

        * def documentNumberIndividual = DocumentNormalizer(CPFGenerator())
        * def dateOfBirth = "2023-03-15T15:06:56.12Z"
        * def profileResponse =
        """
        {
            "profile_id": "#(profileID)",
            "partner_id": "#(uuid())",
            "document_number": "#(documentNumberIndividual)",
            "parent_id": "#(uuid())",
            "name": "John Doe",
            "legacy_id": "#(uuid())",
            "offer_type": "#(BureauRuleOffer)",
            "role_type":"#(RoleType.Customer)",
            "profile_type": "#(ProfileType.Individual)",
            "callback_url": "/callback",
            "individual":{
                "first_name": "John",
                "last_name": "Doe",
                "date_of_birth": "#(dateOfBirth)",
                "date_of_birth_inputted": "#(dateOfBirth)",
                "phones": [
                    {
                        "type" : "RESIDENTIAL",
                        "area_code": "13",
                        "country_code": "55",
                        "number": "888889999"
                    }
                ],
                "bureau_information": {
                    "name": "John Stewart Doe",
                    "date_of_birth": "#(dateOfBirth)"
                },
                "income": 10000,
                "income_currency": "BRL",
                "assets": 100000,
                "assets_currency": "BRL",
                "nationality": "BRA",
                "email": "some@mail.com",
                "pep": false,
                "us_person": false,
                "occupation": "Enginer",
                "foreign_tax_residency": false,
                "country_of_tax_residency": "BRA"
            },
            "email": "Some@mail.com",
            "created_at": "#(funcNow(2))",
            "updated_at": "#(funcNow(2))" 
        }    
        """

        Given url mockURL
        And path '/v1/temis/profile/' + profileID
        And request profileResponse
        When method POST
        Then assert responseStatus == 200

        * def enrichedBirthDate = "23/11/1992"
        * def enrichmentResponse = { situation: 1, name: "Johnatan Stewart Doodle", birth_date: "#(enrichedBirthDate)" }

		Given url mockURL
		And path '/temis-enrichment/individual/' + documentNumberIndividual
		And request enrichmentResponse
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
		Then assert responseStatus == 200
		Then assert response.result == RuleResult.Approved
		Then assert ruleResults[0].name == RuleName.CustomerNotFoundInSerasa
		Then assert ruleResults[0].result == RuleResult.Approved
		Then assert ruleResults[1].name == RuleName.CustomerHasProblemsInSerasa
		Then assert ruleResults[1].result == RuleResult.Approved

        * def expectedFullProfileSnapshot =
        """
        {
            "person":{
                "document_number": "#(profileResponse.document_number)",
                "name": "#(profileResponse.name)",
                "person_type": "#(profileResponse.profile_type)",
                "email": "#(profileResponse.email)",
                "partner_id": "#(profileResponse.partner_id)",
                "offer_type": "#(profileResponse.offer_type)",
                "profile_id": "#(profileID)",
                "entity_id": "#(profileResponse.profile_id)",
                "entity_type": "#(EntityType.Profile)",
                "role_type":"#(profileResponse.role_type)",
                "individual": {
                    "first_name": "John",
                    "last_name": "Doe",
                    "date_of_birth": "#string",
                    "date_of_birth_inputted": "#string",
                    "phones": [
                        {
                            "type":"#(profileResponse.individual.phones[0].type)",
                            "country_code":"#(profileResponse.individual.phones[0].country_code)",
                            "area_code":"#(profileResponse.individual.phones[0].area_code)",
                            "number":"#(profileResponse.individual.phones[0].number)"
                        }
                    ],
                    "bureau_information": {
                        "name":"#(profileResponse.individual.bureau_information.name)",
                        "date_of_birth":"#(profileResponse.individual.bureau_information.date_of_birth)"
                    },
                    "income": "#(profileResponse.individual.income)",
                    "income_currency": "#(profileResponse.individual.income_currency)",
                    "assets": "#(profileResponse.individual.assets)",
                    "assets_currency":"#(profileResponse.individual.assets_currency)",
                    "nationality":"#(profileResponse.individual.nationality)",
                    "pep":"#(profileResponse.individual.pep)",
                    "us_person": "#(profileResponse.individual.us_person)",
                    "occupation" :"#(profileResponse.individual.occupation)",
                    "foreign_tax_residency": "#(profileResponse.individual.foreign_tax_residency)",
                    "country_of_tax_residency": "#(profileResponse.individual.country_of_tax_residency)",
                    "contacts": "##string"
                },
                "company": "##object",
                "enriched_information": {
                    "status": "REGULAR",
                    "name": "#(enrichmentResponse.name)",
                    "birth_date": "#(enrichmentResponse.birth_date)",
                    "providers": null
                },
                "blacklist_status": "##object",
                "watchlist": "##object",
                "addresses": "##[]",
                "contacts":"##[]",
                "notification_recipients": "##[]",
                "documents":"##[]",
                "document_files":"##[]",
                "overrides":"##[]",
                "validation_steps":"##[]",
                "cadastral_validation_config": #(createdCatalogIndividual)
            },
            "profile_id": "#(profileResponse.profile_id)",
            "parent_id": "##(profileResponse.parent_id)",
            "legacy_id": "#(profileResponse.legacy_id)",
            "callback_url":"#(profileResponse.callback_url)",
            "legal_representatives": "##[]",
            "ownership_structure": "##object",
            "board_of_directors": "##[]",
            "created_at": "#string",
            "updated_at": "#string"
        }
        """

        Given url baseURLCompliance
        And path "/profile/" + profileID
        When method GET
        Then match response == expectedFullProfileSnapshot
        

        * def expectedCreatedStateEvent =
        """
        {
            "id": "#string",
			"entity_id": "#(profileID)",
			"entity_type": "#(EntityType.ComplianceState)",
			"event_type": "#(EventType.State.Created)",
			"update_date":"#string",
			"data": {
                "profile_id": "#(profileID)",
                "content": {
                    "state":{
                        "entity_id":"#(profileID)",
                        "engine_name":"#(Engine.Profile)",
                        "result":"#(RuleResult.Created)",
                        "validation_steps_results":"##array",
                        "rule_names": "##array",
                        "pending": false,
                        "request_date": "##string",
                        "execution_time": "##string",
                        "version": "#number",
                        "created_at": "##string",
                        "updated_at": "##string"
                    },
                    "profile": "#(expectedFullProfileSnapshot)",
                    "person": "##null",
                    "contract": "##null"
                },
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
                "profile_id": "#(profileID)",
				"content": {
                    "state":{
                        "entity_id":"#(profileID)",
                        "engine_name":"#(Engine.Profile)",
                        "result":"#(RuleResult.Approved)",
                        "validation_steps_results": [
                            {
                                "step_number": 0,
                                "result": "#(RuleResult.Approved)",
                                "skip_for_approval": false,
                                "rule_results": [
                                    {
                                        "rule_set": "#(RuleSet.Bureau)",
                                        "rule_name": "#(RuleName.CustomerNotFoundInSerasa)",
                                        "result": "#(RuleResult.Approved)",
                                        "pending": false,
                                        "expire_at": null,
                                        "metadata": null,
                                        "tags": null,
                                        "problems": [],
                                        "approved_documents": null
                                    },
                                    {
                                        "rule_set": "#(RuleSet.Bureau)",
                                        "rule_name": "#(RuleName.CustomerHasProblemsInSerasa)",
                                        "result": "#(RuleResult.Approved)",
                                        "pending": false,
                                        "expire_at": null,
                                        "metadata": null,
                                        "tags": null,
                                        "problems": [],
                                        "approved_documents": null
                                    }
                                ]
                            } 
                        ],
                        "rule_names": ["#(RuleName.CustomerNotFoundInSerasa)", "#(RuleName.CustomerHasProblemsInSerasa)"],
                        "pending": false,
                        "request_date": "##string",
                        "execution_time": "##string",
                        "version": "#number",
                        "created_at": "##string",
                        "updated_at": "##string"
                    },
                    "profile": "#(expectedFullProfileSnapshot)",
                    "person": "##null",
                    "contract": "##null"
                },
			}
        }
        """

		Given url mockURL
		And path "/subscribe/state-events/" + profileID
		And retry until response.length == 2
		When method GET
		Then assert responseStatus == 200
		And match response[0] == expectedCreatedStateEvent
		And match response[1] == expectedChangedStateEvent

    Scenario: should send profile company state event when approved in bureau

        * def documentNumberCompany = DocumentNormalizer(CPFGenerator())
        * def dateOfIncorportaion = "1989-02-13"
        * def profileResponse =
        """
        {
            "profile_id": "#(profileID)",
            "partner_id": "#(uuid())",
            "document_number": "#(documentNumberCompany)",
            "parent_id": "#(uuid())",
            "name": "John Doe",
            "legacy_id": "#(uuid())",
            "offer_type": "#(BureauRuleOffer)",
            "role_type":"#(RoleType.Customer)",
            "profile_type": "#(ProfileType.Company)",
            "callback_url": "/callback",
            "company":{
                "legal_name": "Business A",
                "business_name": "Business A",
                "tax_payer_identification": "#(uuid())",
                "company_registration_number":"#(uuid())",
                "date_of_incorporation": "#(dateOfIncorportaion)",
                "place_of_incorporation": "BRA",
                "share_capital": {
                    "amount": 120000,
                    "currency": "BRA"
                },
                "license": "#(uuid())",
                "website": "bussinesA.com",
                "goods_delivery":{
                    "average_delivery_days": 4,
                    "shipping_methods": "corrier",
                    "insurance": false,
                    "tracking_code_available": true
                },
                "assets": 120000,
                "assets_currency": "BRL",
                "annual_income": 1000000,
                "annual_income_currency": "BRL",
                "country_code": "BRA",
                "legal_nature": "2000"
                
            },
            "email": "Some@mail.com",
            "created_at": "#(funcNow(2))",
            "updated_at": "#(funcNow(2))" 
        }    
        """

        Given url mockURL
        And path '/v1/temis/profile/' + profileID
        And request profileResponse
        When method POST
        Then assert responseStatus == 200

        * def openingDate = "22/12/2010"
        * def enrichmentResponse = 
        """
        { 
            "legal_name": "Bussines AAA",
            "opening_date": "#(openingDate)",
            "situation": 2,
            "cnae": "6438701",
            "legal_nature": "2054",
        } 
        """

		Given url mockURL
		And path '/temis-enrichment/legal-entity/' + documentNumberCompany
		And request enrichmentResponse
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
		Then assert responseStatus == 200
		Then assert response.result == RuleResult.Approved
		Then assert ruleResults[0].name == RuleName.CustomerNotFoundInSerasa
		Then assert ruleResults[0].result == RuleResult.Approved
		Then assert ruleResults[1].name == RuleName.CustomerHasProblemsInSerasa
		Then assert ruleResults[1].result == RuleResult.Approved

        * def expectedFullProfileSnapshot =
        """
        {
            "person":{
                "document_number": "#(profileResponse.document_number)",
                "name": "#(profileResponse.name)",
                "person_type": "#(profileResponse.profile_type)",
                "email": "#(profileResponse.email)",
                "partner_id": "#(profileResponse.partner_id)",
                "offer_type": "#(profileResponse.offer_type)",
                "profile_id": "#(profileID)",
                "entity_id": "#(profileResponse.profile_id)",
                "entity_type": "#(EntityType.Profile)",
                "role_type":"#(profileResponse.role_type)",
                "individual": "##object",
                "company": {
                    "legal_name": "#(profileResponse.company.legal_name)",
                    "business_name" : "#(profileResponse.company.business_name)",
                    "tax_payer_identification":"#(profileResponse.company.tax_payer_identification)",
                    "company_registration_number": "#(profileResponse.company.company_registration_number)",
                    "date_of_incorporation":"#string",
                    "place_of_incorporation": "#(profileResponse.company.place_of_incorporation)",
                    "share_capital": {
                        "amount":"#(profileResponse.company.share_capital.amount)",
                        "currency":"#(profileResponse.company.share_capital.currency)"
                    },
                    "license":"#(profileResponse.company.license)",
                    "website" :"#(profileResponse.company.website)",
                    "goods_delivery": {
                        "average_delivery_days": "#(profileResponse.company.goods_delivery.average_delivery_days)",
                        "shipping_methods": "#(profileResponse.company.goods_delivery.shipping_methods)",
                        "insurance": "#(profileResponse.company.goods_delivery.insurance)",
                        "tracking_code_available":"#(profileResponse.company.goods_delivery.tracking_code_available)"
                    },
                    "assets":"#(profileResponse.company.assets)",
                    "assets_currency":"#(profileResponse.company.assets_currency)",
                    "annual_income":"#(profileResponse.company.annual_income)",
                    "annual_income_currency": "#(profileResponse.company.annual_income_currency)",
                    "country_code": "#(profileResponse.company.country_code)",
                    "legal_nature": "#(profileResponse.company.legal_nature)"
                },
                "enriched_information": {
                    "status": "REGULAR",
                    "legal_name": "#(enrichmentResponse.legal_name)",
                    "economic_activity": "#(enrichmentResponse.cnae)",
                    "opening_date": "#(enrichmentResponse.opening_date)",
                    "legal_nature": "#(enrichmentResponse.legal_nature)",
                    "providers": null
                },
                "blacklist_status": "##object",
                "watchlist": "##object",
                "addresses": "##[]",
                "contacts":"##[]",
                "notification_recipients": "##[]",
                "documents":"##[]",
                "document_files":"##[]",
                "overrides":"##[]",
                "validation_steps":"##[]",
                "cadastral_validation_config": #(createdCatalogCompany)
            },
            "profile_id": "#(profileResponse.profile_id)",
            "parent_id": "##(profileResponse.parent_id)",
            "legacy_id": "#(profileResponse.legacy_id)",
            "callback_url":"#(profileResponse.callback_url)",
            "legal_representatives": "##[]",
            "ownership_structure": "##object",
            "board_of_directors": "##[]",
            "created_at": "#string",
            "updated_at": "#string"
        }
        """

        Given url baseURLCompliance
        And path "/profile/" + profileID
        When method GET
        Then match response == expectedFullProfileSnapshot
        

        * def expectedCreatedStateEvent =
        """
        {
            "id": "#string",
			"entity_id": "#(profileID)",
			"entity_type": "#(EntityType.ComplianceState)",
			"event_type": "#(EventType.State.Created)",
			"update_date":"#string",
			"data": {
                "profile_id": "#(profileID)",
                "content": {
                    "state":{
                        "entity_id":"#(profileID)",
                        "engine_name":"#(Engine.Profile)",
                        "result":"#(RuleResult.Created)",
                        "validation_steps_results":"##array",
                        "rule_names": "##array",
                        "pending": false,
                        "request_date": "##string",
                        "execution_time": "##string",
                        "version": "#number",
                        "created_at": "##string",
                        "updated_at": "##string"
                    },
                    "profile": "#(expectedFullProfileSnapshot)",
                    "person": "##null",
                    "contract": "##null"
                },
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
                "profile_id": "#(profileID)",
				"content": {
                    "state":{
                        "entity_id":"#(profileID)",
                        "engine_name":"#(Engine.Profile)",
                        "result":"#(RuleResult.Approved)",
                        "validation_steps_results": [
                            {
                                "step_number": 0,
                                "result": "#(RuleResult.Approved)",
                                "skip_for_approval": false,
                                "rule_results": [
                                    {
                                        "rule_set": "#(RuleSet.Bureau)",
                                        "rule_name": "#(RuleName.CustomerNotFoundInSerasa)",
                                        "result": "#(RuleResult.Approved)",
                                        "pending": false,
                                        "expire_at": null,
                                        "metadata": null,
                                        "tags": null,
                                        "problems": [],
                                        "approved_documents": null
                                    },
                                    {
                                        "rule_set": "#(RuleSet.Bureau)",
                                        "rule_name": "#(RuleName.CustomerHasProblemsInSerasa)",
                                        "result": "#(RuleResult.Approved)",
                                        "pending": false,
                                        "expire_at": null,
                                        "metadata": null,
                                        "tags": null,
                                        "problems": [],
                                        "approved_documents": null
                                    }
                                ]
                            } 
                        ],
                        "rule_names": ["#(RuleName.CustomerNotFoundInSerasa)", "#(RuleName.CustomerHasProblemsInSerasa)"],
                        "pending": false,
                        "request_date": "##string",
                        "execution_time": "##string",
                        "version": "#number",
                        "created_at": "##string",
                        "updated_at": "##string"
                    },
                    "profile": "#(expectedFullProfileSnapshot)",
                    "person": "##null",
                    "contract": "##null"
                },
			}
        }
        """

		Given url mockURL
		And path "/subscribe/state-events/" + profileID
		And retry until response.length == 2
		When method GET
		Then assert responseStatus == 200
		And match response[0] == expectedCreatedStateEvent
		And match response[1] == expectedChangedStateEvent


    Scenario: Should send contract state event when documents are updated

        * def contractID = uuid()
		* def documentID = uuid()
		* def profileID = uuid()
		* def contract = { contract_id: '#(contractID)', profile_id: '#(profileID)', document_id: '#(documentID)'}

        * def expectedContractSnapshot = 
        """
            {
                "contract_id": "#(contractID)",
                "estimated_total_amount": 0,
                "due_time": "##string",
                "installments": 0,
                "correlation_id": "##string",
                "profile_id": "#(profileID)",
                "document_id": "#(documentID)",
                "created_at": "##string",
                "updated_at": "##string"
            }
        """

		Given url mockURL
		And path '/v1/temis/contract/' + contractID
		And request contract
		When method POST

        * def documentResponse = { body: { document_id: '#(documentID)', entity_id: '#(profileID)'}, status: 200}

		Given url mockURL
		And path '/v1/temis/document/' + documentID
		And request documentResponse
		When method POST

		* def documentFileID = uuid()
		* def fileID = uuid()
		* def documentFilesResponse = [{document_file_id: '#(documentFileID)', document_id:'#(documentID)', file_id:'#(fileID)'}]

		Given url mockURL
		And path '/v1/temis/document/' + documentID + '/files'
		And request documentFilesResponse
		When method POST

        * def event = { contract_id: '#(contractID)', entity_id: '#(contractID)', entity_type: 'CONTRACT', event_type: 'CONTRACT_CREATED', parent_type: 'CONTRACT', update_date: '2023-08-03T19:29:06.63556Z'  }
		* string json = event
    	* def result = RegistrationEventsPublisher.publish(json)
		* eval sleep(1000)

		Given url baseURLCompliance
		And path '/state/' + contractID
        And header Content-Type = 'application/json'
        And params ({ only_pending: false })		
		And retry until response.result == RuleResult.Approved
        When method GET
		Then assert responseStatus == 200
		And match response.entity_id == contractID

        * def expectedCreatedStateEvent =
        """
        {
            "id": "#string",
			"entity_id": "#(contractID)",
			"entity_type": "#(EntityType.ComplianceState)",
			"event_type": "#(EventType.State.Created)",
			"update_date":"#string",
			"data": {
                "profile_id": "#(profileID)",
                "content": {
                    "state":{
                        "entity_id":"#(contractID)",
                        "engine_name":"#(Engine.Contract)",
                        "result":"#(RuleResult.Created)",
                        "validation_steps_results":"##array",
                        "rule_names": "##array",
                        "pending": false,
                        "request_date": "##string",
                        "execution_time": "##string",
                        "version": "#number",
                        "created_at": "##string",
                        "updated_at": "##string"
                    },
                    "profile": "##null",
                    "person": "##null",
                    "contract": "#(expectedContractSnapshot)"
                },
			}
        }
        """

        * def expectedChangedStateEvent =
        """
        {
            "id": "#string",
			"entity_id": "#(contractID)",
			"entity_type": "#(EntityType.ComplianceState)",
			"event_type": "#(EventType.State.Changed)",
			"update_date":"#string",
			"data": {
                "profile_id": "#(profileID)",
                "content": {
                    "state":{
                        "entity_id":"#(contractID)",
                        "engine_name":"#(Engine.Contract)",
                        "result":"#(RuleResult.Approved)",
                        "rule_names": "##array",
                        "pending": false,
                        "request_date": "##string",
                        "execution_time": "##string",
                        "version": "#number",
                        "created_at": "##string",
                        "updated_at": "##string"
                    },
                    "profile": "##null",
                    "person": "##null",
                    "contract": "#(expectedContractSnapshot)"
                },
			}
        }
        """

		Given url mockURL
		And path "/subscribe/state-events/" + contractID
		And retry until response.length == 2
		When method GET
		Then assert responseStatus == 200
		And match response[0] == expectedCreatedStateEvent
        And match response[1] contains deep expectedChangedStateEvent


    Scenario: should send director state event when approved in bureau
        * def documentNumberCompany = DocumentNormalizer(CNPJGenerator())
        * def partnerID = uuid()
        * def profileResponse = {profile_id:'#(profileID)', document_number: '#(documentNumberCompany)', offer_type: '#(DirectorRuleOffer)', partner_id: '#(partnerID)', role_type: '#(RoleType.Customer)', profile_type: '#(ProfileType.Company)', company: {legal_nature: "2143"}}

        Given url mockURL
        And path '/v1/temis/profile/' + profileID
        And request profileResponse
        When method POST
        Then assert responseStatus == 200

        * def directorID = uuid()
        * def directorDocumentNumber = CPFGenerator()
        * def directors =  [{director_id : '#(directorID)',document_number: '#(directorDocumentNumber)', profile_id: '#(profileID)', partner_id: '#(profileResponse.partner_id)', full_name: "Nome Completo do Diretor"}]

        Given url mockURL
        And path '/v1/temis/directors'
        And params { profile_id : '#(profileID)' }
        And request directors
        When method POST

        * def enrichmentResponse = { situation: 1 }

        Given url mockURL
        And path '/temis-enrichment/individual/' + directorDocumentNumber
        And request enrichmentResponse
        When method POST

        * def event = { profile_id: '#(profileID)', entity_id: '#(profileID)', entity_type: 'PROFILE', event_type: 'PROFILE_CHANGED', parent_type: 'PROFILE', update_date: '2023-08-03T19:29:06.63556Z'  }
        * string json = event
        
        * def result = RegistrationEventsPublisher.publish(json)

        Given url baseURLCompliance
        And path '/state/' + profileID
        And header Content-Type = 'application/json'
        And retry until response.result == RuleResult.Approved
        When method GET
        And match response.rule_set_result contains deep { set:"#(RuleSet.BoardOfDirectors)", name:"#(RuleName.BoardOfDirectorsComplete)",result:"#(RuleResult.Approved)", pending: false }
        And match response.rule_set_result contains deep { set:"#(RuleSet.BoardOfDirectors)", name:"#(RuleName.BoardOfDirectorsResult)",result:"#(RuleResult.Approved)", pending: false }

        * def expectedBoardOfDirectorProfileSnapshot =
        """
        {
            "board_of_directors": [{
                "director_id": "#(directorID)",
                "role": "",
                "person": {
                    "document_number": "#(directorDocumentNumber)",
                    "name": "#(directors[0].full_name)",
                    "person_type": "#(ProfileType.Individual)",
                    "partner_id": "#(partnerID)",
                    "offer_type": "#(DirectorRuleOffer)",
                    "profile_id": "#(profileID)",
                    "entity_id": "#(directorID)",
                    "entity_type": "#(EntityType.Director)",
                    "role_type": "#(RoleType.Director)",
                    "individual": {
                        "foreign_tax_residency": false
                    },
                    "enriched_information": {
                        "status": "REGULAR",
                        "providers": null
                    },
                    "cadastral_validation_config": #(createdCatalogDirector)
                }
            }]
        }
        """

        Given url baseURLCompliance
        And path "/profile/" + profileID
        When method GET
        Then match response contains deep expectedBoardOfDirectorProfileSnapshot

        * def expectedDirectorCreatedStateEvent = 
        """
        {
            "id": "#string",
            "entity_id": "#(directorID)",
            "entity_type": "#(EntityType.ComplianceState)",
            "event_type": "#(EventType.State.Created)",
            "update_date":"#string",
            "data": {
                "profile_id": "#(profileID)",
                "content": {
                    "person": {
                        "document_number": "#(directorDocumentNumber)",
                        "person_type": "#(ProfileType.Individual)",
                        "offer_type": "#(DirectorRuleOffer)",
                        "partner_id": "#(partnerID)",
                        "profile_id": "#(profileID)",
                        "entity_id": "#(directorID)",
                        "entity_type": "#(EntityType.Director)",
                        "role_type": "#(RoleType.Director)"
                    },
                    "state": "##object",
                    "profile": "##object",
                    "contract": "##object"
                }
            }
        }	
        """

        * def expectedDirectorChangedStateEvent = 
        """
        {
            "id": "#string",
            "entity_id": "#(directorID)",
            "entity_type": "#(EntityType.ComplianceState)",
            "event_type": "#(EventType.State.Changed)",
            "update_date":"#string",
            "data": {
                "profile_id": "#(profileID)",
                "content": {
                    "person": {
                        "document_number": "#(directorDocumentNumber)",
                        "person_type": "#(ProfileType.Individual)",
                        "offer_type": "#(DirectorRuleOffer)",
                        "partner_id": "#(partnerID)",
                        "profile_id": "#(profileID)",
                        "entity_id": "#(directorID)",
                        "entity_type": "#(EntityType.Director)",
                        "role_type": "#(RoleType.Director)"
                    },
                    "state": "##object",
                    "profile": "##object",
                    "contract": "##object"
                }
            }
        }	
        """
        
        Given url mockURL
		And path "/subscribe/state-events/" + directorID
		And retry until response.length == 2
		When method GET
		Then assert responseStatus == 200
		And match response contains deep expectedDirectorCreatedStateEvent
		And match response contains deep expectedDirectorChangedStateEvent


    Scenario: should not resync an state entity id when the id is invalid
        * def entityId = "123"
        * def entityIds = { "ids": ['#(entityId)'] }
    
        Given url baseURLCompliance
        And path '/states/resync'
        And request entityIds
        When method POST
        Then assert responseStatus == 400
    
    Scenario: should not resync an state entity id when not id is provide
        * def entityIds = { "ids": [] }

        Given url baseURLCompliance
        And path '/states/resync'
        And request entityIds
        When method POST
        Then assert responseStatus == 200
        And match response.resynced == null
    
    Scenario: should resync an state entity id when id is valid
        * def createdProfileState1 = call read('classpath:features/apis/state/get.feature@CreateComplianceProfileState')
        * def createdProfileState2 = call read('classpath:features/apis/state/get.feature@CreateComplianceProfileState')
        * def createdProfileState3 = call read('classpath:features/apis/state/get.feature@CreateComplianceProfileState')
        * def createdProfileState4 = call read('classpath:features/apis/state/get.feature@CreateComplianceProfileState')
        * def createdProfileState5 = call read('classpath:features/apis/state/get.feature@CreateComplianceProfileState')

        * def entityIdNotFound = uuid()
        * def entityIds = { "ids": ['#(createdProfileState1.profileID)','#(entityIdNotFound)','#(createdProfileState2.profileID)','#(createdProfileState3.profileID)','#(createdProfileState4.profileID)','#(createdProfileState5.profileID)'] }
      
        * def expectedResynced = 
        """ 
          { 
            "resynced": ["#(createdProfileState1.profileID)","#(createdProfileState2.profileID)","#(createdProfileState3.profileID)","#(createdProfileState4.profileID)","#(createdProfileState5.profileID)"] 
          }
        """
        
        Given url baseURLCompliance
        And path '/states/resync'
        And request entityIds
        When method POST
        Then assert responseStatus == 200
        And match response.resynced contains expectedResynced.resynced
    
        * def expectedResyncStateEvent1 =
        """
        {
            "id": "#string",
            "entity_id": "#(createdProfileState1.profileID)",
            "entity_type": "COMPLIANCE_STATE",
            "event_type": "STATE_RESYNC",
            "update_date": "#string",
            "data": {
                "profile_id": "#(createdProfileState1.profileID)",
                "content": {
                    "state": {
                        "entity_id": "#(createdProfileState1.profileID)",
                        "engine_name": "PROFILE",
                        "result": "APPROVED",
                        "validation_steps_results": "##[]",
                        "rule_names": "##[]",
                        "pending": false,
                        "request_date": "##string",
                        "execution_time": "#string",
                        "version": "#number",
                        "created_at": "#string",
                        "updated_at": "#string"
                    },
                    "profile": {
                        "person": {
                            "document_number": "#string",
                            "name": "#string",
                            "person_type": "INDIVIDUAL",
                            "email": "#string",
                            "partner_id": "#string",
                            "offer_type": "#string",
                            "profile_id": "#(createdProfileState1.profileID)",
                            "entity_id": "#(createdProfileState1.profileID)",
                            "entity_type": "PROFILE",
                            "role_type": "CUSTOMER",
                            "cadastral_validation_config": "##object"
                        },
                        "profile_id": "#(createdProfileState1.profileID)",
                        "created_at": "#string",
                        "updated_at": "#string"
                    },
                    "person": "##object",
                    "contract": "##object"
                }
            }
        }
        """
    
        Given url mockURL
        And path '/subscribe/state-events/' + createdProfileState1.profileID
        And retry until response.length == 3
        When method GET
        Then assert responseStatus == 200
        And match response contains deep expectedResyncStateEvent1
    
        * def expectedResyncStateEvent2 =
        """
        {
            "id": "#string",
            "entity_id": "#(createdProfileState2.profileID)",
            "entity_type": "COMPLIANCE_STATE",
            "event_type": "STATE_RESYNC",
            "update_date": "#string",
            "data": {
                "profile_id": "#(createdProfileState2.profileID)",
                "content": {
                    "state": {
                        "entity_id": "#(createdProfileState2.profileID)",
                        "engine_name": "PROFILE",
                        "result": "APPROVED",
                        "validation_steps_results": "##[]",
                        "rule_names": "##[]",
                        "pending": false,
                        "request_date": "##string",
                        "execution_time": "#string",
                        "version": "#number",
                        "created_at": "#string",
                        "updated_at": "#string"
                    },
                    "profile": {
                        "person": {
                            "document_number": "#string",
                            "name": "#string",
                            "person_type": "INDIVIDUAL",
                            "email": "#string",
                            "partner_id": "#string",
                            "offer_type": "#string",
                            "profile_id": "#(createdProfileState2.profileID)",
                            "entity_id": "#(createdProfileState2.profileID)",
                            "entity_type": "PROFILE",
                            "role_type": "CUSTOMER",
                            "cadastral_validation_config": "##object"
                        },
                        "profile_id": "#(createdProfileState2.profileID)",
                        "created_at": "#string",
                        "updated_at": "#string"
                    },
                    "person": "##object",
                    "contract": "##object"
                }
            }
        }
        """
    
        Given url mockURL
        And path '/subscribe/state-events/' + createdProfileState2.profileID
        And retry until response.length == 3
        When method GET
        Then assert responseStatus == 200
        And match response contains deep expectedResyncStateEvent2

        * def expectedResyncStateEvent3 =
        """
        {
            "id": "#string",
            "entity_id": "#(createdProfileState3.profileID)",
            "entity_type": "COMPLIANCE_STATE",
            "event_type": "STATE_RESYNC",
            "update_date": "#string",
            "data": {
                "profile_id": "#(createdProfileState3.profileID)",
                "content": {
                    "state": {
                        "entity_id": "#(createdProfileState3.profileID)",
                        "engine_name": "PROFILE",
                        "result": "APPROVED",
                        "validation_steps_results": "##[]",
                        "rule_names": "##[]",
                        "pending": false,
                        "request_date": "##string",
                        "execution_time": "#string",
                        "version": "#number",
                        "created_at": "#string",
                        "updated_at": "#string"
                    },
                    "profile": "##object",
                    "person": "##object",
                    "contract": "##object"
                }
            }
        }
        """
    
        Given url mockURL
        And path '/subscribe/state-events/' + createdProfileState3.profileID
        And retry until response.length == 3
        When method GET
        Then assert responseStatus == 200
        And match response contains deep expectedResyncStateEvent3

        * def expectedResyncStateEvent4 =
        """
        {
            "id": "#string",
            "entity_id": "#(createdProfileState4.profileID)",
            "entity_type": "COMPLIANCE_STATE",
            "event_type": "STATE_RESYNC",
            "update_date": "#string",
            "data": {
                "profile_id": "#(createdProfileState4.profileID)",
                "content": {
                    "state": {
                        "entity_id": "#(createdProfileState4.profileID)",
                        "engine_name": "PROFILE",
                        "result": "APPROVED",
                        "validation_steps_results": "##[]",
                        "rule_names": "##[]",
                        "pending": false,
                        "request_date": "##string",
                        "execution_time": "#string",
                        "version": "#number",
                        "created_at": "#string",
                        "updated_at": "#string"
                    },
                    "profile": "##object",
                    "person": "##object",
                    "contract": "##object"
                }
            }
        }
        """
    
        Given url mockURL
        And path '/subscribe/state-events/' + createdProfileState4.profileID
        And retry until response.length == 3
        When method GET
        Then assert responseStatus == 200
        And match response contains deep expectedResyncStateEvent4

        * def expectedResyncStateEvent5 =
        """
        {
            "id": "#string",
            "entity_id": "#(createdProfileState5.profileID)",
            "entity_type": "COMPLIANCE_STATE",
            "event_type": "STATE_RESYNC",
            "update_date": "#string",
            "data": {
                "profile_id": "#(createdProfileState5.profileID)",
                "content": {
                    "state": {
                        "entity_id": "#(createdProfileState5.profileID)",
                        "engine_name": "PROFILE",
                        "result": "APPROVED",
                        "validation_steps_results": "##[]",
                        "rule_names": "##[]",
                        "pending": false,
                        "request_date": "##string",
                        "execution_time": "#string",
                        "version": "#number",
                        "created_at": "#string",
                        "updated_at": "#string"
                    },
                    "profile": "##object",
                    "person": "##object",
                    "contract": "##object"
                }
            }
        }
        """
    
        Given url mockURL
        And path '/subscribe/state-events/' + createdProfileState5.profileID
        And retry until response.length == 3
        When method GET
        Then assert responseStatus == 200
        And match response contains deep expectedResyncStateEvent5

    Scenario: should process a event with the same update_date
        * def documentNumberIndividual = DocumentNormalizer(CPFGenerator())
        * def dateOfBirth = "2023-03-15T15:06:56.12Z"
        * def profileResponse =
        """
        {
            "profile_id": "#(profileID)",
            "partner_id": "#(uuid())",
            "document_number": "#(documentNumberIndividual)",
            "parent_id": "#(uuid())",
            "name": "John Doe",
            "legacy_id": "#(uuid())",
            "offer_type": "#(BureauRuleOffer)",
            "role_type":"#(RoleType.Customer)",
            "profile_type": "#(ProfileType.Individual)",
            "callback_url": "/callback",
            "individual":{
                "first_name": "John",
                "last_name": "Doe",
                "date_of_birth": "#(dateOfBirth)",
                "date_of_birth_inputted": "#(dateOfBirth)",
                "phones": [
                    {
                        "type" : "RESIDENTIAL",
                        "area_code": "13",
                        "country_code": "55",
                        "number": "888889999"
                    }
                ],
                "bureau_information": {
                    "name": "John Stewart Doe",
                    "date_of_birth": "#(dateOfBirth)"
                },
                "income": 10000,
                "income_currency": "BRL",
                "assets": 100000,
                "assets_currency": "BRL",
                "nationality": "BRA",
                "email": "some@mail.com",
                "pep": false,
                "us_person": false,
                "occupation": "Enginer",
                "foreign_tax_residency": false,
                "country_of_tax_residency": "BRA"
            },
            "email": "Some@mail.com",
            "created_at": "#(funcNow(2))",
            "updated_at": "#(funcNow(2))" 
        }    
        """

        Given url mockURL
        And path '/v1/temis/profile/' + profileID
        And request profileResponse
        When method POST
        Then assert responseStatus == 200

        * def enrichedBirthDate = "23/11/1992"
        * def enrichmentResponse = { situation: 1, name: "Johnatan Stewart Doodle", birth_date: "#(enrichedBirthDate)" }

		Given url mockURL
		And path '/temis-enrichment/individual/' + documentNumberIndividual
		And request enrichmentResponse
		When method POST

        * def event = { profile_id: '#(profileID)', entity_id: '#(profileID)', entity_type: 'PROFILE', event_type: 'PROFILE_CREATED', parent_type: 'PROFILE', update_date: '2023-08-03T19:29:06.63556Z'  }
		* string json = event
    	
		* def result = RegistrationEventsPublisher.publish(json)

		Given url baseURLCompliance
		And path '/state/' + profileID
        And header Content-Type = 'application/json'
        And params ({ only_pending: false })
		And retry until response.result == RuleResult.Approved
        When method GET
		Then assert responseStatus == 200

        * def enrichedBirthDate = "23/11/1992"
        * def enrichmentResponse = { situation: 4, name: "Johnatan Stewart Doodle", birth_date: "#(enrichedBirthDate)" }

		Given url mockURL
		And path '/temis-enrichment/individual/' + documentNumberIndividual
		And request enrichmentResponse
		When method POST

        * def event = { profile_id: '#(profileID)', entity_id: '#(profileID)', entity_type: 'PROFILE', event_type: 'PROFILE_UPDATE', parent_type: 'PROFILE', update_date: '2023-08-03T19:29:06.63556Z'  }
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

        Given url mockURL
		And path "/subscribe/state-events/" + profileID
		And retry until response.length == 3
		When method GET
		Then assert responseStatus == 200