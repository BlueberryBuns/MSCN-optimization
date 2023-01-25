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

	populationSize := 1
	sc := solution.CreateSolutionCalculator(&mscnInstance, populationSize)
	population := sc.GeneratePopulation()

	// sc.DisplayPopulation(population)

	incomes := sc.CalculatePopulationIncome(population)

	for idx, income := range incomes {
		fmt.Printf("Total profit for solution %d:\t%.6f\n", idx, income)
	}

	// sc.FixEntities(population)
	// fmt.Printf("%v\n", sc.GenerateRandomSolution())
}
