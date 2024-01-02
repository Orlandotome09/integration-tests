@ignore
Feature: stateful mock server
  for help, see: https://github.com/intuit/karate/wiki/ZIP-Release

  Background:
    * configure cors = true
    * def uuid = function(){ return java.util.UUID.randomUUID() + '' }
    * def pause = function(milliseconds){ java.lang.Thread.sleep(milliseconds) }
    * eval pause(5000)

    * def responseWatchListDocuments = {}
    * def responseWatchList = {title:"SomeToken", name: "Someone", link: "/link", watch: "", other: "", entries: ["someEntry"], sources: []}
    * def responseBureau = { cadastral_situation: { value: {situation: "", situation_code: 0} }}
    * def webPaymentCallbacksReceived = {}
    * def responseCallback = {}

    * def addressesResponses = {}
    * def contactsResponses = {}
    * def notificationRecipientsResponses = {}
    * def accountsResponses = {}
    * def contractResponses = {}
    * def documentResponses = {}
    * def documentsResponses = {}
    * def documentFileResponses = {}
    * def fileResponses = {}
    * def legalRepresentativeResponses = {}
    * def legalRepresentativesResponses = {}
    * def profileResponses = {}
    * def ownershipStructureResponses = {}
    * def partnerResponses = {}
    * def boardOfDirectorsResponses = {}

    * def enrichmentIndividualResponses =  {}
    * def enrichmentLegalEntityResponses =  {}
    * def enrichmentResponses =  {}

    * def internalListResponses = {}

    * def internalListResponsesPEPCOAF = {}

    * def doaResponses = {}

    * def featuresResponses = {}

    * def watchlistCompanyResponse = {}

    * def enrichedPersonResponses = {}

    * def MessageSubscriber = Java.type('bexs.MessageSubscriber')
    * def host = java.lang.System.getenv('PUBSUB_EMULATOR_HOST') != undefined ? java.lang.System.getenv('PUBSUB_EMULATOR_HOST') : 'localhost:8681';
    * def projectId = java.lang.System.getenv('PUBSUB_PROJECT_ID') != undefined ? java.lang.System.getenv('PUBSUB_PROJECT_ID') : 'local-project';
    * def SubscribeNotificationEvents = new MessageSubscriber({'host': host, 'project-id': projectId, 'subscription-id': 'notification_subscription'})
    * def listenNotification = new java.lang.Thread(SubscribeNotificationEvents.listen()).start();

    * def bexDigitalProjectId = java.lang.System.getenv('BEXS_DIGITAL_PROJECT_ID') != undefined ? java.lang.System.getenv('BEXS_DIGITAL_PROJECT_ID') : 'local-project';
    * def SubscribeTreeAdapterEvents = new MessageSubscriber({'host': host, 'project-id': bexDigitalProjectId, 'subscription-id': 'tree_adapter_subscription'})
    * def listenTreeAdapter = new java.lang.Thread(SubscribeTreeAdapterEvents.listen()).start();

    * def SubscribeLimitEvents = new MessageSubscriber({'host': host, 'project-id': bexDigitalProjectId, 'subscription-id': 'limit_subscription'})
    * def listenLimit = new java.lang.Thread(SubscribeLimitEvents.listen()).start();

    * def SubscribeStateEvents = new MessageSubscriber({'host': host, 'project-id': bexDigitalProjectId, 'subscription-id': 'temis-compliance-state-events-v1-test-subscription'})
    * def listenStateEvents = new java.lang.Thread(SubscribeStateEvents.listen()).start();

    * def postgresHost = java.lang.System.getenv('DATABASE_HOST') != undefined ? java.lang.System.getenv('DATABASE_HOST') : 'localhost'
    * def postgresPort = java.lang.System.getenv('DATABASE_PORT') != undefined ? java.lang.System.getenv('DATABASE_PORT') : '5432'
    * def postgresDatabaseName =  java.lang.System.getenv('DATABASE_NAME') != undefined ? java.lang.System.getenv('DATABASE_NAME') : 'temis_compliance'
    * def postgresUserName =  java.lang.System.getenv('DATABASE_USERNAME') != undefined ? java.lang.System.getenv('DATABASE_USERNAME') : 'compliance'
    * def postgresPassword =  java.lang.System.getenv('DATABASE_PASSWORD') != undefined ? java.lang.System.getenv('DATABASE_PASSWORD') : 'compliance'

    * def PostgresReaderWriter = Java.type('bexs.PostgresReaderWriter')
    * def postgresReaderWriter = new PostgresReaderWriter(postgresHost, postgresPort, postgresDatabaseName, postgresUserName, postgresPassword)

    * def cadastralValidationConfigs = []

  #Subscribe-------------------------

  #Notification Command
  Scenario: pathMatches('/subscribe/notification/{entityId}') && methodIs('get')
    * string messageList = SubscribeNotificationEvents.listMessages()
    * print messageList
    * json jsonList = messageList
    * def fun = function (e) { return JSON.parse(e) }
    * def jsonList = karate.map(jsonList, fun)
    * def filter = function(element){ return element.EntityID == pathParams.entityId }
    * def response = karate.filter(jsonList, filter)
    * def responseStatus = 200

  #Tree Adapter Command
  Scenario: pathMatches('/subscribe/tree-adapter/{profile_id}') && methodIs('get')
    * string messageList = SubscribeTreeAdapterEvents.listMessages()
    * print messageList
    * json jsonList = messageList
    * def fun = function (e) { return JSON.parse(e) }
    * def jsonList = karate.map(jsonList, fun)
    * def filter = function(element){ return element.profile_id == pathParams.profile_id }
    * def response = karate.filter(jsonList, filter)
    * def responseStatus = 200

  #Limit Command
  Scenario: pathMatches('/subscribe/limit/{profile_id}') && methodIs('get')
    * string messageList = SubscribeLimitEvents.listMessages()
    * print messageList
    * json jsonList = messageList
    * def fun = function (e) { return JSON.parse(e) }
    * def jsonList = karate.map(jsonList, fun)
    * def filter = function(element){ return element.profile_id == pathParams.profile_id }
    * def response = karate.filter(jsonList, filter)
    * def responseStatus = 200

  #State Events
  Scenario: pathMatches('/subscribe/state-events/{entity_id}') && methodIs('get')
    * string messageList = SubscribeStateEvents.listMessages()
    * print messageList
    * json jsonList = messageList
    * def fun = function (e) { return JSON.parse(e) }
    * def jsonList = karate.map(jsonList, fun)
    * def filter = function(element){ return element.entity_id == pathParams.entity_id || pathParams.entity_id == 'ALL' }
    * def res = karate.filter(jsonList, filter)
    * json response = res
    * def responseStatus = 200

  # Postgres Read Write
  Scenario: pathMatches('/postgres/read/state/{entity_id}') && methodIs('get')
    * string entityId = pathParams.entity_id
    * def queryResponse = postgresReaderWriter.queryStateById(entityId)
    * def responseStatus = queryResponse != "nothing to return" ? 200 : 204
    * def response = queryResponse

  Scenario: pathMatches('/postgres/insert/state/expired') && methodIs('post')
    * def entityID = request.entity_id
    * def result = request.result
    * def insertResponse = postgresReaderWriter.insertExpiredState(entityID, result)
    * def responseStatus = 200
    * def response = insertResponse

  Scenario: pathMatches('/postgres/insert/state/not_expired') && methodIs('post')
    * def entityID = request.entity_id
    * def result = request.result
    * def insertResponse = postgresReaderWriter.insertNotExpiredState(entityID, result)
    * def responseStatus = 200
    * def response = insertResponse

  Scenario: pathMatches('/postgres/clear') && methodIs('post')
    * def tableName = request.table_name
    * def rowsAffected = postgresReaderWriter.deleteRecords(tableName)
    * def responseStatus = 200
    * def response = { rows_affected: "#(rowsAffected)"}

  # Postgres Delete
  Scenario: pathMatches('/postgres/delete/{table_name}/{key_name}/{key_value}') && methodIs('post')
    * def queryResponse = postgresReaderWriter.deleteRecord(pathParams.table_name,pathParams.key_name,pathParams.key_value)
    * def responseStatus = queryResponse != "[]" ? 200 : 204
    * def response = queryResponse

  #Registration-------------------------
  #Address: Search
  Scenario: pathMatches('/v1/temis/addresses') && methodIs('post') && paramValue('profile_id') != null
    * def profile_id = paramValue('profile_id')
    * addressesResponses[profile_id] = request
    * def response = addressesResponses
    * def responseStatus = 200

  Scenario: pathMatches('/v1/temis/addresses') && methodIs('get') && paramValue('profile_id') != null
    * def profile_id = paramValue('profile_id')
    * def response = addressesResponses[profile_id] != null ?  addressesResponses[profile_id] : []
    * def responseStatus = 200

  ##Account: Search
  Scenario: pathMatches('/v1/temis/profile/{profile_id}/accounts') && methodIs('post')
    * accountsResponses[pathParams.profile_id] = request
    * def responseStatus = 200
    * def response = accountsResponses

  Scenario: pathMatches('/v1/temis/profile/{profile_id}/accounts') && methodIs('get')
    * def response = accountsResponses[pathParams.profile_id] != null ? accountsResponses[pathParams.profile_id] : {}
    * def responseStatus = accountsResponses[pathParams.profile_id] != null ? 200 : 500

  #Contacts: Search
  Scenario: pathMatches('/v1/temis/contacts') && methodIs('post') && paramValue('profile_id') != null
    * def profile_id = paramValue('profile_id')
    * contactsResponses[profile_id] = request
    * def response = contactsResponses
    * def responseStatus = 200

  Scenario: pathMatches('/v1/temis/contacts') && methodIs('get') && paramValue('profile_id') != null
    * def profile_id = paramValue('profile_id')
    * def response = contactsResponses[profile_id] != null ?  contactsResponses[profile_id] : []
    * def responseStatus = 200

  #NotificationRecipient: Search
  Scenario: pathMatches('/v1/temis/notification-recipients') && methodIs('post') && paramValue('profile_id') != null
    * def profile_id = paramValue('profile_id')
    * notificationRecipientsResponses[profile_id] = request
    * def response = notificationRecipientsResponses
    * def responseStatus = 200

  Scenario: pathMatches('/v1/temis/notification-recipients') && methodIs('get') && paramValue('profile_id') != null
    * def profile_id = paramValue('profile_id')
    * def response = notificationRecipientsResponses[profile_id] != null ?  notificationRecipientsResponses[profile_id] : []
    * def responseStatus = 200

  #Contract: Get
  Scenario: pathMatches('/v1/temis/contract/{id}') && methodIs('post')
    * contractResponses[pathParams.id] = request
    * def responseStatus = 200
    * def response = contractResponses

  Scenario: pathMatches('/v1/temis/contract/{id}') && methodIs('get')
    * def contract_id = pathParams.id
    * def responseStatus = 200
    * def response = contractResponses[contract_id] != null ? contractResponses[contract_id] : {}

  #Document: Get
  Scenario: pathMatches('/v1/temis/document/{id}') && methodIs('post')
    * documentResponses[pathParams.id] = request
    * def responseStatus = 200
    * def response = documentResponses

  Scenario: pathMatches('/v1/temis/document/{id}') && methodIs('get')
    * def document_id = pathParams.id
    * def responseStatus = documentResponses[document_id].status != null ? documentResponses[document_id].status : 200
    * def response = documentResponses[document_id].body != null ? documentResponses[document_id].body : {}

  #Document: Search
  Scenario: pathMatches('/v1/temis/documents') && methodIs('post') && paramValue('entity_id') != null
    * def entity_id = paramValue('entity_id')
    * documentsResponses[entity_id] = request
    * print(documentsResponses)
    * def response = documentsResponses
    * def responseStatus = 200

  Scenario: pathMatches('/v1/temis/documents') && methodIs('get') && paramValue('entity_id') != null
    * def entity_id = paramValue('entity_id')
    * def response = documentsResponses[entity_id] != null ?  documentsResponses[entity_id] : []
    * def responseStatus = 200

  #DocumentFile: Search
  Scenario: pathMatches('/v1/temis/document/{id}/files') && methodIs('post')
    * documentFileResponses[pathParams.id] = request
    * def responseStatus = 200
    * def response = documentFileResponses

  Scenario: pathMatches('/v1/temis/document/{id}/files') && methodIs('get')
    * def responseStatus = 200
    * def response = documentFileResponses[pathParams.id] != null ? documentFileResponses[pathParams.id] : []

  #File: Get url
  Scenario: pathMatches('/v1/temis/file/{id}/url') && methodIs('post')
    * fileResponses[pathParams.id] = request
    * def responseStatus = 200
    * def response = fileResponses

  Scenario: pathMatches('/v1/temis/file/{id}/url') && methodIs('get')
    * def responseStatus = 200
    * def response = fileResponses[pathParams.id] != null ? fileResponses[pathParams.id] : {}

  #Legal Representative: Search
  Scenario: pathMatches('/v1/temis/legal-representatives') && methodIs('post') && paramValue('profile_id') != null
    * def profile_id = paramValue('profile_id')
    * legalRepresentativesResponses[profile_id] = request
    * def response = legalRepresentativesResponses
    * def responseStatus = 200

  Scenario: pathMatches('/v1/temis/legal-representatives') && methodIs('get') && paramValue('profile_id') != null
    * def profile_id = paramValue('profile_id')
    * def response = legalRepresentativesResponses[profile_id] != null ?  legalRepresentativesResponses[profile_id] : []
    * def responseStatus = 200

  Scenario: pathMatches('/v1/temis/legal-representative/{id}') && methodIs('post')
    * legalRepresentativeResponses[pathParams.id] = request
    * def profile_id = request.profile_id
    * legalRepresentativesResponses[profile_id] = legalRepresentativesResponses[profile_id] != null ? legalRepresentativesResponses[profile_id] : []
    * legalRepresentativesResponses[profile_id].push(request)
    * def response = legalRepresentativeResponses
    * def responseStatus = 200

  Scenario: pathMatches('/v1/temis/legal-representative/{id}') && methodIs('get')
    * def responseStatus = 200
    * def response = legalRepresentativeResponses[pathParams.id] != null ? legalRepresentativeResponses[pathParams.id] : {}

  #Profile: Get
  Scenario: pathMatches('/v1/temis/profile/{id}') && methodIs('post')
    * profileResponses[pathParams.id] = request
    * def response = profileResponses
    * def responseStatus = 200

  Scenario: pathMatches('/v1/temis/profile/{id}') && methodIs('get')
    * def response = profileResponses[pathParams.id] != null ? profileResponses[pathParams.id] : {}
    * def responseStatus = profileResponses[pathParams.id] != null ? 200 : 404

  Scenario: pathMatches('/v1/temis/profile/{id}/sync') && methodIs('put')
    * profileResponses[pathParams.id] = request
    * print(request)
    * def response = profileResponses[pathParams.id]
    * def responseStatus = 200

  #Profile: Create
  Scenario: pathMatches('/v1/temis/profile') && methodIs('post')
    * def key = request.document_number + request.partner_id + request.offer_type + request.role_type
    * profileResponses[key] = request
    * def response = profileResponses
    * def responseStatus = 201

  Scenario: pathMatches('/v1/temis/profile/create/assert') && methodIs('post')
    * def key = request.document_number + request.partner_id + request.offer_type + request.role_type
    * def saved = profileResponses[key] != null ? profileResponses[key] : {}
    * def equals = karate.match(saved, request).pass
    * def response = { passed: "#(equals)", saved: "#(saved)", sent: "#(request)" }
    * def responseStatus = profileResponses[key] != null ? 200 : 204

  #Profile: Find by document number
  Scenario: pathMatches('/v1/temis/profiles/find_by_document_number') && methodIs('post') && paramValue('role_type') != null && paramValue('document_number') != null && paramValue('partner_id') != null && paramValue('parent_id') != null
    * def roleType = paramValue('role_type')
    * def document = paramValue('document_number')
    * def partnerID = paramValue('partner_id')
    * def parentID = paramValue('parent_id')
    * def key = roleType + document + partnerID + parentID
    * profileResponses[key] = request
    * def response = profileResponses
    * def responseStatus = 200

  Scenario: pathMatches('/v1/temis/profiles/find_by_document_number') && methodIs('get') && paramValue('role_type') != null && paramValue('document_number') != null && paramValue('partner_id') != null && paramValue('parent_id') != null
    * def roleType = paramValue('role_type')
    * def document = paramValue('document_number')
    * def partnerID = paramValue('partner_id')
    * def parentID = paramValue('parent_id')
    * def key = roleType + document + partnerID + parentID
    * def response = profileResponses[key].body != null ? profileResponses[key].body : {}
    * def responseStatus = profileResponses[key].status != null ? profileResponses[key].status : 500

  Scenario: pathMatches('/v1/temis/ownership-structure/{profile_id}') && methodIs('post')
    * ownershipStructureResponses[pathParams.profile_id] = request
    * def response = ownershipStructureResponses
    * def responseStatus = 200

  Scenario: pathMatches('/v1/temis/ownership-structure/{profile_id}') && methodIs('get')
    * def response = ownershipStructureResponses[pathParams.profile_id] != null ? ownershipStructureResponses[pathParams.profile_id] : {}
    * def responseStatus = ownershipStructureResponses[pathParams.profile_id] != null ? 200 : 204

  Scenario: pathMatches('/v1/temis/partner/{partner_id}') && methodIs('post')
    * partnerResponses[pathParams.partner_id] = request
    * def response = partnerResponses
    * def responseStatus = 200

  Scenario: pathMatches('/v1/temis/partner/{partner_id}') && methodIs('get')
    * def responseStatus = 200
    * def response = partnerResponses[pathParams.partner_id] != null ? partnerResponses[pathParams.partner_id] : { name: 'XXX', status: 'ACTIVE'}

 #BoardOfDirectors: Search
  Scenario: pathMatches('/v1/temis/directors') && methodIs('post') && paramValue('profile_id') != null
    * def profile_id = paramValue('profile_id')
    * boardOfDirectorsResponses[profile_id] = request
    * def response = boardOfDirectorsResponses
    * def responseStatus = 200

  Scenario: pathMatches('/v1/temis/directors') && methodIs('get') && paramValue('profile_id') != null
    * def profile_id = paramValue('profile_id')
    * def response = boardOfDirectorsResponses[profile_id] != null ?  boardOfDirectorsResponses[profile_id] : []
    * def responseStatus = 200

  #Enrichment ----------------------------------------------------------------
  Scenario: pathMatches('/temis-enrichment/individual/{documentNumber}') && methodIs('post')
    * enrichmentIndividualResponses[pathParams.documentNumber] = request
    * def response = enrichmentIndividualResponses
    * def responseStatus = 200

  Scenario: pathMatches('/temis-enrichment/individual/{documentNumber}') && methodIs('get')
    * print(requestHeaders)
    * def hasHeaders = requestHeaders["Offer-Type"] != null && requestHeaders["Partner-Id"] != null
    * def responseStatus = enrichmentIndividualResponses[pathParams.documentNumber] != null ? 200 : 204
    * def responseStatus = hasHeaders ? responseStatus : 400
    * def response = enrichmentIndividualResponses[pathParams.documentNumber] != null ? enrichmentIndividualResponses[pathParams.documentNumber] : {}
    * def response = hasHeaders ? response : "{\"error\": \"Empty headers :Offer-Type or Partner-Id\"}"

  Scenario: pathMatches('/temis-enrichment/legal-entity/{legalEntityID}') && methodIs('post')
    * enrichmentLegalEntityResponses[pathParams.legalEntityID] = request
    * def response = enrichmentLegalEntityResponses
    * def responseStatus = 200

  Scenario: pathMatches('/temis-enrichment/legal-entity/{legalEntityID}') && methodIs('get')
    * def hasHeaders = requestHeaders["Offer-Type"] != null && requestHeaders["Partner-Id"] != null
    * def responseStatus = enrichmentLegalEntityResponses[pathParams.legalEntityID] != null ? 200 : 204
    * def responseStatus = hasHeaders ? responseStatus : 400
    * def response = enrichmentLegalEntityResponses[pathParams.legalEntityID] != null ? enrichmentLegalEntityResponses[pathParams.legalEntityID] : {}
    * def response = hasHeaders ? response : "{\"error\": \"Empty headers :Offer-Type or Partner-Id\"}"

  Scenario: pathMatches('/temis-enrichment/ownership-structure/{documentNumber}') && methodIs('post')
    * enrichmentResponses[pathParams.documentNumber] = request
    * def response = enrichmentResponses
    * def responseStatus = 200

  Scenario: pathMatches('/temis-enrichment/ownership-structure/{documentNumber}') && methodIs('get')
    * def hasHeaders = requestHeaders["Offer-Type"] != null && requestHeaders["Partner-Id"] != null
    * def responseStatus = enrichmentResponses[pathParams.documentNumber] != null ? 200: 204
    * def responseStatus = hasHeaders ? responseStatus : 400
    * def response = enrichmentResponses[pathParams.documentNumber] != null ? enrichmentResponses[pathParams.documentNumber] : {}
    * response.legal_entity_id = pathParams.documentNumber
    * def response = hasHeaders ? response : "{\"error\": \"Empty headers :Offer-Type or Partner-Id\"}"


  Scenario: pathMatches('/temis-enrichment/enrich/{documentNumber}') && methodIs('post') && paramValue('profile_id') != null && paramValue('person_type') != null && paramValue('offer_type') != null && paramValue('partner_id') != null && paramValue('role_type') != null
    * def key = pathParams.documentNumber + paramValue('profile_id')+ paramValue('person_type') + paramValue('offer_type') + paramValue('partner_id') + paramValue('role_type')
    * enrichedPersonResponses[key] = request
    * def response = enrichedPersonResponses
    * def responseStatus = 200

  Scenario: pathMatches('/temis-enrichment/enrich/{documentNumber}') && methodIs('get') && paramValue('profile_id') != null && paramValue('person_type') != null && paramValue('offer_type') != null && paramValue('partner_id') != null && paramValue('role_type') != null
    * def key = pathParams.documentNumber + paramValue('profile_id')+ paramValue('person_type') + paramValue('offer_type') + paramValue('partner_id') + paramValue('role_type')
    * def response =  enrichedPersonResponses[key] != null ? enrichedPersonResponses[key] : {}
    * def responseStatus = enrichedPersonResponses[key] != null ? 200: 204

  # Temis Restrictive Lists
  Scenario:  pathMatches('/temis-restrictive-lists/internal-list') && methodIs('post') && paramValue('document_number') != null
    * def document = paramValue("document_number")
    * internalListResponses[document] = request
    * def response = internalListResponses
    * def responseStatus = 200

  Scenario:  pathMatches('/temis-restrictive-lists/internal-list') && methodIs('get') && paramValue('document_number') != null
    * def document = paramValue("document_number")
    * def response = internalListResponses[document] != null ? internalListResponses[document] : []
    * def responseStatus = internalListResponses[document] != null ? 200 : 500

  Scenario:  pathMatches('/temis-restrictive-lists/pep/{documentNumber}') && methodIs('post')
    * def document = pathParams.documentNumber
    * internalListResponsesPEPCOAF[document] = request
    * def response = internalListResponsesPEPCOAF
    * def responseStatus = 201

  Scenario:  pathMatches('/temis-restrictive-lists/pep/{documentNumber}') && methodIs('get')
    * def document = pathParams.documentNumber
    * def response = internalListResponsesPEPCOAF[document] != null ? internalListResponsesPEPCOAF[document] : []
    * def responseStatus = internalListResponsesPEPCOAF[document] != null ? 200 : 204

  #DOA
  Scenario: pathMatches('/extraction') && methodIs('put')
    * def resposeStatus = 200
    * doaResponses[request.profile_id] = request
    * def response = doaResponses

  Scenario: pathMatches('/extraction') && methodIs('post')
    * def resposeStatus = 200
    * def requestID = uuid()
    * def saved = doaResponses[request.profile_id] != null ? doaResponses[request.profile_id] : {}
    * def equals = karate.match(request, saved).pass
    * def response = equals ? { message: "Extraction started", request_id: request.profile_id } : {}

  #Watchlist
  Scenario: pathMatches('/watchlist') && methodIs('post')
    * copy responseWatchListCopy = responseWatchList
    * responseWatchListDocuments[request.document_number] = responseWatchListCopy
    * responseWatchListDocuments[request.document_number].sources = request.sources
    * def responseStatus = 200

  #TemisConfig
  Scenario: pathMatches('/temis-config/cadastral-validation-configs') && methodIs('post')
    * karate.log("request", request)
    * karate.appendTo(cadastralValidationConfigs, request)
    * def responseStatus = 200
    * def response = request

  Scenario: pathMatches('/temis-config/cadastral-validation-configs') && methodIs('get') && paramValue('person_type') != null && paramValue('role_type') != null && paramValue('offer') != null
    * def offer = paramValue('offer')
    * def person_type = paramValue('person_type')
    * def role_type = paramValue('role_type')
    * karate.log("params:", offer, person_type, role_type)
    * def filterByQueryParams = function(element){ return element.offer_type == offer && element.person_type == person_type && element.role_type == role_type}
    * def result = karate.filter(cadastralValidationConfigs, filterByQueryParams)
    * def responseStatus = result.length != 0 ? 200 : 404
    * def response = result.length != 0 ? result : {}


  Scenario: pathMatches('/watchlist/{documentNumber}') && methodIs('get')
    * def responseStatus = 200
    * def responseWatchListArray = []
    * set responseWatchListArray[] = responseWatchListDocuments[pathParams.documentNumber] != null ? responseWatchListDocuments[pathParams.documentNumber] : responseWatchList
    * def response =  responseWatchListArray

  Scenario: pathMatches('/watchlist/company/{companyName}') && methodIs('post')
    * watchlistCompanyResponse[pathParams.companyName] = request
    * def response = watchlistCompanyResponse
    * def responseStatus = 200

  Scenario: pathMatches('/watchlist/company') &&  paramValue('companyName') != null  &&  paramValue('countryCode') != null  && methodIs('get')
    * def companyName = paramValue('companyName')
    * def responseStatus = watchlistCompanyResponse[companyName] != null ? 200 : 500
    * def response = watchlistCompanyResponse[companyName] != null ? watchlistCompanyResponse[companyName] : {}

  #Feature Toggle ------------------
  Scenario: pathMatches('/temis/feature-toggle/feature/{name}') && methodIs('post')
    * featuresResponses[pathParams.name] = request
    * def response = featuresResponses
    * def responseStatus = 200

  Scenario: pathMatches('/temis/feature-toggle/feature/{name}') && methodIs('get')
    * def responseStatus = 200
    * def response = featuresResponses[pathParams.name] != null ? featuresResponses[pathParams.name] : { name: "temis-tree-adapter", active: false }

  #Callback
  Scenario: pathMatches('/callback') && methodIs('post')
    * webPaymentCallbacksReceived[request.entity_id] = request.result
    * print request
    * print webPaymentCallbacksReceived
    * def responseStatus = 200
    * def response = ""

  Scenario: pathMatches('/getResponseFromCallback/{entityId}') && methodIs('get')
    * def responseStatus = 200
    * responseCallback.status = webPaymentCallbacksReceived[pathParams.entityId]
    * def response = responseCallback

  #Postgress
  Scenario: pathMatches('/postgres/read') && methodIs('get')
    # TODO

  Scenario: pathMatches('/postgres/write') && methodIs('post')
    # TODO

  Scenario:
    # catch-all
    * def responseStatus = 404
    * def responseHeaders = { 'Content-Type': 'text/html; charset=utf-8' }
    * def response = <html><body>Not Found</body></html>