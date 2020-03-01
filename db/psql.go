package db

import (
	"fmt"

	"github.com/Lefthander/otus-go-antibruteforce/config"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // The postgres driver for sqlx
)

// ConnectDB returns a connection pull to postgres
func ConnectDB(c *config.DBConfig) (*sqlx.DB, error) {
	connInfo := fmt.Sprintf("host=%s port=%s user=%s "+"password=%s dbname=%s sslmode=disable",
		c.DBHost, c.DBPort, c.DBUser, c.DBPass, c.DBName)

	return sqlx.Connect("postgres", connInfo)
}
