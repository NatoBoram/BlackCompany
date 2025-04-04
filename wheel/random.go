package wheel

import "math/rand/v2"

// RandomIn returns a random element from a slice.
func RandomIn[T any](a []T) T {
	return a[rand.IntN(len(a))]
}
