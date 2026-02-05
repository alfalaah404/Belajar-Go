package handlers

import (
	"belajar-go/models"
	"belajar-go/services"
	"encoding/json"
	"net/http"
	"strings"
	"time"
)

type ReportHandler struct {
	service *services.ReportService
}

func NewReportHandler(service *services.ReportService) *ReportHandler {
	return &ReportHandler{service: service}
}

// HandleReports - GET /api/reports, GET /api/reports/hari-ini
func (h *ReportHandler) HandleReports(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Cek apakah path adalah /api/reports/hari-ini
	if strings.HasSuffix(r.URL.Path, "/hari-ini") {
		h.GetTodayReport(w, r)
		return
	}

	// Jika path adalah /api/reports dengan query parameters
	h.GetReportByDateRange(w, r)
}

// GetTodayReport godoc
// @Summary      Get today's report
// @Description  Get report for today's transactions
// @Tags         reports
// @Accept       json
// @Produce      json
// @Success      200  {object}  models.ApiResponse{data=models.ReportResponse}
// @Failure      500  {object}  models.ApiResponse
// @Router       /api/reports/hari-ini [get]
func (h *ReportHandler) GetTodayReport(w http.ResponseWriter, r *http.Request) {
	report, err := h.service.GetTodayReport()
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(models.ApiResponse{
			Status:  "error",
			Message: err.Error(),
			Data:    nil,
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(models.ApiResponse{
		Status:  "success",
		Message: "Success get today's report",
		Data:    report,
	})
}

// GetReportByDateRange godoc
// @Summary      Get report by date range
// @Description  Get report for transactions within a date range
// @Tags         reports
// @Accept       json
// @Produce      json
// @Param        start_date  query     string  true  "Start date (YYYY-MM-DD)"
// @Param        end_date    query     string  true  "End date (YYYY-MM-DD)"
// @Success      200  {object}  models.ApiResponse{data=models.ReportResponse}
// @Failure      400  {object}  models.ApiResponse
// @Failure      500  {object}  models.ApiResponse
// @Router       /api/reports [get]
func (h *ReportHandler) GetReportByDateRange(w http.ResponseWriter, r *http.Request) {
	startDateStr := r.URL.Query().Get("start_date")
	endDateStr := r.URL.Query().Get("end_date")

	if startDateStr == "" || endDateStr == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.ApiResponse{
			Status:  "error",
			Message: "start_date and end_date query parameters are required",
			Data:    nil,
		})
		return
	}

	// Parse dates
	startDate, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.ApiResponse{
			Status:  "error",
			Message: "Invalid start_date format. Use YYYY-MM-DD",
			Data:    nil,
		})
		return
	}

	endDate, err := time.Parse("2006-01-02", endDateStr)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.ApiResponse{
			Status:  "error",
			Message: "Invalid end_date format. Use YYYY-MM-DD",
			Data:    nil,
		})
		return
	}

	// Add 1 day to end_date to include the entire end date
	endDate = endDate.Add(24 * time.Hour)

	report, err := h.service.GetReport(startDate, endDate)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(models.ApiResponse{
			Status:  "error",
			Message: err.Error(),
			Data:    nil,
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(models.ApiResponse{
		Status:  "success",
		Message: "Success get report",
		Data:    report,
	})
}
