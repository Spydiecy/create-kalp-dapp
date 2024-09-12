package main
 
import (
    "encoding/json"
    "errors"
    "fmt"
    "log"
    "strconv"
    "github.com/p2eengineering/kalp-sdk-public/kalpsdk"
)

// Define key names for options
const nameKey = "name"
const symbolKey = "symbol"
const decimalsKey = "decimals"
const totalSupplyKey = "totalSupply"
 
// Define objectType names for prefix
const allowancePrefix = "allowance"
 
// Define key names for options
 
// SmartContract provides functions for transferring tokens between accounts
type SmartContract struct {
    kalpsdk.Contract
}
 
// event provides an organized struct for emitting events
type event struct {
    From  string `json:"from"`
    To    string `json:"to"`
    Value int    `json:"value"`
}

 
// Mint creates new tokens and adds them to minter's account balance
// This function triggers a Transfer event
func (s *SmartContract) Mint(sdk kalpsdk.TransactionContextInterface, amount int) error {
 
    // Check if contract has been intilized first
    initialized, err := checkInitialized(sdk)
    if err != nil {
        return fmt.Errorf("failed to check if contract is already initialized: %v", err)
    }
    if !initialized {
        return fmt.Errorf("contract options need to be set before calling any function, call Initialize() to initialize contract")
    }
 
    // Check minter authorization - this sample assumes mailabs is the central banker with privilege to mint new tokens
    clientMSPID, err := sdk.GetClientIdentity().GetMSPID()
    if err != nil {
        return fmt.Errorf("failed to get MSPID: %v", err)
    }
    if clientMSPID != "mailabs" {
        return fmt.Errorf("client is not authorized to mint new tokens")
    }
 
    // Get ID of submitting client identity
    minter, err := sdk.GetUserID()
    if err != nil {
        return fmt.Errorf("failed to get client id: %v", err)
    }
 
    if amount <= 0 {
        return fmt.Errorf("mint amount must be a positive integer")
    }
 
    currentBalanceBytes, err := sdk.GetState(minter)
    if err != nil {
        return fmt.Errorf("failed to read minter account %s from world state: %v", minter, err)
    }
 
    var currentBalance int
 
    // If minter current balance doesn't yet exist, we'll create it with a current balance of 0
    if currentBalanceBytes == nil {
        currentBalance = 0
    } else {
        currentBalance, _ = strconv.Atoi(string(currentBalanceBytes)) // Error handling not needed since Itoa() was used when setting the account balance, guaranteeing it was an integer.
    }
 
    updatedBalance, err := add(currentBalance, amount)
    if err != nil {
        return err
    }
 
    err = sdk.PutStateWithoutKYC(minter, []byte(strconv.Itoa(updatedBalance)))
    if err != nil {
        return err
    }
 
    // Update the totalSupply
    totalSupplyBytes, err := sdk.GetState(totalSupplyKey)
    if err != nil {
        return fmt.Errorf("failed to retrieve total token supply: %v", err)
    }
 
    var totalSupply int
 
    // If no tokens have been minted, initialize the totalSupply
    if totalSupplyBytes == nil {
        totalSupply = 0
    } else {
        totalSupply, _ = strconv.Atoi(string(totalSupplyBytes)) // Error handling not needed since Itoa() was used when setting the totalSupply, guaranteeing it was an integer.
    }
 
    // Add the mint amount to the total supply and update the state
    totalSupply, err = add(totalSupply, amount)
    if err != nil {
        return err
    }
 
    err = sdk.PutStateWithoutKYC(totalSupplyKey, []byte(strconv.Itoa(totalSupply)))
    if err != nil {
        return err
    }
 
    // Emit the Transfer event
    transferEvent := event{"0x0", minter, amount}
    transferEventJSON, err := json.Marshal(transferEvent)
    if err != nil {
        return fmt.Errorf("failed to obtain JSON encoding: %v", err)
    }
    err = sdk.SetEvent("Transfer", transferEventJSON)
    if err != nil {
        return fmt.Errorf("failed to set event: %v", err)
    }
 
    log.Printf("minter account %s balance updated from %d to %d", minter, currentBalance, updatedBalance)
 
    return nil
}
 
// Burn redeems tokens the minter's account balance
// This function triggers a Transfer event
func (s *SmartContract) Burn(sdk kalpsdk.TransactionContextInterface, amount int) error {
 
    // Check if contract has been intilized first
    initialized, err := checkInitialized(sdk)
    if err != nil {
        return fmt.Errorf("failed to check if contract is already initialized: %v", err)
    }
    if !initialized {
        return fmt.Errorf("Contract options need to be set before calling any function, call Initialize() to initialize contract")
    }
    // Check minter authorization - this sample assumes Org1 is the central banker with privilege to burn new tokens
    clientMSPID, err := sdk.GetClientIdentity().GetMSPID()
    if err != nil {
        return fmt.Errorf("failed to get MSPID: %v", err)
    }
    if clientMSPID != "mailabs" {
        return fmt.Errorf("client is not authorized to mint new tokens")
    }
 
    // Get ID of submitting client identity
    minter, err := sdk.GetUserID()
    if err != nil {
        return fmt.Errorf("failed to get client id: %v", err)
    }
 
    if amount <= 0 {
        return errors.New("burn amount must be a positive integer")
    }
 
    currentBalanceBytes, err := sdk.GetState(minter)
    if err != nil {
        return fmt.Errorf("failed to read minter account %s from world state: %v", minter, err)
    }
 
    var currentBalance int
 
    // Check if minter current balance exists
    if currentBalanceBytes == nil {
        return errors.New("The balance does not exist")
    }
 
    currentBalance, _ = strconv.Atoi(string(currentBalanceBytes)) // Error handling not needed since Itoa() was used when setting the account balance, guaranteeing it was an integer.
 
    updatedBalance, err := sub(currentBalance, amount)
    if err != nil {
        return err
    }
 
    err = sdk.PutStateWithoutKYC(minter, []byte(strconv.Itoa(updatedBalance)))
    if err != nil {
        return err
    }
 
    // Update the totalSupply
    totalSupplyBytes, err := sdk.GetState(totalSupplyKey)
    if err != nil {
        return fmt.Errorf("failed to retrieve total token supply: %v", err)
    }
 
    // If no tokens have been minted, throw error
    if totalSupplyBytes == nil {
        return errors.New("totalSupply does not exist")
    }
 
    totalSupply, _ := strconv.Atoi(string(totalSupplyBytes)) // Error handling not needed since Itoa() was used when setting the totalSupply, guaranteeing it was an integer.
 
    // Subtract the burn amount to the total supply and update the state
    totalSupply, err = sub(totalSupply, amount)
    if err != nil {
        return err
    }
 
    err = sdk.PutStateWithoutKYC(totalSupplyKey, []byte(strconv.Itoa(totalSupply)))
    if err != nil {
        return err
    }
 
    // Emit the Transfer event
    transferEvent := event{minter, "0x0", amount}
    transferEventJSON, err := json.Marshal(transferEvent)
    if err != nil {
        return fmt.Errorf("failed to obtain JSON encoding: %v", err)
    }
    err = sdk.SetEvent("Transfer", transferEventJSON)
    if err != nil {
        return fmt.Errorf("failed to set event: %v", err)
    }
 
    log.Printf("minter account %s balance updated from %d to %d", minter, currentBalance, updatedBalance)
 
    return nil
}
 
// Transfer transfers tokens from client account to recipient account
// recipient account must be a valid clientID as returned by the ClientID() function
// This function triggers a Transfer event
func (s *SmartContract) Transfer(sdk kalpsdk.TransactionContextInterface, recipient string, amount int) error {
 
    // Check if contract has been intilized first
    initialized, err := checkInitialized(sdk)
    if err != nil {
        return fmt.Errorf("failed to check if contract is already initialized: %v", err)
    }
    if !initialized {
        return fmt.Errorf("Contract options need to be set before calling any function, call Initialize() to initialize contract")
    }
 
    // Get ID of submitting client identity
    clientID, err := sdk.GetUserID()
    if err != nil {
        return fmt.Errorf("failed to get client id: %v", err)
    }
 
    err = transferHelper(sdk, clientID, recipient, amount)
    if err != nil {
        return fmt.Errorf("failed to transfer: %v", err)
    }
 
    // Emit the Transfer event
    transferEvent := event{clientID, recipient, amount}
    transferEventJSON, err := json.Marshal(transferEvent)
    if err != nil {
        return fmt.Errorf("failed to obtain JSON encoding: %v", err)
    }
    err = sdk.SetEvent("Transfer", transferEventJSON)
    if err != nil {
        return fmt.Errorf("failed to set event: %v", err)
    }
 
    return nil
}
 
// BalanceOf returns the balance of the given account
func (s *SmartContract) BalanceOf(sdk kalpsdk.TransactionContextInterface, account string) (int, error) {
 
    // Check if contract has been intilized first
    initialized, err := checkInitialized(sdk)
    if err != nil {
        return 0, fmt.Errorf("failed to check if contract is already initialized: %v", err)
    }
    if !initialized {
        return 0, fmt.Errorf("Contract options need to be set before calling any function, call Initialize() to initialize contract")
    }
 
    balanceBytes, err := sdk.GetState(account)
    if err != nil {
        return 0, fmt.Errorf("failed to read from world state: %v", err)
    }
    if balanceBytes == nil {
        return 0, fmt.Errorf("the account %s does not exist", account)
    }
 
    balance, _ := strconv.Atoi(string(balanceBytes)) // Error handling not needed since Itoa() was used when setting the account balance, guaranteeing it was an integer.
 
    return balance, nil
}
 
// ClientAccountBalance returns the balance of the requesting client's account
func (s *SmartContract) ClientAccountBalance(sdk kalpsdk.TransactionContextInterface) (int, error) {
 
    // Check if contract has been intilized first
    initialized, err := checkInitialized(sdk)
    if err != nil {
        return 0, fmt.Errorf("failed to check if contract is already initialized: %v", err)
    }
    if !initialized {
        return 0, fmt.Errorf("Contract options need to be set before calling any function, call Initialize() to initialize contract")
    }
 
    // Get ID of submitting client identity
    clientID, err := sdk.GetUserID()
    if err != nil {
        return 0, fmt.Errorf("failed to get client id: %v", err)
    }
 
    balanceBytes, err := sdk.GetState(clientID)
    if err != nil {
        return 0, fmt.Errorf("failed to read from world state: %v", err)
    }
    if balanceBytes == nil {
        return 0, fmt.Errorf("the account %s does not exist", clientID)
    }
 
    balance, _ := strconv.Atoi(string(balanceBytes)) // Error handling not needed since Itoa() was used when setting the account balance, guaranteeing it was an integer.
 
    return balance, nil
}
 
// ClientAccountID returns the id of the requesting client's account
// In this implementation, the client account ID is the clientId itself
// Users can use this function to get their own account id, which they can then give to others as the payment address
func (s *SmartContract) ClientAccountID(sdk kalpsdk.TransactionContextInterface) (string, error) {
 
    // Check if contract has been intilized first
    initialized, err := checkInitialized(sdk)
    if err != nil {
        return "", fmt.Errorf("failed to check if contract is already initialized: %v", err)
    }
    if !initialized {
        return "", fmt.Errorf("Contract options need to be set before calling any function, call Initialize() to initialize contract")
    }
 
    // Get ID of submitting client identity
    clientAccountID, err := sdk.GetUserID()
    
    if err != nil {
        return "", fmt.Errorf("failed to get client id: %v", err)
    }
 
    return clientAccountID, nil
}
 
// TotalSupply returns the total token supply
func (s *SmartContract) TotalSupply(sdk kalpsdk.TransactionContextInterface) (int, error) {
 
    // Check if contract has been intilized first
    initialized, err := checkInitialized(sdk)
    if err != nil {
        return 0, fmt.Errorf("failed to check if contract is already initialized: %v", err)
    }
    if !initialized {
        return 0, fmt.Errorf("Contract options need to be set before calling any function, call Initialize() to initialize contract")
    }
 
    // Retrieve total supply of tokens from state of smart contract
    totalSupplyBytes, err := sdk.GetState(totalSupplyKey)
    if err != nil {
        return 0, fmt.Errorf("failed to retrieve total token supply: %v", err)
    }
 
    var totalSupply int
 
    // If no tokens have been minted, return 0
    if totalSupplyBytes == nil {
        totalSupply = 0
    } else {
        totalSupply, _ = strconv.Atoi(string(totalSupplyBytes)) // Error handling not needed since Itoa() was used when setting the totalSupply, guaranteeing it was an integer.
    }
 
    log.Printf("TotalSupply: %d tokens", totalSupply)
 
    return totalSupply, nil
}
 
// Approve allows the spender to withdraw from the calling client's token account
// The spender can withdraw multiple times if necessary, up to the value amount
// This function triggers an Approval event
func (s *SmartContract) Approve(sdk kalpsdk.TransactionContextInterface, spender string, value int) error {
 
    // Check if contract has been intilized first
    initialized, err := checkInitialized(sdk)
    if err != nil {
        return fmt.Errorf("failed to check if contract is already initialized: %v", err)
    }
    if !initialized {
        return fmt.Errorf("Contract options need to be set before calling any function, call Initialize() to initialize contract")
    }
 
    // Get ID of submitting client identity
    owner, err := sdk.GetUserID()
    if err != nil {
        return fmt.Errorf("failed to get client id: %v", err)
    }
 
    // Create allowanceKey
    allowanceKey, err := sdk.CreateCompositeKey(allowancePrefix, []string{owner, spender})
    if err != nil {
        return fmt.Errorf("failed to create the composite key for prefix %s: %v", allowancePrefix, err)
    }
 
    // Update the state of the smart contract by adding the allowanceKey and value
    err = sdk.PutStateWithoutKYC(allowanceKey, []byte(strconv.Itoa(value)))
    if err != nil {
        return fmt.Errorf("failed to update state of smart contract for key %s: %v", allowanceKey, err)
    }
 
    // Emit the Approval event
    approvalEvent := event{owner, spender, value}
    approvalEventJSON, err := json.Marshal(approvalEvent)
    if err != nil {
        return fmt.Errorf("failed to obtain JSON encoding: %v", err)
    }
    err = sdk.SetEvent("Approval", approvalEventJSON)
    if err != nil {
        return fmt.Errorf("failed to set event: %v", err)
    }
 
    log.Printf("client %s approved a withdrawal allowance of %d for spender %s", owner, value, spender)
 
    return nil
}
 
// Allowance returns the amount still available for the spender to withdraw from the owner
func (s *SmartContract) Allowance(sdk kalpsdk.TransactionContextInterface, owner string, spender string) (int, error) {
 
    // Check if contract has been intilized first
    initialized, err := checkInitialized(sdk)
    if err != nil {
        return 0, fmt.Errorf("failed to check if contract is already initialized: %v", err)
    }
    if !initialized {
        return 0, fmt.Errorf("Contract options need to be set before calling any function, call Initialize() to initialize contract")
    }
 
    // Create allowanceKey
    allowanceKey, err := sdk.CreateCompositeKey(allowancePrefix, []string{owner, spender})
    if err != nil {
        return 0, fmt.Errorf("failed to create the composite key for prefix %s: %v", allowancePrefix, err)
    }
 
    // Read the allowance amount from the world state
    allowanceBytes, err := sdk.GetState(allowanceKey)
    if err != nil {
        return 0, fmt.Errorf("failed to read allowance for %s from world state: %v", allowanceKey, err)
    }
 
    var allowance int
 
    // If no current allowance, set allowance to 0
    if allowanceBytes == nil {
        allowance = 0
    } else {
        allowance, err = strconv.Atoi(string(allowanceBytes)) // Error handling not needed since Itoa() was used when setting the totalSupply, guaranteeing it was an integer.
    }
 
    log.Printf("The allowance left for spender %s to withdraw from owner %s: %d", spender, owner, allowance)
 
    return allowance, nil
}
 
// TransferFrom transfers the value amount from the "from" address to the "to" address
// This function triggers a Transfer event
func (s *SmartContract) TransferFrom(sdk kalpsdk.TransactionContextInterface, from string, to string, value int) error {
 
    // Check if contract has been intilized first
    initialized, err := checkInitialized(sdk)
    if err != nil {
        return fmt.Errorf("failed to check if contract is already initialized: %v", err)
    }
    if !initialized {
        return fmt.Errorf("Contract options need to be set before calling any function, call Initialize() to initialize contract")
    }
 
    // Get ID of submitting client identity
    spender, err := sdk.GetUserID()
    if err != nil {
        return fmt.Errorf("failed to get client id: %v", err)
    }
 
    // Create allowanceKey
    allowanceKey, err := sdk.CreateCompositeKey(allowancePrefix, []string{from, spender})
    if err != nil {
        return fmt.Errorf("failed to create the composite key for prefix %s: %v", allowancePrefix, err)
    }
 
    // Retrieve the allowance of the spender
    currentAllowanceBytes, err := sdk.GetState(allowanceKey)
    if err != nil {
        return fmt.Errorf("failed to retrieve the allowance for %s from world state: %v", allowanceKey, err)
    }
 
    var currentAllowance int
    currentAllowance, _ = strconv.Atoi(string(currentAllowanceBytes)) // Error handling not needed since Itoa() was used when setting the totalSupply, guaranteeing it was an integer.
 
    // Check if transferred value is less than allowance
    if currentAllowance < value {
        return fmt.Errorf("spender does not have enough allowance for transfer")
    }
 
    // Initiate the transfer
    err = transferHelper(sdk, from, to, value)
    if err != nil {
        return fmt.Errorf("failed to transfer: %v", err)
    }
 
    // Decrease the allowance
    updatedAllowance, err := sub(currentAllowance, value)
    if err != nil {
        return err
    }
 
    err = sdk.PutStateWithoutKYC(allowanceKey, []byte(strconv.Itoa(updatedAllowance)))
    if err != nil {
        return err
    }
 
    // Emit the Transfer event
    transferEvent := event{from, to, value}
    transferEventJSON, err := json.Marshal(transferEvent)
    if err != nil {
        return fmt.Errorf("failed to obtain JSON encoding: %v", err)
    }
    err = sdk.SetEvent("Transfer", transferEventJSON)
    if err != nil {
        return fmt.Errorf("failed to set event: %v", err)
    }
 
    log.Printf("spender %s allowance updated from %d to %d", spender, currentAllowance, updatedAllowance)
 
    return nil
}
 
// Name returns a descriptive name for fungible tokens in this contract
// returns {String} Returns the name of the token
 
func (s *SmartContract) Name(sdk kalpsdk.TransactionContextInterface) (string, error) {
 
    // Check if contract has been intilized first
    initialized, err := checkInitialized(sdk)
    if err != nil {
        return "", fmt.Errorf("failed to check if contract is already initialized: %v", err)
    }
    if !initialized {
        return "", fmt.Errorf("Contract options need to be set before calling any function, call Initialize() to initialize contract")
    }
 
    bytes, err := sdk.GetState(nameKey)
    if err != nil {
        return "", fmt.Errorf("failed to get Name bytes: %s", err)
    }
 
    return string(bytes), nil
}
 
// Symbol returns an abbreviated name for fungible tokens in this contract.
// returns {String} Returns the symbol of the token
 
func (s *SmartContract) Symbol(sdk kalpsdk.TransactionContextInterface) (string, error) {
 
    // Check if contract has been intilized first
    initialized, err := checkInitialized(sdk)
    if err != nil {
        return "", fmt.Errorf("failed to check if contract is already initialized: %v", err)
    }
    if !initialized {
        return "", fmt.Errorf("Contract options need to be set before calling any function, call Initialize() to initialize contract")
    }
 
    bytes, err := sdk.GetState(symbolKey)
    if err != nil {
        return "", fmt.Errorf("failed to get Symbol: %v", err)
    }
 
    return string(bytes), nil
}
 
// Set information for a token and intialize contract.
// param {String} name The name of the token
// param {String} symbol The symbol of the token
// param {String} decimals The decimals used for the token operations
func (s *SmartContract) Initialize(sdk kalpsdk.TransactionContextInterface, name string, symbol string, decimals string) (bool, error) {
 
    // Check minter authorization - this sample assumes Org1 is the central banker with privilege to intitialize contract
    clientMSPID, err := sdk.GetClientIdentity().GetMSPID()
    if err != nil {
        return false, fmt.Errorf("failed to get MSPID: %v", err)
    }
    if clientMSPID != "mailabs" {
        return false, fmt.Errorf("client is not authorized to initialize contract")
    }
 
    // Check contract options are not already set, client is not authorized to change them once intitialized
    bytes, err := sdk.GetState(nameKey)
    if err != nil {
        return false, fmt.Errorf("failed to get Name: %v", err)
    }
    if bytes != nil {
        return false, fmt.Errorf("contract options are already set, client is not authorized to change them")
    }
 
    err = sdk.PutStateWithoutKYC(nameKey, []byte(name))
    if err != nil {
        return false, fmt.Errorf("failed to set token name: %v", err)
    }
 
    err = sdk.PutStateWithoutKYC(symbolKey, []byte(symbol))
    if err != nil {
        return false, fmt.Errorf("failed to set symbol: %v", err)
    }
 
    err = sdk.PutStateWithoutKYC(decimalsKey, []byte(decimals))
    if err != nil {
        return false, fmt.Errorf("failed to set token name: %v", err)
    }
 
    return true, nil
}
 
// Helper Functions
 
// transferHelper is a helper function that transfers tokens from the "from" address to the "to" address
// Dependant functions include Transfer and TransferFrom
func transferHelper(sdk kalpsdk.TransactionContextInterface, from string, to string, value int) error {
 
    if from == to {
        return fmt.Errorf("cannot transfer to and from same client account")
    }
 
    if value < 0 { // transfer of 0 is allowed in ERC-20, so just validate against negative amounts
        return fmt.Errorf("transfer amount cannot be negative")
    }
 
    fromCurrentBalanceBytes, err := sdk.GetState(from)
    if err != nil {
        return fmt.Errorf("failed to read client account %s from world state: %v", from, err)
    }
 
    if fromCurrentBalanceBytes == nil {
        return fmt.Errorf("client account %s has no balance", from)
    }
 
    fromCurrentBalance, _ := strconv.Atoi(string(fromCurrentBalanceBytes)) // Error handling not needed since Itoa() was used when setting the account balance, guaranteeing it was an integer.
 
    if fromCurrentBalance < value {
        return fmt.Errorf("client account %s has insufficient funds", from)
    }
 
    toCurrentBalanceBytes, err := sdk.GetState(to)
    if err != nil {
        return fmt.Errorf("failed to read recipient account %s from world state: %v", to, err)
    }
 
    var toCurrentBalance int
    // If recipient current balance doesn't yet exist, we'll create it with a current balance of 0
    if toCurrentBalanceBytes == nil {
        toCurrentBalance = 0
    } else {
        toCurrentBalance, _ = strconv.Atoi(string(toCurrentBalanceBytes)) // Error handling not needed since Itoa() was used when setting the account balance, guaranteeing it was an integer.
    }
 
    fromUpdatedBalance, err := sub(fromCurrentBalance, value)
    if err != nil {
        return err
    }
 
    toUpdatedBalance, err := add(toCurrentBalance, value)
    if err != nil {
        return err
    }
 
    err = sdk.PutStateWithoutKYC(from, []byte(strconv.Itoa(fromUpdatedBalance)))
    if err != nil {
        return err
    }
 
    err = sdk.PutStateWithoutKYC(to, []byte(strconv.Itoa(toUpdatedBalance)))
    if err != nil {
        return err
    }
 
    log.Printf("client %s balance updated from %d to %d", from, fromCurrentBalance, fromUpdatedBalance)
    log.Printf("recipient %s balance updated from %d to %d", to, toCurrentBalance, toUpdatedBalance)
 
    return nil
}
 
// add two number checking for overflow
func add(b int, q int) (int, error) {
 
    // Check overflow
    var sum int
    sum = q + b
 
    if (sum < q || sum < b) == (b >= 0 && q >= 0) {
        return 0, fmt.Errorf("Math: addition overflow occurred %d + %d", b, q)
    }
 
    return sum, nil
}
 
// Checks that contract options have been already initialized
func checkInitialized(sdk kalpsdk.TransactionContextInterface) (bool, error) {
    tokenName, err := sdk.GetState(nameKey)
    if err != nil {
        return false, fmt.Errorf("failed to get token name: %v", err)
    }
 
    if tokenName == nil {
        return false, nil
    }
 
    return true, nil
}
 
// sub two number checking for overflow
func sub(b int, q int) (int, error) {
 
    // sub two number checking
    if q <= 0 {
        return 0, fmt.Errorf("Error: the subtraction number is %d, it should be greater than 0", q)
    }
    if b < q {
        return 0, fmt.Errorf("Error: the number %d is not enough to be subtracted by %d", b, q)
    }
    var diff int
    diff = b - q
 
    return diff, nil
}