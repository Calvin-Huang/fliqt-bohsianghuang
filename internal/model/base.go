package model

import (
	"time"

	"github.com/rs/xid"
	"gorm.io/gorm"
)

func (base *Base) BeforeCreate(db *gorm.DB) error {
	if base.ID == "" {
		db.Statement.SetColumn("ID", xid.New().String())
	}
	return nil
}

// Base is a base model witch contains common fields for all models like gorm.Model
type Base struct {
	ID        string         `gorm:"not null;primaryKey;type:varchar(20);autoIncrement:false;primary_key;<-:create" json:"id"`
	CreatedAt time.Time      `gorm:"not null;autoCreateTime:nano" json:"created_at"`
	UpdatedAt time.Time      `gorm:"not null;autoUpdateTime:nano" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
