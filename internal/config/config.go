package config

type Config struct {
	Postgres	PostgresConfig	`yaml:"postgres"`
}

type PostgresConfig struct {
	Host	  string		`yaml:"host"`
	Port	  string		`yaml:"port"`
	User	  string		`yaml:"user"`
	Password  string		`yaml:"password"`
}