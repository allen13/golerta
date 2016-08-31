package controllers

import (
	"github.com/allen13/golerta/app/auth"
	"github.com/allen13/golerta/app/models"
	"github.com/labstack/echo"
	"net/http"
)

type AuthController struct {
	Echo         *echo.Echo
	AuthProvider auth.AuthProvider
}

func (ac *AuthController) Init() {
	ac.Echo.Post("/auth/login", ac.LoginHandler)
}

// Handles login request
func (ac *AuthController) LoginHandler(ctx echo.Context) error {
	var loginRequest models.LoginRequest
	err := ctx.Bind(&loginRequest)

	if err != nil || loginRequest.Username == "" || loginRequest.Password == "" {
		return ctx.JSON(http.StatusUnauthorized, models.ErrorResponse("Invalid login request"))
	}

	loginSuccess, token, err := ac.AuthProvider.Authenticate(loginRequest.Username, loginRequest.Password)

	if err != nil || !loginSuccess {
		return ctx.JSON(http.StatusUnauthorized, models.ErrorResponse("Login failed"))
	}

	authToken := models.AuthToken{token}
	return ctx.JSON(http.StatusOK, authToken)
}
