package main

import (
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/jieggii/groshi/internal/config"
	"github.com/jieggii/groshi/internal/currency/exchangerates"
	"github.com/jieggii/groshi/internal/database"
	"github.com/jieggii/groshi/internal/http_server/handlers"
	"github.com/jieggii/groshi/internal/http_server/middlewares"
	"github.com/jieggii/groshi/internal/http_server/validators"
	"github.com/jieggii/groshi/internal/loggers"
	"regexp"
)

func createHTTPRouter(jwtSecretKey string) *gin.Engine {
	router := gin.Default()

	// register validators:
	validatorEngine, ok := binding.Validator.Engine().(*validator.Validate)
	if !ok {
		loggers.Error.Fatalf("could not initialize validator engine")
	}

	// validatorsMap contains all validators and their tags
	validatorsMap := map[string]validator.Func{
		"username": validators.GetRegexValidator(regexp.MustCompile(".{2,}")),
		"password": validators.GetRegexValidator(regexp.MustCompile(".{8,}")),

		"description": validators.GetRegexValidator(regexp.MustCompile(".*")),
		"currency":    validators.GetCurrencyValidator(),
	}

	for validatorTag, validatorFunc := range validatorsMap {
		err := validatorEngine.RegisterValidation(validatorTag, validatorFunc)
		if err != nil {
			loggers.Error.Fatalf("could not register validator %v: %v", validatorTag, err)
		}
	}

	// allow all origins
	router.Use(cors.Default())

	// define and initialize middlewares:
	jwtHandlers := middlewares.NewJWTMiddleware(jwtSecretKey)
	jwtMiddleware := jwtHandlers.MiddlewareFunc()

	// register authorization & authentication routes:
	auth := router.Group("/auth")
	auth.POST("/login", jwtHandlers.LoginHandler)
	auth.POST("/logout", jwtHandlers.LogoutHandler)
	auth.POST("/refresh", jwtHandlers.RefreshHandler)

	// register user routes:
	user := router.Group("/user")
	user.POST("/", handlers.UserCreateHandler)                  // create new user
	user.GET("/", jwtMiddleware, handlers.UserReadHandler)      // read current user
	user.PUT("/", jwtMiddleware, handlers.UserUpdateHandler)    // update current user
	user.DELETE("/", jwtMiddleware, handlers.UserDeleteHandler) // delete current user

	// register transaction routes:
	transaction := router.Group("/transaction")
	transaction.Use(jwtMiddleware)
	transaction.POST("/", handlers.TransactionCreateHandler)        // create new transaction
	transaction.GET("/:uuid", handlers.TransactionReadOneHandler)   // read one transaction
	transaction.GET("/", handlers.TransactionReadManyHandler)       // read multiple transactions for given period
	transaction.GET("/summary", handlers.TransactionReadSummary)    // read summary about transactions for given period
	transaction.PUT("/:uuid", handlers.TransactionUpdateHandler)    // update transaction
	transaction.DELETE("/:uuid", handlers.TransactionDeleteHandler) // delete transaction

	return router
}

func main() {
	loggers.Info.Println("starting groshi")

	// read configuration from environmental variables:
	env := config.ReadEnvVars()

	// initialize database:
	if err := database.InitDatabase(
		env.MongoHost,
		env.MongoPort,
		config.ReadDockerSecret(env.MongoUsernameFiles),
		config.ReadDockerSecret(env.MongoPasswordFile),
		config.ReadDockerSecret(env.MongoDatabaseFile),
	); err != nil {
		loggers.Error.Fatal(err)
	}
	defer func() {
		if err := database.Client.Disconnect(database.Context); err != nil {
			loggers.Error.Fatalf("could not disconnect from the database: %v", err)
		}
	}()

	// initialize exchangerates API client:
	exchangerates.Client.Init(
		config.ReadDockerSecret(env.ExchangeRatesAPIKey),
	)
	router := createHTTPRouter(config.ReadDockerSecret(env.JWTSecretKeyFile))
	socket := fmt.Sprintf("%v:%v", env.Host, env.Port)

	loggers.Info.Printf("starting HTTP server on %v", socket)
	if err := router.Run(socket); err != nil {
		loggers.Error.Fatal(err)
	}
}
