package storage_test

import (
	"os"
	"testing"
)

var (
	databaseURL string
)

func TestMain(t *testing.M) {
	databaseURL = os.Getenv("DATABASE_URL")

	if databaseURL == "" {
		databaseURL = "postgres://selectel:selectel@127.0.0.1:5432/selectel?sslmode=disable"
	}

	os.Exit(t.Run())
}
