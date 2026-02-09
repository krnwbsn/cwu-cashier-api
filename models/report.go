package models

type BestSellingProduct struct {
	Name           string `json:"name"`
	QuantitySold   int    `json:"quantity_sold"`
}

type SalesSummary struct {
	TotalRevenue        int                   `json:"total_revenue"`
	TotalTransactions   int                   `json:"total_transactions"`
	BestSellingProduct *BestSellingProduct   `json:"best_selling_product,omitempty"`
}
