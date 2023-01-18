package main

import (
	"fmt"

	"github.com/rheisen/bconf"
	uuid "github.com/satori/go.uuid"
)

func main() {
	configuration := bconf.NewAppConfig(
		"external_http_api",
		"HTTP API for user authentication and authorization",
	)

	_ = configuration.SetLoaders(
		&bconf.EnvironmentLoader{KeyPrefix: "ext_http_api"},
		&bconf.FlagLoader{},
	)

	_ = configuration.AddFieldSets(
		bconf.NewFieldSetBuilder().Key("app").Fields(
			bconf.NewFieldBuilder().
				Key("id").Type(bconf.String).
				Description("Application identifier").
				DefaultGenerator(
					func() (any, error) {
						return fmt.Sprintf("%s", uuid.NewV4().String()), nil
					},
				).Create(),
			bconf.FB(). // FB() is a shorthand function for NewFieldBuilder()
					Key("session_secret").Type(bconf.String).
					Description("Application secret for session management").
					Sensitive().Required().
					Validator(
					func(fieldValue any) error {
						secret, _ := fieldValue.(string)

						minLength := 20
						if len(secret) < minLength {
							return fmt.Errorf(
								"expected string of minimum %d characters (len=%d)",
								minLength,
								len(secret),
							)
						}

						return nil
					},
				).Create(),
		).Create(),
		bconf.FSB().Key("log").Fields( // FSB() is a shorthand function for NewFieldSetBuilder()
			bconf.FB().
				Key("level").Type(bconf.String).Default("info").
				Description("Logging level").
				Enumeration("debug", "info", "warn", "error").Create(),
			bconf.FB().
				Key("format").Type(bconf.String).Default("json").
				Description("Logging format").
				Enumeration("console", "json").Create(),
			bconf.FB().
				Key("color_enabled").Type(bconf.Bool).Default(true).
				Description("Colored logs when format is 'console'").
				Create(),
		).Create(),
	)

	// Register with the option to handle --help / -h flag set to true
	if errs := configuration.Register(true); len(errs) > 0 {
		// handle configuration load errors
	}

	// returns the log level found in order of: default -> environment -> flag -> user override
	// (based on the loaders set above).
	logLevel, err := configuration.GetString("log", "level")
	if err != nil {
		// handle retrieval error
	}

	fmt.Printf("log-level: %s", logLevel)
}
