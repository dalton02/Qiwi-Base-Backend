package userController

import (
	"api_journal/core/controller/userController/userDto"
	"api_journal/core/server/shared"
	"api_journal/core/service/userService"
	httpkit "api_journal/requester/http"
	"encoding/json"
	"net/http"
	"strconv"
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
func LoginAluno(response http.ResponseWriter, request *http.Request) {
	var user userDto.UserLogin
	json.NewDecoder(request.Body).Decode(&user)
	data, err := getCookiesFromSigaa(user.Login, user.Senha)
	if err != nil {
		httpkit.AppBadRequest(err.Error(), response)
		return
	}
	codigo, _ := strconv.Atoi(data["codigo"])
	userInfo := userDto.UserData{
		Login:  user.Login,
		Curso:  data["curso"],
		Codigo: codigo,
		Nome:   data["nome"],
	}

	_, err = userService.GetAlunoByCodigo(shared.DB, userInfo.Codigo)

	if err != nil {

		idUsuario, err2 := userService.InsertUsuario(shared.DB, userInfo, "aluno", response)

		if err2 != nil {
			return
		}

		_, err3 := userService.InsertAlunoAndRelate(shared.DB, userInfo, idUsuario, response)
		if err3 != nil {
			return
		}

	}
	jwtInfo := map[string]interface{}{
		"login":  userInfo.Login,
		"curso":  userInfo.Curso,
		"codigo": userInfo.Codigo,
		"nome":   userInfo.Nome,
		"perfil": "aluno",
	}

	token, err := httpkit.GenerateJwt(jwtInfo)
	data["token"] = token
	httpkit.GenerateHttpMessage(200, data, "Login bem sucedido", response)
}
