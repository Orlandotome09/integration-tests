#noinspection CucumberUndefinedStep
Feature: Limit Event

	Background:
		* url baseURLCompliance

		* def limitEventOffer = "TEST_OFFER" + uuid() 
		* def offer = { offer_type : '#(limitEventOffer)', product: 'maquininha'}

		Given path '/offers'
		And header Content-Type = 'application/json'
		And request offer
		When method POST
		Then assert responseStatus == 201

		* def documentNumber =  CPFGenerator()
		* def partnerID = uuid()
		* def parentID = uuid()
		* def profileID = uuid()

	Scenario: should send an event to the limit queue

		* def ruleSetConfig = {  bureau: {} }
		* def ruleSetConfig1 = {  under_age: {} }
		* def ruleSetConfig2 = { incomplete: { email_required: true } }
		
		* def params =
		"""
			{
				offer_type: '#(limitEventOffer)',
				role_type: '#(RoleType.Customer)',
				person_type: '#(ProfileType.Individual)',
				steps: [
					{ rules_config: #(ruleSetConfig) },
					{ rules_config: #(ruleSetConfig1) },
					{ rules_config: #(ruleSetConfig2) }
				],
				limit_flag: true
			}
		"""		
		* def catalog = call CreateMultiLevelCatalog params

		Given url mockURL
        And path '/temis-config/cadastral-validation-configs'
        And header Content-Type = 'application/json'
		And request catalog
		When method POST
		Then assert responseStatus == 200

		* def profile =
        """
            {
                profile_id: '#(profileID)',
                partner_id: '#(partnerID)',
                parent_id: '#(parentID)',
                offer_type: '#(limitEventOffer)',
                callback_url: '/url',
                role_type: '#(RoleType.Customer)',
                profile_type: '#(ProfileType.Individual)',
                document_number: '#(documentNumber)',
                individual: {
                    first_name: 'Daniel',
                    last_name: 'Alves',
                    nationality: 'BRA',
                    pep: false
                }
            }
        """

		Given url mockURL
		And path '/v1/temis/profile/' + profileID
		And request profile
		When method POST
		Then assert responseStatus == 200

		* def enrichmentResponse = { situation: 1 }

		Given url mockURL
		And path '/temis-enrichment/individual/' + documentNumber
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

		* def timeNow = call funcNow
		* def event = { profile_id: '#(profileID)', entity_id: '#(profileID)', entity_type: 'PROFILE', event_type: 'PROFILE_CHANGED', parent_type: 'PROFILE', update_date: '#(timeNow)'  }
		* string json = event
		
		* def result = RegistrationEventsPublisher.publish(json)

		* def retryCondition = 
		"""
			function(response){
				return  response.entity_id == profileID &&
						response.rule_set_result.length == 3 &&
						response.rule_set_result[0].result == RuleResult.Approved &&
						response.rule_set_result[1].result == RuleResult.Approved &&
						response.rule_set_result[2].result == RuleResult.Rejected
				
			}
		"""

		Given url baseURLCompliance
		And path "/state/" + profileID
		And retry until retryCondition(response) == true
		When method GET
		Then assert responseStatus == 200
		And assert response.result == RuleResult.Approved

		* def limitMessageEvent1 =
        """
			{
				"event_type" : "PROFILE_APPROVED",
				"profile_id": '#(profile.profile_id)',
				"document_number": '#(profile.document_number)',
				"partner_id": '#(profile.partner_id)',
				"offer_type": '#(limitEventOffer)',
				"person_type": '#(ProfileType.Individual)',
				"role_type": '#(RoleType.Customer)',
				"approved_rules": ['CUSTOMER_NOT_FOUND_IN_SERASA', 'CUSTOMER_HAS_PROBLEMS_IN_SERASA']
			}
		"""

		Given url mockURL
		And path '/subscribe/limit/' + profileID
		And retry until response.length == 1
		When method GET
		Then assert responseStatus == 200
		And match response contains limitMessageEvent1
		
		* profile.individual.date_of_birth = '1991-05-06T00:00:00Z'

		Given url mockURL
		And path '/v1/temis/profile/' + profileID
		And request profile
		When method POST
		Then assert responseStatus == 200

		* def timeNow = call funcNow
		* def event = { profile_id: '#(profileID)', entity_id: '#(profileID)', entity_type: 'PROFILE', event_type: 'PROFILE_CHANGED', parent_type: 'PROFILE', update_date: '#(timeNow)'  }
		* string json = event
		
		* def result = RegistrationEventsPublisher.publish(json)

		* def retryCondition = 
		"""
			function(response){
				return  response.entity_id == profileID &&
						response.rule_set_result.length == 4 &&
						response.rule_set_result[0].result == RuleResult.Approved &&
						response.rule_set_result[1].result == RuleResult.Approved &&
						response.rule_set_result[2].result == RuleResult.Approved &&
						response.rule_set_result[3].result == RuleResult.Incomplete
				
			}
		"""

		Given url baseURLCompliance
		And path "/state/" + profileID
		And retry until retryCondition(response) == true
		When method GET
		Then assert responseStatus == 200
		And assert response.result == RuleResult.Approved

		* def limitMessageEvent2 =
		"""			
			{
				"event_type" : "PROFILE_APPROVED",
				"profile_id": '#(profile.profile_id)',
				"document_number": '#(profile.document_number)',
				"partner_id": '#(profile.partner_id)',
				"offer_type": '#(limitEventOffer)',
				"person_type": '#(ProfileType.Individual)',
				"role_type": '#(RoleType.Customer)',
				"approved_rules": ['CUSTOMER_NOT_FOUND_IN_SERASA', 'CUSTOMER_HAS_PROBLEMS_IN_SERASA', 'CUSTOMER_IS_UNDER_AGE']
			}
		"""
		
		Given url mockURL
		And path '/subscribe/limit/' + profileID
		And retry until response.length == 2
		When method GET
		Then assert responseStatus == 200
		And match response contains limitMessageEvent1
		And match response contains limitMessageEvent2

		* profile.email = 'daniel.alves@mock.test'

		Given url mockURL
		And path '/v1/temis/profile/' + profileID
		And request profile
		When method POST
		Then assert responseStatus == 200

		* def timeNow = call funcNow
		* def event = { profile_id: '#(profileID)', entity_id: '#(profileID)', entity_type: 'PROFILE', event_type: 'PROFILE_CHANGED', parent_type: 'PROFILE', update_date: '#(timeNow)'  }
		* string json = event
		
		* def result = RegistrationEventsPublisher.publish(json)
		
		* def retryCondition = 
		"""
			function(response){
				return  response.entity_id == profileID &&
						response.rule_set_result.length == 4 &&
						response.rule_set_result[0].result == RuleResult.Approved &&
						response.rule_set_result[1].result == RuleResult.Approved &&
						response.rule_set_result[2].result == RuleResult.Approved &&
						response.rule_set_result[3].result == RuleResult.Approved
				
			}
		"""

		Given url baseURLCompliance
		And path "/state/" + profileID
		And retry until retryCondition(response) == true
		When method GET
		Then assert responseStatus == 200
		And assert response.result == RuleResult.Approved
		
		* def limitMessageEvent3 =
		"""
			{
				"event_type" : "PROFILE_APPROVED",
				"profile_id": '#(profile.profile_id)',
				"document_number": '#(profile.document_number)',
				"partner_id": '#(profile.partner_id)',
				"offer_type": '#(limitEventOffer)',
				"person_type": '#(ProfileType.Individual)',
				"role_type": '#(RoleType.Customer)',
				"approved_rules": ['CUSTOMER_NOT_FOUND_IN_SERASA', 'CUSTOMER_HAS_PROBLEMS_IN_SERASA', 'CUSTOMER_IS_UNDER_AGE', 'REQUIRED_FIELDS_NOT_FOUND']
			}
        """

		Given url mockURL
		And path '/subscribe/limit/' + profileID
		And retry until response.length == 3
		When method GET
		Then assert responseStatus == 200
		And match response contains limitMessageEvent1
		And match response contains limitMessageEvent2
		And match response contains limitMessageEvent3

	Scenario: should send an event to the limit queue with documents

		# level 1
		* def ruleSetConfig = { incomplete: {documents_required: [{	"document_type": "CORPORATE_DOCUMENT","file_required": true}] } }

		* def params =
		"""
			{
				offer_type: '#(limitEventOffer)',
				role_type: '#(RoleType.Customer)',
				person_type: '#(ProfileType.Individual)',
				steps: [
					{ rules_config: #(ruleSetConfig) },
				],
				limit_flag: true
			}
		"""		
		* def catalog = call CreateMultiLevelCatalog params

		Given url mockURL
        And path '/temis-config/cadastral-validation-configs'
        And header Content-Type = 'application/json'
		And request catalog
		When method POST
		Then assert responseStatus == 200
		
		* def profile =
        """
            {
                profile_id: '#(profileID)',
                partner_id: '#(partnerID)',
                parent_id: '#(parentID)',
                offer_type: '#(limitEventOffer)',
                callback_url: '/url',
                role_type: '#(RoleType.Customer)',
                profile_type: '#(ProfileType.Individual)',
                document_number: '#(documentNumber)',
                individual: {
                    first_name: 'Daniel',
                    last_name: 'Alves',
                    nationality: 'BRA',
                    pep: false
                }
            }
        """

		Given url mockURL
		And path '/v1/temis/profile/' + profileID
		And request profile
		When method POST
		Then assert responseStatus == 200

		* def enrichmentResponse = { situation: 0 }

		Given url mockURL
		And path '/temis-enrichment/individual/' + documentNumber
		And request enrichmentResponse
		When method POST
		Then assert responseStatus == 200

		* def partnerName = "PartnerXXX"
		* def partnerResponse = { partner_id: '#(partnerID)', name:'#(partnerName)', status: 'ACTIVE'}

		* def documentID1 = uuid()
		* def document1 = { document_id: '#(documentID1)', entity_id: '#(profileID)', type: 'CORPORATE_DOCUMENT'}
		* def documents = []
		* def void = documents.add(document1)

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

		Given url mockURL
		And path '/v1/temis/partner/' + partnerID
		And request partnerResponse
		When method POST
		Then assert responseStatus == 200

		* def timeNow = call funcNow
		* def event = { profile_id: '#(profileID)', entity_id: '#(profileID)', entity_type: 'PROFILE', event_type: 'PROFILE_CHANGED', parent_type: 'PROFILE', update_date: '#(timeNow)'  }
		* string json = event
		
		* def result = RegistrationEventsPublisher.publish(json)

		* def limitMessageEvent1 =
        """
			{
				"event_type" : "PROFILE_APPROVED",
				"profile_id": '#(profile.profile_id)',
				"document_number": '#(profile.document_number)',
				"partner_id": '#(profile.partner_id)',
				"offer_type": '#(limitEventOffer)',
				"person_type": '#(ProfileType.Individual)',
				"role_type": '#(RoleType.Customer)',
				"approved_rules": ['DOCUMENT_NOT_FOUND'],
				"documents": ["CORPORATE_DOCUMENT"]
			}
		"""

		Given url mockURL
		And path '/subscribe/limit/' + profileID
		And retry until response.length == 1
		When method GET
		Then assert responseStatus == 200
		And match response contains limitMessageEvent1

	Scenario: should not send an event to the limit queue when is not APPROVED

		* def ruleSetConfig = { bureau: {} }
		* def params = { offer_type: '#(limitEventOffer)', role_type: '#(RoleType.Customer)', person_type: '#(ProfileType.Individual)', steps: [	{ rules_config: #(ruleSetConfig) } ], limit_flag: true }	
		* def catalog = call CreateMultiLevelCatalog params

		Given url mockURL
        And path '/temis-config/cadastral-validation-configs'
        And header Content-Type = 'application/json'
		And request catalog
		When method POST
		Then assert responseStatus == 200

		* def profile =
        """
            {
                profile_id: '#(profileID)',
                offer_type: '#(limitEventOffer)',
                callback_url: '/url',
                role_type: '#(RoleType.Customer)',
                profile_type: '#(ProfileType.Individual)',
                document_number: '#(documentNumber)',
                individual: {
                    first_name: 'Thiago',
                    last_name: 'Silva',
                    nationality: 'BRA',
                    pep: false
                }
            }
        """

		Given url mockURL
		And path '/v1/temis/profile/' + profileID
		And request profile
		When method POST
		Then assert responseStatus == 200

		* def event = { profile_id: '#(profileID)', entity_id: '#(profileID)', entity_type: 'PROFILE', event_type: 'PROFILE_CHANGED', parent_type: 'PROFILE', update_date: '2023-08-03T19:29:06.63556Z'  }
		* string json = event
		* def result = RegistrationEventsPublisher.publish(json)
		* eval sleep(1000)

		Given url mockURL
		And path '/subscribe/limit/' + profileID
		And retry until response.length == 0
		When method GET
		Then assert responseStatus == 200

	Scenario: should not send an event to the limit queue when is limit flag is false

		* def ruleSetConfig = { bureau: {} }
		* def params = { offer_type: '#(limitEventOffer)', role_type: '#(RoleType.Customer)', person_type: '#(ProfileType.Individual)', steps: [	{ rules_config: #(ruleSetConfig) } ], limit_flag: false }	
		* def catalog = call CreateMultiLevelCatalog params

		Given url mockURL
        And path '/temis-config/cadastral-validation-configs'
        And header Content-Type = 'application/json'
		And request catalog
		When method POST
		Then assert responseStatus == 200

		* def profile =
        """
            {
                profile_id: '#(profileID)',
                offer_type: '#(limitEventOffer)',
                callback_url: '/url',
                role_type: '#(RoleType.Customer)',
                profile_type: '#(ProfileType.Individual)',
                document_number: '#(documentNumber)',
                individual: {
                    first_name: 'Gabriel',
                    last_name: 'Jesus',
                    nationality: 'BRA',
                    pep: false
                }
            }
        """

		Given url mockURL
		And path '/v1/temis/profile/' + profileID
		And request profile
		When method POST
		Then assert responseStatus == 200

		* def enrichmentResponse = { situation: 0 }

		Given url mockURL
		And path '/temis-enrichment/individual/' + documentNumber
		And request enrichmentResponse
		When method POST
		Then assert responseStatus == 200

		* def event = { profile_id: '#(profileID)', entity_id: '#(profileID)', entity_type: 'PROFILE', event_type: 'PROFILE_CHANGED', parent_type: 'PROFILE', update_date: '2023-08-03T19:29:06.63556Z'  }
		* string json = event
		* def result = RegistrationEventsPublisher.publish(json)
		* eval sleep(2000)

		Given url mockURL
		And path '/subscribe/limit/' + profileID
		And retry until response.length == 0
		When method GET
		Then assert responseStatus == 200

	Scenario: should send an event to the limit queue for the business partner role type

		* def params =
		"""
			{
				offer_type: '#(limitEventOffer)',
				role_type: 'BUSINESS_PARTNER',
				person_type: 'COMPANY',
				steps: [
					{},
				],
				limit_flag: true
			}
		"""		
		* def catalog = call CreateMultiLevelCatalog params

		Given url mockURL
        And path '/temis-config/cadastral-validation-configs'
        And header Content-Type = 'application/json'
		And request catalog
		When method POST
		Then assert responseStatus == 200
		
		* def documentNumber =  DocumentNormalizer(CNPJGenerator())
		* def params =
        """
            {
                profile_id: '#(profileID)',
                partner_id: '#(partnerID)',
                parent_id: '#(parentID)',
                offer_type: '#(limitEventOffer)',
                role_type: 'BUSINESS_PARTNER',
                document_number: '#(documentNumber)',
                callback_url: '/url'
            }
        """
		* def profile = call CreateProfileCompany params

		Given url mockURL
		And path '/v1/temis/profile/' + profileID
		And request profile
		When method POST
		Then assert responseStatus == 200

		* def timeNow = call funcNow
		* def event = { profile_id: '#(profileID)', entity_id: '#(profileID)', entity_type: 'PROFILE', event_type: 'PROFILE_CHANGED', parent_type: 'PROFILE', update_date: '#(timeNow)'  }
		* string json = event
		
		* def result = RegistrationEventsPublisher.publish(json)

		* def limitMessageEvent =
        """
			{
				"event_type" : "PROFILE_APPROVED",
				"profile_id": '#(profile.profile_id)',
				"document_number": '#(profile.document_number)',
				"partner_id": '#(profile.partner_id)',
				"offer_type": '#(limitEventOffer)',
				"person_type": 'COMPANY',
				"role_type": 'BUSINESS_PARTNER',
				"approved_rules": []
			}
		"""
		
		Given url mockURL
		And path '/subscribe/limit/' + profileID
		And retry until response.length == 1
		When method GET
		Then assert responseStatus == 200
		And match response contains limitMessageEvent