package storage

import (
	"fmt"

	"log"

	"github.com/mammenj/mandm/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type AdMessageSqlliteStore struct {
	DB *gorm.DB
}

func NewSqliteAdMessageStore() *AdMessageSqlliteStore {
	db, err := gorm.Open(sqlite.Open("matri.db"), &gorm.Config{})
	db.AutoMigrate(models.AdMessages{})
	if err != nil {
		panic("failed to connect database")
	}
	return &AdMessageSqlliteStore{
		DB: db,
	}
}

func (as *AdMessageSqlliteStore) Create(msg *models.AdMessages) (string, error) {
	log.Println("...... AdMessagess ")
	result := as.DB.Create(msg)
	if result.Error != nil {
		return "", result.Error
	}
	log.Println("...... After Create AdMessagess..... ID: ", msg.ID)
	return fmt.Sprint(msg.ID), nil
}

func (as *AdMessageSqlliteStore) Get() ([]models.AdMessages, error) {
	var msg []models.AdMessages
	log.Println("Get AdMessagess")
	result := as.DB.Find(&msg)

	if result.Error != nil {
		return nil, result.Error
	}
	log.Println("...... Total records msg : ", result.RowsAffected)
	return msg, nil
}

func (as *AdMessageSqlliteStore) GetMessagesToID(toId uint) ([]models.AdMessages, error) {
	var msg []models.AdMessages
	log.Println("Get AdMessagess by to_id")
	result := as.DB.Where("to_user = ?", toId).Find(&msg)
	if result.Error != nil {
		return nil, result.Error
	}
	log.Println("...... Total records msg : ", result.RowsAffected)
	return msg, nil
}

func (as *AdMessageSqlliteStore) GetMessagesToIDFromID(toId, fromID uint) ([]models.AdMessages, error) {
	var msg []models.AdMessages
	log.Println("Get GetMessagesToIDFromID by from and to")
	result := as.DB.Where("(to_user = ? and from_user = ?) or (from_user = ? and to_user = ?)", toId, fromID, toId, fromID).Find(&msg)
	log.Println("...... Total records GetMessagesToIDFromID : ", result.RowsAffected)
	if result.Error != nil {
		return nil, result.Error
	}

	return msg, nil
}

func (as *AdMessageSqlliteStore) GetMessagesToIDGroupByFrom(toId uint) ([]models.AdMessages, error) {
	var msg []models.AdMessages
	log.Println("Get GetMessagesToIDGroupByFrom by to_id")

	result := as.DB.Raw("SELECT * FROM ad_messages WHERE to_user = ? and from_user != to_user or from_user = ? GROUP BY from_user", toId, toId).Scan(&msg)
	if result.Error != nil {
		return nil, result.Error
	}
	log.Println("...... Total GetMessagesToIDGroupByFrom records msg : ", result.RowsAffected)
	return msg, nil
}

func (as *AdMessageSqlliteStore) GetMessagesToIDandFrom(toId uint) ([]models.AdMessages, error) {
	var msg []models.AdMessages
	log.Println("Get GetMessagesToIDandFrom by to_id")

	result := as.DB.Raw("SELECT * FROM ad_messages WHERE to_user = ? and from_user != to_user or from_user = ? or from_user = ?  GROUP BY from_user", toId, toId, toId).Scan(&msg)
	if result.Error != nil {
		return nil, result.Error
	}
	log.Println("...... Total GetMessagesToIDandFrom records msg : ", result.RowsAffected)
	return msg, nil
}
