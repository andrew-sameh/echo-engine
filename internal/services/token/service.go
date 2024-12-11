package token

import (
	"github.com/andrew-sameh/echo-engine/internal/config"
	database "github.com/andrew-sameh/echo-engine/internal/database/db"
	"github.com/golang-jwt/jwt/v5"
)

const ExpireCount = 2
const ExpireRefreshCount = 168

type JwtCustomClaims struct {
	ID    int64  `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Role  string `json:"role"`
	jwt.RegisteredClaims
}

type JwtCustomRefreshClaims struct {
	ID int64 `json:"id"`
	jwt.RegisteredClaims
}

type ServiceWrapper interface {
	CreateAccessToken(user *database.User) (accessToken string, exp int64, err error)
	CreateRefreshToken(user *database.User) (t string, err error)
}

type Service struct {
	config *config.Config
}

func NewTokenService(cfg *config.Config) *Service {
	return &Service{
		config: cfg,
	}
}
