function response(params) {
    return {
        legal_entity_id: params.document_number,
        final_beneficiaries_counted: 2,
        shareholding_sum: 100.0,
        shareholders: [
            {
                shareholder_id: params.shareholderID1,
                parent_legal_entity: params.document_number,
                shareholding: 100,
                role: "",
                type: "COMPANY",
                name: "RAZAO SOCIAL",
                document_number: params.shareholder1,
                nationality: "",
                birth_date: "11/12/2011",
                shareholders: [
                    {
                        shareholder_id: params.shareholderID2,
                        parent_legal_entity: params.shareholder1,
                        shareholding: 50,
                        role: "",
                        type: "INDIVIDUAL",
                        name: "sócio um",
                        document_number: params.shareholder2,
                        nationality: "",
                        birth_date: ""
                    },
                    {
                        shareholder_id: params.shareholderID3,
                        parent_legal_entity: params.shareholder1,
                        shareholding: 25,
                        role: "",
                        type: "INDIVIDUAL",
                        name: "sócio dois",
                        document_number: params.shareholder3,
                        nationality: "",
                        birth_date: ""
                    },
                    {
                        shareholder_id: params.shareholderID4,
                        parent_legal_entity: params.shareholder1,
                        shareholding: 25,
                        role: "",
                        type: "INDIVIDUAL",
                        name: "sócio três",
                        document_number: params.shareholder4,
                        nationality: "",
                        birth_date: ""
                    }
                ]
            }
        ]
    };
}