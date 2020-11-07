package httpserver

import (
	"database/sql"
	"fmt"
	"goproject/env"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"

	"github.com/rs/cors"
)

var envVars = env.ParseEnv()

// Server struct
type Server struct {
	DB     *sql.DB
	Router *mux.Router
	// Logger *provider.Logger
}

// New creates a new http server instance
func New() (*Server, error) {
	r := mux.NewRouter()
	s := &Server{
		Router: r,
	}

	s.Router.HandleFunc("/health", handleHealthCheck)

	return s, nil
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.Router.ServeHTTP(w, r)
}

// Run the server instance
func (s *Server) Run() {
	stop := make(chan os.Signal, 1)

	signal.Notify(stop, os.Interrupt)

	//serverAddr := fmt.Sprintf("%s:%s", envVars.Address, envVars.Port)

	h := &http.Server{
		Addr:         ":8080",
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 5 * time.Second,
		Handler:      cors.Default().Handler(s),
	}

	go func() {
		log.Printf("Listening on %s\n", ":8080")

		if err := h.ListenAndServe(); err != nil {
			log.Fatal(err)
		}
	}()

	<-stop

	log.Println("\nGracefully shutting down the server...")
}

func handleHealthCheck(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "All good to Go :)")
}
