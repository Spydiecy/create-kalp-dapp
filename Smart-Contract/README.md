# 🚀 Build Your Own Token Airdrop System on the Kalp Blockchain

### (Challenge-3): Junior

---

## What is this challenge about?

This challenge involves developing a token airdrop system using the Kalp blockchain. You'll create a smart contract in Go that manages the distribution of fungible tokens, allowing users to claim tokens through an airdrop mechanism. This simulates real-world scenarios where tokens are distributed to users for promotional purposes or community engagement.

## What will you learn?

**By participating in this challenge, you will:**

- Gain hands-on experience with the Go programming language.
- Understand the implementation of fungible tokens similar to ERC-20 standards.
- Learn how to develop and deploy smart contracts on the Kalp blockchain.
- Enhance your skills in blockchain development and decentralized applications (dApps).
- Explore concepts such as token minting, balance management, and token transfers.
- Improve your problem-solving and coding abilities in a competitive environment.

---

## Ready to Get Started?

### Let's understand Token Airdrops first!

Imagine a new café opens in town, and to attract customers, they hand out free coffee coupons to people passing by. These coupons can be redeemed for a free coffee, encouraging people to try their offerings. In the blockchain world, this concept translates to **token airdrops**.

A **token airdrop** is a method where blockchain projects distribute free tokens to users' wallets. It's a way to promote the project, incentivize early adopters, and build a community. Users can claim these tokens and later use them within the project's ecosystem or trade them on exchanges.

---

## Checkpoint 0: 📦 Installation

Before you begin, ensure you have the following:

### Step 1: 🖥 Download Go

- [Download Go](https://go.dev/doc/install) and install version `>=1.19` but `<1.20`.

---

**To start the project, follow these steps:**

1. **Clone the repository:**

   ```sh
   git clone https://github.com/Build-Hackathon/junior.git
   ```

2. **Navigate to the project directory:**

   ```sh
   cd airdrop-vending-machine
   ```

3. **Navigate to the smart contract folder:** 

   ```sh
   cd smart-contract
   ```  

3. **Install the dependencies:**

   ```sh
   go mod tidy
   ```

---

After executing the above commands, your folder structure should look like this:

```
airdrop-vending-machine
├── vendor
├── go.mod
├── go.sum
├── main.go
└── krc.go  (Your Airdrop Vending Machine contract file)
```

---

## Checkpoint 1: 🏗 Smart Contract Walkthrough

Let's **open the file `krc.go`** and dive deep into the Airdrop Vending Machine smart contract.

### 1. Initializing the Token Contract

The **Initialize** method sets up your token contract by assigning a name, symbol, and decimals to your fungible token. This method ensures that the token's metadata is properly configured before any transactions occur.

```go
func (s *SmartContract) Initialize(sdk kalpsdk.TransactionContextInterface, name string, symbol string, decimals string) (bool, error) {
    // Authorization check
    clientMSPID, err := sdk.GetClientIdentity().GetMSPID()
    if err != nil {
        return false, fmt.Errorf("failed to get MSPID: %v", err)
    }
    if clientMSPID != "mailabs" {
        return false, fmt.Errorf("client is not authorized to initialize contract")
    }

    // Ensure contract isn't already initialized
    bytes, err := sdk.GetState(nameKey)
    if err != nil {
        return false, fmt.Errorf("failed to get token name: %v", err)
    }
    if bytes != nil {
        return false, fmt.Errorf("contract options are already set")
    }

    // Set token details
    err = sdk.PutStateWithoutKYC(nameKey, []byte(name))
    if err != nil {
        return false, fmt.Errorf("failed to set token name: %v", err)
    }
    err = sdk.PutStateWithoutKYC(symbolKey, []byte(symbol))
    if err != nil {
        return false, fmt.Errorf("failed to set token symbol: %v", err)
    }
    err = sdk.PutStateWithoutKYC(decimalsKey, []byte(decimals))
    if err != nil {
        return false, fmt.Errorf("failed to set token decimals: %v", err)
    }

    return true, nil
}
```

- **Parameters:**

  - `name (string)`: The human-readable name of your token (e.g., "KalpToken").
  - `symbol (string)`: A short symbol representing your token (e.g., "KALP").
  - `decimals (string)`: The number of decimal places the token uses.

- **Return Values:**

  - `(bool, error)`: A boolean indicating success (`true`) or failure (`false`), and an error object if any issues arise.

### 2. 🖨 Claiming Tokens (Minting)

The **Claim** method allows authorized accounts to mint new tokens and assign them to a specific address. This is the core function of the airdrop vending machine, distributing tokens to users who claim them.

```go
func (s *SmartContract) Claim(sdk kalpsdk.TransactionContextInterface, amount int, address string) error {
    // Authorization check
    clientMSPID, err := sdk.GetClientIdentity().GetMSPID()
    if err != nil {
        return fmt.Errorf("failed to get MSPID: %v", err)
    }
    if clientMSPID != "mailabs" {
        return fmt.Errorf("client is not authorized to mint new tokens")
    }

    // Ensure amount is positive
    if amount <= 0 {
        return fmt.Errorf("mint amount must be a positive integer")
    }

    // Update balance of the address
    currentBalanceBytes, err := sdk.GetState(address)
    if err != nil {
        return fmt.Errorf("failed to read account %s: %v", address, err)
    }
    var currentBalance int
    if currentBalanceBytes == nil {
        currentBalance = 0
    } else {
        currentBalance, _ = strconv.Atoi(string(currentBalanceBytes))
    }

    updatedBalance, err := add(currentBalance, amount)
    if err != nil {
        return err
    }
    err = sdk.PutStateWithoutKYC(address, []byte(strconv.Itoa(updatedBalance)))
    if err != nil {
        return err
    }

    // Update total supply
    totalSupplyBytes, err := sdk.GetState(totalSupplyKey)
    if err != nil {
        return fmt.Errorf("failed to retrieve total token supply: %v", err)
    }
    var totalSupply int
    if totalSupplyBytes == nil {
        totalSupply = 0
    } else {
        totalSupply, _ = strconv.Atoi(string(totalSupplyBytes))
    }
    totalSupply, err = add(totalSupply, amount)
    if err != nil {
        return err
    }
    err = sdk.PutStateWithoutKYC(totalSupplyKey, []byte(strconv.Itoa(totalSupply)))
    if err != nil {
        return err
    }

    // Emit Transfer event
    transferEvent := event{"0x0", address, amount}
    transferEventJSON, err := json.Marshal(transferEvent)
    if err != nil {
        return fmt.Errorf("failed to encode event: %v", err)
    }
    err = sdk.SetEvent("Transfer", transferEventJSON)
    if err != nil {
        return fmt.Errorf("failed to set event: %v", err)
    }

    log.Printf("Account %s balance updated to %d", address, updatedBalance)
    return nil
}
```

- **Parameters:**

  - `amount (int)`: The number of tokens to mint.
  - `address (string)`: The recipient's address where tokens will be credited.

- **Return Value:**

  - `error`: An error object if any issues occur; otherwise, `nil` on success.

### 3. 📊 Checking Balances

#### a) BalanceOf

The **BalanceOf** method returns the token balance of a given account.

```go
func (s *SmartContract) BalanceOf(sdk kalpsdk.TransactionContextInterface, account string) (int, error) {
    balanceBytes, err := sdk.GetState(account)
    if err != nil {
        return 0, fmt.Errorf("failed to read account %s: %v", account, err)
    }
    if balanceBytes == nil {
        return 0, fmt.Errorf("account %s does not exist", account)
    }
    balance, _ := strconv.Atoi(string(balanceBytes))
    return balance, nil
}
```

- **Parameters:**

  - `account (string)`: The address of the account to query.

- **Return Values:**

  - `(int, error)`: The balance of the account, and an error object if any issues occur.

#### b) TotalSupply

The **TotalSupply** method returns the total number of tokens in circulation.

```go
func (s *SmartContract) TotalSupply(sdk kalpsdk.TransactionContextInterface) (int, error) {
    totalSupplyBytes, err := sdk.GetState(totalSupplyKey)
    if err != nil {
        return 0, fmt.Errorf("failed to retrieve total token supply: %v", err)
    }
    var totalSupply int
    if totalSupplyBytes == nil {
        totalSupply = 0
    } else {
        totalSupply, _ = strconv.Atoi(string(totalSupplyBytes))
    }
    log.Printf("TotalSupply: %d tokens", totalSupply)
    return totalSupply, nil
}
```

- **Return Values:**

  - `(int, error)`: The total supply of tokens, and an error object if any issues occur.

### 4. 🔄 Transferring Tokens

The **TransferFrom** method allows a spender to transfer tokens from one account to another, given that they have sufficient allowance.

```go
func (s *SmartContract) TransferFrom(sdk kalpsdk.TransactionContextInterface, from string, to string, value int) error {
    // Authorization and allowance checks omitted for brevity

    // Initiate the transfer
    err := transferHelper(sdk, from, to, value)
    if err != nil {
        return fmt.Errorf("failed to transfer: %v", err)
    }

    // Update allowance and emit Transfer event
    // Implementation details omitted
    return nil
}
```

- **Parameters:**

  - `from (string)`: The sender's address.
  - `to (string)`: The recipient's address.
  - `value (int)`: The amount of tokens to transfer.

- **Return Value:**

  - `error`: An error object if any issues occur; otherwise, `nil` on success.

### 5. 🏷 Reading Token Metadata

#### a) Name

Returns the name of the token.

```go
func (s *SmartContract) Name(sdk kalpsdk.TransactionContextInterface) (string, error) {
    bytes, err := sdk.GetState(nameKey)
    if err != nil {
        return "", fmt.Errorf("failed to get token name: %v", err)
    }
    return string(bytes), nil
}
```

- **Return Values:**

  - `(string, error)`: The token name, and an error object if any issues occur.

#### b) Symbol

Returns the symbol of the token.

```go
func (s *SmartContract) Symbol(sdk kalpsdk.TransactionContextInterface) (string, error) {
    bytes, err := sdk.GetState(symbolKey)
    if err != nil {
        return "", fmt.Errorf("failed to get token symbol: %v", err)
    }
    return string(bytes), nil
}
```

- **Return Values:**

  - `(string, error)`: The token symbol, and an error object if any issues occur.

---

## Checkpoint 2: 📀 Deploying the Smart Contract

Now it's time to deploy the contract.

Before you begin, ensure you have an account on the Kalp Studio Platform. You can create an account by following these steps:

1. **Sign Up and Log In to Kalp Studio Platform**

   - [Sign Up and Log In to Kalp Studio Platform](https://doc.kalp.studio/Getting-started/Onboarding/How-to-Sign-Up-and-Log-In-to-Kalp-Studio-Platform/)

After setting up your account, you can deploy the smart contract using Kalp Studio:

2. **Deploy a Smart Contract on Kalp Studio**

   - [Deploy a Smart Contract on Kalp Studio](https://doc.kalp.studio/Dev-documentation/Kalp-DLT/Smart-Contract-Write-Test-Deploy-Interact/Deploy-the-smart-contract/)

   - **Steps:**
     - **Access the Kalp Studio Dashboard.**
     - **Go to Kalp Instant Deployer.**
     - **Click on "Create New" Smart Contract.**
     - **Enter the details:**
       - **Name:** Enter a name for your smart contract.
       - **Category:** Choose a category for your smart contract.
       - **Description:** Optionally, add a description for your smart contract.
     - **Upload your `smart-contract.zip` file.**

> 💡 For a more descriptive deployment guide, please refer to the [detailed documentation](https://doc.kalp.studio/Dev-documentation/Kalp-DLT/Smart-Contract-Write-Test-Deploy-Interact/Deploy-the-smart-contract/).

---

## Checkpoint 3: 🕹️ Interacting with the Smart Contract

Now that we've deployed our token airdrop smart contract on the Kalp blockchain, it's time to interact with it. We'll perform the following actions:

1. **Initialize the Contract**
2. **Claim Tokens**
3. **Check Balance**
4. **Transfer Tokens**

### Prerequisites

Before interacting with the APIs, ensure you have the following:

- **[Download Postman](https://www.postman.com/downloads/)**
- **API Endpoint and API Key**: Obtained from Kalp Studio after deploying your smart contract.
- **Smart Contract Functions and Parameters**: Access to the routes and parameters for your smart contract functions, available under **Check Params** in Kalp Studio.

> 💡 For detailed instructions, refer to the [Interacting with Smart Contracts](https://doc.kalp.studio/Dev-documentation/Kalp-DLT/Smart-Contract-Write-Test-Deploy-Interact/Interacting-with-smart-contract/) guide.

---

### Example: Interacting with the Smart Contract Using Postman

#### Step-by-Step Guide

##### 1. **Initialize the Contract**

The `Initialize` function sets up your token contract by defining the name, symbol, and decimals.

**Instructions:**

1. **Get the API Endpoint for `Initialize`:**

   - In Kalp Studio, navigate to your deployed smart contract.
   - Click on **Check Params** next to the `Initialize` function to view the API route and parameters.

2. **Set Up Postman for a POST Request:**

   - Open Postman and create a new request.
   - Set the request method to **POST**.
   - Paste the API endpoint URL into the address bar.

3. **Add Headers for Authentication:**

   - In the **Headers** tab, add a new key-value pair:
     - **Key:** `x-api`
     - **Value:** *Your API Key* (obtained from Kalp Studio)

4. **Set Up the Request Body:**

   - Select the **Body** tab.
   - Choose **raw** and set the type to **JSON**.
   - Enter the following JSON, replacing the values with your desired token details:

     ```json
     {
       "name": "KalpToken",
       "symbol": "KALP",
       "decimals": "2"
     }
     ```

5. **Send the Request:**

   - Click **Send**.
   - If successful, you should receive a response indicating the contract has been initialized.

**Example Response:**

```json
{
  "success": true,
  "message": "Contract initialized successfully",
  "transactionid": "abcd1234..."
}
```

---

##### 2. **Claim Tokens**

The `Claim` function allows you to mint new tokens and assign them to a user's address.

**Instructions:**

1. **Get the API Endpoint for `Claim`:**

   - In Kalp Studio, under **Check Params**, find the route for the `Claim` function.

2. **Set Up Postman for a POST Request:**

   - Create a new request in Postman.
   - Set the request method to **POST**.
   - Paste the API endpoint URL.

3. **Add Headers for Authentication:**

   - Ensure the `x-api` header is set with your API Key.

4. **Set Up the Request Body:**

   - In the **Body** tab, select **raw** and **JSON**.
   - Enter the JSON data, replacing `"user-address"` with the recipient's address:

     ```json
     {
       "amount": 1000,
       "address": "user-address"
     }
     ```

5. **Send the Request:**

   - Click **Send**.
   - A successful response indicates the tokens have been claimed.

**Example Response:**

```json
{
  "success": true,
  "message": "Tokens claimed successfully",
  "transactionid": "efgh5678..."
}
```

---

##### 3. **Check Balance**

You can retrieve the balance of an account using the `BalanceOf` function.

**Instructions:**

1. **Get the API Endpoint for `BalanceOf`:**

   - In Kalp Studio, under **Check Params**, find the route for the `BalanceOf` function.

2. **Set Up Postman for a GET Request:**

   - Create a new request in Postman.
   - Set the request method to **GET**.
   - Paste the API endpoint URL.

3. **Add Query Parameters:**

   - In Postman, go to the **Params** tab.
   - Add a key `account` with the value of the address you want to check (e.g., `"user-address"`).

4. **Add Headers for Authentication:**

   - Ensure the `x-api` header is set with your API Key.

5. **Send the Request:**

   - Click **Send**.
   - The response will contain the account's balance.

**Example Response:**

```json
{
  "balance": 1000
}
```

---

##### 4. **Transfer Tokens**

The `TransferFrom` function allows you to transfer tokens from one account to another.

**Instructions:**

1. **Get the API Endpoint for `TransferFrom`:**

   - In Kalp Studio, under **Check Params**, find the route for the `TransferFrom` function.

2. **Set Up Postman for a POST Request:**

   - Create a new request in Postman.
   - Set the request method to **POST**.
   - Paste the API endpoint URL.

3. **Add Headers for Authentication:**

   - Ensure the `x-api` header is set with your API Key.

4. **Set Up the Request Body:**

   - In the **Body** tab, select **raw** and **JSON**.
   - Provide the transfer details:

     ```json
     {
       "from": "user-address",
       "to": "recipient-address",
       "value": 500
     }
     ```

5. **Send the Request:**

   - Click **Send**.
   - A successful response indicates the tokens have been transferred.

**Example Response:**

```json
{
  "success": true,
  "message": "Transfer successful",
  "transactionid": "ijkl9012..."
}
```

---

### Tips for Using Postman with Kalp Studio APIs

- **Ensure Correct HTTP Methods:** Use the appropriate HTTP method (GET or POST) as specified in the API documentation.
- **Authentication:** Always include the `x-api` header with your API Key.
- **Content-Type Header:** When sending JSON data in the body, ensure the `Content-Type` header is set to `application/json`. Postman usually sets this automatically when you select **raw** and **JSON** in the Body tab.
- **Error Handling:** If you receive an error, check the response message for details. Common issues include missing parameters or invalid authentication.



---

**Note:** Always keep your API Key secure. Do not share it publicly or commit it to version control systems.

For further exploration, consider implementing additional features like setting up allowances, burning tokens, or creating a user interface for easier interaction. Refer to the smart contract functions outlined in **Checkpoint 1** for guidance.

If you encounter any issues or need more detailed explanations, refer to the [Kalp Studio Documentation](https://doc.kalp.studio/) or seek support from the Kalp community.

---

**Feel free to explore further functionalities and modify the contract to suit your needs. Good luck, and happy coding!**
---

## Example: 🔒 Interacting with the Smart Contract Using Postman

### Prerequisites

1. **After setting up the environment**
2. **[Sign Up and Log In to Kalp Studio Platform](https://doc.kalp.studio/Getting-started/Onboarding/How-to-Sign-Up-and-Log-In-to-Kalp-Studio-Platform/)**
3. **[Deploy the Smart Contract on Kalp Studio](https://doc.kalp.studio/Dev-documentation/Kalp-DLT/Smart-Contract-Write-Test-Deploy-Interact/Deploy-the-smart-contract/)**
4. **[Download Postman](https://www.postman.com/downloads/)**

### Example

After deploying the smart contract in Kalp Studio, an API endpoint will be generated. This API endpoint can be used to interact with the deployed smart contract.

**Example of Generated API Endpoint in Kalp Studio:**

![API Endpoint Example](images/api-endpoint-example.png)

**API Route Details:**

![API Route Details](images/api-route-details.png)

**Generate API Auth Key:**

![Generate API Key](images/generate-api-key.png)

**Copy the API Key for Authorization:**

![Copy API Key](images/copy-api-key.png)

### Using Postman to Interact with the Contract

#### 1. Initialize the Contract

- **Method:** POST
- **URL:** `https://your-api-endpoint/initialize`
- **Headers:**
  - `x-api`: Your API key from Kalp Studio
- **Body (raw JSON):**

```json
{
  "name": "KalpToken",
  "symbol": "KALP",
  "decimals": "2"
}
```

**Response:**

```json
{
  "success": true,
  "message": "Contract initialized successfully",
  "transactionid": "abcd1234..."
}
```

#### 2. Claim Tokens

- **Method:** POST
- **URL:** `https://your-api-endpoint/claim`
- **Headers:**
  - `x-api`: Your API key from Kalp Studio
- **Body (raw JSON):**

```json
{
  "amount": 1000,
  "address": "user-address"
}
```

**Response:**

```json
{
  "success": true,
  "message": "Tokens claimed successfully",
  "transactionid": "efgh5678..."
}
```

#### 3. Check Balance

- **Method:** GET
- **URL:** `https://your-api-endpoint/balanceOf?account=user-address`
- **Headers:**
  - `x-api`: Your API key from Kalp Studio

**Response:**

```json
{
  "balance": 1000
}
```

#### 4. Transfer Tokens

- **Method:** POST
- **URL:** `https://your-api-endpoint/transferFrom`
- **Headers:**
  - `x-api`: Your API key from Kalp Studio
- **Body (raw JSON):**

```json
{
  "from": "user-address",
  "to": "recipient-address",
  "value": 500
}
```

**Response:**

```json
{
  "success": true,
  "message": "Transfer successful",
  "transactionid": "ijkl9012..."
}
```

---

## Conclusion

This README serves as an overview for developing and deploying a token airdrop system using Go on the Kalp blockchain. By completing this challenge, you've gained valuable experience in blockchain development, smart contract programming, and interacting with decentralized applications.

For additional details, refer to the [Kalp SDK documentation](https://doc.kalp.studio/).

---

**Feel free to explore further functionalities and modify the contract to suit your needs. Good luck, and happy coding!**
