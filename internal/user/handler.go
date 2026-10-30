package user

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"shop_echo/internal/domain"
	"shop_echo/internal/middleware"
	"shop_echo/pkg/response"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	userRepo  domain.UserRepository
	tokenRepo domain.TokenRepository
	jwtSecret []byte
}

func NewAuthHandler(userRepo domain.UserRepository, tokenRepo domain.TokenRepository, jwtSecret string) *AuthHandler {
	return &AuthHandler{
		userRepo:  userRepo,
		tokenRepo: tokenRepo,
		jwtSecret: []byte(jwtSecret),
	}
}

// Helper
func (h *AuthHandler) generateTokens(ctx context.Context, user *domain.User) (string, string, error) {
	//1. generate access token (15 menit)
	var roleNames []string
	var permissionNames []string

	permMap := make(map[string]bool)
	for _, role := range user.Roles {
		roleNames = append(roleNames, role.Nama)
		for _, perm := range role.Permission {
			if !permMap[perm.Nama] {
				permMap[perm.Nama] = true
				permissionNames = append(permissionNames, perm.Nama)
			}
		}
	}

	accessClaims := &domain.JWTCustomClaims{
		UserID:     user.ID,
		Email:      user.Email,
		TokenType:  "access",
		Role:       roleNames,
		Permission: permissionNames,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	accessTokenObj := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessToken, err := accessTokenObj.SignedString(h.jwtSecret)
	if err != nil {
		return "", "", err
	}

	// 2. generate refresh token (7 hari)
	refreshDuration := 7 * 24 * time.Hour
	refreshExpiresAt := time.Now().Add(refreshDuration)
	refreshClaims := &domain.JWTCustomClaims{
		UserID:     user.ID,
		Email:      user.Email,
		TokenType:  "refresh",
		Role:       roleNames,
		Permission: permissionNames,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(refreshExpiresAt),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	refreshTokenObj := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshTokenStr, err := refreshTokenObj.SignedString(h.jwtSecret)
	if err != nil {
		return "", "", err
	}

	// 3. Simpan Refresh Token ke Database
	if err := h.tokenRepo.SaveRefreshToken(ctx, user.ID, refreshTokenStr, refreshDuration); err != nil {
		return "", "", err
	}

	return accessToken, refreshTokenStr, nil
}

// POST /api/register
func (h *AuthHandler) Register(c echo.Context) error {
	var req domain.RegisterRequest

	if err := c.Bind(&req); err != nil {
		if ute, ok := err.(*json.UnmarshalTypeError); ok {
			detailErr := fmt.Sprintf("Field '%s' harus bertipe '%s', tetapi dikirim '%s'", ute.Field, ute.Type, ute.Value)
			return response.SendError(c, http.StatusBadRequest, "Tipe data JSON tidak sesuai", []string{detailErr})
		}

		if _, ok := err.(*echo.HTTPError); ok {
			return response.SendError(c, http.StatusBadRequest, "Format JSON tidak valid", []string{"Tipe data tidak valid"})
		}

		return response.SendError(c, http.StatusBadRequest, "Format JSON tidak valid", []string{err.Error()})
	}

	if err := c.Validate(&req); err != nil {
		return response.SendError(c, http.StatusBadRequest, "Validasi gagal", middleware.FormatValidationError(err))
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return response.SendError(c, http.StatusInternalServerError, "Gagal enkripsi password", err.Error())
	}
	user := domain.User{
		Nama:     req.Nama,
		Email:    req.Email,
		Password: string(hashedPassword),
	}
	ctx := c.Request().Context()
	if err := h.userRepo.Create(ctx, &user); err != nil {
		return response.SendError(c, http.StatusInternalServerError, "Gagal membuat user (Email mungkin sudah terdaftar)", err.Error())
	}
	return response.SendSuccess(c, http.StatusCreated, "Registrasi berhasil", user)
}

// POST /api/login
func (h *AuthHandler) Login(c echo.Context) error {
	var req domain.LoginRequest

	if err := c.Bind(&req); err != nil {
		if ute, ok := err.(*json.UnmarshalTypeError); ok {
			detailErr := fmt.Sprintf("Field '%s' harus bertipe '%s', tetapi dikirim '%s'", ute.Field, ute.Type, ute.Value)
			return response.SendError(c, http.StatusBadRequest, "Tipe data JSON tidak sesuai", []string{detailErr})
		}

		if _, ok := err.(*echo.HTTPError); ok {
			return response.SendError(c, http.StatusBadRequest, "Format JSON tidak valid", []string{"Tipe data tidak valid"})
		}

		return response.SendError(c, http.StatusBadRequest, "Format JSON tidak valid", []string{err.Error()})
	}

	if err := c.Validate(&req); err != nil {
		return response.SendError(c, http.StatusBadRequest, "Validasi gagal", middleware.FormatValidationError(err))
	}
	ctx := c.Request().Context()
	user, err := h.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		return response.SendError(c, http.StatusUnauthorized, "Email / password salah", nil)
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return response.SendError(c, http.StatusUnauthorized, "Email atau password salah", nil)
	}

	accessToken, refreshToken, err := h.generateTokens(ctx, user)
	if err != nil {
		return response.SendError(c, http.StatusInternalServerError, "Gagal simpan token", err.Error())
	}

	return response.SendSuccess(c, http.StatusOK, "Login berhasil", map[string]interface{}{
		"tokens": domain.TokenResponse{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
		},
		"user": user,
	})
}

// POST /api/refresh (Dengan Pengecekan Revocation di Dataabase)
func (h *AuthHandler) RefreshToken(c echo.Context) error {
	var req domain.RefreshTokenRequest

	if err := c.Bind(&req); err != nil {
		if ute, ok := err.(*json.UnmarshalTypeError); ok {
			detailErr := fmt.Sprintf("Field '%s' harus bertipe '%s', tetapi dikirim '%s'", ute.Field, ute.Type, ute.Value)
			return response.SendError(c, http.StatusBadRequest, "Tipe data JSON tidak sesuai", []string{detailErr})
		}

		if _, ok := err.(*echo.HTTPError); ok {
			return response.SendError(c, http.StatusBadRequest, "Format JSON tidak valid", []string{"Tipe data tidak valid"})
		}

		return response.SendError(c, http.StatusBadRequest, "Format JSON tidak valid", []string{err.Error()})
	}

	if err := c.Validate(&req); err != nil {
		return response.SendError(c, http.StatusBadRequest, "Validasi gagal", middleware.FormatValidationError(err))
	}

	ctx := c.Request().Context()
	// 1. Cek Keberadaan dan Status token di Database
	_, err := h.tokenRepo.GetUserIDByToken(ctx, req.RefreshToken)
	if err != nil {
		return response.SendError(c, http.StatusUnauthorized, "Refresh token telah dicabut (revoked) atau tidak ditemukan", nil)
	}

	// 2. Parse & Verifikasi JWT Signature
	token, err := jwt.ParseWithClaims(req.RefreshToken, &domain.JWTCustomClaims{}, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("metode signing tidak vlid")
		}
		return h.jwtSecret, nil
	})

	if err != nil || !token.Valid {
		return response.SendError(c, http.StatusUnauthorized, "Refresh token tidak valid atau expired", nil)
	}

	claims, ok := token.Claims.(*domain.JWTCustomClaims)
	if !ok || claims.TokenType != "refresh" {
		return response.SendError(c, http.StatusUnauthorized, "Token bukan Refresh Token valid", nil)
	}

	user, err := h.userRepo.GetByEmail(ctx, claims.Email)
	if err != nil {
		return response.SendError(c, http.StatusUnauthorized, "User tidak ditemukan", nil)
	}

	//3. Revoke Token Lama (Token Rotation)
	_ = h.tokenRepo.RevokeToken(ctx, req.RefreshToken)

	//4. Generat Pasangan Token Baru
	newAccessToken, newRefreshToken, err := h.generateTokens(ctx, user)
	if err != nil {
		return response.SendError(c, http.StatusInternalServerError, "Gagal generate token baru", err.Error())
	}

	return response.SendSuccess(c, http.StatusOK, "Refresh token berhasil", domain.TokenResponse{
		AccessToken:  newAccessToken,
		RefreshToken: newRefreshToken,
	})
}

// POST /api/logout
func (h *AuthHandler) Logout(c echo.Context) error {
	var req domain.LogoutRequest
	if err := c.Bind(&req); err != nil {
		if ute, ok := err.(*json.UnmarshalTypeError); ok {
			detailErr := fmt.Sprintf("Field '%s' harus bertipe '%s', tetapi dikirim '%s'", ute.Field, ute.Type, ute.Value)
			return response.SendError(c, http.StatusBadRequest, "Tipe data JSON tidak sesuai", []string{detailErr})
		}

		if _, ok := err.(*echo.HTTPError); ok {
			return response.SendError(c, http.StatusBadRequest, "Format JSON tidak valid", []string{"Tipe data tidak valid"})
		}

		return response.SendError(c, http.StatusBadRequest, "Format JSON tidak valid", []string{err.Error()})
	}

	if err := c.Validate(&req); err != nil {
		return response.SendError(c, http.StatusBadRequest, "Validasi gagal", middleware.FormatValidationError(err))
	}

	// 1. Pencabutan (Revoke) Token di Database
	ctx := c.Request().Context()
	if err := h.tokenRepo.RevokeToken(ctx, req.RefreshToken); err != nil {
		return response.SendError(c, http.StatusBadRequest, "Gagal melakukan logout / token tidak ditemukan", nil)
	}

	return response.SendSuccess[*domain.User](c, http.StatusOK, "Logout berhasil. Refresh token telah dicabut", nil)
}

// POST /api/logout-all (Logout dari Semua Perangkat)
func (h *AuthHandler) LogoutAll(c echo.Context) error {
	userID := c.Get("user_id").(uint)
	ctx := c.Request().Context()
	if err := h.tokenRepo.RevokeAllUserTokens(ctx, userID); err != nil {
		return response.SendError(c, http.StatusInternalServerError, "Gagal mencabut seluruh sesi login", err.Error())
	}
	return response.SendSuccess[*domain.User](c, http.StatusOK, "Berhasil logout dari seluruh perangkat", nil)
}
