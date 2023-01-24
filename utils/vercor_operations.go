package utils

import "math/rand"

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
