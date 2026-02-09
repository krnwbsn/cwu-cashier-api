package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"strings"

	"cashier-api/database"
	"cashier-api/handlers"
	"cashier-api/repositories"
	"cashier-api/services"

	"github.com/spf13/viper"
)

type Config struct {
	Port   string `mapstructure:"PORT"`
	DBConn string `mapstructure:"DB_CONN"`
}

func loadConfig() Config {
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if _, err := os.Stat(".env"); err == nil {
		viper.SetConfigFile(".env")
		_ = viper.ReadInConfig()
	}

	config := Config{
		Port:   viper.GetString("PORT"),
		DBConn: viper.GetString("DB_CONN"),
	}

	return config
}

func main() {
	migrateFlag := flag.Bool("migrate", false, "Run database migrations and exit")
	seedFlag := flag.Bool("seed", false, "Run database seeding and exit")
	flag.Parse()

	config := loadConfig()

	db, err := database.InitDB(config.DBConn)
	if err != nil {
		fmt.Println("Failed to connect database", err)
		return
	}
	defer db.Close()

	if *migrateFlag {
		if err := database.Migrate(db); err != nil {
			fmt.Println("Failed to run migrations:", err)
			os.Exit(1)
		}
		fmt.Println("Migration command completed successfully")
		return
	}

	if *seedFlag {
		if err := database.Seed(db); err != nil {
			fmt.Println("Failed to seed database:", err)
			os.Exit(1)
		}
		fmt.Println("Seed command completed successfully")
		return
	}

	if err := database.Migrate(db); err != nil {
		fmt.Println("Failed to run migrations", err)
		return
	}

	productRepo := repositories.NewProductRepository(db)
	productService := services.NewProductService(productRepo)
	productHandler := handlers.NewProductHandler(productService)

	categoryRepo := repositories.NewCategoryRepository(db)
	categoryService := services.NewCategoryService(categoryRepo)
	categoryHandler := handlers.NewCategoryHandler(categoryService)

	transactionRepo := repositories.NewTransactionRepository(db)
	transactionService := services.NewTransactionService(transactionRepo)
	transactionHandler := handlers.NewTransactionHandler(transactionService)

	reportRepo := repositories.NewReportRepository(db)
	reportService := services.NewReportService(reportRepo)
	reportHandler := handlers.NewReportHandler(reportService)

	http.HandleFunc("/api/checkout", transactionHandler.HandleCheckout)
	http.HandleFunc("/api/report/today", reportHandler.HandleReportToday)
	http.HandleFunc("/api/report", reportHandler.HandleReport)
	http.HandleFunc("/api/health", handlers.HealthCheckHandler)
	http.HandleFunc("/api/products", productHandler.HandleProducts)
	http.HandleFunc("/api/products/", productHandler.HandleProductByID)
	http.HandleFunc("/api/categories", categoryHandler.HandleCategories)
	http.HandleFunc("/api/categories/", categoryHandler.HandleCategoryByID)

	if config.Port == "" {
		config.Port = "8080"
	}
	addr := "0.0.0.0:" + config.Port
	fmt.Println("Server running at", addr)

	err = http.ListenAndServe(addr, nil)
	if err != nil {
		fmt.Println("Failed to running server", err)
	}
}
