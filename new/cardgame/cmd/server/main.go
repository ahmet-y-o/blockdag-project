package main

import (
	"cardgame/server"
	"flag"
	"fmt"
	"log"
)

func main() {
	// Parse command line flags
	port := flag.String("port", "8080", "Server port")
	flag.Parse()

	// Create and start server
	gameServer := server.NewGameServer(*port)

	fmt.Printf("🎮 Card Battle Game Server starting on port %s...\n", *port)
	fmt.Println("Players can connect using: go run cmd/client/main.go -server localhost:" + *port)

	if err := gameServer.Start(); err != nil {
		log.Fatal("Server error:", err)
	}
}
