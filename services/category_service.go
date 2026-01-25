package services

import (
	"cashier-api/models"
	"cashier-api/repositories"
)

type CategoryServiceInput interface {
	GetAll() ([]models.Category, error)
	Create(category *models.Category) error
	GetByID(id int) (*models.Category, error)
	Update(category *models.Category) error
	Delete(id int) error
}

type categoryService struct {
	repo repositories.CategoryRepositoryInput
}

func NewCategoryService(repo repositories.CategoryRepositoryInput) CategoryServiceInput {
	return &categoryService{repo: repo}
}

func (s *categoryService) GetAll() ([]models.Category, error) {
	return s.repo.GetAll()
}

func (s *categoryService) Create(category *models.Category) error {
	return s.repo.Create(category)
}

func (s *categoryService) GetByID(id int) (*models.Category, error) {
	return s.repo.GetByID(id)
}

func (s *categoryService) Update(category *models.Category) error {
	return s.repo.Update(category)
}

func (s *categoryService) Delete(id int) error {
	return s.repo.Delete(id)
}
