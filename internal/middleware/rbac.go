package middleware

import (
	"net/http"
	"shop_echo/pkg/response"

	"github.com/labstack/echo/v4"
)

// PermissionMiddleware mengecek apakah user memiliki PERMISSION spesifik
func PermissionMiddleware(requiredPermission ...string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			userPermissions, ok := c.Get("user_permissions").([]string)
			if !ok {
				return response.SendError(c, http.StatusForbidden, "Hak akses tidak ditemukan pada token", nil)
			}

			// Buat map pencari 0(1)
			permMap := make(map[string]bool)
			for _, p := range userPermissions {
				permMap[p] = true
			}

			// Cek apakah user memiliki setidaknya salah satu permission yang diizinkan
			for _, reqPerm := range requiredPermission {
				if permMap[reqPerm] {
					return next(c) // Akses Diizinkan
				}
			}

			return response.SendError(c, http.StatusForbidden, "Anda tidak memiliki izin (Permission) untuk melakukan aksi ini", nil)
		}
	}
}

// RoleMiddleware mengecek apakah user memiliki Role spesifik
func RoleMiddleware(requiredRoles ...string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			userRoles, ok := c.Get("user_roles").([]string)
			if !ok {
				return response.SendError(c, http.StatusForbidden, "Hak akses tidak ditemukan pada token", nil)
			}

			// Buat map pencari 0(1)
			permMap := make(map[string]bool)
			for _, p := range userRoles {
				permMap[p] = true
			}

			// Cek apakah user memiliki setidaknya salah satu role yang diizinkan
			for _, reqPerm := range requiredRoles {
				if permMap[reqPerm] {
					return next(c) // Akses Diizinkan
				}
			}

			return response.SendError(c, http.StatusForbidden, "Akses ditolak: Role anda tidak memenuhi syarat", nil)
		}
	}
}
