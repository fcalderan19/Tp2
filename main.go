package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"tp2/errores"
)

func main() {

	archivoUsuarios, err := leerUsuarios()
	if err != nil {
		fmt.Println(err)
		return
	}

	defer archivoUsuarios.Close()

	listaUsuarios := procesarUsuarios(archivoUsuarios)
	usuarioLoggeado := ""

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

			err := validarUsuario(datoTerminal, listaUsuarios)

			if err != nil {
				fmt.Println(err)
				return
			}

			if usuarioLoggeado != "" {
				fmt.Println(errores.ErrorUsuarioLoggeado{})
				return
			}

			fmt.Println("Hola", datoTerminal)
			usuarioLoggeado = datoTerminal

		case "logout":
			if usuarioLoggeado == "" {
				fmt.Println(errores.ErrorUsuarioNoLoggeado{})
				return
			}
			println("Adios")
			usuarioLoggeado = ""
		}

	}
}

func leerUsuarios() (*os.File, error) {
	params := os.Args[1:]

	if len(params) != 2 {
		return nil, errores.ErrorParametros{}
	}

	archivoUsuarios := params[0]

	usuarios, error := os.Open(archivoUsuarios)
	if error != nil {
		return nil, errores.ErrorLeerArchivo{}
	}

	return usuarios, nil
}

func procesarUsuarios(archivo *os.File) []string {
	defer archivo.Close()

	usuarios := make([]string, 0)

	lineas := bufio.NewScanner(archivo)

	for lineas.Scan() {
		linea := lineas.Text()
		usuarios = append(usuarios, linea)
	}

	return usuarios
}

func validarUsuario(usuario string, listaUsuarios []string) error {
	for user := range listaUsuarios {
		if listaUsuarios[user] == usuario {
			return errores.ErrorUsuarioLoggeado{}
		}
	}
	return nil
}
