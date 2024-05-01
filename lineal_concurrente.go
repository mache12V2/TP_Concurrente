package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type partialSums struct {
	sumX  float64
	sumY  float64
	sumXY float64
	sumXX float64
}

func worker(x, y []float64, wg *sync.WaitGroup, ch chan partialSums) {
	defer wg.Done()
	var pSums partialSums
	for i := range x {
		pSums.sumX += x[i]
		pSums.sumY += y[i]
		pSums.sumXY += x[i] * y[i]
		pSums.sumXX += x[i] * x[i]
	}
	ch <- pSums
}

// generateData crea dos slices con valores x e y. Los valores y se generan usando
// una relación lineal simple con algo de ruido.
func generateData(n int) ([]float64, []float64) {
	x := make([]float64, n)
	y := make([]float64, n)
	for i := 0; i < n; i++ {
		x[i] = float64(i)
		y[i] = 2*float64(i) + 5 + rand.Float64()*10 // El ruido se ajusta con rand.Float64()*10
	}
	return x, y
}

func main() {
	x, y := generateData(1000000) // 1 millón de registros
	start := time.Now()

	workers := 4
	size := len(x) / workers
	var wg sync.WaitGroup
	ch := make(chan partialSums, workers)

	for i := 0; i < workers; i++ {
		startIdx := i * size
		endIdx := startIdx + size
		if i == workers-1 {
			endIdx = len(x) // Para manejar cualquier resto si N no es divisible exactamente por workers
		}
		wg.Add(1)
		go worker(x[startIdx:endIdx], y[startIdx:endIdx], &wg, ch)
	}

	go func() {
		wg.Wait()
		close(ch)
	}()

	total := partialSums{}
	for p := range ch {
		total.sumX += p.sumX
		total.sumY += p.sumY
		total.sumXY += p.sumXY
		total.sumXX += p.sumXX
	}

	N := float64(len(x))
	m := (N*total.sumXY - total.sumX*total.sumY) / (N*total.sumXX - total.sumX*total.sumX)
	b := (total.sumY - m*total.sumX) / N

	duration := time.Since(start)
	fmt.Printf("m = %v, b = %v\n", m, b)
	fmt.Printf("Tiempo de ejecución concurrente: %v\n", duration)
}
