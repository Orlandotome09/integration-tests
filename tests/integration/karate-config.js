function fn() {

  var COMPLIANCE_API_URL = java.lang.System.getenv('COMPLIANCE_API_URL');

  if (COMPLIANCE_API_URL == undefined) COMPLIANCE_API_URL = "http://localhost:8193/compliance-int"

  var MOCKS_URL = java.lang.System.getenv('MOCKS_URL');

  if (MOCKS_URL == undefined) MOCKS_URL = "http://localhost:9093"

  var host = java.lang.System.getenv('PUBSUB_PROJECT_HOST') != undefined ? java.lang.System.getenv('PUBSUB_PROJECT_HOST') : 'localhost:8681';

  var MessagePublisher = Java.type('bexs.MessagePublisher')

  var ComplicanceCommandPublisher = new MessagePublisher(host, 'local-project','registration_topic')
  var RegistrationEventsPublisher = new MessagePublisher(host, 'local-project','temis-registration-events')

  var EnrichmentPublisher = new MessagePublisher(host, 'local-project','temis-enrichment-person-events')

  var METRICS_API_URL = java.lang.System.getenv('METRICS_API_URL');
  if (METRICS_API_URL == undefined) METRICS_API_URL = "http://localhost:7777"

  var METRICS_EVENTS_URL = java.lang.System.getenv('METRICS_EVENTS_URL');
  if (METRICS_EVENTS_URL == undefined) METRICS_EVENTS_URL = "http://localhost:7777"

  var schemas = read('classpath:features/_data/schemas.json');

  var config = {
    metricsApiURL: METRICS_API_URL,
    metricsEventsURL: METRICS_EVENTS_URL,
    baseURLCompliance: COMPLIANCE_API_URL,
    mockURL : MOCKS_URL,

    RegistrationEventsPublisher: RegistrationEventsPublisher,
    ComplicanceCommandPublisher: ComplicanceCommandPublisher,
    EnrichmentPublisher: EnrichmentPublisher,
    
    OK : 200,
    CREATED : 201,
    NO_CONTENT : 204,
    BAD_CONTENT : 400,
    NOT_FOUND : 404,

    uuid:uuid,
    sleep:sleep,
    funcNow:now,
    generateDateOfBirthForAge:generateDateOfBirthForAge,
    CNPJGenerator:gerarCNPJ,
    CPFGenerator: gerarCPF,
    DocumentNormalizer: normalizeDocument,
    CreateSingleLevelCatalog: createSingleLevelCatalog,
    CreateMultiLevelCatalog: createMultiLevelCatalog,
    CreateProfileIndividual:createProfileIndividual,
    CreateProfileCompany: createProfileCompany,

    CallbackURL: MOCKS_URL+ "/callback",
    Engine:{
      Profile: "PROFILE",
      Contract: "CONTRACT"
    },
    EntityType:{
      Profile: "PROFILE",
      Contract: "CONTRACT",
      LegalRepresentative: "LEGAL_REPRESENTATIVE",
      Shareholder:"SHAREHOLDER",
      Director: "DIRECTOR",
      ComplianceState: "COMPLIANCE_STATE"
    },
    EventType:{
      State:{
        Created: "STATE_CREATED",
        Changed: "STATE_CHANGED"
      }
    },
    ProfileType:{
      Individual: "INDIVIDUAL",
      Company: "COMPANY",
    },
    RoleType:{
      Customer: "CUSTOMER",
      Merchant: "MERCHANT",
      Shareholder: "SHAREHOLDER",
      LegalRepresentative: "LEGAL_REPRESENTATIVE",
      Director: "DIRECTOR",
      Counterparty: "COUNTERPARTY"
    },
    RuleResult:{
      Created: "CREATED",
      Analysing: "ANALYSING",
      Incomplete: "INCOMPLETE",
      Ignored: "IGNORED",
      Approved: "APPROVED",
      Rejected: "REJECTED",
     
    },
  
    RuleSet:{
      ActivityRisk: "ACTIVITY_RISK",
      Blacklist: "BLACKLIST",
      BoardOfDirectors: "BOARD_OF_DIRECTORS",
      Bureau: "SERASA_BUREAU",
      CafAnalysis: "CAF_ANALYSIS",
      Incomplete: "INCOMPLETE",
      LegalRepresentatives: "LEGAL_REPRESENTATIVES",
      OwnershipStructure: "OWNERSHIP_STRUCTURE",
      Pep: "PEP",
      Watchlist: "WATCHLIST",
      MinimumBilling: "MINIMUM_BILLING",
      MinimumIncome: "MINIMUM_INCOME",
    },
    RuleName:{
      HighRiskActivity: "HIGH_RISK_ACTIVITY",
      Blacklist: "BLACKLIST",
      BoardOfDirectorsResult: "BOARD_OF_DIRECTORS_RESULT",
      BoardOfDirectorsComplete: "BOARD_OF_DIRECTORS_COMPLETE",
      CustomerNotFoundInSerasa: "CUSTOMER_NOT_FOUND_IN_SERASA",
      CustomerHasProblemsInSerasa: "CUSTOMER_HAS_PROBLEMS_IN_SERASA",
      CafAnalysis: "CAF_ANALYSIS_RESULT",
      LegalRepresentativeResult: "LEGAL_REPRESENTATIVES_RESULT",
      RequiredFieldsNotFound: "REQUIRED_FIELDS_NOT_FOUND",
      AddressNotFound: "ADDRESS_NOT_FOUND",
      DocumentNotFound: "DOCUMENT_NOT_FOUND",
      Pep: "PEP",
      Shareholding:"SHAREHOLDING",
      Shareholders: "SHAREHOLDERS",
      Watchlist: "WATCHLIST",
      InsufficientBilling: "INSUFFICIENT_BILLING",
	    InsufficientIncome: "INSUFFICIENT_INCOME",
    },

    profileStateEventContentSchema: schemas.profileStateEventContent,
    personStateEventContentSchema: schemas.personStateEventContent,
  };

  karate.configure("connectTimeout", 15000);
  karate.configure("readTimeout", 15000);
  karate.configure('retry',{ count:10, interval:1500});

  return config;
}

function uuid(){ return java.util.UUID.randomUUID() + '' }
function sleep(milliseconds){ java.lang.Thread.sleep(milliseconds) }
function now(){
  return (new Date()).toJSON();
}

function generateDateOfBirthForAge(age){
  var today = new Date();
  var currentYear = today.getFullYear();
  var dateOfBirth = currentYear - age;
  return (new Date(dateOfBirth, 1, 1)).toJSON();
}

function gerarCNPJ() {
  var n = 9;
  var n1 = randomiza(n);
  var n2 = randomiza(n);
  var n3 = randomiza(n);
  var n4 = randomiza(n);
  var n5 = randomiza(n);
  var n6 = randomiza(n);
  var n7 = randomiza(n);
  var n8 = randomiza(n);
  var n9 = randomiza(n);
  var n10 = 0;
  var n11 = 0
  var n12 = 1;
  var d1 = n12*2+n11*3+n10*4+n9*5+n8*6+n7*7+n6*8+n5*9+n4*2+n3*3+n2*4+n1*5;
  d1 = 11 - ( mod(d1,11) );
  if (d1>=10) d1 = 0;
  var d2 = d1*2+n12*3+n11*4+n10*5+n9*6+n8*7+n7*8+n6*9+n5*2+n4*3+n3*4+n2*5+n1*6;
  d2 = 11 - ( mod(d2,11) );
  if (d2>=10) d2 = 0;
  resultado = ''+n1+n2+'.'+n3+n4+n5+'.'+n6+n7+n8+'/'+n9+n10+n11+n12+'-'+d1+d2;

  return resultado

  function randomiza(n) {
      return Math.round(Math.random() * n);
  }
  
  function mod(dividendo, divisor) {
      return Math.round(dividendo - (Math.floor(dividendo / divisor) * divisor));
  }
}

function gerarCPF() {
  var n = 9;
  var n1 = randomiza(n);
  var n2 = randomiza(n);
  var n3 = randomiza(n);
  var n4 = randomiza(n);
  var n5 = randomiza(n);
  var n6 = randomiza(n);
  var n7 = randomiza(n);
  var n8 = randomiza(n);
  var n9 = randomiza(n);
  var d1 = n9 * 2 + n8 * 3 + n7 * 4 + n6 * 5 + n5 * 6 + n4 * 7 + n3 * 8 + n2 * 9 + n1 * 10;
  d1 = 11 - (mod(d1, 11));
  if (d1 >= 10) d1 = 0;
  var d2 = d1 * 2 + n9 * 3 + n8 * 4 + n7 * 5 + n6 * 6 + n5 * 7 + n4 * 8 + n3 * 9 + n2 * 10 + n1 * 11;
  d2 = 11 - (mod(d2, 11));
  if (d2 >= 10) d2 = 0;
  var cpf = '' + n1 + n2 + n3 + n4 + n5 + n6 + n7 + n8 + n9 + d1 + d2;

  return cpf

  function randomiza(n) {
      return Math.round(Math.random() * n);
  }
  
  function mod(dividendo, divisor) {
      return Math.round(dividendo - (Math.floor(dividendo / divisor) * divisor));
  }
}

function normalizeDocument(input) {
  output = input.replace('/','');
  output = output.replace('.','');
  output = output.replace('.','');
  output = output.replace('-','');
  return output;
}

function createSingleLevelCatalog(params) {

  var catalog = {
      offer_type: params.offer_type,
      role_type: params.role_type,
      person_type: params.person_type,
      validation_steps: [
          {
              step_number: 0,
              skip_for_approval: false,
              rules_config: params.rules_config,
          }
      ],
      product_config: {
          create_bexs_account: params.account_flag,
          tree_integration: false,
          limit_integration: params.limit_flag != undefined ? params.limit_flag : false,
          enrich_profile_with_bureau_data: params.enrich_flag != undefined ? params.enrich_flag : false,
      }
  }

  if (params.partner_id != undefined) {
      catalog.partner_id = params.partner_id
  }

  return catalog
}

function createMultiLevelCatalog(params) {
  var catalog = {
      offer_type: params.offer_type,
      role_type: params.role_type,
      person_type: params.person_type,
      validation_steps: [],
      product_config: {
          create_bexs_account: params.account_flag != undefined ? params.account_flag : false,
          enrich_profile_with_bureau_data: params.enrich_flag != undefined ? params.enrich_flag : false,
          tree_integration: params.tree_flag != undefined ? params.tree_flag : false,
          limit_integration: params.limit_flag != undefined ? params.limit_flag : false
      }
  }

  for (index = 0; index < params.steps.length; index++) {
      if (params.steps[index].rules_config != undefined) {
          step = {
              step_number: index | 0,
              skip_for_approval: false,
              rules_config: params.steps[index].rules_config
          }

          catalog.validation_steps.push(step)
      }
  }

  return catalog
}

function createProfileIndividual(params) {

  return {
      profile_id: params.profile_id,
      partner_id: params.partner_id,
      offer_type: params.offer_type,
      role_type: params.role_type,
      profile_type: "INDIVIDUAL",
      document_number: params.document_number,
      callback_url: params.callback_url,
      email: "teste@teste.com",
      individual: {
          first_name: "Test",
          last_name: "Temis",
          date_of_birth_inputted: "1991-07-30T00:00:00Z",
          phones: [
              {
                  type: "comercial",
                  number: "123",
                  country_code: "55",
                  area_code: "11"
              }
          ]
      },
      parent_id: params.parent_id
  }
}

function createProfileCompany(params) {

  return {
      profile_id: params.profile_id,
      partner_id: params.partner_id,
      offer_type: params.offer_type,
      role_type: params.role_type,
      profile_type: "COMPANY",
      document_number: params.document_number,
      callback_url: params.callback_url,
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
}
