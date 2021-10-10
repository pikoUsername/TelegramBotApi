package tgp

import (
	"net/url"
	"testing"
)

// go test -v

func TestUrlvaluesToMapString(t *testing.T) {
	val := map[string]string{}
	value_text := "exe"

	v := url.Values{}

	v.Add("kek", value_text)

	// urlValuesToMapString(v, val)

	if v, ok := val["kek"]; !ok || v != value_text {
		t.Fatal("value is not correct, converting value is not working")
	}
}
