package repository

import (
	"context"
	"fmt"
	"time"

	"gorm.io/gorm"

	"github.com/JasperRosales/aircraft-system-be/internal/models"
)

type PlanePartRepository struct {
	db *gorm.DB
}

func NewPlanePartRepository(db *gorm.DB) *PlanePartRepository {
	return &PlanePartRepository{db: db}
}

func (r *PlanePartRepository) Create(ctx context.Context, part *models.PlanePart) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	result := r.db.WithContext(ctx).Create(part)
	if result.Error != nil {
		return fmt.Errorf("failed to create plane part: %w", result.Error)
	}

	return nil
}

func (r *PlanePartRepository) GetByID(ctx context.Context, id int64) (*models.PlanePart, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var part models.PlanePart
	result := r.db.WithContext(ctx).First(&part, id)
	if result.Error == gorm.ErrRecordNotFound {
		return nil, nil
	}
	if result.Error != nil {
		return nil, fmt.Errorf("failed to get plane part by id: %w", result.Error)
	}

	return &part, nil
}

func (r *PlanePartRepository) GetBySerialNumber(ctx context.Context, serialNumber string) (*models.PlanePart, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var part models.PlanePart
	result := r.db.WithContext(ctx).Where("serial_number = ?", serialNumber).First(&part)
	if result.Error == gorm.ErrRecordNotFound {
		return nil, nil
	}
	if result.Error != nil {
		return nil, fmt.Errorf("failed to get plane part by serial number: %w", result.Error)
	}

	return &part, nil
}

func (r *PlanePartRepository) GetByPlaneID(ctx context.Context, planeID int64) ([]models.PlanePart, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var parts []models.PlanePart
	result := r.db.WithContext(ctx).Where("plane_id = ?", planeID).Order("id").Find(&parts)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to get plane parts by plane id: %w", result.Error)
	}

	return parts, nil
}

func (r *PlanePartRepository) GetByCategory(ctx context.Context, category string) ([]models.PlanePart, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var parts []models.PlanePart
	result := r.db.WithContext(ctx).Where("category = ?", category).Order("id").Find(&parts)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to get plane parts by category: %w", result.Error)
	}

	return parts, nil
}

func (r *PlanePartRepository) GetByPlaneIDAndCategory(ctx context.Context, planeID int64, category string) ([]models.PlanePart, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var parts []models.PlanePart
	result := r.db.WithContext(ctx).
		Where("plane_id = ? AND category = ?", planeID, category).
		Order("id").
		Find(&parts)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to get plane parts: %w", result.Error)
	}

	return parts, nil
}

func (r *PlanePartRepository) GetNeedingMaintenance(ctx context.Context, thresholdPercent float64) ([]models.PlanePart, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var parts []models.PlanePart
	result := r.db.WithContext(ctx).
		Where("usage_percent >= ?", thresholdPercent).
		Order("usage_percent DESC").
		Find(&parts)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to get parts needing maintenance: %w", result.Error)
	}

	return parts, nil
}

func (r *PlanePartRepository) GetAll(ctx context.Context) ([]models.PlanePart, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var parts []models.PlanePart
	result := r.db.WithContext(ctx).Order("id").Find(&parts)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to get all plane parts: %w", result.Error)
	}

	return parts, nil
}

func (r *PlanePartRepository) Update(ctx context.Context, part *models.PlanePart) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	result := r.db.WithContext(ctx).Save(part)
	if result.Error != nil {
		return fmt.Errorf("failed to update plane part: %w", result.Error)
	}

	return nil
}

func (r *PlanePartRepository) UpdateUsage(ctx context.Context, part *models.PlanePart) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	result := r.db.WithContext(ctx).Model(part).
		Updates(map[string]interface{}{
			"usage_hours": part.UsageHours,
		})
	if result.Error != nil {
		return fmt.Errorf("failed to update usage hours: %w", result.Error)
	}

	return nil
}

func (r *PlanePartRepository) Delete(ctx context.Context, id int64) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	result := r.db.WithContext(ctx).Delete(&models.PlanePart{}, id)
	if result.Error != nil {
		return fmt.Errorf("failed to delete plane part: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("plane part not found")
	}

	return nil
}

func (r *PlanePartRepository) GetByPlaneIDWithDetails(ctx context.Context, planeID int64) ([]models.PlanePart, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var parts []models.PlanePart
	result := r.db.WithContext(ctx).
		Preload("Plane").
		Where("plane_id = ?", planeID).
		Order("id").
		Find(&parts)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to get plane parts with details: %w", result.Error)
	}

	return parts, nil
}
