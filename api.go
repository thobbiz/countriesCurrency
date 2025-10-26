package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

func fetchCountries() ([]CountryApiResponse, error) {
	API_NAME := "https://restcountries.com/v2/all?fields=name,capital,region,population,flag,currencies"
	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	resp, err := client.Get(API_NAME)
	if err != nil {
		err := fmt.Errorf("Could not fetch data from [%s]", API_NAME)
		return nil, err
	}
	log.Printf("✅ External API responded: %d", resp.StatusCode)
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var countryResponseList []CountryApiResponse
	err = json.Unmarshal(body, &countryResponseList)
	if err != nil {
		return nil, err
	}

	return countryResponseList, nil
}

func fetchExchangeRates() (*ExchangeRateResponse, error) {
	resp, err := http.Get("https://open.er-api.com/v6/latest/USD")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var exchangeRates ExchangeRateResponse
	err = json.Unmarshal(body, &exchangeRates)
	if err != nil {
		return nil, err
	}

	return &exchangeRates, nil
}
