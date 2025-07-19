package game

import "cardgame/battle"

// DeckBuilder creates decks
type DeckBuilder struct{}

// NewDeckBuilder creates a new deck builder
func NewDeckBuilder() *DeckBuilder {
	return &DeckBuilder{}
}

// CreateEgyptianDeck creates a 40-card Egyptian-themed deck
func (db *DeckBuilder) CreateEgyptianDeck() []battle.Card {
	return []battle.Card{
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

// CreateGreekDeck creates a 40-card Greek-themed deck
func (db *DeckBuilder) CreateGreekDeck() []battle.Card {
	return []battle.Card{
		// Legendary cards (1-2 copies each)
		{ID: "gr001", Name: "Zeus, King of Olympus", Archetype: battle.ArchetypeGreek, Attack: 3200, Defense: 2400, Cost: 8, Effect: "Deal 500 damage to all enemies", EffectType: "damage"},
		{ID: "gr002", Name: "Athena, Goddess of War", Archetype: battle.ArchetypeGreek, Attack: 2400, Defense: 2600, Cost: 6, Effect: "", EffectType: ""},
		
		// Rare cards (2-3 copies each)
		{ID: "gr003", Name: "Poseidon, Lord of the Seas", Archetype: battle.ArchetypeGreek, Attack: 2800, Defense: 2200, Cost: 7, Effect: "", EffectType: ""},
		{ID: "gr004", Name: "Apollo, God of Light", Archetype: battle.ArchetypeGreek, Attack: 2000, Defense: 2000, Cost: 4, Effect: "Heal 1000 HP", EffectType: "heal"},
		{ID: "gr004", Name: "Apollo, God of Light", Archetype: battle.ArchetypeGreek, Attack: 2000, Defense: 2000, Cost: 4, Effect: "Heal 1000 HP", EffectType: "heal"},
		{ID: "gr005", Name: "Hermes, the Messenger", Archetype: battle.ArchetypeGreek, Attack: 1600, Defense: 1800, Cost: 3, Effect: "Draw 2 cards", EffectType: "draw"},
		{ID: "gr005", Name: "Hermes, the Messenger", Archetype: battle.ArchetypeGreek, Attack: 1600, Defense: 1800, Cost: 3, Effect: "Draw 2 cards", EffectType: "draw"},
		{ID: "gr006", Name: "Ares, God of War", Archetype: battle.ArchetypeGreek, Attack: 2600, Defense: 1800, Cost: 6, Effect: "", EffectType: ""},
		{ID: "gr007", Name: "Hera, Queen of Gods", Archetype: battle.ArchetypeGreek, Attack: 2000, Defense: 2500, Cost: 5, Effect: "", EffectType: ""},
		{ID: "gr007", Name: "Hera, Queen of Gods", Archetype: battle.ArchetypeGreek, Attack: 2000, Defense: 2500, Cost: 5, Effect: "", EffectType: ""},
		{ID: "gr008", Name: "Demeter, Goddess of Harvest", Archetype: battle.ArchetypeGreek, Attack: 1500, Defense: 2300, Cost: 4, Effect: "Gain 2 mana", EffectType: "mana"},
		
		// Common cards (3 copies each)
		{ID: "gr009", Name: "Artemis, the Hunter", Archetype: battle.ArchetypeGreek, Attack: 2100, Defense: 1700, Cost: 4, Effect: "", EffectType: ""},
		{ID: "gr009", Name: "Artemis, the Hunter", Archetype: battle.ArchetypeGreek, Attack: 2100, Defense: 1700, Cost: 4, Effect: "", EffectType: ""},
		{ID: "gr009", Name: "Artemis, the Hunter", Archetype: battle.ArchetypeGreek, Attack: 2100, Defense: 1700, Cost: 4, Effect: "", EffectType: ""},
		{ID: "gr010", Name: "Hephaestus, the Forger", Archetype: battle.ArchetypeGreek, Attack: 1900, Defense: 2100, Cost: 4, Effect: "", EffectType: ""},
		{ID: "gr010", Name: "Hephaestus, the Forger", Archetype: battle.ArchetypeGreek, Attack: 1900, Defense: 2100, Cost: 4, Effect: "", EffectType: ""},
		{ID: "gr010", Name: "Hephaestus, the Forger", Archetype: battle.ArchetypeGreek, Attack: 1900, Defense: 2100, Cost: 4, Effect: "", EffectType: ""},
		{ID: "gr011", Name: "Greek Hoplite", Archetype: battle.ArchetypeGreek, Attack: 1300, Defense: 1700, Cost: 2, Effect: "", EffectType: ""},
		{ID: "gr011", Name: "Greek Hoplite", Archetype: battle.ArchetypeGreek, Attack: 1300, Defense: 1700, Cost: 2, Effect: "", EffectType: ""},
		{ID: "gr011", Name: "Greek Hoplite", Archetype: battle.ArchetypeGreek, Attack: 1300, Defense: 1700, Cost: 2, Effect: "", EffectType: ""},
		{ID: "gr012", Name: "Temple Guardian", Archetype: battle.ArchetypeGreek, Attack: 900, Defense: 2100, Cost: 2, Effect: "", EffectType: ""},
		{ID: "gr012", Name: "Temple Guardian", Archetype: battle.ArchetypeGreek, Attack: 900, Defense: 2100, Cost: 2, Effect: "", EffectType: ""},
		{ID: "gr012", Name: "Temple Guardian", Archetype: battle.ArchetypeGreek, Attack: 900, Defense: 2100, Cost: 2, Effect: "", EffectType: ""},
		{ID: "gr013", Name: "Oracle Priestess", Archetype: battle.ArchetypeGreek, Attack: 1000, Defense: 1500, Cost: 2, Effect: "Draw a card", EffectType: "draw"},
		{ID: "gr013", Name: "Oracle Priestess", Archetype: battle.ArchetypeGreek, Attack: 1000, Defense: 1500, Cost: 2, Effect: "Draw a card", EffectType: "draw"},
		{ID: "gr013", Name: "Oracle Priestess", Archetype: battle.ArchetypeGreek, Attack: 1000, Defense: 1500, Cost: 2, Effect: "Draw a card", EffectType: "draw"},
		
		// Spell/Effect cards (same as Egyptian deck for balance)
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