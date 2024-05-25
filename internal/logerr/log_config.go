package logerr

import (
	"os"
	"time"

	"github.com/rs/zerolog"
)

// Логирует информативные сообщения, вид: 'time' INF 'msg'.
func InfoMsg(msg string) {
	var log = zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.DateTime}).
		Level(zerolog.TraceLevel).
		With().
		Timestamp().
		Logger()
	log.Info().Msg(msg)
}

// Логирует события ошибок, вид:
// 'time' ERR 'filepath':'line' > 'msg' error='err'.
func ErrEvent(msg string, err error) {
	var log = zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.DateTime}).
		Level(zerolog.TraceLevel).
		With().
		Timestamp().
		Caller().
		Logger()
	log.Error().Err(err).Msg(msg)
}

// Логирует события фатальных ошибок, вид:
// 'time' FTL 'filepath':'line' > 'msg' error='err'.
func FatalEvent(msg string, err error) {
	var log = zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.DateTime}).
		Level(zerolog.TraceLevel).
		With().
		Timestamp().
		Caller().
		Logger()
	log.Fatal().Err(err).Msg(msg)
}
