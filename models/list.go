package models

import (
	"time"
)

type List struct {
	ID        uint   `gorm:"primaryKey"`
	Name      string `gorm:"type:varchar(100)"`
	CreatedAt time.Time

	UserID uint   `gorm:"index"`
	Tasks  []Task `gorm:"foreignKey:ListID"`
}
