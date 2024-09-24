package neslang

import (
	httpkit "api_journal/requester/http"
	"fmt"
	"net/http"

	"github.com/rs/cors"
)

type Requests interface {
	Post(handlerFunc func(http.ResponseWriter, *http.Request))
	Get(handlerFunc func(http.ResponseWriter, *http.Request))
	Put(handlerFunc func(http.ResponseWriter, *http.Request))
	Patch(handlerFunc func(http.ResponseWriter, *http.Request))
	Delete(handlerFunc func(http.ResponseWriter, *http.Request))
}

type HandlerRequest[B any, Q any] struct {
	rota       string
	middleware string
	controller func(http.ResponseWriter, *http.Request)
}

func Init(porta string) {
	corsHandler := cors.Default().Handler(http.DefaultServeMux)
	err := http.ListenAndServe(":"+porta, corsHandler)
	fmt.Println("Server running in port:" + porta)
	if err == nil {
		fmt.Println("Erro no servidor: ", err)
	}
}

func Public[B any, Q any](rota string) Requests {
	return &HandlerRequest[B, Q]{rota: rota, middleware: "public"}
}
func Protected[B any, Q any](rota string) Requests {
	return &HandlerRequest[B, Q]{rota: rota, middleware: "protected"}
}

func (r *HandlerRequest[B, Q]) Get(handlerFunc func(http.ResponseWriter, *http.Request)) {
	r.controller = handlerFunc
	http.HandleFunc(r.rota, func(response http.ResponseWriter, request *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				if errStr, ok := err.(string); ok && errStr == "common" {
					return
				}
				fmt.Println(err)
				httpkit.AppInternal("internal error"+err.(string), response)
				return
			}
		}()
		generic[B, Q](response, request, r, "GET")
	})
}
func (r *HandlerRequest[B, Q]) Post(handlerFunc func(http.ResponseWriter, *http.Request)) {
	r.controller = handlerFunc
	http.HandleFunc(r.rota, func(response http.ResponseWriter, request *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				if errStr, ok := err.(string); ok && errStr == "common" {
					return
				}
				fmt.Println(err)
				httpkit.AppInternal("internal error", response)
				return
			}
		}()
		generic[B, Q](response, request, r, "POST")
	})
}
func (r *HandlerRequest[B, Q]) Put(handlerFunc func(http.ResponseWriter, *http.Request)) {
	r.controller = handlerFunc
	http.HandleFunc(r.rota, func(response http.ResponseWriter, request *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				if errStr, ok := err.(string); ok && errStr == "common" {
					return
				}
				fmt.Println(err)
				httpkit.AppInternal("internal error", response)
				return
			}
		}()
		generic[B, Q](response, request, r, "PUT")
	})
}
func (r *HandlerRequest[B, Q]) Patch(handlerFunc func(http.ResponseWriter, *http.Request)) {
	r.controller = handlerFunc
	http.HandleFunc(r.rota, func(response http.ResponseWriter, request *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				if errStr, ok := err.(string); ok && errStr == "common" {
					return
				}
				fmt.Println(err)
				httpkit.AppInternal("internal error", response)
				return
			}
		}()
		generic[B, Q](response, request, r, "PATCH")
	})
}
func (r *HandlerRequest[B, Q]) Delete(handlerFunc func(http.ResponseWriter, *http.Request)) {
	r.controller = handlerFunc
	http.HandleFunc(r.rota, func(response http.ResponseWriter, request *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				if errStr, ok := err.(string); ok && errStr == "common" {
					return
				}
				fmt.Println(err)
				httpkit.AppInternal("internal error", response)
				return
			}
		}()
		generic[B, Q](response, request, r, "DELETE")
	})
}
