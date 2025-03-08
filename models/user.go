package models

import (
	"time"
)

type User struct {
	ID        uint   `gorm:"primary_key"`
	Username  string `gorm:"type:varchar(32)"`
	Password  string `gorm:"type:varchar(100)"`
	CreatedAt time.Time

	Lists []List `gorm:"foreignKey:UserID"`
	Tasks []Task `gorm:"foreignKey:UserID"`
}
