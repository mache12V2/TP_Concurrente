package main

import (
	"fmt"
	"math/rand"
	"sort"
	"time"
)

type partialSums struct {
	sumX  float64
	sumY  float64
	sumXY float64
	sumXX float64
}

func worker(x, y []float64, sem chan struct{}, ch chan partialSums) {
	defer func() {
		<-sem // Libera el semáforo al finalizar
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

func generateData(n int) ([]float64, []float64) {
	x := make([]float64, n)
	y := make([]float64, n)
	for i := 0; i < n; i++ {
		x[i] = float64(i)
		y[i] = 2*float64(i) + 5 + rand.Float64()*10
	}
	return x, y
}

func calculate(x, y []float64, workers int) time.Duration {
	start := time.Now()
	size := len(x) / workers
	sem := make(chan struct{}, workers) // Canal que actúa como semáforo
	ch := make(chan partialSums, workers)

	for i := 0; i < workers; i++ {
		startIdx := i * size
		endIdx := startIdx + size
		if i == workers-1 {
			endIdx = len(x)
		}
		sem <- struct{}{} // Adquiere el semáforo
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

	duration := time.Since(start)
	fmt.Printf("m = %v, b = %v\n", m, b)
	return duration
}

func main() {
	const runs = 1000
	const workers = 4
	const numDataPoints = 1000000 // 1 millón de puntos

	durations := make([]time.Duration, runs)
	for i := 0; i < runs; i++ {
		x, y := generateData(numDataPoints)
		durations[i] = calculate(x, y, workers)
	}

	// Ordenar y calcular la media recortada
	sort.Slice(durations, func(i, j int) bool {
		return durations[i] < durations[j]
	})
	durations = durations[50 : len(durations)-50] // Eliminar los 50 menores y mayores

	var total time.Duration
	for _, duration := range durations {
		total += duration
	}

	average := total / time.Duration(len(durations))
	fmt.Printf("Media recortada de tiempo de ejecución: %v\n", average)
}
