package entity

import "time"

type Content struct {
	ID              uint64    `gorm:"column:id;primaryKey"`
	Title           string    `gorm:"column:title"`
	Genre           string    `gorm:"column:genre"`
	PopularityScore float64   `gorm:"column:popularity_score"`
	CreatedAt       time.Time `gorm:"column:created_at"`
}

func (Content) TableName() string {
	return "content"
}

func (c Content) GetRecencyFactory() float64 {
	daySinceCreation := time.Since(c.CreatedAt).Hours() / 24
	return 1 / (1 + (daySinceCreation / 365))
}

type HistoryContent struct {
	ID        uint64    `gorm:"column:id"`
	Genre     string    `gorm:"column:genre"`
	WatchedAt time.Time `gorm:"column:watched_at"`
}
