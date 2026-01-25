package database

import (
	"database/sql"
	"fmt"
	"log"
)

func Seed(db *sql.DB) error {
	log.Println("Running database seeding...")

	categories := []struct {
		Name        string
		Description string
	}{
		{"Electronics", "Gadgets and devices"},
		{"Clothing", "Apparel and accessories"},
		{"Food & Beverage", "Consumables"},
	}

	for _, c := range categories {
		var exists bool
		checkQuery := "SELECT EXISTS(SELECT 1 FROM category WHERE name = $1)"
		err := db.QueryRow(checkQuery, c.Name).Scan(&exists)
		if err != nil {
			return err
		}

		if !exists {
			_, err := db.Exec("INSERT INTO category (name, description) VALUES ($1, $2)", c.Name, c.Description)
			if err != nil {
				return fmt.Errorf("failed to seed category %s: %w", c.Name, err)
			}
			log.Printf("Seeded category: %s", c.Name)
		} else {
			log.Printf("Category %s already exists, skipping", c.Name)
		}
	}

	var electronicsID int
	err := db.QueryRow("SELECT id FROM category WHERE name = 'Electronics' LIMIT 1").Scan(&electronicsID)
	if err == nil {
		products := []struct {
			Name       string
			Price      float64
			Stock      int
			CategoryID int
		}{
			{"Smartphone", 699.99, 50, electronicsID},
			{"Laptop", 1299.99, 20, electronicsID},
		}

		for _, p := range products {
			var exists bool
			checkQuery := "SELECT EXISTS(SELECT 1 FROM product WHERE name = $1)"
			err := db.QueryRow(checkQuery, p.Name).Scan(&exists)
			if err != nil {
				return err
			}

			if !exists {
				_, err := db.Exec("INSERT INTO product (name, price, stock, category_id) VALUES ($1, $2, $3, $4)",
					p.Name, p.Price, p.Stock, p.CategoryID)
				if err != nil {
					return fmt.Errorf("failed to seed product %s: %w", p.Name, err)
				}
				log.Printf("Seeded product: %s", p.Name)
			} else {
				log.Printf("Product %s already exists, skipping", p.Name)
			}
		}
	}

	log.Println("Database seeding completed")
	return nil
}
