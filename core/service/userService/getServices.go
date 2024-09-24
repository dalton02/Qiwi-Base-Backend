package userService

import (
	"api_journal/core/controller/userController/userDto"
	httpkit "api_journal/requester/http"
	"database/sql"
	"net/http"
)

func GetUserById(db *sql.DB, id int, response http.ResponseWriter) userDto.UserData {
	var user userDto.UserData
	row := db.QueryRow("SELECT login,nome,id from usuarios WHERE id=$1", id)
	err := row.Scan(&user.Login, &user.Nome, &user.Id)
	if err != nil {
		if err == sql.ErrNoRows {
			httpkit.AppNotFound("Usuário não encontrado", response)
		}
		httpkit.AppNotFound("Usuário não encontrado: "+err.Error(), response)
	}
	return user
}
