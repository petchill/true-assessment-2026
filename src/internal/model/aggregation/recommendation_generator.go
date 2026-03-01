package aggregation

import (
	"context"
	"errors"
	"math/rand/v2"
	"recommendation-system/src/internal/model/entity"
	"sort"
	"time"
)

var ErrModelInference = errors.New("model inference failed")

var (
	simulateModelLatencyFn = simulateModelLatency
	randomFloat64Fn        = rand.Float64
)

type PreferenceMap map[string]float64

type RecommendationGenerator struct {
	scoreCalculator *ScoreCalculator
}

func NewRecommendationGenerator(scoreCalculator *ScoreCalculator) *RecommendationGenerator {
	return &RecommendationGenerator{
		scoreCalculator: scoreCalculator,
	}
}

func (rg *RecommendationGenerator) GenerateRecommendationContents(ctx context.Context, watchHistory []entity.HistoryContent, contents []entity.Content, recommendationLimit int) ([]RecommendationContent, error) {
	genreCounts := rg.countGenres(watchHistory)
	scoreCalculator := ScoreCalculator{}
	genrePreferences := rg.calculateGenrePreferenceWeight(genreCounts)
	recommendationContents := []RecommendationContent{}

	if err := simulateModelLatencyFn(ctx); err != nil {
		return recommendationContents, err
	}

	if randomFloat64Fn() < 0.015 {
		return recommendationContents, ErrModelInference
	}

	recommendations := make([]RecommendationContent, 0, len(contents))
	for _, content := range contents {
		score := scoreCalculator.CalculateContentScore(content, genrePreferences)
		recommendations = append(recommendations, RecommendationContent{
			ContentID:       content.ID,
			Title:           content.Title,
			Genre:           content.Genre,
			PopularityScore: content.PopularityScore,
			Score:           score,
		})
	}

	sort.Slice(recommendations, func(i, j int) bool {
		return recommendations[i].Score > recommendations[j].Score
	})

	if len(recommendations) > recommendationLimit {
		recommendations = recommendations[:recommendationLimit]
	}

	return recommendationContents, nil
}

func (rg *RecommendationGenerator) calculateGenrePreferenceWeight(mapCount map[string]int) PreferenceMap {
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

func (rg *RecommendationGenerator) countGenres(watchHistory []entity.HistoryContent) map[string]int {
	genreCounts := make(map[string]int)
	for _, history := range watchHistory {
		genreCounts[history.Genre]++
	}

	return genreCounts
}

func (rg *RecommendationGenerator) normalizeGenreCounts(genreCounts map[string]int) PreferenceMap {
	total := 0
	for _, count := range genreCounts {
		total += count
	}

	preferenceMap := make(PreferenceMap)
	if total == 0 {
		return preferenceMap
	}

	for genre, count := range genreCounts {
		preferenceMap[genre] = float64(count) / float64(total)
	}

	return preferenceMap
}

func simulateModelLatency(ctx context.Context) error {
	delay := time.Duration(rand.IntN(21)+30) * time.Millisecond
	timer := time.NewTimer(delay)
	defer timer.Stop()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-timer.C:
		return nil
	}
}
