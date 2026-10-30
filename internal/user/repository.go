package user

import (
	"context"
	"fmt"
	"shop_echo/internal/domain"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) domain.UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(ctx context.Context, u *domain.User) error {
	return r.db.WithContext(ctx).Create(u).Error
}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	var u domain.User
	if err := r.db.WithContext(ctx).Preload("Roles.Permissions").Where("email = ?", email).First(&u).Error; err != nil {
		return nil, err
	}
	return &u, nil
}

type tokenRepository struct {
	rdb *redis.Client
}

func NewTokenRepository(rdb *redis.Client) domain.TokenRepository {
	return &tokenRepository{rdb: rdb}
}

func (r *tokenRepository) SaveRefreshToken(ctx context.Context, userID uint, tokenStr string, duration time.Duration) error {
	// Key Format: refresh_token<token> -> value: <userID>
	key := fmt.Sprintf("refresh_token:%s", tokenStr)
	err := r.rdb.Set(ctx, key, userID, duration).Err()
	if err != nil {
		return err
	}
	// Simpan juga relasi user -> token agar bisa melakukan Logout All
	userKey := fmt.Sprintf("user_tokens:%d:%s", userID, tokenStr)
	return r.rdb.Set(ctx, userKey, "1", duration).Err()
}

func (r *tokenRepository) GetUserIDByToken(ctx context.Context, tokenStr string) (uint, error) {
	// Key Format: refresh_token<token> -> value: <userID>
	key := fmt.Sprintf("refresh_token:%s", tokenStr)
	val, err := r.rdb.Get(ctx, key).Result()
	if err == redis.Nil {
		return 0, fmt.Errorf("refresh token tidak ditemukan atau sudah kadaluwarsa")
	} else if err != nil {
		return 0, err
	}

	userID, _ := strconv.Atoi(val)
	return uint(userID), nil
}

func (r *tokenRepository) RevokeToken(ctx context.Context, tokenStr string) error {
	// Cari userID dulu dari token
	userID, err := r.GetUserIDByToken(ctx, tokenStr)
	if err == nil {
		userKey := fmt.Sprintf("user_tokens:%d:%s", userID, tokenStr)
		r.rdb.Del(ctx, userKey)
	}

	key := fmt.Sprintf("refresh_token:%s", tokenStr)
	return r.rdb.Del(ctx, key).Err()
}

func (r *tokenRepository) RevokeAllUserTokens(ctx context.Context, userID uint) error {
	pattern := fmt.Sprintf("user_tokens:%d:*", userID)
	iter := r.rdb.Scan(ctx, 0, pattern, 0).Iterator()
	for iter.Next(ctx) {
		userKey := iter.Val()
		var tokenStr string
		fmt.Sscanf(userKey, fmt.Sprintf("user_tokens:%d:%%s", userID), &tokenStr)
		// Hapus key token asli
		tokenKey := fmt.Sprintf("refresh_token:%s", tokenStr)
		r.rdb.Del(ctx, tokenKey)
		r.rdb.Del(ctx, userKey)
	}
	return iter.Err()
}
