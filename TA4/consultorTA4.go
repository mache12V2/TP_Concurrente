package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func main() {
	con, _ := net.Dial("tcp", "calculadorIP:8000")

	br := bufio.NewReader(con)
	br2 := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("Ingrese su solicitud de promedio: ")
		msg, _ := br2.ReadString('\n')

		fmt.Fprint(con, msg)

		resp, _ := br.ReadString('\n')
		fmt.Println("Respuesta: ", resp)

	}
}
