package models

import (
	"fmt"

	"cardgame/constants"

	"github.com/google/uuid"
)

type Card struct {
	ID      string
	CardID  string
	Stats   constants.CardDefinition
	Element string
}

func NewCard(cardID string) (*Card, error) {
	cardDef, exists := constants.AttackCards[cardID]
	if !exists {
		return nil, fmt.Errorf("card with ID %s does not exist", cardID)
	}

	return &Card{
		ID:      uuid.New().String(),
		CardID:  cardID,
		Stats:   cardDef,
		Element: cardDef.Element,
	}, nil
}

func (c *Card) GetDamage() int {
	return c.Stats.Damage
}

func (c *Card) GetManaCost() int {
	return c.Stats.ManaCost
}

func (c *Card) GetType() string {
	return c.Stats.Type
}

func (c *Card) GetName() string {
	return c.Stats.Name
}

func (c *Card) GetElement() string {
	return c.Element
}
