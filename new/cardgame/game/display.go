package game

import (
	"fmt"
	"strings"
	"cardgame/battle"
)

// Display handles all game display functions
type Display struct{}

// NewDisplay creates a new display handler
func NewDisplay() *Display {
	return &Display{}
}

// ClearScreen clears the terminal screen
func (d *Display) ClearScreen() {
	fmt.Print("\033[H\033[2J")
}

// ShowWelcome displays the welcome screen
func (d *Display) ShowWelcome() {
	d.ClearScreen()
	d.ShowBanner()
	fmt.Println("\nGame Rules:")
	fmt.Println("• Each deck has 40 cards")
	fmt.Println("• You start with 5 cards in hand") 
	fmt.Println("• Draw 1 card each turn")
	fmt.Println("• Reduce opponent's HP to 0 to win!")
	fmt.Println("\nDeck Types:")
	fmt.Println("1. " + ColorYellow + "Egyptian Gods" + ColorReset + " (Attack focused - +10% ATK per Egyptian)")
	fmt.Println("2. " + ColorBlue + "Greek Gods" + ColorReset + " (Defense focused - +10% DEF per Greek)")
}

// ShowBanner displays the game banner
func (d *Display) ShowBanner() {
	fmt.Println(ColorCyan + "╔═══════════════════════════════════╗" + ColorReset)
	fmt.Println(ColorCyan + "║" + ColorYellow + "     ⚔️  CARD BATTLE GAME ⚔️      " + ColorCyan + "║" + ColorReset)
	fmt.Println(ColorCyan + "╚═══════════════════════════════════╝" + ColorReset)
}

// ShowGameState displays the current game state
func (d *Display) ShowGameState(game *battle.GameState) {
	d.ClearScreen()
	d.ShowBanner()
	d.ShowTurnInfo(game)
	
	// Show both players
	d.ShowOpponent(game.Player2)
	fmt.Println("\n" + ColorWhite + ThinDivider + ColorReset)
	d.ShowPlayer(game.Player1)
	
	// Show last action
	if game.LastAction != "" {
		fmt.Printf("\n%s► %s%s\n", ColorPurple, game.LastAction, ColorReset)
	}
	
	fmt.Println("\n" + ColorCyan + Divider + ColorReset)
}

// ShowTurnInfo displays turn and phase information
func (d *Display) ShowTurnInfo(game *battle.GameState) {
	phaseColor := d.getPhaseColor(game.Phase)
	fmt.Printf("\nTurn: %s%d%s | Phase: %s%s%s | Current: %s%s%s\n",
		ColorBoldCyan, game.TurnCount, ColorReset,
		phaseColor, game.Phase, ColorReset,
		ColorYellow, game.CurrentTurn, ColorReset)
}

// ShowOpponent displays opponent information
func (d *Display) ShowOpponent(opponent *battle.Player) {
	fmt.Printf("\n%s═══ %s (Opponent) ═══%s\n", ColorRed, opponent.ID, ColorReset)
	d.ShowPlayerStats(opponent, ColorRed)
	d.ShowField("Opponent's Field:", opponent.Field, false)
}

// ShowPlayer displays player information
func (d *Display) ShowPlayer(player *battle.Player) {
	fmt.Printf("\n%s═══ %s (You) ═══%s\n", ColorGreen, player.ID, ColorReset)
	d.ShowPlayerStats(player, ColorGreen)
	d.ShowField("Your Field:", player.Field, true)
	d.ShowHand(player)
}

// ShowPlayerStats displays HP and Mana
func (d *Display) ShowPlayerStats(player *battle.Player, color string) {
	hpBar := d.createHPBar(player.HP, StartingHP)
	manaBar := d.createManaBar(player.Mana, player.MaxMana)
	
	fmt.Printf("HP: %s%s%s %s\n", color, hpBar, ColorReset, d.formatHP(player.HP))
	fmt.Printf("Mana: %s%s%s (%d/%d) | Hand: %d | Deck: %d\n",
		ColorBlue, manaBar, ColorReset,
		player.Mana, player.MaxMana,
		len(player.Hand), len(player.Deck))
}

// ShowField displays cards on the field
func (d *Display) ShowField(title string, cards []battle.Card, showIndex bool) {
	fmt.Printf("\n%s\n", title)
	if len(cards) == 0 {
		fmt.Println(ColorGray + "  (empty field)" + ColorReset)
		return
	}
	
	for i, card := range cards {
		d.ShowFieldCard(card, i, showIndex)
	}
}

// ShowFieldCard displays a single card on the field
func (d *Display) ShowFieldCard(card battle.Card, index int, showIndex bool) {
	color := d.getCardColor(card.Archetype)
	indexStr := ""
	if showIndex {
		indexStr = fmt.Sprintf("[%d] ", index)
	}
	
	fmt.Printf("  %s%s%s%s%s (%s) - ATK: %s%d%s / DEF: %s%d%s",
		indexStr,
		color, card.Name, ColorReset,
		ColorGray, card.Archetype, ColorReset,
		ColorBoldRed, card.Attack, ColorReset,
		ColorBoldBlue, card.Defense, ColorReset)
	
	if card.Effect != "" {
		fmt.Printf(" %s[%s]%s", ColorPurple, card.Effect, ColorReset)
	}
	fmt.Println()
}

// ShowHand displays the player's hand
func (d *Display) ShowHand(player *battle.Player) {
	fmt.Println("\nYour Hand:")
	if len(player.Hand) == 0 {
		fmt.Println(ColorGray + "  (no cards in hand)" + ColorReset)
		return
	}
	
	for i, card := range player.Hand {
		d.ShowHandCard(card, i, player.Mana)
	}
}

// ShowHandCard displays a single card in hand
func (d *Display) ShowHandCard(card battle.Card, index int, playerMana int) {
	color := d.getCardColor(card.Archetype)
	costColor := ColorGreen
	playable := ""
	
	if card.Cost > playerMana {
		costColor = ColorRed
		playable = ColorRed + " (Not enough mana)" + ColorReset
	}
	
	fmt.Printf("  [%d] %s%s%s (Cost: %s%d%s) - ATK: %d / DEF: %d",
		index,
		color, card.Name, ColorReset,
		costColor, card.Cost, ColorReset,
		card.Attack, card.Defense)
	
	if card.Effect != "" {
		fmt.Printf(" - %s%s%s", ColorPurple, card.Effect, ColorReset)
	}
	
	fmt.Printf("%s\n", playable)
}

// ShowCommands displays available commands based on phase
func (d *Display) ShowCommands(phase battle.GamePhase, isPlayerTurn bool) {
	if !isPlayerTurn {
		fmt.Println("\n" + ColorYellow + "⏳ Waiting for opponent's turn..." + ColorReset)
		return
	}
	
	fmt.Println("\n" + ColorBoldCyan + "Available Commands:" + ColorReset)
	
	switch phase {
	case battle.PhaseDrawn:
		fmt.Println("  " + ColorGreen + "draw" + ColorReset + "     - Draw a card")
		
	case battle.PhaseMain:
		fmt.Println("  " + ColorGreen + "play [n]" + ColorReset + " - Play card number n from hand")
		fmt.Println("  " + ColorGreen + "battle" + ColorReset + "   - Enter battle phase")
		fmt.Println("  " + ColorGreen + "end" + ColorReset + "      - End your turn")
		
	case battle.PhaseBattle:
		fmt.Println("  " + ColorGreen + "attack [attacker] [target]" + ColorReset + " - Attack with your card")
		fmt.Println("  " + ColorGray + "                            (use -1 for direct attack)" + ColorReset)
		fmt.Println("  " + ColorGreen + "main" + ColorReset + "     - Return to main phase")
		fmt.Println("  " + ColorGreen + "end" + ColorReset + "      - End your turn")
		
	case battle.PhaseEnd:
		fmt.Println("  " + ColorGreen + "end" + ColorReset + "      - End your turn")
	}
	
	fmt.Println("  " + ColorGreen + "help" + ColorReset + "     - Show this help")
	fmt.Println("  " + ColorGreen + "quit" + ColorReset + "     - Exit the game")
}

// ShowGameOver displays the game over screen
func (d *Display) ShowGameOver(game *battle.GameState) {
	d.ClearScreen()
	fmt.Println(ColorCyan + "╔═══════════════════════════════════╗" + ColorReset)
	fmt.Println(ColorCyan + "║" + ColorBoldRed + "         GAME OVER!              " + ColorCyan + "║" + ColorReset)
	fmt.Println(ColorCyan + "╚═══════════════════════════════════╝" + ColorReset)
	
	winnerColor := ColorGreen
	if game.Winner == "AI" || game.Winner == game.Player2.ID {
		winnerColor = ColorRed
	}
	
	fmt.Printf("\n%s🏆 Winner: %s! 🏆%s\n", winnerColor, game.Winner, ColorReset)
	
	fmt.Println("\n" + ColorBoldCyan + "Final Statistics:" + ColorReset)
	fmt.Printf("├─ Your HP: %s%d%s\n", d.getHPColor(game.Player1.HP), game.Player1.HP, ColorReset)
	fmt.Printf("├─ Opponent HP: %s%d%s\n", d.getHPColor(game.Player2.HP), game.Player2.HP, ColorReset)
	fmt.Printf("├─ Total Turns: %d\n", game.TurnCount)
	fmt.Printf("└─ Cards Played: %d\n", len(game.Player1.Graveyard) + len(game.Player2.Graveyard))
	
	fmt.Println("\n" + ColorYellow + "Thanks for playing!" + ColorReset)
}

// ShowError displays error messages
func (d *Display) ShowError(err error) {
	fmt.Printf("%s❌ Error: %v%s\n", ColorRed, err, ColorReset)
	fmt.Print("Press Enter to continue...")
}

// ShowMessage displays a colored message
func (d *Display) ShowMessage(message string, color string) {
	fmt.Printf("%s%s%s\n", color, message, ColorReset)
}

// Helper functions

func (d *Display) getCardColor(archetype battle.Archetype) string {
	switch archetype {
	case battle.ArchetypeEgyptian:
		return ColorYellow
	case battle.ArchetypeGreek:
		return ColorBlue
	default:
		return ColorWhite
	}
}

func (d *Display) getPhaseColor(phase battle.GamePhase) string {
	switch phase {
	case battle.PhaseDrawn:
		return ColorCyan
	case battle.PhaseMain:
		return ColorGreen
	case battle.PhaseBattle:
		return ColorRed
	case battle.PhaseEnd:
		return ColorGray
	default:
		return ColorWhite
	}
}

func (d *Display) getHPColor(hp int) string {
	if hp > StartingHP*0.7 {
		return ColorGreen
	} else if hp > StartingHP*0.3 {
		return ColorYellow
	}
	return ColorRed
}

func (d *Display) formatHP(hp int) string {
	color := d.getHPColor(hp)
	return fmt.Sprintf("%s%d/%d%s", color, hp, StartingHP, ColorReset)
}

func (d *Display) createHPBar(current, max int) string {
	barLength := 20
	filled := int(float64(current) / float64(max) * float64(barLength))
	if filled < 0 {
		filled = 0
	}
	
	bar := strings.Repeat("█", filled) + strings.Repeat("░", barLength-filled)
	return bar
}

func (d *Display) createManaBar(current, max int) string {
	if max == 0 {
		return ""
	}
	gems := strings.Repeat("◆", current) + strings.Repeat("◇", max-current)
	return gems
}