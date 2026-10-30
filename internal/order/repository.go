package order

import (
	"context"
	"shop_echo/internal/domain"

	"gorm.io/gorm"
)

type orderRepository struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) domain.OrderRepository {
	return &orderRepository{
		db: db,
	}
}

func (r *orderRepository) GetOrderDetailByID(ctx context.Context, orderID uint) (domain.OrderDetailReq, error) {
	var detail domain.OrderDetailReq

	tx := r.db.WithContext(ctx).Table("orders").Select("orders.id AS order_id, users.nama AS customer_name, orders.total_amount, orders.status").Joins("INNER JOIN users ON users.id = orders.user_id").Where("orders.id = ?", orderID).Scan(&detail)

	if tx.Error != nil {
		return detail, tx.Error
	}

	if tx.RowsAffected == 0 {
		return detail, gorm.ErrRecordNotFound
	}

	return detail, nil
}
