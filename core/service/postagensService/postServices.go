package postagensService

import (
	"api_journal/core/controller/postagensController/postagensDto"
	httpkit "api_journal/requester/http"
	"database/sql"
	"net/http"
	"time"

	"github.com/lib/pq"
)

type ReactionType string

const (
	FOGUINHO ReactionType = "FOGUINHO"
	LIKE     ReactionType = "LIKE"
	AMEI     ReactionType = "AMEI"
	ODIEI    ReactionType = "ODIEI"
)

func InsertPost(db *sql.DB, postagem postagensDto.NovaPostagem, response http.ResponseWriter) int {
	var lastInsertID int
	err := db.QueryRow("INSERT INTO postagem (tipo,titulo,conteudo,usuario_id,tags) VALUES ($1, $2, $3,$4,$5) RETURNING id",
		postagem.Tipo, postagem.Titulo, postagem.Conteudo, postagem.UsuarioId, pq.Array(postagem.Tags)).Scan(&lastInsertID)
	if err != nil {
		httpkit.AppConflict("Conflito ao tentar inserir nova postagem: "+err.Error(), response)
	}
	return int(lastInsertID)
}

func InsertPostFromUfca(db *sql.DB, tipo string, titulo string, conteudo string, tags []string, criadoEm time.Time) (int, error) {
	var lastInsertID int
	err := db.QueryRow("INSERT INTO postagem (tipo,titulo,conteudo,tags,criado_em) VALUES ($1, $2, $3,$4,$5) RETURNING id",
		tipo, titulo, conteudo, pq.Array(tags), criadoEm).Scan(&lastInsertID)
	if err != nil {
		return 0, err
	}
	return int(lastInsertID), nil
}

func InsertComment(db *sql.DB, comentario postagensDto.ComentarioData, response http.ResponseWriter) {
	var lastInsertID int
	err := db.QueryRow("INSERT INTO comentario (conteudo, postagem_id, usuario_id) VALUES ($1, $2, $3) RETURNING id", comentario.Conteudo, comentario.PostagemId, comentario.UsuarioId).Scan(&lastInsertID)
	if err != nil {
		httpkit.AppConflict("Erro ao tentar inserir o comentário: "+err.Error(), response)
	}
}

func InsertReaction(db *sql.DB, reacao postagensDto.ReacaoData, response http.ResponseWriter) {
	var lastInsertID int
	switch ReactionType(reacao.Tipo) {
	case LIKE, FOGUINHO, AMEI, ODIEI:
		break
	default:
		httpkit.AppBadRequest("Tipo de reação não é válido: "+reacao.Tipo, response)
		break
	}
	err := db.QueryRow("INSERT INTO reacao (tipo, postagem_id, usuario_id) VALUES ($1, $2, $3) RETURNING id", reacao.Tipo, reacao.PostagemId, reacao.UsuarioId).Scan(&lastInsertID)
	if err != nil {
		httpkit.AppConflict("Erro ao tentar inserir a reação: "+err.Error(), response)
	}
}
