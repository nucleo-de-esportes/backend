package config

type dbConfig struct {
	host string
	user string
	password string
}

type serverConfig struct {
	port string
}

type Config struct {
	db dbConfig
	server serverConfig
}
