package handlers

import (
	"net/http"

	"github.com/andrew-sameh/echo-engine/internal/responses"
	s "github.com/andrew-sameh/echo-engine/internal/server"
	"github.com/andrew-sameh/echo-engine/internal/services/token"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

type UserHandler struct {
	server *s.Server
}

func NewUserHandler(server *s.Server) *UserHandler {
	return &UserHandler{server: server}
}

// ListUsersHandler lists all existing users
//
//	@Summary		List users
//	@Description	get users
//	@Tags			User
//	@Accept			json
//	@Produce		json
//	@Success		200	{array}		[]map[string]string
//	@Failure		500	{object}	error
//	@Security		ApiKeyAuth
//
//	@Router			/users [get]
func (g *UserHandler) ListUsersHandler(c echo.Context) error {

	// user := c.Get("user").(*jwt.Token)
	// claims := user.Claims.(*token.JwtCustomClaims)
	// id := claims.ID
	queries := g.server.DB.Queries()
	users, err := queries.GetAllUsers(c.Request().Context())
	if err != nil {
		res := responses.Response{
			Code:    http.StatusInternalServerError,
			Message: err,
		}
		return res.JSON(c)
	}

	res := responses.Response{
		Code: http.StatusOK,
		Data: users,
	}
	return res.JSON(c)
}

// Get My User
//
//	@Summary		Get my user
//	@Description	get my user
//	@Tags			User
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	map[string]string
//	@Failure		500	{object}	error
//	@Security		ApiKeyAuth
//
//	@Router			/profile [get]
func (g *UserHandler) GetMyUserHandler(c echo.Context) error {
	userToken := c.Get("user").(*jwt.Token)
	claims := userToken.Claims.(*token.JwtCustomClaims)
	id := claims.ID
	queries := g.server.DB.Queries()
	user, err := queries.GetUserById(c.Request().Context(), id)
	if err != nil {
		res := responses.Response{
			Code:    http.StatusInternalServerError,
			Message: err,
		}
		return res.JSON(c)
	}

	res := responses.Response{
		Code:   http.StatusOK,
		Pretty: true,
		Data:   user,
	}
	return res.JSON(c)
}
