const CardTrading = artifacts.require("CardTrading");

module.exports = function(deployer) {
  deployer.deploy(CardTrading);
};