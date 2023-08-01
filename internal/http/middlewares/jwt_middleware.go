package middlewares

import (
	"errors"
	"fmt"
	"github.com/jieggii/groshi/internal/database"
	"github.com/jieggii/groshi/internal/http/error_messages"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
	"time"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

const jwtTimeout = time.Hour
const jwtMaxRefresh = time.Hour
const jwtIdentityKey = "user_uuid"

type jwtCredentials struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type jwtClaims struct {
	UserUUID string `json:"user_uuid"`
}

func NewJWTMiddleware(secretKey string) *jwt.GinJWTMiddleware {
	jwtMiddleware, err := jwt.New(&jwt.GinJWTMiddleware{
		Realm:       "groshi",
		Key:         []byte(secretKey),
		Timeout:     jwtTimeout,
		MaxRefresh:  jwtMaxRefresh,
		IdentityKey: jwtIdentityKey,
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			if v, ok := data.(*jwtClaims); ok {
				return jwt.MapClaims{
					jwtIdentityKey: v.UserUUID,
				}
			}
			return jwt.MapClaims{}
		},

		IdentityHandler: func(c *gin.Context) interface{} {
			claims := jwt.ExtractClaims(c)
			return &jwtClaims{
				UserUUID: claims[jwtIdentityKey].(string),
			}
		},

		Authenticator: func(c *gin.Context) (interface{}, error) {
			credentials := jwtCredentials{}
			if err := c.ShouldBind(&credentials); err != nil {
				return nil, errors.Join(errors.New(error_messages.ErrorInvalidRequestParams), jwt.ErrMissingLoginValues)
			}

			user := database.User{}
			err := database.Users.FindOne(
				database.Context, bson.D{{"username", credentials.Username}},
			).Decode(&user)
			if err != nil {
				if errors.Is(err, mongo.ErrNoDocuments) {
					return nil, errors.Join(errors.New("incorrect username or password"), jwt.ErrFailedAuthentication)
				}
				return nil, errors.Join(errors.New(error_messages.ErrorInternalServerError), jwt.ErrFailedAuthentication)
			}

			err = bcrypt.CompareHashAndPassword(
				[]byte(user.Password), []byte(credentials.Password),
			)
			if err != nil {
				return nil, errors.Join(errors.New("incorrect username or password"), jwt.ErrFailedAuthentication)
			}

			return &jwtClaims{
				UserUUID: user.UUID,
			}, nil
		},

		Authorizator: func(data interface{}, c *gin.Context) bool {
			if v, ok := data.(*jwtClaims); ok {
				user := database.User{}
				err := database.Users.FindOne(
					database.Context, bson.D{{"uuid", v.UserUUID}},
				).Decode(&user)
				if err != nil {
					if errors.Is(err, mongo.ErrNoDocuments) {
						return false
					}
					// todo: log unexpected error
					return false
				}
				c.Set("currentUser", &user)
				return true
			} else {
				return false
			}
		},

		Unauthorized: func(c *gin.Context, code int, message string) {
			c.JSON(code, gin.H{
				"error_message": message,
			})
		},
		TokenLookup:   "header:Authorization",
		TokenHeadName: "Bearer",
		TimeFunc:      time.Now,
	})
	if err != nil {
		panic(fmt.Errorf("could not create jwt middleware instance: %v", err))
	}

	// When you use jwt.New(), the function is already automatically called for checking,
	// which means you don't need to call it again.
	if err := jwtMiddleware.MiddlewareInit(); err != nil {
		panic(fmt.Errorf("could not initialize the jwt middleware: %v", err))
	}

	return jwtMiddleware
}
