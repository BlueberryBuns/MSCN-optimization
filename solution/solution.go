package solution

import (
	"fmt"

	"github.com/BlueberryBuns/MSCN-optimization/entities"
	"github.com/BlueberryBuns/MSCN-optimization/mscn"
	"github.com/BlueberryBuns/MSCN-optimization/utils"
)

type SolutionCalculator struct {
	problemInstance *mscn.MSCN
}

func CreateSolutionCalculator(problemInstance *mscn.MSCN) SolutionCalculator {
	sc := SolutionCalculator{problemInstance: problemInstance}
	fmt.Printf("%v\n", sc.problemInstance.SuppliersStartIndex)
	fmt.Printf("%v\n", sc.problemInstance.FactoriesStartIndex)
	fmt.Printf("%v\n", sc.problemInstance.WarehousesStartIndex)
	fmt.Printf("%v\n", sc.problemInstance.ShopsStartIndex)
	fmt.Printf("XDDDD,%v\n", sc.problemInstance.ShopsStartIndex)

	return sc
}

func (sc *SolutionCalculator) calculateTransportationCost(solution []float32) float32 {
	return 0.0
}

func (sc *SolutionCalculator) calculateContractCost(solution []float32) float32 {
	return 0.0
}

func (sc *SolutionCalculator) calculateProfit(solution []float32) float32 {
	return 0.0
}

func (sc *SolutionCalculator) CalculateIncome(solution []float32) float32 {
	transport_cost := sc.calculateContractCost(solution)
	profit := sc.calculateProfit(solution)
	contract_cost := sc.calculateContractCost(solution)
	return profit - contract_cost - transport_cost
}

func (sc *SolutionCalculator) GenerateRandomSolution() []float32 {
	fmt.Printf("%v\n", sc.problemInstance.DateStarted)
	return []float32{1.1, 1.2}
}

type SolutionGenerator struct {
	mscnInstance *mscn.MSCN
	validator    Validator
}

func generateSupplierConnections() {

}

func generateFactoryConnections() {

}

func generateWarehouseConnections() {

}

func generateShopConnections() {

}

func (sg *SolutionGenerator) generateInitialEntityConnection(entitiesA []entities.IBaseEntity, entitiesB []entities.IBaseEntity, startingIndex int) []float32 {
	var partialSolution []float32
	partiallyValidate := sg.validator.partiallyValidateEntity()
	for indexA, entityA := range entitiesA {
		entitiesB := utils.Shuffle(entitiesB)
		for indexB, entityB := range entitiesB {
			i := 0
			checkedIndex := startingIndex + indexB
			baseIndex := indexA * len(entitiesA)
			start, stop := startingIndex+baseIndex, startingIndex+baseIndex
			for i < 5 {
				partialSolution[checkedIndex+baseIndex] = utils.RandFloatFromRange(0, entityB.GetMaxCapacity())
				if _, ok := partiallyValidate(partialSolution[start:stop], entityA, checkedIndex); ok {
					// @TODO: zrobić coś z tym...
				}

			}
		}
	}
	return partialSolution

}

// func generate_random_solution(mscn_instance mscn.MSCNStructureInput) Solution {
// 	// var supplier_factory_connections_count int = 1
// 	return Solution{}
// }

/*
	1 Producent
	1 Fabryka
	1 Magazyn
	1 Sklep
	Solution: [1.0, 1.0, 1.0] Xd_df, Xf_fm, Xm_ms

	2 Producent
	1 Fabryka
	2 Magazyn
	1 Sklep
	Solution: [ 0.4, 0.6, 0.1, 0.9, 0.3, 0.7 ] P1-F1(0.4), P2-F2(0.6) F1-M1(0.1) F1-M2(0.9), M1-S1(.3), M2-S1(.7)

	2 Producent
	1 Fabryka
	2 Magazyn
	1 Sklep
	Solution: [ 0.4, 0.6, 0.0, 5.2, 0.3, 0.7 ] P1-F1(0.4), P2-F2(0.6) F1-M1(0.1) F1-M2(0.9), M1-S1(.3), M2-S1(.7)

	P * F - pierwszych połączeń Prodycent - fabryka
	F * M - Liczba drugich fabryka - magazyn
	M * S - Liczba połączeń magazyn - sklep
*/
