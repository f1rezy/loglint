package rule4

import (
	"log/slog"
)

func test(password string, token string) {
	slog.Info("user password: " + password) // want "sensitive"
	slog.Info("token: " + token)            // want "sensitive"

	slog.Info("token validated")    // OK
	slog.Info("user authenticated") // OK
}
