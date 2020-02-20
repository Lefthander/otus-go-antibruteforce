package db

import (
	"fmt"

	"github.com/Lefthander/otus-go-antibruteforce/config"
	"github.com/jmoiron/sqlx"
)

// ConnectDB returns a connection pull to postgres
func ConnectDB(c *config.ServiceConfig) (*sqlx.DB, error) {

	connInfo := fmt.Sprintf("host=%s port=%s user=%s "+"password=%s dbname=%s sslmode=disable",
		c.DbHost, c.DBPort, c.DBUser, c.DBPass, c.DBName)

	return sqlx.Connect("postgres", connInfo)

}
