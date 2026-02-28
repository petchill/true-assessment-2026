package main

import (
	"log"
	"math/rand"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type userSeed struct {
	ID               uint64    `gorm:"column:id;primaryKey"`
	Age              int       `gorm:"column:age"`
	Country          string    `gorm:"column:country"`
	SubscriptionType string    `gorm:"column:subscription_type"`
	CreatedAt        time.Time `gorm:"column:created_at"`
}

func (userSeed) TableName() string {
	return "users"
}

func SeedingUsers(db *gorm.DB) []userSeed {

	n := 20 // number of seeding

	countries := []string{"US", "GB", "CA", "AU", "DE"}
	subscriptionTypes := []string{"free", "basic", "premium"}

	users := []userSeed{}
	now := time.Now().UTC()
	for i := 1; i <= n; i++ {
		age := randomIntRange(18, 65)
		country := countries[rand.Intn(len(countries))]
		subscriptionType := selectWeighted(subscriptionTypes, []int{50, 30, 20})
		dateSub := randomIntRange(-30, -1)
		users = append(users, userSeed{
			ID:               uint64(i),
			Age:              age,
			Country:          country,
			SubscriptionType: subscriptionType,
			CreatedAt:        now.AddDate(0, 0, dateSub),
		})
	}

	if err := db.Clauses(clause.OnConflict{DoNothing: true}).Create(&users).Error; err != nil {
		log.Fatalf("seed users: %v", err)
	}

	log.Printf("seeded %d users", len(users))

	return users
}
