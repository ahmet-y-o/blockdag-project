package game

import (
	"cardgame/battle"
	"time"
)

// AIPlayer handles AI decision making
type AIPlayer struct {
	difficulty string
}

// NewAIPlayer creates a new AI player
func NewAIPlayer(difficulty string) *AIPlayer {
	return &AIPlayer{difficulty: difficulty}
}

// MakeDecision makes AI decisions based on game state
func (ai *AIPlayer) MakeDecision(game *battle.GameState, engine *battle.BattleEngine, aiPlayerID string) {
	// Add thinking delay for better UX
	time.Sleep(AIThinkDelay)

	switch game.Phase {
	case battle.PhaseDrawn:
		ai.handleDrawPhase(game, engine, aiPlayerID)
	case battle.PhaseMain:
		ai.handleMainPhase(game, engine, aiPlayerID)
	case battle.PhaseBattle:
		ai.handleBattlePhase(game, engine, aiPlayerID)
	}
}

// handleDrawPhase handles AI decisions during draw phase
func (ai *AIPlayer) handleDrawPhase(game *battle.GameState, engine *battle.BattleEngine, aiPlayerID string) {
	// Always draw
	engine.DrawCard(game.ID, aiPlayerID)
}

// handleMainPhase handles AI decisions during main phase
func (ai *AIPlayer) handleMainPhase(game *battle.GameState, engine *battle.BattleEngine, aiPlayerID string) {
	aiPlayer := ai.getAIPlayer(game, aiPlayerID)

	// Try to play cards
	cardsPlayed := 0
	maxCardsToPlay := 2

	// Sort hand by strategy (play high-cost cards first if possible)
	playableCards := ai.getPlayableCards(aiPlayer)

	for _, cardIndex := range playableCards {
		if cardsPlayed >= maxCardsToPlay {
			break
		}

		if err := engine.PlayCard(game.ID, aiPlayerID, cardIndex); err == nil {
			cardsPlayed++
			time.Sleep(AIActionDelay)
			// Recalculate playable cards since hand changed
			aiPlayer = ai.getAIPlayer(game, aiPlayerID)
			playableCards = ai.getPlayableCards(aiPlayer)
		}
	}

	// Decide whether to enter battle phase
	if len(aiPlayer.Field) > 0 && ai.shouldEnterBattle(game, aiPlayerID) {
		engine.ChangePhase(game.ID, aiPlayerID, battle.PhaseBattle)
	} else {
		engine.EndTurn(game.ID, aiPlayerID)
	}
}

// handleBattlePhase handles AI decisions during battle phase
func (ai *AIPlayer) handleBattlePhase(game *battle.GameState, engine *battle.BattleEngine, aiPlayerID string) {
	aiPlayer := ai.getAIPlayer(game, aiPlayerID)
	opponent := ai.getOpponent(game, aiPlayerID)

	// Attack with all creatures
	for i := 0; i < len(aiPlayer.Field); i++ {
		attacker := aiPlayer.Field[i]

		if len(opponent.Field) > 0 {
			// Choose target based on strategy
			targetIndex := ai.chooseBattleTarget(attacker, opponent.Field)
			engine.Attack(game.ID, aiPlayerID, i, targetIndex)
		} else {
			// Direct attack
			engine.Attack(game.ID, aiPlayerID, i, -1)
		}

		time.Sleep(AIActionDelay)

		// Check if game ended
		game, _ = engine.GetGameState(game.ID)
		if game.GameOver {
			return
		}

		// Refresh player states
		aiPlayer = ai.getAIPlayer(game, aiPlayerID)
		opponent = ai.getOpponent(game, aiPlayerID)
	}

	// End turn after attacks
	engine.EndTurn(game.ID, aiPlayerID)
}

// getPlayableCards returns indices of cards that can be played
func (ai *AIPlayer) getPlayableCards(player *battle.Player) []int {
	var playable []int

	for i, card := range player.Hand {
		if card.Cost <= player.Mana {
			playable = append(playable, i)
		}
	}

	// Sort by cost (higher cost first for better plays)
	for i := 0; i < len(playable)-1; i++ {
		for j := i + 1; j < len(playable); j++ {
			if player.Hand[playable[i]].Cost < player.Hand[playable[j]].Cost {
				playable[i], playable[j] = playable[j], playable[i]
			}
		}
	}

	return playable
}

// shouldEnterBattle decides if AI should enter battle phase
func (ai *AIPlayer) shouldEnterBattle(game *battle.GameState, aiPlayerID string) bool {
	aiPlayer := ai.getAIPlayer(game, aiPlayerID)
	opponent := ai.getOpponent(game, aiPlayerID)

	// Simple strategy: enter battle if we have more creatures or opponent has low HP
	if len(opponent.Field) == 0 {
		return true // Can attack directly
	}

	if opponent.HP < 2000 {
		return true // Go for the kill
	}

	if len(aiPlayer.Field) >= len(opponent.Field) {
		return true // We have board advantage
	}

	// Calculate total attack power
	totalAttack := 0
	for _, card := range aiPlayer.Field {
		totalAttack += card.Attack
	}

	return totalAttack > 2000 // Attack if we have good damage potential
}

// chooseBattleTarget selects the best target for attack
func (ai *AIPlayer) chooseBattleTarget(attacker battle.Card, targets []battle.Card) int {
	// Strategy: Try to destroy cards we can kill
	for i, target := range targets {
		if attacker.Attack > target.Defense {
			return i // Can destroy this target
		}
	}

	// If we can't destroy anything, attack the weakest
	weakestIndex := 0
	weakestDefense := targets[0].Defense

	for i, target := range targets {
		if target.Defense < weakestDefense {
			weakestDefense = target.Defense
			weakestIndex = i
		}
	}

	return weakestIndex
}

// Helper functions
func (ai *AIPlayer) getAIPlayer(game *battle.GameState, aiPlayerID string) *battle.Player {
	if game.Player1.ID == aiPlayerID {
		return game.Player1
	}
	return game.Player2
}

func (ai *AIPlayer) getOpponent(game *battle.GameState, aiPlayerID string) *battle.Player {
	if game.Player1.ID == aiPlayerID {
		return game.Player2
	}
	return game.Player1
}
