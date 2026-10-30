package role

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

type RoleHandler struct {
	repo domain.RoleRepository
}

func NewRoleHandler(repo domain.RoleRepository) *RoleHandler {
	return &RoleHandler{repo: repo}
}

func (h *RoleHandler) CreateRole(c echo.Context) error {
	var req domain.Role

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
		return response.SendError(c, http.StatusInternalServerError, "Gagal menyimpan role", err.Error())
	}
	return response.SendSuccess[*domain.Role](c, http.StatusCreated, "Role berhasil ditmbahkan", nil)
}

func (h *RoleHandler) GetAllRoles(c echo.Context) error {
	// 1. Ambil query paramter 'page' (default: 1)
	page, _ := strconv.Atoi(c.QueryParam("page"))
	if page <= 0 {
		page = 1
	}

	// 2. Ambil query paramter 'limit' (default: 10)
	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	if limit <= 0 {
		limit = 10
	}

	// 3. Panggil repository
	ctx := c.Request().Context()
	roles, totalRows, err := h.repo.GetAll(ctx, page, limit)
	if err != nil {
		return response.SendError(c, http.StatusInternalServerError, "Gagal mengambil daftar role", err.Error())
	}

	// 4. Hitung total halaman (total pages)
	totalPages := int(math.Ceil(float64(totalRows) / float64(limit)))
	meta := response.PaginationMeta{
		Page:       page,
		Limit:      limit,
		TotalRows:  totalRows,
		TotalPages: totalPages,
	}

	// 5. Kirim Response Berpaginasi
	return response.SendPaginateSuccess(c, http.StatusOK, "Berhasil mengambil semua role", roles, meta)
}

func (h *RoleHandler) GetProductByID(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return response.SendError(c, http.StatusBadRequest, "ID harus berupa angka integer", nil)
	}
	ctx := c.Request().Context()
	role, err := h.repo.GetByID(ctx, id)
	if err != nil {
		return response.SendError(c, http.StatusNotFound, fmt.Sprintf("Role ID %d tidak ditemukan", id), nil)
	}
	return response.SendSuccess(c, http.StatusOK, "Berhasil mengambil detail role", role)
}
