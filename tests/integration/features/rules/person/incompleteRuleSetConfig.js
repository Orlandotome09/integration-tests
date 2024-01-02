function getRuleSetConfigIncomplete(params){
    return {
        incomplete: {
            date_of_birth_required: true,
            phone_number_required: true,
            email_required: true,
            address_required: true,
            last_name_required: true,
            documents_required: [
                {
                    document_type: "IDENTIFICATION",
                    file_required: true,
                    pending_on_approval: false
                },
                {
                    document_type: "REGISTRATION_FORM",
                    file_required: true
                }
            ],
            pep_required: true,
        }
    }
}