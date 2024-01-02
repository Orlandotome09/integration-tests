@ignore
Feature: Create Profile

  Background:        

    * def calcLR =
    """
        function(params){
            LRBase = {            
                profile_id: "",
                full_name: "John Doe",                
                email: "aaa@gmail.com",
                phone: "23455",
                nationality: "BRA",
                birth_date: "01/01/2000"
            }
            if (params.data != null) {                
                return params.data
            }
            
            LRBase.profile_id = params.profile_id

            return LRBase            
        }
    """

  Scenario: Generate new Legal Representative    
    * def LRParams = params
    * def LRData = calcLR(LRParams)
    * def legalRepresentativeID = uuid()
    * def cnpjDocument = CNPJGenerator()
	* def cpfDocument = CPFGenerator()		
    * def documentNumber = LRBase.document_number == null ? ( LRBase.profile_type == 'COMPANY' ? cnpjDocument : cpfDocument) : LRBase.document_number
	* LRData.legal_representative_id = LRBase.legal_representative_id == null ? legalRepresentativeID : LRBase.legal_representative_id
	* LRData.document_number = documentNumber

    Given url mockURL
	And path '/v1/temis/legal-representative/' + LRData.legal_representative_id
	And request LRData
	When method POST
	Then assert responseStatus == 200

    * def funcNow = call funcNow    
    * def event = { profile_id: '#(LRData.legal_representative_id.profile_id)', entity_id: '#(LRData.legal_representative_id)', entity_type: 'PROFILE', event_type: 'LEGAL_REPRESENTATIVE_CREATED', parent_type: 'PROFILE', update_date: '#(funcNow)'  }	
    
	* string json = event
	* def result = RegistrationEventsPublisher.publish(json)