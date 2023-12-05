package app

import (
	"os"

	"github.com/streamingfast/logging"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"golang.org/x/term"
)

var log *zap.SugaredLogger

func setupLogger() (err error) {
	switch cfg.Logger {
	case "dev":
		isTTY := term.IsTerminal(int(os.Stderr.Fd()))
		verbosity := 1
		logStdoutWriter := zapcore.Lock(os.Stderr)
		l := zap.New(zapcore.NewCore(logging.NewEncoder(verbosity, isTTY), logStdoutWriter, zapcore.DebugLevel))
		log = l.Sugar().Named("app")
	case "release":
		zapcfg := zap.NewProductionConfig()
		zapcfg.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
		l, err := zap.NewProduction()
		if err != nil {
			return err
		}
		log = l.Sugar().Named("app")
	}

	return nil
}
