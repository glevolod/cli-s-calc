package converter

import (
	"fmt"
	"os"
	"strconv"
)

type Period string

const (
	PeriodHour  Period = "Hour"
	PeriodDay   Period = "Day"
	PeriodMonth Period = "Month"
	PeriodYear  Period = "Year"
)

type Currency string

const (
	CurrencyPLN Currency = "PLN"
	CurrencyEUR Currency = "EUR"
	CurrencyUSD Currency = "USD"
	CurrencyGBP Currency = "GBP"
)

var ValidCurrencies = []Currency{CurrencyPLN, CurrencyEUR, CurrencyUSD, CurrencyGBP}

var ValidPeriods = []Period{PeriodHour, PeriodDay, PeriodMonth, PeriodYear}

type Input struct {
	Amount   float64
	Period   Period
	Currency Currency
}

type Converter struct {
	hoursPerDay  int
	daysPerMonth int
	rates        map[string]float64
	baseCurrency string
}

func NewConverter(rates map[string]float64, baseCurrency string) *Converter {
	hoursPerDay := 8
	if hoursEnv := os.Getenv("S_HOURS_DAY"); hoursEnv != "" {
		if parsed, err := strconv.ParseUint(hoursEnv, 10, 32); err == nil && parsed > 0 {
			hoursPerDay = int(parsed)
		}
	}

	daysPerMonth := 21
	if daysEnv := os.Getenv("S_DAYS_MONTH"); daysEnv != "" {
		if parsed, err := strconv.ParseUint(daysEnv, 10, 32); err == nil && parsed > 0 {
			daysPerMonth = int(parsed)
		}
	}

	return &Converter{
		hoursPerDay:  hoursPerDay,
		daysPerMonth: daysPerMonth,
		rates:        rates,
		baseCurrency: baseCurrency,
	}
}

func (c *Converter) Convert(input Input) map[Period]map[Currency]float64 {
	baseHourly := c.toHourly(input.Amount, input.Period)

	result := make(map[Period]map[Currency]float64)
	for _, period := range ValidPeriods {
		result[period] = make(map[Currency]float64)
		for _, currency := range ValidCurrencies {
			// Convert currency first
			amountInCurrency := baseHourly * c.getRate(string(input.Currency), string(currency))
			// Then convert period
			result[period][currency] = c.fromHourly(amountInCurrency, period)
		}
	}

	return result
}

func (c *Converter) toHourly(amount float64, period Period) float64 {
	switch period {
	case PeriodHour:
		return amount
	case PeriodDay:
		return amount / float64(c.hoursPerDay)
	case PeriodMonth:
		return amount / (float64(c.hoursPerDay) * float64(c.daysPerMonth))
	case PeriodYear:
		return amount / (float64(c.hoursPerDay) * float64(c.daysPerMonth) * 12)
	default:
		return amount
	}
}

func (c *Converter) fromHourly(amount float64, period Period) float64 {
	switch period {
	case PeriodHour:
		return amount
	case PeriodDay:
		return amount * float64(c.hoursPerDay)
	case PeriodMonth:
		return amount * float64(c.hoursPerDay) * float64(c.daysPerMonth)
	case PeriodYear:
		return amount * float64(c.hoursPerDay) * float64(c.daysPerMonth) * 12
	default:
		return amount
	}
}

func (c *Converter) getRate(from, to string) float64 {
	if from == to {
		return 1.0
	}

	if c.rates == nil {
		return 1.0
	}

	var fromRate float64 = 1.0
	if from != c.baseCurrency {
		if rate, ok := c.rates[from]; ok {
			fromRate = rate
		} else {
			return 1.0 // Unknown currency
		}
	}

	var toRate float64 = 1.0
	if to != c.baseCurrency {
		if rate, ok := c.rates[to]; ok {
			toRate = rate
		} else {
			return 1.0 // Unknown currency
		}
	}

	// Convert: fromCurrency -> baseCurrency -> toCurrency
	// If fromRate = 4.25 (PLN), toRate = 1.10 (USD), base = EUR
	// To convert 100 PLN to USD: 100 / 4.25 * 1.10 = 100 * (1.10 / 4.25)
	return toRate / fromRate
}

func ValidateCurrency(currency string) (Currency, error) {
	for _, c := range ValidCurrencies {
		if string(c) == currency || string(c) == fmt.Sprintf("%s%s", string(currency[0]-32), currency[1:]) {
			return c, nil
		}
	}
	return "", fmt.Errorf("invalid currency: %s (supported: PLN, EUR, USD, GBP)", currency)
}

func ValidatePeriod(period string) (Period, error) {
	for _, p := range ValidPeriods {
		if string(p) == period || string(p) == fmt.Sprintf("%s%s", string(period[0]-32), period[1:]) {
			return p, nil
		}
	}
	return "", fmt.Errorf("invalid period: %s (supported: Hour, Day, Month, Year)", period)
}
