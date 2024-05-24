package config

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

func ConnectToPostgreSQL() (*sql.DB, error) {
	fmt.Println("Connecting to PostgreSQL...")
	config, err := LoadConfig("config/conf.json")
	if err != nil {
		return nil, err
	}

	fmt.Println("Configuration loaded successfully!")

	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		config.Database.Host, config.Database.Port, config.Database.User, config.Database.Password, config.Database.DBName, config.Database.SSLMode)

	// Open connection to database
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	fmt.Println("Connection to PostgreSQL opened successfully!")

	// Ping database to verify connection
	err = db.Ping()

	fmt.Println(err)
	if err != nil {
		return nil, err
	}

	log.Println("Connected to PostgreSQL!")

	return db, nil
}

func PrintDBConfig() {
	config, err := LoadConfig("config/conf.json")
	if err != nil {
		fmt.Println("Error loading config:", err)
		return
	}

	fmt.Println("Configuration:", config)
	fmt.Println("Host:", config.Database.Host)
	fmt.Println("Port:", config.Database.Port)
	fmt.Println("User:", config.Database.User)
	fmt.Println("Database Host:", config.Database.Host)
	fmt.Println("Database Password:", config.Database.Password)
}
