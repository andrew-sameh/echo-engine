package handlers

import (
	"fmt"
	"net/http"

	database "github.com/andrew-sameh/echo-engine/internal/database/db"
	"github.com/andrew-sameh/echo-engine/internal/requests"
	"github.com/andrew-sameh/echo-engine/internal/responses"
	s "github.com/andrew-sameh/echo-engine/internal/server"
	tokenservice "github.com/andrew-sameh/echo-engine/internal/services/token"
	"github.com/andrew-sameh/echo-engine/internal/utils"

	"github.com/labstack/echo/v4"

	jwtGo "github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	server *s.Server
}

func NewAuthHandler(server *s.Server) *AuthHandler {
	return &AuthHandler{server: server}
}

// Login
//
//	@Summary		Authenticate a user
//	@Description	Perform user login
//	@ID				user-login
//	@Tags			Authentication
//	@Accept			json
//	@Produce		json
//	@Param			params	body		requests.LoginRequest	true	"User's credentials"
//	@Success		200		{object}	responses.LoginResponse
//	@Failure		401		{object}	responses.Error
//	@Router			/auth/login [post]
func (authHandler *AuthHandler) Login(c echo.Context) error {
	logger := authHandler.server.Logger.Zap
	queries := authHandler.server.DB.Queries()

	loginRequest := new(requests.LoginRequest)

	if err := c.Bind(loginRequest); err != nil {
		return err
	}

	if err := authHandler.server.Echo.Validator.Validate(loginRequest); err != nil {
		logger.Error("error validating request: %v", err, c.Response().Header().Get(echo.HeaderXRequestID))
		res := responses.Response{
			Code:    http.StatusBadRequest,
			Message: "Required fields are empty or not valid",
		}
		return res.JSON(c)
	}

	// if err := loginRequest.Validate(); err != nil {
	// 	// return responses.ErrorResponse(c, http.StatusBadRequest, "Required fields are empty or not valid")
	// 	res := responses.Response{
	// 		Code:    http.StatusBadRequest,
	// 		Message: "Required fields are empty or not valid",
	// 	}
	// 	return res.JSON(c)
	// }

	user, err := queries.GetUserByEmail(c.Request().Context(), loginRequest.Email)

	if err != nil || (bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(loginRequest.Password)) != nil) {
		// return responses.ErrorResponse(c, http.StatusUnauthorized, "Invalid credentials")
		res := responses.Response{
			Code:    http.StatusUnauthorized,
			Message: "Invalid credentials",
		}
		return res.JSON(c)
	}

	tokenService := tokenservice.NewTokenService(authHandler.server.Config)
	accessToken, exp, err := tokenService.CreateAccessToken(&user)
	if err != nil {
		return err
	}
	refreshToken, err := tokenService.CreateRefreshToken(&user)
	if err != nil {
		return err
	}
	resBody := responses.NewLoginResponse(accessToken, refreshToken, exp)

	res := responses.Response{
		Code: http.StatusOK,
		Data: resBody,
	}

	return res.JSON(c)
	// return responses.Response(c, http.StatusOK, res)
}

// RefreshToken
//
//	@Summary		Refresh access token
//	@Description	Perform refresh access token
//	@ID				user-refresh
//	@Tags			Authentication
//	@Accept			json
//	@Produce		json
//	@Param			params	body		requests.RefreshRequest	true	"Refresh token"
//	@Success		200		{object}	responses.LoginResponse
//	@Failure		401		{object}	responses.Error
//	@Router			/auth/refresh [post]
func (authHandler *AuthHandler) RefreshToken(c echo.Context) error {
	queries := authHandler.server.DB.Queries()

	refreshRequest := new(requests.RefreshRequest)
	if err := c.Bind(refreshRequest); err != nil {
		return err
	}

	token, err := jwtGo.Parse(refreshRequest.Token, func(token *jwtGo.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwtGo.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(authHandler.server.Config.Auth.RefreshSecret), nil
	})

	if err != nil {
		// return responses.ErrorResponse(c, http.StatusUnauthorized, err.Error())
		res := responses.Response{
			Code:    http.StatusUnauthorized,
			Message: err.Error(),
		}
		return res.JSON(c)
	}

	claims, ok := token.Claims.(jwtGo.MapClaims)
	if !ok && !token.Valid {
		// return responses.ErrorResponse(c, http.StatusUnauthorized, "Invalid token")
		res := responses.Response{
			Code:    http.StatusUnauthorized,
			Message: "Invalid token",
		}
		return res.JSON(c)
	}

	// user := new(models.User)
	user, err := queries.GetUserById(c.Request().Context(), claims["id"].(int64))

	if err != nil {
		// return responses.ErrorResponse(c, http.StatusUnauthorized, "User not found")
		res := responses.Response{
			Code:    http.StatusUnauthorized,
			Message: "User not found",
		}
		return res.JSON(c)
	}

	tokenService := tokenservice.NewTokenService(authHandler.server.Config)
	accessToken, exp, err := tokenService.CreateAccessToken(&user)
	if err != nil {
		return err
	}
	refreshToken, err := tokenService.CreateRefreshToken(&user)
	if err != nil {
		return err
	}
	resBody := responses.NewLoginResponse(accessToken, refreshToken, exp)

	res := responses.Response{
		Code: http.StatusOK,
		Data: resBody,
	}

	return res.JSON(c)
}

// Register
//
//	@Summary		Register
//	@Description	New user registration
//	@ID				user-register
//	@Tags			Authentication
//	@Accept			json
//	@Produce		json
//	@Param			params	body		requests.RegisterRequest	true	"User's email, user's password"
//	@Success		201		{object}	responses.Data
//	@Failure		400		{object}	responses.Error
//	@Router			/auth/register [post]
func (authHandler *AuthHandler) Register(c echo.Context) error {
	// logger, _ := zap.NewProduction()
	queries := authHandler.server.DB.Queries()

	registerRequest := new(requests.RegisterRequest)

	if err := c.Bind(registerRequest); err != nil {
		return err
	}
	if err := authHandler.server.Echo.Validator.Validate(registerRequest); err != nil {
		// fmt.Errorf("error validating request: %v", err)
		return responses.Response{
			Code:    http.StatusBadRequest,
			Message: "Required fields are empty or not valid",
		}.JSON(c)

	}

	_, err := queries.GetUserByEmail(c.Request().Context(), registerRequest.Email)

	if err == nil {
		res := responses.Response{
			Code:    http.StatusBadRequest,
			Message: "User already exists",
		}
		return res.JSON(c)
	}

	encryptedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(registerRequest.Password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return err
	}
	userParams := database.CreateUserParams{
		FirstName:    registerRequest.FirstName,
		LastName:     registerRequest.LastName,
		Username:     registerRequest.Username,
		Email:        registerRequest.Email,
		Role:         "user",
		PasswordHash: string(encryptedPassword),
		CreatedAt:    utils.PgTimeNow(),
		UpdatedAt:    utils.PgTimeNow(),
	}

	newUser, err := queries.CreateUser(c.Request().Context(), userParams)

	if err != nil {
		res := responses.Response{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
		return res.JSON(c)
	}

	res := responses.Response{
		Code:    http.StatusCreated,
		Data:    newUser,
		Message: "User created successfully",
	}

	return res.JSON(c)
}
