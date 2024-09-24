package userController

import (
	"api_journal/core/controller/userController/userDto"
	"api_journal/core/server/shared"

	"api_journal/core/service/userService"

	jwtkit "api_journal/requester/cripty"
	httpkit "api_journal/requester/http"
	"encoding/json"
	"net/http"
)

// LoginUser autentica o usuário com base no login e senha fornecidos.
// @Summary Autenticação de Usuário
// @Description Autentica um usuário e retorna um token JWT com os dados do usuário.
// @Tags Autenticação
// @Accept json
// @Produce json
// @Param login body userDto.UserLogin true "Dados de Login"
// @Failure 400 {object} map[string]string "Erro ao processar a requisição"
// @Failure 500 {object} map[string]string "Erro interno do servidor"
// @Router /login [post]
func LoginUser(response http.ResponseWriter, request *http.Request) {
	var userLoginRequest userDto.UserLogin
	decoder := json.NewDecoder(request.Body)
	decoder.Decode(&userLoginRequest)
	userData := userService.CheckUsuario(shared.DB, userLoginRequest, response)
	jwtkit.GenerateJwt(userData)
	httpkit.AppSucess("Usuario logado com sucesso", make(map[string]string), response)
}

func SignUser(response http.ResponseWriter, request *http.Request) {
	var userRequest userDto.UserSignin
	decoder := json.NewDecoder(request.Body)
	decoder.Decode(&userRequest)
	lastId := userService.InsertUsuario(shared.DB, userRequest, response)
	userData := userService.GetUserById(shared.DB, lastId, response)
	httpkit.AppSucess("Usuario cadastrado com sucesso", userData, response)
}

func Test(response http.ResponseWriter, request *http.Request) {

}
