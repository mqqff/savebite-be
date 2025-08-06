package entity

import (
	"github.com/google/uuid"
	"time"
)

type Analysis struct {
	ID          uuid.UUID `gorm:"type:char(36);primaryKey"`
	UserID      uuid.UUID `gorm:"not null"`
	ImageURL    string    `gorm:"type:varchar(255);not null"`
	Feedback    string    `gorm:"type:mediumtext;not null"`
	Ingredients []Ingredient
	CreatedAt   time.Time `gorm:"type:timestamp;autoCreateTime"`
}
