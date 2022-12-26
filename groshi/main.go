package main

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/jieggii/groshi/groshi/config"
	"github.com/jieggii/groshi/groshi/database"
	"github.com/jieggii/groshi/groshi/handles"
	"github.com/jieggii/groshi/groshi/handles/jwt"
	"github.com/jieggii/groshi/groshi/loggers"
	"net/http"
)

func startHTTPServer(host string, port int) error {
	// todo: use `r.Use` for middlewares
	r := chi.NewRouter()

	// service handles:
	r.Post("/auth", handles.Auth)

	// users handles:
	r.Route("/user", func(r chi.Router) {
		r.Post("/create", jwt.ValidateJWTMiddleware(handles.UserCreate))
		r.Post("/{username}", jwt.ValidateJWTMiddleware(handles.UserRead))
		r.Post("/{username}/update", jwt.ValidateJWTMiddleware(handles.UserUpdate))
		r.Post("/{username}/delete", jwt.ValidateJWTMiddleware(handles.UserDelete))
	})

	// transaction handles:
	//router.Handle("POST", "/transaction/create", handles.TransactionCreate)
	//router.Handle("GET", "/transaction/:uuid", handles.TransactionRead)
	//router.Handle("PUT", "/transaction/:uuid/update", handles.TransactionUpdate)
	//router.Handle("DELETE", "/transaction/:uuid/delete", handles.TransactionDelete)

	loggers.Info.Printf("Starting HTTP server on %v:%v.\n", host, port)
	err := http.ListenAndServe(fmt.Sprintf("%v:%v", host, port), r)

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
			"failed to connect to PostgreSQL database \"%v\" at %v:%v: %v",
			cfg.PostgresDatabase,
			cfg.PostgresHost,
			cfg.PostgresPort,
			err,
		)
	}
	if err := database.Initialize(cfg.SuperuserUsername, cfg.SuperuserPassword); err != nil {
		return fmt.Errorf(
			"failed to initialize PostgreSQL database \"%v\" at %v:%v: %v",
			cfg.PostgresDatabase,
			cfg.PostgresHost,
			cfg.PostgresPort,
			err,
		)
	}
	return nil
}

func main() {
	loggers.Info.Println("Starting groshi server.")
	cfg := config.ReadFromEnv()
	if err := initializeApp(cfg); err != nil {
		loggers.Fatal.Fatalf("Error initializing groshi: %v.\n", err)
	}
	if err := startHTTPServer(cfg.Host, cfg.Port); err != nil {
		loggers.Fatal.Fatalf("Error starting HTTP server: %v.", err)
	}
}
