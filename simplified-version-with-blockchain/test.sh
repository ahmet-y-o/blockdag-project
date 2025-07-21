#!/bin/bash

# Create battle
echo "Creating battle..."
battle_response=$(curl -s -X POST \
    -H "Content-Type: application/json" \
    -d '{"player1_id":1,"player2_id":2}' \
    http://localhost:8080/battle/create)

echo "Battle response: $battle_response"

# Extract battle ID
battle_id=$(echo $battle_response | grep -o '"battle_id":"[^"]*' | grep -o '[^"]*$')

if [ -z "$battle_id" ]; then
    echo "Failed to get battle ID from response"
    echo "Full response was: $battle_response"
    exit 1
fi

echo "Created battle with ID: $battle_id"

# Player 1 plays a card
echo -e "\nPlayer 1 playing card..."
play_response=$(curl -s -X POST \
    -H "Content-Type: application/json" \
    -d '{"player_id":1,"card_index":0}' \
    "http://localhost:8080/battle/$battle_id/play-card")

echo "Play response: $play_response"

# Get battle state
echo -e "\nGetting battle state..."
state_response=$(curl -s "http://localhost:8080/battle/$battle_id?player_id=1")

echo "Battle state: $state_response"

# List a card
echo "Listing card for sale..."
curl -X POST \
    -H "Content-Type: application/json" \
    -d '{"token_id":"1","price":"1000000000000000000"}' \
    http://localhost:8080/trading/list

# Buy a card
echo -e "\nBuying card..."
curl -X POST \
    -H "Content-Type: application/json" \
    -d '{"token_id":"1"}' \
    http://localhost:8080/trading/buy