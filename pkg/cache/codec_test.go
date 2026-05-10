package cache

import (
	"reflect"
	"testing"
)

func TestJSONCodec_RoundTrip(t *testing.T) {
	var c JSONCodec
	type sample struct {
		Name string `json:"name"`
		N    int    `json:"n"`
	}
	in := sample{Name: "x", N: 42}
	data, err := c.Marshal(in)
	if err != nil {
		t.Fatal(err)
	}
	var out sample
	if err := c.Unmarshal(data, &out); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(in, out) {
		t.Errorf("got %+v", out)
	}
}
