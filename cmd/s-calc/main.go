package main

import (
	"fmt"
	"os"

	"salary-calc/internal/cli"
	"salary-calc/internal/converter"
	"salary-calc/internal/exchangerate"
	"salary-calc/internal/output"
)

func main() {
	flags, args := cli.ParseFlags()

	var input converter.Input
	var err error

	if flags.HasInput() {
		amount, periodStr, currencyStr, ok := flags.GetInput()
		if !ok {
			fmt.Fprintf(os.Stderr, "Error: invalid input\n")
			os.Exit(1)
		}

		period, err := converter.ValidatePeriod(periodStr)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		currency, err := converter.ValidateCurrency(currencyStr)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		input = converter.Input{
			Amount:   amount,
			Period:   period,
			Currency: currency,
		}
	} else if len(args) > 0 {
		amount, periodStr, currencyStr, ok := cli.ParseLegacyFormat(args)
		if ok {
			period, err := converter.ValidatePeriod(periodStr)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error: %v\n", err)
				os.Exit(1)
			}

			currency, err := converter.ValidateCurrency(currencyStr)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error: %v\n", err)
				os.Exit(1)
			}

			input = converter.Input{
				Amount:   amount,
				Period:   period,
				Currency: currency,
			}
		} else {
			amount, period, currency, err := cli.Interactive()
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error: %v\n", err)
				os.Exit(1)
			}

			input = converter.Input{
				Amount:   amount,
				Period:   period,
				Currency: currency,
			}
		}
	} else {
		// Interactive mode
		amount, period, currency, err := cli.Interactive()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		input = converter.Input{
			Amount:   amount,
			Period:   period,
			Currency: currency,
		}
	}

	api, err := exchangerate.NewExchangeRateAPI()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: failed to initialize exchange rate API: %v\n", err)
		os.Exit(1)
	}

	rates, rateInfo, err := api.GetRates(string(input.Currency))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: failed to fetch exchange rates: %v\n", err)
		os.Exit(1)
	}

	conv := converter.NewConverter(rates, string(input.Currency))

	results := conv.Convert(input)

	formatter := output.NewTableFormatter(input.Amount, input.Period, input.Currency, rateInfo)
	table := formatter.Format(results)
	fmt.Print(table)

	if flags.Verbose {
		verbose := output.FormatVerbose(rateInfo, rates)
		fmt.Print(verbose)
	}
}
