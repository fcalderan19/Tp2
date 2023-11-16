package Algogram

import (
	"fmt"
	heap "tdas/cola_prioridad"
	hash "tdas/diccionario"
	errores "tp2/errores"
)

type post struct {
	id        int
	user      string
	contenido string
	likes     []string
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

	VerLikes() []string
}

func CrearPost(id int, user, contenido string, likes []string) Post {
	return &post{id: id, user: user, contenido: contenido, likes: likes}
}

func (p *post) VerPostID() int {
	return p.id
}

func (p *post) VerPublicante() string {
	return p.user
}

func (p *post) VerContenido() string {
	return p.contenido
}

func (p *post) VerLikes() []string {
	return p.likes
}

// Usuarios
type Usuario interface {
	VerProximoPost() *Post

	VerAfinidad() int

	VerFeed() heap.ColaPrioridad[*Post]

	VerUsername() string
}

func Cmp(a, b *Post) int {
	//hacer la funcion que compare por afinidad
	return 0
}

func CrearUsuario(username string, afinidad int, cmp func(a *Post, b *Post) int) Usuario {
	var usuario usuario
	usuario.username = username
	usuario.afinidad = afinidad
	usuario.cmp = cmp
	usuario.feed = heap.CrearHeap[*Post](cmp)
	return &usuario
}

func (user *usuario) VerUsername() string {
	return user.username
}

func (user *usuario) VerAfinidad() int {
	return user.afinidad
}

func (user *usuario) VerFeed() heap.ColaPrioridad[*Post] {
	return user.feed
}

func (user *usuario) VerProximoPost() *Post {
	return user.feed.Desencolar()
}

type Sistema interface {

	//Dado un usuario, este ingresa al sistema. Si ya hay alguien adentro, error
	Loguearse(string) error

	//Desloguea al usuario actualemente loggeado. Si no hay nadie error
	Desloggearse() error

	//El usuario ingresado publica un contenido con un ID asignado y se los muestra en un diccionario
	PublicarPost(int, Post) error

	//Dado un post actual, likea. Si no hay post, error
	LikearPost(int) error

	//Muestra la cantidad de likes que tiene la publicacion x del usuario. Si no hay publicacion, error
	MostrarLikes(int) ([]string, error)

	UsuariosTotales() hash.Diccionario[string, Usuario]

	UsuarioLoggeado() string
}

func CrearSistema() Sistema {
	usuarios := hash.CrearHash[string, Usuario]()
	posteos := hash.CrearHash[int, Post]()
	return &sistema{usuariosTotales: usuarios, postsTotales: posteos, usuarioLogeado: ""}
}

func (system *sistema) Loguearse(user string) error {
	if !system.usuariosTotales.Pertenece(user) {
		return errores.ErrorUsuarioInexistente{}
	} else if system.usuarioLogeado != "" {
		return errores.ErrorUsuarioLoggeado{}
	}
	system.usuarioLogeado = user
	fmt.Println("Hola", user)
	return nil
}

func (system *sistema) Desloggearse() error {
	if system.usuarioLogeado == "" {
		return errores.ErrorUsuarioNoLoggeado{}
	}
	system.usuarioLogeado = ""
	fmt.Println("Adios")
	return nil
}

func (system *sistema) PublicarPost(id int, postNuevo Post) error {
	if system.usuarioLogeado == "" {
		return errores.ErrorUsuarioNoLoggeado{}
	}
	system.postsTotales.Guardar(id, postNuevo)
	fmt.Println("Post publicado")
	return nil
}

func (system *sistema) LikearPost(id int) error {
	var posteo Post
	if system.usuarioLogeado == "" {
		return errores.ErrorLikearPost{}
	}
	for iter := system.postsTotales.Iterador(); iter.HaySiguiente(); iter.Siguiente() {
		Id, posteoGuardado := iter.VerActual()
		if Id == id {
			posteo = posteoGuardado
			break
		}
	}
	if posteo == nil {
		return errores.ErrorLikearPost{}
	}

	user := system.usuarioLogeado

	likes := posteo.VerLikes()
	likes = append(likes, user)

	return nil
}

func (system *sistema) MostrarLikes(id int) ([]string, error) {
	var posteo Post
	var postVacio []string
	if system.usuarioLogeado == "" {
		return postVacio, errores.ErrorMostrarLike{}
	}
	for iter := system.postsTotales.Iterador(); iter.HaySiguiente(); iter.Siguiente() {
		Id, posteoGuardado := iter.VerActual()
		if Id == id {
			posteo = posteoGuardado
			break
		}
	}
	if posteo == nil {
		return postVacio, errores.ErrorMostrarLike{}
	}
	return posteo.VerLikes(), nil
}

func (system *sistema) UsuariosTotales() hash.Diccionario[string, Usuario] {
	return system.usuariosTotales
}

func (system sistema) UsuarioLoggeado() string {
	return system.usuarioLogeado
}
