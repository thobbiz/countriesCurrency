package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Println(" No .env file found, using environment variables from system.")
	}
}

func main() {
	_ = godotenv.Load()

	dsn := os.Getenv("DB_SOURCE")
	if dsn == "" {
		log.Fatal("❌ DB_SOURCE not set")
	}

	db, err := openDB(dsn)
	if err != nil {
		log.Fatal(err)
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
	router.Run(":7070")
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}

func specialErrorResponse(err1 error, err2 error) gin.H {
	return gin.H{"error": err1.Error(), "details": err2.Error()}
}
