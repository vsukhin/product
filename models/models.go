package models

const (
	// DefaultCurrency contains default currency
	DefaultCurrency = "USD"
)

// Product contains product content
type Product struct {
	ID          int                `json:"ID"`
	Name        string             `json:"Name" valid:"Required"`
	Description string             `json:"Description" valid:"Required"`
	Price       float64            `json:"Price" valid:"Required"`
	Tags        []string           `json:"Tags,omitempty"`
	Prices      map[string]float64 `json:"Prices,omitempty"`
}

// ID contains id
type ID struct {
	ID int `json:"ID"`
}

// ValidationError contains validation error
type ValidationError struct {
	Field   string `json:"Field"`
	Message string `json:"Message"`
}
