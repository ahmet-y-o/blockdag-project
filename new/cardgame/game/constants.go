package game

import "time"

// Color codes for terminal display
const (
	ColorReset  = "\033[0m"
	ColorRed    = "\033[31m"
	ColorGreen  = "\033[32m"
	ColorYellow = "\033[33m"
	ColorBlue   = "\033[34m"
	ColorPurple = "\033[35m"
	ColorCyan   = "\033[36m"
	ColorWhite  = "\033[37m"
	ColorGray   = "\033[90m"

	// Bold variants
	ColorBoldRed    = "\033[1;31m"
	ColorBoldGreen  = "\033[1;32m"
	ColorBoldYellow = "\033[1;33m"
	ColorBoldBlue   = "\033[1;34m"
	ColorBoldCyan   = "\033[1;36m"
)

// Game constants
const (
	StartingHP       = 8000
	StartingMana     = 1
	MaxMana          = 10
	StartingHandSize = 5
	DeckSize         = 40
	MaxFieldSize     = 5
	MaxHandSize      = 10

	// Damage/Heal amounts for effects
	LightningDamage = 1500
	HealingAmount   = 1000
	BurnDamage      = 500

	// Card draw amounts
	DrawCardEffect = 1
	DrawTwoEffect  = 2

	// Mana gain
	ManaGainEffect  = 1
	ManaBoostEffect = 2
)

// Timing constants
const (
	AIThinkDelay   = 1 * time.Second
	AIActionDelay  = 500 * time.Millisecond
	AnimationDelay = 300 * time.Millisecond
	MessageDelay   = 2 * time.Second
)

// Display constants
const (
	TerminalWidth = 80
	CardWidth     = 15
	Divider       = "=================================="
	ThinDivider   = "----------------------------------"
)
