function response(params) {
    return {
        legal_entity_id: "63684948000123",
        final_beneficiaries_counted: 3,
        shareholding_sum: 96.0,
        shareholders: [
            {
                parent_legal_entity: "63684948000123",
                shareholding: 80,
                role: "",
                type: "LEGAL_ENTITY",
                name: "Empresa Um",
                document_number: "30085754000152",
                nationality: "",
                birth_date: "22/12/2010",
                shareholders: [
                    {
                        parent_legal_entity: "30085754000152",
                        shareholding: 50,
                        role: "",
                        type: "INDIVIDUAL",
                        name: "Sócio um",
                        document_number: "33983939007",
                        nationality: "",
                        birth_date: ""
                    },
                    {
                        parent_legal_entity: "30085754000152",
                        shareholding: 50,
                        role: "",
                        type: "INDIVIDUAL",
                        name: "Sócio dois",
                        document_number: "22401939059",
                        nationality: "",
                        birth_date: ""
                    }
                ]
            },
            {
                parent_legal_entity: "63684948000123",
                shareholding: 20,
                role: "",
                type: "INDIVIDUAL",
                name: "Sócio um",
                document_number: "33983939007",
                nationality: "",
                birth_date: ""
            },
        ]
    };
}