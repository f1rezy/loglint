package rule1

import (
	"log/slog"
)

func test() {
	slog.Info("Starting server") // want "lowercase"
	slog.Info("server started")  // OK
}
