package service

import (
	"context"
	"errors"
	"time"

	"github.com/jakabrajadenta/go-explore-api/internal/dto"
	"github.com/jakabrajadenta/go-explore-api/internal/model"
	"github.com/jakabrajadenta/go-explore-api/internal/repository"
	"github.com/jakabrajadenta/go-explore-api/pkg/logger"
)

var (
	ErrNotFound         = errors.New("user not found")
	ErrUsernameConflict = errors.New("username already taken")
	ErrEmailConflict    = errors.New("email already registered")
)

type UserService interface {
	GetAll(ctx context.Context, page, perPage int) ([]dto.UserResponse, int, error)
	GetByID(ctx context.Context, id int64) (*dto.UserResponse, error)
	Create(ctx context.Context, req dto.CreateUserRequest) (*dto.UserResponse, error)
	Update(ctx context.Context, id int64, req dto.UpdateUserRequest) (*dto.UserResponse, error)
	Delete(ctx context.Context, id int64) error
}

type userService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{repo: repo}
}

func (s *userService) GetAll(ctx context.Context, page, perPage int) ([]dto.UserResponse, int, error) {
	log := logger.FromCtx(ctx)
	log.Info("service.GetAll", "page", page, "per_page", perPage)

	if page < 1 {
		page = 1
	}
	if perPage < 1 || perPage > 100 {
		perPage = 10
	}

	users, total, err := s.repo.FindAll(ctx, page, perPage)
	if err != nil {
		log.Error("service.GetAll failed", "error", err)
		return nil, 0, err
	}

	result := make([]dto.UserResponse, len(users))
	for i, u := range users {
		result[i] = toResponse(u)
	}
	return result, total, nil
}

func (s *userService) GetByID(ctx context.Context, id int64) (*dto.UserResponse, error) {
	log := logger.FromCtx(ctx)
	log.Info("service.GetByID", "user_id", id)

	user, err := s.repo.FindByID(ctx, id)
	if errors.Is(err, repository.ErrNotFound) {
		return nil, ErrNotFound
	}
	if err != nil {
		log.Error("service.GetByID failed", "user_id", id, "error", err)
		return nil, err
	}
	resp := toResponse(*user)
	return &resp, nil
}

func (s *userService) Create(ctx context.Context, req dto.CreateUserRequest) (*dto.UserResponse, error) {
	log := logger.FromCtx(ctx)
	log.Info("service.Create", "username", req.Username, "email", req.Email)

	if _, err := s.repo.FindByUsername(ctx, req.Username); err == nil {
		return nil, ErrUsernameConflict
	}
	if _, err := s.repo.FindByEmail(ctx, req.Email); err == nil {
		return nil, ErrEmailConflict
	}

	user, err := s.repo.Create(ctx, model.User{
		Username: req.Username,
		Email:    req.Email,
		FullName: req.FullName,
		Phone:    req.Phone,
	})
	if err != nil {
		log.Error("service.Create failed", "username", req.Username, "error", err)
		return nil, err
	}
	resp := toResponse(*user)
	return &resp, nil
}

func (s *userService) Update(ctx context.Context, id int64, req dto.UpdateUserRequest) (*dto.UserResponse, error) {
	log := logger.FromCtx(ctx)
	log.Info("service.Update", "user_id", id, "username", req.Username)

	existing, err := s.repo.FindByID(ctx, id)
	if errors.Is(err, repository.ErrNotFound) {
		return nil, ErrNotFound
	}
	if err != nil {
		log.Error("service.Update failed", "user_id", id, "error", err)
		return nil, err
	}

	if req.Username != existing.Username {
		if _, err := s.repo.FindByUsername(ctx, req.Username); err == nil {
			return nil, ErrUsernameConflict
		}
	}
	if req.Email != existing.Email {
		if _, err := s.repo.FindByEmail(ctx, req.Email); err == nil {
			return nil, ErrEmailConflict
		}
	}

	existing.Username = req.Username
	existing.Email = req.Email
	existing.FullName = req.FullName
	existing.Phone = req.Phone
	if req.IsActive != nil {
		existing.IsActive = *req.IsActive
	}

	updated, err := s.repo.Update(ctx, *existing)
	if err != nil {
		log.Error("service.Update failed", "user_id", id, "error", err)
		return nil, err
	}
	resp := toResponse(*updated)
	return &resp, nil
}

func (s *userService) Delete(ctx context.Context, id int64) error {
	log := logger.FromCtx(ctx)
	log.Info("service.Delete", "user_id", id)

	err := s.repo.Delete(ctx, id)
	if errors.Is(err, repository.ErrNotFound) {
		return ErrNotFound
	}
	if err != nil {
		log.Error("service.Delete failed", "user_id", id, "error", err)
	}
	return err
}

func toResponse(u model.User) dto.UserResponse {
	return dto.UserResponse{
		ID:        u.ID,
		Username:  u.Username,
		Email:     u.Email,
		FullName:  u.FullName,
		Phone:     u.Phone,
		IsActive:  u.IsActive,
		CreatedAt: u.CreatedAt.Format(time.RFC3339),
		UpdatedAt: u.UpdatedAt.Format(time.RFC3339),
	}
}
