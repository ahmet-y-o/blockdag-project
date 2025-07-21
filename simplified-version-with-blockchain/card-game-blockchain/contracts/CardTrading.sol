// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import "@openzeppelin/contracts/token/ERC721/ERC721.sol";
import "@openzeppelin/contracts/security/ReentrancyGuard.sol";
import "@openzeppelin/contracts/access/Ownable.sol";

contract CardTrading is ERC721, ReentrancyGuard, Ownable {
    struct Card {
        string cardId;      // Original game card ID
        string element;
        uint256 damage;
        uint256 manaCost;
        uint256 price;      // Price in wei (0 if not for sale)
    }

    mapping(uint256 => Card) public cards;
    uint256 private _tokenIds;

    event CardListed(uint256 tokenId, uint256 price);
    event CardSold(uint256 tokenId, address from, address to, uint256 price);
    event CardMinted(uint256 tokenId, address to, string cardId);

    constructor() ERC721("GameCard", "CARD") {}

    function mintCard(
        address player,
        string memory cardId,
        string memory element,
        uint256 damage,
        uint256 manaCost
    ) public onlyOwner returns (uint256) {
        _tokenIds++;
        uint256 newTokenId = _tokenIds;
        
        _mint(player, newTokenId);
        cards[newTokenId] = Card(cardId, element, damage, manaCost, 0);
        
        emit CardMinted(newTokenId, player, cardId);
        return newTokenId;
    }

    function listCard(uint256 tokenId, uint256 price) public {
        require(ownerOf(tokenId) == msg.sender, "Not the card owner");
        cards[tokenId].price = price;
        emit CardListed(tokenId, price);
    }

    function buyCard(uint256 tokenId) public payable nonReentrant {
        Card memory card = cards[tokenId];
        require(card.price > 0, "Card not for sale");
        require(msg.value >= card.price, "Insufficient payment");
        
        address seller = ownerOf(tokenId);
        require(seller != msg.sender, "Cannot buy your own card");
        
        _transfer(seller, msg.sender, tokenId);
        
        payable(seller).transfer(msg.value);
        cards[tokenId].price = 0;
        
        emit CardSold(tokenId, seller, msg.sender, msg.value);
    }

    function getCard(uint256 tokenId) public view returns (Card memory) {
        require(_exists(tokenId), "Card does not exist");
        return cards[tokenId];
    }
}