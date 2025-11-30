package cli

import (
	"flag"
	"fmt"
	"os"
	"strconv"
)

type Flags struct {
	Hour     *float64
	Day      *float64
	Month    *float64
	Year     *float64
	Currency string
	Verbose  bool
}

func ParseFlags() (*Flags, []string) {
	flags := &Flags{}

	flags.Hour = flag.Float64("h", 0, "Salary per hour")
	flags.Day = flag.Float64("d", 0, "Salary per day")
	flags.Month = flag.Float64("m", 0, "Salary per month")
	flags.Year = flag.Float64("y", 0, "Salary per year")
	flag.StringVar(&flags.Currency, "c", "EUR", "Currency (PLN, EUR, USD, GBP)")
	flag.BoolVar(&flags.Verbose, "v", false, "Show detailed rate information")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [flags]\n\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Flags:\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nExamples:\n")
		fmt.Fprintf(os.Stderr, "  %s -h=20 EUR\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s -m=5000 -c=USD\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s  (interactive mode)\n", os.Args[0])
	}

	flag.Parse()

	args := flag.Args()
	if len(args) > 0 {
		// Try to parse as currency if not set via flag
		if flags.Currency == "EUR" {
			flags.Currency = args[0]
		}
	}

	return flags, args
}

func (f *Flags) HasInput() bool {
	return (f.Hour != nil && *f.Hour > 0) ||
		(f.Day != nil && *f.Day > 0) ||
		(f.Month != nil && *f.Month > 0) ||
		(f.Year != nil && *f.Year > 0)
}

func (f *Flags) GetInput() (amount float64, period string, currency string, ok bool) {
	if !f.HasInput() {
		return 0, "", "", false
	}

	switch {
	case f.Hour != nil && *f.Hour > 0:
		return *f.Hour, "Hour", f.Currency, true
	case f.Day != nil && *f.Day > 0:
		return *f.Day, "Day", f.Currency, true
	case f.Month != nil && *f.Month > 0:
		return *f.Month, "Month", f.Currency, true
	case f.Year != nil && *f.Year > 0:
		return *f.Year, "Year", f.Currency, true
	}

	return 0, "", "", false
}

func ParseLegacyFormat(args []string) (amount float64, period string, currency string, ok bool) {
	if len(args) < 1 {
		return 0, "", "", false
	}

	for _, arg := range args {
		if len(arg) > 3 && arg[0] == '-' {
			periodChar := arg[1]
			if arg[2] == '=' {
				valStr := arg[3:]
				val, err := strconv.ParseFloat(valStr, 64)
				if err == nil && val > 0 {
					switch periodChar {
					case 'h', 'H':
						period = "Hour"
					case 'd', 'D':
						period = "Day"
					case 'm', 'M':
						period = "Month"
					case 'y', 'Y':
						period = "Year"
					default:
						continue
					}
					amount = val
					ok = true
					break
				}
			}
		}
	}

	for _, arg := range args {
		if len(arg) == 3 {
			currency = arg
			break
		}
	}

	if currency == "" {
		currency = "EUR"
	}

	return amount, period, currency, ok
}
