#!/bin/bash

echo "Starting Card Game Online Test..."

# Start server in background
echo "Starting server..."
go run server/main.go &
SERVER_PID=$!

# Wait for server to start
sleep 2

# Start two clients
echo "Starting Player 1..."
gnome-terminal -- bash -c "go run client/main.go -name TestPlayer1; read"

echo "Starting Player 2..."
gnome-terminal -- bash -c "go run client/main.go -name TestPlayer2; read"

echo "Server PID: $SERVER_PID"
echo "Press Ctrl+C to stop the server"

# Wait for interrupt
trap "kill $SERVER_PID" INT
wait $SERVER_PID
