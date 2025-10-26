package main

import (
	"database/sql"
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
	_ "github.com/go-sql-driver/mysql"
)

func main() {

	// dbUser := os.Getenv("DB_USER")
	// dbPass := os.Getenv("DB_PASS")
	// dbHost := os.Getenv("DB_HOST")
	// dbPort := os.Getenv("DB_PORT")
	dsn := os.Getenv("DB_SOURCE")

	log.Printf("--- DEBUG ---")
	log.Printf("DB_SOURCE: [%s]", dsn)
	log.Printf("-------------")

	// dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbUser, dbPass, dbHost, dbPort, dbName)

	cfg, err := mysql.ParseDSN(dsn)
	if err != nil {
		log.Fatalf("❌ Could not parse DSN: %v", err)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Use 7070 as a fallback for local development
	}

	cfg.Net = "tcp"
	cfg.ParseTime = true

	db, err := openDB(cfg.FormatDSN())
	if err != nil {
		log.Print(err)
	}

	model := &CountryModel{
		DB: db,
	}
	defer db.Close()

	router := gin.Default()
	router.POST("/countries/refresh", model.refreshCountriesHandler)
	router.GET("/countries/:name", model.getCountryHandler)
	router.GET("/countries", model.getCountryWithParamsHandler)
	router.GET("/countries/image", model.GetImageHandler)
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "healthy",
			"port":   os.Getenv("PORT"),
			"time":   time.Now(),
		})
	})
	router.Run(":" + port)
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Print("Failed to open database:", err)
		return nil, err
	}
	if err = db.Ping(); err != nil {
		log.Print("Failed to ping database:", err)
		return nil, err
	}
	log.Println("Database connection successful")
	return db, nil
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}

func specialErrorResponse(err1 error, err2 error) gin.H {
	return gin.H{"error": err1.Error(), "details": err2.Error()}
}
