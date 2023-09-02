package models

import "gorm.io/gorm"

type AdMessages struct {
	gorm.Model
	FromUser uint   `json:"from"`
	ToUser   uint   `json:"to"`
	AdID     string `json:"adID"`
	Message  string `json:"message" `
}
