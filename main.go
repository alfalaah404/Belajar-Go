package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

type Category struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type ApiResponse struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

var categories = []Category{
	{ID: 1, Name: "Elektronik", Description: "Produk elektronik dan gadget"},
	{ID: 2, Name: "Fashion", Description: "Pakaian dan aksesoris"},
	{ID: 3, Name: "Makanan", Description: "Makanan dan minuman"},
}

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
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(ApiResponse{
				Status:  "success",
				Message: "Success get categories",
				Data:    categories,
			})
		} else if r.Method == "POST" {
			// baca data kategori baru dari body
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

			// kasih ID otomatis, terus masukin ke slice
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

	})

	// health check buat ngecek API masih hidup atau ga
	http.HandleFunc("/up", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(ApiResponse{
			Status:  "success",
			Message: "API Running",
			Data:    nil,
		})
	})
	fmt.Println("Server running di localhost:8080")

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("gagal running server")
	}
}
