package model

import (
	"time"
)

type PartsFilter struct {
	Uuids                 []string
	Names                 []string
	Category              []Category
	ManufacturerCountries []string
	Tags                  []string
}

type Part struct {
	// ID            primitive.ObjectID   `bson:"_id,omitempty"`
	Uuid          string               `bson:"uuid"`
	Name          string               `bson:"name"`
	Description   string               `bson:"description"`
	Price         float64              `bson:"price"`
	StockQuantity int64                `bson:"stock_quantity"`
	Category      string               `bson:"category"`
	Dimensions    Dimensions           `bson:"dimensions"`
	Manufacturer  Manufacturer         `bson:"manufacturer"`
	Tags          []string             `bson:"tags"`
	Metadata      map[string]MetaValue `bson:"metadata"`
	CreatedAt     time.Time            `bson:"created_at"`
	UpdatedAt     time.Time            `bson:"updated_at"`
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
	Length float64 `bson:"length"`
	Width  float64 `bson:"width"`
	Height float64 `bson:"height"`
	Weight float64 `bson:"weight"`
}

type Manufacturer struct {
	Name    string `bson:"name"`
	Country string `bson:"country"`
	Website string `bson:"website"`
}

type MetaValue struct {
	StringValue *string  `bson:"string_value,omitempty"`
	Int64Value  *int64   `bson:"int64_value,omitempty"`
	DoubleValue *float64 `bson:"double_value,omitempty"`
	BoolValue   *bool    `bson:"bool_value,omitempty"`
}
