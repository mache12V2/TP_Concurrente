package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"os/exec"
	"path/filepath"
	"project/calculator"
	"strconv"
	"strings"
)

func listarDatos(res http.ResponseWriter, req *http.Request) {
	//tipo de contenido de la respuesta
	res.Header().Set("Content-Type", "application/json")
	// Asegúrate de que los datos necesarios estén cargados
	if len(calculator.Empleados) == 0 {
		calculator.LeerDatos()
	}
	//serializar y codificar el resultado a formato json
	jsonBytes, err := json.MarshalIndent(calculator.Empleados, "", " ")
	if err != nil {
		http.Error(res, fmt.Sprintf("Error al serializar datos: %v", err), http.StatusInternalServerError)
		return
	}
	res.Write(jsonBytes)
	log.Println("Respuesta exitosa!!")
}

func mostrarImagen(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "image/png") // Ajusta el tipo de contenido según tu necesidad
	// Verificar si hay puntos para graficar
	if len(calculator.Puntos) == 0 {
		http.Error(res, "No hay puntos para graficar. Realiza primero algunos cálculos de regresión.", http.StatusBadRequest)
		return
	}

	// Asegúrate de que el CSV esté generado y listo
	if err := calculator.WriteCSV(calculator.Puntos); err != nil {
		http.Error(res, fmt.Sprintf("Error al generar CSV: %v", err), http.StatusInternalServerError)
		return
	}
	pythonPath := filepath.Join("venv", "Scripts", "python")
	// Para Mac y Linux usar: pythonPath := filepath.Join("venv", "bin", "python")

	// Ejecutar el script Python para generar el gráfico
	cmd := exec.Command(pythonPath, "graph_empleados.py")
	stdout, err := cmd.Output()
	if err != nil {
		http.Error(res, fmt.Sprintf("Error al ejecutar el script: %v", err), http.StatusInternalServerError)
		return
	}

	// Escribir la salida del script como respuesta
	res.Write(stdout)
}

func calcularRegresion(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")

	// Leer datos y generar las matrices X e Y
	if len(calculator.Empleados) == 0 {
		calculator.LeerDatos()
	}
	X, Y := calculator.GenerarDataCSV(calculator.Empleados)

	// Realizar el cálculo de regresión
	m, b := calculator.FinalCalc(X, Y, 4)
	// Agregar el punto calculado a Puntos
	calculator.Puntos = append(calculator.Puntos, calculator.Point{Mval: m, Bval: b})

	// Crear la respuesta JSON
	response := map[string]float64{
		"slope":     m,
		"intercept": b,
	}
	//Serializa la respuesta a JSON
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		http.Error(res, "Error al crear la respuesta JSON", http.StatusInternalServerError)
		return
	}

	res.Write(jsonResponse)
}

func predecirSalario(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")

	// Extraer parámetros de la solicitud
	genderStr := req.URL.Query().Get("gender")
	ageStr := req.URL.Query().Get("age")
	PhDStr := req.URL.Query().Get("PhD")

	if genderStr == "" || ageStr == "" || PhDStr == "" {
		http.Error(res, "Parámetros insuficientes", http.StatusBadRequest)
		return
	}

	// Convertir parámetros a los tipos adecuados
	gender, err1 := strconv.ParseInt(genderStr, 10, 64)
	age, err2 := strconv.ParseInt(ageStr, 10, 64)
	PhD, err3 := strconv.ParseInt(PhDStr, 10, 64)

	if err1 != nil || err2 != nil || err3 != nil {
		http.Error(res, "Datos no válidos", http.StatusBadRequest)
		return
	}

	// Asegúrate de que los datos necesarios estén cargados
	if len(calculator.Empleados) == 0 {
		calculator.LeerDatos()
	}

	// Llamar a la función predictSalary con los parámetros adecuados
	predictedSalary := predictSalaryParams(gender, age, PhD)

	response := map[string]float64{
		"predicted_salary": predictedSalary,
	}

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		http.Error(res, "Error al crear la respuesta JSON", http.StatusInternalServerError)
		return
	}

	res.Write(jsonResponse)
}

func predictSalaryParams(gender, age, PhD int64) float64 {
	//Obtener el promedio de todas las características del empleado en relación al salario
	ageRange := calculator.GetAgeRange(age)
	averageSalary := calculator.TotalSalary / float64(calculator.TotalCount)
	averageSalaryPhD := calculator.AvgSalaryByPhD[PhD]
	averageSalaryGender := calculator.AvgSalaryByGender[gender]
	averageSalaryAgeRange := calculator.AvgSalaryByAgeRange[ageRange]
	//Retorna el promedio total
	return (averageSalary + averageSalaryPhD + averageSalaryGender + averageSalaryAgeRange) / 4
}

func ingresarParametros(res http.ResponseWriter, req *http.Request) {
	if req.URL.Path != "/" {
		http.Error(res, "404 not found", http.StatusNotFound)
		return
	}
	switch req.Method {
	case "GET":
		//solicita la página html
		http.ServeFile(res, req, "pagina.html") //Se cambia para angular
	case "POST":
		// Llama a ParseForm para parsear la consulta y actualizar r.PostForm y r.Form
		if err := req.ParseForm(); err != nil {
			http.Error(res, fmt.Sprintf("ParseForm() err: %v", err), http.StatusInternalServerError)
			return
		}

		//Convierte todos los strings obtenidos en ints
		genero, err1 := strconv.Atoi(req.FormValue("gender"))
		edad, err2 := strconv.Atoi(req.FormValue("age"))
		PhD, err3 := strconv.Atoi(req.FormValue("PhD"))
		//Validación de que la conversión se hizo exitosamente
		if err1 != nil || err2 != nil || err3 != nil {
			http.Error(res, "Datos no válidos", http.StatusBadRequest)
			return
		}

		//Enviar los parámetros al servidor
		msg, err := enviarParametros(genero, edad, PhD)
		if err != nil {
			http.Error(res, fmt.Sprintf("Error al enviar parámetros: %v", err), http.StatusInternalServerError)
			return
		}

		fmt.Fprintf(res, "Respuesta del servidor: %s", msg)
	default:
		http.Error(res, "Lo sentimos, solo se aceptan los métodos GET y POST", http.StatusMethodNotAllowed)
	}
}

func enviarParametros(genero, edad, PhD int) (string, error) {
	//puerto al que el nodo manda los parámetros
	remotehost := fmt.Sprintf("localhost:%d", 8000)
	con, err := net.Dial("tcp", remotehost)
	if err != nil {
		return "", fmt.Errorf("fallo en la conexión: %v", err)
	}
	defer con.Close()
	// Crear mensaje a enviar
	mensaje := fmt.Sprintf("%d,%d,%d", genero, edad, PhD)
	fmt.Fprintln(con, mensaje)

	bf := bufio.NewReader(con)
	msg, err := bf.ReadString('\n')

	if err != nil {
		return "", fmt.Errorf("fallo al leer la respuesta: %v", err)
	}
	return strings.TrimSpace(msg), nil
}

func inicializarDatos(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")

	//Inicailizando los datos para el calculator.go
	if len(calculator.Empleados) == 0 {
		calculator.LeerDatos()
	}
	//Haciendo el primer punto consiguiendo m y b
	X, Y := calculator.GenerarDataCSV(calculator.Empleados)
	m, b := calculator.FinalCalc(X, Y, 4)
	calculator.Puntos = append(calculator.Puntos, calculator.Point{Mval: m, Bval: b})

	response := map[string]string{
		"message": "Datos inicializados y primer cálculo realizado",
	}
	json.NewEncoder(res).Encode(response)
}

// Manejador de los endpoints
func manejadorRequest() {
	http.HandleFunc("/listar", listarDatos)
	http.HandleFunc("/grafico", mostrarImagen)
	http.HandleFunc("/salario", predecirSalario)
	http.HandleFunc("/regresion", calcularRegresion)
	http.HandleFunc("/inicializar", inicializarDatos)
	http.HandleFunc("/", ingresarParametros)
	log.Fatal(http.ListenAndServe(":9000", nil))
}

func main() {
	manejadorRequest()
}
