function response(params) {
    return {
        legal_entity_id: params.document_number,
        final_beneficiaries_counted: 3,
        shareholding_sum: 96.0,
        shareholders: [
            {
                shareholder_id: params.shareholderID1,
                parent_legal_entity: params.document_number,
                shareholding: 80,
                role: "",
                type: "COMPANY",
                name: "Empresa Um",
                document_number: params.shareholder1,
                nationality: "",
                birth_date: "22/12/2010",
                shareholders: [
                    {
                        shareholder_id: params.shareholderID2,
                        parent_legal_entity: params.shareholder1,
                        shareholding: 50,
                        role: "",
                        type: "INDIVIDUAL",
                        name: "Sócio um",
                        document_number: params.shareholder2,
                        nationality: "",
                        birth_date: ""
                    },
                    {
                        shareholder_id: params.shareholderID3,
                        parent_legal_entity: params.shareholder1,
                        shareholding: 50,
                        role: "",
                        type: "INDIVIDUAL",
                        name: "Sócio dois",
                        document_number: params.shareholder3,
                        nationality: "",
                        birth_date: ""
                    }
                ]
            },
            {
                shareholder_id: params.shareholderID4,
                parent_legal_entity: params.document_number,
                shareholding: 20,
                role: "",
                type: "INDIVIDUAL",
                name: "Sócio um",
                document_number: params.shareholder4,
                nationality: "",
                birth_date: ""
            },
        ]
    };
}