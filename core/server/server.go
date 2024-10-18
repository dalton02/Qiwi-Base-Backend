package server

import (
	postagemController "api_journal/core/controller/postagensController"
	"api_journal/core/controller/postagensController/postagensDto"
	"api_journal/core/controller/userController"
	"api_journal/core/controller/userController/userDto"
	"api_journal/requester/neslang"
)

func MainServer() {

	neslang.Public[userDto.UserLogin, any]("/auth/login").Post(userController.LoginUser)
	neslang.Public[any, postagensDto.ListagemQuerys]("/postagens/listar").Get(postagemController.GetPostagens)
	neslang.Public[any, postagensDto.PesquisarTituloQuerys]("/postagens/{titulo}").Get(postagemController.GetPostagemByTitle)

	neslang.Protected[postagensDto.NovaPostagem, any]("/postagens/postar").Post(postagemController.PostPostagem)
	neslang.Protected[postagensDto.ComentarioData, any]("/postagens/{postagemId}/comentar").Post(postagemController.PostComentario)
	neslang.Protected[postagensDto.ReacaoData, any]("/postagens/{postagemId}/reagir").Post(postagemController.PostReacao)

	neslang.Public[any, any]("/favicon.ico").Get(doNothing)
	neslang.Init("4000")

}
