Feature: Reprocess profile
    Background:
        * url baseURLCompliance

        * def engineNameProfile = "PROFILE"
        * def eventTypeProfileReprocess = "PROFILE_REPROCESS"

        * def StateApiOffer = "TEST_OFFER" + uuid()
        * def offer = { offer_type : '#(StateApiOffer)', product: 'maquininha'}

        * def ruleSetManualBlock = "MANUAL_BLOCK"
        * def ruleNameManualBlock = "MANUAL_BLOCK"

        * def ruleSetConfig = { manual_block: {} }
    
        Given url mockURL
        And path '/temis-config/cadastral-validation-configs'
        And header Content-Type = 'application/json'
        And def catalog = call CreateSingleLevelCatalog { offer_type: '#(StateApiOffer)',  role_type: '#(RoleType.Customer)', person_type: '#(ProfileType.Individual)', account_flag: false, rules_config: #(ruleSetConfig)}
        And request catalog
        When method POST
    
    Scenario: should reprocess profile
        * def profileID = uuid()
        * def documentNumber = DocumentNormalizer(CPFGenerator())
        * def profile = {profile_id:'#(profileID)', document_number: '#(documentNumber)', offer_type: '#(StateApiOffer)',  role_type: '#(RoleType.Customer)', profile_type: '#(ProfileType.Individual)'}

        Given url mockURL
        And path '/v1/temis/profile/' + profileID
        And request profile
        When method POST
        Then assert responseStatus == 200

        Given url baseURLCompliance
        And path '/state/' + profileID
        And header Content-Type = 'application/json'
        And method GET
        Then assert responseStatus == 404

        * def reprocessRequest = 
        """
        {
            "engine_name": #(engineNameProfile),
            "ids":["#(profileID)"]
        }
        """

        Given url baseURLCompliance
        And path '/states/reprocess'
        And header Content-Type = 'application/json'
        And request reprocessRequest
        When method POST
        Then assert responseStatus == 200
        Then assert response.reprocessed.length === 1

        Given url baseURLCompliance
        And path '/state/' + profileID
        And header Content-Type = 'application/json'
        And params ({ only_pending: false })
        And retry until response.result == RuleResult.Approved
        And method GET
        Then assert response.result == RuleResult.Approved
        And assert response.rule_set_result[0].set == ruleSetManualBlock
        And assert response.rule_set_result[0].name == ruleNameManualBlock
        And assert response.rule_set_result[0].result == RuleResult.Approved
        And assert response.rule_set_result[0].pending == false
        And assert response.rule_set_result[0].metadata == null

    Scenario: should not reprocess profile with invalid profile id
        * def reprocessRequest = 
        """
        {
            "engine_name": #(engineNameProfile),
            "ids":["non valid profile id"]
        }
        """

        Given url baseURLCompliance
        And path '/states/reprocess'
        And header Content-Type = 'application/json'
        And request reprocessRequest
        When method POST
        Then assert responseStatus == 400