package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"belajar-go/database"
	"belajar-go/docs"
	"belajar-go/handlers"
	"belajar-go/models"
	"belajar-go/repositories"
	"belajar-go/services"

	"github.com/spf13/viper"
	httpSwagger "github.com/swaggo/http-swagger"
)

type Config struct {
	Port          string `mapstructure:"PORT"`
	DBConn        string `mapstructure:"DB_CONN"`
	SwaggerHost   string `mapstructure:"SWAGGER_HOST"`
	SwaggerScheme string `mapstructure:"SWAGGER_SCHEME"`
}

// HealthCheck godoc
// @Summary      Health check endpoint
// @Description  Check if the API is running
// @Tags         health
// @Accept       json
// @Produce      json
// @Success      200  {object}  models.ApiResponse
// @Router       /up [get]
func healthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(models.ApiResponse{
		Status:  "success",
		Message: "API Running",
		Data:    nil,
	})
}

// @title           Product & Category API
// @version         1.0
// @description     REST API untuk mengelola produk dan kategori
// @BasePath        /
// @host            localhost:8080
// @schemes         http
func main() {
	// load environment variables
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if _, err := os.Stat(".env"); err == nil {
		viper.SetConfigFile(".env")
		_ = viper.ReadInConfig()
	}

	// load config
	config := Config{
		Port:          viper.GetString("PORT"),
		DBConn:        viper.GetString("DB_CONN"),
		SwaggerHost:   viper.GetString("SWAGGER_HOST"),
		SwaggerScheme: viper.GetString("SWAGGER_SCHEME"),
	}

	// initialize database
	db, err := database.InitDB(config.DBConn)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	// dependency injection for products
	productRepo := repositories.NewProductRepository(db)
	productService := services.NewProductService(productRepo)
	productHandler := handlers.NewProductHandler(productService)

	// dependency injection for categories
	categoryRepo := repositories.NewCategoryRepository(db)
	categoryService := services.NewCategoryService(categoryRepo)
	categoryHandler := handlers.NewCategoryHandler(categoryService)

	// dependency injection for transactions
	transactionRepo := repositories.NewTransactionRepository(db)
	transactionService := services.NewTransactionService(transactionRepo)
	transactionHandler := handlers.NewTransactionHandler(transactionService)

	// dependency injection for reports
	reportRepo := repositories.NewReportRepository(db)
	reportService := services.NewReportService(reportRepo)
	reportHandler := handlers.NewReportHandler(reportService)

	// Setup routes
	http.HandleFunc("/api/products", productHandler.HandleProducts)
	http.HandleFunc("/api/products/", productHandler.HandleProductByID)
	http.HandleFunc("/api/categories", categoryHandler.HandleCategories)
	http.HandleFunc("/api/categories/", categoryHandler.HandleCategoryByID)
	http.HandleFunc("/api/checkout", transactionHandler.Checkout)
	http.HandleFunc("/api/reports/hari-ini", reportHandler.HandleReports)
	http.HandleFunc("/api/reports", reportHandler.HandleReports)

	// Set Swagger host secara dinamis berdasarkan environment
	host := config.SwaggerHost
	if host == "" {
		host = "localhost:8080" // default untuk local development
	}
	docs.SwaggerInfo.Host = host

	// Set scheme (http/https) secara dinamis
	scheme := config.SwaggerScheme
	if scheme == "" {
		scheme = "http" // default untuk local development
	}
	docs.SwaggerInfo.Schemes = []string{scheme}

	// health check buat ngecek API masih hidup atau ga
	http.HandleFunc("/up", healthCheck)

	// Swagger documentation endpoint
	http.HandleFunc("/swagger/", httpSwagger.WrapHandler)

	// Redirect root path to swagger documentation
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			http.Redirect(w, r, "/swagger/index.html", http.StatusMovedPermanently)
			return
		}
		http.NotFound(w, r)
	})

	// Get port from environment variable (untuk Railway, Render, dll)
	port := config.Port
	if port == "" {
		port = config.Port // default untuk local development
	}

	fmt.Printf("Server running on port %s\n", config.Port)
	fmt.Printf("Swagger UI: http://%s/swagger/index.html\n", config.SwaggerHost)

	err = http.ListenAndServe(":"+config.Port, nil)
	if err != nil {
		fmt.Printf("Failed to start server: %v\n", err)
	}
}
