package permission

import (
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func RegisterModule(
	protected *echo.Group,
	db *gorm.DB,
) {
	repo := NewPermissionRepository(db)
	handler := NewPermissionHandler(repo)
	// Protected Routes
	protected.GET("/permissions", handler.GetAllPermissions)
	protected.GET("/permissions/:id", handler.GetPermissionByID)
	protected.POST("/permissions", handler.CreatePermission)
}
