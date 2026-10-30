package permission

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"shop_echo/internal/domain"
	"shop_echo/internal/middleware"
	"shop_echo/pkg/response"
	"strconv"

	"github.com/labstack/echo/v4"
)

type PermissionHandler struct {
	repo domain.PermissionRepository
}

func NewPermissionHandler(repo domain.PermissionRepository) *PermissionHandler {
	return &PermissionHandler{repo: repo}
}

func (h *PermissionHandler) CreatePermission(c echo.Context) error {
	var req domain.Permission

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
	if err := h.repo.Create(ctx, &req); err != nil {
		return response.SendError(c, http.StatusInternalServerError, "Gagal menyimpan permission", err.Error())
	}
	return response.SendSuccess[*domain.Permission](c, http.StatusCreated, "Permission berhasil ditmbahkan", nil)
}

func (h *PermissionHandler) GetAllPermissions(c echo.Context) error {
	var req domain.PermissionQuery
	if err := c.Bind(&req); err != nil {
		if _, ok := err.(*echo.HTTPError); ok {
			return response.SendError(c, http.StatusBadRequest, "Format JSON tidak valid", []string{"Tipe data tidak valid"})
		}

		return response.SendError(c, http.StatusBadRequest, "Format JSON tidak valid", []string{err.Error()})
	}
	if err := c.Validate(&req); err != nil {
		return response.SendError(c, http.StatusBadRequest, "Validasi gagal", middleware.FormatValidationError(err))
	}
	// 1. Ambil query paramter 'page' (default: 1)
	if req.Page <= 0 {
		req.Page = 1
	}

	if req.Limit <= 0 {
		req.Limit = 10
	}

	if req.OrderBy == "" {
		req.OrderBy = "id"
	}

	if req.DirOrder == "" {
		req.DirOrder = "ASC"
	}
	// 3. Panggil repository
	ctx := c.Request().Context()
	permissions, totalRows, err := h.repo.GetAll(ctx, &req)
	if err != nil {
		return response.SendError(c, http.StatusInternalServerError, "Gagal mengambil daftar permission", err.Error())
	}

	// 4. Hitung total halaman (total pages)
	totalPages := int(math.Ceil(float64(totalRows) / float64(req.Limit)))
	meta := response.PaginationMeta{
		Page:       req.Page,
		Limit:      req.Limit,
		TotalRows:  totalRows,
		TotalPages: totalPages,
	}

	// 5. Kirim Response Berpaginasi
	return response.SendPaginateSuccess(c, http.StatusOK, "Berhasil mengambil semua permission", permissions, meta)
}

func (h *PermissionHandler) GetPermissionByID(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return response.SendError(c, http.StatusBadRequest, "ID harus berupa angka integer", nil)
	}
	ctx := c.Request().Context()
	role, err := h.repo.GetByID(ctx, id)
	if err != nil {
		return response.SendError(c, http.StatusNotFound, fmt.Sprintf("Role ID %d tidak ditemukan", id), nil)
	}
	return response.SendSuccess(c, http.StatusOK, "Berhasil mengambil detail permission", role)
}
