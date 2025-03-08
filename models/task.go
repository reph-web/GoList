package models

import (
	"time"
)

type Task struct {
	ID          uint   `gorm:"primaryKey"`
	Description string `gorm:"type:varchar(255)"`
	TaskOrder   uint   `gorm:"unique"`
	Checked     bool
	CreatedAt   time.Time

	UserID uint `gorm:"index"`
	ListID uint `gorm:"index"`
}
