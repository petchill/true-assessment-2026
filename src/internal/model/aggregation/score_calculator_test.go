package aggregation

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_calculateGenrePreferenceWeight(t *testing.T) {
	calculator := ScoreCalculator{}
	type TestCase struct {
		name           string
		param          map[string]int
		expectedResult map[string]float64
	}

	testCases := []TestCase{
		{
			name: "Given only one item, then preference should be 1",
			param: map[string]int{
				"thriller": 10,
			},
			expectedResult: map[string]float64{
				"thriller": 1,
			},
		},
		{
			name: "Given multiple items, then preference should be divided",
			param: map[string]int{
				"thriller": 10,
				"comedy":   30,
			},
			expectedResult: map[string]float64{
				"thriller": 0.25,
				"comedy":   0.75,
			},
		},
	}

	for _, tc := range testCases {
		actual := calculator.calculateGenrePreferenceWeight(tc.param)
		assert.Equal(t, tc.expectedResult, actual)
	}
}
