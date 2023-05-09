package api

import "testing"

func TestInitRequest(t *testing.T) {
	api := newDefaultAPIHandler()
	if err := api.Init(); err != nil {
		t.Fatal(err)
	}
}
