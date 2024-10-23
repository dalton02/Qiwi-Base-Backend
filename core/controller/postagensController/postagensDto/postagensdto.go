package postagensDto

type ConteudoJson struct {
	Tipo     string `json:"tipo" validator:"required"`
	Conteudo string `json:"conteudo" validator:"required"`
}

type PesquisarTituloQuerys struct {
	Tipo string `query:"tipo" validator:"required"`
}

type ListagemQuerys struct {
	Pagina   string   `query:"pagina"   validator:"optional,numericString"`
	Pesquisa string   `query:"pesquisa" validator:"optional"`
	Limite   string   `query:"limite"   validator:"optional,numericString"`
	Tags     []string `query:"tags"   validator:"optional"`
}

type ListagemPostagens struct {
	Postagem     []PostagemDataLista `json:"postagem"`
	Pagina       int                 `json:"pagina"`
	Pesquisa     string              `json:"pesquisa"`
	Limite       int                 `json:"limite"`
	TotalPaginas int                 `json:"TotalPaginas"`
}

type PostagemDataLista struct {
	Titulo      string         `json:"titulo" validator:"required"`
	Tipo        string         `json:"tipo" validator:"required"`
	Conteudo    string         `json:"conteudo" validator:"required"`
	Tags        []string       `json:"tags" validator:"required"`
	Autor       UserPostagem   `json:"autor" validator:"required"`
	Id          int            `json:"id"`
	Comentarios int            `json:"comentarios"`
	Reacoes     map[string]int `json:"reacoes"`
}

type PostagemDataComplete struct {
	Titulo      string                   `json:"titulo" validator:"required"`
	Tipo        string                   `json:"tipo" validator:"required"`
	Conteudo    string                   `json:"conteudo" validator:"required"`
	Tags        []string                 `json:"tags" validator:"required"`
	Autor       UserPostagem             `json:"autor" validator:"required"`
	Id          int                      `json:"id"`
	Comentarios []ComentarioDataComplete `json:"comentarios"`
	Reacoes     map[string]int           `json:"reacoes"`
}

type UserPostagem struct {
	Login string `json:"login"`
	Nome  string `json:"nome"`
	Curso string `json:"curso"`
	Id    int    `json:"id"`
}

type ComentarioDataComplete struct {
	Id       int                      `json:"id" validator:"required"`
	Conteudo string                   `json:"conteudo" validator:"required"`
	CriadoEm string                   `json:"criadoEm"`
	Autor    UserPostagem             `json:"autor" validator:"required"`
	Filhos   []ComentarioDataComplete `json:"filhos"`
}
type ComentarioData struct {
	Conteudo   string `json:"conteudo" validator:"required"`
	CriadoEm   string `json:"criadoEm"`
	UsuarioId  int    `json:"usuarioId" validator:"required"`
	PostagemId int    `json:"postagemId"` //Validado na url
	ParenteId  int    `json:"parenteId" validator:"optional"`
	Id         int    `json:"id"`
}

type ReacaoData struct {
	Tipo       string `json:"tipo" validator:"required"`
	UsuarioId  int    `json:"usuarioId" validator:"required"`
	PostagemId int    `json:"postagemId"` //Validado na url
	Id         int    `json:"id"`
}

type NovaPostagem struct {
	Titulo    string   `json:"titulo" validator:"required"`
	Tipo      string   `json:"tipo" validator:"required"`
	Conteudo  string   `json:"conteudo" validator:"required"`
	Tags      []string `json:"tags" validator:"required"`
	UsuarioId int      `json:"usuarioId" validator:"required"`
}
