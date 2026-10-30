package role

import (
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func RegisterModule(
	protected *echo.Group,
	db *gorm.DB,
) {
	roleRepo := NewRoleRepository(db)
	roleHandler := NewRoleHandler(roleRepo)
	// Protected Routes
	protected.GET("/api/roles", roleHandler.GetAllRoles)
	protected.GET("/api/roles/:id", roleHandler.GetProductByID)
	protected.POST("/roles", roleHandler.CreateRole)
}
