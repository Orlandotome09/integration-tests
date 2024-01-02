Feature: Test Board Of Directors Rule

    Background:
        * url baseURLCompliance

        * def boardOfDirectorsOffer = "TEST_OFFER" + uuid()  
        * def problemCodeDirectorNotApproved = "DIRECTOR_NOT_APPROVED"

        * def offer = { offer_type : '#(boardOfDirectorsOffer)', product: 'maquininha'}

        Given path '/offers'
        And header Content-Type = 'application/json'
        And request offer
        When method POST

        * def ruleSetConfig = { board_of_directors: {} }

        Given url mockURL
        And path '/temis-config/cadastral-validation-configs'
        And header Content-Type = 'application/json'
        And def catalog = call CreateSingleLevelCatalog { offer_type: '#(boardOfDirectorsOffer)', role_type: '#(RoleType.Customer)', person_type: '#(ProfileType.Company)', account_flag: false, rules_config: #(ruleSetConfig) }
        And catalog.product_config.enrich_profile_with_bureau_data = true
        And request catalog
        When method POST	

        * def ruleSetConfig = { bureau: {} }

        Given url mockURL
        And path '/temis-config/cadastral-validation-configs'
        And header Content-Type = 'application/json'
        And def catalog = call CreateSingleLevelCatalog { offer_type: '#(boardOfDirectorsOffer)', role_type: '#(RoleType.Director)', person_type: '#(ProfileType.Individual)', account_flag: false, rules_config: #(ruleSetConfig)}
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

        * def expectedDirector1CreatedStateEvent = 
		"""
		{
			"id": "#string",
			"entity_id": "#(directorID1)",
			"entity_type": "#(EntityType.ComplianceState)",
			"event_type": "#(EventType.State.Created)",
			"update_date":"#string",
			"data": {
                "profile_id":  "#(profileID)",
				"content": "#(personStateEventContentSchema)",
			}
		}	
		"""
        
		* def expectedDirector1ChangedStateEvent = 
		"""
		{
			"id": "#string",
			"entity_id": "#(directorID1)",
			"entity_type": "#(EntityType.ComplianceState)",
			"event_type": "#(EventType.State.Changed)",
			"update_date":"#string",
			"data": {
                "profile_id":  "#(profileID)",
				"content": "#(personStateEventContentSchema)",
			}
		}	
		"""

        * def expectedDirector2CreatedStateEvent = 
        """
        {
            "id": "#string",
            "entity_id": "#(directorID2)",
            "entity_type": "#(EntityType.ComplianceState)",
            "event_type": "#(EventType.State.Created)",
            "update_date":"#string",
            "data": {
                "profile_id":  "#(profileID)",
                "content": "#(personStateEventContentSchema)",
            }
        }	
        """
        
        * def expectedDirector2ChangedStateEvent = 
        """
        {
            "id": "#string",
            "entity_id": "#(directorID2)",
            "entity_type": "#(EntityType.ComplianceState)",
            "event_type": "#(EventType.State.Changed)",
            "update_date":"#string",
            "data": {
                "profile_id":  "#(profileID)",
                "content": "#(personStateEventContentSchema)",
            }
        }	
        """
        
    #
    # Scenarios with Directors registered at temis-registration
    #
    Scenario: Rule result should be set as Rejected if at least one director is not approved in Bureau 

        * def documentNumber = DocumentNormalizer(CNPJGenerator())
        * def profile = {profile_id:'#(profileID)', document_number: '#(documentNumber)', offer_type: '#(boardOfDirectorsOffer)', role_type: '#(RoleType.Customer)', profile_type: '#(ProfileType.Company)', company: {legal_nature: "2143"}}

        Given url mockURL
        And path '/v1/temis/profile/' + profileID
        And request profile
        When method POST
        Then assert responseStatus == 200

        * def directorID1 = uuid()
        * def directorID2 = uuid()
        * def directors =  [{director_id : '#(directorID1)',document_number: '#(CPFGenerator())', profile_id: '#(profileID)'},{director_id : '#(directorID2)',document_number: '#(CPFGenerator())', date_of_birth: "1990-01-01", profile_id: '#(profileID)'}]

        Given url mockURL
        And path '/v1/temis/directors'
        And params { profile_id : '#(profileID)' }
        And request directors
        When method POST

        * def event = { profile_id: '#(profileID)', entity_id: '#(profileID)', entity_type: 'PROFILE', event_type: 'PROFILE_CHANGED', parent_type: 'PROFILE', update_date: '2023-08-03T19:29:06.63556Z'  }
        * string json = event

        * def result = RegistrationEventsPublisher.publish(json)

        Given url baseURLCompliance
        And path '/state/' + profileID
        And header Content-Type = 'application/json'
        And retry until response.result == RuleResult.Rejected
        When method GET
        And match response.rule_set_result contains deep { set:"#(RuleSet.BoardOfDirectors)", name:"#(RuleName.BoardOfDirectorsComplete)",result:"#(RuleResult.Approved)", pending: false }
        And match response.rule_set_result contains deep { set:"#(RuleSet.BoardOfDirectors)", name:"#(RuleName.BoardOfDirectorsResult)",result:"#(RuleResult.Rejected)", pending: false }

        Given url baseURLCompliance
        And path '/state/' + directorID1
        And header Content-Type = 'application/json'
        When method GET
        And match response.rule_set_result contains deep { set:"SERASA_BUREAU", name:"CUSTOMER_NOT_FOUND_IN_SERASA",result:"#(RuleResult.Rejected)", pending: false }

        Given url baseURLCompliance
        And path '/state/' + directorID2
        And header Content-Type = 'application/json'
        When method GET
        And match response.rule_set_result contains deep { set:"SERASA_BUREAU", name:"CUSTOMER_NOT_FOUND_IN_SERASA",result:"#(RuleResult.Rejected)", pending: false }

 
        Given url mockURL
		And path "/subscribe/state-events/" + profileID
		And retry until response.length == 2
		When method GET
		Then assert responseStatus == 200
		And match response contains deep expectedProfileCreatedStateEvent
		And match response contains deep expectedProfileChangedStateEvent

        Given url mockURL
		And path "/subscribe/state-events/" + directorID1
		And retry until response.length == 2
		When method GET
		Then assert responseStatus == 200
		And match response contains deep expectedDirector1CreatedStateEvent
		And match response contains deep expectedDirector1ChangedStateEvent

        Given url mockURL
		And path "/subscribe/state-events/" + directorID2
		And retry until response.length == 2
		When method GET
		Then assert responseStatus == 200
		And match response contains deep expectedDirector2CreatedStateEvent
		And match response contains deep expectedDirector2ChangedStateEvent
        
    Scenario: Rule result should be set as Ignored if Legal Nature is not 2143,2046,1210,3204,3212,3999

        * def documentNumber = DocumentNormalizer(CNPJGenerator())
        * def profile = {profile_id:'#(profileID)', document_number: '#(documentNumber)', offer_type: '#(boardOfDirectorsOffer)', role_type: '#(RoleType.Customer)', profile_type: '#(ProfileType.Company)', company: {legal_nature: "2000"}}

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
        And retry until response.result == RuleResult.Approved
        When method GET
        And match response.rule_set_result contains deep { set:"#(RuleSet.BoardOfDirectors)", name:"#(RuleName.BoardOfDirectorsComplete)",result:"#(RuleResult.Ignored)", pending: false }
        And match response.rule_set_result contains deep { set:"#(RuleSet.BoardOfDirectors)", name:"#(RuleName.BoardOfDirectorsResult)",result:"#(RuleResult.Ignored)", pending: false }

        Given url mockURL
		And path "/subscribe/state-events/" + profileID
		And retry until response.length == 2
		When method GET
		Then assert responseStatus == 200
		And match response contains deep expectedProfileCreatedStateEvent
		And match response contains deep expectedProfileChangedStateEvent

    @CreateProfileWithRegisteredDirectors
    Scenario: Rule result should be set as Approved if has at least one Director and it is approved

        * def documentNumber = DocumentNormalizer(CNPJGenerator())
        * def profile = {profile_id:'#(profileID)', document_number: '#(documentNumber)', offer_type: '#(boardOfDirectorsOffer)', role_type: '#(RoleType.Customer)', profile_type: '#(ProfileType.Company)', company: {legal_nature: "2143"}}

        Given url mockURL
        And path '/v1/temis/profile/' + profileID
        And request profile
        When method POST
        Then assert responseStatus == 200

        * def directorID1 = uuid()
        * def directorID2 = uuid()
        * def directors =  [{director_id : '#(directorID1)',document_number: '#(CPFGenerator())', profile_id: '#(profileID)'},{director_id : '#(directorID2)',document_number: '#(CPFGenerator())', profile_id: '#(profileID)'}]

        Given url mockURL
        And path '/v1/temis/directors'
        And params { profile_id : '#(profileID)' }
        And request directors
        When method POST

        * def enrichmentResponse = { situation: 1 }

        Given url mockURL
        And path '/temis-enrichment/individual/' + directors[0].document_number
        And request enrichmentResponse
        When method POST

        * def enrichmentResponse = { situation: 1 }

        Given url mockURL
        And path '/temis-enrichment/individual/' + directors[1].document_number
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

        Given url baseURLCompliance
        And path '/state/' + directorID1
        And header Content-Type = 'application/json'
        When method GET
        And match response.rule_set_result contains deep { set:"SERASA_BUREAU", name:"CUSTOMER_NOT_FOUND_IN_SERASA",result:"#(RuleResult.Approved)", pending: false }

        Given url baseURLCompliance
        And path '/state/' + directorID2
        And header Content-Type = 'application/json'
        When method GET
        And match response.rule_set_result contains deep { set:"SERASA_BUREAU", name:"CUSTOMER_NOT_FOUND_IN_SERASA",result:"#(RuleResult.Approved)", pending: false }

        Given url mockURL
		And path "/subscribe/state-events/" + profileID
		And retry until response.length == 2
		When method GET
		Then assert responseStatus == 200
		And match response contains deep expectedProfileCreatedStateEvent
		And match response contains deep expectedProfileChangedStateEvent

        Given url mockURL
		And path "/subscribe/state-events/" + directorID1
		And retry until response.length == 2
		When method GET
		Then assert responseStatus == 200
		And match response contains deep expectedDirector1CreatedStateEvent
		And match response contains deep expectedDirector1ChangedStateEvent

        Given url mockURL
		And path "/subscribe/state-events/" + directorID2
		And retry until response.length == 2
		When method GET
		Then assert responseStatus == 200
		And match response contains deep expectedDirector2CreatedStateEvent
		And match response contains deep expectedDirector2ChangedStateEvent

    #
    # Scenarios with enriched Directors
    #
    Scenario: Rule result should be set as Rejected if at least one enriched director is not approved in Bureau

        * def documentNumber = DocumentNormalizer(CNPJGenerator())
        * def profile = {profile_id:'#(profileID)', document_number: '#(documentNumber)', offer_type: '#(boardOfDirectorsOffer)', role_type: '#(RoleType.Customer)', profile_type: '#(ProfileType.Company)', company: {legal_nature: "2143"}}

        Given url mockURL
        And path '/v1/temis/profile/' + profileID
        And request profile
        When method POST
        Then assert responseStatus == 200
        
        * def directors =  [{document_number: '#(CPFGenerator())'},{document_number: '#(CPFGenerator())', profile_id: '#(profileID)'}]

        * def enrichmentResponse = { legal_nature: '2143', board_of_directors: [ { document_number: '#(directors[0].document_number)' }, { document_number: '#(directors[1].document_number)'} ] }

        Given url mockURL
        And path '/temis-enrichment/legal-entity/' + documentNumber
        And request enrichmentResponse
        When method POST

        * def event = { profile_id: '#(profileID)', entity_id: '#(profileID)', entity_type: 'PROFILE', event_type: 'PROFILE_CHANGED', parent_type: 'PROFILE', update_date: '2023-08-03T19:29:06.63556Z'  }
        * string json = event
        
        * def result = RegistrationEventsPublisher.publish(json)

        Given url baseURLCompliance
        And path '/state/' + profileID
        And header Content-Type = 'application/json'
        And retry until response.result == RuleResult.Rejected
        When method GET
        And match response.rule_set_result contains deep { set:"#(RuleSet.BoardOfDirectors)", name:"#(RuleName.BoardOfDirectorsComplete)",result:"#(RuleResult.Approved)", pending: false }
        And match response.rule_set_result contains deep { set:"#(RuleSet.BoardOfDirectors)", name:"#(RuleName.BoardOfDirectorsResult)",result:"#(RuleResult.Rejected)", pending: false }

        * def directorID1 = response.rule_set_result[1].problems[0].detail[0]
        * def directorID2 = response.rule_set_result[1].problems[0].detail[1]

        Given url baseURLCompliance
        And path '/state/' + directorID1
        And header Content-Type = 'application/json'
        When method GET
        And match response.rule_set_result contains deep { set:"SERASA_BUREAU", name:"CUSTOMER_NOT_FOUND_IN_SERASA",result:"#(RuleResult.Rejected)", pending: false }

        Given url baseURLCompliance
        And path '/state/' + directorID2
        And header Content-Type = 'application/json'
        When method GET
        And match response.rule_set_result contains deep { set:"SERASA_BUREAU", name:"CUSTOMER_NOT_FOUND_IN_SERASA",result:"#(RuleResult.Rejected)", pending: false }

        Given url mockURL
		And path "/subscribe/state-events/" + profileID
		And retry until response.length == 2
        When method GET
		Then assert responseStatus == 200
		And match response contains deep expectedProfileCreatedStateEvent
		And match response contains deep expectedProfileChangedStateEvent

        Given url mockURL
		And path "/subscribe/state-events/" + directorID1
		And retry until response.length == 2
		When method GET
		Then assert responseStatus == 200
		And match response contains deep expectedDirector1CreatedStateEvent
		And match response contains deep expectedDirector1ChangedStateEvent

        Given url mockURL
		And path "/subscribe/state-events/" + directorID2
		And retry until response.length == 2
		When method GET
		Then assert responseStatus == 200
		And match response contains deep expectedDirector2CreatedStateEvent
		And match response contains deep expectedDirector2ChangedStateEvent
        
    Scenario: Rule result should be set as Ignored if enriched Legal Nature is not 2143,2046,1210,3204,3212,3999

        * def documentNumber = DocumentNormalizer(CNPJGenerator())
        * def profile = {profile_id:'#(profileID)', document_number: '#(documentNumber)', offer_type: '#(boardOfDirectorsOffer)', role_type: '#(RoleType.Customer)', profile_type: '#(ProfileType.Company)', company: {legal_nature: "2000"}}

        Given url mockURL
        And path '/v1/temis/profile/' + profileID
        And request profile
        When method POST
        Then assert responseStatus == 200

        * def directors =  [{document_number: '#(CPFGenerator())'},{document_number: '#(CPFGenerator())'}]

        * def enrichmentResponse = { legal_nature: '200', board_of_directors: [ { document_number: '#(directors[0].document_number)' }, { document_number: '#(directors[0].document_number)'} ] }

        Given url mockURL
        And path '/temis-enrichment/legal-entity/' + documentNumber
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
        And match response.rule_set_result contains deep { set:"#(RuleSet.BoardOfDirectors)", name:"#(RuleName.BoardOfDirectorsComplete)",result:"#(RuleResult.Ignored)", pending: false }
        And match response.rule_set_result contains deep { set:"#(RuleSet.BoardOfDirectors)", name:"#(RuleName.BoardOfDirectorsResult)",result:"#(RuleResult.Ignored)", pending: false }

        Given url mockURL
		And path "/subscribe/state-events/" + profileID
		And retry until response.length == 2
		When method GET
		Then assert responseStatus == 200
		And match response contains deep expectedProfileCreatedStateEvent
		And match response contains deep expectedProfileChangedStateEvent

    @CreateProfileWithEnrichedDirectors
    Scenario: Rule result should be set as Approved if all enriched Directors are approved

        * def documentNumber = DocumentNormalizer(CNPJGenerator())
        * def profile = {profile_id:'#(profileID)', document_number: '#(documentNumber)', offer_type: '#(boardOfDirectorsOffer)', role_type: '#(RoleType.Customer)', profile_type: '#(ProfileType.Company)', company: {legal_nature: "2143"}}

        Given url mockURL
        And path '/v1/temis/profile/' + profileID
        And request profile
        When method POST
        Then assert responseStatus == 200

        * def directors =  [{document_number: '#(CPFGenerator())'},{document_number: '#(CPFGenerator())'}]

        * def enrichmentResponse = { legal_nature: '2143', board_of_directors: [ { document_number: '#(directors[0].document_number)' }, { document_number: '#(directors[1].document_number)'} ] }

        Given url mockURL
        And path '/temis-enrichment/legal-entity/' + documentNumber
        And request enrichmentResponse
        When method POST

        * def enrichmentResponse = { situation: 1 }

        Given url mockURL
        And path '/temis-enrichment/individual/' + directors[0].document_number
        And request enrichmentResponse
        When method POST

        * def enrichmentResponse = { situation: 1 }

        Given url mockURL
        And path '/temis-enrichment/individual/' + directors[1].document_number
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

        Given url mockURL
		And path "/subscribe/state-events/" + profileID
		And retry until response.length == 2
		When method GET
		Then assert responseStatus == 200
		And match response contains deep expectedProfileCreatedStateEvent
		And match response contains deep expectedProfileChangedStateEvent

        Given url baseURLCompliance
        And path "/profile/" + profileID
        And header Content-Type = 'application/json'
        When method GET
        Then assert responseStatus == 200

        * def profile = response
        * def directorID1 = profile.person.enriched_information.board_of_directors[0].director_id
        * def directorID2 = profile.person.enriched_information.board_of_directors[1].director_id

        Given url mockURL
		And path "/subscribe/state-events/" + directorID1
		And retry until response.length == 2
		When method GET
		Then assert responseStatus == 200
		And match response contains deep expectedDirector1CreatedStateEvent
		And match response contains deep expectedDirector1ChangedStateEvent

        Given url mockURL
		And path "/subscribe/state-events/" + directorID2
		And retry until response.length == 2
		When method GET
		Then assert responseStatus == 200
		And match response contains deep expectedDirector2CreatedStateEvent
		And match response contains deep expectedDirector2ChangedStateEvent

    #
    # Mixes Scenarios (with registered and enriched information)
    #
    @CreateProfileWithEnrichedAndRegisteredDirectors
    Scenario: Rule result should be set as Approved if all enriched Directors are approved (even when there are registered directors not found at bureau)

        * def documentNumber = DocumentNormalizer(CNPJGenerator())
        * def profile = {profile_id:'#(profileID)', document_number: '#(documentNumber)', offer_type: '#(boardOfDirectorsOffer)', role_type: '#(RoleType.Customer)', profile_type: '#(ProfileType.Company)', company: {legal_nature: "2143"}}

        Given url mockURL
        And path '/v1/temis/profile/' + profileID
        And request profile
        When method POST
        Then assert responseStatus == 200

        # Registered directors would not be approved if analyzed, since they do not exists on Bureau
        * def directorID1 = uuid()
        * def directorID2 = uuid()
        * def directors =  [{director_id : '#(directorID1)',document_number: '#(CPFGenerator())', profile_id: '#(profileID)'},{director_id : '#(directorID2)',document_number: '#(CPFGenerator())', profile_id: '#(profileID)'}]

        Given url mockURL
        And path '/v1/temis/directors'
        And params { profile_id : '#(profileID)' }
        And request directors
        When method POST

        # Enriched directors would be approved if analyzed, since they do exists on Bureau
        * def directors =  [{document_number: '#(CPFGenerator())'},{document_number: '#(CPFGenerator())'}]
        * def enrichmentResponse = { legal_nature: '2143', board_of_directors: [ { document_number: '#(directors[0].document_number)' }, { document_number: '#(directors[1].document_number)'} ] }

        Given url mockURL
        And path '/temis-enrichment/legal-entity/' + documentNumber
        And request enrichmentResponse
        When method POST

        * def enrichmentResponse = { situation: 1 }

        Given url mockURL
        And path '/temis-enrichment/individual/' + directors[0].document_number
        And request enrichmentResponse
        When method POST

        * def enrichmentResponse = { situation: 1 }

        Given url mockURL
        And path '/temis-enrichment/individual/' + directors[1].document_number
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

        Given url mockURL
		And path "/subscribe/state-events/" + profileID
		And retry until response.length == 2
		When method GET
		Then assert responseStatus == 200
		And match response contains deep expectedProfileCreatedStateEvent
		And match response contains deep expectedProfileChangedStateEvent

        Given url baseURLCompliance
        And path "/profile/" + profileID
        And header Content-Type = 'application/json'
        When method GET
        Then assert responseStatus == 200

        * def profile = response
        * def directorID1 = profile.person.enriched_information.board_of_directors[0].director_id
        * def directorID2 = profile.person.enriched_information.board_of_directors[1].director_id

        Given url mockURL
		And path "/subscribe/state-events/" + directorID1
		And retry until response.length == 2
		When method GET
		Then assert responseStatus == 200
		And match response contains deep expectedDirector1CreatedStateEvent
		And match response contains deep expectedDirector1ChangedStateEvent

        Given url mockURL
		And path "/subscribe/state-events/" + directorID2
		And retry until response.length == 2
		When method GET
		Then assert responseStatus == 200
		And match response contains deep expectedDirector2CreatedStateEvent
		And match response contains deep expectedDirector2ChangedStateEvent

    Scenario: Rule result should be set as Analysing if has no directors (in both registered and enriched info)

        * def documentNumber = DocumentNormalizer(CNPJGenerator())
        * def profile = {profile_id:'#(profileID)', document_number: '#(documentNumber)', offer_type: '#(boardOfDirectorsOffer)', role_type: '#(RoleType.Customer)', profile_type: '#(ProfileType.Company)', company: {legal_nature: "2143"}}

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
        And retry until response.result == RuleResult.Analysing
        When method GET
        And match response.rule_set_result contains deep { set:"#(RuleSet.BoardOfDirectors)", name:"#(RuleName.BoardOfDirectorsComplete)",result:"#(RuleResult.Analysing)", pending: true }
        And match response.rule_set_result contains deep { set:"#(RuleSet.BoardOfDirectors)", name:"#(RuleName.BoardOfDirectorsResult)",result:"#(RuleResult.Ignored)", pending: false }

        Given url mockURL
		And path "/subscribe/state-events/" + profileID
		And retry until response.length == 2
		When method GET
		Then assert responseStatus == 200
		And match response contains deep expectedProfileCreatedStateEvent
		And match response contains deep expectedProfileChangedStateEvent

    Scenario: Rule result should be set as Ignored if Legal Nature is empty (in both registered and enriched info)

        * def documentNumber = DocumentNormalizer(CNPJGenerator())
        * def profile = {profile_id:'#(profileID)', document_number: '#(documentNumber)', offer_type: '#(boardOfDirectorsOffer)', role_type: '#(RoleType.Customer)', profile_type: '#(ProfileType.Company)', company: {legal_nature: ""}}

        Given url mockURL
        And path '/v1/temis/profile/' + profileID
        And request profile
        When method POST
        Then assert responseStatus == 200

        * def directors =  [{document_number: '#(CPFGenerator())'},{document_number: '#(CPFGenerator())'}]

        * def enrichmentResponse = { legal_nature: '', board_of_directors: [ { document_number: '#(directors[0].document_number)' }, { document_number: '#(directors[0].document_number)'} ] }

        Given url mockURL
        And path '/temis-enrichment/legal-entity/' + documentNumber
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
        And match response.rule_set_result contains deep { set:"#(RuleSet.BoardOfDirectors)", name:"#(RuleName.BoardOfDirectorsComplete)",result:"#(RuleResult.Ignored)", pending: false }
        And match response.rule_set_result contains deep { set:"#(RuleSet.BoardOfDirectors)", name:"#(RuleName.BoardOfDirectorsResult)",result:"#(RuleResult.Ignored)", pending: false }

        Given url mockURL
		And path "/subscribe/state-events/" + profileID
		And retry until response.length == 2
		When method GET
		Then assert responseStatus == 200
		And match response contains deep expectedProfileCreatedStateEvent
		And match response contains deep expectedProfileChangedStateEvent

        Given url baseURLCompliance
        And path "/profile/" + profileID
        And header Content-Type = 'application/json'
        When method GET
        Then assert responseStatus == 200
