package main

import (
	"errors"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/jwtauth/v5"
	"github.com/groshi-project/groshi/internal/database"
	"github.com/groshi-project/groshi/internal/service"
	"github.com/groshi-project/groshi/pkg/jwtauthority"
	"github.com/groshi-project/groshi/pkg/passwdauthority"
	"github.com/jessevdk/go-flags"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

const loggingBaseFlags = log.Ldate | log.Ltime | log.Lmsgprefix

var (
	infoLogger  = log.New(os.Stdout, "[info]: ", loggingBaseFlags)
	fatalLogger = log.New(os.Stderr, "[fatal]: ", loggingBaseFlags|log.Llongfile)
)

// Options provides application options which can be provided both using CLI and environmental variables.
type Options struct {
	General struct {
		Host string `long:"host" env:"GROSHI_HOST" default:"127.0.0.1" description:"host on which groshi will listen for client connections"`
		Port int    `short:"p" long:"port" env:"GROSHI_PORT" default:"8080" description:"port on which groshi will listen for client connections"`
	} `group:"General options"`

	Service struct {
		BcryptCost int `long:"bcrypt-cost" env:"GROSHI_BCRYPT_COST" description:"todo"`

		JWTSecretKey     string `long:"jwt-secret-key" env:"GROSHI_JWT_SECRET_KEY" description:"a secret key which will be used to generate JSON web tokens"`
		JWTSecretKeyFile string `long:"jwt-secret-key-file" env:"GROSHI_JWT_SECRET_KEY_FILE" description:"file containing a secret key which will be used to generate JSON web tokens"`

		JWTTimeToLive time.Duration `long:"jwt-ttl" env:"GROSHI_JWT_TTL" description:"jwt time-to-live" default:"744h"`
	} `group:"Service options"`

	Postgres struct {
		Host string `long:"postgres-host" env:"GROSHI_POSTGRES_HOST" required:"true" description:"host on which postgres is listening for groshi's connection'"`
		Port int    `long:"postgres-port" env:"GROSHI_POSTGRES_PORT" default:"5432" description:"host on which postgres is listening for groshi's connection"`

		User     string `long:"postgres-user" env:"GROSHI_POSTGRES_USER"  description:"todo"`
		UserFile string `long:"postgres-user-file" env:"GROSHI_POSTGRES_USER_FILE" description:"todo"`

		Password     string `long:"postgres-password" env:"GROSHI_POSTGRES_PASSWORD" description:"todo"`
		PasswordFile string `long:"postgres-password-file" env:"GROSHI_POSTGRES_PASSWORD_FILE" description:"todo"`

		Database     string `long:"postgres-database" env:"GROSHI_POSTGRES_DATABASE" description:"todo"`
		DatabaseFile string `long:"postgres-database-file" env:"GROSHI_POSTGRES_DATABASE_FILE" description:"todo"`
	} `group:"PostgreSQL options"`
}

// parseOptionsPair parses option pair. Option pair means option and its "file" pair.
// For example, `--postgres-password` and `--postgres-password-file`.
func parseOptionsPair(cliFlag string, envVar string, value *string, valueFile string) error {
	if *value == "" && valueFile == "" {
		return fmt.Errorf("`%s` ($%s) or `%s-file` ($%s_FILE) is required but not provided", cliFlag, envVar, cliFlag, envVar)
	}

	if *value == "" { // if value was not provided, read value from file:
		bytes, err := os.ReadFile(valueFile)
		if err != nil {
			return err
		}

		content := string(bytes)
		content = strings.Trim(content, "\n ")
		*value = content
	} else if valueFile != "" { // if both value and file are provided:
		return fmt.Errorf("both `%s` ($%s) and `%s-file` ($%s_FILE) are provided, expected only one of them", cliFlag, envVar, cliFlag, envVar)
	}

	return nil
}

// getOptions parses options from CLI and environmental variables.
// Prints error message and terminates program with code 1 on error.
func getOptions() *Options {
	parsingErrors := make([]error, 0)

	var options Options
	parser := flags.NewParser(&options, flags.Default)
	if _, err := parser.Parse(); err != nil {
		var flagsErr *flags.Error
		if errors.As(err, &flagsErr) && errors.Is(flagsErr.Type, flags.ErrHelp) {
			os.Exit(0)
		} else {
			os.Exit(1)
		}
	}

	if err := parseOptionsPair("--jwt-secret-key", "GROSHI_JWT_SECRET_KEY", &options.Service.JWTSecretKey, options.Service.JWTSecretKeyFile); err != nil {
		parsingErrors = append(parsingErrors, err)
	}

	if err := parseOptionsPair("--postgres-user", "GROSHI_POSTGRES_USER", &options.Postgres.User, options.Postgres.UserFile); err != nil {
		parsingErrors = append(parsingErrors, err)
	}

	if err := parseOptionsPair("--postgres-password", "GROSHI_POSTGRES_PASSWORD", &options.Postgres.Password, options.Postgres.PasswordFile); err != nil {
		parsingErrors = append(parsingErrors, err)
	}

	if err := parseOptionsPair("--postgres-database", "GROSHI_POSTGRES_DATABASE", &options.Postgres.Database, options.Postgres.DatabaseFile); err != nil {
		parsingErrors = append(parsingErrors, err)
	}

	if len(parsingErrors) != 0 {
		for _, parsingError := range parsingErrors {
			if _, err := fmt.Fprintln(os.Stderr, parsingError); err != nil {
				panic(err)
			}
		}
		os.Exit(1)
	}

	return &options
}

func getHTTPRouter(groshi *service.Service) *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(time.Duration(30) * time.Second))
	//r.Use(render.SetContentType(render.ContentTypeJSON))

	r.Route("/auth", func(r chi.Router) {
		// public `/auth` routes:
		r.Post("/login", groshi.AuthLogin)
	})

	r.Route("/user", func(r chi.Router) {
		// public `/user` route:
		r.Post("/", groshi.UserCreate)

		// protected `/user` routes:
		r.Group(func(r chi.Router) {
			r.Use(jwtauth.Verifier(groshi.JWTAuthority.JWTAuth))
			r.Use(jwtauth.Authenticator(groshi.JWTAuthority.JWTAuth))

			r.Get("/", groshi.UserGet)
			//r.Put("/", groshi.UserUpdate)
			r.Delete("/", groshi.UserDelete)
		})
	})

	r.Route("/categories", func(r chi.Router) {
		// protected `/categories` routes:
		r.Group(func(r chi.Router) {
			r.Use(jwtauth.Verifier(groshi.JWTAuthority.JWTAuth))
			r.Use(jwtauth.Authenticator(groshi.JWTAuthority.JWTAuth))

			r.Post("/", groshi.CategoriesCreate)
			r.Get("/", groshi.CategoriesGet)
			r.Put("/{uuid}", groshi.CategoriesUpdate)
			r.Delete("/{uuid}", groshi.CategoriesDelete)
		})
	})

	r.Route("/transactions", func(r chi.Router) {
		// protected `/transactions` routes:
		r.Group(func(r chi.Router) {
			r.Use(jwtauth.Verifier(groshi.JWTAuthority.JWTAuth))
			r.Use(jwtauth.Authenticator(groshi.JWTAuthority.JWTAuth))

			r.Post("/", groshi.TransactionsCreate)
			r.Get("/{uuid}", groshi.TransactionsGetOne)
			r.Get("/", groshi.TransactionsGet)
		})
	})

	r.Route("/stats", func(r chi.Router) {
		// protected `/stats` routes:
		r.Group(func(r chi.Router) {
			r.Use(jwtauth.Verifier(groshi.JWTAuthority.JWTAuth))
			r.Use(jwtauth.Authenticator(groshi.JWTAuthority.JWTAuth))

			r.Get("/total", groshi.StatsTotal)
		})
	})

	r.Get("/ping", groshi.Ping)

	return r
}

// startJob starts repeating job which will be run in a goroutine.
//func startJob(job func(args ...any), args []any, interval time.Duration) {
//	go func() {
//		job(args...)
//		for range time.Tick(interval) {
//			job(args...)
//		}
//	}()
//}

func main() {
	// get options provided using CLI and environmental variables:
	options := getOptions()

	infoLogger.Printf("starting groshi")

	// initialize postgres:
	db := database.New()
	if err := db.Connect(
		options.Postgres.Host,
		options.Postgres.Port,
		options.Postgres.User,
		options.Postgres.Password,
		options.Postgres.Database,
	); err != nil {
		fatalLogger.Fatalf("could not connect to the database: %s", err)
	}
	if err := db.Init(); err != nil {
		fatalLogger.Printf("could not initialize database: %s", err)
	}

	// create a groshi service:
	groshi := service.New(
		db,
		jwtauthority.New(options.Service.JWTTimeToLive, options.Service.JWTSecretKey),
		passwdauthority.New(options.Service.BcryptCost),
		log.New(os.Stderr, "[internal server error]: ", loggingBaseFlags|log.Llongfile),
	)

	// create an HTTP router:
	router := getHTTPRouter(groshi)

	// start listening:
	addr := fmt.Sprintf("%s:%d", options.General.Host, options.General.Port)
	infoLogger.Printf("groshi is listening for HTTP requests on %v", addr)
	if err := http.ListenAndServe(addr, router); err != nil {
		fatalLogger.Fatal(err)
	}
}
