package repository

import (
	"context"
	"fmt"
	"time"

	"gorm.io/gorm"

	"github.com/JasperRosales/aircraft-system-be/internal/models"
)

type PlaneRepository struct {
	db *gorm.DB
}

func NewPlaneRepository(db *gorm.DB) *PlaneRepository {
	return &PlaneRepository{db: db}
}

func (r *PlaneRepository) Create(ctx context.Context, plane *models.Plane) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	result := r.db.WithContext(ctx).Create(plane)
	if result.Error != nil {
		return fmt.Errorf("failed to create plane: %w", result.Error)
	}

	return nil
}

func (r *PlaneRepository) GetByID(ctx context.Context, id int64) (*models.Plane, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var plane models.Plane
	result := r.db.WithContext(ctx).First(&plane, id)
	if result.Error == gorm.ErrRecordNotFound {
		return nil, nil
	}
	if result.Error != nil {
		return nil, fmt.Errorf("failed to get plane by id: %w", result.Error)
	}

	return &plane, nil
}

func (r *PlaneRepository) GetByTailNumber(ctx context.Context, tailNumber string) (*models.Plane, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var plane models.Plane
	result := r.db.WithContext(ctx).Where("tail_number = ?", tailNumber).First(&plane)
	if result.Error == gorm.ErrRecordNotFound {
		return nil, nil
	}
	if result.Error != nil {
		return nil, fmt.Errorf("failed to get plane by tail number: %w", result.Error)
	}

	return &plane, nil
}

func (r *PlaneRepository) GetAll(ctx context.Context) ([]models.Plane, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var planes []models.Plane
	result := r.db.WithContext(ctx).Order("id").Find(&planes)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to get all planes: %w", result.Error)
	}

	return planes, nil
}

func (r *PlaneRepository) Update(ctx context.Context, plane *models.Plane) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	result := r.db.WithContext(ctx).Save(plane)
	if result.Error != nil {
		return fmt.Errorf("failed to update plane: %w", result.Error)
	}

	return nil
}

func (r *PlaneRepository) Delete(ctx context.Context, id int64) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	result := r.db.WithContext(ctx).Delete(&models.Plane{}, id)
	if result.Error != nil {
		return fmt.Errorf("failed to delete plane: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("plane not found")
	}

	return nil
}
