package main

import (
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jieggii/groshi/internal/config"
	"github.com/jieggii/groshi/internal/database"
	"github.com/jieggii/groshi/internal/http/handlers"
	"github.com/jieggii/groshi/internal/http/middlewares"
	"github.com/jieggii/groshi/internal/loggers"
)

func createHTTPRouter(jwtSecretKey string) *gin.Engine {
	router := gin.Default()

	// allow all origins
	router.Use(cors.Default())

	// define and initialize middlewares:
	jwtHandlers := middlewares.NewJWTMiddleware(jwtSecretKey)
	jwtMiddleware := jwtHandlers.MiddlewareFunc()

	// authorization & authentication routes:
	auth := router.Group("/auth")
	auth.POST("/login", jwtHandlers.LoginHandler)
	auth.POST("/logout", jwtHandlers.LogoutHandler)
	auth.POST("/refresh", jwtHandlers.RefreshHandler)

	// user routes:
	user := router.Group("/user")
	user.POST("/", handlers.UserCreateHandler)                  // create new user
	user.GET("/", jwtMiddleware, handlers.UserReadHandler)      // read current user
	user.PUT("/", jwtMiddleware, handlers.UserUpdateHandler)    // update current user
	user.DELETE("/", jwtMiddleware, handlers.UserDeleteHandler) // delete current user

	// transaction routes:
	transaction := router.Group("/transaction")
	transaction.POST("/", jwtMiddleware)        // create new transaction
	transaction.GET("/:uuid", jwtMiddleware)    // read transaction
	transaction.GET("/", jwtMiddleware)         // read transactions for given period
	transaction.PUT("/:uuid", jwtMiddleware)    // update transaction
	transaction.DELETE("/:uuid", jwtMiddleware) // delete transaction

	return router
}

func main() {
	loggers.Info.Println("starting groshi")

	env := config.ReadEnvVars()
	if err := database.InitDatabase(
		env.MongoHost,
		env.MongoPort,
		config.ReadDockerSecret(env.MongoUsernameFiles),
		config.ReadDockerSecret(env.MongoPasswordFile),
		config.ReadDockerSecret(env.MongoDatabaseFile),
	); err != nil {
		loggers.Error.Fatal(err)
	}

	router := createHTTPRouter(config.ReadDockerSecret(env.JWTSecretKeyFile))
	socket := fmt.Sprintf("%v:%v", env.Host, env.Port)

	loggers.Info.Printf("starting HTTP server on %v", socket)
	if err := router.Run(socket); err != nil {
		loggers.Error.Fatal(err)
	}
}