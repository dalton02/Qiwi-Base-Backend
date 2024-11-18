package userController

import (
	"api_journal/core/controller/userController/userDto"
	"api_journal/core/server/shared"
	"api_journal/core/service/userService"
	httpkit "api_journal/requester/http"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

// @Summary Cadastro de Usuário Externo
// @Description Autentica um usuário externo e retorna um token JWT com os dados do usuário.
// @Tags Autenticação
// @Accept json
// @Produce json
// @Param cadastro body userDto.UserSignin true "Dados de Cadastro"
// @Failure 400 {object} map[string]string "Erro ao processar a requisição"
// @Failure 500 {object} map[string]string "Erro interno do servidor"
// @Router /auth/cadastro-externo [post]
func CadastroUsuarioExterno(response http.ResponseWriter, request *http.Request) {
	var user userDto.UserSignin
	json.NewDecoder(request.Body).Decode(&user)
	idUsuario, err := userService.InsertUsuarioExterno(shared.DB, user, response)
	if err != nil {
		return
	}
	var httpResposta userDto.UserData
	httpResposta.Id = idUsuario
	httpResposta.Nome = user.Nome
	httpResposta.Login = user.Login
	httpkit.AppSucessCreate("Usuário externo cadastrado com sucesso", httpResposta, response)
}

// @Summary Autenticação de Usuário Externo
// @Description Autentica um usuário externo e retorna um token JWT com os dados do usuário.
// @Tags Autenticação
// @Accept json
// @Produce json
// @Param login body userDto.UserLogin true "Dados de Login"
// @Failure 400 {object} map[string]string "Erro ao processar a requisição"
// @Failure 500 {object} map[string]string "Erro interno do servidor"
// @Router /auth/login-externo [post]
func LoginUsuarioExterno(response http.ResponseWriter, request *http.Request) {
	var user userDto.UserLogin
	json.NewDecoder(request.Body).Decode(&user)
	userData, err := userService.GetUserByLoginPass(shared.DB, user.Login, user.Senha)
	if err != nil {
		httpkit.AppBadRequest("Credenciais incorretas", response)
		return
	}
	jwtInfo := map[string]interface{}{
		"login":  userData.Login,
		"nome":   userData.Nome,
		"id":     userData.Id,
		"perfil": "externo",
	}
	token, err := httpkit.GenerateJwt(jwtInfo)
	resposta := map[string]interface{}{
		"token": token,
	}
	httpkit.AppSucess("Logado com sucesso", resposta, response)

}

// @Summary Autenticação de Usuário Aluno
// @Description Autentica um aluno e retorna um token JWT com os dados do usuário.
// @Tags Autenticação
// @Accept json
// @Produce json
// @Param login body userDto.UserLogin true "Dados de Login"
// @Failure 400 {object} map[string]string "Erro ao processar a requisição"
// @Failure 500 {object} map[string]string "Erro interno do servidor"
// @Router /auth/login [post]
func LoginAluno(response http.ResponseWriter, request *http.Request) {
	var user userDto.UserLogin
	json.NewDecoder(request.Body).Decode(&user)
	data, err := getCookiesFromSigaa(user.Login, user.Senha)
	if err != nil {
		httpkit.AppBadRequest(err.Error(), response)
		return
	}
	codigo, _ := strconv.Atoi(data["codigo"])
	alunoInfo := userDto.AlunoData{
		Login:  user.Login,
		Curso:  data["curso"],
		Codigo: codigo,
		Nome:   data["nome"],
	}
	userInfoSignin := userDto.UserSignin{
		Login: user.Login,
		Senha: "",
		Nome:  alunoInfo.Nome,
	}

	alunoInfo.Id, err = userService.GetAlunoByCodigo(shared.DB, alunoInfo.Codigo)

	//Cadastro de usuario, caso não exista ainda
	if err != nil {
		idUsuario, err2 := userService.InsertUsuarioAluno(shared.DB, userInfoSignin, response)
		alunoInfo.Id = idUsuario
		if err2 != nil {
			return
		}
		_, err3 := userService.InsertAlunoAndRelate(shared.DB, alunoInfo, idUsuario, response)
		if err3 != nil {
			return
		}
	}

	jwtInfo := map[string]interface{}{
		"login":  alunoInfo.Login,
		"curso":  alunoInfo.Curso,
		"codigo": alunoInfo.Codigo,
		"nome":   alunoInfo.Nome,
		"id":     alunoInfo.Id,
		"perfil": "aluno",
	}
	fmt.Println(jwtInfo)
	token, err := httpkit.GenerateJwt(jwtInfo)
	data["token"] = token
	httpkit.GenerateHttpMessage(200, data, "Login bem sucedido", response)
}
