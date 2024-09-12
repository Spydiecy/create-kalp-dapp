package main

import (
	"log"

	"github.com/p2eengineering/kalp-sdk-public/kalpsdk"
)

func main() {

	contract := kalpsdk.Contract{IsPayableContract: false}

	contract.Logger = kalpsdk.NewLogger()
	chaincode, err := kalpsdk.NewChaincode(&SmartContract{contract})
	contract.Logger.Info("My KAPL SDK sm4")

	// Create a new instance of your KalpContractChaincode with your smart contract
	// chaincode, err := kalpsdk.NewChaincode(&SmartContract{kalpsdk.Contract{IsPayableContract: true}})
	// kalpsdk.NewLogger()
	if err != nil {
		log.Panicf("Error creating KalpContractChaincode: %v", err)
	}

	// Start the chaincode
	if err := chaincode.Start(); err != nil {
		log.Panicf("Error starting chaincode: %v", err)
	}
}