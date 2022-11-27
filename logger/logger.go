package logger

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/pkgerrors"
)

var (
	logFileHandle *os.File
	Logger        zerolog.Logger
)

func init() {
	if strings.ToLower(os.Getenv("DEBUG")) == "true" {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	} else {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}

	zerolog.TimeFieldFormat = time.RFC3339
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
	Logger = zerolog.New(os.Stdout).
		With().
		Timestamp().
		Caller().
		Logger()
	// 	DisableColors: true,
	// 	FullTimestamp: true,
	// 	PadLevelText:  false,
	// 	ForceQuote:    false,
	// })
	Logger.Info().Msgf("Starting logger on level %s", Logger.GetLevel())
}

func CloseLogFile() {
	if ^logFileHandle.Fd() == 0 {
		return
	}

	if err := logFileHandle.Close(); err != nil {
		fmt.Println("Couldn't close logging file handle")
		Logger.Error().Err(err).Msg("Couldn't close logging file handle")
	}
}
