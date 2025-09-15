package main_test

import (
	"os"
	main "sms"
	"testing"
)

func TestInitDBWithDSN(t *testing.T) {
	tests := []struct {
		name          string
		dsn           string
		shouldSucceed bool
	}{
		{
			name:          "Successful connection to in-memory db",
			dsn:           ":memory:",
			shouldSucceed: true,
		},
		{
			name:          "Connection to a valid file db",
			dsn:           "test_db.sqlite",
			shouldSucceed: true,
		},
		{
			name:          "Invalid DSN causes ping error",
			dsn:           "file:/path/to/nonexistent/db.sqlite?_auth&_auth_user=test",
			shouldSucceed: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, err := main.InitDBWithDSN(tt.dsn)

			if err != nil {
				if tt.shouldSucceed {
					t.Fatalf("expected no error from sql.Open, but got: %v", err)
				}
				return
			}
			if db == nil {
				t.Fatal("expected a non-nil database connection")
			}
			defer db.Close()
			pingErr := db.Ping()

			if tt.shouldSucceed {
				if pingErr != nil {
					t.Errorf("expected successful ping, but got error: %v", pingErr)
				}
			} else {
				if pingErr == nil {
					t.Fatal("expected an error during ping, but got nil")
				}
			}
			if tt.dsn == "test_db.sqlite" {
				os.Remove(tt.dsn)
			}
		})
	}
}
