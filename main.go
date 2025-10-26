package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

func main() {

	fmt.Println("--- Printing All Environment Variables ---")
	for _, e := range os.Environ() {
		pair := strings.SplitN(e, "=", 2)
		fmt.Println(pair[0], "=", pair[1])
	}
	fmt.Println("------------------------------------")

	dsn := os.Getenv("DB_SOURCE")
	if dsn == "" {
		log.Fatal("❌ DB_SOURCE not set")
	}

	port := os.Getenv("PORT")
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
