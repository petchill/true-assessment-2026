package aggregation

import (
	"recommendation-system/src/internal/model/entity"
	"testing"
	"time"

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
		assert.EqualValues(t, tc.expectedResult, actual)
	}
}

func Test_CalculateContentScore(t *testing.T) {
	calculator := ScoreCalculator{}
	content := entity.Content{
		Genre:           "action",
		PopularityScore: 0.8,
		CreatedAt:       time.Now().AddDate(0, 0, -30),
	}
	preferenceMap := PreferenceMap{
		"action": 0.6,
	}

	baseScore := (content.PopularityScore * 0.4) + (preferenceMap[content.Genre] * 0.35) + (content.GetRecencyFactory() * 0.15)
	actual := calculator.CalculateContentScore(content, preferenceMap)

	assert.GreaterOrEqual(t, actual, baseScore-0.05)
	assert.LessOrEqual(t, actual, baseScore+0.05)
}
