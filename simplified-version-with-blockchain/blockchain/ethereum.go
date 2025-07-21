package blockchain

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"
	"os"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"

	"cardgame/contracts" // This should contain your generated contract bindings
	"cardgame/models"
)

// EthereumClient handles all blockchain interactions
type EthereumClient struct {
	client          *ethclient.Client
	contract        *contracts.CardTrading
	privateKey      *ecdsa.PrivateKey
	contractAddress common.Address
}

// NewEthereumClient creates a new instance of EthereumClient
func NewEthereumClient(nodeURL, contractAddress string) (*EthereumClient, error) {
	// Connect to Ethereum client
	client, err := ethclient.Dial(nodeURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Ethereum client: %v", err)
	}

	// Get private key from environment variable
	privateKeyStr := os.Getenv("OWNER_PRIVATE_KEY")
	if privateKeyStr == "" {
		return nil, fmt.Errorf("OWNER_PRIVATE_KEY not set in environment")
	}

	// Remove "0x" prefix if present
	if len(privateKeyStr) > 2 && privateKeyStr[:2] == "0x" {
		privateKeyStr = privateKeyStr[2:]
	}

	// Parse private key
	privateKey, err := crypto.HexToECDSA(privateKeyStr)
	if err != nil {
		return nil, fmt.Errorf("invalid private key: %v", err)
	}

	// Create contract instance
	address := common.HexToAddress(contractAddress)
	contract, err := contracts.NewCardTrading(address, client)
	if err != nil {
		return nil, fmt.Errorf("failed to create contract instance: %v", err)
	}

	return &EthereumClient{
		client:          client,
		contract:        contract,
		privateKey:      privateKey,
		contractAddress: address,
	}, nil
}

// getTransactOpts prepares transaction options for blockchain interactions
func (ec *EthereumClient) getTransactOpts() (*bind.TransactOpts, error) {
	ctx := context.Background()

	// Get the network ID
	chainID, err := ec.client.ChainID(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get chain ID: %v", err)
	}

	// Get the account details
	publicKey := ec.privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("error casting public key to ECDSA")
	}

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)

	// Get the nonce
	nonce, err := ec.client.PendingNonceAt(ctx, fromAddress)
	if err != nil {
		return nil, fmt.Errorf("failed to get nonce: %v", err)
	}

	// Get gas price
	gasPrice, err := ec.client.SuggestGasPrice(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get gas price: %v", err)
	}

	// Create transaction opts
	auth, err := bind.NewKeyedTransactorWithChainID(ec.privateKey, chainID)
	if err != nil {
		return nil, fmt.Errorf("failed to create transactor: %v", err)
	}

	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0)     // in wei
	auth.GasLimit = uint64(300000) // in units
	auth.GasPrice = gasPrice

	return auth, nil
}

// MintCard mints a new card for a player
func (ec *EthereumClient) MintCard(playerAddress string, card *models.Card) error {
	auth, err := ec.getTransactOpts()
	if err != nil {
		return fmt.Errorf("failed to get transaction options: %v", err)
	}

	address := common.HexToAddress(playerAddress)
	tx, err := ec.contract.MintCard(
		auth,
		address,
		card.CardID,
		card.Element,
		big.NewInt(int64(card.Stats.Damage)),
		big.NewInt(int64(card.Stats.ManaCost)),
	)
	if err != nil {
		return fmt.Errorf("failed to mint card: %v", err)
	}

	return ec.waitForTransaction(tx)
}

// ListCard lists a card for sale
func (ec *EthereumClient) ListCard(tokenId *big.Int, price *big.Int) error {
	auth, err := ec.getTransactOpts()
	if err != nil {
		return fmt.Errorf("failed to get transaction options: %v", err)
	}

	tx, err := ec.contract.ListCard(auth, tokenId, price)
	if err != nil {
		return fmt.Errorf("failed to list card: %v", err)
	}

	return ec.waitForTransaction(tx)
}

// BuyCard purchases a card that is listed for sale
func (ec *EthereumClient) BuyCard(tokenId *big.Int, price *big.Int) error {
	auth, err := ec.getTransactOpts()
	if err != nil {
		return fmt.Errorf("failed to get transaction options: %v", err)
	}

	auth.Value = price // Set the ETH value to send with the transaction

	tx, err := ec.contract.BuyCard(auth, tokenId)
	if err != nil {
		return fmt.Errorf("failed to buy card: %v", err)
	}

	return ec.waitForTransaction(tx)
}

// GetCard retrieves card details by token ID
func (ec *EthereumClient) GetCard(tokenId *big.Int) (*contracts.CardTradingCard, error) {
	card, err := ec.contract.GetCard(nil, tokenId)
	if err != nil {
		return nil, fmt.Errorf("failed to get card: %v", err)
	}
	return &card, nil
}

// GetCardsOfOwner retrieves all cards owned by an address
func (ec *EthereumClient) GetCardsOfOwner(owner string) ([]*big.Int, error) {
	address := common.HexToAddress(owner)
	balance, err := ec.contract.BalanceOf(nil, address)
	if err != nil {
		return nil, fmt.Errorf("failed to get balance: %v", err)
	}

	tokens := make([]*big.Int, 0)
	for i := int64(0); i < balance.Int64(); i++ {
		tokenId, err := ec.contract.TokenOfOwnerByIndex(nil, address, big.NewInt(i))
		if err != nil {
			return nil, fmt.Errorf("failed to get token ID: %v", err)
		}
		tokens = append(tokens, tokenId)
	}

	return tokens, nil
}

// waitForTransaction waits for a transaction to be mined and checks its status
func (ec *EthereumClient) waitForTransaction(tx *types.Transaction) error {
	receipt, err := bind.WaitMined(context.Background(), ec.client, tx)
	if err != nil {
		return fmt.Errorf("failed to wait for transaction to be mined: %v", err)
	}

	if receipt.Status == 0 {
		return fmt.Errorf("transaction failed")
	}

	return nil
}

// GetBalance gets the ETH balance of an address
func (ec *EthereumClient) GetBalance(address string) (*big.Int, error) {
	account := common.HexToAddress(address)
	balance, err := ec.client.BalanceAt(context.Background(), account, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get balance: %v", err)
	}
	return balance, nil
}

// GetContractAddress returns the address of the deployed contract
func (ec *EthereumClient) GetContractAddress() common.Address {
	return ec.contractAddress
}

// Close closes the Ethereum client connection
func (ec *EthereumClient) Close() {
	if ec.client != nil {
		ec.client.Close()
	}
}
