const CardTrading = artifacts.require("CardTrading");

module.exports = async function(callback) {
  try {
    const cardTrading = await CardTrading.deployed();
    const accounts = await web3.eth.getAccounts();

    console.log("Contract deployed at:", cardTrading.address);

    // Mint a card
    const result = await cardTrading.mintCard(
      accounts[1],
      "FIRE_SWORD",
      "FIRE",
      25,
      3,
      { from: accounts[0] }
    );
    console.log("Card minted:", result.logs[0].args.tokenId.toString());

    // List the card for sale
    const tokenId = result.logs[0].args.tokenId;
    const price = web3.utils.toWei("0.1", "ether");
    await cardTrading.listCard(tokenId, price, { from: accounts[1] });
    console.log("Card listed for:", web3.utils.fromWei(price, "ether"), "ETH");

    // Buy the card
    await cardTrading.buyCard(tokenId, {
      from: accounts[2],
      value: price
    });
    console.log("Card bought by account:", accounts[2]);

    // Get card details
    const card = await cardTrading.getCard(tokenId);
    console.log("Card details:", card);

  } catch (error) {
    console.error(error);
  }
  callback();
};