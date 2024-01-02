@ignore
Feature: Create Profile Individual

Background:  

	* def documentNumber = CPFGenerator()

	* def calcProfile =
    """
        function(params, profile_id, cnpjDocument, cpfDocument, profileReturn){			
			profileReturn =	{
				profile_id: "",
				partner_id: "",
				offer_type: "",
				role_type: "",
				profile_type: "",
				document_number: "",
				callback_url: "",			
				email: "teste@teste.com",        
				individual: {
					first_name: "Test",
					last_name: "Temis",
					date_of_birth: "1991-07-30T00:00:00Z",
					phones: [
						{
							type: "comercial",
							number: "123",
							country_code: "55",
							area_code: "11"
						}
					]
				},
				company: {
					legal_name: "Empresa Legal",
					business_name: "Empresa Legal S.A.",
					tax_payer_identification: "1234",
					company_registration_number: "56789",
					date_of_incorporation: "1989-02-13",
					place_of_incorporation: "BR",
					share_capital: {
						amount: 10009,
						currency: "USD"
					},
					license: "Open Source",
					website: "www.bexs.com",
					goods_delivery: {
						average_delivery_days: 5,
						shipping_methods: "BOAT",
						tracking_code_available: true
					}
				}			
			}

            if (params.data != null) {
                profileReturn = params.data
            } else {				
				profileReturn.partner_id = params.partner_id
				profileReturn.offer_type = params.offer_type
				profileReturn.role_type = params.role_type
				profileReturn.callback_url = params.callback_url
				profileReturn.partner_id = params.partner_id
				profileReturn.profile_type = params.profile_type
			}		
												
			
			return profileReturn			            
        }
    """

Scenario: Generate IDS and document, mock data and dispatch compliance created event
  	* def profileParams = params
	* def profileID = uuid()
	* def cnpjDocument = CNPJGenerator()
	* def cpfDocument = CPFGenerator()		
	* def profileData = calcProfile(profileParams, profileID, cnpjDocument, cpfDocument)	
	* def documentNumber = profileData.document_number == null ? ( profileData.profile_type == 'COMPANY' ? cnpjDocument : cpfDocument) : profileData.document_number
	* profileData.profile_id = profileData.profile_id == null ? profileID : profileData.profile_id
	* profileData.document_number = documentNumber
		
	Given url mockURL
	And path '/v1/temis/profile/' + profileID
	And request profileData
	When method POST
	Then assert responseStatus == 200

	* def event = { profile_id: '#(profileID)', entity_id: '#(profileID)', entity_type: 'PROFILE', event_type: 'PROFILE_CREATED', parent_type: 'PROFILE', update_date: '#(timeNow)'  }
	* string json = event
	* def result = RegistrationEventsPublisher.publish(json)	