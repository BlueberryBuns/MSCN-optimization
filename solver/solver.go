package mscn

import (
	"fmt"
	"time"

	"github.com/BlueberryBuns/MSCN-optimization/mscn"
	"github.com/BlueberryBuns/MSCN-optimization/solution"
)

func SolveProblem() {
	mscnInstance, err := mscn.GenerateMscnStructureFromFile("/Users/hulewicz/Private/mscn/sol.yaml")
	if err != nil {
		fmt.Printf("Error occurred, %v", err)
		return
	}

	populationSize := 10000
	sc := solution.CreateSolutionCalculator(&mscnInstance, populationSize)
	populationTimeStart := time.Now()
	population := sc.GeneratePopulation()
	populationElapsed := time.Now().Sub(populationTimeStart)

	// sc.DisplayPopulation(population)
	fixTimeStart := time.Now()
	sc.FixPopulation(population)
	fixElapsed := time.Now().Sub(fixTimeStart)

	incomes := sc.CalculatePopulationIncome(population)

	var bestIncome float32
	var bestSolutionIndex int
	for idx, income := range incomes {
		if bestIncome < income {
			bestIncome = income
			bestSolutionIndex = idx
		}
		// fmt.Printf("Total profit for solution %d:\t%.6f\n", idx, income)
	}

	fmt.Printf("Best income %v\n", bestIncome)
	fmt.Println(population[bestSolutionIndex])
	fmt.Printf("Elapsed population time %v\n", populationElapsed)
	fmt.Printf("Elapsed fixing time %v\n", fixElapsed)

	// sc.FixEntities(population)
	// fmt.Printf("%v\n", sc.GenerateRandomSolution())
}
