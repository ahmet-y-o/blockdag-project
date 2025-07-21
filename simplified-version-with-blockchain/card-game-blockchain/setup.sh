#!/bin/bash

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Function to print status
print_status() {
    echo -e "${GREEN}[*] $1${NC}"
}

print_error() {
    echo -e "${RED}[!] $1${NC}"
}

print_warning() {
    echo -e "${YELLOW}[!] $1${NC}"
}

# Check if necessary tools are installed
check_dependencies() {
    print_status "Checking dependencies..."
    
    # Check Node.js
    if ! command -v node &> /dev/null; then
        print_error "Node.js is not installed. Please install it first."
        exit 1
    fi

    # Check npm
    if ! command -v npm &> /dev/null; then
        print_error "npm is not installed. Please install it first."
        exit 1
    }

    # Check Go
    if ! command -v go &> /dev/null; then
        print_error "Go is not installed. Please install it first."
        exit 1
    }

    # Check Truffle
    if ! command -v truffle &> /dev/null; then
        print_warning "Truffle not found. Installing..."
        npm install -g truffle
    fi

    # Check Ganache
    if ! command -v ganache &> /dev/null; then
        print_warning "Ganache not found. Installing..."
        npm install -g ganache
    fi
}

# Initialize project structure
init_project() {
    print_status "Initializing project structure..."
    
    # Create necessary directories
    mkdir -p contracts migrations scripts test

    # Create .env.example
    cat > .env.example << EOL
# Blockchain Configuration
BLOCKCHAIN_RPC_URL=http://localhost:8545
CONTRACT_ADDRESS=
OWNER_PRIVATE_KEY=

# Network Configuration
NETWORK_ID=1337
GAS_PRICE=20000000000
GAS_LIMIT=6721975

# API Configuration
API_PORT=8080
EOL

    # Create a sample .env file
    cp .env.example .env

    print_warning "Remember to update .env with your actual values!"
}

# Install dependencies
install_dependencies() {
    print_status "Installing dependencies..."
    
    # Install npm packages
    npm install @openzeppelin/contracts
    npm install @truffle/hdwallet-provider
    npm install dotenv
    
    # Install Go packages
    go get github.com/ethereum/go-ethereum
    go get github.com/joho/godotenv
    go mod tidy
}

# Setup smart contract
setup_contract() {
    print_status "Setting up smart contract..."
    
    # Compile contracts
    truffle compile
    
    # Generate Go bindings
    if [ -f build/contracts/CardTrading.json ]; then
        print_status "Generating Go bindings..."
        mkdir -p contracts/generated
        
        # Extract ABI
        cat build/contracts/CardTrading.json | jq -r '.abi' > build/CardTrading.abi
        
        # Generate Go bindings
        if command -v abigen &> /dev/null; then
            abigen --abi=build/CardTrading.abi --pkg=contracts --out=contracts/generated/cardtrading.go
        else
            print_warning "abigen not found. Please install go-ethereum tools to generate bindings."
        fi
    else
        print_error "Contract compilation failed or files not found."
    fi
}

# Create test scripts
create_test_scripts() {
    print_status "Creating test scripts..."
    
    # Create test script for smart contract
    cat > scripts/test_contract.js << EOL
const CardTrading = artifacts.require("CardTrading");

module.exports = async function(callback) {
    try {
        const cardTrading = await CardTrading.deployed();
        const accounts = await web3.eth.getAccounts();
        
        console.log("Contract deployed at:", cardTrading.address);
        
        // Add your test logic here
        
    } catch (error) {
        console.error(error);
    }
    callback();
};
EOL
}

# Main setup
main() {
    print_status "Starting setup..."
    
    check_dependencies
    init_project
    install_dependencies
    setup_contract
    create_test_scripts
    
    print_status "Setup completed!"
    print_warning "Remember to:"
    echo "1. Update .env with your configuration"
    echo "2. Never commit .env or private keys"
    echo "3. Run 'ganache' before testing"
    echo "4. Run 'truffle test' to verify setup"
}

# Run main setup
main#!/bin/bash

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Function to print status
print_status() {
    echo -e "${GREEN}[*] $1${NC}"
}

print_error() {
    echo -e "${RED}[!] $1${NC}"
}

print_warning() {
    echo -e "${YELLOW}[!] $1${NC}"
}

# Check if necessary tools are installed
check_dependencies() {
    print_status "Checking dependencies..."
    
    # Check Node.js
    if ! command -v node &> /dev/null; then
        print_error "Node.js is not installed. Please install it first."
        exit 1
    fi

    # Check npm
    if ! command -v npm &> /dev/null; then
        print_error "npm is not installed. Please install it first."
        exit 1
    }

    # Check Go
    if ! command -v go &> /dev/null; then
        print_error "Go is not installed. Please install it first."
        exit 1
    }

    # Check Truffle
    if ! command -v truffle &> /dev/null; then
        print_warning "Truffle not found. Installing..."
        npm install -g truffle
    fi

    # Check Ganache
    if ! command -v ganache &> /dev/null; then
        print_warning "Ganache not found. Installing..."
        npm install -g ganache
    fi
}

# Initialize project structure
init_project() {
    print_status "Initializing project structure..."
    
    # Create necessary directories
    mkdir -p contracts migrations scripts test

    # Create .env.example
    cat > .env.example << EOL
# Blockchain Configuration
BLOCKCHAIN_RPC_URL=http://localhost:8545
CONTRACT_ADDRESS=
OWNER_PRIVATE_KEY=

# Network Configuration
NETWORK_ID=1337
GAS_PRICE=20000000000
GAS_LIMIT=6721975

# API Configuration
API_PORT=8080
EOL

    # Create a sample .env file
    cp .env.example .env

    print_warning "Remember to update .env with your actual values!"
}

# Install dependencies
install_dependencies() {
    print_status "Installing dependencies..."
    
    # Install npm packages
    npm install @openzeppelin/contracts
    npm install @truffle/hdwallet-provider
    npm install dotenv
    
    # Install Go packages
    go get github.com/ethereum/go-ethereum
    go get github.com/joho/godotenv
    go mod tidy
}

# Setup smart contract
setup_contract() {
    print_status "Setting up smart contract..."
    
    # Compile contracts
    truffle compile
    
    # Generate Go bindings
    if [ -f build/contracts/CardTrading.json ]; then
        print_status "Generating Go bindings..."
        mkdir -p contracts/generated
        
        # Extract ABI
        cat build/contracts/CardTrading.json | jq -r '.abi' > build/CardTrading.abi
        
        # Generate Go bindings
        if command -v abigen &> /dev/null; then
            abigen --abi=build/CardTrading.abi --pkg=contracts --out=contracts/generated/cardtrading.go
        else
            print_warning "abigen not found. Please install go-ethereum tools to generate bindings."
        fi
    else
        print_error "Contract compilation failed or files not found."
    fi
}

# Create test scripts
create_test_scripts() {
    print_status "Creating test scripts..."
    
    # Create test script for smart contract
    cat > scripts/test_contract.js << EOL
const CardTrading = artifacts.require("CardTrading");

module.exports = async function(callback) {
    try {
        const cardTrading = await CardTrading.deployed();
        const accounts = await web3.eth.getAccounts();
        
        console.log("Contract deployed at:", cardTrading.address);
        
        // Add your test logic here
        
    } catch (error) {
        console.error(error);
    }
    callback();
};
EOL
}

# Main setup
main() {
    print_status "Starting setup..."
    
    check_dependencies
    init_project
    install_dependencies
    setup_contract
    create_test_scripts
    
    print_status "Setup completed!"
    print_warning "Remember to:"
    echo "1. Update .env with your configuration"
    echo "2. Never commit .env or private keys"
    echo "3. Run 'ganache' before testing"
    echo "4. Run 'truffle test' to verify setup"
}

# Run main setup
main