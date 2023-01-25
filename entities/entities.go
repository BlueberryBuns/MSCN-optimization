package entities

import (
	"fmt"
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
	UpdateCapacityIn(c float32)
	UpdateCapacityOut(c float32)
	GetEntityType() int
	GetCapacityOut() float32
	GetCapacityIn() float32
	UpdateGlobalInIndexes(inIndex int)
	UpdateGlobalOutIndexes(outIndex int)
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

func (b *BaseEntity) UpdateCapacityIn(c float32) {
	b.capacityIn += c
}

func (b *BaseEntity) GetCapacityOut() float32 {
	return b.capacityOut
}

func (b *BaseEntity) UpdateCapacityOut(c float32) {
	b.capacityOut += c
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

func (b *BaseEntity) UpdateGlobalInIndexes(inIndex int) {
	b.globalInIndexes = append(b.globalInIndexes, inIndex)
}
func (b *BaseEntity) GetGlobalInIndexes() []int {
	return b.globalInIndexes
}
func (b *BaseEntity) UpdateGlobalOutIndexes(outIndex int) {
	b.globalOutIndexes = append(b.globalOutIndexes, outIndex)
}
func (b *BaseEntity) GetGlobalOutIndexes() []int {
	return b.globalOutIndexes
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
