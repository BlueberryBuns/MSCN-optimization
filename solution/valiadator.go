package solution

import (
	"github.com/BlueberryBuns/MSCN-optimization/entities"
	"github.com/BlueberryBuns/MSCN-optimization/mscn"
	"github.com/BlueberryBuns/MSCN-optimization/utils"
)

type Validator struct {
	mscn *mscn.MSCN
}

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
	for supplierIndex, supplier := range v.mscn.Suppliers {
		startIndex := supplierIndex * len(v.mscn.Factories)
		endIndex := supplierIndex*len(v.mscn.Factories) + len(v.mscn.Suppliers)
		// for conn := supplierIndex * len(v.mscn.Factories); conn < v.mscn.FactoriesStartIndex * supplierIndex; conn++ {
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
