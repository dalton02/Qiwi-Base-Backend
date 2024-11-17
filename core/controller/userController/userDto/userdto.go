package userDto

type UserData struct {
	Id    int    `json:"id"`
	Login string `json:"login"`
	Nome  string `json:"nome"`
}

type UserSignin struct {
	Login string `json:"login" validator:"required"`
	Nome  string `json:"nome" validator:"required"`
	Senha string `json:"senha" validator:"strongPassword"`
	Email string `json:"email" validator:"email"`
}

type UserLogin struct {
	Login string `json:"login" validator:"required"`
	Senha string `json:"senha" validator:"required"`
}

type AlunoData struct {
	Login  string `json:"login"`
	Nome   string `json:"nome"`
	Curso  string `json:"curso"`
	Codigo int    `json:"codigo"`
	Id     int    `json:"id"`
}

type UserQuery struct {
	Pagina string `query:"pagina" validator:"optional,numericString"`
}
