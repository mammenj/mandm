package storage

import (
	"fmt"

	"log"

	"github.com/mammenj/mandm/models"
	"gorm.io/driver/sqlite"

	//"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

type MessageSqlliteStore struct {
	DB *gorm.DB
}

func NewSqliteMessageStore() *MessageSqlliteStore {
	db, err := gorm.Open(sqlite.Open("matri.db"), &gorm.Config{})
	db.AutoMigrate(models.Messages{})
	if err != nil {
		panic("failed to connect database")
	}
	return &MessageSqlliteStore{
		DB: db,
	}
}

func (as *MessageSqlliteStore) Create(msg *models.Messages) (string, error) {
	log.Println("Before migrate......Messages")
	//as.DB.AutoMigrate(ad)
	log.Println("...... After migrate Messages")
	result := as.DB.Create(msg)
	if result.Error != nil {
		return "", result.Error
	}
	log.Println("...... After Create AD..... ID: ", msg.ID)
	return fmt.Sprint(msg.ID), nil
}

func (as *MessageSqlliteStore) Get() ([]models.Messages, error) {
	var msg []models.Messages
	log.Println("Get Ads")
	result := as.DB.Find(&msg)

	if result.Error != nil {
		return nil, result.Error
	}
	log.Println("...... Total records msg : ", result.RowsAffected)
	return msg, nil
}


