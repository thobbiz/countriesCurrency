package main

import (
	"database/sql"
	"time"
)

type CountryModel struct {
	DB *sql.DB
}

type Country struct {
	Id              int64     `json:"id"`
	Name            string    `json:"name"`
	Capital         string    `json:"capital"`
	Region          string    `json:"region"`
	Population      int64     `json:"population"`
	CurrencyCode    *string   `json:"currency_code"`
	ExchangeRate    *float64  `json:"exchange_rate"`
	EstimatedGDP    float64   `json:"estimated_gdp"`
	FlagURL         string    `json:"flag_url"`
	LastRefreshedAt time.Time `json:"last_refreshed_at"`
}

type CountryApiResponse struct {
	Name       string                `json:"name"`
	Capital    string                `json:"capital"`
	Region     string                `json:"region"`
	Population int64                 `json:"population"`
	Currencies []CurrencyApiResponse `json:"currencies"`
	FlagURL    string                `json:"flag"`
}

type CurrencyApiResponse struct {
	Name   string `json:"name"`
	Code   string `json:"code"`
	Symbol string `json:"symbol"`
}

type ExchangeRateResponse struct {
	Rates map[string]float64 `json:"rates"`
}

type getCountryRequest struct {
	Name string `uri:"name" binding:"required"`
}

type getCountryWithParamRequest struct {
	Region   string `form:"region" json:"region"`
	Currency string `form:"currency" json:"currency"`
	Sort     string `form:"sort" json:"sort"`
}
