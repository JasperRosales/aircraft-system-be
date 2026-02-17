package models

import (
	"time"
)

type PlanePart struct {
	ID              int64     `json:"id" gorm:"primaryKey;autoIncrement"`
	PlaneID         int64     `json:"plane_id" gorm:"not null;index"`
	PartName        string    `json:"part_name" gorm:"type:varchar(255);not null"`
	SerialNumber    string    `json:"serial_number" gorm:"type:varchar(100);uniqueIndex;not null"`
	Category        string    `json:"category" gorm:"type:varchar(150);not null;index"`
	UsageHours      float64   `json:"usage_hours" gorm:"type:numeric(10,2);default:0"`
	UsageLimitHours float64   `json:"usage_limit_hours" gorm:"type:numeric(10,2);not null"`
	UsagePercent    *float64  `json:"usage_percent" gorm:"-"`
	InstalledAt     time.Time `json:"installed_at" gorm:"autoCreateTime"`
	Plane           *Plane    `json:"plane,omitempty" gorm:"foreignKey:PlaneID"`
}

type CreatePlanePartRequest struct {
	PlaneID         int64   `json:"plane_id" binding:"required"`
	PartName        string  `json:"part_name" binding:"required,min=2,max=255"`
	SerialNumber    string  `json:"serial_number" binding:"required,min=2,max=100"`
	Category        string  `json:"category" binding:"required,min=2,max=150"`
	UsageHours      float64 `json:"usage_hours"`
	UsageLimitHours float64 `json:"usage_limit_hours" binding:"required,gt=0"`
}

type UpdatePlanePartRequest struct {
	PartName        *string  `json:"part_name" binding:"omitempty,min=2,max=255"`
	SerialNumber    *string  `json:"serial_number" binding:"omitempty,min=2,max=100"`
	Category        *string  `json:"category" binding:"omitempty,min=2,max=150"`
	UsageLimitHours *float64 `json:"usage_limit_hours" binding:"omitempty,gt=0"`
}

type UpdatePartUsageRequest struct {
	UsageHours float64 `json:"usage_hours" binding:"required,gte=0"`
}

type PlanePartResponse struct {
	ID              int64     `json:"id"`
	PlaneID         int64     `json:"plane_id"`
	PartName        string    `json:"part_name"`
	SerialNumber    string    `json:"serial_number"`
	Category        string    `json:"category"`
	UsageHours      float64   `json:"usage_hours"`
	UsageLimitHours float64   `json:"usage_limit_hours"`
	UsagePercent    float64   `json:"usage_percent"`
	InstalledAt     time.Time `json:"installed_at"`
}

func (pp *PlanePart) ToResponse() PlanePartResponse {
	resp := PlanePartResponse{
		ID:              pp.ID,
		PlaneID:         pp.PlaneID,
		PartName:        pp.PartName,
		SerialNumber:    pp.SerialNumber,
		Category:        pp.Category,
		UsageHours:      pp.UsageHours,
		UsageLimitHours: pp.UsageLimitHours,
		InstalledAt:     pp.InstalledAt,
	}
	if pp.UsagePercent != nil {
		resp.UsagePercent = *pp.UsagePercent
	} else if pp.UsageLimitHours > 0 {
		resp.UsagePercent = (pp.UsageHours / pp.UsageLimitHours) * 100
	}
	return resp
}

func (pp *PlanePart) ToResponseWithPlane() PlanePartResponse {
	resp := pp.ToResponse()
	if pp.Plane != nil {
		resp.PlaneID = pp.Plane.ID
	}
	return resp
}

type PlanePartsByPlaneQuery struct {
	Category *string `form:"category"`
}

type MaintenanceAlertQuery struct {
	Threshold float64 `form:"threshold" binding:"omitempty,gte=0,lte=100"`
}
