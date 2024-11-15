package server

import (
	httpkit "api_journal/requester/http"
	"fmt"
	"net/http"
)

func middlewareTeste(response http.ResponseWriter, request *http.Request) bool {
	return true
}

func middlewareTeste2(response http.ResponseWriter, request *http.Request) bool {
	fmt.Println("akk")

	httpkit.AppBadRequest("Ultimo middleware mirreu", response)
	return false
}
