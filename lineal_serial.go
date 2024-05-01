package main

import (
	"fmt"
	"math/rand"
	"time"
)

func generateData(n int) ([]float64, []float64) {
	x := make([]float64, n)
	y := make([]float64, n)
	for i := range x {
		x[i] = float64(i)
		y[i] = 2.0*float64(i) + float64(rand.Intn(10)) 
	}
	return x, y
}

func linearRegression(x, y []float64) (float64, float64) {
	sumX, sumY, sumXY, sumXX := 0.0, 0.0, 0.0, 0.0
	N := float64(len(x))

	for i := range x {
		sumX += x[i]
		sumY += y[i]
		sumXY += x[i] * y[i]
		sumXX += x[i] * x[i]
	}

	m := (N*sumXY - sumX*sumY) / (N*sumXX - sumX*sumX)
	b := (sumY - m*sumX) / N
	return m, b
}

func main() {
	x, y := generateData(1000000) 
	start := time.Now()
	m, b := linearRegression(x, y)
	duration := time.Since(start)
	fmt.Printf("m = %v, b = %v\n", m, b)
	fmt.Printf("Tiempo de ejecución: %v\n", duration)
}
