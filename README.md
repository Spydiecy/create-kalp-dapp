# Buildthon2
---

# 🚀 ERC20 Token Smart Contract on Kalp Chain

## Description
This ERC20 smart contract is written in Go and designed for the Kalp blockchain. It enables the creation, transfer, and management of fungible tokens that adhere to the ERC20 standard. Key functionalities include minting, transferring, burning tokens, checking balances, and managing allowances.

## Usage
After deploying the smart contract on the Kalp blockchain, various token operations such as minting, transferring, and burning can be performed using the API endpoint generated during deployment.

---

## Example
This smart contract allows for the creation of a fungible token on the Kalp blockchain, which can be transferred between users, burned, and used to track balances and allowances.

#### What is this contract about?
This contract is about managing fungible tokens on the Kalp blockchain, following the ERC20 standard. Participants will implement features like minting tokens, transferring tokens, burning tokens, and managing allowances. The contract can be applied in real-world use cases such as payments, digital assets, and decentralized finance (DeFi).

#### What will You learn?
**By working on this contract, you will:**

- Understand the ERC20 standard and how it works.
- Gain hands-on experience in smart contract development using the Go programming language.
- Learn how to mint, transfer, and burn tokens on the Kalp blockchain.
- Improve your skills in blockchain-based tokenization and digital asset management.
- Enhance your understanding of smart contract deployment and API integration.

### What is the goal of this contract?
The goal of this contract is to implement a fully functional ERC20 token on the Kalp blockchain, with support for minting, burning, transferring tokens, checking balances, and managing allowances.

#### Executive Summary
The ERC20 Token Smart Contract is designed to follow the ERC20 standard, allowing for the creation and management of fungible tokens on the Kalp blockchain. Users can mint tokens, transfer them between accounts, burn tokens, and manage allowances for other users to spend tokens on their behalf.

---

## Checkpoint 0: 📦 Installation

### Before you begin, ensure you have the following:

- **Step 1. 🖥 [Download-Go](https://go.dev/doc/install) Go version `>=1.19` but `<1.20`.**

---

**To start the project, follow these steps:**
1. Clone the repository:
   ```sh
   git clone https://github.com/yourusername/erc20-token-kalp.git
   ```

2. Navigate to the project directory:
   ```sh
   cd erc20-token-kalp
   ```

3. Install the dependencies:
   ```sh  
   go mod tidy
   ```

---

### Folder Structure:
```sh
Folder erc20-token-kalp
├── vendor
├── go.mod
├── go.sum
├── main.go
└── token.go  (Your ERC20 token contract file)
```

---

## Checkpoint 1: 🏗 ERC20 Contract Functions Walkthrough

### 1. Initialize Token Contract
The `Initialize` function is responsible for setting up the token name, symbol, and decimals for the ERC20 token. This ensures the contract is properly configured before any token operations can be executed.

```go
func (s *SmartContract) Initialize(sdk kalpsdk.TransactionContextInterface, name string, symbol string, decimals string) (bool, error)
```
##### Parameters:
- **name (string)**: The human-readable name of the token (e.g., "KalpToken").
- **symbol (string)**: A short symbol for the token (e.g., "KALP").
- **decimals (string)**: The number of decimal places for the token.

**Return Values:**
- **bool**: Returns true if the contract is successfully initialized.
- **error**: Returns an error if the contract cannot be initialized.

---

### 2. Minting Tokens
The `Mint` function allows authorized users to create new tokens and add them to the minter’s account balance.

```go
func (s *SmartContract) Mint(sdk kalpsdk.TransactionContextInterface, amount int) error
```
##### Parameters:
- **amount (int)**: The number of tokens to mint.

**Return Values:**
- **error**: Returns an error if minting fails.

---

### 3. Transferring Tokens
The `Transfer` function allows tokens to be transferred from one account to another.

```go
func (s *SmartContract) Transfer(sdk kalpsdk.TransactionContextInterface, recipient string, amount int) error
```
##### Parameters:
- **recipient (string)**: The address of the account to receive the tokens.
- **amount (int)**: The number of tokens to transfer.

**Return Values:**
- **error**: Returns an error if the transfer fails.

---

### 4. Burning Tokens
The `Burn` function allows the minter to destroy tokens, reducing the total supply.

```go
func (s *SmartContract) Burn(sdk kalpsdk.TransactionContextInterface, amount int) error
```
##### Parameters:
- **amount (int)**: The number of tokens to burn.

**Return Values:**
- **error**: Returns an error if burning fails.

---

### 5. Checking Balances
The `BalanceOf` function returns the balance of the given account.

```go
func (s *SmartContract) BalanceOf(sdk kalpsdk.TransactionContextInterface, account string) (int, error)
```
##### Parameters:
- **account (string)**: The address of the account to check the balance.

**Return Values:**
- **int**: The balance of the account.
- **error**: Returns an error if the balance cannot be retrieved.

---

### 6. Approving Allowances
The `Approve` function allows an account to grant another account permission to spend a specified number of tokens on its behalf.

```go
func (s *SmartContract) Approve(sdk kalpsdk.TransactionContextInterface, spender string, value int) error
```
##### Parameters:
- **spender (string)**: The address of the account being granted permission.
- **value (int)**: The number of tokens the spender is allowed to spend.

**Return Values:**
- **error**: Returns an error if the approval fails.

---

### 7. Transferring from an Allowed Account
The `TransferFrom` function allows the spender to transfer tokens from one account to another using the allowance granted by `Approve`.

```go
func (s *SmartContract) TransferFrom(sdk kalpsdk.TransactionContextInterface, from string, to string, value int) error
```
##### Parameters:
- **from (string)**: The account from which tokens will be transferred.
- **to (string)**: The account to receive the tokens.
- **value (int)**: The number of tokens to transfer.

**Return Values:**
- **error**: Returns an error if the transfer fails.

---

## Checkpoint 2: 🔒 Deploy Smart Contract

Before you begin, ensure you have the following:
1. **Sign Up and Log In to Kalp Studio Platform**
2. **Deploy the smart contract on the Kalp Studio Platform**

---

## Checkpoint 3: 📊 Total Supply and Token Information

### Total Supply
The `TotalSupply` function returns the total number of tokens that have been minted.

```go
func (s *SmartContract) TotalSupply(sdk kalpsdk.TransactionContextInterface) (int, error)
```

### Name and Symbol
The `Name` and `Symbol` functions return the descriptive name and short symbol of the token, respectively.

```go
func (s *SmartContract) Name(sdk kalpsdk.TransactionContextInterface) (string, error)
func (s *SmartContract) Symbol(sdk kalpsdk.TransactionContextInterface) (string, error)
```

---

## Checkpoint 4: 🔄 Interacting with the Smart Contract

After deploying the smart contract, use tools like Postman or Kalp Studio to interact with the deployed contract. You can make API calls for minting tokens, transferring tokens, and checking balances.

---

### Checkpoint 5: 🔄 API Interaction and Postman Setup

To interact with the deployed ERC20 contract, you can use API endpoints generated by Kalp Studio. Follow the steps below to set up and interact with the smart contract using Postman.

---

#### Step 1: Download and Set Up Postman

- [Download Postman](https://www.postman.com/downloads/) if you haven’t already.
- Install and launch Postman.

---

#### Step 2: Interact with the Smart Contract via API

Once the smart contract is deployed on the Kalp Studio platform, you can use the generated API endpoint to interact with the contract. Here is an example of how to interact with it:

1. **Mint Tokens**
    - Use the generated API route for minting tokens (e.g., `/mint`).
    - Set the request method to POST in Postman.
    - Add the required parameters like the number of tokens to mint in the request body.
    - Example:
      ```json
      {
        "amount": 1000
      }
      ```

2. **Transfer Tokens**
    - Use the generated API route for transferring tokens (e.g., `/transfer`).
    - Set the request method to POST.
    - Example:
      ```json
      {
        "recipient": "0x123456789",
        "amount": 500
      }
      ```

3. **Check Balance**
    - Use the generated API route for checking balances (e.g., `/balanceOf`).
    - Set the request method to GET.
    - Example response:
      ```json
      {
        "balance": 1000
      }
      ```

---

#### Step 3: Generate and Use API Keys

Before interacting with the contract via API, you may need to generate an API key for authentication. Follow these steps:
1. Go to the Kalp Studio platform.
2. Navigate to the API key generation section.
3. Generate a new API key and save it securely.
4. In Postman, add the API key as a header:
   - **Key**: `x-api-key`
   - **Value**: Your generated API key

---

### Example Interaction with Postman

- **Mint Tokens Request**:
  - **URL**: `https://kalpstudio.com/api/v1/mint`
  - **Method**: POST
  - **Headers**:
    - `x-api-key`: Your generated API key
  - **Body** (JSON):
    ```json
    {
      "amount": 1000
    }
    ```

- **Transfer Tokens Request**:
  - **URL**: `https://kalpstudio.com/api/v1/transfer`
  - **Method**: POST
  - **Headers**:
    - `x-api-key`: Your generated API key
  - **Body** (JSON):
    ```json
    {
      "recipient": "0x123456789",
      "amount": 500
    }
    ```

- **Check Balance Request**:
  - **URL**: `https://kalpstudio.com/api/v1/balanceOf`
  - **Method**: GET
  - **Headers**:
    - `x-api-key`: Your generated API key
  - **Params**: `account=0x123456789`

---

## Checkpoint 6: 🔒 Security Considerations

When interacting with the ERC20 contract:
- **API Keys**: Ensure your API keys are stored securely and not shared publicly.
- **Authorization**: Only authorized accounts (e.g., minters) should be allowed to mint and burn tokens.
- **Validations**: The smart contract already has in-built checks for negative values, duplicate accounts, and balance validations.

---

## Checkpoint 7: 📊 Understanding ERC20 Token Properties

- **Name**: The human-readable name of the token, returned by the `Name()` function.
- **Symbol**: The short symbol for the token, returned by the `Symbol()` function.
- **Decimals**: The number of decimals for token operations.
- **Total Supply**: The total number of tokens in circulation, returned by the `TotalSupply()` function.

---

## Example JSON Format for API Interactions

Here's an example of how the JSON structure looks when interacting with the smart contract via API:

```json
{
  "name": "KalpToken",
  "symbol": "KALP",
  "decimals": "18",
  "totalSupply": "1000000"
}
```

You can use this format for minting, transferring, and checking balances as shown in the previous examples.

---

## Conclusion

This README serves as a guide to deploying and interacting with the ERC20 token smart contract on the Kalp blockchain. By following the steps and examples provided, developers can easily mint tokens, manage token transfers, and work with allowances.

Feel free to modify and expand this contract according to your project’s requirements. For additional API integration and testing, refer to the Kalp Studio documentation or your smart contract’s generated API documentation.

--- 

Let me know if you'd like to add or change anything else in this README!

Let me know if you need further customization or help with any specific sections!
