package main

import (
	"database/sql"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"math/rand"
	"os"
	"path/filepath"

	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"
)

func (m *CountryModel) InsertCountry(country CountryApiResponse, exchangeRateResponse *ExchangeRateResponse) error {
	var (
		currencyCode *string
		exRate       *float64
		estimatedGDP *float64
	)

	if len(country.Currencies) > 0 {
		currency := country.Currencies[0]
		currencyCode = &currency.Code

		// Try to get exchange rate
		if exchangeRate, exists := exchangeRateResponse.Rates[*currencyCode]; exists {
			exRate = &exchangeRate

			// Calculate GDP
			randomMultiplier := float64(1000 + rand.Intn(1001))
			gdp := float64(country.Population) * randomMultiplier / *exRate
			estimatedGDP = &gdp
		}
	}

	stmt := `INSERT INTO countries (
		name, capital, region, population,
		currency_code, exchange_rate, estimated_gdp, flag_url
	)
	VALUES(?, ?, ?, ?, ?, ?, ?, ?)
	ON DUPLICATE KEY UPDATE
    capital = VALUES(capital),
    region = VALUES(region),
    population = VALUES(population),
    currency_code = VALUES(currency_code),
    exchange_rate = VALUES(exchange_rate),
    estimated_gdp = VALUES(estimated_gdp),
    flag_url = VALUES(flag_url);
    `

	// Set exchangeRate and EstimatedGDP to nullby not changing their defaut values
	_, err := m.DB.Exec(stmt,
		country.Name,
		country.Capital,
		country.Region,
		country.Population,
		currencyCode,
		exRate,
		estimatedGDP,
		country.FlagURL)
	return err
}

func (m *CountryModel) GetCountry(
	name string,
) (*Country, error) {
	stmt := `SELECT id, name, capital, region, population, currency_code,
 	exchange_rate, estimated_gdp, flag_url, last_refreshed_at
	FROM countries
    WHERE name = ?`

	row := m.DB.QueryRow(stmt, name)

	c := &Country{}

	err := row.Scan(&c.Id, &c.Name, &c.Capital, &c.Region, &c.Population, &c.CurrencyCode, &c.ExchangeRate, &c.EstimatedGDP, &c.FlagURL, &c.LastRefreshedAt)
	if err != nil {
		return nil, err
	}

	return c, nil
}

func (m *CountryModel) getWithParams(params getCountryWithParamRequest) ([]Country, error) {
	query := `SELECT name, capital, region, population, currency_code,
            exchange_rate, COALESCE(estimated_gdp, 0), flag_url, last_refreshed_at
            FROM countries`

	conditions := []string{}

	args := []interface{}{}

	if params.Region != "" {
		conditions = append(conditions, "region = ?")
		args = append(args, params.Region)
	}
	if params.Currency != "" {
		conditions = append(conditions, "currency_code = ?")
		args = append(args, params.Currency)
	}

	if len(conditions) > 0 {
		query += " WHERE " + conditions[0]
		for i := 1; i < len(conditions); i++ {
			query += " AND " + conditions[i]
		}
	}

	// Sorting
	switch params.Sort {
	case "gdp_desc":
		query += " ORDER BY estimated_gdp DESC"
	case "gdp_asc":
		query += " ORDER BY estimated_gdp ASC"
	}

	rows, err := m.DB.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var countries []Country
	i := 1
	for rows.Next() {
		var c Country
		err := rows.Scan(
			&c.Name, &c.Capital, &c.Region, &c.Population,
			&c.CurrencyCode, &c.ExchangeRate, &c.EstimatedGDP, &c.FlagURL, &c.LastRefreshedAt,
		)
		c.Id = int64(i)
		i = i + 1
		if err != nil {
			return nil, err
		}
		countries = append(countries, c)
	}

	return countries, nil
}

func (m *CountryModel) getImage(outPath string) error {

	// Collect Name of top 5 Countries by gdp
	query := `SELECT name
            FROM countries ORDER BY estimated_gdp DESC LIMIT 5 `

	rows, err := m.DB.Query(query)
	if err != nil {
		return err
	}
	defer rows.Close()

	top5 := []string{}
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return err
		}
		top5 = append(top5, name)
	}

	// Collect total number of Countries
	var total int
	if err = m.DB.QueryRow(`SELECT COUNT(*) FROM countries`).Scan(&total); err != nil {
		return err
	}

	// Retrieve max last refreshed at timestamp
	var lastRef sql.NullTime
	if err := m.DB.QueryRow(`SELECT MAX(last_refreshed_at) FROM countries`).Scan(&lastRef); err != nil {
		return fmt.Errorf("last refresh query: %w", err)
	}
	var lastRefStr string
	if lastRef.Valid {
		lastRefStr = lastRef.Time.UTC().Format("2006-01-02 15:04:05 UTC")
	} else {
		lastRefStr = "N/A"
	}

	imgW, imgH := 900, 500
	img := image.NewRGBA(image.Rect(0, 0, imgW, imgH))

	// White background
	draw.Draw(img, img.Bounds(), &image.Uniform{color.White}, image.Point{}, draw.Src)

	// Helper for drawing text
	drawText := func(x, y int, text string) {
		d := &font.Drawer{
			Dst:  img,
			Src:  image.NewUniform(color.Black), // black text
			Face: basicfont.Face7x13,
			Dot:  fixed.P(x, y),
		}
		d.DrawString(text)
	}

	// Add text
	y := 80
	drawText(50, y, fmt.Sprintf("Total no of countries: %d", total))
	y += 40
	drawText(50, y, "Top 5 Countries with Highest GDP:")
	for i, c := range top5 {
		y += 30
		drawText(80, y, fmt.Sprintf("%d. %s", i+1, c))
	}
	y += 50
	drawText(50, y, fmt.Sprintf("Last refreshed at: %s", lastRefStr))

	if err := os.MkdirAll(filepath.Dir(outPath), 0o755); err != nil {
		return err
	}

	// Save image
	f, err := os.Create(outPath)
	if err != nil {
		return err
	}
	defer f.Close()

	return png.Encode(f, img)
}
