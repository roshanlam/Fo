package main

import (
	"fmt"
	"log"
	"mime"
	"net/http"
	"os"

	"github.com/jmoiron/sqlx"
	"github.com/gin-gonic/gin"

	_ "github.com/lib/pq"
)

func main() {
	databaseUrl := os.Getenv("DATABASE_URL")
	if databaseUrl == "" {
		databaseUrl = "postgres://postgres:postgres@localhost/passwordless_social?sslmode=disable"
	}

	db, err := sqlx.Open("postgres", databaseUrl)
	if err != nil {
		log.Fatal(err)
	}

	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected!")

	// Set up our Gin router
	r := gin.Default()

	// Load templates from the templates directory
	r.LoadHTMLGlob("templates/*")
	// Make sure we serve .js files as JavaScript
	mime.AddExtensionType(".js", "application/javascript")
	// Serve static files from the static directory
	r.Static("/static", "./static")

	// Set up two groups, which we’ll be able to apply middleware to later on for easier handling
	unauthenticated := r.Group("/")
	authenticated := r.Group("/")

	// Add a default web page, which will become our homepage
	authenticated.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Hello world!")
	})


	unauthenticated.GET("/login", LoginHandler())

	// Run our web server on port 8080
	log.Fatal(r.Run(":8080"))
}
