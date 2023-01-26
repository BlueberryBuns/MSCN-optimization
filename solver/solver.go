package mscn

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/BlueberryBuns/MSCN-optimization/mscn"
	"github.com/BlueberryBuns/MSCN-optimization/solution"
	"github.com/BlueberryBuns/MSCN-optimization/utils"
	"github.com/gruntwork-io/terratest/modules/random"
)

const (
	TimeConstraint = iota
	EarlyStopping
)

func SolveProblem() {
	mscnInstance, err := mscn.GenerateMscnStructureFromFile("/Users/hulewicz/Private/mscn/solution.yaml")
	if err != nil {
		fmt.Printf("Error occurred, %v", err)
		return
	}

	var crossoverProbability float32 = 0.9
	var differentialWeight float32 = 0.8
	populationSize := 1000
	terminationCriterion := TimeConstraint //alternative: "early stop"
	timeConstraint := 1 * 60               //minutes * seconds
	iterationWithoutImprovement := 1000    // 1000

	isFinishConditionMet := getFinishCriterionFunction(terminationCriterion, timeConstraint, iterationWithoutImprovement)

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

	fmt.Printf("Best income from random solution %v\n", bestIncome)
	fmt.Println(population[bestSolutionIndex])
	fmt.Printf("Elapsed population time %v\n", populationElapsed)
	fmt.Printf("Elapsed fixing time %v\n", fixElapsed)

	fmt.Println("Co jessss")
	differentialEvolution(population, isFinishConditionMet, crossoverProbability, differentialWeight)
	// sc.FixEntities(population)
	// fmt.Printf("%v\n", sc.GenerateRandomSolution())
}

func maximizeCosts() {

}

func crossover() {

}

func differentialEvolution(population [][]float32, isFinishConditionMet func(int) bool, crossoverProbability float32, differentialWeight float32, currentBestScore float32) {
	var iterationsWithoutImprovement int
	var totalIterations int
	var previousBestScore float32 = currentBestScore

	for !isFinishConditionMet(iterationsWithoutImprovement) {
		for i := 0; i < len(population); i++ {
			x, y, z := generateRandomNumbersForDE(len(population), i)
		}

		if previousBestScore < currentBestScore {
			// Log this to csv file
			iterationsWithoutImprovement = -1
		}
		totalIterations++
		iterationsWithoutImprovement++
	}
}

func mutate(s1 []float32, s2 []float32, s3 []float32, differentialWeight float32) []float32 {
	afterMutation := make([]float32, len(s1))
	for i := 0; i < len(afterMutation); i++ {
		if randomNumber := rand.Float32(); randomNumber < 
		afterMutation[i] = s1[i] + (s2[i]-s3[i])*differentialWeight
	}
	return afterMutation
}

func getFinishCriterionFunction(criterion int, timeLimit int, iWithoutImprovementLimit int) func(int) bool {
	startTime := time.Now()

	switch criterion {
	case TimeConstraint:
		return func(iterations int) bool {
			if int(time.Now().Sub(startTime).Seconds()) > timeLimit {
				return true
			}
			return false
		}
	case EarlyStopping:
		return func(iterations int) bool {
			if int(time.Now().Sub(startTime).Seconds()) > timeLimit || iterations >= iWithoutImprovementLimit {
				return true
			}
			return false
		}
	default:
		return func(iterations int) bool {
			fmt.Println("execution finished immediately, invalid criterion was provided!")
			return true
		}

	}
}

func generateRandomNumbersForDE(populationSize int, currentIndex int) (int, int, int) {
	areResultValid := false
	var x, y, z int
	for !areResultValid {
		x = random.Random(0, populationSize)
		y = random.Random(0, populationSize)
		z = random.Random(0, populationSize)

		if currentIndex != x && currentIndex != y && currentIndex != z && x != z && x != y && y != x {
			areResultValid = true
		}
	}
	return x, y, z
}
