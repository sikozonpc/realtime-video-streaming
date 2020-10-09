package api

import (
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

	streamingTransport.NewWS(streaming.Initialize(s.DB), s.Router)

	s.Run()
}
