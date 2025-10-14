package config

import "fmt"

type Config struct {
	Env      string   `yaml:"env" env-default:"dev"`
	Host     string   `yaml:"host"`
	Port     int      `yaml:"port"`
	Database Database `yaml:"database"`
}

type Database struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Name     string `yaml:"name"`
	SSLMode  string `yaml:"ssl_mode"`
}

func (c *Database) GetDatabaseURL() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.User, c.Password, c.Name, c.SSLMode)
}
