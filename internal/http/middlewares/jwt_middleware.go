package middlewares

import (
	"errors"
	"fmt"
	"github.com/jieggii/groshi/internal/database"
	"github.com/jieggii/groshi/internal/http/error_messages"
	"github.com/jieggii/groshi/internal/http/password_hashing"
	"github.com/jieggii/groshi/internal/loggers"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

const jwtTimeout = time.Hour
const jwtMaxRefresh = time.Hour

const jwtIdentityKey = "_id"

// jwtCredentials represents credentials that are used to generate JWT.
type jwtCredentials struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// jwtClaims represents information which will be included in JWT.
type jwtClaims struct {
	UserID string `json:"_id"`
}

// NewJWTMiddleware returns pointer to a new instance of jwt.GinJWTMiddleware.
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
					jwtIdentityKey: v.UserID,
				}
			}
			return jwt.MapClaims{}
		},

		IdentityHandler: func(c *gin.Context) interface{} {
			claims := jwt.ExtractClaims(c)
			return &jwtClaims{
				UserID: claims[jwtIdentityKey].(string),
			}
		},

		Authenticator: func(c *gin.Context) (interface{}, error) {
			credentials := jwtCredentials{}
			if err := c.ShouldBind(&credentials); err != nil {
				return nil, error_messages.ErrorInvalidRequestParams
			}

			user := database.User{}
			err := database.Users.FindOne(
				database.Context, bson.D{{"username", credentials.Username}},
			).Decode(&user)
			if err != nil {
				if errors.Is(err, mongo.ErrNoDocuments) {
					return nil, errors.New("user with such username was not found")
				}
				loggers.Error.Printf("unexpected error: %v", err)
				return nil, errors.New("internal server error")
			}

			if ok := password_hashing.ValidatePassword(credentials.Password, user.Password); !ok {
				return nil, errors.New("incorrect password")
			}
			return &jwtClaims{
				UserID: user.ID.Hex(),
			}, nil

		},

		Authorizator: func(data interface{}, c *gin.Context) bool {
			if v, ok := data.(*jwtClaims); ok {
				user := database.User{}

				userID, err := primitive.ObjectIDFromHex(v.UserID)
				if err != nil {
					loggers.Error.Printf("unexpected error when parsing user ObjectID: %v", err)
					return false
				}

				if err := database.Users.FindOne(
					database.Context, bson.D{{"_id", userID}},
				).Decode(&user); err != nil {
					if !errors.Is(err, mongo.ErrNoDocuments) {
						loggers.Error.Printf("unexpected error: %v", err)
					}
					return false
				}
				c.Set("current_user", &user)
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

	if err := jwtMiddleware.MiddlewareInit(); err != nil {
		panic(fmt.Errorf("could not initialize the jwt middleware: %v", err))
	}

	return jwtMiddleware
}
