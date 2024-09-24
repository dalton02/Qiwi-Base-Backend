package userService

import (
	"api_journal/core/controller/userController/userDto"
	httpkit "api_journal/requester/http"
	"database/sql"
	"net/http"
)

func CheckUsuario(db *sql.DB, user userDto.UserLogin, response http.ResponseWriter) userDto.UserData {
	var Login string
	var Nome string
	var senha string
	var id int
	var userData userDto.UserData
	row := db.QueryRow("SELECT login,nome,senha,id from Usuarios WHERE login=$1", user.Login)
	err := row.Scan(&Login, &Nome, &senha, &id)
	if err != nil {
		if err == sql.ErrNoRows {
			httpkit.AppNotFound("Usuário não encontrado", response)
		}
		httpkit.AppUnauthorized("Credenciais incorretas", response)

	}
	if senha != user.Senha {
		httpkit.AppUnauthorized("Credenciais incorretas", response)
	}
	userData = userDto.UserData{
		Id:    id,
		Nome:  Nome,
		Login: Login,
	}
	return userData
}

func InsertUsuario(db *sql.DB, user userDto.UserSignin, response http.ResponseWriter) int {
	var lastInsertID int

	err := db.QueryRow("INSERT INTO Usuarios (login,nome,senha) VALUES ($1, $2, $3) RETURNING id", user.Login, user.Nome, user.Senha).Scan(&lastInsertID)
	if err != nil {
		httpkit.AppConflict("Conflito ao tentar inserir novo usuário", response)
	}
	return int(lastInsertID)
}
