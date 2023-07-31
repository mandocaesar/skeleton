package model

import (
	"time"

	"gorm.io/gorm"
)

type EntityModel interface {
	TableName() string
}

// BaseModel is the base model for DB entity struct within this project
type BaseModel struct {
	ID        int64           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
	DeletedAt *gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
	CreatedBy *string         `json:"created_by,omitempty"`
	UpdatedBy *string         `json:"updated_by,omitempty"`
	DeletedBy *string         `json:"deleted_by,omitempty"`
}

type EntityMongoModel interface {
	CollectionName() string
}
