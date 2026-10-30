package domain

import (
	"context"
	"time"

	"gorm.io/gorm"
)

type Product struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	UserID    uint           `gorm:"not null" json:"user_id"`
	User      *User          `gorm:"foreignKey:UserID" json:"user,omitempty"` // Relasi Belongs-To
	Nama      string         `gorm:"type:varchar(100);not null" json:"nama"`
	Harga     float64        `gorm:"type:numeric(12,2);not null" json:"harga"`
	Stok      int            `gorm:"not null" json:"stok"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

type ProductRequest struct {
	Nama  string  `json:"nama" validate:"alphanum,required,min=3"`
	Harga float64 `json:"harga" validate:"number,required,gt=0"`
	Stok  int     `json:"stok" validate:"number,required,gte=0"`
}

type ProductRepository interface {
	Create(ctx context.Context, p *Product) error
	GetAll(ctx context.Context, page, limit int) ([]Product, int64, error)
	GetByID(ctx context.Context, id int) (*Product, error)
	Update(ctx context.Context, id int, req *ProductRequest) (*Product, error)
	Delete(ctx context.Context, id int) error
}
