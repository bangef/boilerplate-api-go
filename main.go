package main

import (
	"fmt"
	"log"
	"os"
	"shop_echo/config"
	"shop_echo/internal/domain"
	appMiddleware "shop_echo/internal/middleware"
	"shop_echo/internal/order"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	// 1. Load file env
	if err := godotenv.Load(); err != nil {
		log.Println("Peringatan: file .env tidak ditemukan, menggunakan environment variable sistem")
	}

	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8081" // default jika kosong
	}
	// jwtSecret := os.Getenv("JWT_SECRET")

	// 2. Koneksi DB & Automigrate
	db := config.ConnectDB()
	listTable := []any{&domain.User{}, &domain.Order{}}
	if err := db.AutoMigrate(listTable...); err != nil {
		log.Fatal("Gagal melakukan AutoMigrate:", err)
	}
	// rdb := config.ConnectRedis() // Koneksi Redis

	// 4. Inisialisasi Echo & Custom Validator
	e := echo.New()
	e.Validator = &appMiddleware.CustomValidator{Validator: validator.New()}

	// 5. Global Middlewares
	e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		// 1. Aktifkan data pencatatan yang Anda butuhkan
		LogURI:     true,
		LogStatus:  true,
		LogMethod:  true,
		LogLatency: true,

		// 2. Tulis fungsi untuk mencetak log sesuai format kustom Anda
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			// Mengambil waktu format RFC3339 saat ini
			timeRFC3339 := time.Now().Format(time.RFC3339)

			// Mencetak format sesuai dengan string bawaan Anda:
			// "[${time_rfc3339}] ({method}) {uri} -> Status: ({status} | Latency: ){latency_human}"
			fmt.Printf("[%s] (%s) %s -> Status: (%d | Latency: %s)\n",
				timeRFC3339,
				v.Method,
				v.URI,
				v.Status,
				v.Latency.String(), // Menggunakan .String() untuk format yang mudah dibaca manusia (human-readable)
			)
			return nil
		},
	}))
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())
	e.Use(middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(20)))

	// 6. Routing
	// protected := e.Group("/api", appMiddleware.JWTMiddleware(jwtSecret))
	// user.RegisterModule(e, protected, db, rdb, jwtSecret)
	order.RegisterModule(e, db)

	fmt.Printf("🚀 Server Echo REST API berjalan di http://localhost:%s\n", port)
	e.Logger.Fatal(e.Start(":" + port))
}
