package main

// Product represents a product
type Product struct {
	SKU  string `json:"sku" validate:"required"`
	Name string `json:"name" validate:"required"`
}
