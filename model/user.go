package model

import (
	"gorm.io/gorm"
	"time"
)

type User struct {
	ID            uint   `gorm:"primaryKey;index"`
	Uuid          string `gorm:"type:varchar(100);unique;index"`
	FullName      string `gorm:"type:varchar(100);index"`
	Email         string `gorm:"type:varchar(45);index"`
	PhoneNumber   string `gorm:"type:varchar(13);unique;index"`
	Password      string `gorm:"type:varchar(200)"`
	Role          string `gorm:"type:varchar(15)"`
	Status        string `gorm:"type:varchar(25)"`
	ReferenceID   uint
	ReferenceType string         `gorm:"type:varchar(25)"`
	Reference     any            `gorm:"-"`
	CreatedAt     time.Time      `gorm:"<-:create;type:datetime(0)"`
	UpdatedAt     time.Time      `gorm:"<-:update;type:datetime(0)"`
	DeletedAt     gorm.DeletedAt `gorm:"type:datetime(0);index"`
}
