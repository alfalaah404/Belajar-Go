package handlers

import (
	"belajar-go/models"
	"belajar-go/services"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
)

type CategoryHandler struct {
	service *services.CategoryService
}

func NewCategoryHandler(service *services.CategoryService) *CategoryHandler {
	return &CategoryHandler{service: service}
}

// HandleCategories - GET /api/categories, POST /api/categories
func (h *CategoryHandler) HandleCategories(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.GetAll(w, r)
	case http.MethodPost:
		h.Create(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// GetAll godoc
// @Summary      List all categories
// @Description  Get list of all categories
// @Tags         categories
// @Accept       json
// @Produce      json
// @Success      200  {object}  models.ApiResponse{data=[]models.Category}
// @Failure      500  {object}  models.ApiResponse
// @Router       /api/categories [get]
func (h *CategoryHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	categories, err := h.service.GetAll()
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(models.ApiResponse{
			Status:  "error",
			Message: err.Error(),
			Data:    nil,
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(models.ApiResponse{
		Status:  "success",
		Message: "Success get categories",
		Data:    categories,
	})
}

// Create godoc
// @Summary      Create a new category
// @Description  Create a new category with the provided information
// @Tags         categories
// @Accept       json
// @Produce      json
// @Param        category body models.Category true "Category object"
// @Success      201  {object}  models.ApiResponse{data=models.Category}
// @Failure      400  {object}  models.ApiResponse
// @Router       /api/categories [post]
func (h *CategoryHandler) Create(w http.ResponseWriter, r *http.Request) {
	var category models.Category
	err := json.NewDecoder(r.Body).Decode(&category)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.ApiResponse{
			Status:  "error",
			Message: "Invalid request body",
			Data:    nil,
		})
		return
	}

	err = h.service.Create(&category)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.ApiResponse{
			Status:  "error",
			Message: err.Error(),
			Data:    nil,
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(models.ApiResponse{
		Status:  "success",
		Message: "Success create category",
		Data:    category,
	})
}

// HandleCategoryByID - GET/PUT/DELETE /api/categories/{id}
func (h *CategoryHandler) HandleCategoryByID(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.GetByID(w, r)
	case http.MethodPut:
		h.Update(w, r)
	case http.MethodDelete:
		h.Delete(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// GetByID godoc
// @Summary      Get category by ID
// @Description  Get a single category by its ID
// @Tags         categories
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Category ID"
// @Success      200  {object}  models.ApiResponse{data=models.Category}
// @Failure      400  {object}  models.ApiResponse
// @Failure      404  {object}  models.ApiResponse
// @Router       /api/categories/{id} [get]
func (h *CategoryHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/categories/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.ApiResponse{
			Status:  "error",
			Message: "Invalid Category ID",
			Data:    nil,
		})
		return
	}

	category, err := h.service.GetByID(id)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(models.ApiResponse{
			Status:  "error",
			Message: err.Error(),
			Data:    nil,
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(models.ApiResponse{
		Status:  "success",
		Message: "Success get category",
		Data:    category,
	})
}

// Update godoc
// @Summary      Update category
// @Description  Update an existing category by ID
// @Tags         categories
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Category ID"
// @Param        category body models.Category true  "Category object"
// @Success      200  {object}  models.ApiResponse{data=models.Category}
// @Failure      400  {object}  models.ApiResponse
// @Failure      404  {object}  models.ApiResponse
// @Router       /api/categories/{id} [put]
func (h *CategoryHandler) Update(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/categories/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.ApiResponse{
			Status:  "error",
			Message: "Invalid Category ID",
			Data:    nil,
		})
		return
	}

	var category models.Category
	err = json.NewDecoder(r.Body).Decode(&category)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.ApiResponse{
			Status:  "error",
			Message: "Invalid request body",
			Data:    nil,
		})
		return
	}

	category.ID = id
	err = h.service.Update(&category)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(models.ApiResponse{
			Status:  "error",
			Message: err.Error(),
			Data:    nil,
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(models.ApiResponse{
		Status:  "success",
		Message: "Success update category",
		Data:    category,
	})
}

// Delete godoc
// @Summary      Delete category
// @Description  Delete a category by ID
// @Tags         categories
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Category ID"
// @Success      200  {object}  models.ApiResponse
// @Failure      400  {object}  models.ApiResponse
// @Failure      404  {object}  models.ApiResponse
// @Router       /api/categories/{id} [delete]
func (h *CategoryHandler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/categories/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.ApiResponse{
			Status:  "error",
			Message: "Invalid Category ID",
			Data:    nil,
		})
		return
	}

	err = h.service.Delete(id)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(models.ApiResponse{
			Status:  "error",
			Message: err.Error(),
			Data:    nil,
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(models.ApiResponse{
		Status:  "success",
		Message: "Success delete category",
		Data:    nil,
	})
}
