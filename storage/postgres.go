package storage

import (
	"fmt"
	"database/sql"
	"log"
	
	_ "github.com/lib/pq"
)

type Config struct {
	Host 		string
	Password 	string
	User 		string
	DBName 		string
	SSLMode 	string
}

func NewConnection(config *Config)(*sql.DB, error){
	connStr := fmt.Sprintf(
		"user=%s password=%s dbname=%s host=%s sslmode=%s",
		config.User, config.Password, config.DBName, config.Host, config.SSLMode,
	)

	db, err := sql.Open("postgres", connStr)
	if err != nil{
		log.Fatal(err)
		return db, err
	}
	
	return db, nil
}