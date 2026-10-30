package config

import (
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectDB() *gorm.DB {
	dbHost := os.Getenv("DB_HOST")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dbPort := os.Getenv("DB_PORT")
	dbSSLMode := os.Getenv("DB_SSLMODE")

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=Asia/Jakarta", dbHost, dbUser, dbPassword, dbName, dbPort, dbSSLMode)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Gagal koneksi ke database PostreSQL:", err)
	}

	pgsqlDB, err := db.DB()
	if err != nil {
		log.Fatal("Gagal koneksi ke database PostgreSQL:", err)
	}
	pgsqlDB.SetMaxOpenConns(25)                  // Batasi max koneksi agar tidak menembus batas Postgres (100)
	pgsqlDB.SetMaxIdleConns(25)                  // Samakan dengan MaxOpenConns agar koneksi selalu stand-by
	pgsqlDB.SetConnMaxLifetime(15 * time.Minute) // Ganti koneksi berkala untuk mencegah "broken pipe"
	pgsqlDB.SetConnMaxIdleTime(5 * time.Minute)  // Bersihkan koneksi menganggur jika aplikasi sedang sepi
	return db
}
