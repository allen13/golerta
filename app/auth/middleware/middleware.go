package middleware

import (
	"errors"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
)

type (
	// JWTConfig defines the config for JWT middleware.
	JWTConfig struct {
		// Skipper defines a function to skip middleware.
		Skipper Skipper

		// Signing key to validate token.
		// Required.
		SigningKey interface{} `json:"signing_key"`

		// Signing method, used to check token signing method.
		// Optional. Default value HS256.
		SigningMethod string `json:"signing_method"`

		// Context key to store user information from the token into context.
		// Optional. Default value "user".
		ContextKey string `json:"context_key"`

		// Claims are extendable claims data defining token content.
		// Optional. Default value jwt.MapClaims
		Claims jwt.Claims

		// TokenLookup is a string in the form of "<source>:<name>" that is used
		// to extract token from the request.
		// Optional. Default value "header:Authorization".
		// Possible values:
		// - "header:<name>"
		// - "query:<name>"
		// - "cookie:<name>"
		TokenLookup string `json:"token_lookup"`
	}

	jwtExtractor func(echo.Context) (string, error)

	Skipper func(c echo.Context) bool
)

const (
	bearer = "Bearer"
)

// Algorithims
const (
	AlgorithmHS256 = "HS256"
)

var (
	// DefaultJWTConfig is the default JWT auth middleware config.
	DefaultJWTConfig = JWTConfig{
		Skipper:       DefaultSkipper,
		SigningMethod: AlgorithmHS256,
		ContextKey:    "user",
		TokenLookup:   "header:" + echo.HeaderAuthorization,
		Claims:        jwt.MapClaims{},
	}
)

// JWT returns a JSON Web Token (JWT) auth middleware.
//
// For valid token, it sets the user in context and calls next handler.
// For invalid token, it returns "401 - Unauthorized" error.
// For empty token, it returns "400 - Bad Request" error.
//
// See: https://jwt.io/introduction
// See `JWTConfig.TokenLookup`
func JWT(key []byte) echo.MiddlewareFunc {
	c := DefaultJWTConfig
	c.SigningKey = key
	return JWTWithConfig(c)
}

// JWTWithConfig returns a JWT auth middleware from config.
// See: `JWT()`.
func JWTWithConfig(config JWTConfig) echo.MiddlewareFunc {
	// Defaults
	if config.Skipper == nil {
		config.Skipper = DefaultJWTConfig.Skipper
	}
	if config.SigningKey == nil {
		panic("jwt middleware requires signing key")
	}
	if config.SigningMethod == "" {
		config.SigningMethod = DefaultJWTConfig.SigningMethod
	}
	if config.ContextKey == "" {
		config.ContextKey = DefaultJWTConfig.ContextKey
	}
	if config.Claims == nil {
		config.Claims = DefaultJWTConfig.Claims
	}
	if config.TokenLookup == "" {
		config.TokenLookup = DefaultJWTConfig.TokenLookup
	}

	// Initialize
	extractor := buildExtractor(config.TokenLookup)

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if config.Skipper(c) {
				return next(c)
			}

			auth, err := extractor(c)
			if err != nil {
				return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
			}

			token, err := jwt.Parse(auth, func(t *jwt.Token) (interface{}, error) {
				return config.SigningKey, nil
			})

			if err == nil && token.Valid {
				// Store user information from token into context.
				c.Set(config.ContextKey, token)
				return next(c)
			}
			return echo.ErrUnauthorized
		}
	}
}

func buildExtractor(tokenLookups string) jwtExtractor {
	var extractors []jwtExtractor

	lookups := strings.Split(tokenLookups, ",")
	for _, lookup := range lookups {
		parts := strings.Split(lookup, ":")

		switch parts[0] {
		case "header":
			extractors = append(extractors, jwtFromHeader(parts[1]))
		case "query":
			extractors = append(extractors, jwtFromQuery(parts[1]))
		case "cookie":
			extractors = append(extractors, jwtFromCookie(parts[1]))
		}
	}
	if len(extractors) == 1 {
		return extractors[0]
	}

	return jwtFromExtractors(extractors)
}

// jwtFromHeader returns a `jwtExtractor` that extracts token from request header.
func jwtFromHeader(header string) jwtExtractor {
	return func(c echo.Context) (string, error) {
		auth := c.Request().Header.Get(header)
		l := len(bearer)
		if len(auth) > l+1 && auth[:l] == bearer {
			return auth[l+1:], nil
		}
		return "", errors.New("empty or invalid jwt in request header. ")
	}
}

// jwtFromQuery returns a `jwtExtractor` that extracts token from query string.
func jwtFromQuery(param string) jwtExtractor {
	return func(c echo.Context) (string, error) {
		token := c.QueryParam(param)
		var err error
		if token == "" {
			return "", errors.New("empty jwt in query string. ")
		}
		return token, err
	}
}

// jwtFromCookie returns a `jwtExtractor` that extracts token from named cookie.
func jwtFromCookie(name string) jwtExtractor {
	return func(c echo.Context) (string, error) {
		cookie, err := c.Cookie(name)
		if err != nil {
			return "", errors.New("empty jwt in cookie. ")
		}
		return cookie.Value, nil
	}
}

// jwtFromExtractors returns a `jwtExtractor` that extracts token from header, query, or cookie.
func jwtFromExtractors(extractors []jwtExtractor) jwtExtractor {
	return func(c echo.Context) (string, error) {
		extractorErrors := ""
		for _, extractor := range extractors {
			token, err := extractor(c)
			if err != nil {
				extractorErrors += err.Error()
			}
			if token != "" {
				return token, nil
			}
		}
		return "", errors.New(extractorErrors)
	}
}

func DefaultSkipper(echo.Context) bool {
	return false
}
