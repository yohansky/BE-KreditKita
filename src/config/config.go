package config

import (
	"log"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() {
	var err error

	DSN := os.Getenv("DSN")

	for i := 0; i < 10; i++ {
		DB, err = gorm.Open(mysql.Open(DSN), &gorm.Config{})
		if err == nil {
			log.Println("Connect to DB")
			return
		}
		time.Sleep(2 * time.Second)
	}
	if err != nil {
		panic("couldn't connect with database")
	}
}
