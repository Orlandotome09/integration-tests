function response(params) {
    return {
        legal_entity_id: params.document_number,
        final_beneficiaries_counted: 0,
        shareholding_sum: 0.0,
        shareholders: []
    };
}