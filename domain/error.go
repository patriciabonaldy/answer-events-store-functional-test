package domain

type ErrorResponse struct {
	ErrorCode    string `json:"error_code"`
	ErrorMessage string `json:"error_message"`
}

func (e ErrorResponse) Error() string {
	return e.ErrorMessage
}

type ErrorHttp struct {
	Cause          error  `json:"cause"`
	Message        string `json:"message"`
	ExternalStatus string `json:"error"`
	HTTPStatus     int    `json:"status"`
}

func (e ErrorHttp) Error() string {
	if e.Cause != nil {
		return e.Message + ": " + e.Cause.Error()
	}
	return e.Message
}
