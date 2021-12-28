package env

import (
	"flag"
	"os"
)

// Variables for the env file
type Variables struct {
	Port       string
	Address    string
	DbName     string
	DbUser     string
	DbHost     string
	DbPort     string
	DbPassword string
}

// ParseEnv parses the enviroment variables to run the API
func ParseEnv() Variables {

	var (
		port       = flag.String("port", os.Getenv("PORT"), "The http server port")
		addr       = flag.String("addr", os.Getenv("ADDR"), "The http server address")
		dbname     = flag.String("db-name", os.Getenv("DB-NAME"), "The database name")
		dbuser     = flag.String("db-user", os.Getenv("DB-USER"), "The database user name")
		dbhost     = flag.String("db-host", os.Getenv("DB_HOST"), "The database host")
		dbport     = flag.String("db-port", os.Getenv("DB_PORT"), "The database port")
		dbpassword = flag.String("db-password", os.Getenv("DB-PASSWORD"), "The database  password")
	)

	flag.Parse()

	// TODO: Read from .env file
	// var err error

	// err = godotenv.Load()
	// if err != nil {
	// 	log.Fatalf("error getting env, not comming through %v", err)
	// } else {
	// 	fmt.Println("env variables successfuly loaded")
	// }

	return Variables{*port, *addr, *dbname, *dbuser, *dbhost, *dbport, *dbpassword}
}
