
@ignore
Feature: PROMETHEUS

  Background:
    * url baseURLCompliance

    * def documentNumber = CPFGenerator()

  Scenario: On sucessfull POST should increment 201 metric

    Given url metricsApiURL
    And path "/metrics"
    When method GET
    Then assert responseStatus == OK
    And  def text = response
    And  def before_201_status_count = parseInt(karate.extract(text, 'http_response_status_code\\{code="201",handler="/compliance-int/offers",method="POST"\\} ([0-9]*)', 1)) || 0
    And  def before_201_duration_count = parseInt(karate.extract(text, 'http_request_duration_seconds_count\\{code="201",handler="/compliance-int/offers",method="POST"\\} ([0-9]*)', 1)) || 0

    * def createdOfferType = "TEST_OFFER" + uuid() 
    * def product = 'maquininha'
  
    Given url  baseURLCompliance
    And path '/offers'
    And header Content-Type = 'application/json'
    And def offer = {offer_type: '#(createdOfferType)', product: '#(product)'}
    And request offer
    When method POST
    Then assert responseStatus == CREATED
    And assert response.offer_type == createdOfferType

    Given url metricsApiURL
    And path "/metrics"
    When method GET
    Then assert responseStatus == OK
    And  def text = response
    And  def after_201_status_count = parseInt(karate.extract(text, 'http_response_status_code\\{code="201",handler="/compliance-int/offers",method="POST"\\} ([0-9]*)', 1)) || 0
    And  def after_201_duration_count = parseInt(karate.extract(text, 'http_request_duration_seconds_count\\{code="201",handler="/compliance-int/offers",method="POST"\\} ([0-9]*)', 1)) || 0
    And  match after_201_status_count == before_201_status_count + 1
    And  match after_201_duration_count == before_201_duration_count + 1

  Scenario: On unsucessfull POST should increment 400 metric

    Given url metricsApiURL
    And path "/metrics"
    When method GET
    Then assert responseStatus == OK
    And  def text = response
    And  def before_400_status_count = parseInt(karate.extract(text, 'http_response_status_code\\{code="400",handler="/compliance-int/offers",method="POST"\\} ([0-9]*)', 1)) || 0
    And  def before_400_duration_count = parseInt(karate.extract(text, 'http_request_duration_seconds_count\\{code="400",handler="/compliance-int/offers",method="POST"\\} ([0-9]*)', 1)) || 0

    * def createdOfferType = 'OFFER_' + CPFGenerator()
    * def product = 'maquininha'
  
    Given url  baseURLCompliance
    And path '/offers'
    And header Content-Type = 'application/json'
    And def offer = {offer_type: '', product: ''}
    And request offer
    When method POST
    Then assert responseStatus == BAD_CONTENT

    Given url metricsApiURL
    And path "/metrics"
    When method GET
    Then assert responseStatus == OK
    And  def text = response
    And  def after_400_status_count = parseInt(karate.extract(text, 'http_response_status_code\\{code="400",handler="/compliance-int/offers",method="POST"\\} ([0-9]*)', 1)) || 0
    And  def after_400_duration_count = parseInt(karate.extract(text, 'http_request_duration_seconds_count\\{code="400",handler="/compliance-int/offers",method="POST"\\} ([0-9]*)', 1)) || 0
    And  match after_400_status_count == before_400_status_count + 1
    And  match after_400_duration_count == before_400_duration_count + 1

  Scenario: On Success processing event should increment ACK and Success counters

    Given url metricsEventsURL
    And path "/metrics"
    When method GET
    Then assert responseStatus == OK
    And  def text = response
    And  def before_event_process_duration_seconds_count = parseInt(karate.extract(text, '\nevent_process_duration_seconds_count ([0-9]*)', 1)) || 0
    And  def before_event_waiting_duration_seconds_count = parseInt(karate.extract(text, '\nevent_waiting_duration_seconds_count ([0-9]*)', 1)) || 0
    And  def before_process_result_success_count = parseInt(karate.extract(text, 'events_process_status\\{handshake="ACK",status="success"\\} ([0-9]*)', 1)) || 0
    And  def before_process_result_error_count = parseInt(karate.extract(text, 'events_process_status\\{handshake="ACK",status="error"\\} ([0-9]*)', 1)) || 0
  
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
		* eval sleep(5000)

    Given url metricsEventsURL
    And path "/metrics"
    When method GET
    Then assert responseStatus == OK
    And  def text = response
    And  def after_event_process_duration_seconds_count = parseInt(karate.extract(text, '\nevent_process_duration_seconds_count ([0-9]*)', 1)) || 0
    And  def after_event_waiting_duration_seconds_count = parseInt(karate.extract(text, '\nevent_waiting_duration_seconds_count ([0-9]*)', 1)) || 0
    And  def after_process_result_success_count = parseInt(karate.extract(text, 'events_process_status\\{handshake="ACK",status="success"\\} ([0-9]*)', 1)) || 0
    And  def after_process_result_error_count = parseInt(karate.extract(text, 'events_process_status\\{handshake="ACK",status="error"\\} ([0-9]*)', 1)) || 0
    And  assert after_process_result_error_count   == before_process_result_error_count
    And  assert after_process_result_success_count > before_process_result_success_count 
    And  assert after_event_process_duration_seconds_count > before_event_process_duration_seconds_count
    And  assert after_event_waiting_duration_seconds_count > before_event_waiting_duration_seconds_count 

    
  Scenario: On Error processing event should increment NACK and Error counters

    Given url metricsEventsURL
    And path "/metrics"
    When method GET
    Then assert responseStatus == OK
    And  def text = response
    And  def before_event_process_duration_seconds_count = parseInt(karate.extract(text, '\nevent_process_duration_seconds_count ([0-9]*)', 1)) || 0
    And  def before_event_waiting_duration_seconds_count = parseInt(karate.extract(text, '\nevent_waiting_duration_seconds_count ([0-9]*)', 1)) || 0
    And  def before_process_result_success_count = parseInt(karate.extract(text, 'events_process_status\\{handshake="ACK",status="success"\\} ([0-9]*)', 1)) || 0
    And  def before_process_result_error_count = parseInt(karate.extract(text, 'events_process_status\\{handshake="NACK",status="error"\\} ([0-9]*)', 1)) || 0

    * def contractID = uuid()
		* def event = { contract_id: '#(contractID)', entity_id: '#(contractID)', entity_type: 'CONTRACT', event_type: 'CONTRACT_CREATED', parent_type: 'CONTRACT', update_date: '2023-08-03T19:29:06.63556Z'  }
		* string json = event
    * def result = RegistrationEventsPublisher.publish(json)
    * eval sleep(5000)

    Given url metricsEventsURL
    And path "/metrics"
    When method GET
    Then assert responseStatus == OK
    And  def text = response
    And  def after_event_process_duration_seconds_count = parseInt(karate.extract(text, '\nevent_process_duration_seconds_count ([0-9]*)', 1)) || 0
    And  def after_event_waiting_duration_seconds_count = parseInt(karate.extract(text, '\nevent_waiting_duration_seconds_count ([0-9]*)', 1)) || 0
    And  def after_process_result_success_count = parseInt(karate.extract(text, 'events_process_status\\{handshake="ACK",status="success"\\} ([0-9]*)', 1)) || 0
    And  def after_process_result_error_count = parseInt(karate.extract(text, 'events_process_status\\{handshake="NACK",status="error"\\} ([0-9]*)', 1)) || 0
    And  assert after_process_result_error_count   > before_process_result_error_count
    And  assert after_process_result_success_count >= before_process_result_success_count 
    And  assert after_event_process_duration_seconds_count > before_event_process_duration_seconds_count
    And  assert after_event_waiting_duration_seconds_count > before_event_waiting_duration_seconds_count 