function response(params) {
    return {
        legal_entity_id: params.document_number,
        final_beneficiaries_counted: 1,
        shareholding_sum: 90.0,
        shareholders: [
            {
                parent_legal_entity: params.document_number,
                shareholding: 100,
                role: "",
                type: "COMPANY",
                name: "Cimentos CIA",
                document_number: params.shareholder1,
                nationality: "BRA",
                birth_date: "11/07/1940",
                pep: false,
                shareholders: [
                    {
                        parent_legal_entity: params.shareholder1,
                        shareholding: 50,
                        role: "",
                        type: "INDIVIDUAL",
                        name: "Maria Joaquina",
                        document_number: params.shareholder2,
                        nationality: "BRA",
                        birth_date: "12/12/1990",
                        pep: false,
                    },
                    {
                        parent_legal_entity: params.shareholder1,
                        shareholding: 40,
                        role: "",
                        type: "INDIVIDUAL",
                        name: "Marcio Felix",
                        document_number: params.shareholder3,
                        nationality: "BRA",
                        birth_date: "11/11/1980",
                        pep: false,
                    }
                ]
            }
        ]
    };
}