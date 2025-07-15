# Card Battle Game

A terminal-based card battle game with Egyptian and Greek mythology themes.

## How to Run

1. Make sure you have Go installed (1.16 or later)
2. Copy the main.go content from the artifact to replace the placeholder
3. Run the game:
   ```bash
   go run main.go
   ```

## Game Rules

- Each player starts with 8000 HP
- Players gain 1 mana per turn (max 10)
- Egyptian cards get attack bonuses
- Greek cards get defense bonuses
- Reduce opponent's HP to 0 to win!

## Commands

- `draw` - Draw a card (draw phase only)
- `play [n]` - Play card number n from hand
- `attack [attacker] [target]` - Attack with a creature
- `battle` - Enter battle phase
- `end` - End your turn
- `help` - Show available commands
