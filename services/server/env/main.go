package env

import (
	"flag"
	"os"
)

// Variables for the env file
type Variables struct {
	Port    string
	Address string
}

// ParseEnv parses the environment variables to run the API
func ParseEnv() Variables {
	var (
		port = flag.String("port", os.Getenv("PORT"), "The http server port")
		addr = flag.String("addr", os.Getenv("ADDR"), "The http server address")
	)

	flag.Parse()

	return Variables{*port, *addr}
}
