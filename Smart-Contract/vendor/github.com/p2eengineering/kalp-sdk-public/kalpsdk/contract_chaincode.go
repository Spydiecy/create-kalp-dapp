package kalpsdk

import (
	//Third party Libs
	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/hyperledger/fabric-protos-go/peer"
)

// ChaincodeStubInterface is used by deployable chaincode apps to access and
// modify their ledgers
type ChaincodeStubInterface interface {
	shim.ChaincodeStubInterface
}

// ContractChaincode a struct to meet the chaincode interface and provide routing of calls to contracts
type ContractChaincode struct {
	contractapi.ContractChaincode
}

// Init is called during Instantiate transaction after the chaincode container
// has been established for the first time, passes off details of the request to Invoke
// for handling the request if a function name is passed, otherwise returns shim.Success
func (kc *ContractChaincode) Init(stub ChaincodeStubInterface) peer.Response {
	return kc.ContractChaincode.Init(stub)
}

// Invoke is called to update or query the ledger in a proposal transaction.
func (kc *ContractChaincode) Invoke(stub ChaincodeStubInterface) peer.Response {
	return kc.ContractChaincode.Invoke(stub)
}

// Start starts the chaincode in the fabric network
func (kc *ContractChaincode) Start() error {
	// If Start() is called, we assume this is a standalone chaincode and set
	// up formatted logging.
	setupChaincodeLogging()
	return kc.ContractChaincode.Start()
}
