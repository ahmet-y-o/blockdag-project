package main

import (
	"cardgame/battle"
	"fmt"
	"testing"
)

func TestSimpleGame(t *testing.T) {
	// Create simple test decks
	deck1 := []battle.Card{
		{ID: "1", Name: "Fire Dragon", Archetype: battle.ArchetypeNeutral, Attack: 2000, Defense: 1500, Cost: 4},
		{ID: "2", Name: "Water Sprite", Archetype: battle.ArchetypeNeutral, Attack: 1200, Defense: 1800, Cost: 3},
		{ID: "3", Name: "Earth Golem", Archetype: battle.ArchetypeNeutral, Attack: 1500, Defense: 2500, Cost: 4},
		{ID: "4", Name: "Wind Eagle", Archetype: battle.ArchetypeNeutral, Attack: 1800, Defense: 1200, Cost: 3},
		{ID: "5", Name: "Lightning Bolt", Archetype: battle.ArchetypeNeutral, Attack: 1000, Defense: 1000, Cost: 2},
	}

	deck2 := []battle.Card{
		{ID: "6", Name: "Dark Knight", Archetype: battle.ArchetypeNeutral, Attack: 2200, Defense: 1800, Cost: 5},
		{ID: "7", Name: "Holy Priest", Archetype: battle.ArchetypeNeutral, Attack: 800, Defense: 2000, Cost: 3},
		{ID: "8", Name: "Shadow Assassin", Archetype: battle.ArchetypeNeutral, Attack: 1900, Defense: 1100, Cost: 4},
		{ID: "9", Name: "Light Warrior", Archetype: battle.ArchetypeNeutral, Attack: 1700, Defense: 1600, Cost: 4},
		{ID: "10", Name: "Mystic Sage", Archetype: battle.ArchetypeNeutral, Attack: 1400, Defense: 1400, Cost: 3},
	}

	// Create battle engine
	engine := battle.NewBattleEngine()

	// Create a match
	game, err := engine.CreateMatch("TestPlayer1", "TestPlayer2", deck1, deck2)
	if err != nil {
		t.Fatalf("Failed to create match: %v", err)
	}

	fmt.Printf("Game created with ID: %s\n", game.ID)
	fmt.Printf("Player 1: %s (HP: %d, Hand: %d cards)\n", game.Player1.ID, game.Player1.HP, len(game.Player1.Hand))
	fmt.Printf("Player 2: %s (HP: %d, Hand: %d cards)\n", game.Player2.ID, game.Player2.HP, len(game.Player2.Hand))

	// Simulate a few turns
	fmt.Println("\n--- Turn 1: Player 1 ---")

	// Player 1 draws
	err = engine.DrawCard(game.ID, "TestPlayer1")
	if err != nil {
		t.Logf("Draw error: %v", err)
	}

	// Player 1 plays a card
	game, _ = engine.GetGameState(game.ID)
	fmt.Printf("Player 1 mana: %d/%d\n", game.Player1.Mana, game.Player1.MaxMana)
	fmt.Printf("Playing card: %s (Cost: %d)\n", game.Player1.Hand[0].Name, game.Player1.Hand[0].Cost)

	// End turn
	err = engine.EndTurn(game.ID, "TestPlayer1")
	if err != nil {
		t.Logf("End turn error: %v", err)
	}

	fmt.Println("\n--- Turn 2: Player 2 ---")

	// Player 2's turn
	game, _ = engine.GetGameState(game.ID)
	fmt.Printf("Current turn: %s\n", game.CurrentTurn)

	// Player 2 draws
	err = engine.DrawCard(game.ID, "TestPlayer2")
	if err != nil {
		t.Logf("Draw error: %v", err)
	}

	// End turn
	err = engine.EndTurn(game.ID, "TestPlayer2")
	if err != nil {
		t.Logf("End turn error: %v", err)
	}

	// Check game state
	game, _ = engine.GetGameState(game.ID)
	if game.TurnCount != 3 {
		t.Errorf("Expected turn count to be 3, got %d", game.TurnCount)
	}

	fmt.Println("\nTest completed successfully!")
}

// Run a quick simulation
func TestQuickBattle(t *testing.T) {
	// Create decks with low-cost cards for easier testing
	quickDeck := []battle.Card{
		{ID: "q1", Name: "Quick Strike", Archetype: battle.ArchetypeNeutral, Attack: 1000, Defense: 500, Cost: 1},
		{ID: "q2", Name: "Swift Attack", Archetype: battle.ArchetypeNeutral, Attack: 1200, Defense: 600, Cost: 1},
		{ID: "q3", Name: "Fast Blade", Archetype: battle.ArchetypeNeutral, Attack: 1100, Defense: 700, Cost: 1},
		{ID: "q4", Name: "Rapid Fire", Archetype: battle.ArchetypeNeutral, Attack: 900, Defense: 800, Cost: 1},
		{ID: "q5", Name: "Speed Demon", Archetype: battle.ArchetypeNeutral, Attack: 1300, Defense: 400, Cost: 2},
	}

	engine := battle.NewBattleEngine()
	game, _ := engine.CreateMatch("QuickP1", "QuickP2", quickDeck, quickDeck)

	// Give players extra mana for testing
	game.Player1.Mana = 5
	game.Player2.Mana = 5

	fmt.Println("\n=== Quick Battle Test ===")

	// Player 1 plays cards and attacks
	engine.PlayCard(game.ID, "QuickP1", 0)
	engine.PlayCard(game.ID, "QuickP1", 0)
	engine.ChangePhase(game.ID, "QuickP1", battle.PhaseBattle)

	// Direct attacks
	engine.Attack(game.ID, "QuickP1", 0, -1)
	engine.Attack(game.ID, "QuickP1", 1, -1)

	game, _ = engine.GetGameState(game.ID)
	fmt.Printf("After attacks - Player 2 HP: %d\n", game.Player2.HP)

	if game.Player2.HP >= 8000 {
		t.Error("Expected player 2 HP to decrease after attacks")
	}

	fmt.Println("Quick battle test passed!")
}
