package shared

import "cardgame/battle"

// Message types
const (
	// Client to Server
	MsgSetName      = "setName"
	MsgJoinQueue    = "joinQueue"
	MsgLeaveQueue   = "leaveQueue"
	MsgPlayCard     = "playCard"
	MsgAttack       = "attack"
	MsgEndTurn      = "endTurn"
	MsgChangePhase  = "changePhase"
	MsgDrawCard     = "drawCard"
	
	// Server to Client
	MsgWelcome              = "welcome"
	MsgNameSet              = "nameSet"
	MsgQueueJoined          = "queueJoined"
	MsgQueueLeft            = "queueLeft"
	MsgGameStart            = "gameStart"
	MsgGameUpdate           = "gameUpdate"
	MsgGameOver             = "gameOver"
	MsgError                = "error"
	MsgOpponentDisconnected = "opponentDisconnected"
)

// Message represents a network message
type Message struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}

// ConvertToGameState converts a map to GameState
func ConvertToGameState(data map[string]interface{}) *battle.GameState {
	// This is a simplified conversion - in production you'd want proper error handling
	state := &battle.GameState{
		ID:          data["id"].(string),
		CurrentTurn: data["current_turn"].(string),
		TurnCount:   int(data["turn_count"].(float64)),
		Phase:       battle.GamePhase(data["phase"].(string)),
		GameOver:    data["game_over"].(bool),
		Winner:      data["winner"].(string),
		LastAction:  data["last_action"].(string),
	}
	
	// Convert players
	if p1Data, ok := data["player1"].(map[string]interface{}); ok {
		state.Player1 = ConvertToPlayer(p1Data)
	}
	if p2Data, ok := data["player2"].(map[string]interface{}); ok {
		state.Player2 = ConvertToPlayer(p2Data)
	}
	
	return state
}

// ConvertToPlayer converts a map to Player
func ConvertToPlayer(data map[string]interface{}) *battle.Player {
	player := &battle.Player{
		ID:      data["id"].(string),
		HP:      int(data["hp"].(float64)),
		Mana:    int(data["mana"].(float64)),
		MaxMana: int(data["max_mana"].(float64)),
	}
	
	// Convert cards arrays
	if deckData, ok := data["deck"].([]interface{}); ok {
		player.Deck = ConvertToCards(deckData)
	}
	if handData, ok := data["hand"].([]interface{}); ok {
		player.Hand = ConvertToCards(handData)
	}
	if fieldData, ok := data["field"].([]interface{}); ok {
		player.Field = ConvertToCards(fieldData)
	}
	if graveyardData, ok := data["graveyard"].([]interface{}); ok {
		player.Graveyard = ConvertToCards(graveyardData)
	}
	
	return player
}

// ConvertToCards converts an array of card data to Card slice
func ConvertToCards(data []interface{}) []battle.Card {
	cards := make([]battle.Card, len(data))
	for i, cardData := range data {
		if cardMap, ok := cardData.(map[string]interface{}); ok {
			cards[i] = battle.Card{
				ID:         cardMap["id"].(string),
				Name:       cardMap["name"].(string),
				Archetype:  battle.Archetype(cardMap["archetype"].(string)),
				Attack:     int(cardMap["attack"].(float64)),
				Defense:    int(cardMap["defense"].(float64)),
				Cost:       int(cardMap["cost"].(float64)),
				Effect:     cardMap["effect"].(string),
				EffectType: cardMap["effect_type"].(string),
			}
		}
	}
	return cards
}