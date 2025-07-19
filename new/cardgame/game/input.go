package game

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// InputHandler handles all user input
type InputHandler struct {
	reader *bufio.Reader
}

// NewInputHandler creates a new input handler
func NewInputHandler() *InputHandler {
	return &InputHandler{
		reader: bufio.NewReader(os.Stdin),
	}
}

// GetDeckChoice gets the user's deck selection
func (ih *InputHandler) GetDeckChoice() string {
	fmt.Print("\nChoose your deck (1 or 2): ")
	choice, _ := ih.reader.ReadString('\n')
	return strings.TrimSpace(choice)
}

// WaitForEnter waits for the user to press Enter
func (ih *InputHandler) WaitForEnter(message string) {
	if message != "" {
		fmt.Print(message)
	}
	ih.reader.ReadString('\n')
}

// GetCommand gets a command from the user
func (ih *InputHandler) GetCommand() (string, []string) {
	fmt.Printf("\nPlayer > ")
	input, _ := ih.reader.ReadString('\n')
	input = strings.TrimSpace(input)
	
	parts := strings.Split(input, " ")
	if len(parts) == 0 {
		return "", nil
	}
	
	command := strings.ToLower(parts[0])
	args := parts[1:]
	
	return command, args
}

// ParseCardIndex parses a card index from string
func (ih *InputHandler) ParseCardIndex(arg string) (int, error) {
	index, err := strconv.Atoi(arg)
	if err != nil {
		return -1, fmt.Errorf("invalid card number: %s", arg)
	}
	return index, nil
}

// ParseAttackTargets parses attacker and target indices
func (ih *InputHandler) ParseAttackTargets(args []string) (int, int, error) {
	if len(args) < 2 {
		return -1, -1, fmt.Errorf("usage: attack [attacker index] [target index]")
	}
	
	attackerIndex, err := strconv.Atoi(args[0])
	if err != nil {
		return -1, -1, fmt.Errorf("invalid attacker index: %s", args[0])
	}
	
	targetIndex, err := strconv.Atoi(args[1])
	if err != nil {
		return -1, -1, fmt.Errorf("invalid target index: %s", args[1])
	}
	
	return attackerIndex, targetIndex, nil
}