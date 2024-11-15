package userService

import (
	"api_journal/core/controller/userController/userDto"
	"database/sql"
	"fmt"
)

func GetUserById(db *sql.DB, id int) (userDto.UserData, error) {
	var user userDto.UserData
	row := db.QueryRow("SELECT nome from usuario WHERE id=$1", id)
	err := row.Scan(&user.Nome)
	if err != nil {
		if err == sql.ErrNoRows {
			return user, err
		}
		return user, err
	}
	return user, nil
}

func GetAlunoById(db *sql.DB, id int) (userDto.UserData, error) {
	var user userDto.UserData
	row := db.QueryRow("SELECT curso from aluno WHERE usuario_id =  ", id)
	err := row.Scan(&user.Login)
	if err != nil {
		fmt.Println(err)
		if err == sql.ErrNoRows {
			return user, err
		}
		return user, err
	}
	return user, nil
}

func GetAlunoByCodigo(db *sql.DB, id int) (userDto.UserData, error) {
	var user userDto.UserData
	row := db.QueryRow("SELECT curso from aluno WHERE codigo=$1", id)
	err := row.Scan(&user.Login)
	if err != nil {
		fmt.Println(err)
		if err == sql.ErrNoRows {
			return user, err
		}
		return user, err
	}
	return user, nil
}
