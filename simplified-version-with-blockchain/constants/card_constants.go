package constants

// Card Definitions
type CardDefinition struct {
	Name     string
	Type     string
	Element  string
	Damage   int
	ManaCost int
}

// Attack Cards Database
var AttackCards = map[string]CardDefinition{
	// Fire Element Cards
	"FIRE_SWORD": {
		Name:     "Fire Sword",
		Type:     TYPE_ATTACK,
		Element:  ELEMENT_FIRE,
		Damage:   25,
		ManaCost: 3,
	},
	"INFERNO_BLADE": {
		Name:     "Inferno Blade",
		Type:     TYPE_ATTACK,
		Element:  ELEMENT_FIRE,
		Damage:   35,
		ManaCost: 5,
	},

	// Ice Element Cards
	"ICE_DAGGER": {
		Name:     "Ice Dagger",
		Type:     TYPE_ATTACK,
		Element:  ELEMENT_ICE,
		Damage:   20,
		ManaCost: 2,
	},
	"FROST_SPEAR": {
		Name:     "Frost Spear",
		Type:     TYPE_ATTACK,
		Element:  ELEMENT_ICE,
		Damage:   30,
		ManaCost: 4,
	},

	// Earth Element Cards
	"STONE_FIST": {
		Name:     "Stone Fist",
		Type:     TYPE_ATTACK,
		Element:  ELEMENT_EARTH,
		Damage:   22,
		ManaCost: 2,
	},
	"EARTH_HAMMER": {
		Name:     "Earth Hammer",
		Type:     TYPE_ATTACK,
		Element:  ELEMENT_EARTH,
		Damage:   32,
		ManaCost: 4,
	},

	// Wind Element Cards
	"WIND_SLASH": {
		Name:     "Wind Slash",
		Type:     TYPE_ATTACK,
		Element:  ELEMENT_WIND,
		Damage:   18,
		ManaCost: 2,
	},
	"TORNADO_BLADE": {
		Name:     "Tornado Blade",
		Type:     TYPE_ATTACK,
		Element:  ELEMENT_WIND,
		Damage:   28,
		ManaCost: 4,
	},

	// Thunder Element Cards
	"THUNDER_STRIKE": {
		Name:     "Thunder Strike",
		Type:     TYPE_ATTACK,
		Element:  ELEMENT_THUNDER,
		Damage:   24,
		ManaCost: 3,
	},
	"LIGHTNING_AXE": {
		Name:     "Lightning Axe",
		Type:     TYPE_ATTACK,
		Element:  ELEMENT_THUNDER,
		Damage:   34,
		ManaCost: 5,
	},
}

// Helper function to get all cards
func GetAllCards() map[string]CardDefinition {
	return AttackCards
}

// Helper function to get cards by element
func GetCardsByElement(element string) map[string]CardDefinition {
	cards := make(map[string]CardDefinition)
	for id, card := range AttackCards {
		if card.Element == element {
			cards[id] = card
		}
	}
	return cards
}
