package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

const (
	InternalServerError = "Internal server error"
	NotFound            = "Country not found"
	ValidationFailed    = "Validation failed"
)

// Refresh Countries data in database
func (m *CountryModel) refreshCountriesHandler(ctx *gin.Context) {
	countries, err := fetchCountries()
	if err != nil {
		log.Print(errorResponse(err))
		err1 := fmt.Errorf("External data source unavailable")
		ctx.JSON(http.StatusServiceUnavailable, specialErrorResponse(err1, err))
		return
	}

	exchangeRatesResponse, err := fetchExchangeRates()
	if err != nil {
		log.Print(errorResponse(err))
		ctx.JSON(http.StatusInternalServerError, errorResponse(errors.New(InternalServerError)))
		return
	}

	err = m.Insert(countries, exchangeRatesResponse)
	if err != nil {
		log.Print(errorResponse(err))
		ctx.JSON(http.StatusInternalServerError, errorResponse(errors.New(InternalServerError)))
		return
	}

	err = m.getImage("cache/summary.png")
	if err != nil {
		log.Print(errorResponse(err))
		ctx.JSON(http.StatusInternalServerError, errorResponse(errors.New(InternalServerError)))
		return
	}

	value := "yay success!"

	ctx.JSON(http.StatusOK, value)
}

// Retrieve a country in the database
func (m *CountryModel) getCountryHandler(ctx *gin.Context) {
	var req getCountryRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		log.Print(errorResponse(err))
		ctx.JSON(http.StatusBadRequest, errorResponse(errors.New(ValidationFailed)))
		return
	}

	value, err := m.GetCountry(req.Name)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ctx.JSON(http.StatusNotFound, errorResponse(errors.New(NotFound)))
			return
		}
		log.Print(errorResponse(err))
		ctx.JSON(http.StatusInternalServerError, errorResponse(errors.New(InternalServerError)))
		return
	}

	ctx.JSON(http.StatusOK, value)
}

func (m *CountryModel) getCountryWithParamsHandler(ctx *gin.Context) {
	var req getCountryWithParamRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		log.Print(errorResponse(err))
		ctx.JSON(http.StatusBadRequest, errorResponse(errors.New(ValidationFailed)))
		return
	}

	result, err := m.getWithParams(req)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Print(errorResponse(err))
			ctx.JSON(http.StatusNotFound, errorResponse(errors.New(NotFound)))
			return
		}
		log.Print(errorResponse(err))
		ctx.JSON(http.StatusInternalServerError, errorResponse(errors.New(InternalServerError)))
		return
	}

	if len(result) == 0 {
		ctx.JSON(http.StatusNotFound, errorResponse(errors.New(NotFound)))
		return
	}

	ctx.JSON(http.StatusOK, result)
}

func (m *CountryModel) GetImageHandler(ctx *gin.Context) {
	imagePath := "cache/summary.png"

	_, err := os.Stat(imagePath)
	if os.IsNotExist(err) {
		error := errors.New("Summary image not found")
		log.Print(errorResponse(err))
		ctx.JSON(http.StatusNotFound, errorResponse(error))
		return
	}
	ctx.File(imagePath)
}

func (m *CountryModel) Insert(countriesAPiResponseList []CountryApiResponse, exchangeRatesResponse *ExchangeRateResponse) error {
	for _, country := range countriesAPiResponseList {
		if len(country.Currencies) == 0 {
			err := m.InsertCountry(country, nil)
			if err != nil {
				errLog := fmt.Errorf("Failed to fetch insert countries without currencies: %w", err)
				return errLog
			}
		} else {
			err := m.InsertCountry(country, exchangeRatesResponse)
			if err != nil {
				errLog := fmt.Errorf("Failed to fetch insert countries: %w", err)
				return errLog
			}
		}
	}
	return nil
}
