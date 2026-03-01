package aggregation

import (
	"recommendation-system/src/internal/model/entity"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

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
