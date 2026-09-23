package core

import (
	"encoding/json"
	"testing"
)

func TestCanonicalPayerJSON(t *testing.T) {
	value := Payer{
		Name:           "Juan Perez",
		ExternalID:     "payer-1",
		DocumentType:   "CUIT",
		DocumentNumber: "20234567897",
		AccountIdentifier: &AccountIdentifier{
			Type: "cbu", Value: "0070000000000000000000",
		},
		Institution: &Institution{Type: "bank", Code: "007", Name: "Banco de ejemplo"},
	}
	raw, err := json.Marshal(value)
	if err != nil {
		t.Fatal(err)
	}
	var got map[string]any
	if err = json.Unmarshal(raw, &got); err != nil {
		t.Fatal(err)
	}
	if got["name"] != "Juan Perez" || got["documentNumber"] != "20234567897" {
		t.Fatalf("payer=%s", raw)
	}
}

func TestCanonicalPayerOmitsUnavailableFields(t *testing.T) {
	raw, err := json.Marshal(Payer{Institution: &Institution{Type: "wallet", Name: "MP"}})
	if err != nil {
		t.Fatal(err)
	}
	if string(raw) != `{"institution":{"type":"wallet","name":"MP"}}` {
		t.Fatalf("payer=%s", raw)
	}
}
