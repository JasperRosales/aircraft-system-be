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
}

func NewUserService(repo *repository.UserRepository, jwtSvc *JWTService) *UserService {
	return &UserService{repo: repo, jwtSvc: jwtSvc}
}

type LoginResponse struct {
	User  models.UserResponse `json:"user"`
	Token string              `json:"token"`
}

func (s *UserService) Register(ctx context.Context, req *models.RegisterRequest) (*models.UserResponse, error) {
	existing, err := s.repo.GetByName(ctx, req.Name)
	if err != nil {
		return nil, fmt.Errorf("failed to check existing user: %w", err)
	}
	if existing != nil {
		return nil, errors.New(UserExistsErr)
	}

	hashedPassword, err := util.HashPassword(req.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	user := &models.User{
		Name:     req.Name,
		Password: hashedPassword,
		Role:     "user",
	}

	if err := s.repo.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	resp := user.ToResponse()
	return &resp, nil
}

func (s *UserService) Login(ctx context.Context, req *models.LoginRequest) (*LoginResponse, error) {
	user, err := s.repo.GetByName(ctx, req.Name)
	if err != nil {
		return nil, fmt.Errorf("failed to find user: %w", err)
	}
	if user == nil {
		return nil, errors.New(UserNotFoundErr)
	}

	if !util.CheckPassword(req.Password, user.Password) {
		return nil, errors.New(InvalidPasswordErr)
	}

	token, err := s.jwtSvc.GenerateToken(user.ID, user.Name, user.Role)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	return &LoginResponse{
		User:  user.ToResponse(),
		Token: token,
	}, nil
}

func (s *UserService) GetByID(ctx context.Context, id int64) (*models.UserResponse, error) {
	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	if user == nil {
		return nil, errors.New(UserNotFoundErr)
	}

	resp := user.ToResponse()
	return &resp, nil
}

func (s *UserService) GetAll(ctx context.Context) ([]models.UserResponse, error) {
	users, err := s.repo.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get users: %w", err)
	}

	responses := make([]models.UserResponse, len(users))
	for i, user := range users {
		responses[i] = user.ToResponse()
	}

	return responses, nil
}

func (s *UserService) Update(ctx context.Context, id int64, req *models.UpdateRequest) (*models.UserResponse, error) {
	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	if user == nil {
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
			return nil, fmt.Errorf("failed to hash password: %w", err)
		}
		user.Password = hashedPassword
	}

	if err := s.repo.Update(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	resp := user.ToResponse()
	return &resp, nil
}

func (s *UserService) Delete(ctx context.Context, id int64) error {
	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}
	if user == nil {
		return errors.New(UserNotFoundErr)
	}

	if err := s.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	return nil
}
