package main

import (
	"cardgame/client"
	"flag"
	"fmt"
	"log"
)

func main() {
	// Parse command line flags
	serverAddr := flag.String("server", "localhost:8080", "Server address")
	playerName := flag.String("name", "", "Player name")
	flag.Parse()

	// Create client
	gameClient := client.NewGameClient(*serverAddr)

	// Set player name if provided
	if *playerName != "" {
		gameClient.SetPlayerName(*playerName)
	}

	// Connect to server
	fmt.Printf("🎮 Connecting to Card Battle Game server at %s...\n", *serverAddr)

	if err := gameClient.Connect(); err != nil {
		log.Fatal("Failed to connect:", err)
	}

	// Run the client
	if err := gameClient.Run(); err != nil {
		log.Fatal("Client error:", err)
	}
}
