package errores

type ErrorLeerArchivo struct{}

func (e ErrorLeerArchivo) Error() string {
	return "Error: Lectura de archivos"
}

type ErrorParametros struct{}

func (e ErrorParametros) Error() string {
	return "Error: Faltan parámetros"
}

type ErrorUsuarioLoggeado struct{}

func (e ErrorUsuarioLoggeado) Error() string {
	return "Error: Ya habia un usuario loggeado"
}

type ErrorUsuarioInexistente struct{}

func (e ErrorUsuarioInexistente) Error() string {
	return "Error: usuario no existente"
}

type ErrorUsuarioNoLoggeado struct{}

func (e ErrorUsuarioNoLoggeado) Error() string {
	return "Error: no habia usuario loggeado"
}

type ErrorFinFeed struct{}

func (e ErrorFinFeed) Error() string {
	return "Error: Usuario no loggeado o no hay mas posts para ver"
}

type ErrorLikearPost struct{}

func (e ErrorLikearPost) Error() string {
	return "Error: Usuario no loggeado o Post inexistente"
}

type ErrorMostrarLike struct{}

func (e ErrorMostrarLike) Error() string {
	return "Error: Post inexistente o sin likes"
}
