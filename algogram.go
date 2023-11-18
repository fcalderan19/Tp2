package Algogram

import (
	"fmt"
	heap "tdas/cola_prioridad"
	hash "tdas/diccionario"
	errores "tp2/errores"
)

type post struct {
	id        int
	user      Usuario
	contenido string
	likes     hash.DiccionarioOrdenado[string, Usuario]
}

type usuario struct {
	username string
	afinidad int
	cmp      func(a, b Post) int
	feed     heap.ColaPrioridad[Post] //Hago heap de minimos por afinidad y desencolo (desencolaria el primero que es el que menos valor de afinidad tiene)
}

type sistema struct {
	usuariosTotales hash.Diccionario[string, Usuario]
	postsTotales    hash.Diccionario[int, Post]
	usuarioLogeado  Usuario
}

// Posts
type Post interface {
	VerPostID() int

	VerPublicante() Usuario

	VerContenido()

	VerLikes()

	LikearPost(Usuario)
}

func CrearPost(id int, user Usuario, contenido string, cmp func(a, b string) int) Post {
	likes := hash.CrearABB[string, Usuario](cmp)
	return &post{id: id, user: user, contenido: contenido, likes: likes}
}

func (p *post) VerPostID() int {
	return p.id
}

func (p *post) VerPublicante() Usuario {
	return p.user
}

func (p *post) VerContenido() {
	fmt.Println("Post ID", p.id)
	fmt.Println(p.user, "dijo:", p.contenido)
	fmt.Println("Likes:", p.likes.Cantidad())
}

func (p *post) VerLikes() {
	fmt.Println("El post tiene", p.likes.Cantidad(), "likes:")
	for iter := p.likes.Iterador(); iter.HaySiguiente(); iter.Siguiente() {
		username, _ := iter.VerActual()
		fmt.Println(username)
	}
}

func (p *post) LikearPost(user Usuario) {
	p.likes.Guardar(user.VerUsername(), user)
}

// Usuarios
type Usuario interface {
	VerProximoPost() Post

	VerAfinidad() int

	VerFeed() heap.ColaPrioridad[Post]

	VerUsername() string

	GuardarPostFeedUsuario(Post)
}

func CrearUsuario(username string, afinidad int, cmp func(a, b Post) int) Usuario {
	var usuario usuario
	usuario.username = username
	usuario.afinidad = afinidad
	usuario.cmp = cmp
	usuario.feed = heap.CrearHeap[Post](cmp)
	return &usuario
}

func (user *usuario) VerUsername() string {
	return user.username
}

func (user *usuario) VerAfinidad() int {
	return user.afinidad
}

func (user *usuario) VerFeed() heap.ColaPrioridad[Post] {
	return user.feed
}

func (user *usuario) VerProximoPost() Post {
	return user.feed.Desencolar()
}

func (user *usuario) GuardarPostFeedUsuario(post Post) {
	if user != post.VerPublicante() {
		user.feed.Encolar(post)
	}
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
	MostrarLikes(int) (Post, error)

	UsuariosTotales() hash.Diccionario[string, Usuario]

	UsuarioLoggeado() Usuario

	GuardarPostEnFeed(Post)

	PostsTotales() hash.Diccionario[int, Post]
}

func CrearSistema() Sistema {
	usuarios := hash.CrearHash[string, Usuario]()
	posteos := hash.CrearHash[int, Post]()
	return &sistema{usuariosTotales: usuarios, postsTotales: posteos, usuarioLogeado: nil}
}

func (system *sistema) Loguearse(user string) error {
	if !system.usuariosTotales.Pertenece(user) {
		return errores.ErrorUsuarioInexistente{}
	} else if system.usuarioLogeado != nil {
		return errores.ErrorUsuarioLoggeado{}
	}
	system.usuarioLogeado = system.usuariosTotales.Obtener(user)
	fmt.Println("Hola", user)
	return nil
}

func (system *sistema) Desloggearse() error {
	if system.usuarioLogeado == nil {
		return errores.ErrorUsuarioNoLoggeado{}
	}
	system.usuarioLogeado = nil
	fmt.Println("Adios")
	return nil
}

func (system *sistema) PublicarPost(id int, postNuevo Post) error {
	if system.usuarioLogeado == nil {
		return errores.ErrorUsuarioNoLoggeado{}
	}
	system.postsTotales.Guardar(id, postNuevo)
	fmt.Println("Post publicado")
	return nil
}

func (system *sistema) LikearPost(id int) error {
	var posteo Post
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

	posteo.LikearPost(system.usuarioLogeado)
	fmt.Println("Post likeado")

	return nil
}

func (system *sistema) MostrarLikes(id int) (Post, error) {
	for iter := system.postsTotales.Iterador(); iter.HaySiguiente(); iter.Siguiente() {
		Id, posteoGuardado := iter.VerActual()
		if Id == id {
			return posteoGuardado, nil
		}
	}
	return nil, errores.ErrorMostrarLike{}
}

func (system *sistema) UsuariosTotales() hash.Diccionario[string, Usuario] {
	return system.usuariosTotales
}

func (system sistema) UsuarioLoggeado() Usuario {
	return system.usuarioLogeado
}

func (system *sistema) GuardarPostEnFeed(post Post) {
	for iter := system.usuariosTotales.Iterador(); iter.HaySiguiente(); iter.Siguiente() {
		username, usuario := iter.VerActual()
		if username != post.VerPublicante().VerUsername() {
			usuario.GuardarPostFeedUsuario(post)
		}
	}
}

func (system sistema) PostsTotales() hash.Diccionario[int, Post] {
	return system.postsTotales
}
