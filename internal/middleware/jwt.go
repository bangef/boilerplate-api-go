package middleware

import (
	"fmt"
	"net/http"
	"shop_echo/internal/domain"
	"shop_echo/pkg/response"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

func JWTMiddleware(jwtSecret string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return response.SendError(c, http.StatusUnauthorized, "Header Authorization tidak ditemukan", nil)
			}

			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				return response.SendError(c, http.StatusUnauthorized, "Format token harus 'Bearer <token>'", nil)
			}

			tokenString := parts[1]
			token, err := jwt.ParseWithClaims(tokenString, &domain.JWTCustomClaims{}, func(t *jwt.Token) (any, error) {
				if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("method signing tidak valid")
				}
				return []byte(jwtSecret), nil
			})

			if err != nil || !token.Valid {
				return response.SendError(c, http.StatusUnauthorized, "Token JWT tidak valid atau sudah expired", nil)
			}

			claims, ok := token.Claims.(*domain.JWTCustomClaims)
			if ok {
				c.Set("user_id", claims.UserID)
				c.Set("user_email", claims.Email)
			}

			return next(c)
		}
	}
}
