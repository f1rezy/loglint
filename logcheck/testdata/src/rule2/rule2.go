package rule2

import (
	"log/slog"
)

func test() {
	slog.Info("сервер запущен") // want "lowercase" "english"
	slog.Info("server started") // OK
}
