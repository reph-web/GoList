package database

import (
	"fmt"
	"golist/models"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() {
	dsn := "root@tcp(127.0.0.1:3306)/nita_service?charset=utf8mb4&parseTime=True&loc=Local"
	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("❌ Cannot connect to MySQL:", err)
	}

	DB.AutoMigrate(&models.User{}, &models.List{}, &models.Task{})
	DB = DB.Debug()
	fmt.Println("✅ Connected to MySQL")
}
