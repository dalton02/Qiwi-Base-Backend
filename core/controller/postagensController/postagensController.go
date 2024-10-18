package postagensController

import (
	"api_journal/core/controller/postagensController/postagensDto"
	"api_journal/core/server/shared"
	"api_journal/core/service/postagensService"
	"api_journal/core/service/userService"
	httpkit "api_journal/requester/http"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

func GetPostagemByTitle(response http.ResponseWriter, request *http.Request) {
	params, _ := httpkit.GetUrlParams(request)
	titulo := params.Param["titulo"]
	tipo := request.URL.Query().Get("tipo")
	fmt.Println(titulo)
	listagem := postagensService.GetPostagemByTitle(shared.DB, titulo, tipo, response)
	httpkit.AppSucess("Listagem bem sucedida", listagem, response)
}

func GetPostagens(response http.ResponseWriter, request *http.Request) {
	paginaStr := request.URL.Query().Get("pagina")
	limiteStr := request.URL.Query().Get("limite")
	pesquisa := request.URL.Query().Get("pesquisa")
	pagina, err := strconv.Atoi(paginaStr)
	if err != nil || pagina < 1 {
		pagina = 1
	}
	limite, err := strconv.Atoi(limiteStr)
	if err != nil || limite < 1 {
		limite = 10
	}
	listagem := postagensService.GetPostagens(shared.DB, pagina, limite, pesquisa, "alunoPost", response)
	httpkit.AppSucess("Listagem bem sucedida", listagem, response)
}

func PostPostagem(response http.ResponseWriter, request *http.Request) {
	var postagem postagensDto.NovaPostagem
	json.NewDecoder(request.Body).Decode(&postagem)
	_, err := userService.GetUserById(shared.DB, postagem.UsuarioId, response)
	if err != nil {
		httpkit.AppBadRequest("Usuario com esse id não existe", response)
	}
	idPost := postagensService.InsertPost(shared.DB, postagem, response)
	post, _ := postagensService.GetPostById(shared.DB, idPost, response)
	httpkit.AppSucess("Sucesso", post, response)
}

func PostComentario(response http.ResponseWriter, request *http.Request) {
	var comentario postagensDto.ComentarioData
	json.NewDecoder(request.Body).Decode(&comentario)
	params, _ := httpkit.GetUrlParams(request)
	comentario.PostagemId, _ = strconv.Atoi(params.Param["postagemId"])
	postagensService.InsertComment(shared.DB, comentario, response)
	httpkit.AppSucess("Comentário inserido com sucesso", make(map[string]string), response)
}

func PostReacao(response http.ResponseWriter, request *http.Request) {
	var reacao postagensDto.ReacaoData
	json.NewDecoder(request.Body).Decode(&reacao)
	params, _ := httpkit.GetUrlParams(request)
	reacao.PostagemId, _ = strconv.Atoi(params.Param["postagemId"])
	postagensService.InsertReaction(shared.DB, reacao, response)
	httpkit.AppSucess("Reação inserida com sucesso", make(map[string]string), response)
}
