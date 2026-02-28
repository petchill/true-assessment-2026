package main

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type contentSeed struct {
	ID              uint64    `gorm:"column:id;primaryKey"`
	Title           string    `gorm:"column:title"`
	Genre           string    `gorm:"column:genre"`
	PopularityScore float64   `gorm:"column:popularity_score"`
	CreatedAt       time.Time `gorm:"column:created_at"`
}

func (contentSeed) TableName() string {
	return "content"
}

func SeedingContents(db *gorm.DB) []contentSeed {
	n := 50

	adjectives := []string{"Hidden", "Golden", "Silent", "Neon", "Last", "Broken", "Midnight", "Wild", "Crystal", "Fading"}
	nouns := []string{"Signal", "Empire", "Journey", "Memory", "Storm", "Garden", "Protocol", "Shadow", "Frontier", "Pulse", "หอแต๋วแตก"}
	genres := []string{"action", "drama", "comedy", "sci-fi", "thriller", "documentary"}

	contents := []contentSeed{}
	now := time.Now().UTC()
	for i := 1; i <= n; i++ {
		title := fmt.Sprintf("%s %s", adjectives[rand.Intn(len(adjectives))], nouns[rand.Intn(len(nouns))])
		genre := genres[rand.Intn(len(genres))]
		popularityScore := float64(randomIntRange(50, 100)) / 100
		dateSub := randomIntRange(-30, -1)

		contents = append(contents, contentSeed{
			ID:              uint64(i),
			Title:           fmt.Sprintf("%s %d", title, i),
			Genre:           genre,
			PopularityScore: popularityScore,
			CreatedAt:       now.AddDate(0, 0, dateSub),
		})
	}

	if err := db.Clauses(clause.OnConflict{DoNothing: true}).Create(&contents).Error; err != nil {
		log.Fatalf("seed content: %v", err)
	}

	log.Printf("seeded %d contents", len(contents))
	return contents
}
