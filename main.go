package main

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
)

type Product struct {
	ID    int     `json:"id"`
	Name  string  `json:"name"`
	Price float64 `json:"price"`
	Stock int     `json:"stock"`
}

var products = []Product{
	{ID: 1, Name: "Product 1", Price: 100, Stock: 10},
	{ID: 2, Name: "Product 2", Price: 200, Stock: 20},
	{ID: 3, Name: "Product 3", Price: 300, Stock: 30},
}

var productsMap = make(map[int]Product)
var productIndexMap = make(map[int]int)

func init() {
	for i, p := range products {
		productsMap[p.ID] = p
		productIndexMap[p.ID] = i
	}
}

func respondJSON(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		_ = err
	}
}

func respondError(w http.ResponseWriter, statusCode int, message string) {
	respondJSON(w, statusCode, map[string]string{"error": message})
}

func parseProductID(path string) (int, error) {
	idStr := strings.TrimPrefix(path, "/api/products/")
	return strconv.Atoi(idStr)
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	respondJSON(w, http.StatusOK, map[string]string{
		"status":  "OK",
		"message": "API Running",
	})
}

func getProductByID(w http.ResponseWriter, r *http.Request) {
	id, err := parseProductID(r.URL.Path)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid Product ID")
		return
	}

	p, exists := productsMap[id]
	if !exists {
		respondError(w, http.StatusNotFound, "Product not found")
		return
	}

	respondJSON(w, http.StatusOK, p)
}

func getAllProducts(w http.ResponseWriter, _ *http.Request) {
	respondJSON(w, http.StatusOK, products)
}

func createProduct(w http.ResponseWriter, r *http.Request) {
	var newProduct Product
	if err := json.NewDecoder(r.Body).Decode(&newProduct); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request")
		return
	}

	newProduct.ID = len(products) + 1
	products = append(products, newProduct)
	productsMap[newProduct.ID] = newProduct
	productIndexMap[newProduct.ID] = len(products) - 1

	respondJSON(w, http.StatusCreated, newProduct)
}

func updateProduct(w http.ResponseWriter, r *http.Request) {
	id, err := parseProductID(r.URL.Path)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid Product ID")
		return
	}

	if _, exists := productsMap[id]; !exists {
		respondError(w, http.StatusNotFound, "Product not found")
		return
	}

	index, exists := productIndexMap[id]
	if !exists {
		respondError(w, http.StatusInternalServerError, "Product index not found")
		return
	}

	var updatedProduct Product
	if err := json.NewDecoder(r.Body).Decode(&updatedProduct); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request")
		return
	}

	updatedProduct.ID = id

	products[index] = updatedProduct
	productsMap[id] = updatedProduct

	respondJSON(w, http.StatusOK, updatedProduct)
}

func deleteProduct(w http.ResponseWriter, r *http.Request) {
	id, err := parseProductID(r.URL.Path)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid Product ID")
		return
	}

	if _, exists := productsMap[id]; !exists {
		respondError(w, http.StatusNotFound, "Product not found")
		return
	}

	index, exists := productIndexMap[id]
	if !exists {
		respondError(w, http.StatusInternalServerError, "Product index not found")
		return
	}

	products = append(products[:index], products[index+1:]...)

	for pid, idx := range productIndexMap {
		if idx > index {
			productIndexMap[pid] = idx - 1
		}
	}

	delete(productsMap, id)
	delete(productIndexMap, id)

	respondJSON(w, http.StatusOK, map[string]string{"message": "Product deleted"})
}

func productsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getAllProducts(w, r)
	case http.MethodPost:
		createProduct(w, r)
	default:
		respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
	}
}

func productByIDHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getProductByID(w, r)
	case http.MethodPut:
		updateProduct(w, r)
	case http.MethodDelete:
		deleteProduct(w, r)
	default:
		respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
	}
}

func main() {
	http.HandleFunc("/health", healthHandler)
	http.HandleFunc("/api/produk/", productByIDHandler)
	http.HandleFunc("/api/produk", productsHandler)

	http.ListenAndServe(":8080", nil)
}
