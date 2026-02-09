package services

import (
	"cashier-api/models"
	"cashier-api/repositories"
)

type ProductServiceInput interface {
	GetAll(page, limit, name string) ([]models.Product, error)
	Create(product *models.Product) error
	GetByID(id int) (*models.Product, error)
	Update(product *models.Product) error
	Delete(id int) error
}

type productService struct {
	repo repositories.ProductRepositoryInput
}

func NewProductService(repo repositories.ProductRepositoryInput) ProductServiceInput {
	return &productService{repo: repo}
}

func (s *productService) GetAll(page, limit, name string) ([]models.Product, error) {
	return s.repo.GetAll(page, limit, name)
}

func (s *productService) Create(product *models.Product) error {
	return s.repo.Create(product)
}

func (s *productService) GetByID(id int) (*models.Product, error) {
	return s.repo.GetByID(id)
}

func (s *productService) Update(product *models.Product) error {
	return s.repo.Update(product)
}

func (s *productService) Delete(id int) error {
	return s.repo.Delete(id)
}
