package database

import (
    "gorm.io/gorm"
    "gorm.io/driver/sqlite"
)
var (
	DB  *gorm.DB
)

func ConectaNoBD() {
	db, err := gorm.Open(sqlite.Open("webFilter.db"), &gorm.Config{})
  	if err != nil {
    	panic("failed to connect database")
  	}

	// Criar table via migrate
	//db.AutoMigrate(&models.WebFilter{})

	DB = db
}

func GetDataBase() *gorm.DB {
	return DB
}