package services

import (
	"cashier-api/models"
	"cashier-api/repositories"
)

type ReportService struct {
	repo repositories.ReportRepositoryInput
}

func NewReportService(repo repositories.ReportRepositoryInput) *ReportService {
	return &ReportService{repo: repo}
}

func (s *ReportService) GetSalesSummaryToday() (*models.SalesSummary, error) {
	return s.repo.GetSalesSummaryToday()
}

func (s *ReportService) GetSalesSummaryRange(startDate, endDate string) (*models.SalesSummary, error) {
	return s.repo.GetSalesSummaryRange(startDate, endDate)
}
