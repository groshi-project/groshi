package middlewares

import (
	"fmt"
	"github.com/jieggii/groshi/internal/database"
	"github.com/jieggii/groshi/internal/passhash"
	"time"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

const jwtTimeout = time.Hour
const jwtMaxRefresh = time.Hour
const jwtIdentityKey = "username"

type _credentials struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type _claims struct {
	Username string
}

func NewAuthHandler(secretKey []byte) *jwt.GinJWTMiddleware {
	authMiddleware, err := jwt.New(&jwt.GinJWTMiddleware{
		Realm:       "test zone",
		Key:         secretKey,
		Timeout:     jwtTimeout,
		MaxRefresh:  jwtMaxRefresh,
		IdentityKey: jwtIdentityKey,
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			if v, ok := data.(*_claims); ok {
				return jwt.MapClaims{
					jwtIdentityKey: v.Username,
				}
			}
			return jwt.MapClaims{}
		},
		IdentityHandler: func(c *gin.Context) interface{} {
			claims := jwt.ExtractClaims(c)
			return &_claims{
				Username: claims[jwtIdentityKey].(string),
			}
		},
		Authenticator: func(c *gin.Context) (interface{}, error) {
			var credentials _credentials
			if err := c.ShouldBind(&credentials); err != nil {
				return nil, jwt.ErrMissingLoginValues
			}

			var user database.User
			if err := database.SelectUser(credentials.Username).Scan(database.Ctx, &user); err != nil {
				return nil, jwt.ErrFailedAuthentication // todo
			}

			if !passhash.ValidatePassword(credentials.Password, user.Password) {
				return nil, jwt.ErrFailedAuthentication
			}

			c.Set("current_user", &user)
			return &_claims{
				Username: user.Username,
			}, nil
		},
		Authorizator: func(data interface{}, c *gin.Context) bool {
			return true
			//if v, ok := data.(*_claims); ok && v.Username == "admin" {
			//	return true
			//}
			//
			//return false
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
	if err := authMiddleware.MiddlewareInit(); err != nil {
		panic(fmt.Errorf("could not initialize the jwt middleware: %v", err))
	}

	return authMiddleware
}
