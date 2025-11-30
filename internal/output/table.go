package output

import (
	"fmt"
	"strings"
	"time"

	"salary-calc/internal/converter"
	"salary-calc/internal/exchangerate"
)

type TableFormatter struct {
	originalAmount   float64
	originalPeriod   converter.Period
	originalCurrency converter.Currency
	rateInfo         *exchangerate.RateInfo
}

func NewTableFormatter(amount float64, period converter.Period, currency converter.Currency, rateInfo *exchangerate.RateInfo) *TableFormatter {
	return &TableFormatter{
		originalAmount:   amount,
		originalPeriod:   period,
		originalCurrency: currency,
		rateInfo:         rateInfo,
	}
}

func (tf *TableFormatter) Format(results map[converter.Period]map[converter.Currency]float64) string {
	var sb strings.Builder

	periodWidth := 8
	currencyWidth := 13

	sb.WriteString("┌")
	sb.WriteString(strings.Repeat("─", periodWidth))
	for range converter.ValidCurrencies {
		sb.WriteString("┬")
		sb.WriteString(strings.Repeat("─", currencyWidth))
	}
	sb.WriteString("┐\n")

	// Header row
	sb.WriteString("│")
	sb.WriteString(padCenter("Period", periodWidth))
	for _, currency := range converter.ValidCurrencies {
		sb.WriteString("│")
		sb.WriteString(padCenter(string(currency), currencyWidth))
	}
	sb.WriteString("│\n")

	sb.WriteString("├")
	sb.WriteString(strings.Repeat("─", periodWidth))
	for range converter.ValidCurrencies {
		sb.WriteString("┼")
		sb.WriteString(strings.Repeat("─", currencyWidth))
	}
	sb.WriteString("┤\n")

	// Data rows
	for _, period := range converter.ValidPeriods {
		sb.WriteString("│")
		sb.WriteString(padLeft(string(period), periodWidth))
		for _, currency := range converter.ValidCurrencies {
			sb.WriteString("│")
			value := results[period][currency]
			formattedValue := formatNumber(value)
			isOriginal := period == tf.originalPeriod && currency == tf.originalCurrency
			if isOriginal {
				formattedValue = formattedValue + " ⭐"
			}
			sb.WriteString(padRight(formattedValue, currencyWidth))
		}
		sb.WriteString("│\n")
	}

	// Footer
	sb.WriteString("└")
	sb.WriteString(strings.Repeat("─", periodWidth))
	for range converter.ValidCurrencies {
		sb.WriteString("┴")
		sb.WriteString(strings.Repeat("─", currencyWidth))
	}
	sb.WriteString("┘\n")

	sb.WriteString("\n⭐ Original input: ")
	sb.WriteString(formatNumber(tf.originalAmount))
	sb.WriteString(" ")
	sb.WriteString(string(tf.originalCurrency))
	sb.WriteString("/")
	sb.WriteString(strings.ToLower(string(tf.originalPeriod)))
	sb.WriteString("\n")

	if tf.rateInfo != nil {
		sb.WriteString("\nRate source: ")
		sb.WriteString(tf.rateInfo.Source)
		sb.WriteString("\nLast updated: ")
		sb.WriteString(tf.rateInfo.Timestamp.Format("2006-01-02 15:04:05 UTC"))
		sb.WriteString("\nCache expires: ")
		sb.WriteString(tf.rateInfo.ExpiresAt.Format("2006-01-02 15:04:05 UTC"))
		sb.WriteString("\n")
	}

	return sb.String()
}

func padLeft(s string, width int) string {
	if len(s) >= width {
		return s[:width]
	}
	return strings.Repeat(" ", width-len(s)) + s
}

func padRight(s string, width int) string {
	if len(s) >= width {
		return s[:width]
	}
	return s + strings.Repeat(" ", width-len(s))
}

func padCenter(s string, width int) string {
	if len(s) >= width {
		return s[:width]
	}
	padding := width - len(s)
	left := padding / 2
	right := padding - left
	return strings.Repeat(" ", left) + s + strings.Repeat(" ", right)
}

func formatNumber(n float64) string {
	formatted := fmt.Sprintf("%.2f", n)

	parts := strings.Split(formatted, ".")
	intPart := parts[0]

	if len(intPart) > 3 {
		var result strings.Builder
		start := len(intPart) % 3
		if start > 0 {
			result.WriteString(intPart[:start])
			if start < len(intPart) {
				result.WriteString(",")
			}
		}
		for i := start; i < len(intPart); i += 3 {
			if i > start {
				result.WriteString(",")
			}
			result.WriteString(intPart[i : i+3])
		}
		intPart = result.String()
	}

	if len(parts) > 1 {
		return intPart + "." + parts[1]
	}
	return intPart
}

func FormatVerbose(rateInfo *exchangerate.RateInfo, rates map[string]float64) string {
	if rateInfo == nil {
		return ""
	}

	var sb strings.Builder
	sb.WriteString("\n--- Exchange Rate Details ---\n")
	sb.WriteString(fmt.Sprintf("Source: %s\n", rateInfo.Source))
	sb.WriteString(fmt.Sprintf("Fetched at: %s\n", rateInfo.Timestamp.Format(time.RFC3339)))
	sb.WriteString(fmt.Sprintf("Expires at: %s\n", rateInfo.ExpiresAt.Format(time.RFC3339)))
	sb.WriteString("\nCurrent rates:\n")
	for currency, rate := range rates {
		sb.WriteString(fmt.Sprintf("  %s: %.4f\n", currency, rate))
	}
	return sb.String()
}
