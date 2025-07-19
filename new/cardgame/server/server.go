package server

import (
	"cardgame/battle"
	"cardgame/shared"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// GameServer manages all online games
type GameServer struct {
	port       string
	games      map[string]*OnlineGame
	players    map[string]*Player
	matchQueue []*Player
	upgrader   websocket.Upgrader
	mu         sync.RWMutex
}

// Player represents a connected player
type Player struct {
	ID         string
	Name       string
	Conn       *websocket.Conn
	GameID     string
	DeckChoice string
	mu         sync.Mutex
}

// OnlineGame represents an online game session
type OnlineGame struct {
	ID         string
	Engine     *battle.BattleEngine
	State      *battle.GameState
	Player1    *Player
	Player2    *Player
	Spectators []*Player
	mu         sync.RWMutex
}

// NewGameServer creates a new game server
func NewGameServer(port string) *GameServer {
	return &GameServer{
		port:    port,
		games:   make(map[string]*OnlineGame),
		players: make(map[string]*Player),
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true // Allow all origins in development
			},
		},
	}
}

// Start starts the game server
func (gs *GameServer) Start() error {
	// Set up routes
	http.HandleFunc("/ws", gs.handleWebSocket)
	http.HandleFunc("/status", gs.handleStatus)

	// Start matchmaking goroutine
	go gs.runMatchmaking()

	return http.ListenAndServe(":"+gs.port, nil)
}

// handleWebSocket handles new WebSocket connections
func (gs *GameServer) handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := gs.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket upgrade error: %v", err)
		return
	}

	// Create new player
	playerID := generateID()
	player := &Player{
		ID:   playerID,
		Conn: conn,
	}

	gs.mu.Lock()
	gs.players[playerID] = player
	gs.mu.Unlock()

	// Send welcome message
	gs.sendToPlayer(player, shared.Message{
		Type: shared.MsgWelcome,
		Data: map[string]interface{}{
			"playerID": playerID,
			"message":  "Welcome to Card Battle Game!",
		},
	})

	// Handle player messages
	go gs.handlePlayer(player)
}

// handlePlayer handles messages from a player
func (gs *GameServer) handlePlayer(player *Player) {
	defer func() {
		gs.disconnectPlayer(player)
		player.Conn.Close()
	}()

	for {
		var msg shared.Message
		err := player.Conn.ReadJSON(&msg)
		if err != nil {
			log.Printf("Read error for player %s: %v", player.ID, err)
			return
		}

		gs.processMessage(player, msg)
	}
}

// processMessage processes a message from a player
func (gs *GameServer) processMessage(player *Player, msg shared.Message) {
	switch msg.Type {
	case shared.MsgSetName:
		gs.handleSetName(player, msg)
	case shared.MsgJoinQueue:
		gs.handleJoinQueue(player, msg)
	case shared.MsgLeaveQueue:
		gs.handleLeaveQueue(player)
	case shared.MsgPlayCard:
		gs.handlePlayCard(player, msg)
	case shared.MsgAttack:
		gs.handleAttack(player, msg)
	case shared.MsgEndTurn:
		gs.handleEndTurn(player)
	case shared.MsgChangePhase:
		gs.handleChangePhase(player, msg)
	case shared.MsgDrawCard:
		gs.handleDrawCard(player)
	}
}

// handleSetName handles player name setting
func (gs *GameServer) handleSetName(player *Player, msg shared.Message) {
	data := msg.Data.(map[string]interface{})
	player.Name = data["name"].(string)

	gs.sendToPlayer(player, shared.Message{
		Type: shared.MsgNameSet,
		Data: map[string]interface{}{
			"name": player.Name,
		},
	})
}

// handleJoinQueue handles player joining matchmaking queue
func (gs *GameServer) handleJoinQueue(player *Player, msg shared.Message) {
	data := msg.Data.(map[string]interface{})
	player.DeckChoice = data["deck"].(string)

	gs.mu.Lock()

	// Check if player is already in queue
	for _, p := range gs.matchQueue {
		if p.ID == player.ID {
			gs.mu.Unlock()
			return
		}
	}

	gs.matchQueue = append(gs.matchQueue, player)
	queueSize := len(gs.matchQueue)
	gs.mu.Unlock()

	gs.sendToPlayer(player, shared.Message{
		Type: shared.MsgQueueJoined,
		Data: map[string]interface{}{
			"position": queueSize,
		},
	})
}

// handleLeaveQueue handles player leaving matchmaking queue
func (gs *GameServer) handleLeaveQueue(player *Player) {
	gs.mu.Lock()
	defer gs.mu.Unlock()

	for i, p := range gs.matchQueue {
		if p.ID == player.ID {
			gs.matchQueue = append(gs.matchQueue[:i], gs.matchQueue[i+1:]...)
			break
		}
	}

	gs.sendToPlayer(player, shared.Message{
		Type: shared.MsgQueueLeft,
	})
}

// runMatchmaking runs the matchmaking loop
func (gs *GameServer) runMatchmaking() {
	for {
		gs.mu.Lock()
		if len(gs.matchQueue) >= 2 {
			// Get first two players
			player1 := gs.matchQueue[0]
			player2 := gs.matchQueue[1]
			gs.matchQueue = gs.matchQueue[2:]

			// Create game
			gs.createGame(player1, player2)
		}
		gs.mu.Unlock()

		// Small delay to prevent busy waiting
		<-time.After(100 * time.Millisecond)
	}
}

// createGame creates a new game between two players
func (gs *GameServer) createGame(player1, player2 *Player) {
	gameID := generateID()

	// Create decks based on player choices
	deckBuilder := &DeckBuilder{}
	var deck1, deck2 []battle.Card

	if player1.DeckChoice == "egyptian" {
		deck1 = deckBuilder.CreateEgyptianDeck()
	} else {
		deck1 = deckBuilder.CreateGreekDeck()
	}

	if player2.DeckChoice == "egyptian" {
		deck2 = deckBuilder.CreateEgyptianDeck()
	} else {
		deck2 = deckBuilder.CreateGreekDeck()
	}

	// Create battle engine and game
	engine := battle.NewBattleEngine()
	gameState, err := engine.CreateMatch(player1.ID, player2.ID, deck1, deck2)
	if err != nil {
		log.Printf("Error creating game: %v", err)
		return
	}

	// Create online game
	onlineGame := &OnlineGame{
		ID:      gameID,
		Engine:  engine,
		State:   gameState,
		Player1: player1,
		Player2: player2,
	}

	// Register game
	gs.mu.Lock()
	gs.games[gameID] = onlineGame
	player1.GameID = gameID
	player2.GameID = gameID
	gs.mu.Unlock()

	// Notify players
	gs.sendToPlayer(player1, shared.Message{
		Type: shared.MsgGameStart,
		Data: map[string]interface{}{
			"gameID":       gameID,
			"playerNum":    1,
			"opponentName": player2.Name,
			"gameState":    gameState,
		},
	})

	gs.sendToPlayer(player2, shared.Message{
		Type: shared.MsgGameStart,
		Data: map[string]interface{}{
			"gameID":       gameID,
			"playerNum":    2,
			"opponentName": player1.Name,
			"gameState":    gameState,
		},
	})

	log.Printf("Game %s started: %s vs %s", gameID, player1.Name, player2.Name)
}

// Game action handlers

func (gs *GameServer) handlePlayCard(player *Player, msg shared.Message) {
	game := gs.getPlayerGame(player)
	if game == nil {
		return
	}

	data := msg.Data.(map[string]interface{})
	cardIndex := int(data["cardIndex"].(float64))

	game.mu.Lock()
	err := game.Engine.PlayCard(game.State.ID, player.ID, cardIndex)
	if err != nil {
		game.mu.Unlock()
		gs.sendError(player, err.Error())
		return
	}

	// Get updated state
	game.State, _ = game.Engine.GetGameState(game.State.ID)
	game.mu.Unlock()

	// Broadcast to both players
	gs.broadcastGameState(game)
}

func (gs *GameServer) handleAttack(player *Player, msg shared.Message) {
	game := gs.getPlayerGame(player)
	if game == nil {
		return
	}

	data := msg.Data.(map[string]interface{})
	attackerIndex := int(data["attackerIndex"].(float64))
	targetIndex := int(data["targetIndex"].(float64))

	game.mu.Lock()
	err := game.Engine.Attack(game.State.ID, player.ID, attackerIndex, targetIndex)
	if err != nil {
		game.mu.Unlock()
		gs.sendError(player, err.Error())
		return
	}

	// Get updated state
	game.State, _ = game.Engine.GetGameState(game.State.ID)
	game.mu.Unlock()

	// Broadcast to both players
	gs.broadcastGameState(game)

	// Check for game over
	if game.State.GameOver {
		gs.handleGameOver(game)
	}
}

func (gs *GameServer) handleEndTurn(player *Player) {
	game := gs.getPlayerGame(player)
	if game == nil {
		return
	}

	game.mu.Lock()
	err := game.Engine.EndTurn(game.State.ID, player.ID)
	if err != nil {
		game.mu.Unlock()
		gs.sendError(player, err.Error())
		return
	}

	// Get updated state
	game.State, _ = game.Engine.GetGameState(game.State.ID)
	game.mu.Unlock()

	// Broadcast to both players
	gs.broadcastGameState(game)
}

func (gs *GameServer) handleChangePhase(player *Player, msg shared.Message) {
	game := gs.getPlayerGame(player)
	if game == nil {
		return
	}

	data := msg.Data.(map[string]interface{})
	phase := battle.GamePhase(data["phase"].(string))

	game.mu.Lock()
	err := game.Engine.ChangePhase(game.State.ID, player.ID, phase)
	if err != nil {
		game.mu.Unlock()
		gs.sendError(player, err.Error())
		return
	}

	// Get updated state
	game.State, _ = game.Engine.GetGameState(game.State.ID)
	game.mu.Unlock()

	// Broadcast to both players
	gs.broadcastGameState(game)
}

func (gs *GameServer) handleDrawCard(player *Player) {
	game := gs.getPlayerGame(player)
	if game == nil {
		return
	}

	game.mu.Lock()
	err := game.Engine.DrawCard(game.State.ID, player.ID)
	if err != nil {
		game.mu.Unlock()
		gs.sendError(player, err.Error())
		return
	}

	// Get updated state
	game.State, _ = game.Engine.GetGameState(game.State.ID)
	game.mu.Unlock()

	// Broadcast to both players
	gs.broadcastGameState(game)
}

// Helper functions

func (gs *GameServer) getPlayerGame(player *Player) *OnlineGame {
	gs.mu.RLock()
	defer gs.mu.RUnlock()

	if player.GameID == "" {
		return nil
	}

	return gs.games[player.GameID]
}

func (gs *GameServer) broadcastGameState(game *OnlineGame) {
	game.mu.RLock()
	state := game.State
	game.mu.RUnlock()

	msg := shared.Message{
		Type: shared.MsgGameUpdate,
		Data: map[string]interface{}{
			"gameState": state,
		},
	}

	gs.sendToPlayer(game.Player1, msg)
	gs.sendToPlayer(game.Player2, msg)
}

func (gs *GameServer) handleGameOver(game *OnlineGame) {
	// Determine winner name
	winnerName := ""
	if game.State.Winner == game.Player1.ID {
		winnerName = game.Player1.Name
	} else {
		winnerName = game.Player2.Name
	}

	msg := shared.Message{
		Type: shared.MsgGameOver,
		Data: map[string]interface{}{
			"winner":     game.State.Winner,
			"winnerName": winnerName,
		},
	}

	gs.sendToPlayer(game.Player1, msg)
	gs.sendToPlayer(game.Player2, msg)

	// Clean up game
	gs.mu.Lock()
	delete(gs.games, game.ID)
	game.Player1.GameID = ""
	game.Player2.GameID = ""
	gs.mu.Unlock()
}

func (gs *GameServer) sendToPlayer(player *Player, msg shared.Message) {
	player.mu.Lock()
	defer player.mu.Unlock()

	if err := player.Conn.WriteJSON(msg); err != nil {
		log.Printf("Error sending to player %s: %v", player.ID, err)
	}
}

func (gs *GameServer) sendError(player *Player, errMsg string) {
	gs.sendToPlayer(player, shared.Message{
		Type: shared.MsgError,
		Data: map[string]interface{}{
			"error": errMsg,
		},
	})
}

func (gs *GameServer) disconnectPlayer(player *Player) {
	gs.mu.Lock()
	defer gs.mu.Unlock()

	// Remove from players map
	delete(gs.players, player.ID)

	// Remove from queue
	for i, p := range gs.matchQueue {
		if p.ID == player.ID {
			gs.matchQueue = append(gs.matchQueue[:i], gs.matchQueue[i+1:]...)
			break
		}
	}

	// Handle game disconnection
	if player.GameID != "" {
		if game, exists := gs.games[player.GameID]; exists {
			// Notify opponent
			var opponent *Player
			if game.Player1.ID == player.ID {
				opponent = game.Player2
			} else {
				opponent = game.Player1
			}

			gs.sendToPlayer(opponent, shared.Message{
				Type: shared.MsgOpponentDisconnected,
				Data: map[string]interface{}{
					"message": "Your opponent has disconnected",
				},
			})

			// Clean up game
			delete(gs.games, player.GameID)
		}
	}

	log.Printf("Player %s disconnected", player.ID)
}

// handleStatus handles status endpoint
func (gs *GameServer) handleStatus(w http.ResponseWriter, r *http.Request) {
	gs.mu.RLock()
	status := map[string]interface{}{
		"players":     len(gs.players),
		"games":       len(gs.games),
		"queueLength": len(gs.matchQueue),
	}
	gs.mu.RUnlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(status)
}

// generateID generates a unique ID
func generateID() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}
