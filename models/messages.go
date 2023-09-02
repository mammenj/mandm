package models

import (
	// "time"

	//"github.com/google/uuid"
	"gorm.io/gorm"
)

type Messages struct {
	gorm.Model
	Name    string `json:"name"`
	Email   string `json:"email"`
	Subject string `json:"subject"`
	Message string `json:"message" `
}

// BeforeCreate BeforeCreate
/* func (u *User) BeforeCreate(scope *gorm.DB) error {
	//scope.Statement.SetColumn("UpdatedAt", time.Now())
	scope.Statement.SetColumn("UID", uuid.New().String())
	return nil
} */

// BeforeUpdate BeforeUpdate
/* func (u *User) BeforeUpdate(scope *gorm.DB) error {
	scope.Statement.SetColumn("UpdatedAt", time.Now())
	return nil
} */
