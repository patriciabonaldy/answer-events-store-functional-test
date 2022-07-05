package domain

type AnswerRequest struct {
	MerchantOperationReference string `json:"merchant_operation_reference"`
	ProcessorID                string `json:"processor_id"`
	CountryCode                string `json:"country_code"`
	MerchantID                 string `json:"merchant_id"`
	Amount                     int32  `json:"amount"`
}

type Response struct {
	AcquirerTransactionID string            `json:"acquirer_transaction_id"`
	ResponseCode          string            `json:"response_code"`
	ResponseMessage       string            `json:"response_message"`
	AuthorizationCode     string            `json:"authorization_code,omitempty"`
	ICCRelatedData        string            `json:"icc_related_data,omitempty"`
	References            map[string]string `json:"references,omitempty"`
}
