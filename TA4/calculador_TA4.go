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

type partialSums struct {
	sumX  float64
	sumY  float64
	sumXY float64
	sumXX float64
}

type Empleado struct {
	salary float64
	gender int64
	age    int64
	PhD    int64
}

var (
	empleados []Empleado
)

func worker(x, y []float64, sem chan struct{}, ch chan partialSums) {
	defer func() {
		<-sem
	}()
	var pSums partialSums
	for i := range x {
		pSums.sumX += x[i]
		pSums.sumY += y[i]
		pSums.sumXY += x[i] * y[i]
		pSums.sumXX += x[i] * x[i]
	}
	ch <- pSums
}

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
	}
	fmt.Printf("Número de empleados: %d\n", len(empleados))
}

func generarDataCSV(empleados []Empleado) ([]float64, []float64) {
	n := len(empleados)
	X := make([]float64, n)
	Y := make([]float64, n)
	rand.Seed(time.Now().UnixNano())

	for i, empleado := range empleados {
		k := 3 + rand.Float64()*2
		X[i] = empleado.salary + k
		Y[i] = float64(empleado.age) + k
	}
	return X, Y
}

func calculate(x, y []float64, workers int) time.Duration {
	start := time.Now()
	size := len(x) / workers
	sem := make(chan struct{}, workers)
	ch := make(chan partialSums, workers)

	for i := 0; i < workers; i++ {
		startIdx := i * size
		endIdx := startIdx + size
		if i == workers-1 {
			endIdx = len(x)
		}
		sem <- struct{}{}
		go worker(x[startIdx:endIdx], y[startIdx:endIdx], sem, ch)
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

	N := float64(len(x))
	m := (N*total.sumXY - total.sumX*total.sumY) / (N*total.sumXX - total.sumX*total.sumX)
	b := (total.sumY - m*total.sumX) / N

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

func startCalc(runs int) (float64, float64, time.Duration) {
	const workers = 4

	leerDatos()
	x, y := generarDataCSV(empleados)

	durations := make([]time.Duration, runs)
	var totalDuration time.Duration

	for i := 0; i < runs; i++ {
		durations[i] = calculate(x, y, workers)
		totalDuration += durations[i]
	}

	trimmedMean := calcularMediaRecortada(durations)
	finalM, finalB := finalCalc(x, y, workers)

	return finalM, finalB, trimmedMean
}

func finalCalc(x, y []float64, workers int) (float64, float64) {
	size := len(x) / workers
	sem := make(chan struct{}, workers)
	ch := make(chan partialSums, workers)

	for i := 0; i < workers; i++ {
		startIdx := i * size
		endIdx := startIdx + size
		if i == workers-1 {
			endIdx = len(x)
		}
		sem <- struct{}{}
		go worker(x[startIdx:endIdx], y[startIdx:endIdx], sem, ch)
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

	N := float64(len(x))
	m := (N*total.sumXY - total.sumX*total.sumY) / (N*total.sumXX - total.sumX*total.sumX)
	b := (total.sumY - m*total.sumX) / N

	return m, b
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
