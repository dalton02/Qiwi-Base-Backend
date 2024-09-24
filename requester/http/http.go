package httpkit

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type HttpMessage[T any] struct {
	StatusCode int    `json:"statusCode"`
	Data       T      `json:"data"`
	Message    string `json:"message"`
}

func GenerateHttpMessage[T any](statusCode int, data T, message string, response http.ResponseWriter) {
	var dataResponse = HttpMessage[T]{
		StatusCode: statusCode,
		Data:       data,
		Message:    message,
	}
	json.Marshal(dataResponse)
	response.Header().Set("Content-Type", "application/json")
	response.WriteHeader(statusCode)
	if err := json.NewEncoder(response).Encode(dataResponse); err != nil {
		http.Error(response, err.Error(), http.StatusInternalServerError)
	}
	panic("common")

}

func GenerateErrorHttpMessage(statusCode int, message string, response http.ResponseWriter) {
	var dataResponse = HttpMessage[any]{
		StatusCode: statusCode,
		Data:       nil,
		Message:    message,
	}
	json.Marshal(dataResponse)
	response.Header().Set("Content-Type", "application/json")
	response.WriteHeader(statusCode)
	if err := json.NewEncoder(response).Encode(dataResponse); err != nil {
		http.Error(response, err.Error(), http.StatusInternalServerError)
	}

	panic("common")

}

func GenerateErrorHttpMessageNonPanic(statusCode int, message string, response http.ResponseWriter) {
	var dataResponse = HttpMessage[any]{
		StatusCode: statusCode,
		Data:       nil,
		Message:    message,
	}
	fmt.Println(dataResponse)
	json.Marshal(dataResponse)
	response.Header().Set("Content-Type", "application/json")
	response.WriteHeader(statusCode)
	if err := json.NewEncoder(response).Encode(dataResponse); err != nil {
		http.Error(response, err.Error(), http.StatusInternalServerError)
	}
}
