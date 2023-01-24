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

	fmt.Printf("%v", mscnInstance)

	sc := solution.CreateSolutionCalculator(&mscnInstance)
	population := sc.GeneratePopulation(5)

	sc.DisplayPopulation(population)

	// fmt.Printf("%v\n", sc.GenerateRandomSolution())
}
