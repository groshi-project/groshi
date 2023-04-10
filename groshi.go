package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/jieggii/groshi/internal/config"
	"github.com/jieggii/groshi/internal/database"
	"github.com/jieggii/groshi/internal/http/handlers"
	"github.com/jieggii/groshi/internal/http/middlewares"
	"github.com/jieggii/groshi/internal/http/validators"
	"github.com/jieggii/groshi/internal/loggers"
)

func initHTTPRouter() *gin.Engine {
	router := gin.Default()

	// register validators
	validatorConfig, ok := binding.Validator.Engine().(*validator.Validate)
	if !ok {
		panic("could not get validator config")
	}

	validatorsMap := map[string]validator.Func{
		"username":                validators.Username,
		"password":                validators.Password,
		"transaction_description": validators.TransactionDescription,
		"currency":                validators.Currency,
	}

	for tag, fn := range validatorsMap {
		if err := validatorConfig.RegisterValidation(tag, fn); err != nil {
			panic(fmt.Errorf("could not register validator %v: %v", tag, err))
		}
	}

	// define and initialize middlewares
	authHandler := middlewares.NewAuthHandler([]byte("test 123"))
	authMiddleware := authHandler.MiddlewareFunc()

	// authorization & authentication
	auth := router.Group("/auth")
	auth.POST("/login", authHandler.LoginHandler)
	auth.POST("/logout", authHandler.LogoutHandler)
	auth.POST("/refresh", authHandler.RefreshHandler)

	// user
	user := router.Group("/user")
	//user.Use(authHandler.MiddlewareFunc())
	user.POST("/", handlers.UserCreate) // create new user
	user.GET("/", authMiddleware)       // read current user
	user.PUT("/", authMiddleware)       // update current user
	user.DELETE("/", authMiddleware)    // delete current user

	// transactions
	transactions := router.Group("/transactions")
	transactions.Use(authHandler.MiddlewareFunc())
	transactions.POST("/")        // create new transaction
	transactions.GET("/:uuid")    // get transaction
	transactions.GET("/")         // get transactions for given period
	transactions.PUT("/:uuid")    // update transaction
	transactions.DELETE("/:uuid") // delete transaction

	return router
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
	loggers.Info.Println("starting groshi")

	cfg := config.ReadFromEnv()
	if err := initDatabase(cfg); err != nil {
		loggers.Error.Fatal(err)
	}

	router := initHTTPRouter()
	address := fmt.Sprintf("%v:%v", cfg.Host, cfg.Port)

	loggers.Info.Printf("running HTTP server on %v", address)
	if err := router.Run(fmt.Sprintf(address)); err != nil {
		loggers.Error.Fatal(err)
	}
}
