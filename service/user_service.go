package service

import (
	"context"
	"errors"
	"github.com/simabdi/auth-service/model"
	"github.com/simabdi/auth-service/policies"
	"github.com/simabdi/auth-service/repository"
	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	Login(ctx context.Context, request policies.LoginRequest) (*model.User, error)
	GetByUuid(ctx context.Context, uuid string) (*model.User, error)
}

type userService struct {
	userRepository repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) UserService {
	return &userService{
		userRepository: userRepo,
	}
}

func (s *userService) Login(ctx context.Context, request policies.LoginRequest) (*model.User, error) {
	email := request.Email
	password := request.Password

	user, err := s.userRepository.GetByEmail(ctx, email)
	if err != nil {
		return user, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return user, err
	}

	return user, nil
}

func (s *userService) GetByUuid(ctx context.Context, uuid string) (*model.User, error) {
	if uuid == "" {
		return nil, errors.New("uuid is required")
	}
	return s.userRepository.GetByUuid(ctx, uuid)
}
