package controllers

import (
	"net/http"

	"github.com/allen13/golerta/app/auth"
	"github.com/allen13/golerta/app/models"
	"github.com/labstack/echo"
)

// AuthController for the /auth endpoint
type AuthController struct {
	Echo         *echo.Echo
	AuthProvider auth.AuthProvider
}

// Init the controller
func (ac *AuthController) Init() {
	ac.Echo.POST("/auth/login", ac.LoginHandler)
}

// LoginHandler handles login request
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

	authToken := models.AuthToken{Token: token}
	return ctx.JSON(http.StatusOK, authToken)
}
