package domain

import (
	"context"
	"time"
)

type Order struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	UserID      uint      `gorm:"not null" json:"user_id"`
	User        User      `gorm:"foreignKey:UserID;" json:"user"`
	TotalAmount float64   `gorm:"type:numeric(12,2)" json:"total_amount"`
	Status      string    `gorm:"type:varchar(20)" json:"status"`
	CreatedAt   time.Time `json:"created_at"`
}

type OrderRepository interface {
	GetOrderDetailByID(ctx context.Context, orderID uint) (OrderDetailReq, error)
}

type OrderDetailReq struct {
	OrderID      uint    `json:"order_id" validate:"numeric,required"`
	CustomerName string  `json:"customer_name" validate:"alphanum,required"`
	TotalAmount  float64 `json:"total_amount" validate:"numeric,required,gt=0"`
	Status       string  `json:"status" validate:"alphanum,required"`
}

type BulkOrderRequest struct {
	OrderIDs []uint `json:"order_ids" validate:"required,min=1,dive,numeric,gt=0"`
}

type OrderJob struct {
	OrderID uint
}

type OrderResult struct {
	Data OrderDetailReq
	Err  error
}
