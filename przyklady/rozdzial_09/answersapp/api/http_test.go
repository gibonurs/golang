package api

import (
	"net/http"
	"testing"
)

func TestPathParams(t *testing.T) {
	r, err := http.NewRequest("GET", "1/2/3/4/5", nil)
	if err != nil {
		t.Errorf("NewRequest: %s", err)
	}
	params := pathParams(r, "one/two/three/four")
	if len(params) != 4 {
    t.Errorf("oczekiwano 4 parametrów, a uzyskano %d: %v", len(params), params)
	}
	for k, v := range map[string]string{
    "jeden":  "1",
    "dwa":    "2",
    "trzy":   "3",
    "cztery": "4",
	} {
		if params[k] != v {
			t.Errorf("%s: %s != %s", k, params[k], v)
		}
	}
  params = pathParams(r, "jeden/dwa/trzy/cztery/pięć/sześć")
  if len(params) != 5 {
    t.Errorf("oczekiwano 5 parametrów, a uzyskano %d: %v", len(params), params)
  }
  for k, v := range map[string]string{
    "jeden":  "1",
    "dwa":    "2",
    "trzy":   "3",
    "cztery": "4",
    "pięć":   "5",
	} {
		if params[k] != v {
			t.Errorf("%s: %s != %s", k, params[k], v)
		}
	}
}
