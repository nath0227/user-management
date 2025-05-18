package middleware

import (
	"net/http"
	"strings"
	"user-management/response"

	jwt "github.com/golang-jwt/jwt/v5"
	echo "github.com/labstack/echo/v4"
)

func AuthMiddleware(secret string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			var tokenStr string

			// 1. Try cookie
			if cookie, err := c.Cookie("token"); err == nil {
				tokenStr = cookie.Value
			} else {
				// 2. Try Authorization header
				authHeader := c.Request().Header.Get("Authorization")
				if strings.HasPrefix(authHeader, "Bearer ") {
					tokenStr = strings.TrimPrefix(authHeader, "Bearer ")
				}
			}

			if tokenStr == "" {
				return echo.NewHTTPError(response.Unauthorized().WithHTTPStatus())
			}

			// 3. Parse and verify token
			token, err := jwt.ParseWithClaims(tokenStr, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, echo.NewHTTPError(http.StatusUnauthorized, "Unexpected signing method")
				}
				return []byte(secret), nil
			})
			if err != nil || !token.Valid {
				return echo.NewHTTPError(response.Unauthorized().WithHTTPStatus())
			}

			return next(c)
		}
	}
}
