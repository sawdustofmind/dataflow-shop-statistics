package dto

const SuccessStatus = "success"

type APISuccess struct {
	Status string `json:"status"`
}

var APISuccessStatus = &APISuccess{Status: SuccessStatus}
