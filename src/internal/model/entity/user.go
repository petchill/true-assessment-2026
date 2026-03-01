package entity

import "time"

type User struct {
	ID               uint64    `gorm:"column:id;primaryKey"`
	Age              int       `gorm:"column:age"`
	Country          string    `gorm:"column:country"`
	SubscriptionType string    `gorm:"column:subscription_type"`
	CreatedAt        time.Time `gorm:"column:created_at"`
}

func (User) TableName() string {
	return "users"
}
