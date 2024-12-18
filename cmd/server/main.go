package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"

	middleware "github.com/oapi-codegen/gin-middleware"

	"github.com/joremysh/fliqt/api"
	"github.com/joremysh/fliqt/internal/handler"
	"github.com/joremysh/fliqt/internal/repository"
	"github.com/joremysh/fliqt/pkg/cache"
	"github.com/joremysh/fliqt/pkg/database"
)

func NewServer(petStore *handler.HRSystem, port string) *http.Server {
	swagger, err := api.GetSwagger()

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading swagger spec\n: %s", err)
		os.Exit(1)
	}

	// Clear out the servers array in the swagger spec, that skips validating
	// that server names match. We don't know how this thing will be run.
	swagger.Servers = nil

	// This is how you set up a basic gin router
	r := gin.Default()

	// Use our validation middleware to check all requests against the
	// OpenAPI schema.
	r.Use(middleware.OapiRequestValidator(swagger))

	// We now register our petStore above as the handler for the interface
	api.RegisterHandlers(r, petStore)

	s := &http.Server{
		Handler: r,
		Addr:    net.JoinHostPort("0.0.0.0", port),
	}
	return s
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	dsn := os.Getenv("DSN")
	redisHost := os.Getenv("REDIS_HOST")
	redisPort := os.Getenv("REDIS_PORT")

	gdb, err := database.NewDatabase(dsn)
	if err != nil {
		log.Fatal(err.Error())
	}
	err = repository.Migrate(gdb)
	if err != nil {
		log.Fatal(err.Error())
	}

	redisClient, err := cache.NewRedisClient(redisHost + ":" + redisPort)
	if err != nil {
		log.Fatal(err.Error())
	}

	hrSystem := handler.NewHRSystem(gdb, redisClient)
	s := NewServer(hrSystem, port)

	log.Fatal(s.ListenAndServe())
}
