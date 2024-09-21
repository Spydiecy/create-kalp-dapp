// Package kalpsdk provides a set of APIs for developing applications on the Kalptantra blockchain network.
// It aims to simplify the process of interacting with the blockchain network and offers a range of functionalities
// to enhance the development experience.
//
// The kalpsdk package aims to provide developers with a streamlined and efficient development experience when building
// applications on the Kalptantra blockchain network. By offering a range of convenient functionalities and integrating
// with existing Hyperledger Fabric packages, it simplifies the implementation of smart contracts and interaction with
// the blockchain network.
//
// Note: This overview provides a high-level description of the kalpsdk package and its main components. For detailed
// information on specific types, methods, and usage examples, please refer to the package documentation and code comments.
package kalpsdk

import (
	//Standard Libs
	"encoding/json"
	"fmt"
	"time"

	//Third party Libs
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/hyperledger/fabric-contract-api-go/metadata"
)

// Contract defines functions for setting and getting before, after and unknown transactions
// and name. Can be embedded in structs to quickly ensure their definition meets the
// ContractInterface.
type Contract struct {
	Logger            *ChaincodeLogger
	IsPayableContract bool
	contractapi.Contract
}

// PaymentTracker represents the payment tracking information associated with a transaction on the Kalptantra blockchain network.
// The struct is used for storing and retrieving payment-related data.
type PaymentMetaData struct {
	Amount                 float64   `json:”amount”`                        // Amount of the payment
	CurrencyCode           string    `json:”currencyCode”`                  // Currency Code of the payment
	paymentTimestamp       time.Time `json:"paymentTimestamp"`              // Timestamp of the Payment
	ApplicationReferenceId string    `json:"applicationReferenceId"`        // ID of the application or uuid of Payment Engine
	IsPaymentEngineUsed    bool      `json:"isPaymentEngineUsed,omitempty"` // If Payment Engine used or not, default value should be true
}

type PaymentTracker struct {
	TransactionId        string          `json:"transactionId"`        // The ID of the transaction.
	DocType              string          `json:"DocType"`              // The type of the document it must be PAYMENT-INFO.
	PaymentTransactionID string          `json:”paymentTransactionId"` // The reference number of the payment.
	PaymentGatewayName   string          `json:"paymentGatewayName`    // The Name of the payment gateway.
	PaymentMetaData      PaymentMetaData `json:"paymentMetaData"`      // Additional metadata related to the payment.
	AssetInfo            interface{}     `json:"assetInfo,omitempty”`  // Information about the associated asset.
	AssetId              string          `json:"id,omitempty”`
	AssetDocType         string          `json:"docType,omitempty”`
}

// NewChaincode creates a new chaincode using the contracts passed as arguments. Each of the passed contracts
// is parsed, and details about their structure and public functions are stored for use by the chaincode.
// The function will return an error if the contracts are invalid, such as having public functions that take illegal types.
// A system contract is added to the chaincode, which provides functionality for getting the metadata of the chaincode.
// The generated metadata is a JSON-formatted MetadataContractChaincode containing each contract's name and details of
// its public functions and the types they tahe actual functions; instead, they are labeled as param0, param1, ..., paramN.
// If a file nake in and return. The parameter names recorded in the metadata do not
// match those used in tmed contract-metadata/metadata.json exists, it will overwrite the generated metadata. The contents of this
// file must validate against the schema.
//
// By default, the transaction serializer for the contract is set to JSONSerializer. This can be updated by changing
// the TransactionSerializer property.
//
// Parameters:
//   - contracts: The contracts implementing the chaincode functionality.
//
// Returns:
//   - *ContractChaincode: The initialized ContractChaincode instance.
//   - error: An error if there was a failure in creating the chaincode.
func NewChaincode(contracts ...contractapi.ContractInterface) (*ContractChaincode, error) {
	chaincode, err := contractapi.NewChaincode(contracts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create chaincode: %v", err)
	}

	return &ContractChaincode{ContractChaincode: *chaincode}, nil
}

// GetInfo returns the information about the contract that can be used in metadata.
// It retrieves the InfoMetadata object associated with the contract.
//
// Returns:
//   - metadata.InfoMetadata: The metadata containing information about the contract.
func (c *Contract) GetInfo() metadata.InfoMetadata {
	return c.Info
}

// GetUnknownTransaction returns the current unknown transaction set for the contract.
// It retrieves the unknownTransaction interface{} object associated with the contract.
//
// Returns:
//   - interface{}: The unknown transaction set for the contract, which may be nil.
func (c *Contract) GetUnknownTransaction() interface{} {
	return c.UnknownTransaction
}

// GetBeforeTransaction returns the current set beforeTransaction, may be nil
func (c *Contract) GetBeforeTransaction() interface{} {
	return c.BeforeTransaction
}

// GetAfterTransaction returns the current set afterTransaction, which is a function to be executed after each transaction.
// The returned function takes two parameters: the transaction context and the result of the transaction.
// It performs post-transaction operations such as payment processing or data persistence.
// If the contract is a payable contract, the payment operation is implemented based on the transaction type.
func (c *Contract) GetAfterTransaction() interface{} {
	fmt.Println("GetAfterTransaction Called once while install chaincode")
	c.Logger = NewLogger()
	setupChaincodeLogging()

	assetMap := make(map[string]interface{})
	// afterFunction is an anonymous function that will be executed after each transaction
	afterFunction := func(ctx TransactionContextInterface) error {
		c.Logger.Println("After Transaction:", ctx.GetTxID(), "IsPayable is:", c.IsPayableContract)
		if c.IsPayableContract {

			// Retrieve the function name and arguments of the executed transaction
			fnName, args := ctx.GetFunctionAndParameters()
			c.Logger.Println("args:", fnName, args)

			var paymentTracker PaymentTracker
			inputData := args[len(args)-1]

			err := json.Unmarshal([]byte(inputData), &paymentTracker)
			if err != nil {
				return err
			}

			err = json.Unmarshal([]byte(inputData), &assetMap)
			if err != nil {
				return fmt.Errorf("failed to parse input data of asset: %v", err)
			}

			if !c.CheckPaymentDetails(paymentTracker) {
				return fmt.Errorf("payment transaction does not have valid amount or currencycode!")
			}

			// Update the paymentTracker fields
			paymentTracker.DocType = "PAYMENT-INFO"
			paymentTracker.TransactionId = ctx.GetTxID()

			if Id, err := assetMap["id"]; err {
				paymentTracker.AssetId = Id.(string)
			}
			if DocType, err := assetMap["docType"]; err {
				paymentTracker.AssetDocType = DocType.(string)
			}

			// Marshal the updated paymentTracker object into JSON
			paymentData, err := json.Marshal(paymentTracker)
			if err != nil {
				return err
			}

			// Put the paymentData to the ledger using PutStateWithKyc
			err = ctx.PutStateWithKYC(paymentTracker.TransactionId, paymentData)
			if err != nil {
				return err
			}
			c.Logger.Println("Successfully triggered the After operations for the transaction")
		}
		return nil
	}
	return afterFunction
}

// GetName returns the name of the contract.
// GetName retrieves the name associated with the contract.
//
// Returns:
//   - string: The name of the contract.

func (c *Contract) CheckPaymentDetails(PaymentDetails PaymentTracker) bool {
	if PaymentDetails.PaymentMetaData.Amount < 1 || PaymentDetails.PaymentMetaData.CurrencyCode == "" {
		return false
	}
	return true
}

func (c *Contract) GetName() string {
	return c.Name
}

// GetTransactionContextHandler returns the current transaction context handler set for the contract.
// If no transaction context handler has been set, a new TransactionContext will be returned.
//
// Returns:
//   - contractapi.SettableTransactionContextInterface: The transaction context handler for the contract.
func (c *Contract) GetTransactionContextHandler() contractapi.SettableTransactionContextInterface {
	if c.TransactionContextHandler == nil {
		return new(TransactionContext)
	}

	return c.TransactionContextHandler
}
