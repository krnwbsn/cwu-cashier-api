package repositories

import (
	"cashier-api/models"
	"database/sql"
)

type ReportRepositoryInput interface {
	GetSalesSummaryToday() (*models.SalesSummary, error)
	GetSalesSummaryRange(startDate, endDate string) (*models.SalesSummary, error)
}

type ReportRepository struct {
	db *sql.DB
}

func NewReportRepository(db *sql.DB) ReportRepositoryInput {
	return &ReportRepository{db: db}
}

func (repo *ReportRepository) GetSalesSummaryToday() (*models.SalesSummary, error) {
	return repo.getSalesSummary("t.created_at::date = CURRENT_DATE")
}

func (repo *ReportRepository) GetSalesSummaryRange(startDate, endDate string) (*models.SalesSummary, error) {
	return repo.getSalesSummary("t.created_at::date >= $1 AND t.created_at::date <= $2", startDate, endDate)
}

func (repo *ReportRepository) getSalesSummary(whereClause string, args ...interface{}) (*models.SalesSummary, error) {
	summary := &models.SalesSummary{}

	err := repo.db.QueryRow(
		`SELECT COALESCE(SUM(t."total_amount"), 0), COUNT(t.id) FROM "transaction" t WHERE `+whereClause,
		args...,
	).Scan(&summary.TotalRevenue, &summary.TotalTransactions)
	if err != nil {
		return nil, err
	}

	query := `
		SELECT p.name, COALESCE(SUM(td.quantity), 0) AS qty
		FROM transaction_details td
		JOIN product p ON p.id = td.product_id
		JOIN "transaction" t ON t.id = td.transaction_id
		WHERE ` + whereClause + `
		GROUP BY p.id, p.name
		ORDER BY qty DESC
		LIMIT 1
	`
	rows, err := repo.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if rows.Next() {
		var productName string
		var qty int
		if err := rows.Scan(&productName, &qty); err != nil {
			return nil, err
		}
		summary.BestSellingProduct = &models.BestSellingProduct{Name: productName, QuantitySold: qty}
	}

	return summary, nil
}
