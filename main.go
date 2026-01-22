package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	_ "belajar-go/docs"

	httpSwagger "github.com/swaggo/http-swagger"
)

// Category represents a product category
type Category struct {
	ID          int    `json:"id" example:"1"`
	Name        string `json:"name" example:"Elektronik"`
	Description string `json:"description" example:"Produk elektronik dan gadget"`
}

// ApiResponse represents the standard API response format
type ApiResponse struct {
	Status  string      `json:"status" example:"success"`
	Message string      `json:"message" example:"Success"`
	Data    interface{} `json:"data"`
}

var categories = []Category{
	{ID: 1, Name: "Elektronik", Description: "Produk elektronik dan gadget"},
	{ID: 2, Name: "Fashion", Description: "Pakaian dan aksesoris"},
	{ID: 3, Name: "Makanan", Description: "Makanan dan minuman"},
}

// GetCategoryByID godoc
// @Summary      Get category by ID
// @Description  Get a single category by its ID
// @Tags         categories
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Category ID"
// @Success      200  {object}  ApiResponse{data=Category}
// @Failure      400  {object}  ApiResponse
// @Failure      404  {object}  ApiResponse
// @Router       /categories/{id} [get]
func getCategoryByID(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/categories/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ApiResponse{
			Status:  "error",
			Message: "Invalid Category ID",
			Data:    nil,
		})
		return
	}

	for _, c := range categories {
		if c.ID == id {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(ApiResponse{
				Status:  "success",
				Message: "Success get category",
				Data:    c,
			})
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotFound)
	json.NewEncoder(w).Encode(ApiResponse{
		Status:  "error",
		Message: "Category not found",
		Data:    nil,
	})
}

// UpdateCategory godoc
// @Summary      Update category
// @Description  Update an existing category by ID
// @Tags         categories
// @Accept       json
// @Produce      json
// @Param        id   path      int      true  "Category ID"
// @Param        category body Category true  "Category object"
// @Success      200  {object}  ApiResponse{data=Category}
// @Failure      400  {object}  ApiResponse
// @Failure      404  {object}  ApiResponse
// @Router       /categories/{id} [put]
func updateCategory(w http.ResponseWriter, r *http.Request) {
	// ambil id dari URL
	idStr := strings.TrimPrefix(r.URL.Path, "/categories/")

	// konversi string jadi int
	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ApiResponse{
			Status:  "error",
			Message: "Invalid Category ID",
			Data:    nil,
		})
		return
	}

	// baca data baru dari body request
	var updatedCategory Category
	err = json.NewDecoder(r.Body).Decode(&updatedCategory)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ApiResponse{
			Status:  "error",
			Message: "Invalid request body",
			Data:    nil,
		})
		return
	}

	// cari kategori yang mau diupdate, terus ganti datanya
	for i := range categories {
		if categories[i].ID == id {
			updatedCategory.ID = id
			categories[i] = updatedCategory

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(ApiResponse{
				Status:  "success",
				Message: "Success update category",
				Data:    updatedCategory,
			})
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotFound)
	json.NewEncoder(w).Encode(ApiResponse{
		Status:  "error",
		Message: "Category not found",
		Data:    nil,
	})
}

// DeleteCategory godoc
// @Summary      Delete category
// @Description  Delete a category by ID
// @Tags         categories
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Category ID"
// @Success      200  {object}  ApiResponse
// @Failure      400  {object}  ApiResponse
// @Failure      404  {object}  ApiResponse
// @Router       /categories/{id} [delete]
func deleteCategory(w http.ResponseWriter, r *http.Request) {
	// ambil id dari URL
	idStr := strings.TrimPrefix(r.URL.Path, "/categories/")
	// ubah jadi integer
	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ApiResponse{
			Status:  "error",
			Message: "Invalid Category ID",
			Data:    nil,
		})
		return
	}
	// cari kategori yang mau dihapus
	for i, c := range categories {
		if c.ID == id {
			// hapus dengan cara gabungin data sebelum & sesudah index
			categories = append(categories[:i], categories[i+1:]...)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(ApiResponse{
				Status:  "success",
				Message: "Success delete category",
				Data:    nil,
			})

			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotFound)
	json.NewEncoder(w).Encode(ApiResponse{
		Status:  "error",
		Message: "Category not found",
		Data:    nil,
	})
}

// GetCategories godoc
// @Summary      List all categories
// @Description  Get list of all categories
// @Tags         categories
// @Accept       json
// @Produce      json
// @Success      200  {object}  ApiResponse{data=[]Category}
// @Router       /categories [get]
func getCategories(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ApiResponse{
		Status:  "success",
		Message: "Success get categories",
		Data:    categories,
	})
}

// CreateCategory godoc
// @Summary      Create a new category
// @Description  Create a new category with the provided information
// @Tags         categories
// @Accept       json
// @Produce      json
// @Param        category body Category true "Category object"
// @Success      201  {object}  ApiResponse{data=Category}
// @Failure      400  {object}  ApiResponse
// @Router       /categories [post]
func createCategory(w http.ResponseWriter, r *http.Request) {
	var newCategory Category
	err := json.NewDecoder(r.Body).Decode(&newCategory)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ApiResponse{
			Status:  "error",
			Message: "Invalid request body",
			Data:    nil,
		})
		return
	}

	newCategory.ID = len(categories) + 1
	categories = append(categories, newCategory)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(ApiResponse{
		Status:  "success",
		Message: "Success create category",
		Data:    newCategory,
	})
}

// HealthCheck godoc
// @Summary      Health check endpoint
// @Description  Check if the API is running
// @Tags         health
// @Accept       json
// @Produce      json
// @Success      200  {object}  ApiResponse
// @Router       /up [get]
func healthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ApiResponse{
		Status:  "success",
		Message: "API Running",
		Data:    nil,
	})
}

// @title           Category API
// @version         1.0
// @description     API untuk mengelola kategori produk
// @host      localhost:8080
// @BasePath  /
func main() {

	// endpoint buat detail kategori (GET, PUT, DELETE)
	http.HandleFunc("/categories/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			getCategoryByID(w, r)
		} else if r.Method == "PUT" {
			updateCategory(w, r)
		} else if r.Method == "DELETE" {
			deleteCategory(w, r)
		}
	})

	// endpoint buat list & create kategori
	http.HandleFunc("/categories", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			getCategories(w, r)
		} else if r.Method == "POST" {
			createCategory(w, r)
		}
	})

	// health check buat ngecek API masih hidup atau ga
	http.HandleFunc("/up", healthCheck)

	// Swagger documentation endpoint
	http.HandleFunc("/swagger/", httpSwagger.WrapHandler)
	fmt.Println("Server running di localhost:8080")

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("gagal running server")
	}
}
