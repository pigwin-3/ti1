package config

import (
	"database/sql"
	"fmt"
	"log"
	"time"

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

    // Set connection pool settings
    db.SetMaxOpenConns(25) // Maximum number of open connections to the database
    db.SetMaxIdleConns(25) // Maximum number of connections in the idle connection pool
    db.SetConnMaxLifetime(1 * time.Hour) // Maximum amount of time a connection may be reused

	fmt.Println("Connection to PostgreSQL opened successfully :D")

    // Ping database to verify connection
    err = db.Ping()

    fmt.Println(err)
    if err != nil {
        return nil, err
    }

    log.Println("Connected to PostgreSQL!")

    return db, nil
}

func DisconnectFromPostgreSQL(db *sql.DB) error {
	fmt.Println("Disconnecting from PostgreSQL...")
	if err := db.Close(); err != nil {
		return err
	}
	log.Println("Disconnected from PostgreSQL!")
	return nil
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
	fmt.Println("Database User:", config.Database.User)
	fmt.Println("Database Password:", config.Database.Password)
}
