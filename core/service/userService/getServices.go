package userService

import (
	"api_journal/core/controller/userController/userDto"
	"database/sql"
)

func GetUserById(db *sql.DB, id int) (userDto.UserData, error) {
	var user userDto.UserData
	row := db.QueryRow("SELECT login,nome,codigo from usuario WHERE id=$1", id)
	err := row.Scan(&user.Login, &user.Nome, &user.Codigo)
	if err != nil {
		if err == sql.ErrNoRows {
			return user, err
		}
		return user, err
	}
	return user, nil
}

func GetUserByCodigo(db *sql.DB, id int) (userDto.UserData, error) {
	var user userDto.UserData
	row := db.QueryRow("SELECT login,nome,codigo from usuario WHERE codigo=$1", id)
	err := row.Scan(&user.Login, &user.Nome, &user.Codigo)
	if err != nil {
		if err == sql.ErrNoRows {
			return user, err
		}
		return user, err
	}
	return user, nil
}
