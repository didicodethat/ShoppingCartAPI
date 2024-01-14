package database

import (
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Make sure to only use it after calling the database.SetupDB() function
var DB *gorm.DB

// sets up the DB and returns, but also as a side effect sets the DB global variable
func SetupDB() (*gorm.DB, error) {

	// pretty much temporary config
	dsn := "host=localhost user=postgres password=postgres dbname=shopping_cart port=5432"

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		return nil, err
	}

	err = db.AutoMigrate(&User{}, &ShoppingList{}, &ListItem{})

	if err != nil {
		return nil, err
	}

	DB = db

	logger := log.Default()
	logger.Println("Successfully set the DB global variable (database.DB)")
	return db, nil
}
