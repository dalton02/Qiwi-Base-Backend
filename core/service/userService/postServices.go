package userService

import (
	"api_journal/core/controller/userController/userDto"
	httpkit "api_journal/requester/http"
	"database/sql"
	"net/http"
)

func InsertUsuario(db *sql.DB, user userDto.UserData, response http.ResponseWriter) int {
	var lastInsertID int

	err := db.QueryRow("INSERT INTO usuario (login,nome,curso,codigo) VALUES ($1, $2, $3,$4) RETURNING id", user.Login, user.Nome, user.Curso, user.Codigo).Scan(&lastInsertID)
	if err != nil {
		httpkit.AppConflict("Conflito ao tentar inserir novo usuário: "+err.Error(), response)
	}
	return int(lastInsertID)
}
