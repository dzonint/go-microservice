package data

import "testing"

func TestChecksValidation(t *testing.T) {
	p := &Product{
		Name:  "Dzonint",
		Price: 1.99,
		SKU:   "abc-abcd-abcde",
	}

	err := p.Validate()

	if err != nil {
		t.Fatal(err)
	}
}
