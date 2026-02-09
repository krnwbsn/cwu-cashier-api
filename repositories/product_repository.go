package repositories

import (
	"cashier-api/models"
	"database/sql"
	"errors"
	"fmt"
)

type ProductRepositoryInput interface {
	GetAll(page, limit, name string) ([]models.Product, error)
	Create(product *models.Product) error
	GetByID(id int) (*models.Product, error)
	Update(product *models.Product) error
	Delete(id int) error
}

type productRepository struct {
	db *sql.DB
}

func NewProductRepository(db *sql.DB) ProductRepositoryInput {
	return &productRepository{db: db}
}

func (repo *productRepository) GetAll(page, limit, name string) ([]models.Product, error) {
	fmt.Println(page, limit, name)
	query := `
		SELECT p.id, p.name, p.price, p.stock, COALESCE(p.category_id, 0), COALESCE(c.name, '') 
		FROM product p 
		LEFT JOIN category c ON p.category_id = c.id
	`
	args := []interface{}{}
	argNum := 1
	if name != "" {
		query += " WHERE p.name ILIKE $" + fmt.Sprint(argNum)
		args = append(args, "%"+name+"%")
		argNum++
	}
	query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", argNum, argNum+1)
	args = append(args, limit, page)

	rows, err := repo.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	products := make([]models.Product, 0)
	for rows.Next() {
		var p models.Product
		err := rows.Scan(&p.ID, &p.Name, &p.Price, &p.Stock, &p.CategoryID, &p.CategoryName)
		if err != nil {
			return nil, err
		}
		products = append(products, p)
	}

	return products, nil
}

func (repo *productRepository) Create(product *models.Product) error {
	query := "INSERT INTO product (name, price, stock, category_id) VALUES ($1, $2, $3, $4) RETURNING id"
	err := repo.db.QueryRow(query, product.Name, product.Price, product.Stock, product.CategoryID).Scan(&product.ID)
	if err != nil {
		return err
	}
	return nil
}

func (repo *productRepository) GetByID(id int) (*models.Product, error) {
	query := `
		SELECT p.id, p.name, p.price, p.stock, COALESCE(p.category_id, 0), COALESCE(c.name, '') 
		FROM product p 
		LEFT JOIN category c ON p.category_id = c.id 
		WHERE p.id = $1
	`

	var p models.Product
	err := repo.db.QueryRow(query, id).Scan(&p.ID, &p.Name, &p.Price, &p.Stock, &p.CategoryID, &p.CategoryName)
	if err == sql.ErrNoRows {
		return nil, errors.New("product not found")
	}
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (repo *productRepository) Update(product *models.Product) error {
	query := "UPDATE product SET name = $1, price = $2, stock = $3, category_id = $4 WHERE id = $5"
	result, err := repo.db.Exec(query, product.Name, product.Price, product.Stock, product.CategoryID, product.ID)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return errors.New("product not found")
	}

	return nil
}

func (repo *productRepository) Delete(id int) error {
	query := "DELETE FROM product WHERE id = $1"
	result, err := repo.db.Exec(query, id)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return errors.New("product not found")
	}

	return nil
}
