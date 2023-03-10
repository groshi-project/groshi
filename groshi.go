package main

import (
	"fmt"
	"github.com/jieggii/groshi/internal/config"
	"github.com/jieggii/groshi/internal/database"
	"github.com/jieggii/groshi/internal/http/handles"
	"github.com/jieggii/groshi/internal/http/jwt"
	"github.com/jieggii/groshi/internal/http/middleware"
	"github.com/jieggii/groshi/internal/loggers"
	"net/http"
)

func startHTTPServer(host string, port int) error {
	mux := http.NewServeMux()

	// user handles:
	mux.HandleFunc(
		"/user/auth", middleware.Middleware(false, handles.UserAuth),
	)
	mux.HandleFunc(
		"/user/create", middleware.Middleware(true, handles.UserCreate),
	)
	mux.HandleFunc(
		"/user/read", middleware.Middleware(true, handles.UserRead),
	)
	mux.HandleFunc(
		"/user/update", middleware.Middleware(true, handles.UserUpdate),
	)
	mux.HandleFunc(
		"/user/delete", middleware.Middleware(true, handles.UserDelete),
	)

	// transaction handles:
	mux.HandleFunc(
		"/transaction/create", middleware.Middleware(true, handles.TransactionCreate),
	)
	mux.HandleFunc(
		"/transaction/read", middleware.Middleware(true, handles.TransactionRead),
	)
	mux.HandleFunc(
		"/transaction/update", middleware.Middleware(true, handles.TransactionUpdate),
	)
	mux.HandleFunc(
		"/transaction/delete", middleware.Middleware(true, handles.TransactionDelete),
	)

	loggers.Info.Printf("starting HTTP server on %v:%v.\n", host, port)
	err := http.ListenAndServe(fmt.Sprintf("%v:%v", host, port), mux)

	return err
}

func initializeApp(cfg *config.Config) error {
	jwt.SecretKey = cfg.JWTSecretKey

	if err := database.Connect(
		cfg.PostgresHost,
		cfg.PostgresPort,
		cfg.PostgresUser,
		cfg.PostgresPassword,
		cfg.PostgresDatabase,
	); err != nil {
		return fmt.Errorf(
			"failed to connect to PostgreSQL database \"%v\" at %v:%v (%v)",
			cfg.PostgresDatabase,
			cfg.PostgresHost,
			cfg.PostgresPort,
			err,
		)
	}
	if err := database.Initialize(); err != nil {
		return fmt.Errorf(
			"failed to initialize PostgreSQL database \"%v\" at %v:%v (%v)",
			cfg.PostgresDatabase,
			cfg.PostgresHost,
			cfg.PostgresPort,
			err,
		)
	}
	return nil
}

func main() {
	loggers.Info.Println("Starting groshi server...")
	cfg := config.ReadFromEnv()
	if err := initializeApp(cfg); err != nil {
		loggers.Fatal.Fatal(err)
	}
	if err := startHTTPServer(cfg.Host, cfg.Port); err != nil {
		loggers.Fatal.Fatal(err)
	}
}
