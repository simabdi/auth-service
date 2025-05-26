package service

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/golang-jwt/jwt/v5"
	log "github.com/sirupsen/logrus"
	"os"
	"strconv"
	"time"
)

type JwtService interface {
	GenerateToken(ctx context.Context, uuid string) (string, error)
	VerifyToken(ctx context.Context, tokenStr string) (*jwt.Token, error)
}

type jwtService struct {
	secretKey []byte
	redis     *redis.Client
}

func NewJwtService(redis *redis.Client) JwtService {
	secretStr := os.Getenv("JWT_SECRET_KEY")
	if secretStr == "" {
		panic("JWT_SECRET_KEY is required")
	}

	secretKey, err := base64.StdEncoding.DecodeString(secretStr)
	if err != nil {
		panic("invalid base64 JWT_SECRET_KEY: " + err.Error())
	}

	return &jwtService{
		secretKey: secretKey,
		redis:     redis,
	}
}

func (j *jwtService) GenerateToken(ctx context.Context, uuid string) (string, error) {
	lifeTimeStr := os.Getenv("LIFETIME")
	lifeTime, err := strconv.Atoi(lifeTimeStr)
	if err != nil {
		return "", errors.New("invalid LIFETIME value")
	}
	ttl := time.Duration(lifeTime) * time.Second
	jti := generateTokenID()
	claims := jwt.MapClaims{
		"uuid": uuid,
		"jti":  jti,
		"exp":  time.Now().Add(ttl).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString(j.secretKey)
	if err != nil {
		return "", err
	}

	key := fmt.Sprintf("auth:session:%s", uuid)
	err = j.redis.Set(ctx, key, jti, ttl).Err()
	if err != nil {
		log.WithFields(log.Fields{
			"message": "❌ Failed",
			"error":   err,
		}).Info("Failed to save jti to Redis")
	}

	return tokenStr, nil
}

func (j *jwtService) VerifyToken(ctx context.Context, tokenStr string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return j.secretKey, nil
	})
	if err != nil || !token.Valid {
		return nil, errors.New("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid claims")
	}

	uuid, _ := claims["uuid"].(string)
	jti, _ := claims["jti"].(string)

	if uuid == "" || jti == "" {
		return nil, errors.New("missing uuid or jti")
	}

	key := fmt.Sprintf("auth:session:%s", uuid)
	savedJti, err := j.redis.Get(ctx, key).Result()
	if err != nil || savedJti != jti {
		return nil, errors.New("token revoked or invalid session")
	}

	return token, nil
}

func generateTokenID() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}
