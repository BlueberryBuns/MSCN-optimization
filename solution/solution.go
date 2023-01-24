package solution

import (
	"fmt"

	"github.com/BlueberryBuns/MSCN-optimization/entities"
	"github.com/BlueberryBuns/MSCN-optimization/mscn"
	"github.com/BlueberryBuns/MSCN-optimization/utils"
)

type SolutionCalculator struct {
	*mscn.MSCN
}

func CreateSolutionCalculator(problemInstance *mscn.MSCN) SolutionCalculator {
	sc := SolutionCalculator{problemInstance}
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
	contract_cost := sc.calculateContractCost(solution)
	profit := sc.calculateProfit(solution)
	return profit - contract_cost - transport_cost
}

func (sc *SolutionCalculator) GenerateRandomSolution() []float32 {
	fmt.Printf("%v\n", sc.DateStarted)
	solution := sc.generateSupplierConnections()
	return solution
}

func (sc *SolutionCalculator) GeneratePopulation(population_size int) [][]float32 {
	population := make([][]float32, 0)

	for i := 0; i < population_size; i++ {
		specimen := sc.GenerateRandomSolution() //:= make([]float32, sc.SupplierCount*sc.FactoryCount+sc.FactoryCount*sc.WarehousesCount+sc.WarehousesCount*sc.ShopsCount)
		population = append(population, [][]float32{specimen}...)
	}

	return population
}

type SolutionGenerator struct {
	mscnInstance *mscn.MSCN
	validator    Validator
}

func (sc *SolutionCalculator) iterateOverConstraints(higherOrderEntities []entities.IBaseEntity, lowerOrderEntities []entities.IBaseEntity, provisioning []float32) []float32 {
	partial_solution := make([]float32, len(lowerOrderEntities)*len(higherOrderEntities))
	for _, hoe := range higherOrderEntities {
		for _, loe := range lowerOrderEntities {
			partial_solution[(hoe.GetIndex()*len(lowerOrderEntities) + loe.GetIndex())] = sc.InitializeConnection(hoe, loe, provisioning, len(lowerOrderEntities))
			fmt.Printf("Saving at: %v, length of partial_solution %v\n", (hoe.GetIndex()*len(lowerOrderEntities) + loe.GetIndex()), len(partial_solution))
		}
	}
	fmt.Printf("%v\n", partial_solution)
	return partial_solution
}

func (sc *SolutionCalculator) InitializeConnection(hoe entities.IBaseEntity, loe entities.IBaseEntity, provisioning []float32, lowerOrderEntitiesLength int) float32 {
	min_index := (loe.GetIndex() + hoe.GetIndex()*lowerOrderEntitiesLength) * 2
	max_index := min_index + 1
	min := provisioning[min_index]
	max := provisioning[max_index]
	connection_value := utils.RandFloatFromRange(min, max)

	new_capacity := hoe.GetCurrentCapacity() - connection_value
	hoe.SetCurrentCapacity(new_capacity)
	loe.SetCurrentCapacity(connection_value)

	fmt.Printf("set new connection value: %v, %v --->>> %v\t\t new %s capacity: %v\t\t Saving at: %v\t", min, max, connection_value, hoe.GetEncodedRepresentation(), hoe.GetCurrentCapacity(), loe.GetIndex()+hoe.GetIndex()*lowerOrderEntitiesLength)
	return connection_value
}

func (sc *SolutionCalculator) generateSupplierConnections() []float32 {
	solution := make([]float32, 0)

	fmt.Println("Suppliers")
	partial_solution := sc.iterateOverConstraints(sc.Suppliers, sc.Factories, sc.MinMaxSupplierFactoryProvisioning)
	solution = append(solution, partial_solution...)
	fmt.Printf("%v\n", solution)

	fmt.Println("Factories")
	partial_solution = sc.iterateOverConstraints(sc.Factories, sc.Warehouses, sc.MinMaxFactoryWarehouseProvisioning)
	solution = append(solution, partial_solution...)
	fmt.Printf("%v\n", solution)

	fmt.Println("Warehouses")
	partial_solution = sc.iterateOverConstraints(sc.Warehouses, sc.Shops, sc.MinMaxWarehouseShopProvisioning)
	solution = append(solution, partial_solution...)
	fmt.Printf("%v\n", solution)

	return solution
}

func (sc *SolutionCalculator) DisplayPopulation(population [][]float32) {
	fmt.Printf("Length of array %d\n", len(population))
	for idx, row := range population {
		fmt.Printf("%d: %v\n", idx, row)
	}
}

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
