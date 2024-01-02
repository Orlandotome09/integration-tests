Feature: Temis Query Event

    Background:
        * url baseURLCompliance

        * def offer = "TEST_OFFER" + uuid() 

        * def offerBody = { offer_type : '#(offer)', product: 'maquininha'}

        Given path '/offers'
        And header Content-Type = 'application/json'
        And request offerBody
        When method POST

        * def ruleSetConfig = { bureau: {} }

       * def catalog = call CreateSingleLevelCatalog { offer_type: '#(offer)', role_type: '#(RoleType.Customer)', person_type: '#(ProfileType.Individual)', account_flag: false, rules_config: #(ruleSetConfig)}

        Given url mockURL
        And path '/temis-config/cadastral-validation-configs'
        And header Content-Type = 'application/json'
        And request catalog
        When method POST

        Given url mockURL
        And path '/temis-config/cadastral-validation-configs'
        And header Content-Type = 'application/json'
        And def catalog = call CreateSingleLevelCatalog { offer_type: '#(offer)', role_type: '#(RoleType.LegalRepresentative)', person_type: '#(ProfileType.Individual)', account_flag: false,rules_config: #(ruleSetConfig)}
        And request catalog
        When method POST

        Given url mockURL
        And path '/temis-config/cadastral-validation-configs'
        And header Content-Type = 'application/json'
        And def catalog = call CreateSingleLevelCatalog { offer_type: '#(offer)', role_type: '#(RoleType.Shareholder)', person_type: '#(ProfileType.Individual)', account_flag: false, rules_config: #(ruleSetConfig)}
        And request catalog
        When method POST

    Scenario: should send an event to temis query
		* def documentNumber =  CPFGenerator()
		* def partnerID = uuid()
		* def profileID = uuid()
		* def profile =
        """
            {
                profile_id: '#(profileID)',
                partner_id: '#(partnerID)',
                offer_type: '#(offer)',
                callback_url: '/url',
                role_type: '#(RoleType.Customer)',
                profile_type: '#(ProfileType.Individual)',
                document_number: '#(documentNumber)',
                individual: {
                    first_name: 'Daniel',
                    last_name: 'Alves',
                    date_of_birth: '1991-05-06T00:00:00Z',
                    nationality: 'BRA',
                    pep: true
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

		Given url mockURL
		And path '/v1/temis/partner/' + partnerID
		And request partnerResponse
		When method POST
		Then assert responseStatus == 200

		* def event = { profile_id: '#(profileID)', entity_id: '#(profileID)', entity_type: 'PROFILE', event_type: 'PROFILE_CHANGED', parent_type: 'PROFILE', update_date: '2023-08-03T19:29:06.63556Z'  }
		* string json = event
		* def result = RegistrationEventsPublisher.publish(json)
		* eval sleep(1000)
