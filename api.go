package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"strconv"
	"strings"
)

func listarDatos(res http.ResponseWriter, req *http.Request) {
	//tipo de contenido de la respuesta
	res.Header().Set("Content-Type", "application/json")
	//serializar y codificar el resultado a formato json
	jsonBytes, err := json.MarshalIndent(empleados, "", " ")
	if err != nil {
		http.Error(res, fmt.Sprintf("Error al serializar datos: %v", err), http.StatusInternalServerError)
		return
	}
	res.Write(jsonBytes)
	log.Println("Respuesta exitosa!!")
}
func mostrarImagen(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")
	//crear la gráfica en python
}
func predecirSalario(res http.ResponseWriter, req *http.Request) {
	//funcion a realizar
	res.Header().Set("Content-Type", "application/json")

}
func ingresarParametros(res http.ResponseWriter, req *http.Request) {
	if req.URL.Path != "/" {
		http.Error(res, "404 not found", http.StatusNotFound)
	}
	switch req.Method {
	case "GET":
		//solocita la página html
		http.ServeFile(res, req, "pagina.html")
	case "POST":
		//Llama a ParseForm para parsear la consulta y actualizar r.PostForm y r.Form
		if err := req.ParseForm(); err != nil {
			http.Error(res, fmt.Sprintf("ParseForm() err: %v", err), http.StatusInternalServerError)
			return
		}

		genero, err1 := strconv.Atoi(req.FormValue("gender"))
		edad, err2 := strconv.Atoi(req.FormValue("age"))
		Phd, err3 := strconv.Atoi(req.FormValue("Phd"))
		if err1 != nil || err2 != nil || err3 != nil {
			http.Error(res, "Datos no validos", http.StatusBadRequest)
			return
		}
		msg, err := enviarParametros(genero, edad, Phd)
		if err != nil {
			http.Error(res, fmt.Sprintf("Error al enviar parámetros: %v", err), http.StatusInternalServerError)
		}

		fmt.Fprintln(res, "Respuesta del servidor: %s", msg)
	default:
		http.Error(res, "Lo sentimos, solo se acepto el método GET y POST", http.StatusMethodNotAllowed)
	}
}

func enviarParametros(genero, edad, Phd int) (string, error) {
	//puerto al que el nodo manda los parámetros
	remotehost := fmt.Sprintf("localhost:%d", 8000)
	con, err := net.Dial("tcp", remotehost)
	if err != nil {
		return "", fmt.Errorf("fallo en la conexión: %v", err)
	}
	defer con.Close()
	//Crear mensaje a enviar
	mensaje := fmt.Sprintf("%d,%d,%d", genero, edad, Phd)
	fmt.Fprintln(con, mensaje)

	bf := bufio.NewReader(con)
	msg, err := bf.ReadString('\n')

	if err != nil {
		return "", fmt.Errorf("fallo al leer la respuesta: %v", err)
	}
	return strings.TrimSpace(msg), nil
}

func manejadorRequest() {
	http.HandleFunc("/listar", listarDatos)
	http.HandleFunc("/grafico", mostrarImagen)
	http.HandleFunc("/salario", predecirSalario)
	http.HandleFunc("/", ingresarParametros)
	log.Fatal(http.ListenAndServe(":9000", nil))
}

func main() {
	manejadorRequest()
}
