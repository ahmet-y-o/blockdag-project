package battle

import (
	//"encoding/json"
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// Card represents a game card
type Card struct {
	ID         string    `json:"id"`
	Name       string    `json:"name"`
	Archetype  Archetype `json:"archetype"`
	Attack     int       `json:"attack"`
	Defense    int       `json:"defense"`
	Cost       int       `json:"cost"`
	Effect     string    `json:"effect"`
	EffectType string    `json:"effect_type"`
}

// Archetype represents card archetypes
type Archetype string

const (
	ArchetypeEgyptian Archetype = "egyptian"
	ArchetypeGreek    Archetype = "greek"
	ArchetypeNeutral  Archetype = "neutral"
)

// Player represents a player in the game
type Player struct {
	ID             string                `json:"id"`
	Name           string                `json:"name"`
	HP             int                   `json:"hp"`
	Mana           int                   `json:"mana"`
	MaxMana        int                   `json:"max_mana"`
	Deck           []Card                `json:"deck"`
	Hand           []Card                `json:"hand"`
	Field          []Card                `json:"field"`
	Graveyard      []Card                `json:"graveyard"`
	ArchetypeBonus map[Archetype]float32 `json:"archetype_bonus"`
}

// GameState represents the current state of the game
type GameState struct {
	ID          string    `json:"id"`
	Player1     *Player   `json:"player1"`
	Player2     *Player   `json:"player2"`
	CurrentTurn string    `json:"current_turn"`
	TurnCount   int       `json:"turn_count"`
	Phase       GamePhase `json:"phase"`
	Winner      string    `json:"winner"`
	GameOver    bool      `json:"game_over"`
	LastAction  string    `json:"last_action"`
}

// GamePhase represents different phases of a turn
type GamePhase string

const (
	PhaseDrawn  GamePhase = "draw"
	PhaseMain   GamePhase = "main"
	PhaseBattle GamePhase = "battle"
	PhaseEnd    GamePhase = "end"
)

// BattleEngine manages the game logic
type BattleEngine struct {
	games     map[string]*GameState
	mu        sync.RWMutex
	callbacks map[string]func(*GameState)
}

// NewBattleEngine creates a new battle engine instance
func NewBattleEngine() *BattleEngine {
	return &BattleEngine{
		games:     make(map[string]*GameState),
		callbacks: make(map[string]func(*GameState)),
	}
}

// CreateMatch creates a new match between two players
func (be *BattleEngine) CreateMatch(player1ID, player2ID string, deck1, deck2 []Card) (*GameState, error) {
	be.mu.Lock()
	defer be.mu.Unlock()

	gameID := fmt.Sprintf("game_%d", time.Now().Unix())

	// Initialize players
	p1 := &Player{
		ID:             player1ID,
		HP:             8000,
		Mana:           10,
		MaxMana:        15,
		Deck:           shuffleDeck(deck1),
		Hand:           []Card{},
		Field:          []Card{},
		Graveyard:      []Card{},
		ArchetypeBonus: make(map[Archetype]float32),
	}

	p2 := &Player{
		ID:             player2ID,
		HP:             8000,
		Mana:           10,
		MaxMana:        15,
		Deck:           shuffleDeck(deck2),
		Hand:           []Card{},
		Field:          []Card{},
		Graveyard:      []Card{},
		ArchetypeBonus: make(map[Archetype]float32),
	}

	// Draw initial hands
	drawCards(p1, 5)
	drawCards(p2, 5)

	game := &GameState{
		ID:          gameID,
		Player1:     p1,
		Player2:     p2,
		CurrentTurn: player1ID,
		TurnCount:   1,
		Phase:       PhaseMain,
		GameOver:    false,
	}

	be.games[gameID] = game
	be.notifyStateChange(game)

	return game, nil
}

// DrawCard handles drawing a card for the current player
func (be *BattleEngine) DrawCard(gameID, playerID string) error {
	be.mu.Lock()
	defer be.mu.Unlock()

	game, exists := be.games[gameID]
	if !exists {
		return fmt.Errorf("game not found")
	}

	if game.CurrentTurn != playerID {
		return fmt.Errorf("not your turn")
	}

	if game.Phase != PhaseDrawn {
		return fmt.Errorf("can only draw during draw phase")
	}

	player := be.getPlayer(game, playerID)
	if player == nil {
		return fmt.Errorf("player not found")
	}

	if drawCards(player, 1) {
		game.Phase = PhaseMain
		be.notifyStateChange(game)
	}

	return nil
}

// PlayCard plays a card from hand to field
func (be *BattleEngine) PlayCard(gameID, playerID string, cardIndex int) error {
	be.mu.Lock()
	defer be.mu.Unlock()

	game, exists := be.games[gameID]
	if !exists {
		return fmt.Errorf("game not found")
	}

	if game.CurrentTurn != playerID {
		return fmt.Errorf("not your turn")
	}

	if game.Phase != PhaseMain {
		return fmt.Errorf("can only play cards during main phase")
	}

	player := be.getPlayer(game, playerID)
	if player == nil {
		return fmt.Errorf("player not found")
	}

	if cardIndex < 0 || cardIndex >= len(player.Hand) {
		return fmt.Errorf("invalid card index")
	}

	card := player.Hand[cardIndex]

	// Check mana cost
	if card.Cost > player.Mana {
		return fmt.Errorf("insufficient mana")
	}

	// Apply archetype bonuses
	be.applyArchetypeBonus(player, &card)

	// Move card from hand to field
	player.Hand = append(player.Hand[:cardIndex], player.Hand[cardIndex+1:]...)
	player.Field = append(player.Field, card)
	player.Mana -= card.Cost

	// Apply card effects
	be.applyCardEffect(game, player, card)

	game.LastAction = fmt.Sprintf("%s played %s", playerID, card.Name)
	be.notifyStateChange(game)

	return nil
}

// Attack executes an attack with a card
func (be *BattleEngine) Attack(gameID, playerID string, attackerIndex, targetIndex int) error {
	be.mu.Lock()
	defer be.mu.Unlock()

	game, exists := be.games[gameID]
	if !exists {
		return fmt.Errorf("game not found")
	}

	if game.CurrentTurn != playerID {
		return fmt.Errorf("not your turn")
	}

	if game.Phase != PhaseBattle {
		return fmt.Errorf("can only attack during battle phase")
	}

	attacker := be.getPlayer(game, playerID)
	defender := be.getOpponent(game, playerID)

	if attackerIndex < 0 || attackerIndex >= len(attacker.Field) {
		return fmt.Errorf("invalid attacker index")
	}

	attackCard := attacker.Field[attackerIndex]

	// Direct attack to player
	if targetIndex == -1 {
		if len(defender.Field) > 0 {
			return fmt.Errorf("cannot attack directly when opponent has cards")
		}
		defender.HP -= attackCard.Attack
		game.LastAction = fmt.Sprintf("%s attacked directly for %d damage", attackCard.Name, attackCard.Attack)
	} else {
		// Attack a card
		if targetIndex < 0 || targetIndex >= len(defender.Field) {
			return fmt.Errorf("invalid target index")
		}

		targetCard := defender.Field[targetIndex]

		// Battle calculation
		if attackCard.Attack > targetCard.Defense {
			// Destroy target card
			defender.Field = append(defender.Field[:targetIndex], defender.Field[targetIndex+1:]...)
			defender.Graveyard = append(defender.Graveyard, targetCard)
			game.LastAction = fmt.Sprintf("%s destroyed %s", attackCard.Name, targetCard.Name)
		} else if attackCard.Attack < targetCard.Defense {
			// Destroy attacker
			attacker.Field = append(attacker.Field[:attackerIndex], attacker.Field[attackerIndex+1:]...)
			attacker.Graveyard = append(attacker.Graveyard, attackCard)
			game.LastAction = fmt.Sprintf("%s was destroyed by %s", attackCard.Name, targetCard.Name)
		} else {
			// Both destroyed
			attacker.Field = append(attacker.Field[:attackerIndex], attacker.Field[attackerIndex+1:]...)
			attacker.Graveyard = append(attacker.Graveyard, attackCard)
			defender.Field = append(defender.Field[:targetIndex], defender.Field[targetIndex+1:]...)
			defender.Graveyard = append(defender.Graveyard, targetCard)
			game.LastAction = "Both cards destroyed"
		}
	}

	// Check win condition
	if defender.HP <= 0 {
		game.GameOver = true
		game.Winner = playerID
	}

	be.notifyStateChange(game)
	return nil
}

// EndTurn ends the current player's turn
func (be *BattleEngine) EndTurn(gameID, playerID string) error {
	be.mu.Lock()
	defer be.mu.Unlock()

	game, exists := be.games[gameID]
	if !exists {
		return fmt.Errorf("game not found")
	}

	if game.CurrentTurn != playerID {
		return fmt.Errorf("not your turn")
	}

	// Switch turn
	if game.CurrentTurn == game.Player1.ID {
		game.CurrentTurn = game.Player2.ID
	} else {
		game.CurrentTurn = game.Player1.ID
	}

	// Increment turn count
	game.TurnCount++

	// Reset phase
	game.Phase = PhaseDrawn

	// Increase mana for new turn player
	currentPlayer := be.getPlayer(game, game.CurrentTurn)
	if currentPlayer.MaxMana < 10 {
		currentPlayer.MaxMana++
	}
	currentPlayer.Mana = currentPlayer.MaxMana

	game.LastAction = fmt.Sprintf("%s ended turn", playerID)
	be.notifyStateChange(game)

	return nil
}

// ChangePhase changes the game phase
func (be *BattleEngine) ChangePhase(gameID, playerID string, phase GamePhase) error {
	be.mu.Lock()
	defer be.mu.Unlock()

	game, exists := be.games[gameID]
	if !exists {
		return fmt.Errorf("game not found")
	}

	if game.CurrentTurn != playerID {
		return fmt.Errorf("not your turn")
	}

	game.Phase = phase
	be.notifyStateChange(game)

	return nil
}

// GetGameState returns the current game state
func (be *BattleEngine) GetGameState(gameID string) (*GameState, error) {
	be.mu.RLock()
	defer be.mu.RUnlock()

	game, exists := be.games[gameID]
	if !exists {
		return nil, fmt.Errorf("game not found")
	}

	return game, nil
}

// RegisterCallback registers a callback for state changes
func (be *BattleEngine) RegisterCallback(gameID string, callback func(*GameState)) {
	be.mu.Lock()
	defer be.mu.Unlock()

	be.callbacks[gameID] = callback
}

// Helper functions

func (be *BattleEngine) getPlayer(game *GameState, playerID string) *Player {
	if game.Player1.ID == playerID {
		return game.Player1
	} else if game.Player2.ID == playerID {
		return game.Player2
	}
	return nil
}

func (be *BattleEngine) getOpponent(game *GameState, playerID string) *Player {
	if game.Player1.ID == playerID {
		return game.Player2
	} else if game.Player2.ID == playerID {
		return game.Player1
	}
	return nil
}

func (be *BattleEngine) notifyStateChange(game *GameState) {
	if callback, exists := be.callbacks[game.ID]; exists {
		go callback(game)
	}
}

func (be *BattleEngine) applyArchetypeBonus(player *Player, card *Card) {
	// Count archetype cards on field
	archetypeCount := make(map[Archetype]int)
	for _, fieldCard := range player.Field {
		archetypeCount[fieldCard.Archetype]++
	}

	// Egyptian bonus: +10% attack for each Egyptian card
	if card.Archetype == ArchetypeEgyptian {
		bonus := float32(archetypeCount[ArchetypeEgyptian]) * 0.1
		card.Attack = int(float32(card.Attack) * (1 + bonus))
	}

	// Greek bonus: +10% defense for each Greek card
	if card.Archetype == ArchetypeGreek {
		bonus := float32(archetypeCount[ArchetypeGreek]) * 0.1
		card.Defense = int(float32(card.Defense) * (1 + bonus))
	}
}

func (be *BattleEngine) applyCardEffect(game *GameState, player *Player, card Card) {
	switch card.EffectType {
	case "draw":
		drawCards(player, 1)
	case "damage":
		opponent := be.getOpponent(game, player.ID)
		opponent.HP -= 500
	case "heal":
		player.HP += 500
		if player.HP > 8000 {
			player.HP = 8000
		}
	case "mana":
		player.Mana += 1
	}
}

func shuffleDeck(deck []Card) []Card {
	shuffled := make([]Card, len(deck))
	copy(shuffled, deck)
	rand.Shuffle(len(shuffled), func(i, j int) {
		shuffled[i], shuffled[j] = shuffled[j], shuffled[i]
	})
	return shuffled
}

func drawCards(player *Player, count int) bool {
	for i := 0; i < count; i++ {
		if len(player.Deck) == 0 {
			return false // Deck out condition
		}
		player.Hand = append(player.Hand, player.Deck[0])
		player.Deck = player.Deck[1:]
	}
	return true
}

// API Response type
type APIResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}
