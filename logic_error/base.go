package logic_error

type Errors struct {
	// Code is the error code
	Code int `json:"code"`
	// Message is the error message
	Message string `json:"message"`
}

func (e *Errors) Error() string {
	return e.Message
}

func New(code int, message string) *Errors {
	return &Errors{
		Code:    code,
		Message: message,
	}
}
