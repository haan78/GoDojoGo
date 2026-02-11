package service

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	data "GoDojoGo/data"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v5"

	globals "GoDojoGo/deff"
)

type loginRequestType struct {
	User     string `json:"user"`
	Password string `json:"password"`
}

type tokenType struct {
	Token    string `json:"token"`
	Duration int16  `json:"duration"`
	Role     string `json:"role"`
	Name     string `json:"name"`
}

type CustomClaims struct {
	UserId int64  `json:"user_id"`
	Name   string `json:"name"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

const ClaimsKey = "claims"

type HandlerFuncType = func(next echo.HandlerFunc) echo.HandlerFunc

func createCustomClaims(ud *data.GetUserType) CustomClaims {
	return CustomClaims{
		UserId: ud.UserId,
		Name:   ud.Name,
		Role:   ud.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
}

func generateToken(udata *data.GetUserType) (string, error) {

	/*if err := godotenv.Load(); err != nil {
		return "", err
	}*/

	jwtSecret := globals.Settings.JWT_SECRET
	if jwtSecret == "" {
		return "", errors.New("no connection string in .env file")
	}

	claims := createCustomClaims(udata)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(jwtSecret))
}

func validateTokenString(tokenString string) (*CustomClaims, error) {
	claims := &CustomClaims{}
	jwtSecret := globals.Settings.JWT_SECRET
	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		// Enforce expected signing method to prevent alg=none or confusion attacks
		if t.Method != jwt.SigningMethodHS256 {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(jwtSecret), nil
	})

	if err == nil {
		if token.Valid {
			return claims, nil
		} else {
			return nil, RaiseServiceError(401, "invalid token")
		}
	} else {
		return nil, err
	}
}

func getTokenString(r *http.Request) (string, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader != "" {
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) == 2 && strings.EqualFold(parts[0], "Bearer") {
			tokenString := strings.TrimSpace(parts[1])
			if tokenString != "" {
				return tokenString, nil
			} else {
				return "", RaiseServiceError(401, "empty bearer token")
			}
		} else {
			return "", RaiseServiceError(401, "invalid Authorization header format")
		}
	} else {
		return "", RaiseServiceError(401, "missing Authorization header")
	}
}

func tokenValidate(r *http.Request) (*CustomClaims, error) {
	token, err := getTokenString(r)
	if err == nil {
		claims, err := validateTokenString(token)
		if err == nil {
			return claims, nil
		} else {
			return nil, err
		}
	} else {
		return nil, err
	}
}

func tokenAuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c *echo.Context) error {
		claims, err := tokenValidate(c.Request())
		if err != nil {
			// you can return your own structured error here if you want
			return RaiseServiceError(http.StatusUnauthorized, "invalid or missing token")
		}

		// store claims for handlers to optionally use later
		c.Set(ClaimsKey, claims)
		return next(c)
	}
}

func noAuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c *echo.Context) error {
		claims := createCustomClaims(&data.GetUserType{
			UserId: 0,
			Name:   "Nobody",
			EMail:  "no@email.com",
			Role:   "Public",
		})
		c.Set(ClaimsKey, claims)
		return next(c)
	}
}

func GetSecLevel(level int) HandlerFuncType {
	if level == 1 {
		return tokenAuthMiddleware
	} else {
		return noAuthMiddleware
	}
}

func CreateTokenReq(c *echo.Context) error {
	var req loginRequestType
	if err := c.Bind(&req); err != nil {
		return RaiseServiceError(http.StatusBadRequest, "invalid request body")
	}

	user, err := data.GetUser(req.User, req.Password, false)

	if err != nil {
		return RaiseServiceError(http.StatusBadRequest, "wrong credential: "+err.Error())
	}

	tokenString, err := generateToken(user)
	if err != nil {
		return RaiseServiceError(http.StatusInternalServerError, "failed to generate token")
	}

	result := &tokenType{
		Token:    tokenString,
		Duration: 3600,
		Role:     user.Role,
		Name:     user.Name,
	}

	return c.JSON(http.StatusOK, result)
}
