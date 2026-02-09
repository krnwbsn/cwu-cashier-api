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

	createProductTable := `
	CREATE TABLE IF NOT EXISTS product (
		id SERIAL PRIMARY KEY,
		name VARCHAR(100) NOT NULL,
		price INT NOT NULL,
		stock INT NOT NULL,
		category_id INT REFERENCES category(id) ON DELETE SET NULL,
		created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
	);`
	if _, err := db.Exec(createProductTable); err != nil {
		return fmt.Errorf("failed to create product table: %w", err)
	}

	createTransactionTable := `
	CREATE TABLE IF NOT EXISTS transaction (
		id SERIAL PRIMARY KEY,
		total_amount INT NOT NULL,
		created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
	);`
	if _, err := db.Exec(createTransactionTable); err != nil {
		return fmt.Errorf("failed to create transaction table: %w", err)
	}

	createTransactionDetailTable := `
	CREATE TABLE IF NOT EXISTS transaction_details (
		id SERIAL PRIMARY KEY,
		transaction_id INT REFERENCES transaction(id) ON DELETE CASCADE,
		product_id INT REFERENCES product(id),
		quantity INT NOT NULL,
		subtotal INT NOT NULL
	);`
	if _, err := db.Exec(createTransactionDetailTable); err != nil {
		return fmt.Errorf("failed to create transaction details table: %w", err)
	}
	log.Println("transaction details table created successfully")

	var totalAmountExists bool
	checkTotalAmountQuery := `
	SELECT EXISTS (
		SELECT 1 FROM information_schema.columns
		WHERE table_schema = 'public' AND table_name = 'transaction' AND column_name = 'total_amount'
	);`
	if err := db.QueryRow(checkTotalAmountQuery).Scan(&totalAmountExists); err != nil {
		return fmt.Errorf("failed to check total_amount column: %w", err)
	}
	if !totalAmountExists {
		log.Println("Adding total_amount column to transaction table...")
		if _, err := db.Exec(`ALTER TABLE "transaction" ADD COLUMN total_amount INT NOT NULL DEFAULT 0`); err != nil {
			return fmt.Errorf("failed to add total_amount column: %w", err)
		}
		log.Println("total_amount column added successfully")
	}

	var createdAtExists bool
	checkCreatedAtQuery := `
	SELECT EXISTS (
		SELECT 1 FROM information_schema.columns
		WHERE table_schema = 'public' AND table_name = 'transaction' AND column_name = 'created_at'
	);`
	if err := db.QueryRow(checkCreatedAtQuery).Scan(&createdAtExists); err != nil {
		return fmt.Errorf("failed to check created_at column: %w", err)
	}
	if !createdAtExists {
		log.Println("Adding created_at column to transaction table...")
		if _, err := db.Exec(`ALTER TABLE "transaction" ADD COLUMN created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP`); err != nil {
			return fmt.Errorf("failed to add created_at column: %w", err)
		}
		log.Println("created_at column added successfully")
	}

	// Make all transaction columns nullable except id, total_amount, created_at
	// so checkout can INSERT only (total_amount) without violating NOT NULL on extra columns.
	rows, err := db.Query(`
		SELECT column_name FROM information_schema.columns
		WHERE table_schema = 'public' AND table_name = 'transaction' AND is_nullable = 'NO'
		AND column_name NOT IN ('id', 'total_amount', 'created_at')
	`)
	if err != nil {
		return fmt.Errorf("failed to list transaction columns: %w", err)
	}
	var cols []string
	for rows.Next() {
		var col string
		if err := rows.Scan(&col); err != nil {
			rows.Close()
			return fmt.Errorf("failed to scan column name: %w", err)
		}
		cols = append(cols, col)
	}
	rows.Close()
	for _, col := range cols {
		// Quote column name for reserved words / safety
		q := fmt.Sprintf(`ALTER TABLE "transaction" ALTER COLUMN "%s" DROP NOT NULL`, col)
		if _, err := db.Exec(q); err != nil {
			log.Printf("Warning: could not make transaction.%s nullable: %v", col, err)
		} else {
			log.Printf("transaction.%s is now nullable", col)
		}
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
