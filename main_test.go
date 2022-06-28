package main

import "testing"

func init() {
	testing_mode = true
	nameStartsWith = "sp"
	limit = 10
}

func TestGetMarvelData(t *testing.T) {
	got := getMarvelData()
	want := 0

	if got.Offset != want {
		t.Errorf("got %q, wanted %q", got.Offset, want)
	}
}
