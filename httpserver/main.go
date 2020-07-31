package httpserver

import (
	"database/sql"
	"fmt"
	"github.com/gorilla/mux"
	"goproject/env"
	"goproject/postgresdb"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

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
	// Open connection to database...
	db, err := postgresdb.New(envVars)
	if err != nil {
		return nil, err
	}

	r := mux.NewRouter()
	s := &Server{
		Router: r,
		DB:     db,
	}

	// Run migrations...

	s.Router.HandleFunc("/health", handleHealthCheck)

	fs := http.FileServer(http.Dir("static/"))
	s.Router.Handle("/static/", http.StripPrefix("/static/", fs))

	return s, nil
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.Router.ServeHTTP(w, r)
}

// Run the server instance
func (s *Server) Run() {
	stop := make(chan os.Signal, 1)

	signal.Notify(stop, os.Interrupt)

	serverAddr := fmt.Sprintf("%s:%s", envVars.Address, envVars.Port)

	h := &http.Server{
		Addr:         serverAddr,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 5 * time.Second,
		Handler:       cors.Default().Handler(s),
	}

	go func() {
		log.Printf("Listening on %s\n", serverAddr)

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
