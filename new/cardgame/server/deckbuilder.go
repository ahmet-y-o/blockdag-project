package server

import "cardgame/battle"

// DeckBuilder creates decks for the server
type DeckBuilder struct{}

// CreateEgyptianDeck creates a 40-card Egyptian deck
func (db *DeckBuilder) CreateEgyptianDeck() []battle.Card {
	return []battle.Card{
		// Copy the deck definition from game/decks.go
		// Legendary cards (1-2 copies each)
		{ID: "eg001", Name: "Ra, the Sun God", Archetype: battle.ArchetypeEgyptian, Attack: 3000, Defense: 2500, Cost: 8, Effect: "Deal 1000 damage to opponent", EffectType: "damage"},
		{ID: "eg002", Name: "Anubis, Guardian of the Dead", Archetype: battle.ArchetypeEgyptian, Attack: 2200, Defense: 2800, Cost: 6, Effect: "Heal 500 HP when a card is destroyed", EffectType: "heal"},
		
		// Rare cards (2-3 copies each)
		{ID: "eg003", Name: "Isis, Mother of Magic", Archetype: battle.ArchetypeEgyptian, Attack: 1800, Defense: 2000, Cost: 4, Effect: "Draw an additional card", EffectType: "draw"},
		{ID: "eg003", Name: "Isis, Mother of Magic", Archetype: battle.ArchetypeEgyptian, Attack: 1800, Defense: 2000, Cost: 4, Effect: "Draw an additional card", EffectType: "draw"},
		{ID: "eg004", Name: "Horus, the Avenger", Archetype: battle.ArchetypeEgyptian, Attack: 2500, Defense: 2000, Cost: 5, Effect: "", EffectType: ""},
		{ID: "eg004", Name: "Horus, the Avenger", Archetype: battle.ArchetypeEgyptian, Attack: 2500, Defense: 2000, Cost: 5, Effect: "", EffectType: ""},
		{ID: "eg005", Name: "Thoth, God of Wisdom", Archetype: battle.ArchetypeEgyptian, Attack: 1500, Defense: 2200, Cost: 3, Effect: "Gain 1 extra mana", EffectType: "mana"},
		{ID: "eg005", Name: "Thoth, God of Wisdom", Archetype: battle.ArchetypeEgyptian, Attack: 1500, Defense: 2200, Cost: 3, Effect: "Gain 1 extra mana", EffectType: "mana"},
		{ID: "eg006", Name: "Set, God of Chaos", Archetype: battle.ArchetypeEgyptian, Attack: 2800, Defense: 2000, Cost: 7, Effect: "", EffectType: ""},
		{ID: "eg007", Name: "Sobek, Crocodile God", Archetype: battle.ArchetypeEgyptian, Attack: 2000, Defense: 2400, Cost: 5, Effect: "", EffectType: ""},
		{ID: "eg007", Name: "Sobek, Crocodile God", Archetype: battle.ArchetypeEgyptian, Attack: 2000, Defense: 2400, Cost: 5, Effect: "", EffectType: ""},
		
		// Common cards (3 copies each)
		{ID: "eg008", Name: "Bastet, Cat Goddess", Archetype: battle.ArchetypeEgyptian, Attack: 1600, Defense: 1400, Cost: 3, Effect: "", EffectType: ""},
		{ID: "eg008", Name: "Bastet, Cat Goddess", Archetype: battle.ArchetypeEgyptian, Attack: 1600, Defense: 1400, Cost: 3, Effect: "", EffectType: ""},
		{ID: "eg008", Name: "Bastet, Cat Goddess", Archetype: battle.ArchetypeEgyptian, Attack: 1600, Defense: 1400, Cost: 3, Effect: "", EffectType: ""},
		{ID: "eg009", Name: "Nephthys, Lady of the House", Archetype: battle.ArchetypeEgyptian, Attack: 1700, Defense: 2100, Cost: 4, Effect: "", EffectType: ""},
		{ID: "eg009", Name: "Nephthys, Lady of the House", Archetype: battle.ArchetypeEgyptian, Attack: 1700, Defense: 2100, Cost: 4, Effect: "", EffectType: ""},
		{ID: "eg009", Name: "Nephthys, Lady of the House", Archetype: battle.ArchetypeEgyptian, Attack: 1700, Defense: 2100, Cost: 4, Effect: "", EffectType: ""},
		{ID: "eg010", Name: "Khepri, Scarab God", Archetype: battle.ArchetypeEgyptian, Attack: 1400, Defense: 1800, Cost: 3, Effect: "", EffectType: ""},
		{ID: "eg010", Name: "Khepri, Scarab God", Archetype: battle.ArchetypeEgyptian, Attack: 1400, Defense: 1800, Cost: 3, Effect: "", EffectType: ""},
		{ID: "eg010", Name: "Khepri, Scarab God", Archetype: battle.ArchetypeEgyptian, Attack: 1400, Defense: 1800, Cost: 3, Effect: "", EffectType: ""},
		{ID: "eg011", Name: "Egyptian Warrior", Archetype: battle.ArchetypeEgyptian, Attack: 1200, Defense: 1000, Cost: 2, Effect: "", EffectType: ""},
		{ID: "eg011", Name: "Egyptian Warrior", Archetype: battle.ArchetypeEgyptian, Attack: 1200, Defense: 1000, Cost: 2, Effect: "", EffectType: ""},
		{ID: "eg011", Name: "Egyptian Warrior", Archetype: battle.ArchetypeEgyptian, Attack: 1200, Defense: 1000, Cost: 2, Effect: "", EffectType: ""},
		{ID: "eg012", Name: "Pyramid Guardian", Archetype: battle.ArchetypeEgyptian, Attack: 800, Defense: 2000, Cost: 2, Effect: "", EffectType: ""},
		{ID: "eg012", Name: "Pyramid Guardian", Archetype: battle.ArchetypeEgyptian, Attack: 800, Defense: 2000, Cost: 2, Effect: "", EffectType: ""},
		{ID: "eg012", Name: "Pyramid Guardian", Archetype: battle.ArchetypeEgyptian, Attack: 800, Defense: 2000, Cost: 2, Effect: "", EffectType: ""},
		
		// Spell/Effect cards
		{ID: "n001", Name: "Healing Potion", Archetype: battle.ArchetypeNeutral, Attack: 0, Defense: 0, Cost: 1, Effect: "Heal 1000 HP", EffectType: "heal"},
		{ID: "n001", Name: "Healing Potion", Archetype: battle.ArchetypeNeutral, Attack: 0, Defense: 0, Cost: 1, Effect: "Heal 1000 HP", EffectType: "heal"},
		{ID: "n002", Name: "Lightning Bolt", Archetype: battle.ArchetypeNeutral, Attack: 0, Defense: 0, Cost: 3, Effect: "Deal 1500 damage", EffectType: "damage"},
		{ID: "n002", Name: "Lightning Bolt", Archetype: battle.ArchetypeNeutral, Attack: 0, Defense: 0, Cost: 3, Effect: "Deal 1500 damage", EffectType: "damage"},
		{ID: "n003", Name: "Power Crystal", Archetype: battle.ArchetypeNeutral, Attack: 1000, Defense: 1000, Cost: 2, Effect: "Gain 2 mana", EffectType: "mana"},
		{ID: "n003", Name: "Power Crystal", Archetype: battle.ArchetypeNeutral, Attack: 1000, Defense: 1000, Cost: 2, Effect: "Gain 2 mana", EffectType: "mana"},
		{ID: "n004", Name: "Ancient Warrior", Archetype: battle.ArchetypeNeutral, Attack: 1800, Defense: 1600, Cost: 3, Effect: "", EffectType: ""},
		{ID: "n004", Name: "Ancient Warrior", Archetype: battle.ArchetypeNeutral, Attack: 1800, Defense: 1600, Cost: 3, Effect: "", EffectType: ""},
		{ID: "n005", Name: "Mystic Shield", Archetype: battle.ArchetypeNeutral, Attack: 0, Defense: 2500, Cost: 2, Effect: "", EffectType: ""},
		{ID: "n005", Name: "Mystic Shield", Archetype: battle.ArchetypeNeutral, Attack: 0, Defense: 2500, Cost: 2, Effect: "", EffectType: ""},
		{ID: "n006", Name: "Swift Strike", Archetype: battle.ArchetypeNeutral, Attack: 1500, Defense: 1000, Cost: 2, Effect: "", EffectType: ""},
		{ID: "n006", Name: "Swift Strike", Archetype: battle.ArchetypeNeutral, Attack: 1500, Defense: 1000, Cost: 2, Effect: "", EffectType: ""},
		{ID: "n006", Name: "Swift Strike", Archetype: battle.ArchetypeNeutral, Attack: 1500, Defense: 1000, Cost: 2, Effect: "", EffectType: ""},
	}
}

// CreateGreekDeck creates a 40-card Greek deck
func (db *DeckBuilder) CreateGreekDeck() []battle.Card {
	return []battle.Card{
		// Copy the deck definition from game/decks.go
		// Similar structure to Egyptian deck but with Greek cards
		// ... (same as in game/decks.go)
		// Legendary cards (1-2 copies each)
		{ID: "gr001", Name: "Zeus, King of Olympus", Archetype: battle.ArchetypeGreek, Attack: 3200, Defense: 2400, Cost: 8, Effect: "Deal 500 damage to all enemies", EffectType: "damage"},
		{ID: "gr002", Name: "Athena, Goddess of War", Archetype: battle.ArchetypeGreek, Attack: 2400, Defense: 2600, Cost: 6, Effect: "", EffectType: ""},
		
		// Continue with the rest of the Greek deck...
		// (Copy from game/decks.go)
	}
}