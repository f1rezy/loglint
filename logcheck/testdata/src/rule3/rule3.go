package rule3

import (
	"log/slog"
)

func test() {
	slog.Info("server started 🚀") // want "english" "forbidden"
	slog.Info("what happened???") // want "excessive punctuation"
	slog.Info("server started")   // OK
}
