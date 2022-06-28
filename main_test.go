package main

import (
	"testing"
	"time"

	"github.com/allegro/bigcache/v3"
)

func init() {
	testing_mode = true
	nameStartsWith = "sp"
	limit = 10
	cache, _ = bigcache.NewBigCache(bigcache.DefaultConfig(10 * time.Minute))
}

func TestGetMarvelData(t *testing.T) {
	got := getMarvelData()
	want := 0

	if got.Offset != want {
		t.Errorf("got %q, wanted %q", got.Offset, want)
	}
}
