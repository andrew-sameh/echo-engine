package token

import (
	"time"

	database "github.com/andrew-sameh/echo-engine/internal/database/db"
	"github.com/golang-jwt/jwt/v5"
)

func (tokenService *Service) CreateAccessToken(user *database.User) (t string, expired int64, err error) {
	exp := time.Now().Add(time.Hour * ExpireCount)
	claims := &JwtCustomClaims{
		user.ID,
		user.FirstName + " " + user.LastName,
		user.Email,
		string(user.Role),
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(exp),
		},
	}
	expired = exp.Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, err = token.SignedString([]byte(tokenService.config.Auth.AccessSecret))
	if err != nil {
		return
	}

	return
}
