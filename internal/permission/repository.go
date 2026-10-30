package permission

import (
	"context"
	"shop_echo/internal/domain"

	"gorm.io/gorm"
)

type permissionRepository struct {
	db *gorm.DB
}

func NewPermissionRepository(db *gorm.DB) domain.PermissionRepository {
	return &permissionRepository{
		db: db,
	}
}

func (r *permissionRepository) Create(ctx context.Context, p *domain.Permission) error {
	return r.db.WithContext(ctx).Create(p).Error
}

func (r *permissionRepository) GetAll(ctx context.Context, req *domain.PermissionQuery) ([]domain.Permission, int64, error) {
	var permissions []domain.Permission
	var totalRows int64
	// 1. Hitung total seluruh baris data (tanpa limit & offset)
	req.Search = "%" + req.Search + "%"
	query := r.db.WithContext(ctx).Model(&domain.Permission{}).Count(&totalRows)
	if req.Search != "%%" {
		query = query.Where("nama ILIKE ?", req.Search)
	}
	if err := query.Error; err != nil {
		return nil, 0, err
	}

	// 2. Hitung nilai ofset
	offset := (req.Page - 1) * req.Limit

	// 3. Query data menggunakan Limti dan Offset
	// 3.1 Preloaded("User") akan otomatis mengambil relasi data User
	if err := r.db.WithContext(ctx).Limit(req.Limit).Offset(offset).Order(req.OrderBy+" "+req.DirOrder).Where("nama ILIKE ?", req.Search).Find(&permissions).Error; err != nil {
		return nil, 0, err
	}
	return permissions, totalRows, nil
}

func (r *permissionRepository) GetByID(ctx context.Context, id int) (*domain.Permission, error) {
	var permission domain.Permission
	if err := r.db.WithContext(ctx).First(&permission, id).Error; err != nil {
		return nil, err
	}
	return &permission, nil
}

func (r *permissionRepository) Update(ctx context.Context, id int, req *domain.PermissionRequest) (*domain.Permission, error) {
	p, err := r.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	p.Nama = req.Nama
	p.Desskripsi = req.Deskripsi

	if err := r.db.WithContext(ctx).Save(p).Error; err != nil {
		return nil, err
	}
	return p, nil
}

func (r *permissionRepository) Delete(ctx context.Context, id int) error {
	result := r.db.WithContext(ctx).Delete(&domain.Permission{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
