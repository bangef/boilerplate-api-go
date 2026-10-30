package order

import (
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func RegisterModule(
	public *echo.Echo,
	db *gorm.DB,
) {
	orderRepo := NewOrderRepository(db)
	orderHandler := NewOrderHandler(orderRepo)
	public.GET("/api/orders", orderHandler.ProcessBulkOrders)
}
