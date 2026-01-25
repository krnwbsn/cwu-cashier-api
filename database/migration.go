package database

import (
	"database/sql"
	"fmt"
	"log"
)

func Migrate(db *sql.DB) error {
	log.Println("Running database migrations...")

	createCategoryTable := `
	CREATE TABLE IF NOT EXISTS category (
		id SERIAL PRIMARY KEY,
		name VARCHAR(100) NOT NULL,
		description TEXT
	);`
	if _, err := db.Exec(createCategoryTable); err != nil {
		return fmt.Errorf("failed to create category table: %w", err)
	}

	var columnExists bool
	checkColumnQuery := `
	SELECT EXISTS (
		SELECT 1 
		FROM information_schema.columns 
		WHERE table_name='product' AND column_name='category_id'
	);`
	if err := db.QueryRow(checkColumnQuery).Scan(&columnExists); err != nil {
		return fmt.Errorf("failed to check category_id column: %w", err)
	}

	if !columnExists {
		log.Println("Adding category_id column to product table...")
		addCategoryColumn := `ALTER TABLE product ADD COLUMN category_id INT;`
		if _, err := db.Exec(addCategoryColumn); err != nil {
			return fmt.Errorf("failed to add category_id column: %w", err)
		}

		addFKConstraint := `
		ALTER TABLE product 
		ADD CONSTRAINT fk_category 
		FOREIGN KEY (category_id) 
		REFERENCES category(id) 
		ON DELETE SET NULL;`
		if _, err := db.Exec(addFKConstraint); err != nil {
			return fmt.Errorf("failed to add foreign key constraint: %w", err)
		}
		log.Println("category_id column added successfully")
	}

	var priceType string
	checkPriceTypeQuery := `
	SELECT data_type 
	FROM information_schema.columns 
	WHERE table_name='product' AND column_name='price' AND table_schema='public';`

	if err := db.QueryRow(checkPriceTypeQuery).Scan(&priceType); err != nil {
		log.Printf("Warning: failed to check price column type: %v", err)
	} else {
		if priceType == "integer" || priceType == "bigint" || priceType == "smallint" {
			log.Println("Migrating price column from integer to numeric...")
			alterPriceQuery := `ALTER TABLE product ALTER COLUMN price TYPE NUMERIC(10, 2);`
			if _, err := db.Exec(alterPriceQuery); err != nil {
				return fmt.Errorf("failed to alter price column type: %w", err)
			}
			log.Println("price column migrated to numeric successfully")
		}
	}

	log.Println("Database migrations completed")
	return nil
}
