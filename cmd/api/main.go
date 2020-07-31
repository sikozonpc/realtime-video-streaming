package api

import (
	"goproject/auth"
	authTransport "goproject/auth/transport"
	"goproject/httpserver"
	"goproject/streaming"
	"goproject/streaming/hub"
	streamingTransport "goproject/streaming/transport"
	"log"
)

// Run the http server
func Run() {
	go hub.Instance.Run()

	s, err := httpserver.New()
	if err != nil {
		log.Fatal(err.Error())
	}

	defer s.DB.Close()

	authTransport.NewHTTP(auth.Initialize(s.DB), s.Router)
	streamingTransport.NewWS(streaming.Initialize(s.DB), s.Router)

	s.Run()
}
