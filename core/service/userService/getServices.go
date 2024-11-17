package userService

import (
	"api_journal/core/controller/userController/userDto"
	"database/sql"
	"fmt"
)

func GetUserByLoginPass(db *sql.DB, login string, senha string) (userDto.UserData, error) {
	var user userDto.UserData
	user.Login = login
	row := db.QueryRow("SELECT nome,id from usuario WHERE login=$1 AND senha=$2", login, senha)
	err := row.Scan(&user.Nome, &user.Id)
	if err != nil {
		if err == sql.ErrNoRows {
			return user, err
		}
		return user, err
	}

	return user, nil
}

func GetUserById(db *sql.DB, id int) (userDto.AlunoData, error) {
	var user userDto.AlunoData
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

func GetAlunoById(db *sql.DB, id int) (userDto.AlunoData, error) {
	var user userDto.AlunoData
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

func GetAlunoByCodigo(db *sql.DB, id int) (int, error) {
	var userId int
	row := db.QueryRow("SELECT id_usuario from aluno WHERE codigo=$1", id)
	fmt.Println("DASDSA")
	err := row.Scan(&userId)
	if err != nil {
		if err == sql.ErrNoRows {
			return userId, err
		}
		return userId, err
	}
	return userId, nil
}
