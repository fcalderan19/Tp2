package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	sistema "tp2/algogram"
	errores "tp2/errores"
)

const (
	login        = "login"
	logout       = "logout"
	publicarPost = "publicar"
	verFeed      = "ver_siguiente_feed"
	likear       = "likear_post"
	verLikes     = "mostrar_likes"
)

func main() {

	archivoUsuarios, err := leerUsuarios()
	if err != nil {
		fmt.Println(err)
		return
	}

	defer archivoUsuarios.Close()

	algogram := procesarUsuarios(archivoUsuarios)

	IDpost := 0

	terminal := bufio.NewScanner(os.Stdin)

	for terminal.Scan() {
		input := terminal.Text()
		input = strings.TrimSpace(input)
		params := strings.Fields(input)

		if len(params) == 0 {
			continue
		}

		cmd := params[0]
		data := params[1:]
		datoTerminal := strings.Join(data, " ")

		switch cmd {
		case login: //O(1) hacerlo con diccionarios

			err := algogram.Loguearse(datoTerminal)

			if err != nil {
				fmt.Println(err)
				continue
			}

		case logout:
			err := algogram.Desloggearse()

			if err != nil {
				fmt.Println(err)
				continue
			}

		case publicarPost:
			user := algogram.UsuarioLoggeado()
			if user == nil {
				fmt.Println(errores.ErrorUsuarioNoLoggeado{})
				continue
			}
			postNuevo := sistema.CrearPost(IDpost, user, datoTerminal, cmpString)
			err := algogram.PublicarPost(IDpost, postNuevo)
			if err != nil {
				fmt.Print(err)
				continue
			}
			algogram.GuardarPostEnFeed(postNuevo)
			IDpost++

		case verFeed:
			user := algogram.UsuarioLoggeado()
			if user == nil || user.VerFeed().EstaVacia() {
				fmt.Println(errores.ErrorFinFeed{})
				continue
			}

			post := user.VerProximoPost()
			post.VerContenido()

		case likear:
			ID, err := strconv.ParseInt(datoTerminal, 10, 64)
			if err != nil {
				continue
			}
			user := algogram.UsuarioLoggeado()

			if int(ID) > IDpost-1 || user == nil {
				fmt.Println(errores.ErrorLikearPost{})
			}

			algogram.LikearPost(int(ID))

		case verLikes:
			ID, err := strconv.ParseInt(datoTerminal, 10, 64)
			if err != nil {
				continue
			}
			if int(ID) > IDpost-1 {
				fmt.Println(errores.ErrorMostrarLike{})
				continue
			}
			post, err := algogram.MostrarLikes(int(ID))
			if err != nil {
				fmt.Println(errores.ErrorMostrarLike{})
				continue
			}
			post.VerLikes()
		}
	}

}

func leerUsuarios() (*os.File, error) {
	params := os.Args[1:]

	if len(params) != 1 {
		return nil, errores.ErrorParametros{}
	}

	archivoUsuarios := params[0]

	usuarios, error := os.Open(archivoUsuarios)
	if error != nil {
		return nil, errores.ErrorLeerArchivo{}
	}

	return usuarios, nil
}

func procesarUsuarios(archivo *os.File) sistema.Sistema {
	defer archivo.Close()

	system := sistema.CrearSistema()
	usuarios := system.UsuariosTotales()

	lineas := bufio.NewScanner(archivo)
	afinidad := 0

	for lineas.Scan() {
		linea := lineas.Text()
		afinidad++
		user := sistema.CrearUsuario(linea, afinidad, cmpPost)
		usuarios.Guardar(linea, user)
	}

	return system
}

func cmpPost(a, b sistema.Post) int {

	if a.VerPublicante().VerAfinidad() > b.VerPublicante().VerAfinidad() {
		return 1
	} else if a.VerPublicante().VerAfinidad() < b.VerPublicante().VerAfinidad() {
		return -1
	} else {
		if a.VerPostID() > b.VerPostID() {
			return 1
		} else {
			return -1
		}
	}
}

func cmpString(s1, s2 string) int {
	if s1 == s2 {
		return 0
	} else if s1 < s2 {
		return -1
	}
	return 1
}

/*
func validarUsuario(usuario string, listaUsuarios TDA.Diccionario[string, sistema.Usuario]) error {
	if !listaUsuarios.Pertenece(usuario) {
		return errores.ErrorUsuarioInexistente{}
	}
	return nil
}
*/
