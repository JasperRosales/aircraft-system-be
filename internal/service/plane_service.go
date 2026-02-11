package service

import (
	"context"
	"errors"
	"fmt"

	"go.uber.org/zap"

	"github.com/JasperRosales/aircraft-system-be/internal/models"
	"github.com/JasperRosales/aircraft-system-be/internal/repository"
	"github.com/JasperRosales/aircraft-system-be/internal/util"
)

const (
	PlaneNotFoundErr = "plane not found"
	PlaneExistsErr   = "plane with this tail number already exists"
)

type PlaneService struct {
	planeRepo *repository.PlaneRepository
	logger    *util.Logger
}

func NewPlaneService(planeRepo *repository.PlaneRepository, logger *util.Logger) *PlaneService {
	return &PlaneService{
		planeRepo: planeRepo,
		logger:    logger,
	}
}

func (s *PlaneService) CreatePlane(ctx context.Context, req *models.CreatePlaneRequest) (*models.PlaneResponse, error) {
	s.logger.Info("PlaneService: Creating new plane",
		zap.String("tail_number", req.TailNumber),
		zap.String("model", req.Model),
	)

	existing, err := s.planeRepo.GetByTailNumber(ctx, req.TailNumber)
	if err != nil {
		s.logger.Error("PlaneService: Failed to check existing plane",
			zap.String("tail_number", req.TailNumber),
			zap.Error(err),
		)
		return nil, fmt.Errorf("failed to check existing plane: %w", err)
	}
	if existing != nil {
		s.logger.Warn("PlaneService: Plane with tail number already exists",
			zap.String("tail_number", req.TailNumber),
		)
		return nil, errors.New(PlaneExistsErr)
	}

	plane := &models.Plane{
		TailNumber: req.TailNumber,
		Model:      req.Model,
	}

	if err := s.planeRepo.Create(ctx, plane); err != nil {
		s.logger.Error("PlaneService: Failed to create plane",
			zap.String("tail_number", req.TailNumber),
			zap.Error(err),
		)
		return nil, fmt.Errorf("failed to create plane: %w", err)
	}

	s.logger.Info("PlaneService: Plane created successfully",
		zap.Int64("plane_id", plane.ID),
		zap.String("tail_number", plane.TailNumber),
	)

	resp := plane.ToResponse()
	return &resp, nil
}

func (s *PlaneService) GetPlane(ctx context.Context, id int64) (*models.PlaneResponse, error) {
	s.logger.Info("PlaneService: GetPlane",
		zap.Int64("plane_id", id),
	)

	plane, err := s.planeRepo.GetByID(ctx, id)
	if err != nil {
		s.logger.Error("PlaneService: Failed to get plane",
			zap.Int64("plane_id", id),
			zap.Error(err),
		)
		return nil, fmt.Errorf("failed to get plane: %w", err)
	}
	if plane == nil {
		s.logger.Warn("PlaneService: Plane not found",
			zap.Int64("plane_id", id),
		)
		return nil, errors.New(PlaneNotFoundErr)
	}

	resp := plane.ToResponse()
	return &resp, nil
}

func (s *PlaneService) GetPlaneByTailNumber(ctx context.Context, tailNumber string) (*models.PlaneResponse, error) {
	s.logger.Info("PlaneService: GetPlaneByTailNumber",
		zap.String("tail_number", tailNumber),
	)

	plane, err := s.planeRepo.GetByTailNumber(ctx, tailNumber)
	if err != nil {
		s.logger.Error("PlaneService: Failed to get plane by tail number",
			zap.String("tail_number", tailNumber),
			zap.Error(err),
		)
		return nil, fmt.Errorf("failed to get plane: %w", err)
	}
	if plane == nil {
		s.logger.Warn("PlaneService: Plane not found",
			zap.String("tail_number", tailNumber),
		)
		return nil, errors.New(PlaneNotFoundErr)
	}

	resp := plane.ToResponse()
	return &resp, nil
}

func (s *PlaneService) GetAllPlanes(ctx context.Context) ([]models.PlaneResponse, error) {
	s.logger.Info("PlaneService: GetAllPlanes")

	planes, err := s.planeRepo.GetAll(ctx)
	if err != nil {
		s.logger.Error("PlaneService: Failed to get planes",
			zap.Error(err),
		)
		return nil, fmt.Errorf("failed to get planes: %w", err)
	}

	s.logger.Info("PlaneService: GetAllPlanes successful",
		zap.Int("count", len(planes)),
	)

	responses := make([]models.PlaneResponse, len(planes))
	for i, plane := range planes {
		responses[i] = plane.ToResponse()
	}

	return responses, nil
}

func (s *PlaneService) UpdatePlane(ctx context.Context, id int64, req *models.UpdatePlaneRequest) (*models.PlaneResponse, error) {
	s.logger.Info("PlaneService: UpdatePlane",
		zap.Int64("plane_id", id),
	)

	plane, err := s.planeRepo.GetByID(ctx, id)
	if err != nil {
		s.logger.Error("PlaneService: Failed to get plane",
			zap.Int64("plane_id", id),
			zap.Error(err),
		)
		return nil, fmt.Errorf("failed to get plane: %w", err)
	}
	if plane == nil {
		s.logger.Warn("PlaneService: Plane not found",
			zap.Int64("plane_id", id),
		)
		return nil, errors.New(PlaneNotFoundErr)
	}

	if req.TailNumber != nil {
		if *req.TailNumber != plane.TailNumber {
			existing, err := s.planeRepo.GetByTailNumber(ctx, *req.TailNumber)
			if err != nil {
				s.logger.Error("PlaneService: Failed to check existing plane",
					zap.String("tail_number", *req.TailNumber),
					zap.Error(err),
				)
				return nil, fmt.Errorf("failed to check existing plane: %w", err)
			}
			if existing != nil {
				s.logger.Warn("PlaneService: Plane with tail number already exists",
					zap.String("tail_number", *req.TailNumber),
				)
				return nil, errors.New(PlaneExistsErr)
			}
		}
		plane.TailNumber = *req.TailNumber
	}
	if req.Model != nil {
		plane.Model = *req.Model
	}

	if err := s.planeRepo.Update(ctx, plane); err != nil {
		s.logger.Error("PlaneService: Failed to update plane",
			zap.Int64("plane_id", id),
			zap.Error(err),
		)
		return nil, fmt.Errorf("failed to update plane: %w", err)
	}

	s.logger.Info("PlaneService: Update successful",
		zap.Int64("plane_id", id),
	)

	resp := plane.ToResponse()
	return &resp, nil
}

func (s *PlaneService) DeletePlane(ctx context.Context, id int64) error {
	s.logger.Info("PlaneService: DeletePlane",
		zap.Int64("plane_id", id),
	)

	plane, err := s.planeRepo.GetByID(ctx, id)
	if err != nil {
		s.logger.Error("PlaneService: Failed to get plane",
			zap.Int64("plane_id", id),
			zap.Error(err),
		)
		return fmt.Errorf("failed to get plane: %w", err)
	}
	if plane == nil {
		s.logger.Warn("PlaneService: Plane not found",
			zap.Int64("plane_id", id),
		)
		return errors.New(PlaneNotFoundErr)
	}

	if err := s.planeRepo.Delete(ctx, id); err != nil {
		s.logger.Error("PlaneService: Failed to delete plane",
			zap.Int64("plane_id", id),
			zap.Error(err),
		)
		return fmt.Errorf("failed to delete plane: %w", err)
	}

	s.logger.Info("PlaneService: Delete successful",
		zap.Int64("plane_id", id),
	)

	return nil
}

func (s *PlaneService) GetPlaneWithParts(ctx context.Context, id int64) (*models.PlaneResponse, []models.PlanePartResponse, error) {
	s.logger.Info("PlaneService: GetPlaneWithParts",
		zap.Int64("plane_id", id),
	)

	plane, err := s.planeRepo.GetByID(ctx, id)
	if err != nil {
		s.logger.Error("PlaneService: Failed to get plane",
			zap.Int64("plane_id", id),
			zap.Error(err),
		)
		return nil, nil, fmt.Errorf("failed to get plane: %w", err)
	}
	if plane == nil {
		s.logger.Warn("PlaneService: Plane not found",
			zap.Int64("plane_id", id),
		)
		return nil, nil, errors.New(PlaneNotFoundErr)
	}

	planeResp := plane.ToResponse()
	return &planeResp, nil, nil
}
