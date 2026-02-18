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
		"name", req.Name,
	)

	existing, err := s.repo.GetByName(ctx, req.Name)
	if err != nil {
		s.logger.Error("UserService: Failed to check existing user",
			"name", req.Name,
			"error", err,
		)
		return nil, fmt.Errorf("failed to check existing user: %w", err)
	}
	if existing != nil {
		s.logger.Warn("UserService: User already exists",
			"name", req.Name,
		)
		return nil, errors.New(UserExistsErr)
	}

	hashedPassword, err := util.HashPassword(req.Password)
	if err != nil {
		s.logger.Error("UserService: Failed to hash password",
			"error", err,
		)
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	user := &models.User{
		Name:     req.Name,
		Password: hashedPassword,
		Role:     req.Role,
	}

	// Default to "user" if role is not provided or invalid
	if user.Role == "" {
		user.Role = "user"
	}

	if err := s.repo.Create(ctx, user); err != nil {
		s.logger.Error("UserService: Failed to create user",
			"name", req.Name,
			"error", err,
		)
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	s.logger.Info("UserService: User registered successfully",
		"user_id", user.ID,
		"name", user.Name,
	)

	resp := user.ToResponse()
	return &resp, nil
}

func (s *UserService) Login(ctx context.Context, req *models.LoginRequest) (*LoginResponse, error) {
	s.logger.Info("UserService: Login attempt",
		"name", req.Name,
	)

	user, err := s.repo.GetByName(ctx, req.Name)
	if err != nil {
		s.logger.Error("UserService: Failed to find user",
			"name", req.Name,
			"error", err,
		)
		return nil, fmt.Errorf("failed to find user: %w", err)
	}
	if user == nil {
		s.logger.Warn("UserService: User not found",
			"name", req.Name,
		)
		return nil, errors.New(UserNotFoundErr)
	}

	if !util.CheckPassword(req.Password, user.Password) {
		s.logger.Warn("UserService: Invalid password",
			"name", req.Name,
		)
		return nil, errors.New(InvalidPasswordErr)
	}

	token, err := s.jwtSvc.GenerateToken(user.ID, user.Name, user.Role)
	if err != nil {
		s.logger.Error("UserService: Failed to generate token",
			"user_id", user.ID,
			"error", err,
		)
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	s.logger.Info("UserService: Login successful",
		"user_id", user.ID,
		"name", user.Name,
		"role", user.Role,
	)

	return &LoginResponse{
		User:  user.ToResponse(),
		Token: token,
	}, nil
}

func (s *UserService) GetByID(ctx context.Context, id int64) (*models.UserResponse, error) {
	s.logger.Info("UserService: GetByID",
		"user_id", id,
	)

	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		s.logger.Error("UserService: Failed to get user",
			"user_id", id,
			"error", err,
		)
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	if user == nil {
		s.logger.Warn("UserService: User not found",
			"user_id", id,
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
			"error", err,
		)
		return nil, fmt.Errorf("failed to get users: %w", err)
	}

	s.logger.Info("UserService: GetAll successful",
		"count", len(users),
	)

	responses := make([]models.UserResponse, len(users))
	for i, user := range users {
		responses[i] = user.ToResponse()
	}

	return responses, nil
}

func (s *UserService) Update(ctx context.Context, id int64, req *models.UpdateRequest) (*models.UserResponse, error) {
	s.logger.Info("UserService: Update",
		"user_id", id,
	)

	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		s.logger.Error("UserService: Failed to get user",
			"user_id", id,
			"error", err,
		)
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	if user == nil {
		s.logger.Warn("UserService: User not found",
			"user_id", id,
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
				"error", err,
			)
			return nil, fmt.Errorf("failed to hash password: %w", err)
		}
		user.Password = hashedPassword
	}

	if err := s.repo.Update(ctx, user); err != nil {
		s.logger.Error("UserService: Failed to update user",
			"user_id", id,
			"error", err,
		)
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	s.logger.Info("UserService: Update successful",
		"user_id", id,
	)

	resp := user.ToResponse()
	return &resp, nil
}

func (s *UserService) Delete(ctx context.Context, id int64) error {
	s.logger.Info("UserService: Delete",
		"user_id", id,
	)

	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		s.logger.Error("UserService: Failed to get user",
			"user_id", id,
			"error", err,
		)
		return fmt.Errorf("failed to get user: %w", err)
	}
	if user == nil {
		s.logger.Warn("UserService: User not found",
			"user_id", id,
		)
		return errors.New(UserNotFoundErr)
	}

	if err := s.repo.Delete(ctx, id); err != nil {
		s.logger.Error("UserService: Failed to delete user",
			"user_id", id,
			"error", err,
		)
		return fmt.Errorf("failed to delete user: %w", err)
	}

	s.logger.Info("UserService: Delete successful",
		"user_id", id,
	)

	return nil
}

func (s *UserService) GetMe(ctx context.Context, userID int64) (*models.UserResponse, error) {
	s.logger.Info("UserService: GetMe",
		"user_id", userID,
	)

	user, err := s.repo.GetByID(ctx, userID)
	if err != nil {
		s.logger.Error("UserService: Failed to get user",
			"user_id", userID,
			"error", err,
		)
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	if user == nil {
		s.logger.Warn("UserService: User not found",
			"user_id", userID,
		)
		return nil, errors.New(UserNotFoundErr)
	}

	s.logger.Info("UserService: GetMe successful",
		"user_id", userID,
		"name", user.Name,
		"role", user.Role,
	)

	resp := user.ToResponse()
	return &resp, nil
}
