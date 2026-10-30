package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"shop_echo/config"
	"shop_echo/internal/domain"
	"time"

	"github.com/joho/godotenv"
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

	// 2. Koneksi DB & Automigrate
	db := config.ConnectDB()
	listTable := []any{&domain.User{}, &domain.Order{}}
	if err := db.AutoMigrate(listTable...); err != nil {
		log.Fatal("Gagal melakukan AutoMigrate:", err)
	}

	users := []domain.User{
		{
			Nama:  "Budi Santoso",
			Email: "budi@example.com",
		},
		{
			Nama:  "Siti Aminah",
			Email: "siti@example.com",
		},
		{
			Nama:  "Eko Kurniawan",
			Email: "eko@example.com",
		},
		{
			Nama:  "Rian Hidayat",
			Email: "rian@example.com",
		},
		{
			Nama:  "Dewi Lestari",
			Email: "dewi@example.com",
		},
	}

	for i := range users {
		db.Where(domain.User{Email: users[i].Email}).FirstOrCreate(&users[i])
	}

	statuses := []string{"PENDING", "PAID", "SHIPPED", "CANCELLED"}
	amounts := []float64{50000, 125000, 250000, 500000, 750000, 1200000}
	for i := 1; i <= 1500; i++ {
		randomUser := users[rand.Intn(len(users))]
		order := domain.Order{
			UserID:      randomUser.ID,
			TotalAmount: amounts[rand.Intn(len(amounts))],
			Status:      statuses[rand.Intn(len(statuses))],
			CreatedAt:   time.Now().Add(-time.Duration(rand.Intn(72)) * time.Hour),
		}
		db.Where(domain.Order{ID: uint(i)}).FirstOrCreate(&order)
	}
	fmt.Println("🎉 Seeding Selesai!")
}
