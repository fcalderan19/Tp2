package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	sistema "tp2/algogram"
	errores "tp2/errores"
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
		case "login": //O(1) hacerlo con diccionarios

			err := algogram.Loguearse(datoTerminal)

			if err != nil {
				fmt.Println(err)
				continue
			}

		case "logout":
			err := algogram.Desloggearse()

			if err != nil {
				fmt.Println(err)
				continue
			}

		case "publicar":
			likes := make([]string, 0)
			postNuevo := sistema.CrearPost(IDpost, algogram.UsuarioLoggeado(), datoTerminal, likes)
			err := algogram.PublicarPost(IDpost, postNuevo)

			if err != nil {
				fmt.Print(err)
				continue
			}

			IDpost++
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
		user := sistema.CrearUsuario(linea, afinidad, sistema.Cmp)
		usuarios.Guardar(linea, user)
	}

	return system
}

/*
func validarUsuario(usuario string, listaUsuarios TDA.Diccionario[string, sistema.Usuario]) error {
	if !listaUsuarios.Pertenece(usuario) {
		return errores.ErrorUsuarioInexistente{}
	}
	return nil
}
*/
