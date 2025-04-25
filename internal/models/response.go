package models

import (
	"encoding/json"
	"net/http"
)

type ErrorResponse struct {
	Status int    `json:"status"`
	Error  string `json:"error"`
}

type DataResponse struct {
	Status int `json:"status"`
	Data   any `json:"data"`
}

func NewDataResponse(status int, data any) *DataResponse {
	return &DataResponse{
		Status: status,
		Data:   data,
	}
}

func NewErrorResponse(status int, err string) *ErrorResponse {
	return &ErrorResponse{
		Status: status,
		Error:  err,
	}
}

func ResponseWithJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(payload)
}
