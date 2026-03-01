package aggregation

type ScoreCalculator struct {
}

// func (c *ScoreCalculator) CalculateScore()

func (c *ScoreCalculator) calculateGenrePreferenceWeight(mapCount map[string]int) map[string]float64 {
	sum := 0
	result := make(map[string]float64)
	for _, v := range mapCount {
		sum += v
	}

	for k, v := range mapCount {
		result[k] = float64(v) / float64(sum)
	}
	return result
}
