// SPDX-License-Identifier: MIT
pragma solidity ^0.8.25;

import "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import "@openzeppelin/contracts/access/Ownable.sol";

contract Airdrop is Ownable {
    IERC20 public token;

    event TokensDistributed(address indexed user, uint256 amount);

    constructor(IERC20 _token) Ownable(msg.sender) {
        token = _token;
    }

    // Function to distribute tokens to multiple users
    function distributeTokens(address[] calldata recipients, uint256[] calldata amounts) external onlyOwner {
        require(recipients.length == amounts.length, "Airdrop: Recipients and amounts length mismatch");

        for (uint256 i = 0; i < recipients.length; i++) {
            require(token.transfer(recipients[i], amounts[i]), "Airdrop: Failed to transfer tokens");
            emit TokensDistributed(recipients[i], amounts[i]);
        }
    }

    // Function to withdraw tokens from the contract
    function withdrawTokens(uint256 amount) external onlyOwner {
        require(token.transfer(owner(), amount), "Airdrop: Failed to transfer tokens");
    }

    // Function to withdraw ETH from the contract
    function withdrawETH(uint256 amount) external onlyOwner {
        payable(owner()).transfer(amount);
    }

    // Function to get the token balance of the contract
    function getTokenBalance() external view returns (uint256) {
        return token.balanceOf(address(this));
    }

    // Update the token address
    function updateToken(IERC20 _token) external onlyOwner {
        token = _token;
    }    
}
