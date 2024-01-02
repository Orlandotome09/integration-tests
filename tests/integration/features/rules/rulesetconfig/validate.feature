Feature: Approve registration of offers that do not have validation steps

    Background:
        
        * url baseURLCompliance

		* def documentNumber = CPFGenerator()
		* def companyDocumentNumber = CNPJGenerator()
		* def ruleOffer = "TEST_OFFER" + uuid()  
        
        Given path '/offers'
		And header Content-Type = 'application/json'
		And def offer = {offer_type: '#(ruleOffer)', product : 'maquininha'}
		And request offer
		When method POST

	Scenario: Should approve individual entity since catalog rules is empty

		* def profileID = uuid()
		* def documentNumber = DocumentNormalizer(CPFGenerator())
		* def profile = { profile_id: '#(profileID)', offer_type: '#(ruleOffer)', role_type: '#(RoleType.Customer)', profile_type: '#(ProfileType.Individual)', document_number: '#(documentNumber)' }

		Given url mockURL
		And path '/v1/temis/profile/' + profileID
		And request profile
		When method POST
		Then assert responseStatus == 200

		* def params =
		"""
			{
				offer_type: '#(ruleOffer)',
				role_type: '#(RoleType.Customer)',
				person_type: '#(ProfileType.Individual)',
				steps: []
			}
		"""     
		* def catalog = call CreateMultiLevelCatalog params

		Given url mockURL
        And path '/temis-config/cadastral-validation-configs'
		And header Content-Type = 'application/json'
		And request catalog
		When method POST

		* def validator = { entity_type: '#(EntityType.Profile)' }

		Given url baseURLCompliance
		Given path '/check/' + profileID
		And header Content-Type = 'application/json'
		And request validator
		And method POST
		Then assert response.result == RuleResult.Approved
		And assert response.entity_id == profileID
		And assert response.entity_type == EntityType.Profile

		Given url baseURLCompliance
		Given path '/state/' + profileID
		And header Content-Type = 'application/json'
		When method GET
		Then assert responseStatus == 200
		And match response.rule_set_result == '#notpresent'


	Scenario: Should approve legal entity since catalog rules is empty

		* def companyProfileID = uuid()
		* def companyDocumentNumber = DocumentNormalizer(CNPJGenerator())
		* def companyProfile = { profile_id:'#(companyProfileID)', offer_type: '#(ruleOffer)', role_type: '#(RoleType.Customer)', profile_type: 'COMPANY', document_number: '#(companyDocumentNumber)' }

		Given url mockURL
		And path '/v1/temis/profile/' + companyProfileID
		And request companyProfile
		When method POST
		Then assert responseStatus == 200

		* def params =
		"""
			{
				offer_type: '#(ruleOffer)',
				
				role_type: '#(RoleType.Customer)',
				person_type: 'COMPANY',
				steps: []
			}
		"""     
		* def catalog = call CreateMultiLevelCatalog params

		Given url mockURL
        And path '/temis-config/cadastral-validation-configs'
		And header Content-Type = 'application/json'
		And request catalog
		When method POST

		* def validator = { entity_type: '#(EntityType.Profile)' }

		Given url baseURLCompliance
		Given path '/check/' + companyProfileID
		And header Content-Type = 'application/json'
		And request validator
		And method POST
		Then assert response.result == RuleResult.Approved
		And assert response.entity_id == companyProfileID
		And assert response.entity_type == EntityType.Profile

		Given url baseURLCompliance
		Given path '/state/' + companyProfileID
		And header Content-Type = 'application/json'
		When method GET
		Then assert responseStatus == 200
		And match response.rule_set_result == '#notpresent'