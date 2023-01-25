package utils

import (
	"math/rand"

	"github.com/BlueberryBuns/MSCN-optimization/entities"
)

var entitiesMapping = map[int]string{
	0: "supplier",
	1: "factory",
	2: "warehouse",
	3: "shop",
}

func Sum[T int | float32](s []T) T {
	var sum T = 0
	for _, elem := range s {
		sum += elem
	}
	return sum
}

func Shuffle[T any](s []T) []T {
	rand.Shuffle(len(s), func(i, j int) { s[i], s[j] = s[j], s[i] })
	return s
}

func RandFloatFromRange(min, max float32) float32 {
	return min + rand.Float32()*(max-min)
}

func slice_assign(s1 []float32, s2 []float32, start_at int) []float32 {
	for i := 0; i < len(s2); i++ {
		s1[i+start_at] = s2[i]
	}
	return s1
}

func CopyEntities(src []entities.IBaseEntity) []entities.IBaseEntity {
	newEntities := make([]entities.IBaseEntity, len(src))
	var err error = nil

	for index := 0; index < len(src); index++ {
		newEntities[index], err = entities.EntityFactory(entitiesMapping[src[index].GetEntityType()], src[index].GetIndex(), src[index].GetMaxCapacity(), src[index].GetSetupCost())
		if err != nil {
		}
	}

	return newEntities
}
