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

func slice_assign(s1 []float32, s2 []float32, start_at int) []float32 {
	for i := 0; i < len(s2); i++ {
		s1[i+start_at] = s2[i]
	}
	return s1
}
