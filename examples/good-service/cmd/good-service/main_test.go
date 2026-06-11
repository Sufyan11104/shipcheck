package main

import "testing"

func TestGoodServiceFixture(t *testing.T) {
	t.Setenv("APP_PORT", "8080")
}
