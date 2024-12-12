package random

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRandomUniqueInts(t *testing.T) {
	Init()

	t.Run("panics if count is greater than the range", func(t *testing.T) {
		assert.PanicsWithValue(t, "can't generate more unique random integers than the range of possible values", func() {
			RandomInts(1, 5, 5)
		})
	})

	t.Run("returns all possible values when count equals the range", func(t *testing.T) {
		result := RandomInts(0, 100, 100)
		expected := make([]int, 100)
		for i := 0; i < 100; i++ {
			expected[i] = i
		}
		assert.ElementsMatch(t, expected, result)

		fmt.Println(result)
	})
}
