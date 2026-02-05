package services

import (
	"belajar-go/models"
	"belajar-go/repositories"
	"time"
)

type ReportService struct {
	repo *repositories.ReportRepository
}

func NewReportService(repo *repositories.ReportRepository) *ReportService {
	return &ReportService{repo: repo}
}

func (s *ReportService) GetReport(startDate, endDate time.Time) (*models.ReportResponse, error) {
	return s.repo.GetReport(startDate, endDate)
}

func (s *ReportService) GetTodayReport() (*models.ReportResponse, error) {
	now := time.Now()
	startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	endOfDay := startOfDay.Add(24 * time.Hour)
	
	return s.repo.GetReport(startOfDay, endOfDay)
}
