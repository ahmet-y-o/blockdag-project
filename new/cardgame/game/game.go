package game

import (
	"fmt"
	"time"
	"cardgame/battle"
)

// Game represents the main game controller
type Game struct {
	engine      *battle.BattleEngine
	display     *Display
	input       *InputHandler
	ai          *AIPlayer
	deckBuilder *DeckBuilder
	gameState   *battle.GameState
}

// NewGame creates a new game instance
func NewGame() *Game {
	return &Game{
		engine:      battle.NewBattleEngine(),
		display:     NewDisplay(),
		input:       NewInputHandler(),
		ai:          NewAIPlayer("normal"),
		deckBuilder: NewDeckBuilder(),
	}
}

// Run starts and runs the game
func (g *Game) Run() error {
	// Show welcome screen
	g.display.ShowWelcome()
	
	// Get player deck choice
	playerDeck, aiDeck := g.selectDecks()
	
	// Create the match
	gameState, err := g.engine.CreateMatch("Player", "AI", playerDeck, aiDeck)
	if err != nil {
		return fmt.Errorf("failed to create match: %v", err)
	}
	g.gameState = gameState
	
	// Main game loop
	g.runGameLoop()
	
	// Show game over screen
	g.display.ShowGameOver(g.gameState)
	
	return nil
}

// selectDecks handles deck selection
func (g *Game) selectDecks() ([]battle.Card, []battle.Card) {
	choice := g.input.GetDeckChoice()
	
	var playerDeck, aiDeck []battle.Card
	
	if choice == "2" {
		playerDeck = g.deckBuilder.CreateGreekDeck()
		aiDeck = g.deckBuilder.CreateEgyptianDeck()
		g.display.ShowMessage("You chose Greek Gods!", ColorGreen)
		g.display.ShowMessage("AI opponent will use Egyptian Gods!", ColorRed)
	} else {
		playerDeck = g.deckBuilder.CreateEgyptianDeck()
		aiDeck = g.deckBuilder.CreateGreekDeck()
		g.display.ShowMessage("You chose Egyptian Gods!", ColorYellow)
		g.display.ShowMessage("AI opponent will use Greek Gods!", ColorBlue)
	}
	
	g.input.WaitForEnter("\nPress Enter to start the game...")
	return playerDeck, aiDeck
}

// runGameLoop runs the main game loop
func (g *Game) runGameLoop() {
	for !g.gameState.GameOver {
		// Update game state
		g.gameState, _ = g.engine.GetGameState(g.gameState.ID)
		
		// Display current state
		g.display.ShowGameState(g.gameState)
		
		// Check whose turn it is
		if g.gameState.CurrentTurn == "AI" {
			g.handleAITurn()
		} else {
			g.handlePlayerTurn()
		}
	}
}

// handleAITurn handles the AI's turn
func (g *Game) handleAITurn() {
	g.display.ShowCommands(g.gameState.Phase, false)
	g.ai.MakeDecision(g.gameState, g.engine, "AI")
}

// handlePlayerTurn handles the player's turn
func (g *Game) handlePlayerTurn() {
	// Auto-draw at start of turn
	if g.gameState.Phase == battle.PhaseDrawn {
		g.display.ShowMessage("Drawing card...", ColorGreen)
		g.engine.DrawCard(g.gameState.ID, "Player")
		time.Sleep(500 * time.Millisecond)
		return
	}
	
	// Show available commands
	g.display.ShowCommands(g.gameState.Phase, true)
	
	// Get and process player command
	command, args := g.input.GetCommand()
	g.processCommand(command, args)
}

// processCommand processes player commands
func (g *Game) processCommand(command string, args []string) {
	var err error
	
	switch command {
	case "play":
		err = g.handlePlayCommand(args)
		
	case "attack":
		err = g.handleAttackCommand(args)
		
	case "battle":
		err = g.engine.ChangePhase(g.gameState.ID, "Player", battle.PhaseBattle)
		
	case "main":
		err = g.engine.ChangePhase(g.gameState.ID, "Player", battle.PhaseMain)
		
	case "end":
		err = g.engine.EndTurn(g.gameState.ID, "Player")
		
	case "help":
		g.display.ShowCommands(g.gameState.Phase, true)
		g.input.WaitForEnter("Press Enter to continue...")
		
	case "quit":
		g.gameState.GameOver = true
		g.gameState.Winner = "AI"
		
	default:
		g.display.ShowMessage("Unknown command. Type 'help' for available commands.", ColorRed)
		g.input.WaitForEnter("Press Enter to continue...")
	}
	
	if err != nil {
		g.display.ShowError(err)
		g.input.WaitForEnter("")
	}
}

// handlePlayCommand handles the play card command
func (g *Game) handlePlayCommand(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("usage: play [card number]")
	}
	
	cardIndex, err := g.input.ParseCardIndex(args[0])
	if err != nil {
		return err
	}
	
	return g.engine.PlayCard(g.gameState.ID, "Player", cardIndex)
}

// handleAttackCommand handles the attack command
func (g *Game) handleAttackCommand(args []string) error {
	attackerIndex, targetIndex, err := g.input.ParseAttackTargets(args)
	if err != nil {
		return err
	}
	
	return g.engine.Attack(g.gameState.ID, "Player", attackerIndex, targetIndex)
}