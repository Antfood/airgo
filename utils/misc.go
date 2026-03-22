package utils

/*
Map transforms a slice of type T to a slice of type R using the provided function.

Example:

	nums := []int{1, 2, 3}
	doubled := Map(nums, func(n int) int { return n * 2 })

	// Output: doubled = [2, 4, 6]
*/
func Map[T any, R any](slice []T, fn func(T) R) []R {
	result := make([]R, len(slice))

	for i, val := range slice {
		result[i] = fn(val)
	}

	return result
}
