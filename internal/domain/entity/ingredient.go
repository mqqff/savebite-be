package entity

import "github.com/google/uuid"

type Ingredient struct {
	ID         uint      `gorm:"primaryKey;autoIncrement"`
	AnalysisID uuid.UUID `gorm:"not null"`
	Name       string    `gorm:"type:varchar(255);not null"`
	Status     string    `gorm:"type:enum('usable','unusable');not null"`
}
