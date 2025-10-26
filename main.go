package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

func main() {

	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASS")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbUser, dbPass, dbHost, dbPort, dbName)

	port := os.Getenv("DB_PORT")
	if port == "" {
		port = "7070" // Use 7070 as a fallback for local development
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
	router.Run(":" + port)
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
