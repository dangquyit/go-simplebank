package util

// Constants for all supported currencies

const (
	USD = "USD"
	VND = "VND"
	EUR = "EUR"
)

// IsSupportedCurrency return true if the currency is supported
func IsSupportedCurrency(currency string) bool {
	switch currency {
	case USD, VND, EUR:
		return true
	}
	return false
}
