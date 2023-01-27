package mscn

import (
	"bufio"
	"fmt"
	"math"
	"math/rand"
	"os"
	"path/filepath"
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
	instanceSpecifier := "" //"mscn_100_56_40_40"
	currentTime := time.Now()
	iPath := fmt.Sprintf("./solution_%s.yaml", instanceSpecifier)
	oPath := fmt.Sprintf("./result_%s_%s.csv", instanceSpecifier, currentTime.Format("2006-01-02 15:04:05"))
	instancePath, _ := filepath.Abs(iPath)
	resultPath, _ := filepath.Abs(oPath)
	mscnInstance, err := mscn.GenerateMscnStructureFromFile(instancePath)
	if err != nil {
		fmt.Printf("Error occurred, %v", err)
		return
	}

	f, err := os.Create(resultPath)
	if err != nil {
		panic("Error while opening file!\n")
	}
	w := bufio.NewWriter(f)
	defer f.Close()

	var crossoverProbability float32 = 0.9
	var differentialWeight float32 = 0.8
	populationSize := 1000
	terminationCriterion := TimeConstraint //alternative: "early stop"
	timeConstraint := 1 * 4 * 1            //minutes * seconds
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

	deStartTime := time.Now()
	getElapsedTime := getTimeElapsedCounter(deStartTime)

	differentialEvolution(&sc, population, isFinishConditionMet, crossoverProbability, differentialWeight, bestIncome, incomes, w, getElapsedTime)
	// sc.FixEntities(population)
	// fmt.Printf("%v\n", sc.GenerateRandomSolution())
}

func differentialEvolution(sc *solution.SolutionCalculator, population [][]float32, isFinishConditionMet func(int) bool, crossoverProbability float32, differentialWeight float32, currentBestScore float32, incomes []float32, w *bufio.Writer, getElapsedTime func() float32) {
	var iterationsWithoutImprovement int
	var totalIterations int
	var previousBestScore float32 = currentBestScore
	var bestSpecimen []float32
	var newSpecimen []float32

	// var newSpecimenBeforeRepair []float32
	// var oldSpecimen []float32
	var buf = make([]byte, 1024)
	copy(buf, "T elapsed [s],Best solution, epoch, epoch since best\n")
	w.Write(buf)
	w.Flush()
	buf = make([]byte, 1024)
	var elapsedSeconds float32

	for !isFinishConditionMet(iterationsWithoutImprovement) {
		for i := 0; i < len(population); i++ {
			x, y, z := generateRandomNumbersForDE(len(population), i)
			// fmt.Printf("Population %d: %v\n", i, population[i])
			// fmt.Printf("Population length %d\n", len(population))
			// newSpecimen = make([]float32, len(population[i]))
			// oldSpecimen = make([]float32, len(population[i]))
			// newSpecimenBeforeRepair = make([]float32, len(population[i]))

			// copy(oldSpecimen, population[i])
			newSpecimen = recomb(population[i], population[x], population[y], population[z], differentialWeight, crossoverProbability)
			// copy(newSpecimenBeforeRepair, newSpecimen)
			sc.FixEntities(newSpecimen, i)
			newSpecimenIncome := sc.CalculateIncome(newSpecimen, i)
			if newSpecimenIncome > incomes[i] {
				copy(population[i], newSpecimen)
				incomes[i] = newSpecimenIncome

				// Log to csv in the future
				if currentBestScore < newSpecimenIncome {
					// fmt.Printf("Got, ya new best score: %v,\nOld Specimen:\t%v\nSpecimen: \t%v\nBefore repair\t%v\n\n", newSpecimenIncome, oldSpecimen, newSpecimen, newSpecimenBeforeRepair)
					// fmt.Printf("Got, ya new best score: %v\n", newSpecimenIncome)
					bestSpecimen = make([]float32, len(population[i]))
					copy(bestSpecimen, newSpecimen)
					currentBestScore = newSpecimenIncome
				}
			}
		}

		if previousBestScore < currentBestScore {
			// Log this to csv file
			elapsedSeconds = getElapsedTime()
			copy(buf, fmt.Sprintf("%f,%f,%d,%d\n", elapsedSeconds, currentBestScore, totalIterations, iterationsWithoutImprovement))
			w.Write(buf)
			iterationsWithoutImprovement = -1
			previousBestScore = currentBestScore
		}
		totalIterations++
		iterationsWithoutImprovement++
	}

	fmt.Printf("Best Specimen found: %v", bestSpecimen)
	fmt.Printf("With the profit of: %v", currentBestScore)

}

func recomb(x []float32, s1 []float32, s2 []float32, s3 []float32, differentialWeight float32, crossoverProbability float32) []float32 {
	afterMutation := make([]float32, len(s1))
	for i := 0; i < len(afterMutation); i++ {
		if randomNumber := rand.Float32(); randomNumber < crossoverProbability {
			afterMutation[i] = s1[i] + (s2[i]-s3[i])*differentialWeight
			if afterMutation[i] == 0 && randomNumber < (1-crossoverProbability)*(1-differentialWeight)*0.01 {
				afterMutation[i] = utils.RandFloatFromRange(1, 5)
			}
		} else {
			afterMutation[i] = x[i]
		}
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

func getTimeElapsedCounter(startTime time.Time) func() float32 {
	s := startTime

	return func() float32 {
		return float32(math.Floor(float64(time.Since(s))) / float64(time.Second))
	}
}

func generateRandomNumbersForDE(populationSize int, currentIndex int) (int, int, int) {
	areResultValid := false
	var x, y, z int
	for !areResultValid {
		x = random.Random(0, populationSize-1)
		y = random.Random(0, populationSize-1)
		z = random.Random(0, populationSize-1)

		if currentIndex != x && currentIndex != y && currentIndex != z && x != z && x != y && y != x {
			areResultValid = true
		}
	}
	return x, y, z
}
