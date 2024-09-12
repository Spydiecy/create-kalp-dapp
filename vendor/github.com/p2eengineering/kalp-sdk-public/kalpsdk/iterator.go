package kalpsdk

import (
	//Third party Libs
	"github.com/hyperledger/fabric-chaincode-go/shim"
)

// StateQueryIteratorInterface allows a chaincode to iterate over a set of key/value pairs returned by range and execute query.
type StateQueryIteratorInterface interface {
	shim.StateQueryIteratorInterface
}

// HistoryQueryIteratorInterface allows a chaincode to iterate over a set of
// key/value pairs returned by a history query.
type HistoryQueryIteratorInterface interface {
	shim.HistoryQueryIteratorInterface
}
