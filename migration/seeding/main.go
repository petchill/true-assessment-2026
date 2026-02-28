package main

import (
	"log"
	"math/rand"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func randomIntRange(min, max int) int {
	return min + rand.Intn(max-min+1)
}

func selectWeighted(items []string, weights []int) string {
	var totalWeight int
	for _, w := range weights {
		totalWeight += w
	}

	randomVal := rand.Intn(totalWeight) // rand.IntN generates a number in [0, totalWeight)

	var cursor int
	for i, w := range weights {
		cursor += w
		if cursor > randomVal {
			return items[i]
		}
	}

	// Should not be reached if totalWeight > 0
	return ""
}

func main() {
	databaseURL := os.Getenv("DATABASE_URL")

	db, err := gorm.Open(postgres.Open(databaseURL), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	if err := db.Exec("TRUNCATE TABLE user_watch_history, content, users RESTART IDENTITY CASCADE").Error; err != nil {
		log.Fatalf("truncate tables: %v", err)
	}

	users := SeedingUsers(db)
	contents := SeedingContents(db)
	SeedingUserWatchHistory(db, users, contents)
}
