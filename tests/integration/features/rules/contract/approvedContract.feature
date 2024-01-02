Feature: Contract Rule Incomplete

    Background:
        * url baseURLCompliance

	    * def ContractRuleOffer = "TEST_OFFER" + uuid()  

		* def ruleSetIncompleteContract = "INCOMPLETE_CONTRACT"

		* def ruleNameDocumentNotFound = "DOCUMENT_NOT_FOUND"
		* def ruleNameFileNotFound = "FILE_NOT_FOUND"

		* def ProblemCodeInvoiceIsRequired = "INVOICE_IS_REQUIRED"
		* def ProblemCodeInvoiceDocumentNotFound = "INVOICE_DOCUMENT_NOT_FOUND"
		* def ProblemCodeInvoiceAssociatedToAnotherProfile = "INVOICE_ASSOCIATED_TO_ANOTHER_PROFILE"
		* def ProblemCodeInvoiceFileNotFound = "INVOICE_FILE_NOT_FOUND"
		* def ProblemCodeInvoiceNotFound = "STATE_NOT_FOUND"		
		* def ProblemCodeProfileNotApproved = "PROFILE_NOT_APPROVED"		
	
    Scenario: Should not find invoice      
		
		* def contractID = uuid()
		* def profileID = uuid()
		* def contract = { contract_id: '#(contractID)', profile_id: '#(profileID)', document_id : ''}

		Given url mockURL
		And path '/v1/temis/contract/' + contractID
		And request contract
		When method POST
		Then assert responseStatus == 200	

		* def event = { contract_id: '#(contractID)', entity_id: '#(contractID)', entity_type: 'CONTRACT', event_type: 'CONTRACT_CREATED', parent_type: 'CONTRACT', update_date: '2023-08-03T19:29:06.63556Z'  }
		* string json = event
    	
		* def result = RegistrationEventsPublisher.publish(json)

		* def expectedResponse =
		"""
			{
				result: '#(RuleResult.Approved)'	
			}
		"""

		Given url baseURLCompliance
		And path '/state/' + contractID
        And header Content-Type = 'application/json'
        And params ({ only_pending: false })		
		And retry until response.result == RuleResult.Approved
        When method GET
		Then assert responseStatus == 200
		And match response contains deep expectedResponse

	Scenario: Should not find document

		* def contractID = uuid()
		* def documentID = uuid()
		* def profileID = uuid()
		* def contract = { contract_id: '#(contractID)', profile_id: '#(profileID)', document_id: '#(documentID)'}

		Given url mockURL
		And path '/v1/temis/contract/' + contractID
		And request contract
		When method POST

		* def documentResponse = { body: { document_id: '#(documentID)', entity_id: '#(profileID)'}, status: 404}

		Given url mockURL
		And path '/v1/temis/document/' + documentID
		And request documentResponse
		When method POST

		* def event = { contract_id: '#(contractID)', entity_id: '#(contractID)', entity_type: 'CONTRACT', event_type: 'CONTRACT_CREATED', parent_type: 'CONTRACT', update_date: '2023-08-03T19:29:06.63556Z'  }
		* string json = event
    	
		* def result = RegistrationEventsPublisher.publish(json)

		* def metadata = 'Invoice Document: ' + documentID + ' not found'
		* def expectedResponse =
		"""
			{
				result: '#(RuleResult.Approved)'
			}
		"""

		Given url baseURLCompliance
		And path '/state/' + contractID
        And header Content-Type = 'application/json'
        And params ({ only_pending: false })		
		And retry until response.result == RuleResult.Approved
        When method GET
		Then assert responseStatus == 200
		And match response contains deep expectedResponse

	Scenario: Should find document associated with another profile

		* def contractID = uuid()
		* def documentID = uuid()
		* def profileID = uuid()
		* def contract = { contract_id: '#(contractID)', profile_id: '#(profileID)', document_id: '#(documentID)'}

		Given url mockURL
		And path '/v1/temis/contract/' + contractID
		And request contract
		When method POST

		* def anotherProfileID = uuid()
		* def documentResponse = { body: { document_id: '#(documentID)', entity_id: '#(anotherProfileID)'}, status: 200}

		Given url mockURL
		And path '/v1/temis/document/' + documentID
		And request documentResponse
		When method POST

		* def event = { contract_id: '#(contractID)', entity_id: '#(contractID)', entity_type: 'CONTRACT', event_type: 'CONTRACT_CREATED', parent_type: 'CONTRACT', update_date: '2023-08-03T19:29:06.63556Z'  }
		* string json = event
    	
		* def result = RegistrationEventsPublisher.publish(json)

		* def expectedResponse =
		"""
			{
				result: '#(RuleResult.Approved)'
			}
		"""

		Given url baseURLCompliance
		And path '/state/' + contractID
        And header Content-Type = 'application/json'
        And params ({ only_pending: false })		
		And retry until response.result == RuleResult.Approved
        When method GET
		Then assert responseStatus == 200
		And match response contains deep expectedResponse

	Scenario: Should not find file for invoice

		* def contractID = uuid()
		* def documentID = uuid()
		* def profileID = uuid()
		* def contract = { contract_id: '#(contractID)', profile_id: '#(profileID)', document_id: '#(documentID)'}

		Given url mockURL
		And path '/v1/temis/contract/' + contractID
		And request contract
		When method POST

		* def documentResponse = { body: { document_id: '#(documentID)', entity_id: '#(profileID)'}, status: 200}

		Given url mockURL
		And path '/v1/temis/document/' + documentID
		And request documentResponse
		When method POST

		* def documentFilesResponse = []

		Given url mockURL
		And path '/v1/temis/document/' + documentID + '/files'
		And request documentFilesResponse
		When method POST

		* def event = { contract_id: '#(contractID)', entity_id: '#(contractID)', entity_type: 'CONTRACT', event_type: 'CONTRACT_CREATED', parent_type: 'CONTRACT', update_date: '2023-08-03T19:29:06.63556Z'  }
		* string json = event
    	
		* def result = RegistrationEventsPublisher.publish(json)

		* def expectedResponse =
		"""
			{
			 result: '#(RuleResult.Approved)'
			}
		"""

		Given url baseURLCompliance
		And path '/state/' + contractID
        And header Content-Type = 'application/json'
        And params ({ only_pending: false })		
		And retry until response.result == RuleResult.Approved
        When method GET
		Then assert responseStatus == 200
		And match response contains deep expectedResponse

	Scenario: Should find document and invoice file

		* def contractID = uuid()
		* def documentID = uuid()
		* def profileID = uuid()
		* def contract = { contract_id: '#(contractID)', profile_id: '#(profileID)', document_id: '#(documentID)'}

		Given url mockURL
		And path '/v1/temis/contract/' + contractID
		And request contract
		When method POST

		* def documentResponse = { body: { document_id: '#(documentID)', entity_id: '#(profileID)'}, status: 200}

		Given url mockURL
		And path '/v1/temis/document/' + documentID
		And request documentResponse
		When method POST

		* def documentFileID = uuid()
		* def fileID = uuid()
		* def documentFilesResponse = [{document_file_id: '#(documentFileID)', document_id:'#(documentID)', file_id:'#(fileID)'}]

		Given url mockURL
		And path '/v1/temis/document/' + documentID + '/files'
		And request documentFilesResponse
		When method POST

		* def event = { contract_id: '#(contractID)', entity_id: '#(contractID)', entity_type: 'CONTRACT', event_type: 'CONTRACT_CREATED', parent_type: 'CONTRACT', update_date: '2023-08-03T19:29:06.63556Z'  }
		* string json = event
    	* def result = RegistrationEventsPublisher.publish(json)
		* eval sleep(1000)

		* def expectedResponse =
		"""
			{
				result: '#(RuleResult.Approved)'			
			}
		"""

		Given url baseURLCompliance
		And path '/state/' + contractID
        And header Content-Type = 'application/json'
        And params ({ only_pending: false })		
		And retry until response.result == RuleResult.Approved
        When method GET
		Then assert responseStatus == 200
		And match response contains deep expectedResponse