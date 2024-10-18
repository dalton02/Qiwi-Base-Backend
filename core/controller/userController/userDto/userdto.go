package userDto

type UserLogin struct {
	Login string `json:"login" validator:"required"`
	Senha string `json:"senha" validator:"required"`
}

type UserQuery struct {
	Pagina string `query:"pagina" validator:"optional,numericString"`
}
type UserSignin struct {
	Login string `json:"login" validator:"optional"`
	Nome  string `json:"nome" validator:"required"`
	Senha string `json:"senha" validator:"required"`
}

type UserData struct {
	Login  string
	Nome   string
	Curso  string
	Codigo int
}
