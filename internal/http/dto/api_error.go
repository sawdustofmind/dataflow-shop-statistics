package dto

const ErrorStatus = "error"

type APIError struct {
	Status string `json:"status"`
	Error  string `json:"error"`
}

func NewAPIError(err error) APIError {
	return APIError{
		Status: ErrorStatus,
		Error:  err.Error(),
	}
}
