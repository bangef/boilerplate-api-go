package response

import "github.com/labstack/echo/v4"

type APIResponse[T any] struct {
	Code    int    `json:"code"`
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    T      `json:"data,omitempty"`
	Error   any    `json:"error,omitempty"`
}

type PaginationMeta struct {
	Page       int   `json:"page"`
	Limit      int   `json:"limit"`
	TotalRows  int64 `json:"total_rows"`
	TotalPages int   `json:"total_pages"`
}

type PaginatedResponse[T any] struct {
	Code    int            `json:"code"`
	Status  string         `json:"status"`
	Message string         `json:"message"`
	Data    []T            `json:"data"`
	Meta    PaginationMeta `json:"meta"`
}

func SendSuccess[T any](c echo.Context, statusCode int, message string, data T) error {
	return c.JSON(statusCode, APIResponse[T]{
		Code:    statusCode,
		Status:  "SUCCEES",
		Message: message,
		Data:    data,
	})
}

func SendPaginateSuccess[T any](c echo.Context, statusCode int, message string, data []T, meta PaginationMeta) error {
	return c.JSON(statusCode, PaginatedResponse[T]{
		Code:    statusCode,
		Status:  "SUCCEES",
		Message: message,
		Data:    data,
		Meta:    meta,
	})
}

func SendError(c echo.Context, statusCode int, message string, errDetail any) error {
	return c.JSON(statusCode, APIResponse[any]{
		Code:    statusCode,
		Status:  "ERROR",
		Message: message,
		Error:   errDetail,
	})
}
