package server

import (
	postagemController "api_journal/core/controller/postagensController"
	"api_journal/core/controller/postagensController/postagensDto"
	"api_journal/core/controller/userController"
	"api_journal/core/controller/userController/userDto"
	"api_journal/requester/neslang"
)

func MainServer() {

	//Ao final de cada rota, você pode ir retornando infinitas funções de request, obrigatoriamente devem retornar um booleano informando o sucesso
	// da operação
	neslang.Public[userDto.UserLogin, any]("/auth/login").
		Post(userController.LoginAluno)
	neslang.Public[userDto.UserLogin, any]("/auth/login-externo").
		Post(userController.LoginAluno)

	neslang.Public[any, postagensDto.ListagemQuerys]("/postagens/listar").
		Get(postagemController.GetPostagens)
	neslang.Public[any, postagensDto.PesquisarTituloQuerys]("/postagens/{titulo}").
		Get(postagemController.GetPostagemByTitle)

	//Em rotas protected você tem como overload de parametros, as permissões de perfis de usuario que vão acessar as rotas
	neslang.Protected[postagensDto.NovaPostagem, any]("/postagens/postar", "aluno").
		Post(postagemController.PostPostagem)
	neslang.Protected[postagensDto.ComentarioData, any]("/postagens/{postagemId}/comentar", "aluno").
		Post(postagemController.PostComentario, postagemController.PostByParamExiste)
	neslang.Protected[postagensDto.ReacaoData, any]("/postagens/{postagemId}/reagir").
		Post(postagemController.PostReacao, postagemController.PostByParamExiste)

	// neslang.Public[any, any]("/favicon.ico").Get(doNothing)
	neslang.Init("4000")

}
