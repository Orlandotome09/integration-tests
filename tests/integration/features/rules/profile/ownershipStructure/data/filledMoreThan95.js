function response(params) {
    return {
        legal_entity_id: "74385930000147",
        final_beneficiaries_counted: 3,
        shareholding_sum: 96.0,
        shareholders: [
            {
                shareholder_id: 'b71c8c85-4730-42a5-8d2c-f87d03c7c0e5',
                parent_legal_entity: "74385930000147",
                shareholding: 80,
                role: "",
                type: "COMPANY",
                name: "Empresa Um",
                document_number: "31551591000119",
                nationality: "",
                birth_date: "22/12/2010",
                shareholders: [
                    {
                        shareholder_id: '88f19545-f295-4559-88ce-e992860c095f',
                        parent_legal_entity: "31551591000119",
                        shareholding: 50,
                        role: "",
                        type: "INDIVIDUAL",
                        name: "Sócio um",
                        document_number: "84721895038",
                        nationality: "",
                        birth_date: ""
                    },
                    {
                        shareholder_id: '251f8e8b-754d-456b-91ee-cbca9a14bdf0',
                        parent_legal_entity: "31551591000119",
                        shareholding: 50,
                        role: "",
                        type: "INDIVIDUAL",
                        name: "Sócio dois",
                        document_number: "09522116025",
                        nationality: "",
                        birth_date: ""
                    }
                ]
            },
            {
                shareholder_id: 'd004c2d3-181b-46aa-87de-b311632334f2',
                parent_legal_entity: "74385930000147",
                shareholding: 20,
                role: "",
                type: "INDIVIDUAL",
                name: "Sócio um",
                document_number: "84721895038",
                nationality: "",
                birth_date: ""
            },
        ]
    };
}