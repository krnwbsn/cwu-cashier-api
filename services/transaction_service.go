package services

import (
	"cashier-api/models"
	"cashier-api/repositories"
)

type TransactionService struct {
	repo repositories.TransactionRepositoryInput
}

func NewTransactionService(repo repositories.TransactionRepositoryInput) *TransactionService {
	return &TransactionService{repo: repo}
}

func (s *TransactionService) Checkout(items []models.CheckoutItem, useLock bool) (*models.Transaction, error) {
	return s.repo.CreateTransaction(items)
}
