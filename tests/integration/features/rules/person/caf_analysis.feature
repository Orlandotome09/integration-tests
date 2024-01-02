Feature: CAF Rule

    Background:
        * url baseURLCompliance

        * def CAFRuleOffer = "TEST_OFFER" + uuid() 

        * def cafStatusApproved = "APPROVED"
        * def cafStatusRejected = "REJECTED"
        * def cafStatusAnalysing = "ANALYSING"

        * def offer = { offer_type : '#(CAFRuleOffer)', product: 'maquininha' }

        Given path '/offers'
        And header Content-Type = 'application/json'
        And request offer
        When method POST

        * def ruleSetConfig = { caf: {} }

        Given url mockURL
        And path '/temis-config/cadastral-validation-configs'
        And header Content-Type = 'application/json'
        And def catalog = call CreateSingleLevelCatalog { offer_type: '#(CAFRuleOffer)', role_type: '#(RoleType.Customer)', person_type: '#(ProfileType.Individual)', account_flag: false, rules_config: #(ruleSetConfig)}
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

    @CreateProfileEnrichedCAFProvider
    Scenario: should approve when caf approved

        * def partnerID = uuid()
        * def requestID = uuid()
        * def documentNumber = DocumentNormalizer(CPFGenerator())
        * def profile = {profile_id:'#(profileID)', partner_id: "#(partnerID)", document_number: '#(documentNumber)', offer_type: '#(CAFRuleOffer)',  role_type: '#(RoleType.Customer)', profile_type: '#(ProfileType.Individual)'}

        Given url mockURL
        And path '/v1/temis/profile/' + profileID
        And request profile
        When method POST
        Then assert responseStatus == 200

        * def enrichedPersonResponse = 
        """
        {
            "type": "#(ProfileType.Individual)",
            "name": "JOAO DA SILVA",
            "individual":{
                "birth_date": "25/12/1980"
            },
            "providers":[
                {
                    "provider_name": "CAF_ENRICHER",
                    "request_id": "#(requestID)",
                    "status": "#(cafStatusApproved)",
                }
            ]
        }
        """

        Given url mockURL
        And path '/temis-enrichment/enrich/' + documentNumber
        And param profile_id = profileID
        And param person_type = ProfileType.Individual
        And param offer_type = CAFRuleOffer
        And param partner_id = partnerID
        And param role_type = RoleType.Customer
        And header Content-Type = 'application/json'
        And request enrichedPersonResponse
        When method POST
        Then assert responseStatus == 200

        * def event = { id: '#(profileID)'}
        * string json = event
        
        * def result = EnrichmentPublisher.publish(json)

        * def expectedCAFAnalysisResult =
        """
        {
            "set": "#(RuleSet.CafAnalysis)",
            "name": "#(RuleName.CafAnalysis)",
            "result": "#(RuleResult.Approved)"
        }
        """

        Given url baseURLCompliance
        And path '/state/' + profileID
        And header Content-Type = 'application/json'
        And params ({ only_pending: false })		
        And retry until response.result == RuleResult.Approved
        When method GET
        Then assert responseStatus == 200
        And match response.rule_set_result contains deep expectedCAFAnalysisResult

        Given url mockURL
		And path "/subscribe/state-events/" + profileID
		And retry until response.length == 2
		When method GET
		Then assert responseStatus == 200
		And match response contains deep expectedCreatedStateEvent
		And match response contains deep expectedChangedStateEvent

    Scenario: should reject when caf rejected

        * def partnerID = uuid()
        * def requestID = uuid()
        * def documentNumber = DocumentNormalizer(CPFGenerator())
        * def profile = {profile_id:'#(profileID)', partner_id: "#(partnerID)", document_number: '#(documentNumber)', offer_type: '#(CAFRuleOffer)',  role_type: '#(RoleType.Customer)', profile_type: '#(ProfileType.Individual)'}

        Given url mockURL
        And path '/v1/temis/profile/' + profileID
        And request profile
        When method POST
        Then assert responseStatus == 200

        * def enrichedPersonResponse = 
        """
        {
            "providers":[
                {
                    "provider_name": "CAF_ENRICHER",
                    "request_id": "#(requestID)",
                    "status": "#(cafStatusRejected)",
                }
            ]
        }
        """

        Given url mockURL
        And path '/temis-enrichment/enrich/' + documentNumber
        And param profile_id = profileID
        And param person_type = ProfileType.Individual
        And param offer_type = CAFRuleOffer
        And param partner_id = partnerID
        And param role_type = RoleType.Customer
        And header Content-Type = 'application/json'
        And request enrichedPersonResponse
        When method POST
        Then assert responseStatus == 200

        * def event = { id: '#(profileID)'}
        * string json = event
        
        * def result = EnrichmentPublisher.publish(json)

        * def expectedCAFAnalysisResult =
        """
        {
            "set": "#(RuleSet.CafAnalysis)",
            "name": "#(RuleName.CafAnalysis)",
            "result": "#(RuleResult.Rejected)"
        }
        """

        Given url baseURLCompliance
        And path '/state/' + profileID
        And header Content-Type = 'application/json'
        And params ({ only_pending: false })		
        And retry until response.result == RuleResult.Rejected
        When method GET
        Then assert responseStatus == 200
        And match response.rule_set_result contains deep expectedCAFAnalysisResult

        Given url mockURL
		And path "/subscribe/state-events/" + profileID
		And retry until response.length == 2
		When method GET
		Then assert responseStatus == 200
		And match response contains deep expectedCreatedStateEvent
		And match response contains deep expectedChangedStateEvent

    Scenario: should set to analysis when caf analysing

        * def partnerID = uuid()
        * def requestID = uuid()
        * def documentNumber = DocumentNormalizer(CPFGenerator())
        * def profile = {profile_id:'#(profileID)', partner_id: "#(partnerID)", document_number: '#(documentNumber)', offer_type: '#(CAFRuleOffer)',  role_type: '#(RoleType.Customer)', profile_type: '#(ProfileType.Individual)'}

        Given url mockURL
        And path '/v1/temis/profile/' + profileID
        And request profile
        When method POST
        Then assert responseStatus == 200

        * def enrichedPersonResponse = 
        """
        {
            "providers":[
                {
                    "provider_name": "CAF_ENRICHER",
                    "request_id": "#(requestID)",
                    "status": "#(cafStatusAnalysing)",
                }
            ]
        }
        """

        Given url mockURL
        And path '/temis-enrichment/enrich/' + documentNumber
        And param profile_id = profileID
        And param person_type = ProfileType.Individual
        And param offer_type = CAFRuleOffer
        And param partner_id = partnerID
        And param role_type = RoleType.Customer
        And header Content-Type = 'application/json'
        And request enrichedPersonResponse
        When method POST
        Then assert responseStatus == 200

        * def event = { id: '#(profileID)'}
        * string json = event

        * def result = EnrichmentPublisher.publish(json)

        * def expectedCAFAnalysisResult =
        """
        {
            "set": "#(RuleSet.CafAnalysis)",
            "name": "#(RuleName.CafAnalysis)",
            "result": "#(RuleResult.Analysing)",
            "pending": true,
            "problems":[
                {
                    "code": "CAF_ANALYSIS_PENDING",
                    "detail": "##array",
                }
            ]
        }
        """



        Given url baseURLCompliance
        And path '/state/' + profileID
        And header Content-Type = 'application/json'
        And params ({ only_pending: false })		
        And retry until response.result == RuleResult.Analysing
        When method GET
        Then assert responseStatus == 200
        And match response.rule_set_result contains deep expectedCAFAnalysisResult

        Given url mockURL
		And path "/subscribe/state-events/" + profileID
		And retry until response.length == 2
		When method GET
		Then assert responseStatus == 200
		And match response contains deep expectedCreatedStateEvent
		And match response contains deep expectedChangedStateEvent

    Scenario: should set to analysis when not found caf analysis

        * def partnerID = uuid()
        * def requestID = uuid()
        * def documentNumber = DocumentNormalizer(CPFGenerator())
        * def profile = {profile_id:'#(profileID)', partner_id: "#(partnerID)", document_number: '#(documentNumber)', offer_type: '#(CAFRuleOffer)',  role_type: '#(RoleType.Customer)', profile_type: '#(ProfileType.Individual)'}

        Given url mockURL
        And path '/v1/temis/profile/' + profileID
        And request profile
        When method POST
        Then assert responseStatus == 200

        * def enrichedPersonResponse = 
        """
        {
            "providers":[
                {
                    "provider_name": "BUREAU_ENRICHER",
                    "request_id": "#(requestID)",
                    "status": "#(cafStatusAnalysing)",
                }
            ]
        }
        """

        Given url mockURL
        And path '/temis-enrichment/enrich/' + documentNumber
        And param profile_id = profileID
        And param person_type = ProfileType.Individual
        And param offer_type = CAFRuleOffer
        And param partner_id = partnerID
        And param role_type = RoleType.Customer
        And header Content-Type = 'application/json'
        And request enrichedPersonResponse
        When method POST
        Then assert responseStatus == 200

        * def event = { id: '#(profileID)'}
        * string json = event
        
        * def result = EnrichmentPublisher.publish(json)

        * def expectedCAFAnalysisResult =
        """
        {
            "set": "#(RuleSet.CafAnalysis)",
            "name": "#(RuleName.CafAnalysis)",
            "result": "#(RuleResult.Analysing)",
            "pending": true,
            "problems": [
                {
                    "code": "NOT_FOUND_CAF_ANALYSIS",
                    "detail": ""
                }
            ]
        }
        """

        Given url baseURLCompliance
        And path '/state/' + profileID
        And header Content-Type = 'application/json'
        And params ({ only_pending: false })		
        And retry until response.result == RuleResult.Analysing
        When method GET
        Then assert responseStatus == 200
        And match response.rule_set_result contains deep expectedCAFAnalysisResult

        Given url mockURL
		And path "/subscribe/state-events/" + profileID
		And retry until response.length == 2
		When method GET
		Then assert responseStatus == 200
		And match response contains deep expectedCreatedStateEvent
		And match response contains deep expectedChangedStateEvent

    Scenario: should set to analysis when not found enriched information

        * def partnerID = uuid()
        * def requestID = uuid()
        * def documentNumber = DocumentNormalizer(CPFGenerator())
        * def profile = {profile_id:'#(profileID)', partner_id: "#(partnerID)", document_number: '#(documentNumber)', offer_type: '#(CAFRuleOffer)',  role_type: '#(RoleType.Customer)', profile_type: '#(ProfileType.Individual)'}

        Given url mockURL
        And path '/v1/temis/profile/' + profileID
        And request profile
        When method POST
        Then assert responseStatus == 200

        * def event = { id: '#(profileID)'}
        * string json = event
        
        * def result = EnrichmentPublisher.publish(json)

        * def expectedCAFAnalysisResult =
        """
        {
            "set": "#(RuleSet.CafAnalysis)",
            "name": "#(RuleName.CafAnalysis)",
            "result": "#(RuleResult.Analysing)",
            "pending": true,
            "problems": [
                {
                    "code": "NOT_FOUND_ENRICHED_INFORMATION",
                    "detail": ""
                }
            ]
        }
        """

        Given url baseURLCompliance
        And path '/state/' + profileID
        And header Content-Type = 'application/json'
        And params ({ only_pending: false })		
        And retry until response.result == RuleResult.Analysing
        When method GET
        Then assert responseStatus == 200
        And match response.rule_set_result contains deep expectedCAFAnalysisResult

        Given url mockURL
		And path "/subscribe/state-events/" + profileID
		And retry until response.length == 2
		When method GET
		Then assert responseStatus == 200
		And match response contains deep expectedCreatedStateEvent
		And match response contains deep expectedChangedStateEvent

