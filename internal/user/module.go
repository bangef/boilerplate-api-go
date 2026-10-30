// package user

// import (
// 	"github.com/labstack/echo/v4"
// 	"github.com/redis/go-redis/v9"
// 	"gorm.io/gorm"
// )

// func RegisterModule(
// 	public *echo.Echo,
// 	protected *echo.Group,
// 	db *gorm.DB,
// 	rdb *redis.Client,
// 	jwtSecret string,
// ) {
// 	userRepo := NewUserRepository(db)
// 	tokenRepo := NewTokenRepository(rdb)
// 	userHandler := NewAuthHandler(userRepo, tokenRepo, jwtSecret)
// 	// Public Routes
// 	public.POST("/api/register", userHandler.Register)
// 	public.POST("/api/login", userHandler.Login)
// 	public.POST("/api/refresh", userHandler.RefreshToken)
// 	public.POST("/api/logout", userHandler.Logout)
// 	// Protected Routes
// 	protected.POST("/logout-all", userHandler.LogoutAll)
// }
