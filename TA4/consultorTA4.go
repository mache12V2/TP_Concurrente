package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func main() {
	con, err := net.Dial("tcp", "UserIP:8000")
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
