package main

import (
	"database/sql"
	"flag"
	"log"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	dsn := flag.String("dsn", "code:abifoluwa3#@tcp(localhost:3307)/countries?parseTime=true", "MySQL data source name")
	flag.Parse()

	db, err := openDB(*dsn)
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
