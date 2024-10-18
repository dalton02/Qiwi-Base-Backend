package neslang

import (
	httpkit "api_journal/requester/http"
	"api_journal/requester/validator"
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

func validation[B any, Q any](response http.ResponseWriter, request *http.Request) bool {
	body, err := ioutil.ReadAll(request.Body)
	if err != nil {
		httpkit.GenerateErrorHttpMessage(400, "Erro ao ler o corpo da requisição", response)
	}
	request.Body = ioutil.NopCloser(bytes.NewBuffer(body))

	var dataR B

	if len(body) == 0 {
		body = []byte("{}")
	}

	jsonString := body
	json.Unmarshal(body, &dataR)
	var jsonData map[string]interface{}
	json.Unmarshal([]byte(jsonString), &jsonData)
	keysByLevel := make(map[int][]string)

	errValidacao, hasError := validator.CheckPropretys[B](dataR, validator.ExtractKeysByLevel(jsonData, 1, keysByLevel))
	params, has, maping := extractQueryParams[Q](request)
	errQuerys := ""
	hasErrorQ := false
	if has {
		errQuerys, hasErrorQ = validator.CheckPropretys[Q](params, validator.QueryMap(maping))
	}
	if hasError || hasErrorQ {
		httpkit.GenerateErrorHttpMessage(400, errValidacao+errQuerys, response)
		return false
	}

	return true
}

func generic[B any, Q any](response http.ResponseWriter, request *http.Request, r *HandlerRequest[B, Q], typeRequest string) {
	isSameRequest(typeRequest, request, response)
	if r.middleware == "public" {
		public[B, Q](response, request)
	} else if r.middleware == "protected" {
		protected[B, Q](response, request)
	}
	params, err := extractParams(r.rota, request.URL.Path)
	if err == nil {
		ctx := context.WithValue(request.Context(), "params", params)
		r.controller(response, request.WithContext(ctx))
	}
	r.controller(response, request)
}

func public[B any, Q any](response http.ResponseWriter, request *http.Request) {
	validation[B, Q](response, request)
}

func protected[B any, Q any](response http.ResponseWriter, request *http.Request) {
	validation[B, Q](response, request)
	auth := request.Header.Get("Authorization")
	auth = httpkit.GetBearerToken(auth)
	httpkit.GetJwtInfo(auth, response)
}
