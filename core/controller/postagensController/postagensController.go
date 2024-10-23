package postagensController

import (
	"api_journal/core/controller/postagensController/postagensDto"
	"api_journal/core/server/shared"
	"api_journal/core/service/postagensService"
	"api_journal/core/service/userService"
	httpkit "api_journal/requester/http"
	"encoding/json"
	"net/http"
	"strconv"
)

type ReactionType string

const (
	FOGUINHO ReactionType = "FOGUINHO"
	LIKE     ReactionType = "LIKE"
	AMEI     ReactionType = "AMEI"
	ODIEI    ReactionType = "ODIEI"
)

func GetPostagemByTitle(response http.ResponseWriter, request *http.Request) {
	params, _ := httpkit.GetUrlParams(request)
	titulo := params.Param["titulo"]
	tipo := request.URL.Query().Get("tipo")
	listagem, err, status := postagensService.GetPostagemByTitle(shared.DB, titulo, tipo)
	if err != nil {
		httpkit.GenerateErrorHttpMessage(status, err.Error(), response)
		return
	}
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
	listagem, err := postagensService.GetPostagens(shared.DB, pagina, limite, pesquisa, "alunoPost")
	if err != nil {
		httpkit.AppBadRequest(err.Error(), response)
		return
	}
	httpkit.AppSucess("Listagem bem sucedida", listagem, response)
}

func PostPostagem(response http.ResponseWriter, request *http.Request) {
	var postagem postagensDto.NovaPostagem
	json.NewDecoder(request.Body).Decode(&postagem)
	_, err := userService.GetUserById(shared.DB, postagem.UsuarioId)
	if err != nil {
		httpkit.AppBadRequest("Usuario com esse id não existe", response)
		return
	}
	idPost, err := postagensService.InsertPost(shared.DB, postagem)
	if err != nil {
		httpkit.AppBadRequest(err.Error(), response)
		return
	}
	post, err := postagensService.GetPostById(shared.DB, idPost)
	if err != nil {
		httpkit.AppBadRequest(err.Error(), response)
		return
	}
	httpkit.AppSucess("Sucesso", post, response)
}

func PostComentario(response http.ResponseWriter, request *http.Request) {
	var comentario postagensDto.ComentarioData
	json.NewDecoder(request.Body).Decode(&comentario)
	params, _ := httpkit.GetUrlParams(request)
	comentario.PostagemId, _ = strconv.Atoi(params.Param["postagemId"])
	jsonSchema := httpkit.GetJsonSchema[postagensDto.ComentarioData](request)
	hasParent := false
	for i := 0; i < len(jsonSchema[1]); i++ {
		if jsonSchema[1][i] == "parenteId" {
			hasParent = true
		}
	}
	_, err := postagensService.InsertComment(shared.DB, comentario, hasParent)
	if err != nil {
		httpkit.AppBadRequest(err.Error(), response)
		return
	}
	httpkit.AppSucess("Comentário inserido com sucesso", make(map[string]string), response)
}

func PostReacao(response http.ResponseWriter, request *http.Request) {
	var reacao postagensDto.ReacaoData
	json.NewDecoder(request.Body).Decode(&reacao)
	params, _ := httpkit.GetUrlParams(request)
	reacao.PostagemId, _ = strconv.Atoi(params.Param["postagemId"])
	switch ReactionType(reacao.Tipo) {
	case LIKE, FOGUINHO, AMEI, ODIEI:
		break
	default:
		httpkit.AppBadRequest("Tipo de reação não é válido: "+reacao.Tipo, response)
		return
	}
	_, err := postagensService.InsertReaction(shared.DB, reacao)
	if err != nil {
		httpkit.AppBadRequest(err.Error(), response)
		return
	}
	httpkit.AppSucess("Reação inserida com sucesso", make(map[string]string), response)
}
