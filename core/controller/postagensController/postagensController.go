package postagensController

import (
	"api_journal/core/controller/postagensController/postagensDto"
	"api_journal/core/server/shared"
	"api_journal/core/service/postagensService"
	httpkit "api_journal/requester/http"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
)

type ReactionType string

const (
	FOGUINHO ReactionType = "FOGUINHO"
	LIKE     ReactionType = "LIKE"
	AMEI     ReactionType = "AMEI"
	ODIEI    ReactionType = "ODIEI"
)

// @Summary Pegar postagem por titulo
// @Description Pega uma postagem por titulo
// @Tags Postagens
// @Param titulo path string true "Título da postagem"
// @Param tipo query string  true "Tipo de postagem: (alunoPost,ufcaPost)"
// @Produce json
// @Success 200 {object} postagensDto.PostagemDataComplete "Postagem a mostra"
// @Failure 400 {object} map[string]string "Erro ao processar a requisição"
// @Failure 500 {object} map[string]string "Erro interno do servidor"
// @Router /postagens/{titulo} [get]
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

// @Summary Pegar listagem de postagens
// @Description Gera uma lista paginada com as postagens do blog
// @Tags Postagens
// @Param pagina query string false "Pagina"
// @Param limite query string false "Limite por página"
// @Param pesquisa query string false "Pesquisa"
// @Produce json
// @Success 200 {object} postagensDto.ListagemPostagens  "Listagem a mostra:"
// @Failure 400 {object} map[string]string "Erro ao processar a requisição"
// @Failure 500 {object} map[string]string "Erro interno do servidor"
// @Router /postagens/listar [get]
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

// @Summary Criar uma nova postagem
// @Description Cria uma nova postagem no blog
// @Tags Postagens
// @Accept json
// @Param Authorization header string true "Bearer token"
// @Param postagem body postagensDto.NovaPostagem true "Dados da nova postagem"
// @Produce json
// @Success 201 {object} postagensDto.PostagemDataComplete "Postagem criada com sucesso"
// @Failure 400 {object} map[string]string "Erro ao processar a requisição"
// @Failure 500 {object} map[string]string "Erro interno do servidor"
// @Router /postagens/postar [post]
func PostPostagem(response http.ResponseWriter, request *http.Request) {
	var postagem postagensDto.NovaPostagem
	json.NewDecoder(request.Body).Decode(&postagem)

	dataToken, _ := httpkit.GetDataToken(request)
	idToken, ok := dataToken["id"].(float64)
	if ok {
		postagem.UsuarioId = int(idToken)
	}

	idPost, err := postagensService.InsertPost(shared.DB, postagem)
	if err != nil {
		if strings.Contains(err.Error(), "titulo_key") {
			httpkit.AppBadRequest("Já existe uma postagem com mesmo titulo", response)
			return
		}
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

// @Summary Adicionar um comentário a uma postagem
// @Description Adiciona um comentário a uma postagem existente
// @Tags Comentários
// @Accept json
// @Param Authorization header string true "Bearer token"
// @Param postagemId path int true "ID da postagem"
// @Param comentario body postagensDto.ComentarioData true "Dados do comentário"
// @Produce json
// @Success 201 {object} map[string]string "Comentário inserido com sucesso"
// @Failure 400 {object} map[string]string "Erro ao processar a requisição"
// @Failure 500 {object} map[string]string "Erro interno do servidor"
// @Router /postagens/{postagemId}/comentar [post]
func PostComentario(response http.ResponseWriter, request *http.Request) {
	var comentario postagensDto.ComentarioData
	json.NewDecoder(request.Body).Decode(&comentario)
	params, _ := httpkit.GetUrlParams(request)
	comentario.PostagemId, _ = strconv.Atoi(params.Param["postagemId"])

	dataToken, _ := httpkit.GetDataToken(request)
	idToken, ok := dataToken["id"].(float64)
	if ok {
		comentario.UsuarioId = int(idToken)
	}

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

// @Summary Adicionar ou atualizar uma reação a uma postagem
// @Description Adiciona ou atualiza uma reação (LIKE, FOGUINHO, AMEI, ODIEI) para uma postagem específica
// @Tags Reações
// @Accept json
// @Param postagemId path int true "ID da postagem"
// @Param Authorization header string true "Bearer token"
// @Param reacao body postagensDocDto.ReacaoData true "Dados da reação"
// @Produce json
// @Success 201 {object} map[string]string "Reação inserida com sucesso"
// @Success 200 {object} map[string]string "Reação atualizada com sucesso"
// @Failure 400 {object} map[string]string "Erro ao processar a requisição"
// @Failure 500 {object} map[string]string "Erro interno do servidor"
// @Router /postagens/{postagemId}/reagir [post]
func PostReacao(response http.ResponseWriter, request *http.Request) {
	var reacao postagensDto.ReacaoData
	json.NewDecoder(request.Body).Decode(&reacao)

	params, _ := httpkit.GetUrlParams(request)
	reacao.PostagemId, _ = strconv.Atoi(params.Param["postagemId"])

	dataToken, _ := httpkit.GetDataToken(request)
	idToken, ok := dataToken["id"].(float64)
	if ok {
		reacao.UsuarioId = int(idToken)
	}

	switch ReactionType(reacao.Tipo) {
	case LIKE, FOGUINHO, AMEI, ODIEI:
		break
	default:
		httpkit.AppBadRequest("Tipo de reação não é válido - Opções disponiveis: (LIKE,FOGUINHO,AMEI,ODIEI)", response)
		return
	}

	_, err := postagensService.GetReaction(shared.DB, reacao)

	if err == nil {
		_, errUpdate := postagensService.UpdateReaction(shared.DB, reacao)

		if errUpdate != nil {
			httpkit.AppBadRequest(err.Error(), response)
			return
		}
		httpkit.AppSucess("Reação atualizada com sucesso", make(map[string]string), response)
		return
	}

	_, err = postagensService.InsertReaction(shared.DB, reacao)
	if err != nil {
		httpkit.AppBadRequest(err.Error(), response)
		return
	}
	httpkit.AppSucessCreate("Reação inserida com sucesso", make(map[string]string), response)
}

func PostByParamExiste(response http.ResponseWriter, request *http.Request) bool {

	params, _ := httpkit.GetUrlParams(request)
	postId, _ := strconv.Atoi(params.Param["postagemId"])
	existe := postagensService.GetPostagemExiste(shared.DB, postId)
	if !existe {
		httpkit.AppNotFound("Postagem não encontrada", response)
		return false
	}
	return true
}
