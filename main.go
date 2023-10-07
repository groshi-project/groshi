package main

import (
	"fmt"
	"net/http"
	"regexp"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/groshi-project/groshi/docs"
	"github.com/groshi-project/groshi/internal/config"
	"github.com/groshi-project/groshi/internal/currency/exchangerates"
	"github.com/groshi-project/groshi/internal/database"
	"github.com/groshi-project/groshi/internal/handlers"
	"github.com/groshi-project/groshi/internal/loggers"
	"github.com/groshi-project/groshi/internal/middlewares"
	"github.com/groshi-project/groshi/internal/validators"
	swaggofiles "github.com/swaggo/files"
	swagger "github.com/swaggo/gin-swagger"
)

//	@title			groshi HTTP API documentation
//	@version		0.1.0
//	@description	📉 groshi - goddamn simple tool to keep track of your finances.
//	@license.name	Licensed under MIT license.
//	@license.url	https://github.com/groshi-project/groshi/tree/master/LICENSE

func createHTTPRouter(jwtSecretKey string, enableDebug bool, enableSwagger bool) *gin.Engine {
	// set appropriate gin mode
	if enableDebug {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

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

		"description":       validators.GetRegexValidator(regexp.MustCompile(".*")),
		"currency":          validators.GetCurrencyValidator(false),
		"optional_currency": validators.GetCurrencyValidator(true),

		"nonzero_time": validators.NonzeroTimeValidator,
	}

	for validatorTag, validatorFunc := range validatorsMap {
		err := validatorEngine.RegisterValidation(validatorTag, validatorFunc)
		if err != nil {
			loggers.Error.Fatalf("could not register validator %v: %v", validatorTag, err)
		}
	}

	// setup cross-origin resource sharing
	corsConfig := cors.Config{
		AllowAllOrigins: true,
		AllowMethods:    []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodOptions},
		AllowHeaders: []string{
			"Authorization",
			"Content-Type",
		},
	}
	router.Use(cors.New(corsConfig))

	// define and initialize middlewares:
	jwtHandlers := middlewares.NewJWTMiddleware(jwtSecretKey)
	jwtMiddleware := jwtHandlers.MiddlewareFunc()

	// register authorization & authentication routes:
	auth := router.Group("/auth")
	auth.POST("/login", jwtHandlers.LoginHandler)
	auth.POST("/refresh", jwtHandlers.RefreshHandler)

	// register user routes:
	user := router.Group("/user")
	user.POST("", handlers.UserCreateHandler)                  // create new user
	user.GET("", jwtMiddleware, handlers.UserReadHandler)      // read current user
	user.PUT("", jwtMiddleware, handlers.UserUpdateHandler)    // update current user
	user.DELETE("", jwtMiddleware, handlers.UserDeleteHandler) // delete current user

	// register transactions routes:
	transactions := router.Group("/transactions")
	transactions.Use(jwtMiddleware)
	transactions.POST("", handlers.TransactionsCreateHandler)         // create new transaction
	transactions.GET("", handlers.TransactionsReadManyHandler)        // read multiple transactions for given period
	transactions.GET("/:uuid", handlers.TransactionsReadOneHandler)   // read one transaction
	transactions.PUT("/:uuid", handlers.TransactionsUpdateHandler)    // update transaction
	transactions.DELETE("/:uuid", handlers.TransactionsDeleteHandler) // delete transaction
	transactions.GET("/summary", handlers.TransactionsReadSummary)    // read summary about transactions for given period

	// register currencies route:
	currencies := router.Group("/currencies")
	currencies.GET("", handlers.CurrenciesRead) // read available currencies

	// register Swagger documentation route if needed:
	if enableSwagger {
		docs.SwaggerInfo.BasePath = ""
		router.GET("/docs/*any", swagger.WrapHandler(swaggofiles.Handler))
	} else {
		loggers.Warning.Println("swagger UI feature is disabled, no documentation will be shown at `/docs/index.html` route")
	}

	return router
}

func main() {
	loggers.Info.Printf("starting groshi")

	// read configuration from environmental variables:
	env := config.ReadEnvVars()

	// initialize database:
	if err := database.InitDatabase(
		env.MongoHost,
		env.MongoPort,
		config.ReadDockerSecret(env.MongoUsernameFile),
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

	router := createHTTPRouter(config.ReadDockerSecret(env.JWTSecretKeyFile), env.Debug, env.Swagger)
	socket := fmt.Sprintf("%v:%v", env.Host, env.Port)

	loggers.Info.Printf("starting HTTP server on %v", socket)
	if err := router.Run(socket); err != nil {
		loggers.Error.Fatal(err)
	}
}
