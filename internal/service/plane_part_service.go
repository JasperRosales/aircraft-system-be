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
	PlanePartNotFoundErr = "plane part not found"
	PlanePartExistsErr   = "plane part with this serial number already exists"
	InvalidUsageHoursErr = "usage hours cannot exceed limit"
	PlaneNotMatchErr     = "plane part does not belong to this plane"
	PlaneNotFoundErrPart = "plane not found"
)

type PlanePartService struct {
	planeRepo     *repository.PlaneRepository
	planePartRepo *repository.PlanePartRepository
	logger        *util.Logger
}

func NewPlanePartService(planeRepo *repository.PlaneRepository, planePartRepo *repository.PlanePartRepository, logger *util.Logger) *PlanePartService {
	return &PlanePartService{
		planeRepo:     planeRepo,
		planePartRepo: planePartRepo,
		logger:        logger,
	}
}

func (s *PlanePartService) AddPart(ctx context.Context, req *models.CreatePlanePartRequest) (*models.PlanePartResponse, error) {
	s.logger.Info("PlanePartService: Adding new part to plane",
		zap.Int64("plane_id", req.PlaneID),
		zap.String("part_name", req.PartName),
		zap.String("serial_number", req.SerialNumber),
	)

	plane, err := s.planeRepo.GetByID(ctx, req.PlaneID)
	if err != nil {
		s.logger.Error("PlanePartService: Failed to verify plane",
			zap.Int64("plane_id", req.PlaneID),
			zap.Error(err),
		)
		return nil, fmt.Errorf("failed to verify plane: %w", err)
	}
	if plane == nil {
		s.logger.Warn("PlanePartService: Plane not found",
			zap.Int64("plane_id", req.PlaneID),
		)
		return nil, errors.New(PlaneNotFoundErrPart)
	}

	existing, err := s.planePartRepo.GetBySerialNumber(ctx, req.SerialNumber)
	if err != nil {
		s.logger.Error("PlanePartService: Failed to check existing part",
			zap.String("serial_number", req.SerialNumber),
			zap.Error(err),
		)
		return nil, fmt.Errorf("failed to check existing part: %w", err)
	}
	if existing != nil {
		s.logger.Warn("PlanePartService: Part with serial number already exists",
			zap.String("serial_number", req.SerialNumber),
		)
		return nil, errors.New(PlanePartExistsErr)
	}

	part := &models.PlanePart{
		PlaneID:         req.PlaneID,
		PartName:        req.PartName,
		SerialNumber:    req.SerialNumber,
		Category:        req.Category,
		UsageHours:      req.UsageHours,
		UsageLimitHours: req.UsageLimitHours,
	}

	if err := s.planePartRepo.Create(ctx, part); err != nil {
		s.logger.Error("PlanePartService: Failed to create part",
			zap.String("serial_number", req.SerialNumber),
			zap.Error(err),
		)
		return nil, fmt.Errorf("failed to create part: %w", err)
	}

	s.logger.Info("PlanePartService: Part added successfully",
		zap.Int64("part_id", part.ID),
		zap.Int64("plane_id", req.PlaneID),
		zap.String("serial_number", req.SerialNumber),
	)

	resp := part.ToResponse()
	return &resp, nil
}

func (s *PlanePartService) GetPart(ctx context.Context, id int64) (*models.PlanePartResponse, error) {
	s.logger.Info("PlanePartService: GetPart",
		zap.Int64("part_id", id),
	)

	part, err := s.planePartRepo.GetByID(ctx, id)
	if err != nil {
		s.logger.Error("PlanePartService: Failed to get part",
			zap.Int64("part_id", id),
			zap.Error(err),
		)
		return nil, fmt.Errorf("failed to get part: %w", err)
	}
	if part == nil {
		s.logger.Warn("PlanePartService: Part not found",
			zap.Int64("part_id", id),
		)
		return nil, errors.New(PlanePartNotFoundErr)
	}

	resp := part.ToResponse()
	return &resp, nil
}

func (s *PlanePartService) GetPartsByPlane(ctx context.Context, planeID int64, category *string) ([]models.PlanePartResponse, error) {
	s.logger.Info("PlanePartService: GetPartsByPlane",
		zap.Int64("plane_id", planeID),
	)

	// Verify plane exists
	plane, err := s.planeRepo.GetByID(ctx, planeID)
	if err != nil {
		s.logger.Error("PlanePartService: Failed to verify plane",
			zap.Int64("plane_id", planeID),
			zap.Error(err),
		)
		return nil, fmt.Errorf("failed to verify plane: %w", err)
	}
	if plane == nil {
		s.logger.Warn("PlanePartService: Plane not found",
			zap.Int64("plane_id", planeID),
		)
		return nil, errors.New(PlaneNotFoundErrPart)
	}

	var parts []models.PlanePart
	if category != nil && *category != "" {
		parts, err = s.planePartRepo.GetByPlaneIDAndCategory(ctx, planeID, *category)
		if err != nil {
			s.logger.Error("PlanePartService: Failed to get parts by category",
				zap.Int64("plane_id", planeID),
				zap.String("category", *category),
				zap.Error(err),
			)
			return nil, fmt.Errorf("failed to get parts: %w", err)
		}
	} else {
		parts, err = s.planePartRepo.GetByPlaneID(ctx, planeID)
		if err != nil {
			s.logger.Error("PlanePartService: Failed to get parts",
				zap.Int64("plane_id", planeID),
				zap.Error(err),
			)
			return nil, fmt.Errorf("failed to get parts: %w", err)
		}
	}

	s.logger.Info("PlanePartService: GetPartsByPlane successful",
		zap.Int64("plane_id", planeID),
		zap.Int("count", len(parts)),
	)

	responses := make([]models.PlanePartResponse, len(parts))
	for i, part := range parts {
		responses[i] = part.ToResponse()
	}

	return responses, nil
}

func (s *PlanePartService) GetAllParts(ctx context.Context) ([]models.PlanePartResponse, error) {
	s.logger.Info("PlanePartService: GetAllParts")

	parts, err := s.planePartRepo.GetAll(ctx)
	if err != nil {
		s.logger.Error("PlanePartService: Failed to get all parts",
			zap.Error(err),
		)
		return nil, fmt.Errorf("failed to get all parts: %w", err)
	}

	s.logger.Info("PlanePartService: GetAllParts successful",
		zap.Int("count", len(parts)),
	)

	responses := make([]models.PlanePartResponse, len(parts))
	for i, part := range parts {
		responses[i] = part.ToResponse()
	}

	return responses, nil
}

func (s *PlanePartService) UpdatePart(ctx context.Context, id int64, req *models.UpdatePlanePartRequest) (*models.PlanePartResponse, error) {
	s.logger.Info("PlanePartService: UpdatePart",
		zap.Int64("part_id", id),
	)

	part, err := s.planePartRepo.GetByID(ctx, id)
	if err != nil {
		s.logger.Error("PlanePartService: Failed to get part",
			zap.Int64("part_id", id),
			zap.Error(err),
		)
		return nil, fmt.Errorf("failed to get part: %w", err)
	}
	if part == nil {
		s.logger.Warn("PlanePartService: Part not found",
			zap.Int64("part_id", id),
		)
		return nil, errors.New(PlanePartNotFoundErr)
	}

	if req.PartName != nil {
		part.PartName = *req.PartName
	}
	if req.Category != nil {
		part.Category = *req.Category
	}
	if req.SerialNumber != nil {
		if *req.SerialNumber != part.SerialNumber {
			existing, err := s.planePartRepo.GetBySerialNumber(ctx, *req.SerialNumber)
			if err != nil {
				s.logger.Error("PlanePartService: Failed to check existing part",
					zap.String("serial_number", *req.SerialNumber),
					zap.Error(err),
				)
				return nil, fmt.Errorf("failed to check existing part: %w", err)
			}
			if existing != nil {
				s.logger.Warn("PlanePartService: Part with serial number already exists",
					zap.String("serial_number", *req.SerialNumber),
				)
				return nil, errors.New(PlanePartExistsErr)
			}
		}
		part.SerialNumber = *req.SerialNumber
	}
	if req.UsageLimitHours != nil {
		part.UsageLimitHours = *req.UsageLimitHours
	}

	if err := s.planePartRepo.Update(ctx, part); err != nil {
		s.logger.Error("PlanePartService: Failed to update part",
			zap.Int64("part_id", id),
			zap.Error(err),
		)
		return nil, fmt.Errorf("failed to update part: %w", err)
	}

	s.logger.Info("PlanePartService: UpdatePart successful",
		zap.Int64("part_id", id),
	)

	resp := part.ToResponse()
	return &resp, nil
}

func (s *PlanePartService) UpdatePartUsage(ctx context.Context, id int64, req *models.UpdatePartUsageRequest) (*models.PlanePartResponse, error) {
	s.logger.Info("PlanePartService: UpdatePartUsage",
		zap.Int64("part_id", id),
		zap.Float64("new_usage_hours", req.UsageHours),
	)

	part, err := s.planePartRepo.GetByID(ctx, id)
	if err != nil {
		s.logger.Error("PlanePartService: Failed to get part",
			zap.Int64("part_id", id),
			zap.Error(err),
		)
		return nil, fmt.Errorf("failed to get part: %w", err)
	}
	if part == nil {
		s.logger.Warn("PlanePartService: Part not found",
			zap.Int64("part_id", id),
		)
		return nil, errors.New(PlanePartNotFoundErr)
	}

	if req.UsageHours > part.UsageLimitHours {
		s.logger.Warn("PlanePartService: Usage hours exceeds limit",
			zap.Int64("part_id", id),
			zap.Float64("usage_hours", req.UsageHours),
			zap.Float64("limit_hours", part.UsageLimitHours),
		)
		return nil, errors.New(InvalidUsageHoursErr)
	}

	part.UsageHours = req.UsageHours

	if err := s.planePartRepo.UpdateUsage(ctx, part); err != nil {
		s.logger.Error("PlanePartService: Failed to update usage",
			zap.Int64("part_id", id),
			zap.Error(err),
		)
		return nil, fmt.Errorf("failed to update usage: %w", err)
	}

	s.logger.Info("PlanePartService: UpdatePartUsage successful",
		zap.Int64("part_id", id),
		zap.Float64("usage_percent", part.UsagePercent),
	)

	resp := part.ToResponse()
	return &resp, nil
}

func (s *PlanePartService) DeletePart(ctx context.Context, id int64) error {
	s.logger.Info("PlanePartService: DeletePart",
		zap.Int64("part_id", id),
	)

	part, err := s.planePartRepo.GetByID(ctx, id)
	if err != nil {
		s.logger.Error("PlanePartService: Failed to get part",
			zap.Int64("part_id", id),
			zap.Error(err),
		)
		return fmt.Errorf("failed to get part: %w", err)
	}
	if part == nil {
		s.logger.Warn("PlanePartService: Part not found",
			zap.Int64("part_id", id),
		)
		return errors.New(PlanePartNotFoundErr)
	}

	if err := s.planePartRepo.Delete(ctx, id); err != nil {
		s.logger.Error("PlanePartService: Failed to delete part",
			zap.Int64("part_id", id),
			zap.Error(err),
		)
		return fmt.Errorf("failed to delete part: %w", err)
	}

	s.logger.Info("PlanePartService: DeletePart successful",
		zap.Int64("part_id", id),
	)

	return nil
}

// ============= Maintenance Monitoring =============

func (s *PlanePartService) GetPartsNeedingMaintenance(ctx context.Context, thresholdPercent float64) ([]models.PlanePartResponse, error) {
	s.logger.Info("PlanePartService: GetPartsNeedingMaintenance",
		zap.Float64("threshold", thresholdPercent),
	)

	parts, err := s.planePartRepo.GetNeedingMaintenance(ctx, thresholdPercent)
	if err != nil {
		s.logger.Error("PlanePartService: Failed to get parts needing maintenance",
			zap.Error(err),
		)
		return nil, fmt.Errorf("failed to get parts: %w", err)
	}

	s.logger.Info("PlanePartService: GetPartsNeedingMaintenance successful",
		zap.Int("count", len(parts)),
	)

	responses := make([]models.PlanePartResponse, len(parts))
	for i, part := range parts {
		responses[i] = part.ToResponse()
	}

	return responses, nil
}

func (s *PlanePartService) GetPlaneWithParts(ctx context.Context, id int64) (*models.PlaneResponse, []models.PlanePartResponse, error) {
	s.logger.Info("PlanePartService: GetPlaneWithParts",
		zap.Int64("plane_id", id),
	)

	plane, err := s.planeRepo.GetByID(ctx, id)
	if err != nil {
		s.logger.Error("PlanePartService: Failed to get plane",
			zap.Int64("plane_id", id),
			zap.Error(err),
		)
		return nil, nil, fmt.Errorf("failed to get plane: %w", err)
	}
	if plane == nil {
		s.logger.Warn("PlanePartService: Plane not found",
			zap.Int64("plane_id", id),
		)
		return nil, nil, errors.New(PlaneNotFoundErrPart)
	}

	parts, err := s.planePartRepo.GetByPlaneIDWithDetails(ctx, id)
	if err != nil {
		s.logger.Error("PlanePartService: Failed to get plane parts",
			zap.Int64("plane_id", id),
			zap.Error(err),
		)
		return nil, nil, fmt.Errorf("failed to get parts: %w", err)
	}

	partResponses := make([]models.PlanePartResponse, len(parts))
	for i, part := range parts {
		partResponses[i] = part.ToResponse()
	}

	planeResp := plane.ToResponse()
	return &planeResp, partResponses, nil
}
