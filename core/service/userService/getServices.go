package userService

import (
	"api_journal/core/controller/userController/userDto"
	httpkit "api_journal/requester/http"
	"database/sql"
	"net/http"
)

func GetUserById(db *sql.DB, id int, response http.ResponseWriter) (userDto.UserData, error) {
	var user userDto.UserData
	row := db.QueryRow("SELECT login,nome,codigo from usuario WHERE id=$1", id)
	err := row.Scan(&user.Login, &user.Nome, &user.Codigo)
	if err != nil {
		if err == sql.ErrNoRows {
			return user, err
		}
		httpkit.AppNotFound("Erro no banco ao tenta buscar usuario:"+err.Error(), response)
	}
	return user, nil
}

func GetUserByCodigo(db *sql.DB, id int, response http.ResponseWriter) (userDto.UserData, error) {
	var user userDto.UserData
	row := db.QueryRow("SELECT login,nome,codigo from usuario WHERE codigo=$1", id)
	err := row.Scan(&user.Login, &user.Nome, &user.Codigo)
	if err != nil {
		if err == sql.ErrNoRows {
			return user, err
		}
		httpkit.AppNotFound("Erro no banco ao tenta buscar usuario:"+err.Error(), response)
	}
	return user, nil
}
