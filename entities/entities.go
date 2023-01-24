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
	SetCurrentCapacity(c float32)
	GetCurrentCapacity() float32
}

type BaseEntity struct {
	entityType      int
	index           int
	maxCapacity     float32
	setupCost       float32
	currentCapacity float32
}

func (b *BaseEntity) GetCurrentCapacity() float32 {
	return b.currentCapacity
}

func (b *BaseEntity) SetCurrentCapacity(c float32) {
	b.currentCapacity = c
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
			entityType:      SupplierType,
			index:           index,
			maxCapacity:     maxCapacity,
			setupCost:       setupCost,
			currentCapacity: maxCapacity,
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

func entityFactory(entityType string, index int, maxCapacity float32, setupCost float32) (IBaseEntity, error) {
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
		entities[index], err = entityFactory(entityType, index, capacities[index], setupCosts[index])
		if err != nil {
			return nil, fmt.Errorf("[Error]: %v while creating entities slice of type %s\n", err, entityType)
		}
	}
	return entities, nil

}
