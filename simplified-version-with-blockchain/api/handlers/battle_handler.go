package handlers

import (
	"cardgame/models"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

type BattleHandler struct {
	ActiveBattles map[string]*models.Battle
}

// Request/Response structures
type CreateBattleRequest struct {
	Player1ID int `json:"player1_id"`
	Player2ID int `json:"player2_id"`
}

type BattleResponse struct {
	BattleID     string                `json:"battle_id"`
	CurrentState models.BattleState    `json:"current_state"`
	Round        int                   `json:"round"`
	Player1Info  PlayerInfo            `json:"player1"`
	Player2Info  PlayerInfo            `json:"player2"`
	PlayedCards  map[string][]CardInfo `json:"played_cards,omitempty"`
	Winner       *PlayerInfo           `json:"winner,omitempty"`
	IsGameOver   bool                  `json:"is_game_over"`
}

type PlayerInfo struct {
	ID       int        `json:"id"`
	Mana     int        `json:"mana"`
	MaxMana  int        `json:"max_mana"`
	HandSize int        `json:"hand_size"`
	DeckSize int        `json:"deck_size"`
	Hand     []CardInfo `json:"hand,omitempty"` // Only shown to the owner
}

type CardInfo struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Type     string `json:"type"`
	Damage   int    `json:"damage"`
	ManaCost int    `json:"mana_cost"`
}

type PlayCardRequest struct {
	PlayerID  int `json:"player_id"`
	CardIndex int `json:"card_index"`
}

func NewBattleHandler() *BattleHandler {
	return &BattleHandler{
		ActiveBattles: make(map[string]*models.Battle),
	}
}

// Helper method to create battle response
func (bh *BattleHandler) createBattleResponse(battle *models.Battle, requestingPlayerID int) BattleResponse {
	response := BattleResponse{
		BattleID:     battle.ID,
		CurrentState: battle.CurrentState,
		Round:        battle.Round,
		Player1Info:  bh.createPlayerInfo(battle.Player1, requestingPlayerID == battle.Player1.ID),
		Player2Info:  bh.createPlayerInfo(battle.Player2, requestingPlayerID == battle.Player2.ID),
		PlayedCards:  make(map[string][]CardInfo),
		IsGameOver:   battle.Winner != nil,
	}

	// Add played cards information
	for playerID, cards := range battle.PlayedCards {
		cardInfos := make([]CardInfo, len(cards))
		for i, card := range cards {
			cardInfos[i] = bh.createCardInfo(card)
		}
		response.PlayedCards[playerID] = cardInfos
	}

	// Add winner information if game is over
	if battle.Winner != nil {
		winnerInfo := bh.createPlayerInfo(battle.Winner, true)
		response.Winner = &winnerInfo
	}

	return response
}

func (bh *BattleHandler) createPlayerInfo(player *models.Player, includeHand bool) PlayerInfo {
	info := PlayerInfo{
		ID:       player.ID,
		Mana:     player.Mana,
		MaxMana:  player.MaxMana,
		HandSize: len(player.Hand),
		DeckSize: len(player.Deck),
	}

	if includeHand {
		info.Hand = make([]CardInfo, len(player.Hand))
		for i, card := range player.Hand {
			info.Hand[i] = bh.createCardInfo(card)
		}
	}

	return info
}

func (bh *BattleHandler) createCardInfo(card models.Card) CardInfo {
	return CardInfo{
		ID:       card.ID,
		Name:     card.GetName(),
		Type:     card.GetType(),
		Damage:   card.GetDamage(),
		ManaCost: card.GetManaCost(),
	}
}

func (bh *BattleHandler) PlayCard(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	battleID := vars["battleID"] // Changed from r.URL.Query().Get("battle_id")

	battle, exists := bh.ActiveBattles[battleID]
	if !exists {
		http.Error(w, `{"error": "Battle not found"}`, http.StatusNotFound)
		return
	}

	var req PlayCardRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error": "Invalid request body"}`, http.StatusBadRequest)
		return
	}

	if err := bh.validateAndPlayCard(battle, req); err != nil {
		http.Error(w, fmt.Sprintf(`{"error": "%s"}`, err.Error()), http.StatusBadRequest)
		return
	}

	response := bh.createBattleResponse(battle, req.PlayerID)
	json.NewEncoder(w).Encode(response)
}

func (bh *BattleHandler) GetBattleState(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	battleID := vars["battleID"] // Changed from r.URL.Query().Get("battle_id")
	playerID := r.URL.Query().Get("player_id")

	battle, exists := bh.ActiveBattles[battleID]
	if !exists {
		http.Error(w, `{"error": "Battle not found"}`, http.StatusNotFound)
		return
	}

	var pid int
	fmt.Sscanf(playerID, "%d", &pid)

	response := bh.createBattleResponse(battle, pid)
	json.NewEncoder(w).Encode(response)
}

func (bh *BattleHandler) validateAndPlayCard(battle *models.Battle, req PlayCardRequest) error {
	// Validate battle state
	if battle.CurrentState == models.BattleStateBattle ||
		battle.CurrentState == models.BattleStateEnd {
		return fmt.Errorf("cannot play cards in current state: %s", battle.CurrentState)
	}

	// Validate current player's turn
	if (battle.CurrentState == models.BattleStatePlayer1 && req.PlayerID != battle.Player1.ID) ||
		(battle.CurrentState == models.BattleStatePlayer2 && req.PlayerID != battle.Player2.ID) {
		return fmt.Errorf("not player's turn")
	}

	player := battle.Player1
	if req.PlayerID == battle.Player2.ID {
		player = battle.Player2
	}

	// Play the card
	card, err := player.PlayCard(req.CardIndex)
	if err != nil {
		return err
	}

	// Record played card
	battle.PlayedCards[fmt.Sprintf("%d", player.ID)] = append(
		battle.PlayedCards[fmt.Sprintf("%d", player.ID)],
		card,
	)

	// Update battle state
	if player.ID == battle.Player1.ID {
		battle.CurrentState = models.BattleStatePlayer2
	} else {
		battle.CurrentState = models.BattleStatePlayer1
	}

	return nil
}

func (bh *BattleHandler) shouldStartBattlePhase(battle *models.Battle) bool {
	// Check if both players have played their cards
	return len(battle.PlayedCards[fmt.Sprintf("%d", battle.Player1.ID)]) > 0 &&
		len(battle.PlayedCards[fmt.Sprintf("%d", battle.Player2.ID)]) > 0
}

// New methods for battle management
func (bh *BattleHandler) EndBattle(battleID string) error {
	delete(bh.ActiveBattles, battleID)
	return nil
}

func (bh *BattleHandler) CleanupInactiveBattles(maxAge time.Duration) {
	now := time.Now()
	for id, battle := range bh.ActiveBattles {
		if now.Sub(battle.CreatedAt) > maxAge {
			delete(bh.ActiveBattles, id)
		}
	}
}

func (bh *BattleHandler) CreateBattle(w http.ResponseWriter, r *http.Request) {
	var req CreateBattleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, fmt.Sprintf(`{"error": "%s"}`, err.Error()), http.StatusBadRequest)
		return
	}

	// Validate player IDs
	if req.Player1ID == req.Player2ID {
		http.Error(w, `{"error": "Players must be different"}`, http.StatusBadRequest)
		return
	}

	player1 := models.NewPlayer(req.Player1ID)
	player2 := models.NewPlayer(req.Player2ID)

	// Add cards to players' decks
	if err := addInitialCards(player1); err != nil {
		http.Error(w, fmt.Sprintf(`{"error": "Failed to initialize player 1 deck: %v"}`, err), http.StatusInternalServerError)
		return
	}
	if err := addInitialCards(player2); err != nil {
		http.Error(w, fmt.Sprintf(`{"error": "Failed to initialize player 2 deck: %v"}`, err), http.StatusInternalServerError)
		return
	}

	battle := models.NewBattle(player1, player2)
	if err := battle.Start(); err != nil {
		http.Error(w, fmt.Sprintf(`{"error": "%s"}`, err.Error()), http.StatusInternalServerError)
		return
	}

	bh.ActiveBattles[battle.ID] = battle

	response := bh.createBattleResponse(battle, req.Player1ID)
	json.NewEncoder(w).Encode(response)
}

// Helper function to add initial cards to a player's deck
func addInitialCards(player *models.Player) error {
	// List of all available cards
	cardTypes := []string{
		"FIRE_SWORD", "INFERNO_BLADE",
		"ICE_DAGGER", "FROST_SPEAR",
		"STONE_FIST", "EARTH_HAMMER",
		"WIND_SLASH", "TORNADO_BLADE",
		"THUNDER_STRIKE", "LIGHTNING_AXE",
	}

	// Add 3 copies of each card to make a full deck
	for _, cardID := range cardTypes {
		for i := 0; i < 3; i++ {
			card, err := models.NewCard(cardID)
			if err != nil {
				return fmt.Errorf("failed to create card %s: %v", cardID, err)
			}
			player.AddCard(*card)
		}
	}

	fmt.Printf("Player %d deck size: %d\n", player.ID, len(player.Deck))
	player.ShuffleDeck()
	return nil
}
