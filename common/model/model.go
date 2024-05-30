package model

import (
	"gorm.io/gorm"
	"time"
)

type Model struct {
	ID        uint           `gorm:"primarykey"`
	CreatedAt time.Time      `json:"created_at" gorm:"index,autoCreateTime"`
	UpdatedAt time.Time      `json:"updated_at" gorm:"autoUpdateTime:nano"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}
