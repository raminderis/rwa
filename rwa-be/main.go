package main

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5"
)

type PostgresConfig struct {
	Host     string
	Port     string
	Username string
	Password string
	Database string
	SSLMode  string
}

func (cfg PostgresConfig) DatabaseURL() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.Database, cfg.SSLMode)
}

func main() {
	cfg := PostgresConfig{
		Host:     "localhost",
		Port:     "5432",
		Username: "baloo",
		Password: "junglebook",
		Database: "rwabe",
		SSLMode:  "disable",
	}
	// fmt.Println("Postgres Config:", cfg)
	// export DATABASE_URL="postgres://baloo:junglebook@localhost:5432/rwabe"
	conn, err := pgx.Connect(context.Background(), cfg.DatabaseURL())
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(context.Background())

	err = conn.Ping(context.Background())
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to ping database: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Connected to database successfully!")
	_, err = conn.Exec(context.Background(), `
		CREATE TABLE IF NOT EXISTS dces (
			stockname TEXT, 
			stockid TEXT,
			stocktinker TEXT);
		CREATE TABLE IF NOT EXISTS dcec (
			stockid TEXT,
			smartcontractaddress TEXT,
			supportedchainids TEXT);
	`)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to create table: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Table created successfully!")

	// Insert a new record into the dces table
	stockname := "Apple Inc."
	stockid := "1"
	stocktinker := "AAPL"
	_, err = conn.Exec(context.Background(), "INSERT INTO dces (stockname, stockid, stocktinker) VALUES ($1, $2, $3)", stockname, stockid, stocktinker)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to insert record: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Record inserted successfully!")

	smartcontractaddress := "0xD771a71E5bb303da787b4ba2ce559e39dc6eD85c"
	supportedchainids := "11155111"
	_, err = conn.Exec(context.Background(), "INSERT INTO dcec (stockid, smartcontractaddress, supportedchainids) VALUES ($1, $2, $3)", stockid, smartcontractaddress, supportedchainids)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to insert record: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Record inserted successfully!")
}
