package xero

import (
	"encoding/json"
	"os"
	"testing"
)

func TestAccountsType(t *testing.T) {

	b, err := os.ReadFile("testdata/accounts.json")
	if err != nil {
		t.Fatal(err)
	}

	var ar AccountResponse
	if err := json.Unmarshal(b, &ar); err != nil {
		t.Fatal(err)
	}
	if got, want := len(ar.Accounts), 90; got != want {
		t.Errorf("got %d accounts, want %d", got, want)
	}
}
