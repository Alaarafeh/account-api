package login

import (
	"context"
	"log"
	"net/http"
	"register-api/infrastructure/db"
	"register-api/models"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
)

var jwtKey = []byte("secret_key")

type Handler struct {
	dbConnection *db.Connection
	validator    *validator.Validate
}

func NewHandler(dbConnection *db.Connection, validate *validator.Validate) *Handler {
	return &Handler{dbConnection: dbConnection, validator: validate}
}

func (h *Handler) GetCredential() gin.HandlerFunc {
	return func(c *gin.Context) {

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		var user models.User
		c.BindJSON(&user)
		defer cancel()


		// Check if the email and password match a user in the MongoDB database
		collection := h.dbConnection.GetUsers()
		var result models.User

		err := collection.FindOne(ctx, bson.M{"email": user.Email}).Decode(&result)

		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": "Invalid email or password",
			})
			return
		}

		err = bcrypt.CompareHashAndPassword([]byte(result.Password), []byte(user.Password))
		
		log.Println(result.Password)

		if err !=  nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": "Invalid email or password",
			})
			return
		}

		 // create JWT token
		 expirationTime := time.Now().Add(24 * time.Hour)
		 claims := &jwt.StandardClaims{
			 ExpiresAt: expirationTime.Unix(),
			 Subject:   user.Email,
		 }
		 token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		 tokenString, err := token.SignedString(jwtKey)
		 log.Println(tokenString)
		 if err != nil {
			c.JSON(http.StatusInternalServerError, "Failed to sign token")
			 return
		 }

		 c.SetCookie("token", tokenString, int(time.Hour/time.Second), "/", "localhost", false, true)
		 c.JSON(http.StatusOK, "Cookie set")
	}
}