package handlers

import (
	"encoding/json"
	"math/big"
	"net/http"

	"cardgame/blockchain"
)

type TradingHandler struct {
	ethClient *blockchain.EthereumClient
}

type ListCardRequest struct {
	TokenId string `json:"token_id"`
	Price   string `json:"price"` // In wei
}

type BuyCardRequest struct {
	TokenId string `json:"token_id"`
}

func NewTradingHandler(ethClient *blockchain.EthereumClient) *TradingHandler {
	return &TradingHandler{
		ethClient: ethClient,
	}
}

func (th *TradingHandler) ListCard(w http.ResponseWriter, r *http.Request) {
	var req ListCardRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	tokenId := new(big.Int)
	tokenId.SetString(req.TokenId, 10)

	price := new(big.Int)
	price.SetString(req.Price, 10)

	if err := th.ethClient.ListCard(tokenId, price); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"status":  "success",
		"message": "Card listed successfully",
	})
}

func (th *TradingHandler) BuyCard(w http.ResponseWriter, r *http.Request) {
	var req BuyCardRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	tokenId := new(big.Int)
	tokenId.SetString(req.TokenId, 10)

	if err := th.ethClient.BuyCard(tokenId, nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"status":  "success",
		"message": "Card purchased successfully",
	})
}
