package server

import (
	postagemController "api_journal/core/controller/postagensController"
	"api_journal/core/controller/postagensController/postagensDto"
	"api_journal/core/controller/userController"
	"api_journal/core/controller/userController/userDto"
	"api_journal/requester/neslang"
)

func MainServer() {

	neslang.Public[userDto.UserSignin, any]("/auth/cadastro-externo").
		Post(userController.CadastroUsuarioExterno)

	neslang.Public[userDto.UserLogin, any]("/auth/login").
		Post(userController.LoginAluno)

	neslang.Public[userDto.UserLogin, any]("/auth/login-externo").
		Post(userController.LoginUsuarioExterno)

	neslang.Public[any, postagensDto.ListagemQuerys]("/postagens/listar").
		Get(postagemController.GetPostagens)

	neslang.Public[any, postagensDto.PesquisarTituloQuerys]("/postagens/{titulo}").
		Get(postagemController.GetPostagemByTitle)

	neslang.Protected[postagensDto.NovaPostagem, any]("/postagens/postar", "aluno").
		Post(postagemController.PostPostagem)

	neslang.Protected[postagensDto.ComentarioData, any]("/postagens/{postagemId}/comentar", "aluno").
		Post(postagemController.PostComentario,
			postagemController.PostByParamExiste)

	neslang.Protected[postagensDto.ReacaoData, any]("/postagens/{postagemId}/reagir", "aluno").
		Post(postagemController.PostReacao,
			postagemController.PostByParamExiste)

	neslang.Init("4000")

}
