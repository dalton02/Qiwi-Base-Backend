package postagensService

import (
	"api_journal/core/controller/postagensController/postagensDto"
	"database/sql"
	"fmt"
	"time"

	"github.com/lib/pq"
)

func InsertPost(db *sql.DB, postagem postagensDto.NovaPostagem) (int, error) {
	var lastInsertID int
	err := db.QueryRow("INSERT INTO postagem (tipo,titulo,conteudo,usuario_id,tags) VALUES ($1, $2, $3,$4,$5) RETURNING id",
		postagem.Tipo, postagem.Titulo, postagem.Conteudo, postagem.UsuarioId, pq.Array(postagem.Tags)).Scan(&lastInsertID)
	if err != nil {
		return 0, fmt.Errorf("Erro ao inserir postagem: " + err.Error())
	}
	return int(lastInsertID), nil
}

func InsertPostFromUfca(db *sql.DB, tipo string, titulo string, conteudo string, tags []string, criadoEm time.Time) (int, error) {
	var lastInsertID int
	err := db.QueryRow("INSERT INTO postagem (tipo,titulo,conteudo,tags,criado_em) VALUES ($1, $2, $3,$4,$5) RETURNING id",
		tipo, titulo, conteudo, pq.Array(tags), criadoEm).Scan(&lastInsertID)
	if err != nil {
		return 0, fmt.Errorf("Erro ao inserir postagem: " + err.Error())
	}
	return int(lastInsertID), nil
}

func FilhoPaiPost(db *sql.DB, parenteId int, postagemId int) bool {
	var id int
	err := db.QueryRow("SELECT id FROM comentarios WHERE id=$1 AND postagem_id=$2", parenteId, postagemId).Scan(&id)

	if err != nil {
		return false
	}
	return true

}

func InsertComment(db *sql.DB, comentario postagensDto.ComentarioData, hasParent bool) (int, error) {
	var lastInsertID int
	var parenteId *int
	if hasParent {
		parenteId = &comentario.ParenteId
		valid := FilhoPaiPost(db, *parenteId, comentario.PostagemId)
		if !valid {
			return 0, fmt.Errorf("Comentario ID não existe")
		}
	} else {
		parenteId = nil
	}

	err := db.QueryRow("INSERT INTO comentario (conteudo, postagem_id, usuario_id, parente_id) VALUES ($1, $2, $3, $4) RETURNING id",
		comentario.Conteudo, comentario.PostagemId, comentario.UsuarioId, parenteId).Scan(&lastInsertID)

	if err != nil {
		return 0, err
	}

	return lastInsertID, nil
}

func InsertReaction(db *sql.DB, reacao postagensDto.ReacaoData) (int, error) {
	var lastInsertID int

	err := db.QueryRow("INSERT INTO reacao (tipo, postagem_id, usuario_id) VALUES ($1, $2, $3) RETURNING id", reacao.Tipo, reacao.PostagemId, reacao.UsuarioId).Scan(&lastInsertID)
	if err != nil {
		return 0, fmt.Errorf("Erro ao inserir nova reação")
	}
	return int(lastInsertID), nil
}
