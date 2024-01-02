function fn() {

    var ENV = (java.lang.System.getenv('ENVIRONMENT') == 'production') ? "prd" : "dev"

    var URL_REGISTRATION_INT = "https://api-"+ENV+".bexs.com.br/compliance-int"
    var URL_REGISTRATION_EXT = "https://api-"+ENV+".bexs.com.br/compliance"
    var URL_REGISTRATION_EXT_PUBLIC = "https://api-"+ENV+"-ext.bexs.com.br/compliance"

    var config = {
        baseURLInternal: URL_REGISTRATION_INT,
        baseURLExternal: URL_REGISTRATION_EXT,
        baseURLPublic: URL_REGISTRATION_EXT_PUBLIC,
    };

    karate.configure("connectTimeout", 15000);
    karate.configure("readTimeout", 15000);


    return config;
}