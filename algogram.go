package Algogram

import (
	heap "tdas/cola_prioridad"
	hash "tdas/diccionario"
	errores "tp2/errores"
)

type post struct {
	id        int
	user      string
	contenido string
	likes     []Usuario
}

type usuario struct {
	username string
	afinidad int
	cmp      func(a, b *Post) int
	feed     heap.ColaPrioridad[*Post] //Hago heap de minimos por afinidad y desencolo (desencolaria el primero que es el que menos valor de afinidad tiene)
}

type sistema struct {
	usuariosTotales hash.Diccionario[string, Usuario]
	postsTotales    hash.Diccionario[int, Post]
	usuarioLogeado  string
}

// Posts
type Post interface {
	VerPostID() int

	VerPublicante() string

	VerContenido() string

	VerLikes() []Usuario
}

func CrearPost(id int, user, contenido string, likes []Usuario) Post {
	return &post{id: id, user: user, contenido: contenido, likes: likes}
}

func (p post) VerPostID() int {
	return p.id
}

func (p post) VerPublicante() string {
	return p.user
}

func (p post) VerContenido() string {
	return p.contenido
}

func (p post) VerLikes() []Usuario {
	return p.likes
}

// Usuarios
type Usuario interface {
	VerProximoPost() *Post

	VerAfinidad() int

	VerFeed() heap.ColaPrioridad[*Post]

	VerUsername() string
}

func CrearUsuario(username string, afinidad int, cmp func(a *Post, b *Post) int) Usuario {
	var usuario usuario
	usuario.username = username
	usuario.afinidad = afinidad
	usuario.cmp = cmp
	usuario.feed = heap.CrearHeap[*Post](cmp)
	return usuario
}

func (user usuario) VerUsername() string {
	return user.username
}

func (user usuario) VerAfinidad() int {
	return user.afinidad
}

func (user usuario) VerFeed() heap.ColaPrioridad[*Post] {
	return user.feed
}

func (user usuario) VerProximoPost() *Post {
	return user.feed.Desencolar()
}

type Sistema interface {

	//Dado un usuario, este ingresa al sistema. Si ya hay alguien adentro, error
	Loguearse(string) error

	//Desloguea al usuario actualemente loggeado. Si no hay nadie error
	Desloggearse() error

	//El usuario ingresado publica un contenido con un ID asignado y se los muestra en un diccionario
	PublicarPost(Post) error

	//Dado un post actual, likea. Si no hay post, error
	LikearPost(int, Post) error

	//Muestra la cantidad de likes que tiene la publicacion x del usuario. Si no hay publicacion, error
	MostrarLikes(int, Post) error
}

func (system sistema) Loguearse(user string) error {
	if !system.usuariosTotales.Pertenece(user) {
		return errores.ErrorUsuarioInexistente{}
	}
	system.usuarioLogeado = user
	return nil
}

func (system sistema) Desloggearse() error {
	if system.usuarioLogeado == "" {
		return errores.ErrorUsuarioNoLoggeado{}
	}
	return nil
}

func (system sistema) PublicarPost(postNuevo Post) error {

}
