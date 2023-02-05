package main

import (
	"log"
	"register-api/infrastructure/db"
	"register-api/ui/rest/login"
	"register-api/ui/rest/register"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func main() {
    r := gin.Default()

    // read db config from environment
	dbConfig, err := db.NewConfigFromEnvironment()
	if err != nil {
		log.Fatal(err)
	}

	// setup db connection
	dbConnection, err := db.NewConnection(dbConfig)
	if err != nil {
		log.Fatal(err)
	}

	// setup validator
	validate := validator.New()

	loginHandler := login.NewHandler(dbConnection, validate)
	login.RegisterRoutes(r, loginHandler)

	registerHandler := register.NewHandler(dbConnection, validate)
	register.RegisterRoutes(r, registerHandler)

    r.Run()
}
