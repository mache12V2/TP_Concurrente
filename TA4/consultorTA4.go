package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {
	userIP := bufio.NewReader(os.Stdin)
	fmt.Print("Ingese el ip del servidor: ")
	dir, _ := userIP.ReadString('\n')
	dir = strings.TrimSpace(dir)

	con, err := net.Dial("tcp", dir + ":8000")
	if err != nil {
		fmt.Println("Error al conectar al servidor:", err)
		return
	}
	defer con.Close()

	br := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("Ingrese accion a realizar\nPosibles inputs - 'promedio': ")
		msg, _ := br.ReadString('\n')

		fmt.Fprint(con, msg)

		resp, err := bufio.NewReader(con).ReadString('\n')
		if err != nil {
			fmt.Println("Error al recibir respuesta:", err)
			return
		}

		fmt.Println(resp)
	}
}
