package solution

import (
	"github.com/BlueberryBuns/MSCN-optimization/entities"
	"github.com/BlueberryBuns/MSCN-optimization/mscn"
	"github.com/BlueberryBuns/MSCN-optimization/utils"
)

type Validator struct {
	*mscn.MSCN
	decrementalValue float32
}

func (v *Validator) FixHOEntities(partialSolution []float32, hoEntities []entities.IBaseEntity, loEntities []entities.IBaseEntity) {
	v.UpdateCapacitiesIn(hoEntities, loEntities, partialSolution)
	v.UpdateCapacitiesOut(hoEntities, loEntities, partialSolution)

	for i := 0; i < len(hoEntities); i++ {
		for hoEntities[i].GetCapacityOut()-hoEntities[i].GetCapacityIn() < 0 {

		}
	}
}

func (v *Validator) FixLOEntities(partialSolution []float32, hoEntities []entities.IBaseEntity, loEntities []entities.IBaseEntity) {
	v.UpdateCapacitiesIn(hoEntities, loEntities, partialSolution)
	v.UpdateCapacitiesOut(hoEntities, loEntities, partialSolution)

	for i := 0; i < len(hoEntities); i++ {
		for hoEntities[i].GetCapacityOut()-hoEntities[i].GetCapacityIn() < 0 {

		}
	}
}

func (v *Validator) UpdateCapacitiesIn(hoes []entities.IBaseEntity, loes []entities.IBaseEntity, partialSolution []float32) {
	if loes[0].GetEntityType() == entities.SupplierType {
		return
	}

	for i := 0; i < len(loes); i++ {
		loes[i].UpdateCapacityIn(calculateCapacityIn(partialSolution, loes[i], hoes, len(loes)))
	}
}

func calculateCapacityIn(partialSolution []float32, loe entities.IBaseEntity, hoes []entities.IBaseEntity, loesLength int) float32 {
	var capacityIn float32
	for _, hoe := range hoes {
		capacityIn += partialSolution[loesLength*hoe.GetIndex()+loe.GetIndex()]
	}

	return capacityIn
}

func (v *Validator) UpdateCapacitiesOut(hoes []entities.IBaseEntity, loes []entities.IBaseEntity, partialSolution []float32) {
	if loes[0].GetEntityType() == entities.SupplierType {
		return
	}

	for i := 0; i < len(loes); i++ {
		loes[i].UpdateCapacityIn(calculateCapacityOut(partialSolution, loes[i], hoes, len(loes)))
	}
}

func calculateCapacityOut(partialSolution []float32, hoe entities.IBaseEntity, loes []entities.IBaseEntity, hoesLength int) float32 {
	var capacityOut float32
	for _, loe := range loes {
		capacityOut += partialSolution[hoesLength*loe.GetIndex()+hoe.GetIndex()]
	}

	return capacityOut
}

// func (v *Validator) calculateCapacityOut(partialSolution []float32, loes []entities.IBaseEntity, hoe entities.IBaseEntity, hoesLength int) float32 {
// 	if loe.GetEntityType() == entities.SupplierType {
// 		return loe.GetCapacityOut()
// 	}

// 	var capacityOut float32
// 	for _, hoe := range loes {
// 		capacityOut += partialSolution[hoesLength*hoe.GetIndex()+hoe.GetIndex()]
// 	}

// 	return capacityOut
// }

// func checkLowerOrderCapacity(solution []float32, loEntities []entities.IBaseEntity, offset int) {
// 	if loEntities == nil {
// 		return
// 	}
// 	for i := 0; i < len(loEntities); i++ {
// 		if loEntities[i].GetCurrentCapacity() > loEntities[i].GetMaxCapacity() {
// 			loEntities[i].SetCurrentCapacity(loEntities[i].GetMaxCapacity())

// 		}
// 	}
// }

// func (v *Validator) fixConnection(solution []float32)

func (v *Validator) partiallyValidateEntity() func(partialSolution []float32, checkedEntity entities.IBaseEntity, checkedIndex int) (int, bool) { // you should pass already a sliced version
	var cache = make(map[string]float32)
	return func(partialSolution []float32, checkedEntity entities.IBaseEntity, checkedIndex int) (int, bool) {
		var currentCapacity float32 = 0
		if ret, ok := cache[checkedEntity.GetEncodedRepresentation()]; ok {
			currentCapacity = ret
		}
		currentCapacity += partialSolution[checkedIndex]
		if currentCapacity > checkedEntity.GetMaxCapacity() {
			return checkedIndex, false
		}
		cache[checkedEntity.GetEncodedRepresentation()] = currentCapacity
		return checkedIndex + 1, true
	}
}

func (v *Validator) isSupplierConnectionValid(solution []float32) bool {
	for supplierIndex, supplier := range v.Suppliers {
		startIndex := supplierIndex * len(v.Factories)
		endIndex := supplierIndex*len(v.Factories) + len(v.Suppliers)
		// for conn := supplierIndex * len(v.Factories); conn < v.FactoriesStartIndex * supplierIndex; conn++ {
		// 	totalSupplierCapacity +=
		// }
		totalSupplierCapacity := utils.Sum(solution[startIndex:endIndex])
		if totalSupplierCapacity > supplier.GetMaxCapacity() {
			return false
		}
	}
	return true
}
func (v *Validator) isFactoryInputConnectionValid(solution []float32) bool {
	return true
}
func (v *Validator) isFactoryOutputConnectionValid(solution []float32) bool {
	return true
}
func (v *Validator) isWarehouseInputConnectionValid(solution []float32) bool {
	return true
}
func (v *Validator) isWarehouseOutputConnectionValid(solution []float32) bool {
	return true
}
func (v *Validator) isShopInputConnectionValid(solution []float32) bool {
	return true
}
func (v *Validator) isShopOutputConnectionValid(solution []float32) bool {
	return true
}
