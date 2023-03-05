package register

import (
	"context"
	"net/http"
	"register-api/infrastructure/db"
	"register-api/models"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
)

type Handler struct {
	dbConnection *db.Connection
	validator    *validator.Validate
}

func NewHandler(dbConnection *db.Connection, validate *validator.Validate) *Handler {
	return &Handler{dbConnection: dbConnection, validator: validate}
}

func (h *Handler) CreateUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		var user models.User
		c.BindJSON(&user)
		defer cancel()

		// check if user is null
		if user.Email == "" || user.Password == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Error User or Password cannot be empty",
			})
			return
		} 

		// Hash the password
		password := []byte(user.Password)
		hashedPassword, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": "Error hashing password",
			})
			return
		}

		// Save the user to the MongoDB databaseh.dbConnection.GetUsers()
		collection := h.dbConnection.GetUsers()
		_, err = collection.InsertOne(ctx, bson.M{"email": user.Email, "password": hashedPassword})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Error creating user",
				"error":   err.Error(),
			})
			return
		}
		c.Header("Access-Control-Allow-Origin","*")
		c.JSON(http.StatusOK, gin.H{
			"message": "User created successfully!",
		})
	}
}
