Feature: ManualValidation Rule

    Background:
        * url baseURLCompliance

        * def ManualValidationRuleOffer = "TEST_OFFER" + uuid()  
        * def TestRulePartner = "TestManualValidationRulePartner"

        * def ruleSetManualValidation = "MANUAL_VALIDATION"
        * def ruleNameManualValidation = "PROFILE_DATA"

        * def offer = { offer_type : '#(ManualValidationRuleOffer)', product: 'maquininha'}

        Given path '/offers'
        And header Content-Type = 'application/json'
        And request offer
        When method POST

        * def ruleSetConfig = { manual_validation: {} }

        Given url mockURL
        And path '/temis-config/cadastral-validation-configs'
        And header Content-Type = 'application/json'
        And def catalogCompany = call CreateSingleLevelCatalog { offer_type: '#(ManualValidationRuleOffer)', role_type: '#(RoleType.Merchant)', person_type: '#(ProfileType.Company)', account_flag: false, rules_config: #(ruleSetConfig)}
        And request catalogCompany
        When method POST
        Then match response.validation_steps[*].rules_config contains deep ruleSetConfig

    Scenario: Should be analysing when profile company merchant is created 
        * def documentNumber =  CNPJGenerator()
        * def documentNormalized = DocumentNormalizer(documentNumber)
        * def profileID = uuid()
        * def profile = { profile_id: '#(profileID)', offer_type: '#(ManualValidationRuleOffer)', role_type: '#(RoleType.Merchant)', profile_type: '#(ProfileType.Company)', document_number: '#(documentNormalized)' }

        Given url mockURL
        And path '/v1/temis/profile/' + profileID
        And request profile
        When method POST
        Then assert responseStatus == 200

        * def event = { profile_id: '#(profileID)', entity_id: '#(profileID)', update_date: '2023-08-03T19:29:06.63556Z'  }
        * string json = event
        * def result = RegistrationEventsPublisher.publish(json)
        * eval sleep(3000)

        Given url baseURLCompliance
        And path '/state/' + profileID
        And header Content-Type = 'application/json'
        And retry until response.result == RuleResult.Analysing
        When method GET
        * def ruleResults = response.rule_set_result
        Then assert responseStatus == 200
        Then assert response.result == RuleResult.Analysing
        Then assert ruleResults[0].set == ruleSetManualValidation
        Then assert ruleResults[0].name == ruleNameManualValidation
        Then assert ruleResults[0].result == RuleResult.Analysing


    Scenario: Should override the profile state 
        * def documentNumber =  CNPJGenerator()
        * def documentNormalized = DocumentNormalizer(documentNumber)
        * def profileID = uuid()
        * def profile = { profile_id: '#(profileID)', offer_type: '#(ManualValidationRuleOffer)', role_type: '#(RoleType.Merchant)', profile_type: '#(ProfileType.Company)', document_number: '#(documentNormalized)' }

        Given url mockURL
        And path '/v1/temis/profile/' + profileID
        And request profile
        When method POST
        Then assert responseStatus == 200

        * def event = { profile_id: '#(profileID)', entity_id: '#(profileID)', update_date: '2023-08-03T19:29:06.63556Z'  }
        * string json = event
        * def result = RegistrationEventsPublisher.publish(json)
        * eval sleep(3000)

        Given url baseURLCompliance
        And path '/state/' + profileID
        And header Content-Type = 'application/json'
        And retry until response.result == RuleResult.Analysing
        When method GET
        * def ruleResults = response.rule_set_result
        Then assert responseStatus == 200
        Then assert response.result == RuleResult.Analysing
        Then assert ruleResults[0].set == ruleSetManualValidation
        Then assert ruleResults[0].name == ruleNameManualValidation
        Then assert ruleResults[0].result == RuleResult.Analysing
        Then assert ruleResults[0].pending == true

        * def override = 
        """
            {
                "entity_id": "#(profileID)",
                "entity_type": "#(Engine.Profile)",
                "rule_set": "#(ruleSetManualValidation)",
                "rule_name": "#(ruleNameManualValidation)",
                "result": "#(RuleResult.Approved)",
                "comments": "approved",
                "author": "me"
            }
        """

        Given url baseURLCompliance
        And path '/override'
        And header Content-Type = 'application/json'
        And request override
        When method POST
        Then assert responseStatus == 200

        Given url baseURLCompliance
        And path '/state/' + profileID
        And header Content-Type = 'application/json'
        And retry until response.result == RuleResult.Approved
        When method GET
        * def ruleResults = response.rule_set_result
        Then assert responseStatus == 200
        Then assert response.result == RuleResult.Approved
        Then assert ruleResults[0].set == ruleSetManualValidation
        Then assert ruleResults[0].name == ruleNameManualValidation
        Then assert ruleResults[0].result == RuleResult.Approved
        Then assert ruleResults[0].pending == false