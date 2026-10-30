package role

import (
	"context"
	"shop_echo/internal/domain"

	"gorm.io/gorm"
)

type roleRepository struct {
	db *gorm.DB
}

func NewRoleRepository(db *gorm.DB) domain.RoleRepository {
	return &roleRepository{
		db: db,
	}
}

func (r *roleRepository) Create(ctx context.Context, u *domain.Role) error {
	return r.db.WithContext(ctx).Create(u).Error
}

func (r *roleRepository) GetAll(ctx context.Context, page int, limit int) ([]domain.Role, int64, error) {
	var roles []domain.Role
	var totalRows int64
	// 1. Hitung total seluruh baris data (tanpa limit & offset)
	if err := r.db.WithContext(ctx).Model(&domain.Role{}).Count(&totalRows).Error; err != nil {
		return nil, 0, err
	}

	// 2. Hitung nilai ofset
	offset := (page - 1) * limit

	// 3. Query data menggunakan Limti dan Offset
	// 3.1 Preloaded("User") akan otomatis mengambil relasi data User
	if err := r.db.WithContext(ctx).Limit(limit).Offset(offset).Order("id desc").Find(&roles).Error; err != nil {
		return nil, 0, err
	}
	return roles, totalRows, nil
}

func (r *roleRepository) GetByID(ctx context.Context, id int) (*domain.Role, error) {
	var role domain.Role
	if err := r.db.WithContext(ctx).First(&role, id).Error; err != nil {
		return nil, err
	}
	return &role, nil
}

func (r *roleRepository) Update(ctx context.Context, id int, req *domain.RoleRequest) (*domain.Role, error) {
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

func (r *roleRepository) Delete(ctx context.Context, id int) error {
	result := r.db.WithContext(ctx).Delete(&domain.Role{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
