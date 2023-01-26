package solution

import (
	"github.com/BlueberryBuns/MSCN-optimization/entities"
	"github.com/BlueberryBuns/MSCN-optimization/mscn"
)

type Validator struct {
	*mscn.MSCN
	decrementalValue float32
}

func (v *Validator) FixEntities(solution []float32, hoEntities []entities.IBaseEntity, loEntities []entities.IBaseEntity, provisioning []float32) {
	// lowerOrderEntities := utils.Shuffle(loEntities)
	v.fixOverflows(solution, hoEntities, loEntities)
	v.fixOutOfBoundConnections(solution, hoEntities, loEntities, provisioning)

}

// func updateSolution(solution []float32)
func (v *Validator) fixOutOfBoundConnections(solution []float32, hoEntities []entities.IBaseEntity, loEntities []entities.IBaseEntity, provisioning []float32) {
	for _, e := range hoEntities {
		for loeIndex, globalIndex := range e.GetGlobalOutIndexes() {

			minIndex := (e.GetIndex()*len(loEntities) + loeIndex) * 2
			maxIndex := minIndex + 1

			if solution[globalIndex] < provisioning[minIndex] {
				// fmt.Printf("Solution: %v exceeds min capacity (%v), at index %v\n", solution[globalIndex], provisioning[minIndex], minIndex)
				solution[globalIndex] = 0.0
				e.UpdateCapacityOut(solution)
				loEntities[loeIndex].UpdateCapacityIn(solution)
				loEntities[loeIndex].UpdateCapacityOut(solution)
			}

			if solution[globalIndex] > provisioning[maxIndex] {
				solution[globalIndex] = provisioning[maxIndex]
				e.UpdateCapacityOut(solution)
				loEntities[loeIndex].UpdateCapacityIn(solution)
				loEntities[loeIndex].UpdateCapacityOut(solution)
			}
		}
	}
	v.updateCapacitiesIn(solution, hoEntities)
	v.updateCapacitiesIn(solution, loEntities)
	v.updateCapacitiesOut(solution, hoEntities)
	v.updateCapacitiesOut(solution, loEntities)
	// fmt.Printf("Provisioning: %v\n", provisioning)
}

func (v *Validator) fixOverflows(solution []float32, hoEntities []entities.IBaseEntity, loEntities []entities.IBaseEntity) {
	v.updateCapacitiesIn(solution, hoEntities)
	v.updateCapacitiesIn(solution, loEntities)
	v.updateCapacitiesOut(solution, hoEntities)
	v.updateCapacitiesOut(solution, loEntities)
	var nullify bool
	for _, h := range hoEntities {
		for (h.GetCapacityIn() - h.GetCapacityOut()) < 0.0 {
			nullify = false
			if h.GetCapacityIn() == 0 {
				nullify = true
			}

			// fmt.Printf("Current capacity %v, %v\n", h.GetCapacityIn(), h.GetCapacityOut())
			for loIndex, globalIndex := range h.GetGlobalOutIndexes() {
				if nullify {
					// fmt.Printf(">>>>>>>>>Before %v\n", solution)
					solution[globalIndex] = 0
					h.UpdateCapacityOut(solution)
					loEntities[loIndex].UpdateCapacityIn(solution)
					// fmt.Printf(">>>>>>>>>After %v\n", solution)
					continue
				}

				solution[globalIndex] = solution[globalIndex] * 0.9
				h.UpdateCapacityOut(solution)
				loEntities[loIndex].UpdateCapacityIn(solution)
				// fmt.Printf("solution %v\n", solution[loIndex])
			}
		}
		// fmt.Printf("Redundancy in capacity: %v\n", h.GetCapacityIn()-h.GetCapacityOut())
	}
}

func (v *Validator) updateCapacitiesIn(solution []float32, eSlice []entities.IBaseEntity) {
	if eSlice[0].GetEntityType() == entities.SupplierType {
		return
	}

	for _, e := range eSlice {
		e.UpdateCapacityIn(solution)
	}
}

func (v *Validator) updateCapacitiesOut(solution []float32, eSlice []entities.IBaseEntity) {
	if eSlice[0].GetEntityType() == entities.ShopType {
		return
	}

	for _, e := range eSlice {
		e.UpdateCapacityOut(solution)
	}
}
