package controllers

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/jackc/pgx/v5"
	"rwa.abc/views"
)

type Stock struct {
	Name      string
	Templates struct {
		New   views.Template
		Query views.Template
	}
}

type StockData struct {
	StockName   string
	StockId     string
	StockTinker string
}

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

func (s Stock) NewQuery(w http.ResponseWriter, r *http.Request) {
	// Parse the form data
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}
	// Render the template with the parsed form data
	s.Templates.New.Execute(w, nil)
}

func (s Stock) StockQuery(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}
	data := StockData{
		StockTinker: r.FormValue("ticker"),
		StockName:   "",
		StockId:     "",
	}
	//pg query

	//fmt.Fprint(w, "Stock Name You want to query is: ", stockName)

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
	// Query the database
	rows, err := conn.Query(context.Background(), "SELECT * FROM dces WHERE stocktinker = $1", data.StockTinker)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Query failed: %v\n", err)
		os.Exit(1)
	}
	defer rows.Close()
	// Iterate through the rows and print the results
	for rows.Next() {
		err := rows.Scan(&data.StockName, &data.StockId, &data.StockTinker)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Row scan failed: %v\n", err)
			os.Exit(1)
		}
	}
	s.Templates.Query.Execute(w, data)
}
