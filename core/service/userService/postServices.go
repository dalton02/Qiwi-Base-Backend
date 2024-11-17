package userService

import (
	"api_journal/core/controller/userController/userDto"
	httpkit "api_journal/requester/http"
	"database/sql"
	"net/http"
	"strings"
)

func InsertUsuarioExterno(db *sql.DB, user userDto.UserSignin, response http.ResponseWriter) (int, error) {
	var lastInsertID int

	err := db.QueryRow("INSERT INTO usuario (login,nome,senha,email,perfil) VALUES ($1, $2, $3,$4,$5) RETURNING id", user.Login, user.Nome, user.Senha, user.Email, "externo").Scan(&lastInsertID)
	if err != nil {
		if strings.Contains(err.Error(), "usuario_email") {
			httpkit.AppConflict("Conflito ao tentar inserir novo usuário: Email já cadastrado", response)
		}
		if strings.Contains(err.Error(), "usuario_login") {
			httpkit.AppConflict("Conflito ao tentar inserir novo usuário: Login já é cadastrado no sistema", response)
		}
		return 0, err
	}
	return int(lastInsertID), nil
}
func InsertUsuarioAluno(db *sql.DB, user userDto.UserSignin, response http.ResponseWriter) (int, error) {
	var lastInsertID int

	err := db.QueryRow("INSERT INTO usuario (login,nome,perfil) VALUES ($1, $2, $3,$4) RETURNING id", user.Login, user.Nome, "aluno").Scan(&lastInsertID)
	if err != nil {
		httpkit.AppConflict("Conflito ao tentar inserir novo usuário: "+err.Error(), response)
		return 0, err
	}
	return int(lastInsertID), nil
}

func InsertAlunoAndRelate(db *sql.DB, user userDto.AlunoData, idUsuario int, response http.ResponseWriter) (int, error) {
	var lastInsertID int

	err := db.QueryRow("INSERT INTO aluno (curso,codigo,id_usuario) VALUES ($1, $2, $3) RETURNING id", user.Curso, user.Codigo, idUsuario).Scan(&lastInsertID)
	if err != nil {

		httpkit.AppConflict("Conflito ao tentar inserir novo aluno: "+err.Error(), response)
		return 0, err
	}
	db.QueryRow("UPDATE usuario SET id_aluno=$1 WHERE id=$2", lastInsertID, idUsuario)
	return int(lastInsertID), nil
}
