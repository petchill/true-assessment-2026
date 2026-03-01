package aggregation

import (
	"math/rand/v2"
	"recommendation-system/src/internal/model/entity"
)

type ScoreCalculator struct {
}

func (c *ScoreCalculator) CalculateContentScore(content entity.Content, preferenceMap PreferenceMap) float64 {
	popularityComponent := content.PopularityScore * 0.4
	genreBoost := preferenceMap[content.Genre] * 0.35
	recencyComponent := content.GetRecencyFactory() * 0.15
	randomNoise := (rand.Float64() - 0.5) * 0.1
	return popularityComponent + genreBoost + recencyComponent + randomNoise
}
