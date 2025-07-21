package models

import (
	"fmt"
	"math/rand"
	"time"
)

type Player struct {
	ID         int
	Experience float64
	Deck       []Card
	Hand       []Card
	Money      float64
	Mana       int
	MaxMana    int
	Wins       []Game
	Losses     []Game
}

type Game struct {
	GameID     int
	Date       time.Time
	OpponentID int
	Score      int
}

// Constructor function to create a new player
func NewPlayer(id int) *Player {
	return &Player{
		ID:         id,
		Experience: 0,
		Deck:       make([]Card, 0),
		Hand:       make([]Card, 0), // Initialize hand
		Money:      0,
		Mana:       0,
		MaxMana:    0,
		Wins:       make([]Game, 0),
		Losses:     make([]Game, 0),
	}
}

// Basic player methods
func (p *Player) AddCard(card Card) {
	p.Deck = append(p.Deck, card)
}

func (p *Player) AddWin(game Game) {
	p.Wins = append(p.Wins, game)
}

func (p *Player) AddLoss(game Game) {
	p.Losses = append(p.Losses, game)
}

func (p *Player) GetWinCount() int {
	return len(p.Wins)
}

func (p *Player) GetLossCount() int {
	return len(p.Losses)
}

// Mana management methods
func (p *Player) RefillMana() {
	p.Mana = p.MaxMana
}

func (p *Player) IncrementMaxMana(amount int) {
	p.MaxMana += amount
	p.Mana = p.MaxMana
}

// Card management methods
func (p *Player) DrawCard() (Card, error) {
	if len(p.Deck) == 0 {
		return Card{}, fmt.Errorf("deck is empty")
	}

	card := p.Deck[0]
	p.Deck = p.Deck[1:]
	p.Hand = append(p.Hand, card)
	return card, nil
}

func (p *Player) PlayCard(cardIndex int) (Card, error) {
	if cardIndex >= len(p.Hand) {
		return Card{}, fmt.Errorf("invalid card index")
	}

	card := p.Hand[cardIndex]
	if card.GetManaCost() > p.Mana {
		return Card{}, fmt.Errorf("not enough mana")
	}

	// Remove the card from hand and subtract mana
	p.Hand = append(p.Hand[:cardIndex], p.Hand[cardIndex+1:]...)
	p.Mana -= card.GetManaCost()

	return card, nil
}

func (p *Player) ShuffleDeck() {
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(p.Deck), func(i, j int) {
		p.Deck[i], p.Deck[j] = p.Deck[j], p.Deck[i]
	})
}

// Utility methods
func (p *Player) GetHandSize() int {
	return len(p.Hand)
}

func (p *Player) GetDeckSize() int {
	return len(p.Deck)
}

// Check if player can play a specific card
func (p *Player) CanPlayCard(cardIndex int) bool {
	if cardIndex >= len(p.Hand) {
		return false
	}
	return p.Hand[cardIndex].GetManaCost() <= p.Mana
}

// Reset player for new game
func (p *Player) ResetForNewGame() {
	p.Hand = make([]Card, 0)
	p.Mana = 0
	p.MaxMana = 0
	p.ShuffleDeck()
}
