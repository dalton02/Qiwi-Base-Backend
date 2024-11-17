package server

import (
	httpkit "api_journal/requester/http"
	"fmt"
	"net/http"
)

func testFormData(response http.ResponseWriter, request *http.Request) {

	httpkit.AppSucess("Sucesso ao enviar arquivo", make(map[string]string), response)

}

func middlewareTeste2(response http.ResponseWriter, request *http.Request) bool {
	fmt.Println("akk")

	httpkit.AppBadRequest("Ultimo middleware mirreu", response)
	return false
}
