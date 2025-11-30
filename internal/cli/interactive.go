package cli

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"salary-calc/internal/converter"
)

// Interactive prompts user for input
func Interactive() (float64, converter.Period, converter.Currency, error) {
	reader := bufio.NewReader(os.Stdin)

	period, err := promptPeriod(reader)
	if err != nil {
		return 0, "", "", err
	}

	amount, err := promptAmount(reader)
	if err != nil {
		return 0, "", "", err
	}

	currency, err := promptCurrency(reader)
	if err != nil {
		return 0, "", "", err
	}

	return amount, period, currency, nil
}

func promptPeriod(reader *bufio.Reader) (converter.Period, error) {
	for {
		fmt.Print("Select period [Hour/Day/Month/Year]: ")
		input, err := reader.ReadString('\n')
		if err != nil {
			return "", err
		}

		input = strings.TrimSpace(input)
		if input == "" {
			continue
		}

		period, err := converter.ValidatePeriod(input)
		if err == nil {
			return period, nil
		}

		fmt.Printf("Error: %v\n", err)
	}
}

func promptAmount(reader *bufio.Reader) (float64, error) {
	for {
		fmt.Print("Enter amount: ")
		input, err := reader.ReadString('\n')
		if err != nil {
			return 0, err
		}

		input = strings.TrimSpace(input)
		if input == "" {
			continue
		}

		amount, err := parseFloat(input)
		if err != nil {
			fmt.Printf("Error: invalid number: %v\n", err)
			continue
		}

		if amount <= 0 {
			fmt.Println("Error: amount must be positive")
			continue
		}

		return amount, nil
	}
}

func promptCurrency(reader *bufio.Reader) (converter.Currency, error) {
	for {
		fmt.Print("Select currency [PLN/EUR/USD/GBP]: ")
		input, err := reader.ReadString('\n')
		if err != nil {
			return "", err
		}

		input = strings.TrimSpace(input)
		if input == "" {
			continue
		}

		currency, err := converter.ValidateCurrency(input)
		if err == nil {
			return currency, nil
		}

		fmt.Printf("Error: %v\n", err)
	}
}

func parseFloat(s string) (float64, error) {
	s = strings.TrimSpace(s)
	var result float64
	_, err := fmt.Sscanf(s, "%f", &result)
	return result, err
}
