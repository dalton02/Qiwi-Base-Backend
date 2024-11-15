package postagensService

import (
	"api_journal/core/controller/postagensController/postagensDto"
	"database/sql"

	"github.com/lib/pq"
)

func GetPostagemExiste(db *sql.DB, id int) bool {
	var titulo string
	row := db.QueryRow("SELECT titulo from postagem WHERE id=$1", id)
	err := row.Scan(&titulo)
	if err != nil {
		return false
	}
	return true
}
func GetPostById(db *sql.DB, id int) (postagensDto.NovaPostagem, error) {
	var postagem postagensDto.NovaPostagem
	row := db.QueryRow("SELECT tipo,titulo,conteudo,tags,usuario_id from postagem WHERE id=$1", id)
	err := row.Scan(&postagem.Tipo, &postagem.Titulo, &postagem.Conteudo, pq.Array(&postagem.Tags), &postagem.UsuarioId)
	if err != nil {
		if err == sql.ErrNoRows {
			return postagem, err
		}
		return postagem, err
	}
	return postagem, nil
}

func GetPostagemByTitle(db *sql.DB, titulo string, tipo string) (postagensDto.PostagemDataComplete, error, int) {
	var postagem postagensDto.PostagemDataComplete
	var query string
	if tipo == "ufca-reportagem" {
		query = queryPostagemUfca
	}
	query = queryPostagemNormal
	rows, err := db.Query(query, titulo, tipo)
	if err != nil {
		return postagem, err, 400
	}
	defer rows.Close()

	postagem.Tags = []string{}
	postagem.Reacoes = make(map[string]int)
	postagem.Comentarios = []postagensDto.ComentarioDataComplete{}

	if tipo == "ufca-reportagem" {
		return scanPostagemUfca(rows, postagem)
	}
	return scanPostagemNormal(rows, postagem)

}

func GetPostagens(db *sql.DB, pagina int, limite int, pesquisa string, tipo string) (postagensDto.ListagemPostagens, error) {
	var listagem postagensDto.ListagemPostagens
	var total int
	err := db.QueryRow("SELECT COUNT(*) FROM postagem WHERE titulo ILIKE $1", "%"+pesquisa+"%").Scan(&total)
	if err != nil {
		return listagem, err
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
		return listagem, err
	}
	defer rows.Close()

	for rows.Next() {
		var postagem postagensDto.PostagemDataLista
		var tags pq.StringArray
		var reacoesCount int
		var reacaoTipo sql.NullString // Para lidar com reações que podem ser nulas
		rows.Scan(&postagem.Id, &postagem.Titulo, &postagem.Tipo, &postagem.Conteudo, &tags,
			&postagem.Comentarios, &reacoesCount, &reacaoTipo, &postagem.Autor.Nome, &postagem.Autor.Curso, &postagem.Autor.Login, &postagem.Autor.Id)

		postagem.Tags = tags
		postagem.Reacoes = make(map[string]int)

		if reacaoTipo.Valid {
			postagem.Reacoes[reacaoTipo.String] = reacoesCount
		}

		listagem.Postagem = append(listagem.Postagem, postagem)
	}
	listagem.Pagina = pagina
	listagem.Limite = limite
	listagem.Pesquisa = pesquisa
	return listagem, nil
}
