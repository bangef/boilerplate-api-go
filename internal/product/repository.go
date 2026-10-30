package product

import (
	"context"
	"shop_echo/internal/domain"

	"gorm.io/gorm"
)

type productRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) domain.ProductRepository {
	return &productRepository{
		db: db,
	}
}

func (r *productRepository) Create(ctx context.Context, p *domain.Product) error {
	return r.db.WithContext(ctx).Create(p).Error
}

func (r *productRepository) GetAll(ctx context.Context, page int, limit int) ([]domain.Product, int64, error) {
	var products []domain.Product
	var totalRows int64
	// 1. Hitung total seluruh baris data (tanpa limit & offset)
	if err := r.db.WithContext(ctx).Model(&domain.Product{}).Count(&totalRows).Error; err != nil {
		return nil, 0, err
	}

	// 2. Hitung nilai ofset
	offset := (page - 1) * limit

	// 3. Query data menggunakan Limti dan Offset
	// 3.1 Preloaded("User") akan otomatis mengambil relasi data User
	if err := r.db.WithContext(ctx).Preload("User").Limit(limit).Offset(offset).Order("id desc").Find(&products).Error; err != nil {
		return nil, 0, err
	}
	return products, totalRows, nil
}

func (r *productRepository) GetByID(ctx context.Context, id int) (*domain.Product, error) {
	var p domain.Product
	// 1. Preloaded("User") akan otomatis mengambil relasi data User
	if err := r.db.WithContext(ctx).Preload("User").First(&p, id).Error; err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *productRepository) Update(ctx context.Context, id int, req *domain.ProductRequest) (*domain.Product, error) {
	p, err := r.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	p.Nama = req.Nama
	p.Harga = req.Harga
	p.Stok = req.Stok

	if err := r.db.WithContext(ctx).Save(p).Error; err != nil {
		return nil, err
	}
	return p, nil
}

func (r *productRepository) Delete(ctx context.Context, id int) error {
	result := r.db.WithContext(ctx).Delete(&domain.Product{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
