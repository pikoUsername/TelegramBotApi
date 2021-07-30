package utils_test

import (
	"net/url"
	"testing"

	"github.com/pikoUsername/tgp/utils"
)

func TestURLToMapString(t *testing.T) {
	v := &url.Values{}
	base_value := "key"
	v.Add(base_value, base_value)
	params := make(map[string]string)
	utils.UrlValuesToMapString(*v, params)
	b, ok := params[base_value]
	if !ok {
		t.Error("cannot get string by key")
		t.Fail()
	}

	if b != base_value {
		t.Error("values from params and real value is different")
		t.Fail()
	}
	t.Log(params)
}
