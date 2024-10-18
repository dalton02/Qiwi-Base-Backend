package postagensService

import (
	"api_journal/core/controller/postagensController/postagensDto"
	httpkit "api_journal/requester/http"
	"database/sql"
	"fmt"
	"net/http"

	"github.com/lib/pq"
)

func GetPostById(db *sql.DB, id int, response http.ResponseWriter) (postagensDto.NovaPostagem, error) {
	var postagem postagensDto.NovaPostagem
	row := db.QueryRow("SELECT tipo,titulo,conteudo,tags,usuario_id from postagem WHERE id=$1", id)
	err := row.Scan(&postagem.Tipo, &postagem.Titulo, &postagem.Conteudo, pq.Array(&postagem.Tags), &postagem.UsuarioId)
	if err != nil {
		if err == sql.ErrNoRows {
			return postagem, err
		}
		httpkit.AppNotFound("Erro no banco ao tenta buscar postagem:"+err.Error(), response)
	}
	return postagem, nil
}
func GetPostagemByTitle(db *sql.DB, titulo string, tipo string, response http.ResponseWriter) postagensDto.PostagemDataComplete {
	var postagem postagensDto.PostagemDataComplete

	query := `SELECT p.id, p.titulo, p.tipo, p.conteudo, p.tags,
       comentario.conteudo, usuarioComentario.nome, usuarioComentario.curso, usuarioComentario.login, comentario.criado_em,
       reacao.tipo, reacao.reacoes_count,
       usuarioPost.nome, usuarioPost.curso, usuarioPost.login, usuarioPost.id
	   FROM postagem p
	   LEFT JOIN (
           SELECT postagem_id, COUNT(*) AS comentarios_count
           FROM comentario
           GROUP BY postagem_id
       ) comentario_count ON p.id = comentario_count.postagem_id
       LEFT JOIN (
           SELECT postagem_id, tipo, COUNT(*) AS reacoes_count
           FROM reacao
           GROUP BY postagem_id, tipo
       ) reacao ON p.id = reacao.postagem_id
       LEFT JOIN usuario usuarioPost ON p.usuario_id = usuarioPost.id
       LEFT JOIN (
           SELECT conteudo, usuario_id, criado_em, postagem_id
           FROM comentario
       ) comentario ON p.id = comentario.postagem_id
       LEFT JOIN usuario usuarioComentario ON comentario.usuario_id = usuarioComentario.id
	   WHERE p.titulo = $1
       AND p.tipo = $2;`

	rows, err := db.Query(query, titulo, tipo)
	if err != nil {
		fmt.Println(err)
		httpkit.AppBadRequest("Erro ao buscar postagens", response)
		return postagem
	}
	defer rows.Close()

	if rows.Next() {
		var tags pq.StringArray
		var comentarioConteudo sql.NullString
		var comentarioUsuarioNome sql.NullString
		var comentarioUsuarioCurso sql.NullString
		var comentarioUsuarioLogin sql.NullString
		var comentarioCriadoEm sql.NullString
		var reacoesCount sql.NullInt64
		var reacaoTipo sql.NullString

		err := rows.Scan(&postagem.Id, &postagem.Titulo, &postagem.Tipo, &postagem.Conteudo, &tags,
			&comentarioConteudo, &comentarioUsuarioNome, &comentarioUsuarioCurso, &comentarioUsuarioLogin, &comentarioCriadoEm,
			&reacaoTipo, &reacoesCount,
			&postagem.Autor.Nome, &postagem.Autor.Curso, &postagem.Autor.Login, &postagem.Autor.Id)

		if err != nil {
			httpkit.AppBadRequest("Erro ao ler dados da postagem", response)
			return postagem
		}

		postagem.Tags = tags
		postagem.Reacoes = make(map[string]int)
		postagem.Comentarios = []postagensDto.ComentarioDataComplete{}

		if reacaoTipo.Valid {
			postagem.Reacoes[reacaoTipo.String] = int(reacoesCount.Int64)
		}

		if comentarioConteudo.Valid {
			comentario := postagensDto.ComentarioDataComplete{
				Conteudo: comentarioConteudo.String,
				CriadoEm: comentarioCriadoEm.String,
				Autor: postagensDto.UserPostagem{
					Login: comentarioUsuarioLogin.String,
					Id:    1,
					Nome:  comentarioUsuarioNome.String,
					Curso: comentarioUsuarioCurso.String,
				},
			}
			postagem.Comentarios = append(postagem.Comentarios, comentario)
		}
	} else {
		httpkit.AppNotFound("Postagem não encontrada", response)
	}

	return postagem
}

func GetPostagens(db *sql.DB, pagina int, limite int, pesquisa string, tipo string, response http.ResponseWriter) postagensDto.ListagemPostagens {
	var listagem postagensDto.ListagemPostagens
	var total int
	err := db.QueryRow("SELECT COUNT(*) FROM postagem WHERE titulo ILIKE $1", "%"+pesquisa+"%").Scan(&total)
	if err != nil {
		httpkit.AppBadRequest("Erro ao contar postagens:", response)
	}
	listagem.TotalPaginas = (total + limite - 1) / limite

	query := `SELECT p.id, p.titulo, p.tipo, p.conteudo, p.tags,
       COALESCE(c.comentarios_count, 0) AS comentarios_count,
       COALESCE(r.reacoes_count, 0) AS reacoes_count,
       r.tipo AS reacao_tipo,
       u.nome AS usuario_nome, u.curso AS usuario_curso, u.login AS usuario_login,u.id AS usuario_id
FROM postagem p
LEFT JOIN (
    SELECT postagem_id, COUNT(*) AS comentarios_count
    FROM comentario
    GROUP BY postagem_id
) c ON p.id = c.postagem_id
LEFT JOIN (
    SELECT postagem_id, tipo, COUNT(*) AS reacoes_count
    FROM reacao
    GROUP BY postagem_id, tipo
) r ON p.id = r.postagem_id
LEFT JOIN usuario u ON p.usuario_id = u.id
WHERE p.titulo ILIKE $1
  AND p.tipo = $2
ORDER BY p.id DESC
LIMIT $3 OFFSET $4;`

	rows, err := db.Query(query, "%"+pesquisa+"%", tipo, limite, (pagina-1)*limite)
	if err != nil {
		httpkit.AppBadRequest("Erro ao buscar postagens", response)
	}
	defer rows.Close()

	// Laço para ler as postagens
	for rows.Next() {
		var postagem postagensDto.PostagemDataLista
		var tags pq.StringArray
		var reacoesCount int
		var reacaoTipo sql.NullString // Para lidar com reações que podem ser nulas
		rows.Scan(&postagem.Id, &postagem.Titulo, &postagem.Tipo, &postagem.Conteudo, &tags,
			&postagem.Comentarios, &reacoesCount, &reacaoTipo, &postagem.Autor.Nome, &postagem.Autor.Curso, &postagem.Autor.Login, &postagem.Autor.Id)

		postagem.Tags = tags
		postagem.Reacoes = make(map[string]int)

		// Adicionando a contagem de reações
		if reacaoTipo.Valid {
			postagem.Reacoes[reacaoTipo.String] = reacoesCount
		}

		// Adicionar a postagem à lista
		listagem.Postagem = append(listagem.Postagem, postagem)
	}
	return listagem
}
