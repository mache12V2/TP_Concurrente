package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"log"
	"math/rand"
	"net"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

type Empleado struct {
	salary float64
	gender int64
	age    int64
	PhD    int64
}

type partialSums struct {
	sumX  float64
	sumY  float64
	sumXY float64
	sumXX float64
}

type Point struct {
	mval float64
	bval float64
}

type AgeRange string

const (
	Age18To24 AgeRange = "18-24"
	Age25To49 AgeRange = "25-49"
	Age50To77 AgeRange = "50-77"
)

func getAgeRange(age int64) AgeRange {
	switch {
	case age >= 18 && age <= 24:
		return Age18To24
	case age >= 25 && age <= 49:
		return Age25To49
	case age >= 50 && age <= 77:
		return Age50To77
	default:
		return ""
	}
}

var (
	empleados           []Empleado
	puntos              []Point
	totalSalary         float64
	totalCount          int64
	avgSalaryByPhD      = make(map[int64]float64)
	countByPhD          = make(map[int64]int64)
	avgSalaryByGender   = make(map[int64]float64)
	countByGender       = make(map[int64]int64)
	avgSalaryByAgeRange = make(map[AgeRange]float64)
	countByAgeRange     = make(map[AgeRange]int64)
)

func leerDatos() {
	url := "https://raw.githubusercontent.com/mache12V2/TP_Concurrente/main/TA3/Datasets/Updated_Expanded_Salary_Data.csv"
	resp, err := http.Get(url)
	if err != nil {
		log.Fatalf("Error al realizar la solicitud HTTP: %v", err)
	}
	fmt.Println("Se encontró el archivo CSV con éxito")
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Fatalf("Error al obtener el archivo CSV: %v", resp.Status)
	}

	reader := csv.NewReader(resp.Body)

	_, err = reader.Read()
	if err != nil {
		fmt.Println("Error al leer la primera fila:", err)
		return
	}

	csvLines, err := reader.ReadAll()
	if err != nil {
		log.Fatalf("Error al leer el archivo CSV: %v", err)
	}
	fmt.Printf("Número de filas en el archivo CSV: %d\n", len(csvLines))

	for _, line := range csvLines {
		salary, err := strconv.ParseFloat(line[0], 64)
		gender, err1 := strconv.Atoi(line[1])
		age, err2 := strconv.Atoi(line[2])
		PhD, err3 := strconv.Atoi(line[3])

		if err != nil || err1 != nil || err2 != nil || err3 != nil {
			fmt.Println("Error en la lectura de los datos")
			return
		}

		emp := Empleado{
			salary: salary,
			gender: int64(gender),
			age:    int64(age),
			PhD:    int64(PhD),
		}

		empleados = append(empleados, emp)
		totalSalary += emp.salary
		totalCount++

		avgSalaryByPhD[emp.PhD] += emp.salary
		countByPhD[emp.PhD]++

		avgSalaryByGender[emp.gender] += emp.salary
		countByGender[emp.gender]++

		ageRange := getAgeRange(emp.age)
		avgSalaryByAgeRange[ageRange] += emp.salary
		countByAgeRange[ageRange]++
	}

	fmt.Printf("Número de empleados: %d\n", len(empleados))

	for key, total := range avgSalaryByPhD {
		avgSalaryByPhD[key] = total / float64(countByPhD[key])
	}
	for key, total := range avgSalaryByGender {
		avgSalaryByGender[key] = total / float64(countByGender[key])
	}
	for key, total := range avgSalaryByAgeRange {
		avgSalaryByAgeRange[key] = total / float64(countByAgeRange[key])
	}
}

// Modificada para crear una matriz de variables predictoras (X) y una variable objetivo (Y)
func generarDataCSV(empleados []Empleado) ([][]float64, []float64) {
	n := len(empleados)
	X := make([][]float64, n) // Matriz para las variables predictoras
	Y := make([]float64, n)   // Vector para la variable objetivo
	rand.Seed(time.Now().UnixNano())

	for i, empleado := range empleados {
		k := 3 + rand.Float64()*2
		X[i] = []float64{float64(empleado.gender), float64(empleado.age), float64(empleado.PhD)} // Incluye todas las variables predictoras
		Y[i] = empleado.salary + k
	}
	return X, Y
}

// Modificada para manejar una matriz de variables predictoras y un vector de variables objetivo
func worker(X [][]float64, Y []float64, sem chan struct{}, ch chan partialSums) {
	defer func() {
		<-sem
	}()
	var pSums partialSums
	for i := range X {
		for j := range X[i] {
			pSums.sumX += X[i][j]            // Suma de las variables predictoras
			pSums.sumXY += X[i][j] * Y[i]    // Producto de las variables predictoras y la variable objetivo
			pSums.sumXX += X[i][j] * X[i][j] // Producto de las variables predictoras
		}
		pSums.sumY += Y[i] // Suma de la variable objetivo
	}
	ch <- pSums
}

// Modificada para manejar una matriz de variables predictoras y un vector de variables objetivo
func calculate(X [][]float64, Y []float64, workers int) time.Duration {
	start := time.Now()
	size := len(X) / workers
	sem := make(chan struct{}, workers)
	ch := make(chan partialSums, workers)

	for i := 0; i < workers; i++ {
		startIdx := i * size
		endIdx := startIdx + size
		if i == workers-1 {
			endIdx = len(X)
		}
		sem <- struct{}{}
		go worker(X[startIdx:endIdx], Y[startIdx:endIdx], sem, ch)
	}

	desv_x := rand.Float64()
	desv_y := rand.Float64()

	total := partialSums{}
	for i := 0; i < workers; i++ {
		p := <-ch
		total.sumX += p.sumX + desv_x
		total.sumY += p.sumY + desv_y
		total.sumXY += p.sumXY
		total.sumXX += p.sumXX
	}
	close(ch)

	N := float64(len(X))
	m := (N*total.sumXY - total.sumX*total.sumY) / (N*total.sumXX - total.sumX*total.sumX)
	b := (total.sumY - m*total.sumX) / N

	if len(puntos) <= 1000000 {
		punto_n := Point{mval: m, bval: b}
		puntos = append(puntos, punto_n)
	}

	fmt.Printf("m = %v, b = %v\n", m, b)
	return time.Since(start)
}

func calcularMediaRecortada(durations []time.Duration) time.Duration {
	sort.Slice(durations, func(i, j int) bool {
		return durations[i] < durations[j]
	})
	durations = durations[50 : len(durations)-50]

	var total time.Duration
	for _, duration := range durations {
		total += duration
	}

	average := total / time.Duration(len(durations))
	return average
}

// Modificada para manejar una matriz de variables predictoras y un vector de variables objetivo
func startCalc(runs int) (float64, float64, time.Duration) {
	const workers = 4

	leerDatos()
	X, Y := generarDataCSV(empleados)

	durations := make([]time.Duration, runs)
	var totalDuration time.Duration

	for i := 0; i < runs; i++ {
		durations[i] = calculate(X, Y, workers)
		totalDuration += durations[i]
	}

	trimmedMean := calcularMediaRecortada(durations)
	finalM, finalB := finalCalc(X, Y, workers)

	return finalM, finalB, trimmedMean
}

// Modificada para manejar una matriz de variables predictoras y un vector de variables objetivo
func finalCalc(X [][]float64, Y []float64, workers int) (float64, float64) {
	size := len(X) / workers
	sem := make(chan struct{}, workers)
	ch := make(chan partialSums, workers)

	for i := 0; i < workers; i++ {
		startIdx := i * size
		endIdx := startIdx + size
		if i == workers-1 {
			endIdx = len(X)
		}
		sem <- struct{}{}
		go worker(X[startIdx:endIdx], Y[startIdx:endIdx], sem, ch)
	}

	total := partialSums{}
	for i := 0; i < workers; i++ {
		p := <-ch
		total.sumX += p.sumX
		total.sumY += p.sumY
		total.sumXY += p.sumXY
		total.sumXX += p.sumXX
	}
	close(ch)

	N := float64(len(X))
	m := (N*total.sumXY - total.sumX*total.sumY) / (N*total.sumXX - total.sumX*total.sumX)
	b := (total.sumY - m*total.sumX) / N

	return m, b
}

func writeCSV(puntos []Point) {
	file, err := os.Create("points_Empleados.csv")
	if err != nil {
		log.Fatal("Error creating file:", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	headers := []string{"m", "b"}
	if err := writer.Write(headers); err != nil {
		log.Fatal("Error writing headers:", err)
	}

	for i := 0; i < len(puntos); i++ {
		record := []string{strconv.FormatFloat(puntos[i].mval, 'f', -1, 64), strconv.FormatFloat(puntos[i].bval, 'f', -1, 64)}
		if err := writer.Write(record); err != nil {
			log.Fatal("Error writing record to CSV:", err)
		}
	}
}

// Predicción de salario según el input del usuario y promedios de la base de datos
func predictSalary(con net.Conn) {
	reader := bufio.NewReader(con)

	genderStr, _ := reader.ReadString('\n')
	genderStr = strings.TrimSpace(genderStr)
	gender, err := strconv.ParseInt(genderStr, 10, 64)
	if err != nil {
		fmt.Fprintln(con, "Error al leer el género")
		return
	}

	ageStr, _ := reader.ReadString('\n')
	ageStr = strings.TrimSpace(ageStr)
	age, err := strconv.ParseInt(ageStr, 10, 64)
	if err != nil {
		fmt.Fprintln(con, "Error al leer la edad")
		return
	}

	PhDStr, _ := reader.ReadString('\n')
	PhDStr = strings.TrimSpace(PhDStr)
	PhD, err := strconv.ParseInt(PhDStr, 10, 64)
	if err != nil {
		fmt.Fprintln(con, "Error al leer el PhD")
		return
	}

	ageRange := getAgeRange(age)
	averageSalary := totalSalary / float64(totalCount)
	averageSalaryPhD := avgSalaryByPhD[PhD]
	averageSalaryGender := avgSalaryByGender[gender]
	averageSalaryAgeRange := avgSalaryByAgeRange[ageRange]

	predictedSalary := (averageSalary + averageSalaryPhD + averageSalaryGender + averageSalaryAgeRange) / 4
	fmt.Fprintf(con, "El salario predicho es: %.2f\n", predictedSalary)
}

func manejador(con net.Conn) {
	defer con.Close()
	br := bufio.NewReader(con)

	for {
		datos, err := br.ReadString('\n')
		if err != nil {
			fmt.Println("Error leyendo data:", err)
			return
		}
		datos = strings.TrimSpace(datos)
		fmt.Println("Recibido:", datos)

		if datos == "promedio" {
			const runs = 1000
			m, b, trimmedMean := startCalc(runs)
			fmt.Fprintf(con, "m = %v, b = %v, Media recortada de tiempo de ejecución: %v\n", m, b, trimmedMean)
			writeCSV(puntos)
		} else if datos == "predecir" {
			predictSalary(con)
		} else {
			fmt.Fprintln(con, "Comando desconocido")
		}
	}
}

func main() {
	userIP := bufio.NewReader(os.Stdin)
	fmt.Print("Ingese el ip del cliente: ")
	dir, _ := userIP.ReadString('\n')
	dir = strings.TrimSpace(dir)

	ls, err := net.Listen("tcp", dir+":8000")
	if err != nil {
		fmt.Println("Fallo en la comunicación ", err.Error())
		os.Exit(1)
	}
	defer ls.Close()

	for {
		con, err := ls.Accept()
		if err != nil {
			fmt.Println("Fallo en la conexión ", err.Error())
		}
		go manejador(con)
	}
}
