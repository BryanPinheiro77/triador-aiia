package model

import "time"

type Analysis struct {
	ID uint `gorm:"primaryKey"`

	CandidateName   string `gorm:"not null"`
	Skills          string `gorm:"type:text"`
	YearsExperience int
	FitScore        int
	Summary         string `gorm:"type:text"`

	CreatedAt time.Time
}
