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

type HistoryContent struct {
	ID        uint64    `gorm:"column:id"`
	Genre     string    `gorm:"column:genre"`
	WatchedAt time.Time `gorm:"column:watched_at"`
}
