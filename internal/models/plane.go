package models

import (
	"time"
)

type Plane struct {
	ID         int64     `json:"id" gorm:"primaryKey;autoIncrement"`
	TailNumber string    `json:"tail_number" gorm:"type:varchar(50);uniqueIndex;not null"`
	Model      string    `json:"model" gorm:"type:varchar(100);not null"`
	CreatedAt  time.Time `json:"created_at" gorm:"autoCreateTime"`
}

type CreatePlaneRequest struct {
	TailNumber string `json:"tail_number" binding:"required,min=2,max=50"`
	Model      string `json:"model" binding:"required,min=2,max=100"`
}

type UpdatePlaneRequest struct {
	TailNumber *string `json:"tail_number" binding:"omitempty,min=2,max=50"`
	Model      *string `json:"model" binding:"omitempty,min=2,max=100"`
}

type PlaneResponse struct {
	ID         int64     `json:"id"`
	TailNumber string    `json:"tail_number"`
	Model      string    `json:"model"`
	CreatedAt  time.Time `json:"created_at"`
}

func (p *Plane) ToResponse() PlaneResponse {
	return PlaneResponse{
		ID:         p.ID,
		TailNumber: p.TailNumber,
		Model:      p.Model,
		CreatedAt:  p.CreatedAt,
	}
}
