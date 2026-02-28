package main

import (
	"log"
	"math/rand"
	"time"

	"gorm.io/gorm"
)

type userWatchHistorySeed struct {
	ID        uint64    `gorm:"column:id;primaryKey"`
	UserID    uint64    `gorm:"column:user_id"`
	ContentID uint64    `gorm:"column:content_id"`
	WatchedAt time.Time `gorm:"column:watched_at"`
}

func (userWatchHistorySeed) TableName() string {
	return "user_watch_history"
}

func SeedingUserWatchHistory(db *gorm.DB, users []userSeed, contents []contentSeed) {
	n := 200
	histories := make([]userWatchHistorySeed, 0, n)
	now := time.Now().UTC()
	for i := 1; i <= n; i++ {
		user := users[rand.Intn(len(users))]
		content := contents[rand.Intn(len(contents))]
		histories = append(histories, userWatchHistorySeed{
			ID:        uint64(i),
			UserID:    user.ID,
			ContentID: content.ID,
			WatchedAt: now.AddDate(0, 0, -randomIntRange(1, 30)).Add(-time.Duration(rand.Intn(24)) * time.Hour),
		})
	}

	if err := db.Create(&histories).Error; err != nil {
		log.Fatalf("seed user watch history: %v", err)
	}

	log.Printf("seeded %d user watch history rows", len(histories))
}
