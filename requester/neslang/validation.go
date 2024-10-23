package neslang

import (
	httpkit "api_journal/requester/http"
	"api_journal/requester/validator"
	"bytes"
	"context"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
)

func readBody(request *http.Request, response http.ResponseWriter) ([]byte, error) {
	var buf bytes.Buffer
	tee := io.TeeReader(request.Body, &buf)
	body, err := ioutil.ReadAll(tee)
	if err != nil {
		return body, err
	}
	if len(body) == 0 {
		body = []byte("{}")
	}
	request.Body = ioutil.NopCloser(bytes.NewBuffer(buf.Bytes()))
	return body, nil
}

func validation[B any, Q any](response http.ResponseWriter, request *http.Request) (bool, *http.Request) {
	body, err := readBody(request, response)
	if err != nil {
		httpkit.GenerateErrorHttpMessage(400, "Erro ao ler o corpo da requisição", response)
		return false, request
	}
	var dataR B
	jsonString := body
	json.Unmarshal(body, &dataR)
	var jsonData map[string]interface{}
	json.Unmarshal([]byte(jsonString), &jsonData)
	keysByLevel := make(map[int][]string)
	ctx := context.WithValue(request.Context(), "original_body", body)
	request = request.WithContext(ctx)
	errValidacao, hasError := validator.CheckPropretys[B](dataR, validator.ExtractKeysByLevel(jsonData, 1, keysByLevel))
	params, has, maping := extractQueryParams[Q](request)
	errQuerys := ""
	hasErrorQ := false
	if has {
		errQuerys, hasErrorQ = validator.CheckPropretys[Q](params, validator.QueryMap(maping))
	}
	if hasError || hasErrorQ {
		httpkit.GenerateErrorHttpMessage(400, errValidacao+errQuerys, response)
		return false, request
	}

	return true, request
}

func generic[B any, Q any](response http.ResponseWriter, request *http.Request, r *HandlerRequest[B, Q], typeRequest string) {
	isSameRequest(typeRequest, request, response)
	var valid bool
	if r.middleware == "public" {
		valid, request = public[B, Q](response, request)
	} else if r.middleware == "protected" {
		valid, request = protected[B, Q](response, request)
	}
	if !valid {
		return
	}
	params, err := extractParams(r.rota, request.URL.Path)
	if err == nil {
		ctx := context.WithValue(request.Context(), "params", params)
		r.controller(response, request.WithContext(ctx))
		return
	}
	r.controller(response, request)
}

func public[B any, Q any](response http.ResponseWriter, request *http.Request) (bool, *http.Request) {
	valid, request := validation[B, Q](response, request)
	return valid, request
}

func protected[B any, Q any](response http.ResponseWriter, request *http.Request) (bool, *http.Request) {
	valid, request := validation[B, Q](response, request)
	auth := request.Header.Get("Authorization")
	auth = httpkit.GetBearerToken(auth)
	_, err := httpkit.GetJwtInfo(auth)
	if err != nil {
		valid = false
		httpkit.AppForbidden("Token invalido/Expirado, faça login novamente", response)
	}
	return valid, request
}
