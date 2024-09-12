package kalpsdk

import (
	//Custom Build Libs
	res "github.com/p2eengineering/kalp-sdk-public/response"

	//Third party Libs
	"github.com/hyperledger/fabric-chaincode-go/pkg/cid"
	"github.com/hyperledger/fabric-chaincode-go/shim"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type TransactionContextInterface interface {
	// PutStateWithKYC puts the specified `key` and `value` into the transaction's
	// writeset as a data-write proposal, only if the user has completed KYC.
	// If the user has not completed KYC, an error is returned.
	// The data is not immediately written to the ledger, but instead, it becomes part of
	// the transaction proposal and will be committed if the transaction is validated successfully.
	PutStateWithKYC(key string, value []byte) error

	// PutStateWithoutKYC puts the specified `key` and `value` into the transaction's
	// writeset as a data-write proposal without requiring KYC verification.
	// The data is not immediately written to the ledger, but instead, it becomes part of
	// the transaction proposal and will be committed if the transaction is validated successfully.
	// This function does not enforce KYC restrictions, allowing any user to write data to
	// the ledger without completing KYC. Use this function with caution as it may bypass
	// security and compliance measures. It is recommended to use PutStateWithKYC instead,
	// which enforces KYC restrictions and provides an additional layer of security.
	PutStateWithoutKYC(key string, value []byte) error

	// GetKYC checks if a user has completed KYC on our network by invoking the KycExists function on the kyc chaincode
	// for the given user ID in the universalkyc channel.
	GetKYC(userId string) (bool, error)

	// PutKYC records the KYC information associated with a user.
	// It invokes the "kyc" chaincode with the specified parameters to create a KYC record.
	// This function should only be used by administrators to create KYC records. It invokes
	// the "kyc" chaincode with the specified parameters using the "universalkyc" chaincode name.
	PutKYC(id string, kycId string, kycHash string) error

	// DelStateWithoutKYC records the specified `key` to be deleted in the writeset of
	// the transaction proposal. The `key` and its value will be deleted from
	// the ledger when the transaction is validated and successfully committed.
	// This function does not require KYC verification, allowing any user
	// to delete data from the ledger without KYC restrictions. Use this function
	// with caution as it may bypass security and compliance measures.
	// It is recommended to use DelStateWithKYC instead, which enforces KYC restrictions
	// and provides an additional layer of security.
	DelStateWithoutKYC(key string) error

	// DelStateWithKYC records the specified `key` to be deleted in the writeset of
	// the transaction proposal. The `key` and its value will be deleted from
	// the ledger when the transaction is validated and successfully committed.
	// It requires the user to have completed the KYC process before
	// deleting the state. This ensures that only authorized users can delete
	// data from the ledger, providing an additional layer of security and compliance.
	DelStateWithKYC(key string) error

	// GetState returns the value of the specified `key` from the
	// ledger. Note that GetState doesn't read data from the writeset, which
	// has not been committed to the ledger. In other words, GetState doesn't
	// consider data modified by PutState that has not been committed.
	// If the key does not exist in the state database, (nil, nil) is returned.
	GetState(key string) ([]byte, error)

	// SetEvent allows the chaincode to set an event on the response to the
	// proposal to be included as part of a transaction. The event will be
	// available within the transaction in the committed block regardless of the
	// validity of the transaction.
	// Only a single event can be included in a transaction, and must originate
	// from the outer-most invoked chaincode in chaincode-to-chaincode scenarios.
	// The marshaled ChaincodeEvent will be available in the transaction's ChaincodeAction.events field.
	SetEvent(name string, payload []byte) error

	// GetTxID returns the transaction ID of the transaction proposal. The transaction ID is
	// unique per transaction and per client. It can be used to uniquely identify and track a specific
	// transaction within the blockchain network.
	GetTxID() string

	// GetChannelID returns the channel the proposal is sent to for chaincode to process.
	// This would be the channel_id of the transaction proposal
	GetChannelID() string

	// GetUserID retrieves the name of the minter from the CA certificate embedded in the client identity.
	// It returns the user ID extracted from the client identity and an error if there was a failure in
	// reading or extracting the user ID.
	GetUserID() (string, error)

	// InvokeChaincode locally calls the specified chaincode `Invoke` using the
	// same transaction context. It allows one chaincode to invoke another chaincode
	// within the same transaction. If the called chaincode is on the same channel as
	// the calling chaincode, the called chaincode's read set and write set are added to
	// the calling transaction. If the called chaincode is on a different channel, only the
	// response from the called chaincode is returned to the calling chaincode. Any state changes
	// made by the called chaincode will not affect the ledger. Essentially, the called
	// chaincode on a different channel acts like a `Query`, and its read set and
	// write set are not applied during the state validation checks in the subsequent
	// commit phase. Only the calling chaincode's read set and write set are applied
	// to the transaction.  If the `channel` parameter is empty, it is assumed that the caller's channel is used.
	InvokeChaincode(chaincodeName string, args [][]byte, channel string) res.Response

	// CreateCompositeKey combines the given `attributes` to form a composite
	// key. The objectType and attributes are expected to have only valid utf8
	// strings and should not contain U+0000 (nil byte) and U+10FFFF
	// (biggest and unallocated code point).
	// The resulting composite key can be used as the key in PutState().
	CreateCompositeKey(objectType string, attributes []string) (string, error)

	// SplitCompositeKey splits the specified key into attributes on which the
	// composite key was formed. Composite keys found during range queries
	// or partial composite key queries can therefore be split into their
	// composite parts.
	SplitCompositeKey(compositeKey string) (string, []string, error)

	// GetStateByPartialCompositeKey queries the state in the ledger based on
	// a given partial composite key. This function returns an iterator
	// which can be used to iterate over all composite keys whose prefix matches
	// the given partial composite key. However, if the number of matching composite
	// keys is greater than the totalQueryLimit (defined in core.yaml), this iterator
	// cannot be used to fetch all matching keys (results will be limited by the totalQueryLimit).
	// The `objectType` and attributes are expected to have only valid utf8 strings and
	// should not contain U+0000 (nil byte) and U+10FFFF (biggest and unallocated code point).
	// See related functions SplitCompositeKey and CreateCompositeKey.
	// Call Close() on the returned StateQueryIteratorInterface object when done.
	// The query is re-executed during validation phase to ensure result set
	// has not changed since transaction endorsement (phantom reads detected). This function should be used only for
	// a partial composite key. For a full composite key, an iter with empty response
	// would be returned.
	GetStateByPartialCompositeKey(objectType string, keys []string) (StateQueryIteratorInterface, error)

	// GetStateByRange returns a range iterator over a set of keys in the
	// ledger. The iterator can be used to iterate over all keys
	// between the startKey (inclusive) and endKey (exclusive).
	// However, if the number of keys between startKey and endKey is greater than the
	// totalQueryLimit (defined in core.yaml), this iterator cannot be used
	// to fetch all keys (results will be capped by the totalQueryLimit).
	// The keys are returned by the iterator in lexical order. Note
	// that startKey and endKey can be empty string, which implies unbounded range
	// query on start or end.
	// Call Close() on the returned StateQueryIteratorInterface object when done.
	// The query is re-executed during validation phase to ensure result set
	// has not changed since transaction endorsement (phantom reads detected).
	GetStateByRange(startKey string, endKey string) (StateQueryIteratorInterface, error)

	// GetQueryResult performs a "rich" query against a state database. It is
	// only supported for state databases that support rich query,
	// e.g.CouchDB. The query string is in the native syntax
	// of the underlying state database. An iterator is returned
	// which can be used to iterate over all keys in the query result set.
	// However, if the number of keys in the query result set is greater than the
	// totalQueryLimit (defined in core.yaml), this iterator cannot be used
	// to fetch all keys in the query result set (results will be limited by
	// the totalQueryLimit).
	// The query is NOT re-executed during validation phase, phantom reads are
	// not detected. That is, other committed transactions may have added,
	// updated, or removed keys that impact the result set, and this would not
	// be detected at validation/commit time.  Applications susceptible to this
	// should therefore not use GetQueryResult as part of transactions that update
	// ledger, and should limit use to read-only chaincode operations.
	GetQueryResult(query string) (StateQueryIteratorInterface, error)

	// GetHistoryForKey returns a history of key values across time.
	// For each historic key update, the historic value and associated
	// transaction id and timestamp are returned. The timestamp is the
	// timestamp provided by the client in the proposal header.
	// GetHistoryForKey requires peer configuration
	// core.ledger.history.enableHistoryDatabase to be true.
	// The query is NOT re-executed during validation phase, phantom reads are
	// not detected. That is, other committed transactions may have updated
	// the key concurrently, impacting the result set, and this would not be
	// detected at validation/commit time. Applications susceptible to this
	// should therefore not use GetHistoryForKey as part of transactions that
	// update ledger, and should limit use to read-only chaincode operations.
	// Starting in Fabric v2.0, the GetHistoryForKey chaincode API
	// will return results from newest to oldest in terms of ordered transaction
	// height (block height and transaction height within block).
	// This will allow applications to efficiently iterate through the top results
	// to understand recent changes to a key.
	GetHistoryForKey(key string) (HistoryQueryIteratorInterface, error)

	// GetTxTimestamp returns the timestamp when the transaction was created. This
	// is taken from the transaction ChannelHeader, therefore it will indicate the
	// client's timestamp and will have the same value across all endorsers.
	GetTxTimestamp() (*timestamppb.Timestamp, error)

	// GetFunctionAndParameters returns the first argument as the function
	// name and the rest of the arguments as parameters in a string array.
	// Only use GetFunctionAndParameters if the client passes arguments intended
	// to be used as strings.
	GetFunctionAndParameters() (string, []string)

	// ValidateCreateTokenTransaction checks if the contract has been initialized, if the operator is authorized
	// to create the token, and if the token with the given ID and document type is already minted. Returns an error
	// if any of the checks fail, or nil if the transaction is valid.
	ValidateCreateTokenTransaction(id string, docType string, account []string) error

	// ClientIdentity represents information about the identity that submitted the transaction
	GetClientIdentity() cid.ClientIdentity
}

// TransactionContext is a basic transaction context to be used in contracts,
// containing minimal required functionality use in contracts as part of
// chaincode. Provides access to the stub and clientIdentity of a transaction.
// If a contract implements the ContractInterface using the Contract struct then
// this is the default transaction context that will be used.
type TransactionContext struct {
	stub           shim.ChaincodeStubInterface
	clientIdentity cid.ClientIdentity
}

// SetStub stores the passed stub in the transaction context
func (ctx *TransactionContext) SetStub(stub shim.ChaincodeStubInterface) {
	ctx.stub = stub
}

// SetClientIdentity stores the passed stub in the transaction context
func (ctx *TransactionContext) SetClientIdentity(ci cid.ClientIdentity) {
	ctx.clientIdentity = ci
}

// GetStub returns the current set stub
func (ctx *TransactionContext) GetStub() shim.ChaincodeStubInterface {
	return ctx.stub
}

// GetClientIdentity returns the current set client identity
func (ctx *TransactionContext) GetClientIdentity() cid.ClientIdentity {
	return ctx.clientIdentity
}
