package main

import (
	"math/rand"
	"time"

	s "github.com/BlueberryBuns/MSCN-optimization/solver"
)

func main() {
	rand.Seed(time.Now().UnixNano())
	s.SolveProblem()

}
