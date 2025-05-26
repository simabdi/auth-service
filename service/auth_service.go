package service

import (
	"context"
	"errors"
	"github.com/golang-jwt/jwt/v5"
)

type AuthInfo struct {
	UserID  uint
	Uuid    string
	RefType string
	RefID   uint
}

type AuthService interface {
	GenerateToken(ctx context.Context, uuid string) (string, error)
	ParseToken(ctx context.Context, token string) (*AuthInfo, error)
}

type authService struct {
	jwtService  JwtService
	userService UserService
}

func NewAuthService(jwt JwtService, user UserService) AuthService {
	return &authService{
		jwtService:  jwt,
		userService: user,
	}
}

func (a *authService) GenerateToken(ctx context.Context, uuid string) (string, error) {
	return a.jwtService.GenerateToken(ctx, uuid)
}

func (a *authService) ParseToken(ctx context.Context, token string) (*AuthInfo, error) {
	validToken, err := a.jwtService.VerifyToken(ctx, token)
	if err != nil {
		return nil, err
	}

	claims, ok := validToken.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid token claims")
	}

	uuidStr, ok := claims["uuid"].(string)
	if !ok {
		return nil, errors.New("uuid not found in token")
	}

	// cari user dari UUID
	user, err := a.userService.GetByUuid(ctx, uuidStr)
	if err != nil {
		return nil, err
	}

	return &AuthInfo{
		UserID:  user.ID,
		Uuid:    user.Uuid,
		RefType: user.ReferenceType,
		RefID:   user.ReferenceID,
	}, nil
}
