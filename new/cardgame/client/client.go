package client

import (
	"bufio"
	"cardgame/battle"
	"cardgame/game"
	"cardgame/shared"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/gorilla/websocket"
)

// GameClient handles the client side of the game
type GameClient struct {
	serverAddr string
	conn       *websocket.Conn
	playerID   string
	playerName string
	playerNum  int
	gameID     string
	gameState  *battle.GameState
	display    *game.Display
	input      *bufio.Reader
	mu         sync.RWMutex
	connected  bool
	inGame     bool
	inQueue    bool
}

// NewGameClient creates a new game client
func NewGameClient(serverAddr string) *GameClient {
	return &GameClient{
		serverAddr: serverAddr,
		display:    game.NewDisplay(),
		input:      bufio.NewReader(os.Stdin),
	}
}

// SetPlayerName sets the player name
func (gc *GameClient) SetPlayerName(name string) {
	gc.playerName = name
}

// Connect connects to the game server
func (gc *GameClient) Connect() error {
	url := "ws://" + gc.serverAddr + "/ws"
	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		return fmt.Errorf("failed to connect: %v", err)
	}

	gc.conn = conn
	gc.connected = true

	// Start message handler
	go gc.handleMessages()

	return nil
}

// Run runs the client main loop
func (gc *GameClient) Run() error {
	// Wait for welcome message
	for gc.playerID == "" {
		// Small delay
	}

	// Get player name if not set
	if gc.playerName == "" {
		gc.display.ClearScreen()
		fmt.Println(game.ColorCyan + "=====================================" + game.ColorReset)
		fmt.Println(game.ColorCyan + "    CARD BATTLE GAME - ONLINE       " + game.ColorReset)
		fmt.Println(game.ColorCyan + "=====================================" + game.ColorReset)
		fmt.Print("\nEnter your name: ")
		name, _ := gc.input.ReadString('\n')
		gc.playerName = strings.TrimSpace(name)
	}

	// Send name to server
	gc.sendMessage(shared.Message{
		Type: shared.MsgSetName,
		Data: map[string]interface{}{
			"name": gc.playerName,
		},
	})

	// Main menu loop
	for gc.connected {
		if gc.inGame {
			gc.runGame()
		} else {
			gc.showMainMenu()
		}
	}

	return nil
}

// showMainMenu shows the main menu
func (gc *GameClient) showMainMenu() {
	gc.display.ClearScreen()
	fmt.Println(game.ColorCyan + "=====================================" + game.ColorReset)
	fmt.Println(game.ColorCyan + "         MAIN MENU                  " + game.ColorReset)
	fmt.Println(game.ColorCyan + "=====================================" + game.ColorReset)
	fmt.Printf("\nWelcome, %s%s%s!\n", game.ColorGreen, gc.playerName, game.ColorReset)

	if gc.inQueue {
		fmt.Println(game.ColorYellow + "\n⏳ Waiting for opponent..." + game.ColorReset)
		fmt.Println("\nOptions:")
		fmt.Println("1. Leave Queue")
		fmt.Println("2. Quit")
	} else {
		fmt.Println("\nOptions:")
		fmt.Println("1. Find Match")
		fmt.Println("2. Quit")
	}

	fmt.Print("\nChoice: ")
	choice, _ := gc.input.ReadString('\n')
	choice = strings.TrimSpace(choice)

	switch choice {
	case "1":
		if gc.inQueue {
			gc.leaveQueue()
		} else {
			gc.joinQueue()
		}
	case "2":
		gc.connected = false
		gc.conn.Close()
	}
}

// joinQueue joins the matchmaking queue
func (gc *GameClient) joinQueue() {
	// Select deck
	gc.display.ClearScreen()
	fmt.Println(game.ColorCyan + "=====================================" + game.ColorReset)
	fmt.Println(game.ColorCyan + "        DECK SELECTION              " + game.ColorReset)
	fmt.Println(game.ColorCyan + "=====================================" + game.ColorReset)
	fmt.Println("\n1. Egyptian Gods Deck (Attack focused)")
	fmt.Println("2. Greek Gods Deck (Defense focused)")

	fmt.Print("\nChoose your deck (1 or 2): ")
	choice, _ := gc.input.ReadString('\n')
	choice = strings.TrimSpace(choice)

	deck := "egyptian"
	if choice == "2" {
		deck = "greek"
	}

	// Send join queue message
	gc.sendMessage(shared.Message{
		Type: shared.MsgJoinQueue,
		Data: map[string]interface{}{
			"deck": deck,
		},
	})
}

// leaveQueue leaves the matchmaking queue
func (gc *GameClient) leaveQueue() {
	gc.sendMessage(shared.Message{
		Type: shared.MsgLeaveQueue,
	})
	gc.inQueue = false
}

// runGame runs the game loop
func (gc *GameClient) runGame() {
	gc.mu.RLock()
	state := gc.gameState
	gc.mu.RUnlock()

	if state == nil {
		return
	}

	// Display game state
	gc.displayGameState()

	// Check if it's our turn
	isOurTurn := (gc.playerNum == 1 && state.CurrentTurn == state.Player1.ID) ||
		(gc.playerNum == 2 && state.CurrentTurn == state.Player2.ID)

	if !isOurTurn {
		fmt.Println("\n" + game.ColorYellow + "Waiting for opponent's move..." + game.ColorReset)
		// Just wait for updates
		return
	}

	// Auto-draw at start of turn
	if state.Phase == battle.PhaseDrawn {
		fmt.Println("\n" + game.ColorGreen + "Drawing card..." + game.ColorReset)
		gc.sendMessage(shared.Message{
			Type: shared.MsgDrawCard,
		})
		return
	}

	// Show commands and get input
	gc.display.ShowCommands(state.Phase, true)

	fmt.Printf("\n%s > ", gc.playerName)
	input, _ := gc.input.ReadString('\n')
	input = strings.TrimSpace(input)

	parts := strings.Split(input, " ")
	if len(parts) == 0 {
		return
	}

	command := strings.ToLower(parts[0])
	gc.processCommand(command, parts[1:])
}

// processCommand processes player commands
func (gc *GameClient) processCommand(command string, args []string) {
	switch command {
	case "play":
		if len(args) < 1 {
			fmt.Println(game.ColorRed + "Usage: play [card number]" + game.ColorReset)
			return
		}
		cardIndex, err := strconv.Atoi(args[0])
		if err != nil {
			fmt.Println(game.ColorRed + "Invalid card number" + game.ColorReset)
			return
		}
		gc.sendMessage(shared.Message{
			Type: shared.MsgPlayCard,
			Data: map[string]interface{}{
				"cardIndex": cardIndex,
			},
		})

	case "attack":
		if len(args) < 2 {
			fmt.Println(game.ColorRed + "Usage: attack [attacker] [target]" + game.ColorReset)
			return
		}
		attackerIndex, err1 := strconv.Atoi(args[0])
		targetIndex, err2 := strconv.Atoi(args[1])
		if err1 != nil || err2 != nil {
			fmt.Println(game.ColorRed + "Invalid indices" + game.ColorReset)
			return
		}
		gc.sendMessage(shared.Message{
			Type: shared.MsgAttack,
			Data: map[string]interface{}{
				"attackerIndex": attackerIndex,
				"targetIndex":   targetIndex,
			},
		})

	case "battle":
		gc.sendMessage(shared.Message{
			Type: shared.MsgChangePhase,
			Data: map[string]interface{}{
				"phase": string(battle.PhaseBattle),
			},
		})

	case "main":
		gc.sendMessage(shared.Message{
			Type: shared.MsgChangePhase,
			Data: map[string]interface{}{
				"phase": string(battle.PhaseMain),
			},
		})

	case "end":
		gc.sendMessage(shared.Message{
			Type: shared.MsgEndTurn,
		})

	case "help":
		// Help is already shown

	case "quit":
		gc.inGame = false
		gc.gameState = nil
	}
}

// displayGameState displays the current game state
func (gc *GameClient) displayGameState() {
	gc.mu.RLock()
	state := gc.gameState
	playerNum := gc.playerNum
	gc.mu.RUnlock()

	if state == nil {
		return
	}

	gc.display.ClearScreen()
	fmt.Println(game.ColorCyan + "=====================================" + game.ColorReset)
	fmt.Println(game.ColorCyan + "      CARD BATTLE GAME - ONLINE     " + game.ColorReset)
	fmt.Println(game.ColorCyan + "=====================================" + game.ColorReset)
	fmt.Printf("Turn: %d | Phase: %s\n", state.TurnCount, state.Phase)

	// Determine which player we are
	var ourPlayer, opponentPlayer *battle.Player
	if playerNum == 1 {
		ourPlayer = state.Player1
		opponentPlayer = state.Player2
	} else {
		ourPlayer = state.Player2
		opponentPlayer = state.Player1
	}

	currentTurnIndicator := ""
	if state.CurrentTurn == ourPlayer.ID {
		currentTurnIndicator = game.ColorGreen + " (Your turn)" + game.ColorReset
	} else {
		currentTurnIndicator = game.ColorRed + " (Opponent's turn)" + game.ColorReset
	}
	fmt.Printf("Current Turn: %s%s\n\n", state.CurrentTurn, currentTurnIndicator)

	// Display Opponent
	fmt.Printf("%s=== Opponent ===%s\n", game.ColorRed, game.ColorReset)
	fmt.Printf("HP: %s%d%s | Mana: %d/%d\n", game.ColorRed, opponentPlayer.HP, game.ColorReset, opponentPlayer.Mana, opponentPlayer.MaxMana)
	fmt.Printf("Hand: %d cards | Deck: %d cards\n", len(opponentPlayer.Hand), len(opponentPlayer.Deck))

	fmt.Println("\nOpponent's Field:")
	if len(opponentPlayer.Field) == 0 {
		fmt.Println("  (empty)")
	} else {
		for i, card := range opponentPlayer.Field {
			color := gc.getCardColor(card.Archetype)
			fmt.Printf("  [%d] %s%s%s (%s) - ATK: %d / DEF: %d\n",
				i, color, card.Name, game.ColorReset, card.Archetype, card.Attack, card.Defense)
		}
	}

	fmt.Println("\n" + game.ColorWhite + "-------------------------------------" + game.ColorReset)

	// Display Our Player
	fmt.Printf("\n%s=== You ===%s\n", game.ColorGreen, game.ColorReset)
	fmt.Printf("HP: %s%d%s | Mana: %d/%d\n", game.ColorGreen, ourPlayer.HP, game.ColorReset, ourPlayer.Mana, ourPlayer.MaxMana)
	fmt.Printf("Deck: %d cards\n", len(ourPlayer.Deck))

	fmt.Println("\nYour Field:")
	if len(ourPlayer.Field) == 0 {
		fmt.Println("  (empty)")
	} else {
		for i, card := range ourPlayer.Field {
			color := gc.getCardColor(card.Archetype)
			fmt.Printf("  [%d] %s%s%s (%s) - ATK: %d / DEF: %d\n",
				i, color, card.Name, game.ColorReset, card.Archetype, card.Attack, card.Defense)
		}
	}

	fmt.Println("\nYour Hand:")
	for i, card := range ourPlayer.Hand {
		color := gc.getCardColor(card.Archetype)
		canPlay := ""
		if card.Cost > ourPlayer.Mana {
			canPlay = game.ColorRed + " (Not enough mana)" + game.ColorReset
		}

		effectStr := ""
		if card.Effect != "" {
			effectStr = fmt.Sprintf(" - %s", card.Effect)
		}

		fmt.Printf("  [%d] %s%s%s (Cost: %d) - ATK: %d / DEF: %d%s%s\n",
			i, color, card.Name, game.ColorReset, card.Cost, card.Attack, card.Defense, effectStr, canPlay)
	}

	if state.LastAction != "" {
		fmt.Printf("\n%sLast Action: %s%s\n", game.ColorPurple, state.LastAction, game.ColorReset)
	}

	fmt.Println("\n" + game.ColorCyan + "=====================================" + game.ColorReset)
}

// handleMessages handles messages from the server
func (gc *GameClient) handleMessages() {
	for gc.connected {
		var msg shared.Message
		err := gc.conn.ReadJSON(&msg)
		if err != nil {
			fmt.Printf("\n%sDisconnected from server%s\n", game.ColorRed, game.ColorReset)
			gc.connected = false
			return
		}

		gc.processServerMessage(msg)
	}
}

// processServerMessage processes a message from the server
func (gc *GameClient) processServerMessage(msg shared.Message) {
	switch msg.Type {
	case shared.MsgWelcome:
		data := msg.Data.(map[string]interface{})
		gc.playerID = data["playerID"].(string)

	case shared.MsgNameSet:
		// Name confirmed

	case shared.MsgQueueJoined:
		gc.inQueue = true
		data := msg.Data.(map[string]interface{})
		position := int(data["position"].(float64))
		fmt.Printf("\n%sJoined queue at position %d%s\n", game.ColorGreen, position, game.ColorReset)

	case shared.MsgQueueLeft:
		gc.inQueue = false

	case shared.MsgGameStart:
		gc.handleGameStart(msg)

	case shared.MsgGameUpdate:
		gc.handleGameUpdate(msg)

	case shared.MsgGameOver:
		gc.handleGameOver(msg)

	case shared.MsgError:
		data := msg.Data.(map[string]interface{})
		fmt.Printf("\n%sError: %s%s\n", game.ColorRed, data["error"], game.ColorReset)
		fmt.Print("Press Enter to continue...")
		gc.input.ReadString('\n')

	case shared.MsgOpponentDisconnected:
		fmt.Printf("\n%sYour opponent has disconnected!%s\n", game.ColorRed, game.ColorReset)
		fmt.Print("Press Enter to return to menu...")
		gc.input.ReadString('\n')
		gc.inGame = false
		gc.gameState = nil
	}
}

// handleGameStart handles game start message
func (gc *GameClient) handleGameStart(msg shared.Message) {
	data := msg.Data.(map[string]interface{})
	gc.gameID = data["gameID"].(string)
	gc.playerNum = int(data["playerNum"].(float64))
	opponentName := data["opponentName"].(string)

	// Convert game state
	stateData := data["gameState"].(map[string]interface{})
	gc.gameState = shared.ConvertToGameState(stateData)

	gc.inGame = true
	gc.inQueue = false

	fmt.Printf("\n%sMatch found! You are playing against %s%s\n",
		game.ColorGreen, opponentName, game.ColorReset)
	fmt.Print("Press Enter to start...")
	gc.input.ReadString('\n')
}

// handleGameUpdate handles game state update
func (gc *GameClient) handleGameUpdate(msg shared.Message) {
	data := msg.Data.(map[string]interface{})
	stateData := data["gameState"].(map[string]interface{})

	gc.mu.Lock()
	gc.gameState = shared.ConvertToGameState(stateData)
	gc.mu.Unlock()
}

// handleGameOver handles game over message
func (gc *GameClient) handleGameOver(msg shared.Message) {
	data := msg.Data.(map[string]interface{})
	winnerID := data["winner"].(string)
	winnerName := data["winnerName"].(string)

	gc.display.ClearScreen()
	fmt.Println(game.ColorCyan + "=====================================" + game.ColorReset)
	fmt.Println(game.ColorCyan + "           GAME OVER!              " + game.ColorReset)
	fmt.Println(game.ColorCyan + "=====================================" + game.ColorReset)

	if (gc.playerNum == 1 && winnerID == gc.gameState.Player1.ID) ||
		(gc.playerNum == 2 && winnerID == gc.gameState.Player2.ID) {
		fmt.Printf("\n%sCongratulations! You won!%s\n", game.ColorGreen, game.ColorReset)
	} else {
		fmt.Printf("\n%s%s wins!%s\n", game.ColorRed, winnerName, game.ColorReset)
	}

	fmt.Println("\nFinal Stats:")
	fmt.Printf("Your HP: %d\n", gc.getOurPlayer().HP)
	fmt.Printf("Opponent HP: %d\n", gc.getOpponentPlayer().HP)
	fmt.Printf("Total Turns: %d\n", gc.gameState.TurnCount)

	fmt.Print("\nPress Enter to return to menu...")
	gc.input.ReadString('\n')

	gc.inGame = false
	gc.gameState = nil
}

// Helper functions

func (gc *GameClient) sendMessage(msg shared.Message) {
	if err := gc.conn.WriteJSON(msg); err != nil {
		fmt.Printf("\n%sError sending message: %v%s\n", game.ColorRed, err, game.ColorReset)
	}
}

func (gc *GameClient) getOurPlayer() *battle.Player {
	if gc.playerNum == 1 {
		return gc.gameState.Player1
	}
	return gc.gameState.Player2
}

func (gc *GameClient) getOpponentPlayer() *battle.Player {
	if gc.playerNum == 1 {
		return gc.gameState.Player2
	}
	return gc.gameState.Player1
}

func (gc *GameClient) getCardColor(archetype battle.Archetype) string {
	switch archetype {
	case battle.ArchetypeEgyptian:
		return game.ColorYellow
	case battle.ArchetypeGreek:
		return game.ColorBlue
	default:
		return game.ColorWhite
	}
}
