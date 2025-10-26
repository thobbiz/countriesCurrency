## CountriesCurrecny

CountriesCurrency is a backend service built with Golang that collects information about countries and their exchange rate and stores in a database.

## Features:
- rate limiting using `golang.org/x/time/rate`
- error handling

## Stack and Tools:
- Golang
- Gin Web Framework
- net/http

## Get Started
### Prerequisites
- G0 1.21+
- Internet Connection
- Modules:
  ```bash
  go get github.com/gin-gonic/gin

## Usage
- Start the server:
  ```bash
  go run main.go
- Refresh countries database:
  ```bash
  curl http://localhost:7070/countries/refresh

- Get a Country:
  ```bash
  curl http://localhost:7070/countries/:country_name
  
- Get a Country by Filtering:
  ```bash
  curl http://localhost:7070/countries?region=Africa
