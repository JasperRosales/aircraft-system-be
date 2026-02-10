package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/JasperRosales/aircraft-system-be/internal/models"
	"github.com/JasperRosales/aircraft-system-be/internal/repository"
	"github.com/JasperRosales/aircraft-system-be/internal/util"
	"go.uber.org/zap"
)

const (
	UserNotFoundErr    = "user not found"
	UserExistsErr      = "user already exists"
	InvalidPasswordErr = "invalid password"
)

type UserService struct {
	repo   *repository.UserRepository
	jwtSvc *JWTService
	logger *util.Logger
}

func NewUserService(repo *repository.UserRepository, jwtSvc *JWTService, logger *util.Logger) *UserService {
	return &UserService{repo: repo, jwtSvc: jwtSvc, logger: logger}
}

type LoginResponse struct {
	User  models.UserResponse `json:"user"`
	Token string              `json:"token"`
}

func (s *UserService) Register(ctx context.Context, req *models.RegisterRequest) (*models.UserResponse, error) {
	s.logger.Info("UserService: Registering new user",
		zap.String("name", req.Name),
	)

	existing, err := s.repo.GetByName(ctx, req.Name)
	if err != nil {
		s.logger.Error("UserService: Failed to check existing user",
			zap.String("name", req.Name),
			zap.Error(err),
		)
		return nil, fmt.Errorf("failed to check existing user: %w", err)
	}
	if existing != nil {
		s.logger.Warn("UserService: User already exists",
			zap.String("name", req.Name),
		)
		return nil, errors.New(UserExistsErr)
	}

	hashedPassword, err := util.HashPassword(req.Password)
	if err != nil {
		s.logger.Error("UserService: Failed to hash password",
			zap.Error(err),
		)
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	user := &models.User{
		Name:     req.Name,
		Password: hashedPassword,
		Role:     "user",
	}

	if err := s.repo.Create(ctx, user); err != nil {
		s.logger.Error("UserService: Failed to create user",
			zap.String("name", req.Name),
			zap.Error(err),
		)
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	s.logger.Info("UserService: User registered successfully",
		zap.Int64("user_id", user.ID),
		zap.String("name", user.Name),
	)

	resp := user.ToResponse()
	return &resp, nil
}

func (s *UserService) Login(ctx context.Context, req *models.LoginRequest) (*LoginResponse, error) {
	s.logger.Info("UserService: Login attempt",
		zap.String("name", req.Name),
	)

	user, err := s.repo.GetByName(ctx, req.Name)
	if err != nil {
		s.logger.Error("UserService: Failed to find user",
			zap.String("name", req.Name),
			zap.Error(err),
		)
		return nil, fmt.Errorf("failed to find user: %w", err)
	}
	if user == nil {
		s.logger.Warn("UserService: User not found",
			zap.String("name", req.Name),
		)
		return nil, errors.New(UserNotFoundErr)
	}

	if !util.CheckPassword(req.Password, user.Password) {
		s.logger.Warn("UserService: Invalid password",
			zap.String("name", req.Name),
		)
		return nil, errors.New(InvalidPasswordErr)
	}

	token, err := s.jwtSvc.GenerateToken(user.ID, user.Name, user.Role)
	if err != nil {
		s.logger.Error("UserService: Failed to generate token",
			zap.Int64("user_id", user.ID),
			zap.Error(err),
		)
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	s.logger.Info("UserService: Login successful",
		zap.Int64("user_id", user.ID),
		zap.String("name", user.Name),
		zap.String("role", user.Role),
	)

	return &LoginResponse{
		User:  user.ToResponse(),
		Token: token,
	}, nil
}

func (s *UserService) GetByID(ctx context.Context, id int64) (*models.UserResponse, error) {
	s.logger.Info("UserService: GetByID",
		zap.Int64("user_id", id),
	)

	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		s.logger.Error("UserService: Failed to get user",
			zap.Int64("user_id", id),
			zap.Error(err),
		)
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	if user == nil {
		s.logger.Warn("UserService: User not found",
			zap.Int64("user_id", id),
		)
		return nil, errors.New(UserNotFoundErr)
	}

	resp := user.ToResponse()
	return &resp, nil
}

func (s *UserService) GetAll(ctx context.Context) ([]models.UserResponse, error) {
	s.logger.Info("UserService: GetAll")

	users, err := s.repo.GetAll(ctx)
	if err != nil {
		s.logger.Error("UserService: Failed to get users",
			zap.Error(err),
		)
		return nil, fmt.Errorf("failed to get users: %w", err)
	}

	s.logger.Info("UserService: GetAll successful",
		zap.Int("count", len(users)),
	)

	responses := make([]models.UserResponse, len(users))
	for i, user := range users {
		responses[i] = user.ToResponse()
	}

	return responses, nil
}

func (s *UserService) Update(ctx context.Context, id int64, req *models.UpdateRequest) (*models.UserResponse, error) {
	s.logger.Info("UserService: Update",
		zap.Int64("user_id", id),
	)

	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		s.logger.Error("UserService: Failed to get user",
			zap.Int64("user_id", id),
			zap.Error(err),
		)
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	if user == nil {
		s.logger.Warn("UserService: User not found",
			zap.Int64("user_id", id),
		)
		return nil, errors.New(UserNotFoundErr)
	}

	if req.Name != "" {
		user.Name = req.Name
	}
	if req.Role != "" {
		user.Role = req.Role
	}
	if req.Password != "" {
		hashedPassword, err := util.HashPassword(req.Password)
		if err != nil {
			s.logger.Error("UserService: Failed to hash password",
				zap.Error(err),
			)
			return nil, fmt.Errorf("failed to hash password: %w", err)
		}
		user.Password = hashedPassword
	}

	if err := s.repo.Update(ctx, user); err != nil {
		s.logger.Error("UserService: Failed to update user",
			zap.Int64("user_id", id),
			zap.Error(err),
		)
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	s.logger.Info("UserService: Update successful",
		zap.Int64("user_id", id),
	)

	resp := user.ToResponse()
	return &resp, nil
}

func (s *UserService) Delete(ctx context.Context, id int64) error {
	s.logger.Info("UserService: Delete",
		zap.Int64("user_id", id),
	)

	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		s.logger.Error("UserService: Failed to get user",
			zap.Int64("user_id", id),
			zap.Error(err),
		)
		return fmt.Errorf("failed to get user: %w", err)
	}
	if user == nil {
		s.logger.Warn("UserService: User not found",
			zap.Int64("user_id", id),
		)
		return errors.New(UserNotFoundErr)
	}

	if err := s.repo.Delete(ctx, id); err != nil {
		s.logger.Error("UserService: Failed to delete user",
			zap.Int64("user_id", id),
			zap.Error(err),
		)
		return fmt.Errorf("failed to delete user: %w", err)
	}

	s.logger.Info("UserService: Delete successful",
		zap.Int64("user_id", id),
	)

	return nil
}
