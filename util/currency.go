package util

// Constants for all supported currencies
const (
	USD = "USD"
	EUR = "EUR"
	GBP = "GBP"
)

// IsSupportedCurrency returns true if the currency is supported
func IsSupportedCurrency(currency string) bool {
	switch currency {
	case USD, EUR, GBP:
		return true
	}
	return false
}
