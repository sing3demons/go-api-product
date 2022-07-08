package middleware

import (
	"log"
	"os"
	"time"

	"github.com/sing3demons/app/database"
	"github.com/sing3demons/app/models"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type login struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

var identityKey = "sub"
var exp = time.Hour * 72

// Authenticate - publicLoginHandler
// Login godoc
// @Summary login
// @Description login by form user
// @Tags auth
// @Accept  json
// @Produce  json
// @Param login body login true "login"
// @Success 200 string token
// @Router /api/v1/auth/login [post]
func Authenticate() *jwt.GinJWTMiddleware {
	authMiddleware, err := jwt.New(&jwt.GinJWTMiddleware{
		// secret key
		Key:           []byte(os.Getenv("SECRET_KEY")),
		Timeout:       exp,
		MaxRefresh:    exp,
		IdentityKey:   identityKey,
		TokenLookup:   "header: Authorization",
		TokenHeadName: "Bearer",

		IdentityHandler: func(c *gin.Context) interface{} {
			var user models.User
			claims := jwt.ExtractClaims(c)
			id := claims[identityKey]

			db := database.GetDB()
			if err := db.First(&user, uint(id.(float64))).Error; err != nil {
				log.Printf("error: %v", err)
				return nil
			}

			return &user
		},

		// login => user
		Authenticator: func(c *gin.Context) (interface{}, error) {
			var form login
			var user models.User

			if err := c.ShouldBindJSON(&form); err != nil {
				return nil, jwt.ErrMissingLoginValues
			}

			db := database.GetDB()
			if err := db.Where("email = ?", form.Email).First(&user).Error; err != nil {
				log.Printf("error: %v", err)
				return nil, jwt.ErrFailedAuthentication
			}

			if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(form.Password)); err != nil {
				return nil, jwt.ErrFailedAuthentication
			}

			return &user, nil
		},

		// user => payload(sub) => JWT
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			if v, ok := data.(*models.User); ok {
				claims := jwt.MapClaims{
					identityKey: v.ID,
				}

				return claims
			}

			return jwt.MapClaims{}
		},

		// handle error
		Unauthorized: func(c *gin.Context, code int, message string) {
			c.JSON(code, gin.H{
				"error": message,
			})
		},
	})

	if err != nil {
		log.Fatal("JWT Error:" + err.Error())
	}

	return authMiddleware
}
