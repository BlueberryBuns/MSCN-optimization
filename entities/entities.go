package entities

import (
	"fmt"
	"math"
	"strconv"
)

const (
	SupplierType = iota
	FactoryType
	WarehouseType
	ShopType
)

type IBaseEntity interface {
	GetIndex() int
	GetMaxCapacity() float32
	GetSetupCost() float32
	GetEncodedRepresentation() string
	UpdateCapacityIn(solution []float32)
	UpdateCapacityOut(solution []float32)
	GetEntityType() int
	GetCapacityOut() float32
	GetCapacityIn() float32
	AddGlobalInIndex(inIndex int)
	AddGlobalOutIndex(outIndex int)
	RemoveGlobalOutIndex(removedValue int)
	RemoveGlobalInIndex(removedValue int)
	GetGlobalOutIndexes() []int
	GetGlobalInIndexes() []int
}

type BaseEntity struct {
	entityType       int
	index            int
	maxCapacity      float32
	setupCost        float32
	capacityIn       float32
	capacityOut      float32
	globalInIndexes  []int
	globalOutIndexes []int
}

func (b *BaseEntity) GetEntityType() int {
	return b.entityType
}

func (b *BaseEntity) GetCapacityIn() float32 {
	return b.capacityIn
}

func (b *BaseEntity) UpdateCapacityIn(solution []float32) {
	var totalCapacityIn float32
	for _, globalIndex := range b.globalInIndexes {
		totalCapacityIn += solution[globalIndex]
	}
	b.capacityIn = float32(math.Floor(float64(totalCapacityIn*10e6)) / 10e6)
}

func (b *BaseEntity) GetCapacityOut() float32 {
	return b.capacityOut
}

func (b *BaseEntity) UpdateCapacityOut(solution []float32) {
	var totalCapacityOut float32
	for _, globalIndex := range b.globalOutIndexes {
		totalCapacityOut += solution[globalIndex]
	}
	b.capacityOut = float32(math.Floor(float64(totalCapacityOut*10e6)) / 10e6)
}

func (b *BaseEntity) GetIndex() int {
	return b.index
}

func (b *BaseEntity) GetMaxCapacity() float32 {
	return b.maxCapacity
}

func (b *BaseEntity) GetSetupCost() float32 {
	return b.setupCost
}

func (b *BaseEntity) GetEncodedRepresentation() string {
	return strconv.Itoa(b.entityType) + strconv.Itoa(b.index)
}

func (b *BaseEntity) AddGlobalInIndex(inIndex int) {
	b.globalInIndexes = append(b.globalInIndexes, inIndex)
}

func (b *BaseEntity) AddGlobalOutIndex(outIndex int) {
	b.globalOutIndexes = append(b.globalOutIndexes, outIndex)
}

func (b *BaseEntity) GetGlobalInIndexes() []int {
	return b.globalInIndexes
}

func (b *BaseEntity) GetGlobalOutIndexes() []int {
	return b.globalOutIndexes
}

func (b *BaseEntity) RemoveGlobalOutIndex(removedValue int) {
	newIndexes := make([]int, 0)
	for _, index := range b.globalOutIndexes {
		if index == removedValue {
			continue
		}

		newIndexes = append(newIndexes, index)
	}
	b.globalOutIndexes = newIndexes
}

func (b *BaseEntity) RemoveGlobalInIndex(removedValue int) {
	newIndexes := make([]int, 0)
	for _, index := range b.globalInIndexes {
		if index == removedValue {
			newIndexes = append(newIndexes, -1)
			continue
		}

		newIndexes = append(newIndexes, index)
	}
	b.globalInIndexes = newIndexes
}

type Factory struct {
	BaseEntity
}

type Shop struct {
	BaseEntity
}

type Supplier struct {
	BaseEntity
}

type Warehouse struct {
	BaseEntity
}

func createSupplierEntity(index int, maxCapacity float32, setupCost float32) IBaseEntity {
	return &Supplier{
		BaseEntity: BaseEntity{
			entityType:  SupplierType,
			index:       index,
			maxCapacity: maxCapacity,
			setupCost:   setupCost,
			capacityIn:  maxCapacity,
		},
	}
}

func createFactoryEntity(index int, maxCapacity float32, setupCost float32) IBaseEntity {
	return &Factory{
		BaseEntity: BaseEntity{
			entityType:  FactoryType,
			index:       index,
			maxCapacity: maxCapacity,
			setupCost:   setupCost,
		},
	}
}

func createWarehouseEntity(index int, maxCapacity float32, setupCost float32) IBaseEntity {
	return &Warehouse{
		BaseEntity: BaseEntity{
			entityType:  WarehouseType,
			index:       index,
			maxCapacity: maxCapacity,
			setupCost:   setupCost,
		},
	}
}

func createShopEntity(index int, maxCapacity float32, profitPerProduct float32) IBaseEntity {
	return &Shop{
		BaseEntity: BaseEntity{
			entityType:  ShopType,
			index:       index,
			maxCapacity: maxCapacity,
			setupCost:   profitPerProduct,
		},
	}
}

func EntityFactory(entityType string, index int, maxCapacity float32, setupCost float32) (IBaseEntity, error) {
	switch entityType {
	case "supplier":
		return createSupplierEntity(index, maxCapacity, setupCost), nil
	case "factory":
		return createFactoryEntity(index, maxCapacity, setupCost), nil
	case "warehouse":
		return createWarehouseEntity(index, maxCapacity, setupCost), nil
	case "shop":
		return createShopEntity(index, maxCapacity, setupCost), nil
	}
	return nil, fmt.Errorf("Wrong type of entity was passed")
}

func EntitiesFactory(entityType string, entityCount int, capacities []float32, setupCosts []float32) ([]IBaseEntity, error) {
	entities := make([]IBaseEntity, entityCount)
	var err error
	for index := 0; index < entityCount; index++ {
		entities[index], err = EntityFactory(entityType, index, capacities[index], setupCosts[index])
		if err != nil {
			return nil, fmt.Errorf("[Error]: %v while creating entities slice of type %s\n", err, entityType)
		}
	}
	return entities, nil

}
