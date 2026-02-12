package service

import (
	"context"
	"errors"
	"fmt"

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
		"tail_number", req.TailNumber,
		"model", req.Model,
	)

	existing, err := s.planeRepo.GetByTailNumber(ctx, req.TailNumber)
	if err != nil {
		s.logger.Error("PlaneService: Failed to check existing plane",
			"tail_number", req.TailNumber,
			"error", err,
		)
		return nil, fmt.Errorf("failed to check existing plane: %w", err)
	}
	if existing != nil {
		s.logger.Warn("PlaneService: Plane with tail number already exists",
			"tail_number", req.TailNumber,
		)
		return nil, errors.New(PlaneExistsErr)
	}

	plane := &models.Plane{
		TailNumber: req.TailNumber,
		Model:      req.Model,
	}

	if err := s.planeRepo.Create(ctx, plane); err != nil {
		s.logger.Error("PlaneService: Failed to create plane",
			"tail_number", req.TailNumber,
			"error", err,
		)
		return nil, fmt.Errorf("failed to create plane: %w", err)
	}

	s.logger.Info("PlaneService: Plane created successfully",
		"plane_id", plane.ID,
		"tail_number", plane.TailNumber,
	)

	resp := plane.ToResponse()
	return &resp, nil
}

func (s *PlaneService) GetPlane(ctx context.Context, id int64) (*models.PlaneResponse, error) {
	s.logger.Info("PlaneService: GetPlane",
		"plane_id", id,
	)

	plane, err := s.planeRepo.GetByID(ctx, id)
	if err != nil {
		s.logger.Error("PlaneService: Failed to get plane",
			"plane_id", id,
			"error", err,
		)
		return nil, fmt.Errorf("failed to get plane: %w", err)
	}
	if plane == nil {
		s.logger.Warn("PlaneService: Plane not found",
			"plane_id", id,
		)
		return nil, errors.New(PlaneNotFoundErr)
	}

	resp := plane.ToResponse()
	return &resp, nil
}

func (s *PlaneService) GetPlaneByTailNumber(ctx context.Context, tailNumber string) (*models.PlaneResponse, error) {
	s.logger.Info("PlaneService: GetPlaneByTailNumber",
		"tail_number", tailNumber,
	)

	plane, err := s.planeRepo.GetByTailNumber(ctx, tailNumber)
	if err != nil {
		s.logger.Error("PlaneService: Failed to get plane by tail number",
			"tail_number", tailNumber,
			"error", err,
		)
		return nil, fmt.Errorf("failed to get plane: %w", err)
	}
	if plane == nil {
		s.logger.Warn("PlaneService: Plane not found",
			"tail_number", tailNumber,
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
			"error", err,
		)
		return nil, fmt.Errorf("failed to get planes: %w", err)
	}

	s.logger.Info("PlaneService: GetAllPlanes successful",
		"count", len(planes),
	)

	responses := make([]models.PlaneResponse, len(planes))
	for i, plane := range planes {
		responses[i] = plane.ToResponse()
	}

	return responses, nil
}

func (s *PlaneService) UpdatePlane(ctx context.Context, id int64, req *models.UpdatePlaneRequest) (*models.PlaneResponse, error) {
	s.logger.Info("PlaneService: UpdatePlane",
		"plane_id", id,
	)

	plane, err := s.planeRepo.GetByID(ctx, id)
	if err != nil {
		s.logger.Error("PlaneService: Failed to get plane",
			"plane_id", id,
			"error", err,
		)
		return nil, fmt.Errorf("failed to get plane: %w", err)
	}
	if plane == nil {
		s.logger.Warn("PlaneService: Plane not found",
			"plane_id", id,
		)
		return nil, errors.New(PlaneNotFoundErr)
	}

	if req.TailNumber != nil {
		if *req.TailNumber != plane.TailNumber {
			existing, err := s.planeRepo.GetByTailNumber(ctx, *req.TailNumber)
			if err != nil {
				s.logger.Error("PlaneService: Failed to check existing plane",
					"tail_number", *req.TailNumber,
					"error", err,
				)
				return nil, fmt.Errorf("failed to check existing plane: %w", err)
			}
			if existing != nil {
				s.logger.Warn("PlaneService: Plane with tail number already exists",
					"tail_number", *req.TailNumber,
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
			"plane_id", id,
			"error", err,
		)
		return nil, fmt.Errorf("failed to update plane: %w", err)
	}

	s.logger.Info("PlaneService: Update successful",
		"plane_id", id,
	)

	resp := plane.ToResponse()
	return &resp, nil
}

func (s *PlaneService) DeletePlane(ctx context.Context, id int64) error {
	s.logger.Info("PlaneService: DeletePlane",
		"plane_id", id,
	)

	plane, err := s.planeRepo.GetByID(ctx, id)
	if err != nil {
		s.logger.Error("PlaneService: Failed to get plane",
			"plane_id", id,
			"error", err,
		)
		return fmt.Errorf("failed to get plane: %w", err)
	}
	if plane == nil {
		s.logger.Warn("PlaneService: Plane not found",
			"plane_id", id,
		)
		return errors.New(PlaneNotFoundErr)
	}

	if err := s.planeRepo.Delete(ctx, id); err != nil {
		s.logger.Error("PlaneService: Failed to delete plane",
			"plane_id", id,
			"error", err,
		)
		return fmt.Errorf("failed to delete plane: %w", err)
	}

	s.logger.Info("PlaneService: Delete successful",
		"plane_id", id,
	)

	return nil
}

func (s *PlaneService) GetPlaneWithParts(ctx context.Context, id int64) (*models.PlaneResponse, []models.PlanePartResponse, error) {
	s.logger.Info("PlaneService: GetPlaneWithParts",
		"plane_id", id,
	)

	plane, err := s.planeRepo.GetByID(ctx, id)
	if err != nil {
		s.logger.Error("PlaneService: Failed to get plane",
			"plane_id", id,
			"error", err,
		)
		return nil, nil, fmt.Errorf("failed to get plane: %w", err)
	}
	if plane == nil {
		s.logger.Warn("PlaneService: Plane not found",
			"plane_id", id,
		)
		return nil, nil, errors.New(PlaneNotFoundErr)
	}

	planeResp := plane.ToResponse()
	return &planeResp, nil, nil
}
