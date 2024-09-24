package server

import (
	"api_journal/core/controller/userController"
	"api_journal/core/controller/userController/userDto"
	httpkit "api_journal/requester/http"
	"api_journal/requester/neslang"
	"net/http"
)

func test(response http.ResponseWriter, request *http.Request) {
	params, _ := httpkit.GetUrlParams(request)
	httpkit.AppSucess("Sucesso", params, response)
}

func MainServer() {

	neslang.Public[userDto.UserLogin, any]("/auth/login").Post(userController.LoginUser)
	neslang.Public[userDto.UserSignin, any]("/auth/sign").Post(userController.SignUser)
	neslang.Protected[userDto.UserLogin, any]("/teste/{id}").Post(test)
	neslang.Public[any, any]("/favicon.ico").Get(doNothing)
	neslang.Init("4000")

}
