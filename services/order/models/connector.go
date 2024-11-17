package models

import "gorm.io/gorm"

func Connect(db *gorm.DB) {
	db.AutoMigrate(&Order{}, &LineItem{})
}
