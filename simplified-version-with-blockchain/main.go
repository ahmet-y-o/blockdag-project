package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"

	"cardgame/api/handlers"
	"cardgame/blockchain"
)

func init() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found")
	}
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func main() {
	// Get configuration from environment variables
	nodeURL := getEnvOrDefault("BLOCKCHAIN_RPC_URL", "http://localhost:8545")
	contractAddress := getEnvOrDefault("CONTRACT_ADDRESS", "")
	port := getEnvOrDefault("API_PORT", "8080")

	// Validate contract address
	if contractAddress == "" {
		log.Fatal("CONTRACT_ADDRESS must be set in environment variables")
	}

	// Initialize Ethereum client
	ethClient, err := blockchain.NewEthereumClient(nodeURL, contractAddress)
	if err != nil {
		log.Fatalf("Failed to initialize Ethereum client: %v", err)
	}

	router := mux.NewRouter()

	// Initialize handlers
	battleHandler := handlers.NewBattleHandler()
	tradingHandler := handlers.NewTradingHandler(ethClient)

	// Battle routes
	router.HandleFunc("/battle/create", battleHandler.CreateBattle).Methods("POST")
	router.HandleFunc("/battle/{battleID}/play-card", battleHandler.PlayCard).Methods("POST")
	router.HandleFunc("/battle/{battleID}", battleHandler.GetBattleState).Methods("GET")

	// Trading routes
	router.HandleFunc("/trading/list", tradingHandler.ListCard).Methods("POST")
	router.HandleFunc("/trading/buy", tradingHandler.BuyCard).Methods("POST")

	// Middleware
	router.Use(loggingMiddleware)
	router.Use(jsonContentTypeMiddleware)

	// CORS middleware
	router.Use(corsMiddleware)

	// Health check
	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}).Methods("GET")

	// Start server
	log.Printf("Server starting on :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}

// Middleware for logging
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s", r.Method, r.RequestURI)
		next.ServeHTTP(w, r)
	})
}

// Middleware for JSON content type
func jsonContentTypeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

// Middleware for CORS
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
