package model

import "time"

type PartsFilter struct {
	Uuids                 []string
	Names                 []string
	Category              []Category
	ManufacturerCountries []string
	Tags                  []string
}

type Part struct {
	Uuid          string
	Name          string
	Description   string
	Price         float64
	StockQuantity int64
	Category      Category
	Dimensions    Dimensions
	Manufacturer  Manufacturer
	Tags          []string
	Metadata      map[string]MetaValue
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

type Category string

const (
	CategoryUnknown  Category = "CATEGORY_UNKNOWN"
	CategoryEngine   Category = "CATEGORY_ENGINE"
	CategoryFuel     Category = "CATEGORY_FUEL"
	CategoryPorthole Category = "CATEGORY_PORTHOLE"
	CategoryWing     Category = "CATEGORY_WING"
)

type Dimensions struct {
	Length float64
	Width  float64
	Height float64
	Weight float64
}

type Manufacturer struct {
	Name    string
	Country string
	Website string
}

type MetaValue struct {
	StringValue *string
	Int64Value  *int64
	DoubleValue *float64
	BoolValue   *bool
}
