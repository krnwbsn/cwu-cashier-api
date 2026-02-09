package handlers

import (
	"cashier-api/models"
	"cashier-api/services"
	"cashier-api/utils"
	"encoding/json"
	"net/http"
)

type ProductHandler struct {
	service services.ProductServiceInput
}

func NewProductHandler(service services.ProductServiceInput) *ProductHandler {
	return &ProductHandler{service: service}
}

func (h *ProductHandler) HandleProducts(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.GetAll(w, r)
	case http.MethodPost:
		h.Create(w, r)
	default:
		utils.Error(w, http.StatusMethodNotAllowed, "Method not allowed")
	}
}

func (h *ProductHandler) HandleProductByID(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.GetByID(w, r)
	case http.MethodPut:
		h.Update(w, r)
	case http.MethodDelete:
		h.Delete(w, r)
	default:
		utils.Error(w, http.StatusMethodNotAllowed, "Method not allowed")
	}
}

func (h *ProductHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	page := query.Get("page")
	limit := query.Get("limit")
	name := query.Get("name")

	products, err := h.service.GetAll(page, limit, name)
	if err != nil {
		utils.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.JSON(w, http.StatusOK, products)
}

func (h *ProductHandler) Create(w http.ResponseWriter, r *http.Request) {
	var product models.Product
	err := json.NewDecoder(r.Body).Decode(&product)
	if err != nil {
		utils.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	err = h.service.Create(&product)
	if err != nil {
		utils.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.JSON(w, http.StatusCreated, product)
}

func (h *ProductHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := utils.GetIDFromPath(r, "/api/products/")
	if err != nil {
		utils.Error(w, http.StatusBadRequest, "Invalid product ID")
		return
	}

	product, err := h.service.GetByID(id)
	if err != nil {
		utils.Error(w, http.StatusNotFound, err.Error())
		return
	}

	utils.JSON(w, http.StatusOK, product)
}

func (h *ProductHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := utils.GetIDFromPath(r, "/api/products/")
	if err != nil {
		utils.Error(w, http.StatusBadRequest, "Invalid product ID")
		return
	}

	var product models.Product
	err = json.NewDecoder(r.Body).Decode(&product)
	if err != nil {
		utils.Error(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	product.ID = id
	err = h.service.Update(&product)
	if err != nil {
		utils.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.JSON(w, http.StatusOK, product)
}

func (h *ProductHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := utils.GetIDFromPath(r, "/api/products/")
	if err != nil {
		utils.Error(w, http.StatusBadRequest, "Invalid product ID")
		return
	}

	err = h.service.Delete(id)
	if err != nil {
		utils.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.JSON(w, http.StatusOK, map[string]string{"message": "product deleted"})
}
