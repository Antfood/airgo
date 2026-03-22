package utils

import (
	"testing"

	. "github.com/Antfood/airgo/testutils/testutils"
)

func TestMap(t *testing.T) {
	names := []string{"John", "Mary", "Peter"}

	result := Map(names, func(item string) int {
		return len(item)
	})

	Equals(t, []int{4, 4, 5}, result)
}
