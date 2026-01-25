package handlers

import (
	"cashier-api/models"
	"cashier-api/services"
	"cashier-api/utils"
	"encoding/json"
	"net/http"
)

type CategoryHandler struct {
	service services.CategoryServiceInput
}

func NewCategoryHandler(service services.CategoryServiceInput) *CategoryHandler {
	return &CategoryHandler{service: service}
}

func (h *CategoryHandler) HandleCategories(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.GetAll(w, r)
	case http.MethodPost:
		h.Create(w, r)
	default:
		utils.Error(w, http.StatusMethodNotAllowed, "Method not allowed")
	}
}

func (h *CategoryHandler) HandleCategoryByID(w http.ResponseWriter, r *http.Request) {
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

func (h *CategoryHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	categories, err := h.service.GetAll()
	if err != nil {
		utils.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.JSON(w, http.StatusOK, categories)
}

func (h *CategoryHandler) Create(w http.ResponseWriter, r *http.Request) {
	var category models.Category
	if err := json.NewDecoder(r.Body).Decode(&category); err != nil {
		utils.Error(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if err := h.service.Create(&category); err != nil {
		utils.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.JSON(w, http.StatusCreated, category)
}

func (h *CategoryHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := utils.GetIDFromPath(r, "/api/categories/")
	if err != nil {
		utils.Error(w, http.StatusBadRequest, "Invalid category ID")
		return
	}

	category, err := h.service.GetByID(id)
	if err != nil {
		utils.Error(w, http.StatusNotFound, err.Error())
		return
	}

	utils.JSON(w, http.StatusOK, category)
}

func (h *CategoryHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := utils.GetIDFromPath(r, "/api/categories/")
	if err != nil {
		utils.Error(w, http.StatusBadRequest, "Invalid category ID")
		return
	}

	var category models.Category
	if err := json.NewDecoder(r.Body).Decode(&category); err != nil {
		utils.Error(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	category.ID = id

	if err := h.service.Update(&category); err != nil {
		utils.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.JSON(w, http.StatusOK, category)
}

func (h *CategoryHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := utils.GetIDFromPath(r, "/api/categories/")
	if err != nil {
		utils.Error(w, http.StatusBadRequest, "Invalid category ID")
		return
	}

	if err := h.service.Delete(id); err != nil {
		utils.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
