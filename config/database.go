package database

import (
	"denis/first/helpers"
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() {
	var err error
	DB, err = gorm.Open(mysql.Open("root@tcp(127.0.0.1:3306)/golang?charset=utf8mb4&parseTime=True&loc=Local"), &gorm.Config{})
	helpers.PanicIfError(err)
	fmt.Println("Connected to database!")
}
