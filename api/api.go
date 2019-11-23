package api

import (
	"fmt"
)

// Dictionary can look up a word
type Dictionary interface {
	Search(string)
}

// GetService gets dictionary service by name
func GetService(name string) (Dictionary, error) {
	switch name {
	case "iciba":
		return new(iciba), nil
	case "youdao":
		return new(youdao), nil
	default:
		return new(iciba), fmt.Errorf("Dictionary '%s' is not found, fall back to Iciba", name)
	}
}
