package main

import (
	"database/sql"
	"time"
)

type CountryModel struct {
	DB *sql.DB
}

type Country struct {
	Id              int64
	Name            string
	Capital         string
	Region          string
	Population      int64
	CurrencyCode    *string
	ExchangeRate    *float64
	EstimatedGDP    float64
	FlagURL         string
	LastRefreshedAt time.Time
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
