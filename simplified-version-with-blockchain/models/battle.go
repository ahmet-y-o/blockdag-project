package models

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

type BattleState string

const (
	BattleStateInit    BattleState = "INIT"
	BattleStatePlayer1 BattleState = "PLAYER1_TURN"
	BattleStatePlayer2 BattleState = "PLAYER2_TURN"
	BattleStateBattle  BattleState = "BATTLE_PHASE"
	BattleStateEnd     BattleState = "END"
)

type Battle struct {
	ID           string
	Player1      *Player
	Player2      *Player
	CurrentState BattleState
	Round        int
	PlayedCards  map[string][]Card // Maps player ID to their played cards
	CreatedAt    time.Time
	Winner       *Player
}

func NewBattle(player1, player2 *Player) *Battle {
	battle := &Battle{
		ID:           uuid.New().String(),
		Player1:      player1,
		Player2:      player2,
		CurrentState: BattleStateInit,
		Round:        1,
		PlayedCards:  make(map[string][]Card),
		CreatedAt:    time.Now(),
	}

	// Initialize played cards slices
	battle.PlayedCards[fmt.Sprintf("%d", player1.ID)] = make([]Card, 0)
	battle.PlayedCards[fmt.Sprintf("%d", player2.ID)] = make([]Card, 0)

	return battle
}

func (b *Battle) Start() error {
	// Initialize players
	b.Player1.ResetForNewGame()
	b.Player2.ResetForNewGame()

	// Set initial mana
	b.Player1.MaxMana = 1
	b.Player1.Mana = 1
	b.Player2.MaxMana = 1
	b.Player2.Mana = 1

	// Draw initial hands
	for i := 0; i < 5; i++ {
		if _, err := b.Player1.DrawCard(); err != nil {
			return fmt.Errorf("failed to draw card for player 1: %v", err)
		}
		if _, err := b.Player2.DrawCard(); err != nil {
			return fmt.Errorf("failed to draw card for player 2: %v", err)
		}
	}

	b.CurrentState = BattleStatePlayer1
	return nil
}

func (b *Battle) StartBattlePhase() error {
	if b.CurrentState == BattleStateBattle {
		return fmt.Errorf("battle phase already in progress")
	}

	// Calculate damage for both players
	p1Cards := b.PlayedCards[fmt.Sprintf("%d", b.Player1.ID)]
	p2Cards := b.PlayedCards[fmt.Sprintf("%d", b.Player2.ID)]

	p1Damage := calculateTotalDamage(p1Cards)
	p2Damage := calculateTotalDamage(p2Cards)

	// Determine round winner
	if p1Damage > p2Damage {
		b.Player1.Experience += 10
	} else if p2Damage > p1Damage {
		b.Player2.Experience += 10
	} else {
		// Draw - both players get some experience
		b.Player1.Experience += 5
		b.Player2.Experience += 5
	}

	// Clear played cards
	b.PlayedCards[fmt.Sprintf("%d", b.Player1.ID)] = make([]Card, 0)
	b.PlayedCards[fmt.Sprintf("%d", b.Player2.ID)] = make([]Card, 0)

	// Start new round
	b.Round++

	// Increment max mana (up to 10)
	if b.Player1.MaxMana < 10 {
		b.Player1.MaxMana++
	}
	if b.Player2.MaxMana < 10 {
		b.Player2.MaxMana++
	}

	// Refill mana
	b.Player1.RefillMana()
	b.Player2.RefillMana()

	// Draw new cards
	if _, err := b.Player1.DrawCard(); err != nil {
		b.checkGameOver()
		return err
	}
	if _, err := b.Player2.DrawCard(); err != nil {
		b.checkGameOver()
		return err
	}

	// Check if game should end
	if b.checkGameOver() {
		return nil
	}

	// Continue to next round
	b.CurrentState = BattleStatePlayer1
	return nil
}

func (b *Battle) checkGameOver() bool {
	// Check if either player is out of cards (both deck and hand)
	p1NoCards := len(b.Player1.Deck) == 0 && len(b.Player1.Hand) == 0
	p2NoCards := len(b.Player2.Deck) == 0 && len(b.Player2.Hand) == 0

	if p1NoCards || p2NoCards {
		b.CurrentState = BattleStateEnd

		// Determine winner based on experience points
		if b.Player1.Experience > b.Player2.Experience {
			b.Winner = b.Player1
		} else if b.Player2.Experience > b.Player1.Experience {
			b.Winner = b.Player2
		}
		// If equal experience, game ends in a draw (Winner remains nil)

		return true
	}

	return false
}

func calculateTotalDamage(cards []Card) int {
	totalDamage := 0
	for _, card := range cards {
		totalDamage += card.GetDamage()
	}
	return totalDamage
}

// Helper method to check if a player can play
func (b *Battle) CanPlayerPlay(playerID int) bool {
	if b.CurrentState != BattleStatePlayer1 && b.CurrentState != BattleStatePlayer2 {
		return false
	}

	if b.CurrentState == BattleStatePlayer1 && playerID == b.Player1.ID {
		return true
	}

	if b.CurrentState == BattleStatePlayer2 && playerID == b.Player2.ID {
		return true
	}

	return false
}

// Get current player
func (b *Battle) GetCurrentPlayer() *Player {
	if b.CurrentState == BattleStatePlayer1 {
		return b.Player1
	}
	if b.CurrentState == BattleStatePlayer2 {
		return b.Player2
	}
	return nil
}
