package product

import (
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func RegisterModule(
	public *echo.Echo,
	protected *echo.Group,
	db *gorm.DB,
) {
	productRepo := NewProductRepository(db)
	productHandler := NewProductHandler(productRepo)
	// Public Routes
	public.GET("/api/products", productHandler.GetAllProducts)
	public.GET("/api/products/:id", productHandler.GetProductByID)

	// Protected Routes
	protected.POST("/products", productHandler.CreateProduct)
}
