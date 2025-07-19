# Card Battle Game - Online Multiplayer Version

## Project Structure
```
cardgame/
├── server/
│   ├── main.go          # Server entry point
│   ├── server.go        # Server core logic
│   └── deckbuilder.go   # Server deck builder
├── client/
│   ├── main.go          # Client entry point
│   └── client.go        # Client core logic
├── shared/
│   └── messages.go      # Shared message types
├── battle/
│   └── engine.go        # Battle engine (from original)
└── game/
    ├── constants.go     # Game constants
    ├── display.go       # Display functions
    └── ... (other game files)
```

## How to Run

### 1. Start the Server
```bash
go run server/main.go -port 8080
```

### 2. Connect Clients
In separate terminals:
```bash
# Player 1
go run client/main.go -server localhost:8080 -name "Player1"

# Player 2
go run client/main.go -server localhost:8080 -name "Player2"
```

## Features

- Real-time multiplayer over WebSocket
- Automatic matchmaking
- Deck selection (Egyptian/Greek)
- Disconnect handling
- Game state synchronization
- Turn-based gameplay

## Network Architecture

- Server manages all game logic
- Clients send commands and receive state updates
- WebSocket for real-time communication
- JSON message protocol

## Playing Online

1. Both players connect to server
2. Choose deck and join queue
3. Server matches players automatically
4. Play the game in real-time!

## Running on Different Machines

1. Server machine:
   ```bash
   go run server/main.go -port 8080
   ```

2. Client machines:
   ```bash
   go run client/main.go -server <SERVER_IP>:8080
   ```

Replace <SERVER_IP> with the server's IP address.
