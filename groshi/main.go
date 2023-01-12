package main

import (
	"fmt"
	"github.com/jieggii/groshi/groshi/auth/jwt"
	"github.com/jieggii/groshi/groshi/config"
	"github.com/jieggii/groshi/groshi/database"
	"github.com/jieggii/groshi/groshi/handles"
	"github.com/jieggii/groshi/groshi/loggers"
	"net/http"
)

func startHTTPServer(host string, port int) error {
	mux := http.NewServeMux()

	// user handles:
	mux.HandleFunc("/user/auth", middleware(false, handles.UserAuth))
	mux.HandleFunc("/user/create", middleware(true, handles.UserCreate))
	mux.HandleFunc("/user/read", middleware(false, handles.UserRead))
	mux.HandleFunc("/user/delete", middleware(true, handles.UserDelete))

	// transaction handles:
	mux.HandleFunc("/transaction/create", middleware(true, handles.TransactionCreate))
	mux.HandleFunc("/transaction/read", middleware(true, handles.TransactionRead))
	mux.HandleFunc("/transaction/update", middleware(true, handles.TransactionUpdate))
	mux.HandleFunc("/transaction/delete", middleware(true, handles.TransactionDelete))

	loggers.Info.Printf("Starting HTTP server on %v:%v.\n", host, port)
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
