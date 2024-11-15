package userService

import (
	"api_journal/core/controller/userController/userDto"
	httpkit "api_journal/requester/http"
	"database/sql"
	"net/http"
)

func InsertUsuario(db *sql.DB, user userDto.UserData, perfil string, response http.ResponseWriter) (int, error) {
	var lastInsertID int

	err := db.QueryRow("INSERT INTO usuario (login,nome,perfil) VALUES ($1, $2, $3) RETURNING id", user.Login, user.Nome, perfil).Scan(&lastInsertID)
	if err != nil {
		httpkit.AppConflict("Conflito ao tentar inserir novo usuário: "+err.Error(), response)
		return 0, err
	}
	return int(lastInsertID), nil
}

func InsertAlunoAndRelate(db *sql.DB, user userDto.UserData, idUsuario int, response http.ResponseWriter) (int, error) {
	var lastInsertID int

	err := db.QueryRow("INSERT INTO aluno (curso,codigo,id_usuario) VALUES ($1, $2, $3) RETURNING id", user.Curso, user.Codigo, idUsuario).Scan(&lastInsertID)
	if err != nil {

		httpkit.AppConflict("Conflito ao tentar inserir novo aluno: "+err.Error(), response)
		return 0, err
	}
	db.QueryRow("UPDATE usuario SET id_aluno=$1 WHERE id=$2", lastInsertID, idUsuario)
	return int(lastInsertID), nil
}
