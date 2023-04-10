package main

import (
	"fmt"
	"github.com/jieggii/groshi/internal/http/handlers"
	"github.com/jieggii/groshi/internal/http/middlewares"

	"github.com/gin-gonic/gin"
	"github.com/jieggii/groshi/internal/config"
	"github.com/jieggii/groshi/internal/database"
	"github.com/jieggii/groshi/internal/loggers"
)

func startHTTPServer(host string, port int) error {
	r := gin.Default()

	// middlewares
	authHandler := middlewares.NewAuthHandler([]byte("test 123"))
	authMiddleware := authHandler.MiddlewareFunc()

	// authorization & authentication
	auth := r.Group("/auth")
	auth.POST("/login", authHandler.LoginHandler)
	auth.POST("/logout", authHandler.LogoutHandler)
	auth.POST("/refresh", authHandler.RefreshHandler)

	// user
	user := r.Group("/user")
	//user.Use(authHandler.MiddlewareFunc())
	user.POST("/", handlers.UserCreate) // create new user
	user.GET("/", authMiddleware)       // read current user
	user.PUT("/", authMiddleware)       // update current user
	user.DELETE("/", authMiddleware)    // delete current user

	// transactions
	transactions := r.Group("/transactions")
	transactions.Use(authHandler.MiddlewareFunc())
	transactions.POST("/")        // create new transaction
	transactions.GET("/:uuid")    // get transaction
	transactions.GET("/")         // get transactions for given period
	transactions.PUT("/:uuid")    // update transaction
	transactions.DELETE("/:uuid") // delete transaction

	loggers.Info.Printf("starting HTTP server on %v:%v.\n", host, port)

	return r.Run(fmt.Sprintf("%v:%v", host, port))
}

func initDatabase(cfg *config.Config) error {
	if err := database.Connect(
		cfg.PostgresHost,
		cfg.PostgresPort,
		cfg.PostgresUser,
		cfg.PostgresPassword,
		cfg.PostgresDatabase,
	); err != nil {
		return fmt.Errorf(
			"could not connect to PostgreSQL database \"%v\" at %v:%v (%v)",
			cfg.PostgresDatabase,
			cfg.PostgresHost,
			cfg.PostgresPort,
			err,
		)
	}

	if err := database.Init(); err != nil {
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
	loggers.Info.Println("starting groshi server")

	cfg := config.ReadFromEnv()
	if err := initDatabase(cfg); err != nil {
		loggers.Error.Fatal(err)
	}

	if err := startHTTPServer(cfg.Host, cfg.Port); err != nil {
		loggers.Error.Fatal(err)
	}
}
