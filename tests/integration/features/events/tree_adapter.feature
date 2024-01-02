#noinspection CucumberUndefinedStep
Feature: Send an event to Tree Adapter

	Background:
		* url baseURLCompliance

		* def treeAdapterEventOffer = "TEST_OFFER" + uuid() 

		* def serasaIndividualRegularStatus = 1
		* def serasaIndividualNotRegularStatus = 9
		* def serasaCompanyRegularStatus = 2
		* def serasaCompanyNotRegularStatus = 9

		* def offer = { offer_type : '#(treeAdapterEventOffer)', product: 'maquininha'}

		Given path '/offers'
		And header Content-Type = 'application/json'
		And request offer
		When method POST
		Then assert responseStatus == 201

		* def ruleSetConfig = { bureau: {} }

		* def catalog = call CreateSingleLevelCatalog { offer_type: '#(treeAdapterEventOffer)', role_type: '#(RoleType.Customer)', person_type: '#(ProfileType.Individual)', account_flag: false, rules_config: #(ruleSetConfig)}
		* catalog.product_config.tree_integration = true

		Given url mockURL
        And path '/temis-config/cadastral-validation-configs'
        And header Content-Type = 'application/json'
		And request catalog
		When method POST
		Then assert responseStatus == 200

		* def catalogLR = call CreateSingleLevelCatalog { offer_type: '#(treeAdapterEventOffer)', role_type: '#(RoleType.LegalRepresentative)', person_type: '#(ProfileType.Individual)', account_flag: false, rules_config: #(ruleSetConfig) }
		* catalogLR.product_config.tree_integration = true

		Given url mockURL
        And path '/temis-config/cadastral-validation-configs'
        And header Content-Type = 'application/json'
		And request catalogLR
		When method POST
		Then assert responseStatus == 200

		* def catalogCompany = call CreateSingleLevelCatalog { offer_type: '#(treeAdapterEventOffer)', role_type: '#(RoleType.Customer)', person_type: '#(ProfileType.Company)', account_flag: false, rules_config: #(ruleSetConfig)}
		* catalogCompany.product_config.tree_integration = true

		Given url mockURL
        And path '/temis-config/cadastral-validation-configs'
        And header Content-Type = 'application/json'
		And request catalogCompany
		When method POST
		Then assert responseStatus == 200

		* def catalogCompanyLR = call CreateSingleLevelCatalog { offer_type: '#(treeAdapterEventOffer)', role_type: '#(RoleType.LegalRepresentative)', person_type: '#(ProfileType.Company)', account_flag: false, rules_config: #(ruleSetConfig)}
		* catalogCompanyLR.product_config.tree_integration = true

		Given url mockURL
        And path '/temis-config/cadastral-validation-configs'
        And header Content-Type = 'application/json'
		And request catalogCompanyLR
		When method POST
		Then assert responseStatus == 200

		* def documentNumberIndividual =  CPFGenerator()
		* def documentNumberCompany = DocumentNormalizer(CNPJGenerator())
		* def partnerID = uuid()
		* def parentID = uuid()
		* def profileID = uuid()

	Scenario: should send a profile individual event to tree adapter

		* def profile =
        """
            {
                profile_id: '#(profileID)',
                partner_id: '#(partnerID)',
                parent_id: '#(parentID)',
                offer_type: '#(treeAdapterEventOffer)',
                callback_url: '/url',
                role_type: '#(RoleType.Customer)',
                profile_type: '#(ProfileType.Individual)',
                document_number: '#(documentNumberIndividual)',
                individual: {
                    first_name: 'Daniel',
                    last_name: 'Alves',
                    date_of_birth: '1991-05-06T00:00:00Z',
                    nationality: 'BRA',
                    pep: true,
					income: 120.0,
					assets: 265000
                }
            }
        """

		Given url mockURL
		And path '/v1/temis/profile/' + profileID
		And request profile
		When method POST
		Then assert responseStatus == 200

		* def notificationRecipients =
        """
            [
                {
					notification_type: "POST_WARNINGS",
					email_to: "postWarningsEmailTo@teste.com.br",
					copy_email: "postWarningsCopyEmail@teste.com.br",
				},
				{
					notification_type: "SENT_OP",
					email_to:          "sentOPEmailTo@teste.com.br",
					copy_email:        "sentOPCopyEmail@teste.com.br",
				},
			]
        """

		Given url mockURL
		And path '/v1/temis/notification-recipients'
		And params ({ profile_id: profileID })
		And request notificationRecipients
		When method POST
		Then assert responseStatus == 200

		* def enrichmentResponse = { situation: "#(serasaIndividualRegularStatus)" }

		Given url mockURL
		And path '/temis-enrichment/individual/' + documentNumberIndividual
		And request enrichmentResponse
		When method POST
		Then assert responseStatus == 200

		* def partnerName = "PartnerXXX"
		* def partnerResponse = { partner_id: '#(partnerID)', name:'#(partnerName)', status: 'ACTIVE'}

		Given url mockURL
		And path '/v1/temis/partner/' + partnerID
		And request partnerResponse
		When method POST
		Then assert responseStatus == 200

		* def accounts = [{account_number: '11', account_digit: '1', agency_number: '22', agency_digit: '2', bank_code: '3'}, {account_number: '44', account_digit: '4', agency_number: '55', agency_digit: '5', bank_code: '6'}]

		Given url mockURL
		And path '/v1/temis/profile/' + profileID + '/accounts'
		And request accounts
		When method POST
		Then assert responseStatus == 200

		* def addresses = 
		"""
			[
				{ 
					address_id: '#(uuid())',
					profile_id: '#(profileID)',
					type: 'RESIDENTIAL',
					zip_code: '1234',
					street: 'Rua Roberto Menotti',
					number: '567',
					complement: 'Próximo ao Salão de Beleza',
					neighborhood: 'Vila Albuquerque',
					city: 'São Paulo',
					state_code: 'SP',
					country_code: 'BRA',
					updated_at: '2019-05-06T00:00:00Z'
				},
				{ 
					address_id: '#(uuid())',
					profile_id: '#(profileID)',
					type: 'RESIDENTIAL',
					zip_code: '5555',
					street: 'Rua Arapongas',
					number: '5555',
					complement: 'ap 999 conj 7',
					neighborhood: 'Vila Yara',
					city: 'Belo Horizonte',
					state_code: 'MG',
					country_code: 'BRA',
					updated_at: '2021-05-06T00:00:00Z'
				}
			]
		"""

		Given url mockURL
		And path '/v1/temis/addresses'
		And params ({ profile_id: profileID })
		And request addresses
		When method POST

		* def contacts = 
		"""
			[
				{
					"contact_id": '#(uuid())',
					"profile_id": '#(profileID)',
					"name": 'RESIDENTIAL',
					"document_number": '123434',
					"email": 'contact@test1.com',
					"phones": [
							{ "number": "01123456789" },
							{ "number": "0" }
					],
					"nationality": 'BRA'
				}
			]
		"""

		Given url mockURL
		And path '/v1/temis/contacts'
		And params ({ profile_id: profileID })
		And request contacts
		When method POST

		* def event = { profile_id: '#(profileID)', entity_id: '#(profileID)', entity_type: 'PROFILE', event_type: 'PROFILE_CHANGED', parent_type: 'PROFILE', update_date: '2023-08-03T19:29:06.63556Z'  }
		* string json = event
		
		* def result = RegistrationEventsPublisher.publish(json)

		Given url baseURLCompliance
		And path '/state/' + profileID
		And header Content-Type = 'application/json'
		And params ({ only_pending: false })
		And retry until response.rule_set_result.length == 2
		When method GET
		And assert response.result == RuleResult.Approved
		And assert response.rule_set_result[0].result == RuleResult.Approved
		And assert response.rule_set_result[1].result == RuleResult.Approved

		* def name = profile.individual.first_name + ' ' + profile.individual.last_name
		* def birthDate = '06051991'
		* def treeAdapterMessage =
        """
        {
            "profile_id": '#(profile.profile_id)',
            "partner": '#(partnerName)',
			"profile_created_date": "#notnull",
			"negotiation": "MESA",
            "person": {
                "id": '#(profile.document_number)',
                "type": '#(profile.profile_type)',
                "name": '#(name)',
                "birth_date": '#(birthDate)',
                "nationality": '#(profile.individual.nationality)',
				"digital_sign": false,
				"accounts": [{"number":"111","agency_number":"222","bank_code":"3"},{"number":"444","agency_number":"555","bank_code":"6"}],
				"addresses": [
					{
						"type": "#(addresses[1].type)",
						"street": "#(addresses[1].street)",
						"number": "#(addresses[1].number)",
						"complement": "#(addresses[1].complement)",
						"neighborhood": "#(addresses[1].neighborhood)",
						"city": "#(addresses[1].city)",
						"state": "#(addresses[1].state_code)",
						"country": "#(addresses[1].country_code)",
						"code": "#(addresses[1].zip_code)",
					},
					{
						"type": "#(addresses[0].type)",
						"street": "#(addresses[0].street)",
						"number": "#(addresses[0].number)",
						"complement": "#(addresses[0].complement)",
						"neighborhood": "#(addresses[0].neighborhood)",
						"city": "#(addresses[0].city)",
						"state": "#(addresses[0].state_code)",
						"country": "#(addresses[0].country_code)",
						"code": "#(addresses[0].zip_code)",
					}
				
				],
				"contacts": [
					{
						"name": "#(contacts[0].name)",
						"email": "#(contacts[0].email)",
						"phones": [
							{ "number": 01123456789 },
							{ "number": 0 }
						]
					}
				],
				"notification_recipients": '#(notificationRecipients)',
                "individual": {
                    "pep": '#(profile.individual.pep)',
					"income": '#(profile.individual.income)',
					"assets": '#(profile.individual.assets)'
                }
            },
			"date": "#notnull",
        }
        """
		
		Given url mockURL
		And path '/subscribe/tree-adapter/' + profileID
		And retry until response.length == 1
		When method GET
		Then assert responseStatus == 200
		* string json = response
		And match response contains deep treeAdapterMessage
		And assert response[0].date != null
		And match response[0].person.addresses[0] == treeAdapterMessage.person.addresses[0]
		And match response[0].person.addresses[1] == treeAdapterMessage.person.addresses[1]
		And match response[0].person.contacts[0] == treeAdapterMessage.person.contacts[0]

	Scenario: should send a profile individual event to tree adapter without accounts and addresses

		* def profile =
        """
            {
                profile_id: '#(profileID)',
                partner_id: '#(partnerID)',
                parent_id: '#(parentID)',
                offer_type: '#(treeAdapterEventOffer)',
                callback_url: '/url',
                role_type: '#(RoleType.Customer)',
                profile_type: '#(ProfileType.Individual)',
                document_number: '#(documentNumberIndividual)',
                individual: {
                    first_name: 'Daniel',
                    last_name: 'Alves',
                    date_of_birth: '1991-05-06T00:00:00Z',
                    nationality: 'BRA',
                    pep: true,
					income: 120.0,
					assets: 265000
                }
            }
        """

		Given url mockURL
		And path '/v1/temis/profile/' + profileID
		And request profile
		When method POST
		Then assert responseStatus == 200

		* def enrichmentResponse = { situation: "#(serasaIndividualRegularStatus)" }

		Given url mockURL
		And path '/temis-enrichment/individual/' + documentNumberIndividual
		And request enrichmentResponse
		When method POST
		Then assert responseStatus == 200

		* def partnerName = "PartnerXXX"
		* def partnerResponse = { partner_id: '#(partnerID)', name:'#(partnerName)', status: 'ACTIVE'}

		Given url mockURL
		And path '/v1/temis/partner/' + partnerID
		And request partnerResponse
		When method POST
		Then assert responseStatus == 200

		* def accounts = []

		Given url mockURL
		And path '/v1/temis/profile/' + profileID + '/accounts'
		And request accounts
		When method POST
		Then assert responseStatus == 200

		* def addresses = []

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
		And retry until response.rule_set_result.length == 2
		When method GET
		And assert response.result == RuleResult.Approved
		And assert response.rule_set_result[0].result == RuleResult.Approved
		And assert response.rule_set_result[1].result == RuleResult.Approved

		* def name = profile.individual.first_name + ' ' + profile.individual.last_name
		* def birthDate = '06051991'
		* def treeAdapterMessage =
        """
        {
            "profile_id": '#(profile.profile_id)',
            "partner": '#(partnerName)',
			"profile_created_date": "#notnull",
			"negotiation": "MESA",
            "person": {
                "id": '#(profile.document_number)',
                "type": '#(profile.profile_type)',
                "name": '#(name)',
                "birth_date": '#(birthDate)',
                "nationality": '#(profile.individual.nationality)',
				"digital_sign": false,
                "individual": {
                    "pep": '#(profile.individual.pep)',
					"income": '#(profile.individual.income)',
					"assets": '#(profile.individual.assets)'
                }
            },
			"date": "#notnull",
        }
        """
		
		Given url mockURL
		And path '/subscribe/tree-adapter/' + profileID
		And retry until response.length == 1
		When method GET
		Then assert responseStatus == 200
		And match response contains deep treeAdapterMessage
		And assert response[0].date != null

	Scenario: should not send a profle individual event to tree adapter, when it is not APPROVED

		* def profile = { profile_id: '#(profileID)', offer_type: '#(treeAdapterEventOffer)', role_type: '#(RoleType.Customer)', profile_type: '#(ProfileType.Individual)', document_number: '#(documentNumberIndividual)' }

		Given url mockURL
		And path '/v1/temis/profile/' + profileID
		And request profile
		When method POST
		Then assert responseStatus == 200

		* def enrichmentResponse = { situation: "#(serasaIndividualNotRegularStatus)" }

		Given url mockURL
		And path '/temis-enrichment/individual/' + documentNumberIndividual
		And request enrichmentResponse
		When method POST
		Then assert responseStatus == 200

		* def event = { profile_id: '#(profileID)', entity_id: '#(profileID)', entity_type: 'PROFILE', event_type: 'PROFILE_CHANGED', parent_type: 'PROFILE', update_date: '2023-08-03T19:29:06.63556Z'  }
		* string json = event
		* def result = RegistrationEventsPublisher.publish(json)
		* eval sleep(1000)

		Given url baseURLCompliance
		And path '/state/' + profileID
        And header Content-Type = 'application/json'
        And params ({ only_pending: false })
		And retry until response.result == RuleResult.Rejected
        When method GET
		And assert response.result == RuleResult.Rejected
		
		Given url mockURL
		And path '/subscribe/tree-adapter/' + profileID
		And retry until response.length == 0
		When method GET
		Then assert responseStatus == 200
	
	Scenario: should not send a profile individual event to tree adapter, when role type is not #(RoleType.Customer)

		* def profile =
        """
            {
                profile_id: '#(profileID)',
                offer_type: '#(treeAdapterEventOffer)',
                role_type: '#(RoleType.LegalRepresentative)',
                profile_type: '#(ProfileType.Individual)',
                document_number: '#(documentNumberIndividual)'
            }
        """

		Given url mockURL
		And path '/v1/temis/profile/' + profileID
		And request profile
		When method POST
		Then assert responseStatus == 200

		* def enrichmentResponse = { situation: "#(serasaIndividualRegularStatus)" }

		Given url mockURL
		And path '/temis-enrichment/individual/' + documentNumberIndividual
		And request enrichmentResponse
		When method POST
		Then assert responseStatus == 200

		* def event = { profile_id: '#(profileID)', entity_id: '#(profileID)', entity_type: 'PROFILE', event_type: 'PROFILE_CHANGED', parent_type: 'PROFILE', update_date: '2023-08-03T19:29:06.63556Z'  }
		* string json = event
		* def result = RegistrationEventsPublisher.publish(json)
		* eval sleep(1000)

		Given url baseURLCompliance
		And path '/state/' + profileID
        And header Content-Type = 'application/json'
        And params ({ only_pending: false })
		And retry until response.result == RuleResult.Approved
        When method GET
		And assert response.result == RuleResult.Approved

		Given url mockURL
		And path '/subscribe/tree-adapter/' + profileID
		And retry until response.length == 0
		When method GET
		Then assert responseStatus == 200

	Scenario: should send a profile company event to tree adapter

		* def profile =
        """
            {
                profile_id: '#(profileID)',
                partner_id: '#(partnerID)',
                parent_id: '#(parentID)',
                offer_type: '#(treeAdapterEventOffer)',
                callback_url: '/url',
                role_type: '#(RoleType.Customer)',
                profile_type: '#(ProfileType.Company)',
                document_number: '#(documentNumberCompany)',
                company: {
                    legal_name: 'Empresa SA',
					date_of_incorporation: '2018-08-08',
					place_of_incorporation: 'BRA',
					annual_income: 120.0
                }
            }
        """

		Given url mockURL
		And path '/v1/temis/profile/' + profileID
		And request profile
		When method POST
		Then assert responseStatus == 200

		* def notificationRecipients =
        """
            [
                {
					notification_type: "POST_WARNINGS",
					email_to: 		   "postWarningsEmailTo@teste.com.br",
					copy_email: 	   "postWarningsCopyEmail@teste.com.br",
				},
				{
					notification_type: "SENT_OP",
					email_to:          "sentOPEmailTo@teste.com.br",
					copy_email:        "sentOPCopyEmail@teste.com.br",
				},
			]
        """

		Given url mockURL
		And path '/v1/temis/notification-recipients'
		And params ({ profile_id: profileID })
		And request notificationRecipients
		When method POST
		Then assert responseStatus == 200

		* def enrichmentResponse = { situation: "#(serasaCompanyRegularStatus)" }

		Given url mockURL
		And path '/temis-enrichment/legal-entity/' + documentNumberCompany
		And request enrichmentResponse
		When method POST
		Then assert responseStatus == 200

		* def partnerName = "PartnerXXX"
		* def partnerResponse = { partner_id: '#(partnerID)', name:'#(partnerName)', status: 'ACTIVE'}

		Given url mockURL
		And path '/v1/temis/partner/' + partnerID
		And request partnerResponse
		When method POST
		Then assert responseStatus == 200

		* def accounts = [{account_number: '11', account_digit: '1', agency_number: '22', agency_digit: '2', bank_code: '3'}, {account_number: '44', account_digit: '4', agency_number: '55', agency_digit: '5', bank_code: '6'}]

		Given url mockURL
		And path '/v1/temis/profile/' + profileID + '/accounts'
		And request accounts
		When method POST
		Then assert responseStatus == 200

		* def addresses = 
		"""
			[
				{ 
					address_id: '#(uuid())',
					profile_id: '#(profileID)',
					type: 'RESIDENTIAL',
					zip_code: '777',
					street: 'Rua Campos Neto',
					number: '333',
					complement: 'Próximo ao Parque',
					neighborhood: 'Vila Aurea',
					city: 'Rio de Janeiro',
					state_code: 'RJ',
					country_code: 'BRA',
					updated_at: '2011-05-06T00:00:00Z'
				},
				{ 
					address_id: '#(uuid())',
					profile_id: '#(profileID)',
					type: 'RESIDENTIAL',
					zip_code: '5555',
					street: 'Rua Arapongas',
					number: '5555',
					complement: 'ap 999 conj 7',
					neighborhood: 'Vila Yara',
					city: 'Belo Horizonte',
					state_code: 'MG',
					country_code: 'BRA',
					updated_at: '2019-05-06T00:00:00Z'
				}
			]
		"""

		Given url mockURL
		And path '/v1/temis/addresses'
		And params ({ profile_id: profileID })
		And request addresses
		When method POST

		* def contacts = 
		"""
			[
				{
					"contact_id": '#(uuid())',
					"profile_id": '#(profileID)',
					"name": 'RESIDENTIAL',
					"document_number": '123434',
					"email": 'contact@test1.com',
					"phones": [
							{"number": "1123456789"},
							{"number": "0" }
					],
					"nationality": 'BRA'
				}
			]
		"""

		Given url mockURL
		And path '/v1/temis/contacts'
		And params ({ profile_id: profileID })
		And request contacts
		When method POST


		* def event = { profile_id: '#(profileID)', entity_id: '#(profileID)', entity_type: 'PROFILE', event_type: 'PROFILE_CHANGED', parent_type: 'PROFILE', update_date: '2023-08-03T19:29:06.63556Z'  }
		* string json = event
		
		* def result = RegistrationEventsPublisher.publish(json)

		Given url baseURLCompliance
		And path '/state/' + profileID
		And header Content-Type = 'application/json'
		And params ({ only_pending: false })
		And retry until response.rule_set_result.length == 2
		When method GET
		And assert response.result == RuleResult.Approved
		And assert response.rule_set_result[0].result == RuleResult.Approved
		And assert response.rule_set_result[1].result == RuleResult.Approved

		* def birthDate = '08082018'
		* def monthlyIncome = profile.company.annual_income / 12
		* def treeAdapterMessage =
        """
        {
            "profile_id": '#(profile.profile_id)',
            "partner": '#(partnerName)',
			"profile_created_date": '#notnull',
			"negotiation": "MESA",
            "person": {
                "id": '#(profile.document_number)',
                "type": '#(profile.profile_type)',
                "name": '#(profile.company.legal_name)',
                "birth_date": '#(birthDate)',
                "nationality": '#(profile.company.place_of_incorporation)',
				"digital_sign": false,
				"accounts": [{"number":"111","agency_number":"222","bank_code":"3"},{"number":"444","agency_number":"555","bank_code":"6"}],
				"addresses": [
					{
						"type": "#(addresses[1].type)",
						"street": "#(addresses[1].street)",
						"number": "#(addresses[1].number)",
						"complement": "#(addresses[1].complement)",
						"neighborhood": "#(addresses[1].neighborhood)",
						"city": "#(addresses[1].city)",
						"state": "#(addresses[1].state_code)",
						"country": "#(addresses[1].country_code)",
						"code": "#(addresses[1].zip_code)",
					},
					{
						"type": "#(addresses[0].type)",
						"street": "#(addresses[0].street)",
						"number": "#(addresses[0].number)",
						"complement": "#(addresses[0].complement)",
						"neighborhood": "#(addresses[0].neighborhood)",
						"city": "#(addresses[0].city)",
						"state": "#(addresses[0].state_code)",
						"country": "#(addresses[0].country_code)",
						"code": "#(addresses[0].zip_code)",
					}
				],
				"contacts": [
					{
						"name": "#(contacts[0].name)",
						"email": "#(contacts[0].email)",
						"phones": [
							{"number": 1123456789},
							{"number": 0 }
						]
					}
				],
				"notification_recipients": '#(notificationRecipients)',
                "legal_entity": {
                    "size": 'NAO_INFORMADO',
					"business_activity": 'A',
					"monthly_income": '#(monthlyIncome)'
                }
            },
			"date": '#notnull'
        }
        """

		Given url mockURL
		And path '/subscribe/tree-adapter/' + profileID
		And retry until response.length == 1
		When method GET
		Then assert responseStatus == 200
		And match response contains deep treeAdapterMessage
		And assert response[0].date != null
		And match response[0].person.addresses[0] == treeAdapterMessage.person.addresses[0]
		And match response[0].person.addresses[1] == treeAdapterMessage.person.addresses[1]
		And match response[0].person.contacts[0] == treeAdapterMessage.person.contacts[0]

	Scenario: should send a profile company event to tree adapter without accounts and addresses

		* def profile =
        """
            {
                profile_id: '#(profileID)',
                partner_id: '#(partnerID)',
                parent_id: '#(parentID)',
                offer_type: '#(treeAdapterEventOffer)',
                callback_url: '/url',
                role_type: '#(RoleType.Customer)',
                profile_type: '#(ProfileType.Company)',
                document_number: '#(documentNumberCompany)',
                company: {
                    legal_name: 'Empresa SA',
					date_of_incorporation: '2018-08-08',
					place_of_incorporation: 'BRA',
					annual_income: 120.0
                }
            }
        """

		Given url mockURL
		And path '/v1/temis/profile/' + profileID
		And request profile
		When method POST
		Then assert responseStatus == 200

		* def enrichmentResponse = { situation: "#(serasaCompanyRegularStatus)" }

		Given url mockURL
		And path '/temis-enrichment/legal-entity/' + documentNumberCompany
		And request enrichmentResponse
		When method POST
		Then assert responseStatus == 200

		* def partnerName = "PartnerXXX"
		* def partnerResponse = { partner_id: '#(partnerID)', name:'#(partnerName)', status: 'ACTIVE'}

		Given url mockURL
		And path '/v1/temis/partner/' + partnerID
		And request partnerResponse
		When method POST
		Then assert responseStatus == 200

		* def accounts = []

		Given url mockURL
		And path '/v1/temis/profile/' + profileID + '/accounts'
		And request accounts
		When method POST
		Then assert responseStatus == 200

		* def addresses = []

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
		And retry until response.rule_set_result.length == 2
		When method GET
		And assert response.result == RuleResult.Approved
		And assert response.rule_set_result[0].result == RuleResult.Approved
		And assert response.rule_set_result[1].result == RuleResult.Approved

		* def birthDate = '08082018'
		* def monthlyIncome = profile.company.annual_income / 12
		* def treeAdapterMessage =
        """
        {
            "profile_id": '#(profile.profile_id)',
            "partner": '#(partnerName)',
			"profile_created_date": '#notnull',
			"negotiation": "MESA",
            "person": {
                "id": '#(profile.document_number)',
                "type": '#(profile.profile_type)',
                "name": '#(profile.company.legal_name)',
                "birth_date": '#(birthDate)',
                "nationality": '#(profile.company.place_of_incorporation)',
				"digital_sign": false,
                "legal_entity": {
                    "size": 'NAO_INFORMADO',
					"business_activity": 'A',
					"monthly_income": '#(monthlyIncome)'
                }
            },
			"date": '#notnull'
        }
        """

		Given url mockURL
		And path '/subscribe/tree-adapter/' + profileID
		And retry until response.length == 1
		When method GET
		Then assert responseStatus == 200
		And match response contains deep treeAdapterMessage
		And assert response[0].date != null

	Scenario: should not send a profle company event to tree adapter, when it is not APPROVED

		* def profile = { profile_id: '#(profileID)', offer_type: '#(treeAdapterEventOffer)', role_type: '#(RoleType.Customer)', profile_type: '#(ProfileType.Company)', document_number: '#(documentNumberCompany)' }

		Given url mockURL
		And path '/v1/temis/profile/' + profileID
		And request profile
		When method POST
		Then assert responseStatus == 200

		* def enrichmentResponse = { situation: "#(serasaCompanyNotRegularStatus)" }

		Given url mockURL
		And path '/temis-enrichment/legal-entity/' + documentNumberCompany
		And request enrichmentResponse
		When method POST
		Then assert responseStatus == 200

		* def event = { profile_id: '#(profileID)', entity_id: '#(profileID)', entity_type: 'PROFILE', event_type: 'PROFILE_CHANGED', parent_type: 'PROFILE', update_date: '2023-08-03T19:29:06.63556Z'  }
		* string json = event
		* def result = RegistrationEventsPublisher.publish(json)
		* eval sleep(1000)

		Given url baseURLCompliance
		And path '/state/' + profileID
        And header Content-Type = 'application/json'
        And params ({ only_pending: false })
		And retry until response.result == RuleResult.Rejected
        When method GET
		And assert response.result == RuleResult.Rejected

		Given url mockURL
		And path '/subscribe/tree-adapter/' + profileID
		And retry until response.length == 0
		When method GET
		Then assert responseStatus == 200
	
	Scenario: should not send a profile company event to tree adapter, when role type is not #(RoleType.Customer)

		* def profile =
        """
            {
                profile_id: '#(profileID)',
                offer_type: '#(treeAdapterEventOffer)',
                role_type: '#(RoleType.LegalRepresentative)',
                profile_type: '#(ProfileType.Company)',
                document_number: '#(documentNumberCompany)'
            }
        """

		Given url mockURL
		And path '/v1/temis/profile/' + profileID
		And request profile
		When method POST
		Then assert responseStatus == 200

		* def enrichmentResponse = { situation: "#(serasaCompanyRegularStatus)" }

		Given url mockURL
		And path '/temis-enrichment/legal-entity/' + documentNumberCompany
		And request enrichmentResponse
		When method POST
		Then assert responseStatus == 200

		* def event = { profile_id: '#(profileID)', entity_id: '#(profileID)', entity_type: 'PROFILE', event_type: 'PROFILE_CHANGED', parent_type: 'PROFILE', update_date: '2023-08-03T19:29:06.63556Z'  }
		* string json = event
		* def result = RegistrationEventsPublisher.publish(json)
		* eval sleep(1000)

		Given url baseURLCompliance
		And path '/state/' + profileID
        And header Content-Type = 'application/json'
        And params ({ only_pending: false })
		And retry until response.result == RuleResult.Approved
        When method GET
		And assert response.result == RuleResult.Approved

		Given url mockURL
		And path '/subscribe/tree-adapter/' + profileID
		And retry until response.length == 0
		When method GET
		Then assert responseStatus == 200
