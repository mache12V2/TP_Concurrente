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
	fmt.Print("Ingrese el ip del servidor: ")
	dir, _ := userIP.ReadString('\n')
	dir = strings.TrimSpace(dir)

	con, err := net.Dial("tcp", dir+":8000")
	if err != nil {
		fmt.Println("Error al conectar al servidor:", err)
		return
	}
	defer con.Close()

	br := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("Ingrese accion a realizar\nPosibles inputs - 'promedio', 'predecir': ")
		msg, _ := br.ReadString('\n')
		msg = strings.TrimSpace(msg)

		if msg == "predecir" {
			// Ask for gender
			fmt.Print("Ingrese el género (0 para femenino, 1 para masculino): ")
			gender, _ := br.ReadString('\n')
			gender = strings.TrimSpace(gender)

			// Ask for age
			fmt.Print("Ingrese la edad: ")
			age, _ := br.ReadString('\n')
			age = strings.TrimSpace(age)

			// Ask for PhD status
			fmt.Print("Ingrese el PhD (0 para no, 1 para sí): ")
			PhD, _ := br.ReadString('\n')
			PhD = strings.TrimSpace(PhD)

			// Send the 'predecir' command and the inputs to the server
			fmt.Fprint(con, "predecir\n")
			fmt.Fprint(con, gender+"\n")
			fmt.Fprint(con, age+"\n")
			fmt.Fprint(con, PhD+"\n")

			// Receive and print the response from the server
			resp, err := bufio.NewReader(con).ReadString('\n')
			if err != nil {
				fmt.Println("Error al recibir respuesta:", err)
				return
			}

			fmt.Println(resp)

		} else if msg == "promedio" {
			fmt.Fprint(con, msg+"\n")

			resp, err := bufio.NewReader(con).ReadString('\n')
			if err != nil {
				fmt.Println("Error al recibir respuesta:", err)
				return
			}

			fmt.Println(resp)
		} else {
			fmt.Println("Comando desconocido. Intente nuevamente.")
		}
	}
}
