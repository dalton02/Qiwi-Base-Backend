package postagensService

import (
	"api_journal/core/controller/postagensController/postagensDto"
	"database/sql"
	"fmt"

	"github.com/lib/pq"
)

const queryPostagemNormal = `SELECT 
       p.id, p.titulo, p.tipo, p.conteudo,p.tags,
       comentario.conteudo,comentario.criado_em,comentario.id,
       usuarioComentario.nome,usuarioComentario.curso,usuarioComentario.login, usuarioComentario.id,
       reacao.tipo,reacao.reacoes_count,
       usuarioPost.nome,usuarioPost.curso,usuarioPost.login,usuarioPost.id,
       filhos.conteudo, filhos.criado_em,filhos.id,
       usuarioFilho.nome,usuarioFilho.curso,usuarioFilho.login,usuarioFilho.id
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
       LEFT JOIN comentario ON p.id = comentario.postagem_id AND comentario.parente_id IS NULL
       LEFT JOIN usuario usuarioComentario ON comentario.usuario_id = usuarioComentario.id
       LEFT JOIN comentario filhos ON comentario.id = filhos.parente_id
       LEFT JOIN usuario usuarioFilho ON filhos.usuario_id = usuarioFilho.id
	   WHERE p.titulo = $1
       AND p.tipo = $2;`

const queryPostagemUfca = `SELECT 
p.id, p.titulo, p.tipo, p.conteudo,p.tags,
comentario.conteudo,comentario.criado_em,comentario.id,
usuarioComentario.nome,usuarioComentario.curso,usuarioComentario.login, usuarioComentario.id,
reacao.tipo,reacao.reacoes_count,
filhos.conteudo, filhos.criado_em,filhos.id,
usuarioFilho.nome,usuarioFilho.curso,usuarioFilho.login,usuarioFilho.id
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
LEFT JOIN comentario ON p.id = comentario.postagem_id AND comentario.parente_id IS NULL
LEFT JOIN usuario usuarioComentario ON comentario.usuario_id = usuarioComentario.id
LEFT JOIN comentario filhos ON comentario.id = filhos.parente_id
LEFT JOIN usuario usuarioFilho ON filhos.usuario_id = usuarioFilho.id
WHERE p.titulo = $1
AND p.tipo = $2;`

func scanPostagemNormal(rows *sql.Rows, postagem postagensDto.PostagemDataComplete) (postagensDto.PostagemDataComplete, error, int) {
	found := false
	for rows.Next() {
		found = true
		var tags pq.StringArray
		var comentario postagensDto.ComentarioDataComplete
		var reacoesCount sql.NullInt64
		var reacaoTipo sql.NullString
		var ComConteudo sql.NullString
		var ComCriadoEm sql.NullString
		var ComId sql.NullInt32
		var ComAutorNome sql.NullString
		var ComAutorCurso sql.NullString
		var ComAutorLogin sql.NullString
		var ComAutorId sql.NullInt32

		var filhoComConteudo sql.NullString
		var filhoComCriadoEm sql.NullString
		var filhoComId sql.NullInt32
		var filhoComAutorNome sql.NullString
		var filhoComAutorCurso sql.NullString
		var filhoComAutorLogin sql.NullString
		var filhoComAutorId sql.NullInt32

		err := rows.Scan(&postagem.Id, &postagem.Titulo, &postagem.Tipo, &postagem.Conteudo, &tags,
			&ComConteudo, &ComCriadoEm, &ComId, &ComAutorNome, &ComAutorCurso, &ComAutorLogin, &ComAutorId,
			&reacaoTipo, &reacoesCount,
			&postagem.Autor.Nome, &postagem.Autor.Curso, &postagem.Autor.Login, &postagem.Autor.Id,
			&filhoComConteudo, &filhoComCriadoEm, &filhoComId, &filhoComAutorNome, &filhoComAutorCurso, &filhoComAutorLogin, &filhoComAutorId)

		if err != nil {
			return postagem, err, 400
		}

		postagem.Tags = tags

		if reacaoTipo.Valid {
			postagem.Reacoes[reacaoTipo.String] = int(reacoesCount.Int64)
		}
		var filho postagensDto.ComentarioDataComplete

		comentario.Conteudo = ComConteudo.String
		comentario.CriadoEm = ComCriadoEm.String
		comentario.Id = int(ComId.Int32)
		comentario.Autor.Curso = ComAutorCurso.String
		comentario.Autor.Login = ComAutorLogin.String
		comentario.Autor.Nome = ComAutorNome.String
		comentario.Autor.Id = int(ComAutorId.Int32)

		if filhoComAutorId.Valid {
			filho = postagensDto.ComentarioDataComplete{
				Conteudo: filhoComConteudo.String,
				Id:       int(filhoComId.Int32),
				CriadoEm: filhoComCriadoEm.String,
				Autor: postagensDto.UserPostagem{
					Login: filhoComAutorLogin.String,
					Id:    int(filhoComAutorId.Int32),
					Nome:  filhoComAutorNome.String,
					Curso: filhoComAutorCurso.String,
				},
			}
			comentario.Filhos = append(comentario.Filhos, filho)
		}
		if ComId.Valid {
			postagem.Comentarios = append(postagem.Comentarios, comentario)
		}
	}
	if !found {
		return postagem, fmt.Errorf("Nenhuma postagem encontrada"), 404
	}
	return postagem, nil, 200
}

func scanPostagemUfca(rows *sql.Rows, postagem postagensDto.PostagemDataComplete) (postagensDto.PostagemDataComplete, error, int) {
	found := false
	fmt.Println(rows.Columns())
	for rows.Next() {
		found = true
		var tags pq.StringArray
		var comentario postagensDto.ComentarioDataComplete
		var reacoesCount sql.NullInt64
		var reacaoTipo sql.NullString
		var ComConteudo sql.NullString
		var ComCriadoEm sql.NullString
		var ComId sql.NullInt32
		var ComAutorNome sql.NullString
		var ComAutorCurso sql.NullString
		var ComAutorLogin sql.NullString
		var ComAutorId sql.NullInt32

		var filhoComConteudo sql.NullString
		var filhoComCriadoEm sql.NullString
		var filhoComId sql.NullInt32
		var filhoComAutorNome sql.NullString
		var filhoComAutorCurso sql.NullString
		var filhoComAutorLogin sql.NullString
		var filhoComAutorId sql.NullInt32
		var tmp sql.NullString
		err := rows.Scan(&postagem.Id, &postagem.Titulo, &postagem.Tipo, &postagem.Conteudo, &tags,
			&ComConteudo, &ComCriadoEm, &ComId, &ComAutorNome, &ComAutorCurso, &ComAutorLogin, &ComAutorId,
			&reacaoTipo, &reacoesCount,
			&tmp, &tmp, &tmp, &tmp,
			&filhoComConteudo, &filhoComCriadoEm, &filhoComId, &filhoComAutorNome, &filhoComAutorCurso, &filhoComAutorLogin, &filhoComAutorId)

		if err != nil {
			return postagem, err, 400
		}
		postagem.Tags = tags

		if reacaoTipo.Valid {
			postagem.Reacoes[reacaoTipo.String] = int(reacoesCount.Int64)
		}
		var filho postagensDto.ComentarioDataComplete
		if filhoComAutorId.Valid {
			filho = postagensDto.ComentarioDataComplete{
				Conteudo: filhoComConteudo.String,
				Id:       int(filhoComId.Int32),
				CriadoEm: filhoComCriadoEm.String,
				Autor: postagensDto.UserPostagem{
					Login: filhoComAutorLogin.String,
					Id:    int(filhoComAutorId.Int32),
					Nome:  filhoComAutorNome.String,
					Curso: filhoComAutorCurso.String,
				},
			}
			comentario.Filhos = append(comentario.Filhos, filho)
		}
		if ComId.Valid {
			postagem.Comentarios = append(postagem.Comentarios, comentario)
		}
	}
	if !found {
		return postagem, fmt.Errorf("Nenhuma postagem encontrada"), 404
	}
	return postagem, nil, 200
}
