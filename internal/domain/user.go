package domain

import (
	"context"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

type User struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Nama      string         `gorm:"type:varchar(100);not null" json:"nama"`
	Email     string         `gorm:"type:varchar(100);not null;unique" json:"email"`
	Password  string         `gorm:"not null" json:"-"` // agar hash password tidak muncul di response
	Product   []Product      `gorm:"foreignKey:UserID" json:"products,omitempty"`
	Roles     []Role         `gorm:"many2many:user_roles;" json:"roles,omitempty"` // Relasi M2M to Role
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

type RegisterRequest struct {
	Nama     string `json:"nama" validate:"required,min=3"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"email,required"`
	Password string `json:"password" validate:"alphanum,required"`
}

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

type LogoutRequest struct {
	RefreshToken string `json:"refresh_token" validatPe:"required"`
}

type RefreshToken struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `gorm:"not null;index" json:"user_id"`
	User      *User     `gorm:"foreignKey:UserId;" json:"-"`
	Token     string    `gorm:"type:text;not null;uniqueIndex" json:"token"`
	IsRevoke  bool      `gorm:"default:false" json:"is_revoked"`
	ExpiresAt time.Time `gorm:"not null" json:"expires_at"`
	CreatedAt time.Time `json:"created_at"`
}

type JWTCustomClaims struct {
	UserID     uint     `json:"user_id"`
	Email      string   `json:"email"`
	Role       []string `json:"roles"`
	Permission []string `json:"permissions"`
	TokenType  string   `json:"token_type"` // access atau refresh
	jwt.RegisteredClaims
}

type UserRepository interface {
	Create(ctx context.Context, user *User) error
	GetByEmail(ctx context.Context, email string) (*User, error)
}

type TokenRepository interface {
	SaveRefreshToken(ctx context.Context, userID uint, tokenStr string, duration time.Duration) error
	GetUserIDByToken(ctx context.Context, tokenStr string) (uint, error)
	RevokeToken(ctx context.Context, tokenStr string) error
	RevokeAllUserTokens(ctx context.Context, userID uint) error
}
