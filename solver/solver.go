package mscn

import (
	"fmt"

	"github.com/BlueberryBuns/MSCN-optimization/mscn"
	"github.com/BlueberryBuns/MSCN-optimization/solution"
)

func SolveProblem() {
	mscnInstance, err := mscn.GenerateMscnStructureFromFile("/Users/hulewicz/Private/mscn/solution.yaml")
	if err != nil {
		fmt.Printf("Error occurred, %v", err)
		return
	}

	solutionCalculator := solution.CreateSolutionCalculator(&mscnInstance)

	fmt.Printf("%v\n", solutionCalculator.GenerateRandomSolution())
}
