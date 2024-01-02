Feature: OVERRIDE API

  	Background:
    * url baseURLCompliance
    
	* def ruleSetLegalRepresentatives = "LEGAL_REPRESENTATIVES"
	* def ruleNameLegalRepresentativeResult = "LEGAL_REPRESENTATIVES_RESULT"
    * def ruleSetSerasaBureau = "SERASA_BUREAU"
    * def ruleNameCustomerHasProblemsInSerasa = "CUSTOMER_HAS_PROBLEMS_IN_SERASA"

    Scenario: Should successfully create an override     
   
        * def offerType = "TEST_OFFER" + uuid()
        * def ruleSetConfig = { bureau: {} }
        * def cadastralValidationConfigs = call CreateSingleLevelCatalog { offer_type: '#(offerType)', role_type: '#(RoleType.Customer)', person_type: '#(ProfileType.Individual)', account_flag: false,  rules_config: #(ruleSetConfig)}

        Given url baseURLCompliance
        And path '/offers'
        And header Content-Type = 'application/json'
        And def offer = {offer_type: '#(offerType)', product: 'maquininha'}
        And request offer
        When method POST
        Then assert responseStatus == 201

        Given url mockURL
        And path '/temis-config/cadastral-validation-configs'
        And header Content-Type = 'application/json'
        And request cadastralValidationConfigs
        When method POST
        Then assert responseStatus == 200
        
        * def documentNumber =  CPFGenerator()
        * def profileID = uuid()
        * def profile = { profile_id: '#(profileID)', offer_type: '#(offerType)', role_type: '#(RoleType.Customer)', profile_type: '#(ProfileType.Individual)', document_number: '#(documentNumber)' }

        Given url mockURL
        And path '/v1/temis/profile/' + profileID
        And request profile
        When method POST
        Then assert responseStatus == 200

        * def enrichmentResponse = { situation: 4 }

        Given url mockURL
        And path '/temis-enrichment/individual/' + documentNumber
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
        Then assert responseStatus == 200
        * def ruleResults = response.rule_set_result
        And assert ruleResults[1].set == ruleSetSerasaBureau
        And assert ruleResults[1].name == ruleNameCustomerHasProblemsInSerasa

        * def override = {entity_id: '#(profileID)', entity_type: '#(Engine.Profile)', rule_set: '#(ruleSetSerasaBureau)', rule_name: '#(ruleNameCustomerHasProblemsInSerasa)', result: '#(RuleResult.Approved)', comments: 'it is approved now', author: 'Henrique'}
    
        Given url baseURLCompliance
        And path '/override'
        And header Content-Type = 'application/json'
        And request override
        When method POST
        Then assert responseStatus == 200

        Given url baseURLCompliance
		And path '/state/' + profileID
        And header Content-Type = 'application/json'
        And params ({ only_pending: false })		
        When method GET
        Then assert responseStatus == 200
        * def ruleResults = response.rule_set_result
        And assert ruleResults[1].set == ruleSetSerasaBureau
        And assert ruleResults[1].name == ruleNameCustomerHasProblemsInSerasa
		
    Scenario: Should create and delete an override

        * def offerType = "TEST_OFFER" + uuid()
        * def ruleSetConfig = { bureau: {} }
        * def cadastralValidationConfigs = call CreateSingleLevelCatalog { offer_type: '#(offerType)', role_type: '#(RoleType.Customer)', person_type: '#(ProfileType.Individual)', account_flag: false,  rules_config: #(ruleSetConfig)}

        Given url baseURLCompliance
        And path '/offers'
        And header Content-Type = 'application/json'
        And def offer = {offer_type: '#(offerType)', product: 'maquininha'}
        And request offer
        When method POST
        Then assert responseStatus == 201
        
        Given url mockURL
        And path '/temis-config/cadastral-validation-configs'
        And header Content-Type = 'application/json'
        And request cadastralValidationConfigs
        When method POST
        Then assert responseStatus == 200

        * def documentNumber =  CPFGenerator()
        * def profileID = uuid()
        * def profile = { profile_id: '#(profileID)', offer_type: '#(offerType)', role_type: '#(RoleType.Customer)', profile_type: '#(ProfileType.Individual)', document_number: '#(documentNumber)' }

        Given url mockURL
        And path '/v1/temis/profile/' + profileID
        And request profile
        When method POST
        Then assert responseStatus == 200

        * def override = {entity_id: '#(profileID)', entity_type: '#(Engine.Profile)', rule_set: '#(ruleSetSerasaBureau)', rule_name: '#(ruleNameCustomerHasProblemsInSerasa)', result: '#(RuleResult.Approved)', comments: 'Will be deleted', author: 'Henrique'}

        Given url baseURLCompliance
        And path '/override'
        And header Content-Type = 'application/json'
        And request override
        When method POST
        Then assert responseStatus == 200

        Given url baseURLCompliance
        And path '/override'
        And header Content-Type = 'application/json'
        And request override
        When method DELETE
        Then assert responseStatus == 200
		
    Scenario: Should override a legal representative for approving profile   
        
        * def offerType = "TEST_OFFER" + uuid()
        * def ruleSetConfig = { legal_representative: {} }
        * def cadastralValidationConfigs = call CreateSingleLevelCatalog { offer_type: '#(offerType)', role_type: '#(RoleType.Customer)', person_type: '#(ProfileType.Company)', account_flag: false,  rules_config: #(ruleSetConfig)}

        Given url baseURLCompliance
        And path '/offers'
        And header Content-Type = 'application/json'
        And def offer = {offer_type: '#(offerType)', product: 'maquininha'}
        And request offer
        When method POST
        Then assert responseStatus == 201

        Given url mockURL
        And path '/temis-config/cadastral-validation-configs'
        And header Content-Type = 'application/json'
        And request cadastralValidationConfigs
        When method POST
        Then assert responseStatus == 200

        * def ruleSetConfigBureau = { bureau: {} }
        * def cadastralValidationConfigs = call CreateSingleLevelCatalog { offer_type: '#(offerType)', role_type: '#(RoleType.LegalRepresentative)', person_type: '#(ProfileType.Individual)', account_flag: false,  rules_config: #(ruleSetConfigBureau)}

        Given url mockURL
        And path '/temis-config/cadastral-validation-configs'
        And header Content-Type = 'application/json'
        And request cadastralValidationConfigs
        When method POST
        Then assert responseStatus == 200

        * def companyProfileID = uuid()
		* def companyDocumentNumber = DocumentNormalizer(CNPJGenerator())
		* def companyProfile = {profile_id:'#(companyProfileID)', offer_type: '#(offerType)', document_number: '#(companyDocumentNumber)', role_type: '#(RoleType.Customer)', profile_type: '#(ProfileType.Company)'}

		Given url mockURL
		And path '/v1/temis/profile/' + companyProfileID
		And request companyProfile
		When method POST
		Then assert responseStatus == 200

        * def firstLegalRepresentativeProfileID = uuid()
		* def firstLegalRepresentativeDocumentNumber = CPFGenerator()
        * def secondLegalRepresentativeProfileID = uuid()
		* def secondLegalRepresentativeDocumentNumber = CPFGenerator()
		* def legalRepresentativesResponse = [ {legal_representative_id: '#(firstLegalRepresentativeProfileID)', profile_id: '#(companyProfileID)', document_number: '#(firstLegalRepresentativeDocumentNumber)'}, {legal_representative_id: '#(secondLegalRepresentativeProfileID)', profile_id: '#(companyProfileID)', document_number: '#(secondLegalRepresentativeDocumentNumber)'} ]

		Given url mockURL
		And path '/v1/temis/legal-representatives'
		And params { profile_id : '#(companyProfileID)' }
		And request legalRepresentativesResponse
		When method POST
        Then assert responseStatus == 200

        * def enrichmentResponse = { situation: 1 }
        
        Given url mockURL
		And path '/temis-enrichment/individual/' + firstLegalRepresentativeDocumentNumber
		And request enrichmentResponse
		When method POST
        Then assert responseStatus == 200

        * def enrichmentResponse = { situation: 4 }
        
        Given url mockURL
		And path '/temis-enrichment/individual/' + secondLegalRepresentativeDocumentNumber
		And request enrichmentResponse
		When method POST
        Then assert responseStatus == 200

        * def event = { profile_id: '#(companyProfileID)', entity_id: '#(companyProfileID)', entity_type: 'PROFILE', event_type: 'PROFILE_CHANGED', parent_type: 'PROFILE', update_date: '2023-08-03T19:29:06.63556Z'  }
		* string json = event
    	* def result = RegistrationEventsPublisher.publish(json)
		* eval sleep(1000)

        Given url baseURLCompliance
		And path '/state/' + companyProfileID
        And header Content-Type = 'application/json'
        And params ({ only_pending: false })		
		And retry until response.result == RuleResult.Rejected
        When method GET
        Then assert responseStatus == 200

        * def override = {entity_id: '#(secondLegalRepresentativeProfileID)', entity_type: '#(Engine.Profile)', parent_id: '#(companyProfileID)', rule_set: '#(ruleSetSerasaBureau)', rule_name: '#(ruleNameCustomerHasProblemsInSerasa)', result: '#(RuleResult.Approved)', comments: 'it is approved now', author: 'Henrique'}
    
        Given url baseURLCompliance
        And path '/override'
        And header Content-Type = 'application/json'
        And request override
        When method POST
        Then assert responseStatus == 200

        Given url baseURLCompliance
		And path '/state/' + secondLegalRepresentativeProfileID
        And header Content-Type = 'application/json'
        And params ({ only_pending: false })		
        When method GET
        Then assert responseStatus == 200
        * def ruleResults = response.rule_set_result
        And assert ruleResults[1].set == ruleSetSerasaBureau
        And assert ruleResults[1].name == ruleNameCustomerHasProblemsInSerasa

        Given url baseURLCompliance
		And path '/state/' + companyProfileID
        And header Content-Type = 'application/json'
        And params ({ only_pending: false })		
        When method GET
        Then assert responseStatus == 200
        * def ruleResults = response.rule_set_result
        And assert ruleResults[0].set == ruleSetLegalRepresentatives
        And assert ruleResults[0].name == ruleNameLegalRepresentativeResult

    Scenario: should not Override when missing required fields     
   
        * def offerType = "TEST_OFFER" + uuid()
        * def ruleSetConfig = { bureau: {} }
        * def cadastralValidationConfigs = call CreateSingleLevelCatalog { offer_type: '#(offerType)', role_type: '#(RoleType.Customer)', person_type: '#(ProfileType.Individual)', account_flag: false,  rules_config: #(ruleSetConfig)}

        Given url baseURLCompliance
        And path '/offers'
        And header Content-Type = 'application/json'
        And def offer = {offer_type: '#(offerType)', product: 'maquininha'}
        And request offer
        When method POST
        Then assert responseStatus == 201

        Given url mockURL
        And path '/temis-config/cadastral-validation-configs'
        And header Content-Type = 'application/json'
        And request cadastralValidationConfigs
        When method POST
        Then assert responseStatus == 200
        
        * def documentNumber =  CPFGenerator()
        * def profileID = uuid()
        * def profile = { profile_id: '#(profileID)', offer_type: '#(offerType)', role_type: '#(RoleType.Customer)', profile_type: '#(ProfileType.Individual)', document_number: '#(documentNumber)' }

        Given url mockURL
        And path '/v1/temis/profile/' + profileID
        And request profile
        When method POST
        Then assert responseStatus == 200

        * def enrichmentResponse = { situation: 2 }

        Given url mockURL
        And path '/temis-enrichment/individual/' + documentNumber
        And request enrichmentResponse
        When method POST
        Then assert responseStatus == 200

        * def override = {}
        * def expectedError = "EntityID is required, EntityType is required, RuleName is required, RuleSet is required, Result is required"

        Given url baseURLCompliance
        And path '/override'
        And header Content-Type = 'application/json'
        And request override
        When method POST
        Then assert responseStatus == 400
        Then assert response.error == expectedError 