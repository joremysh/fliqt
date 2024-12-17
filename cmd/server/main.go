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
	"github.com/joremysh/fliqt/internal/repository"
	"github.com/joremysh/fliqt/pkg/database"
)

func NewServer(petStore *api.HRSystem, port string) *http.Server {
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
	if dsn == "" {
		dsn = "user:password@tcp(localhost:3306)/hrs?collation=utf8_unicode_ci&parseTime=true&loc=Asia%2FTaipei&multiStatements=true"
	}

	gdb, err := database.NewDatabase(dsn)
	if err != nil {
		log.Fatal(err.Error())
	}
	err = repository.Migrate(gdb)
	if err != nil {
		log.Fatal(err.Error())
	}

	// Create an instance of our handler which satisfies the generated interface
	hrSystem := api.NewHRSystem(gdb)
	s := NewServer(hrSystem, port)
	// And we serve HTTP until the world ends.
	log.Fatal(s.ListenAndServe())
}
