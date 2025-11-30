# Salary Calculator (s-calc)

A Go console application that converts salary amounts between different currencies (PLN, EUR, USD, GBP) and time periods (Hour, Day, Month, Year).

## Features

- **Currency Conversion**: Supports PLN, EUR, USD, GBP
- **Time Period Conversion**: Hour, Day, Month, Year
- **Multiple Input Methods**: Command-line flags or interactive prompts
- **Exchange Rate Caching**: Caches rates for 24 hours to reduce API calls
- **Beautiful Table Output**: Formatted table with highlighted original input
- **Rate Metadata**: Shows rate source, timestamp, and cache expiration

## Installation

```bash
go build -o s-calc ./cmd/s-calc
```

Or install globally:

```bash
go install ./cmd/s-calc
```

## Usage

### Command-Line Flags

```bash
# Hourly salary
s-calc -h=20 -c=EUR

# Monthly salary
s-calc -m=5000 -c=USD

# Yearly salary
s-calc -y=60000 -c=PLN

# With verbose output
s-calc -h=25 -c=GBP -v
```

### Interactive Mode

If no flags are provided, the application will prompt for input:

```bash
s-calc
```

Example interaction:
```
Select period [Hour/Day/Month/Year]: Hour
Enter amount: 20
Select currency [PLN/EUR/USD/GBP]: EUR
```


Also supports the format:

```bash
s-calc -h=20 EUR
```

## Output Example

```
┌─────────┬─────────────┬─────────────┬─────────────┬─────────────┐
│ Period  │ PLN         │ EUR         │ USD         │ GBP         │
├─────────┼─────────────┼─────────────┼─────────────┼─────────────┤
│ Hour    │ 85.00       │ 20.00 ⭐    │ 22.00       │ 17.00       │
│ Day     │ 680.00      │ 160.00      │ 176.00      │ 136.00      │
│ Month   │ 14,743.40   │ 3,466.80    │ 3,813.48    │ 2,945.78    │
│ Year    │ 176,920.80  │ 41,601.60   │ 45,761.76   │ 35,349.36   │
└─────────┴─────────────┴─────────────┴─────────────┴─────────────┘

⭐ Original input: 20.00 EUR/hour

Rate source: exchangerate-api.com
Last updated: 2024-01-15 10:30:00 UTC
Cache expires: 2024-01-16 10:30:00 UTC
```

## Configuration

### Environment Variables

- `S_CALC_CACHE_TTL`: Cache TTL in hours (default: 24)
- `S_CALC_CACHE_DIR`: Custom cache directory path
- `S_HOURS_DAY`: Working hours per day (default: 8)
- `S_DAYS_MONTH`: Working days per month (default: 21.67)

### Cache Location

- **Unix/Linux/macOS**: `~/.cache/s-calc/rates-{currency}.json`
- **Windows**: `%LOCALAPPDATA%\s-calc\rates-{currency}.json`

## Exchange Rate Sources

The application uses the following APIs (in order of preference):

1. **exchangerate-api.com** (primary) - Free tier: 1,500 requests/month
2. **exchangerate.host** (fallback) - Free, no API key required

Rates are cached for 24 hours to minimize API calls.

## Conversion Logic

### Time Periods

- **Hour → Day**: Multiply by working hours per day (default: 8)
- **Day → Month**: Multiply by working days per month (default: 21)
- **Month → Year**: Multiply by 12

### Currency Conversion

Exchange rates are fetched from external APIs and cached locally. The application automatically handles currency conversions using the latest available rates.

## Error Handling

The application handles various error scenarios:

- Network errors: Falls back to cached rates if available
- Invalid API responses: Tries fallback API or uses expired cache
- Invalid input: Shows clear error messages with examples
- Cache errors: Continues without cache, attempts to create cache directory

## Development

### Project Structure

```
salary-calc/
├── cmd/
│   └── s-calc/
│       └── main.go          # Entry point
├── internal/
│   ├── converter/
│   │   └── converter.go      # Conversion logic
│   ├── exchangerate/
│   │   ├── api.go            # API client
│   │   └── cache.go          # Caching logic
│   ├── cli/
│   │   ├── flags.go          # Flag parsing
│   │   └── interactive.go    # Interactive prompts
│   └── output/
│       └── table.go          # Table formatting
├── go.mod
└── README.md
```

### Building

```bash
go build -o s-calc ./cmd/s-calc
```

## License

MIT

