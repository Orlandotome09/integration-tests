Feature: State Api

  Background:
    * url baseURLCompliance

    * def profileCreate = function(data){ return karate.repeat(data.repeat, function(i){ return karate.call("../data_driven/create_profile.feature", data)}) }

    * def StateApiOffer = "TEST_OFFER" + uuid() 
    * def StateApiOfferProblem = "TEST_OFFER" + uuid() 
    * def StateApiOfferPending = "TEST_OFFER" + uuid() 
    * def StateApiOfferPending2 = "TEST_OFFER" + uuid() 
    * def StateApiOfferPending3 = "TEST_OFFER" + uuid() 

    * def ruleSetManualBlock = "MANUAL_BLOCK"

    * def ruleNameManualBlock = "MANUAL_BLOCK"
    * def ruleNameCustomerNotFoundInSerasa = "CUSTOMER_NOT_FOUND_IN_SERASA"

    * def offer = { offer_type : '#(StateApiOffer)', product: 'maquininha'}
    * def offerProblem = { offer_type : '#(StateApiOfferProblem)', product: 'maquininha'}
    * def offerPending = { offer_type : '#(StateApiOfferPending)', product: 'maquininha'}
    * def offerPending2 = { offer_type : '#(StateApiOfferPending2)', product: 'maquininha'}
    * def offerPending3 = { offer_type : '#(StateApiOfferPending3)', product: 'maquininha'}

    Given path '/offers'
    And header Content-Type = 'application/json'
    And request offer
    When method POST

    Given path '/offers'
    And header Content-Type = 'application/json'
    And request offerProblem
    When method POST

    Given path '/offers'
    And header Content-Type = 'application/json'
    And request offerPending
    When method POST

    Given path '/offers'
    And header Content-Type = 'application/json'
    And request offerPending2
    When method POST

    Given path '/offers'
    And header Content-Type = 'application/json'
    And request offerPending3
    When method POST

    * def ruleSetConfig = { manual_block: {} }
    * def ruleSetConfigProblem = { manual_block: {}, incomplete:{ last_name_required : true, email_required: true } }
    * def ruleSetConfigPending = { pep: {}, watchlist: {  want_pep_tag: true, wanted_sources: [] } }

    Given url mockURL
    And path '/temis-config/cadastral-validation-configs'
    And header Content-Type = 'application/json'
    And def catalog = call CreateSingleLevelCatalog { offer_type: '#(StateApiOffer)',  role_type: '#(RoleType.Customer)', person_type: '#(ProfileType.Individual)', account_flag: false, rules_config: #(ruleSetConfig)}
    And request catalog
    When method POST

    Given url mockURL
    And path '/temis-config/cadastral-validation-configs'
    And header Content-Type = 'application/json'
    And def catalog = call CreateSingleLevelCatalog { offer_type: '#(StateApiOfferProblem)',  role_type: '#(RoleType.Customer)', person_type: '#(ProfileType.Individual)', account_flag: false, rules_config: #(ruleSetConfigProblem)}
    And request catalog
    When method POST

    Given url mockURL
    And path '/temis-config/cadastral-validation-configs'
    And header Content-Type = 'application/json'
    And def catalog = call CreateSingleLevelCatalog { offer_type: '#(StateApiOfferPending)',  role_type: '#(RoleType.Customer)', person_type: '#(ProfileType.Individual)', account_flag: false, rules_config: #(ruleSetConfigPending)}
    And request catalog
    When method POST

    Given url mockURL
    And path '/temis-config/cadastral-validation-configs'
    And header Content-Type = 'application/json'
    And def catalog = call CreateSingleLevelCatalog { offer_type: '#(StateApiOfferPending2)',  role_type: '#(RoleType.Customer)', person_type: '#(ProfileType.Individual)', account_flag: false, rules_config: #(ruleSetConfigPending)}
    And request catalog
    When method POST

    Given url mockURL
    And path '/temis-config/cadastral-validation-configs'
    And header Content-Type = 'application/json'
    And def catalog = call CreateSingleLevelCatalog { offer_type: '#(StateApiOfferPending3)',  role_type: '#(RoleType.Customer)', person_type: '#(ProfileType.Individual)', account_flag: false, rules_config: #(ruleSetConfigPending)}
    And request catalog
    When method POST

  Scenario: Get state for profile

    * def profileID = uuid()
    * def documentNumber = DocumentNormalizer(CPFGenerator())
    * def profile = {profile_id:'#(profileID)', document_number: '#(documentNumber)', offer_type: '#(StateApiOffer)',  role_type: '#(RoleType.Customer)', profile_type: '#(ProfileType.Individual)'}

    Given url mockURL
    And path '/v1/temis/profile/' + profileID
    And request profile
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
    And method GET
    Then assert response.result == RuleResult.Approved
    And assert response.rule_set_result[0].set == ruleSetManualBlock
    And assert response.rule_set_result[0].name == ruleNameManualBlock
    And assert response.rule_set_result[0].result == RuleResult.Approved
    And assert response.rule_set_result[0].pending == false
    And assert response.rule_set_result[0].metadata == null

  Scenario: Get compliance result for profile

    * def profileID = uuid()
    * def documentNumber = DocumentNormalizer(CPFGenerator())
    * def profile = {profile_id:'#(profileID)', document_number: '#(documentNumber)', offer_type: '#(StateApiOffer)',  role_type: '#(RoleType.Customer)', profile_type: '#(ProfileType.Individual)'}

    Given url mockURL
    And path '/v1/temis/profile/' + profileID
    And request profile
    When method POST
    Then assert responseStatus == 200

    * def event = { profile_id: '#(profileID)', entity_id: '#(profileID)', entity_type: 'PROFILE', event_type: 'PROFILE_CHANGED', parent_type: 'PROFILE', update_date: '2023-08-03T19:29:06.63556Z'  }
    * string json = event
    * def result = RegistrationEventsPublisher.publish(json)
    * eval sleep(1000)

    Given url baseURLCompliance
    And path '/result/' + profileID
    And header Content-Type = 'application/json'
    And retry until response.result == RuleResult.Approved
    And method GET
    Then assert response.result == RuleResult.Approved
    And assert response.entity_id == profileID
    And assert response.entity_type == EntityType.Profile

  @CreateComplianceProfileState
  Scenario: Execute compliance check for profile

    * def profileID = uuid()
    * def documentNumber = DocumentNormalizer(CPFGenerator())
    * def profile = {profile_id:'#(profileID)', document_number: '#(documentNumber)', offer_type: '#(StateApiOffer)',  role_type: '#(RoleType.Customer)', profile_type: '#(ProfileType.Individual)'}

    Given url mockURL
    And path '/v1/temis/profile/' + profileID
    And request profile
    When method POST
    Then assert responseStatus == 200

    Given url baseURLCompliance
    And path '/check/' + profileID
    And header Content-Type = 'application/json'
    And request { entity_type: '#(Engine.Profile)'}
    And method POST
    Then assert response.result == RuleResult.Approved
    And assert response.entity_id == profileID
    And assert response.entity_type == EntityType.Profile

  Scenario: Confirm problems codes in check return

    * def profileID = uuid()
    * def documentNumber = DocumentNormalizer(CPFGenerator())
    * def profile = {profile_id:'#(profileID)', document_number: '#(documentNumber)', offer_type: '#(StateApiOfferProblem)',  role_type: '#(RoleType.Customer)', profile_type: '#(ProfileType.Individual)'}

    Given url mockURL
    And path '/v1/temis/profile/' + profileID
    And request profile
    When method POST
    Then assert responseStatus == 200

    Given url baseURLCompliance
    And path '/check/' + profileID
    And header Content-Type = 'application/json'
    And request { entity_type: '#(Engine.Profile)'}
    And method POST
    Then assert response.result == RuleResult.Incomplete
    And assert response.entity_id == profileID
    And assert response.entity_type == EntityType.Profile
    And assert response.detailed_result[0].details[0].code == "EMAIL_REQUIRED"
    And assert response.detailed_result[0].details[1].code == "LAST_NAME_REQUIRED"

  Scenario: should not compliance check when missing engine name in request body

    * def profileID = uuid()
    * def documentNumber = DocumentNormalizer(CPFGenerator())
    * def profile = {profile_id:'#(profileID)', document_number: '#(documentNumber)', offer_type: '#(StateApiOfferProblem)',  role_type: '#(RoleType.Customer)', profile_type: '#(ProfileType.Individual)'}
    * def expectedError = " is an invalid Engine name"

    Given url mockURL
    And path '/v1/temis/profile/' + profileID
    And request profile
    When method POST
    Then assert responseStatus == 200

    Given url baseURLCompliance
    And path '/check/' + profileID
    And header Content-Type = 'application/json'
    And request {}
    And method POST
    Then assert responseStatus == 400
    Then assert response.error == expectedError

  Scenario: should not compliance check when engine name is invalid

    * def profileID = uuid()
    * def documentNumber = DocumentNormalizer(CPFGenerator())
    * def profile = {profile_id:'#(profileID)', document_number: '#(documentNumber)', offer_type: '#(StateApiOfferProblem)',  role_type: '#(RoleType.Customer)', profile_type: '#(ProfileType.Individual)'}
    * def invalidEngineName = "invalid"
    * def expectedError = invalidEngineName +" is an invalid Engine name"

    Given url mockURL
    And path '/v1/temis/profile/' + profileID
    And request profile
    When method POST
    Then assert responseStatus == 200

    Given url baseURLCompliance
    And path '/check/' + profileID
    And header Content-Type = 'application/json'
    And request {entity_type: '#(invalidEngineName)'}
    And method POST
    Then assert responseStatus == 400
    Then assert response.error == expectedError
