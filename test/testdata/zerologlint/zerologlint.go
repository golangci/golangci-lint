//golangcitest:args -Ezerologlint
package zerologlint

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func bad() {
	log.Error() // want "must be dispatched by Msg or Send method"
	log.Info()  // want "must be dispatched by Msg or Send method"
	log.Fatal() // want "must be dispatched by Msg or Send method"
	log.Debug() // want "must be dispatched by Msg or Send method"
	log.Warn()  // want "must be dispatched by Msg or Send method"

	var err error
	log.Error().Err(err)                                 // want "must be dispatched by Msg or Send method"
	log.Error().Err(err).Str("foo", "bar").Int("foo", 1) // want "must be dispatched by Msg or Send method"

	logger := log.Error() // want "must be dispatched by Msg or Send method"
	logger.Err(err).Str("foo", "bar").Int("foo", 1)

	// include zerolog.Dict()
	log.Info(). // want "must be dispatched by Msg or Send method"
			Str("foo", "bar").
			Dict("dict", zerolog.Dict().
				Str("bar", "baz").
				Int("n", 1),
		)

	// conditional
	logger2 := log.Info() // want "must be dispatched by Msg or Send method"
	if err != nil {
		logger2 = log.Error() // want "must be dispatched by Msg or Send method"
	}
	logger2.Str("foo", "bar")
}

func ok() {
	log.Fatal().Send()
	log.Panic().Msg("")
	log.Debug().Send()
	log.Info().Msg("")
	log.Warn().Send()
	log.Error().Msg("")

	log.Error().Str("foo", "bar").Send()
	var err error
	log.Error().Err(err).Str("foo", "bar").Int("foo", 1).Msg("")

	logger := log.Error()
	logger.Send()

	// include zerolog.Dict()
	log.Info().
		Str("foo", "bar").
		Dict("dict", zerolog.Dict().
			Str("bar", "baz").
			Int("n", 1),
		).Send()

	// conditional
	logger2 := log.Info()
	if err != nil {
		logger2 = log.Error()
	}
	logger2.Send()
}
