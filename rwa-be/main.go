package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/andybalholm/brotli"
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
			stockname TEXT UNIQUE, 
			stockid TEXT UNIQUE, 
			stocktinker TEXT);
		CREATE TABLE IF NOT EXISTS dcec (
			stockid TEXT PRIMARY KEY,
			smartcontractaddress TEXT,
			supportedchainids INTEGER);
	`)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to create table: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Table created successfully!")

	// Insert a new record into the dces table
	// given a URL send GET request to the URL and get the response
	url := "https://api-enterprise.sandbox.dinari.com/api/v1/stocks"

	// dont use token in url, add bearer token in header
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to create GET request: %v\n", err)
		os.Exit(1)
	}
	req.Header.Set("Authorization", "Bearer xGL0hv2V1th-C--ZCXx_mLU00wztW3w8KQ3lGOQi7N8")
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Accept-Encoding", "gzip, deflate, br")
	req.Header.Set("Connection", "keep-alive")

	// Send the GET request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to send GET request: %v\n", err)
		os.Exit(1)
	}
	defer resp.Body.Close()
	// Check the response status code
	if resp.StatusCode != 200 {
		fmt.Fprintf(os.Stderr, "GET request failed with status code: %d\n %s\n", resp.StatusCode, resp.Status)
		os.Exit(1)
	}
	fmt.Fprint(os.Stdout, "GET request succeeded with status code: ", resp.StatusCode, "\n")
	fmt.Println("Content-Type:", resp.Header.Get("Content-Type"))

	// Handle brotli encoding
	var reader io.Reader = resp.Body
	if resp.Header.Get("Content-Encoding") == "br" {
		// Use a Brotli reader to decompress the response body
		reader = brotli.NewReader(resp.Body)
	}

	// // Read the response body
	// body, _ := io.ReadAll(reader)
	// if err != nil {
	// 	fmt.Fprintf(os.Stderr, "Unable to read response body: %v\n", err)
	// 	os.Exit(1)
	// }
	//fmt.Println("Response body:", string(body))

	type stockData struct {
		StockName   string `json:"name"`
		StockID     string `json:"id"`
		StockTicker string `json:"symbol"`
	}
	var stockDataList []stockData

	err = json.NewDecoder(reader).Decode(&stockDataList)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to decode JSON response: %v\n", err)
		os.Exit(1)
	}

	// // Print the stock data
	for _, stock := range stockDataList {
		fmt.Printf("Stock Name: %s\n", stock.StockName)
		fmt.Printf("Stock ID: %s\n", stock.StockID)
		fmt.Printf("Stock Ticker: %s\n", stock.StockTicker)
		_, err = conn.Exec(context.Background(), "INSERT INTO dces (stockname, stockid, stocktinker) VALUES ($1, $2, $3)", stock.StockName, stock.StockID, stock.StockTicker)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable to insert record: %v\n", err)
		}
		fmt.Println("Record inserted successfully!")
	}
	// list of chainids that are supported
	type chainID struct {
		name      string
		title     string
		chainId   int
		networkId int
	}
	var chainIDList []chainID

	chainIDList = append(chainIDList, chainID{"Ethereum Sepolia", "Ethereum Testnet Sepolia", 11155111, 11155111})
	chainIDList = append(chainIDList, chainID{"Base Sepolia Testnet", "Base Sepolia Testnet", 84532, 84532})
	chainIDList = append(chainIDList, chainID{"Arbitrum Sepolia", "Arbitrum Sepolia Rollup Testnet", 421614, 421614})
	chainIDList = append(chainIDList, chainID{"Blast Sepolia Testnet", "Blast Sepolia Testnet", 168587773, 168587773})
	chainIDList = append(chainIDList, chainID{"Plume Testnet", "Plume Sepolia L2 Rollup Testnet", 98867, 98867})

	type Stock struct {
		Id     string `json:"id"`
		Symbol string `json:"symbol"`
	}
	type Token struct {
		Address  string `json:"address"`
		ChainId  int    `json:"chain_id"`
		Decimals int    `json:"decimals"`
	}
	type ChainIDData struct {
		Stock Stock `json:"stock"`
		Token Token `json:"token"`
	}
	var ChainIDDataList []ChainIDData

	for _, chainId := range chainIDList {
		fmt.Printf("Supported Chain ID: %d\n", chainId.chainId)
		fmt.Printf("Supported Network ID: %d\n", chainId.networkId)
		fmt.Printf("Supported Chain Name: %s\n", chainId.name)
		fmt.Printf("Supported Chain Title: %s\n", chainId.title)

		url := "https://api-enterprise.sandbox.dinari.com/api/v1/tokens/" + fmt.Sprint(chainId.chainId)
		fmt.Println("URL:", url)
		// dont use token in url, add bearer token in header
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable to create GET request: %v\n", err)
			os.Exit(1)
		}
		req.Header.Set("Authorization", "Bearer xGL0hv2V1th-C--ZCXx_mLU00wztW3w8KQ3lGOQi7N8")
		req.Header.Set("Accept", "*/*")
		req.Header.Set("Accept-Encoding", "gzip, deflate, br")
		req.Header.Set("Connection", "keep-alive")
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable to send GET request: %v\n", err)
			os.Exit(1)
		}
		defer resp.Body.Close()
		// Check the response status code
		if resp.StatusCode != 200 {
			fmt.Fprintf(os.Stderr, "GET request failed with status code: %d\n %s\n", resp.StatusCode, resp.Status)
			os.Exit(1)
		}
		fmt.Fprint(os.Stdout, "GET request succeeded with status code: ", resp.StatusCode, "\n")
		fmt.Println("Content-Type:", resp.Header.Get("Content-Type"))

		// Handle brotli encoding
		var reader io.Reader = resp.Body
		if resp.Header.Get("Content-Encoding") == "br" {
			// Use a Brotli reader to decompress the response body
			reader = brotli.NewReader(resp.Body)
		}
		err = json.NewDecoder(reader).Decode(&ChainIDDataList)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable to decode JSON response: %v\n", err)
			os.Exit(1)
		}
		for _, stockdata := range ChainIDDataList {
			fmt.Printf("Stock ID: %s\n", stockdata.Stock.Id)
			fmt.Printf("Stock Symbol: %s\n", stockdata.Stock.Symbol)
			fmt.Printf("Stock Address: %s\n", stockdata.Token.Address)
			fmt.Printf("Stock Chain ID: %d\n", stockdata.Token.ChainId)
			fmt.Printf("Stock Decimals: %d\n", stockdata.Token.Decimals)
			_, err = conn.Exec(context.Background(), "INSERT INTO dcec (stockid, smartcontractaddress, supportedchainids) VALUES ($1, $2, $3)", stockdata.Stock.Id, stockdata.Token.Address, stockdata.Token.ChainId)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Unable to insert record: %v\n", err)
			}
			fmt.Println("Record inserted successfully!")
		}
	}
}
