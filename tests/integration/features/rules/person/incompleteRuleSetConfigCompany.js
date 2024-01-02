function getRuleSetConfigIncomplete(params){
    return {
        incomplete: {
            address_required: true,
            documents_required: [
                {
                    document_type: "CORPORATE_DOCUMENT",
                    file_required: true
                },
                {
                   document_type: "APPOINTMENT_DOCUMENT",
                   document_sub_type: "MINUTES_OF_ELECTION",
                   file_required: true,
                   conditions: [ 
                        {
                            field_name: "LEGAL_NATURE",
                            values: ["2046","2143"],
                        }
                   ]
                },
                {
                   document_type: "CONSTITUTION_DOCUMENT",
                   document_sub_type: "STATUTE_SOCIAL",
                   file_required: true,
                   conditions: [
                        { 
                            field_name: "LEGAL_NATURE",
                            values: ["2046","2143"]
                        }                   
                   ]
                }
            ]
        }
    }
}