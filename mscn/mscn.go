package mscn

import (
	"fmt"
	"io/ioutil"
	"log"
	"time"

	"github.com/BlueberryBuns/MSCN-optimization/entities"
	"gopkg.in/yaml.v3"
)

type MSCN struct {
	DateStarted          string
	SuppliersStartIndex  int
	FactoriesStartIndex  int
	WarehousesStartIndex int
	ShopsStartIndex      int
	Suppliers            []entities.IBaseEntity
	Factories            []entities.IBaseEntity
	Warehouses           []entities.IBaseEntity
	Shops                []entities.IBaseEntity
}

type MSCNStructureInput struct {
	SupplierCount                      int       `yaml:"D"`
	FactoryCount                       int       `yaml:"F"`
	WarehousesCount                    int       `yaml:"M"`
	ShopsCount                         int       `yaml:"S"`
	SupplierMaxCapacity                []float32 `yaml:"sd"`
	FactoryMaxCapacity                 []float32 `yaml:"sf"`
	WarehouseMaxCapacity               []float32 `yaml:"sm"`
	ShopMarketDemand                   []float32 `yaml:"ss"` // Max market Demand
	TransportCostSupplierToFactory     []float32 `yaml:"cd"`
	TransportCostFactoryToWarehouse    []float32 `yaml:"cf"`
	TransportCostWarehouseToStore      []float32 `yaml:"cm"`
	SetupCostSupplier                  []float32 `yaml:"ud"`
	SetupCostFactory                   []float32 `yaml:"uf"`
	SetupCostWarehouse                 []float32 `yaml:"um"`
	ShopIncomePerProduct               []float32 `yaml:"p"`
	MinMaxSupplierFactoryProvisioning  []float32 `yaml:"xdminmax"`
	MinMaxFactoryWarehouseProvisioning []float32 `yaml:"xfminmax"`
}

func GenerateMscnStructureFromFile(filename string) (MSCN, error) {
	s := readStructureFromFile(filename)
	suppliers, err := entities.EntitiesFactory("supplier", s.SupplierCount, s.SupplierMaxCapacity, s.SetupCostSupplier)

	if err != nil {
		return MSCN{}, fmt.Errorf("%v", err)
	}

	factories, err := entities.EntitiesFactory("factory", s.FactoryCount, s.FactoryMaxCapacity, s.SetupCostFactory)

	if err != nil {
		return MSCN{}, fmt.Errorf("%v", err)
	}

	warehouses, err := entities.EntitiesFactory("warehouse", s.WarehousesCount, s.WarehouseMaxCapacity, s.SetupCostWarehouse)

	if err != nil {
		return MSCN{}, fmt.Errorf("%v", err)
	}

	shops, err := entities.EntitiesFactory("shop", s.ShopsCount, s.ShopMarketDemand, s.ShopIncomePerProduct)

	if err != nil {
		return MSCN{}, fmt.Errorf("%v", err)
	}

	formatedDate := time.Now().Format("2006-01-02 15:04:05")
	suppliersStartIndex := 0
	factoriesStartIndex := s.SupplierCount
	warehousesStartIndex := s.FactoryCount + factoriesStartIndex
	shopsStartIndex := s.WarehousesCount + warehousesStartIndex

	return MSCN{
			formatedDate,
			suppliersStartIndex,
			factoriesStartIndex,
			warehousesStartIndex,
			shopsStartIndex,
			suppliers,
			factories,
			warehouses,
			shops},
		nil
}

func readStructureFromFile(filename string) MSCNStructureInput {
	var MSCNStruct = &MSCNStructureInput{}
	yaml_file, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Panic(err)
	}

	parsing_error := yaml.Unmarshal(yaml_file, MSCNStruct)
	if parsing_error != nil {
		log.Fatal(err)
	}
	return *MSCNStruct
}
