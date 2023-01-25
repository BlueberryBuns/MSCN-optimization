package solution

import (
	"fmt"

	"github.com/BlueberryBuns/MSCN-optimization/entities"
	"github.com/BlueberryBuns/MSCN-optimization/mscn"
	"github.com/BlueberryBuns/MSCN-optimization/utils"
)

type SolutionCalculator struct {
	*mscn.MSCN
	suppliers      [][]entities.IBaseEntity
	factories      [][]entities.IBaseEntity
	warehouses     [][]entities.IBaseEntity
	shops          [][]entities.IBaseEntity
	populationSize int
}

type costCalculator func(partialSolution []float32, hoe entities.IBaseEntity, loe entities.IBaseEntity, loeLength int, transportCost []float32) float32

func CreateSolutionCalculator(problemInstance *mscn.MSCN, populationSize int) SolutionCalculator {
	suppliers := make([][]entities.IBaseEntity, populationSize)
	factories := make([][]entities.IBaseEntity, populationSize)
	warehouses := make([][]entities.IBaseEntity, populationSize)
	shops := make([][]entities.IBaseEntity, populationSize)

	for i := 0; i < populationSize; i++ {
		suppliers[i] = make([]entities.IBaseEntity, problemInstance.SupplierCount)
		factories[i] = make([]entities.IBaseEntity, problemInstance.FactoryCount)
		warehouses[i] = make([]entities.IBaseEntity, problemInstance.WarehousesCount)
		shops[i] = make([]entities.IBaseEntity, problemInstance.ShopsCount)

		copy(suppliers[i], utils.CopyEntities(problemInstance.Suppliers))
		copy(factories[i], utils.CopyEntities(problemInstance.Factories))
		copy(warehouses[i], utils.CopyEntities(problemInstance.Warehouses))
		copy(shops[i], utils.CopyEntities(problemInstance.Shops))
	}

	sc := SolutionCalculator{problemInstance, suppliers, factories, warehouses, shops, populationSize}

	return sc
}

func (sc *SolutionCalculator) calculateTotalTransportationCost(solution []float32, populationIndex int) float32 {
	var totalTransportationCost float32

	first_index := sc.SuppliersStartIndex
	last_index := sc.FactoryCount * sc.SupplierCount
	totalTransportationCost += sc.calculateCost(solution[first_index:last_index], sc.suppliers[populationIndex], sc.factories[populationIndex], singleTransportCost, sc.TransportCostSupplierToFactory)
	first_index = last_index
	last_index = last_index + sc.FactoryCount*sc.WarehousesCount
	totalTransportationCost += sc.calculateCost(solution[first_index:last_index], sc.factories[populationIndex], sc.warehouses[populationIndex], singleTransportCost, sc.TransportCostFactoryToWarehouse)
	first_index = last_index
	last_index = last_index + sc.WarehousesCount*sc.ShopsCount
	totalTransportationCost += sc.calculateCost(solution[first_index:last_index], sc.warehouses[populationIndex], sc.shops[populationIndex], singleTransportCost, sc.TransportCostWarehouseToStore)
	fmt.Printf("Total Transport Cost:\t%.4f\n", totalTransportationCost)
	return totalTransportationCost
}

func (sc *SolutionCalculator) calculateTotalContractCost(solution []float32, populationIndex int) float32 {
	var totalContractCost float32

	first_index := sc.SuppliersStartIndex
	last_index := sc.FactoryCount * sc.SupplierCount
	totalContractCost += sc.calculateCost(solution[first_index:last_index], sc.suppliers[populationIndex], sc.factories[populationIndex], singleContractCost, sc.TransportCostSupplierToFactory)
	first_index = last_index
	last_index = last_index + sc.FactoryCount*sc.WarehousesCount
	totalContractCost += sc.calculateCost(solution[first_index:last_index], sc.factories[populationIndex], sc.warehouses[populationIndex], singleContractCost, sc.TransportCostFactoryToWarehouse)
	first_index = last_index
	last_index = last_index + sc.WarehousesCount*sc.ShopsCount
	totalContractCost += sc.calculateCost(solution[first_index:last_index], sc.warehouses[populationIndex], sc.shops[populationIndex], singleContractCost, sc.TransportCostWarehouseToStore)
	fmt.Printf("Total Contract Cost:\t%.4f\n", totalContractCost)
	return totalContractCost
}

func (sc *SolutionCalculator) calculateCost(partialSolution []float32, higherOrderEntities []entities.IBaseEntity, lowerOrderEntities []entities.IBaseEntity, callable costCalculator, transportCost []float32) float32 {
	var totalCost float32
	for _, hoe := range higherOrderEntities {
		for _, loe := range lowerOrderEntities {
			totalCost += callable(partialSolution, hoe, loe, len(lowerOrderEntities), transportCost)
		}
	}

	return totalCost
}

func singleContractCost(partialSolution []float32, hoe entities.IBaseEntity, loe entities.IBaseEntity, loeLength int, transportCost []float32) float32 {
	index := hoe.GetIndex()*loeLength + loe.GetIndex()
	if connectionValue := partialSolution[index]; connectionValue != 0 {
		return hoe.GetSetupCost()
	}

	return 0.0
}

func singleTransportCost(partialSolution []float32, hoe entities.IBaseEntity, loe entities.IBaseEntity, loeLength int, transportCost []float32) float32 {
	index := hoe.GetIndex()*loeLength + loe.GetIndex()
	// fmt.Printf("index: %v, \ttransport cost: %v, \t partialSolution: %v\n", index, transportCost, partialSolution)
	return partialSolution[index] * transportCost[index]
}

func (sc *SolutionCalculator) calculateProfit(shops []entities.IBaseEntity) float32 {
	var totalIncome float32

	for _, shop := range shops {
		totalIncome += (shop.GetCurrentCapacity() * shop.GetSetupCost())
	}

	return totalIncome
}

func (sc *SolutionCalculator) CalculatePopulationIncome(population [][]float32) []float32 {
	incomes := make([]float32, len(population))
	for solutionIndex, solution := range population {
		incomes[solutionIndex] = sc.CalculateIncome(solution, solutionIndex)
	}

	return incomes
}

func (sc *SolutionCalculator) CalculateIncome(solution []float32, populationIndex int) float32 {
	transport_cost := sc.calculateTotalTransportationCost(solution, populationIndex)
	contract_cost := sc.calculateTotalContractCost(solution, populationIndex)
	profit := sc.calculateProfit(sc.shops[populationIndex])
	return profit - contract_cost - transport_cost
}

func (sc *SolutionCalculator) GenerateRandomSolution(populationIndex int) []float32 {
	fmt.Printf("%v\n", sc.DateStarted)
	solution := sc.generateSupplierConnections(populationIndex)
	return solution
}

func (sc *SolutionCalculator) GeneratePopulation() [][]float32 {
	population := make([][]float32, 0)

	for i := 0; i < sc.populationSize; i++ {
		specimen := sc.GenerateRandomSolution(i)
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
			// fmt.Printf("Saving at: %v, length of partial_solution %v\n", (hoe.GetIndex()*len(lowerOrderEntities) + loe.GetIndex()), len(partial_solution))
		}
	}
	// fmt.Printf("%v\n", partial_solution)
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

	// fmt.Printf("set new connection value: %v, %v --->>> %v\t\t new %s capacity: %v\t\t Saving at: %v\t", min, max, connection_value, hoe.GetEncodedRepresentation(), hoe.GetCurrentCapacity(), loe.GetIndex()+hoe.GetIndex()*lowerOrderEntitiesLength)
	return connection_value
}

func (sc *SolutionCalculator) generateSupplierConnections(populationIndex int) []float32 {
	solution := make([]float32, 0)

	_suppliers := sc.suppliers[populationIndex]
	_factories := sc.factories[populationIndex]
	_warehouses := sc.warehouses[populationIndex]
	_shops := sc.shops[populationIndex]

	for i, supplier := range sc.Suppliers {
		fmt.Printf("Supplier %d after population %d: %v\n", i, populationIndex, supplier.GetCurrentCapacity())
	}

	partialSolution := sc.iterateOverConstraints(_suppliers, _factories, sc.MinMaxSupplierFactoryProvisioning)
	solution = append(solution, partialSolution...)

	// fmt.Printf("%v\n", solution)

	// fmt.Println("Factories")
	partialSolution = sc.iterateOverConstraints(_factories, _warehouses, sc.MinMaxFactoryWarehouseProvisioning)
	solution = append(solution, partialSolution...)
	// fmt.Printf("%v\n", solution)

	// fmt.Println("Warehouses")
	partialSolution = sc.iterateOverConstraints(_warehouses, _shops, sc.MinMaxWarehouseShopProvisioning)
	solution = append(solution, partialSolution...)
	// fmt.Printf("%v\n", solution)

	return solution
}

func (sc *SolutionCalculator) DisplayPopulation(population [][]float32) {
	fmt.Printf("Length of array %d\n", len(population))
	for idx, row := range population {
		fmt.Printf("%d: %v\n", idx, row)
	}
}
