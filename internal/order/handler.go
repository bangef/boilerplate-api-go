package order

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"shop_echo/internal/domain"
	"shop_echo/internal/middleware"
	"shop_echo/pkg/response"
	"sync"
	"time"

	"github.com/labstack/echo/v4"
)

type OrderHandler struct {
	repo domain.OrderRepository
}

func NewOrderHandler(repo domain.OrderRepository) *OrderHandler {
	return &OrderHandler{repo: repo}
}

func (h *OrderHandler) ProcessBulkOrders(c echo.Context) error {
	var req domain.BulkOrderRequest

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

	ctx, cancle := context.WithTimeout(c.Request().Context(), 5*time.Second)
	defer cancle()

	totalJobs := len(req.OrderIDs)

	workerCount := 5

	if totalJobs < workerCount {
		workerCount = totalJobs
	}

	jobs := make(chan domain.OrderJob, totalJobs)
	results := make(chan domain.OrderResult, totalJobs)
	var wg sync.WaitGroup

	for w := 1; w <= workerCount; w++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			for job := range jobs {
				select {
				case <-ctx.Done():
					results <- domain.OrderResult{
						Err: fmt.Errorf("worker [%d]: timeout.canceled untuk Order ID [%d]", workerID, job.OrderID),
					}
					continue
				default:
				}

				detail, err := h.repo.GetOrderDetailByID(ctx, job.OrderID)
				if err != nil {
					results <- domain.OrderResult{
						Err: fmt.Errorf("order ID %d error: %v", job.OrderID, err),
					}
				}

				results <- domain.OrderResult{Data: detail}
			}
		}(w)
	}

	for _, id := range req.OrderIDs {
		jobs <- domain.OrderJob{OrderID: id}
	}
	close(jobs)

	go func() {
		wg.Wait()
		close(results)
	}()

	var processedOrders []domain.OrderDetailReq
	var errorList []string

	for res := range results {
		if res.Err != nil {
			errorList = append(errorList, res.Err.Error())
			continue
		}
		processedOrders = append(processedOrders, res.Data)
	}

	return response.SendSuccess[any](c, http.StatusOK, "Data Order berhasil didapatkan", processedOrders)
}
