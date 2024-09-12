package kalpsdk

import (
	//Standard Libs
	"encoding/base64"
	"fmt"
	"strconv"
	"strings"

	//Third party Libs
	"google.golang.org/protobuf/types/known/timestamppb"
)

// GetChannelName retrieves the name of the channel associated with the transaction context.
// It returns the channel name as a string and an error if the channel ID is empty or retrieval fails.
func (ctx *TransactionContext) GetChannelName() (string, error) {
	// Get the channel ID using the transaction context's stub.
	channelID := ctx.GetStub().GetChannelID()
	// Check if the channel ID is empty.
	if channelID == "" {
		// If the channel ID is empty, return an error indicating the failure to retrieve the channel name.
		return "", fmt.Errorf("failed to get channelName: %v", channelID)
	}
	// If the channel ID is not empty, return it as the channel name along with no errors.
	return channelID, nil
}

// GetKYC checks if a user has completed KYC on our network by invoking the KycExists function
// on the kyc chaincode for the given user ID in the universalkyc channel.
//
// Parameters:
//   - userId: The ID of the user to check for KYC completion.
//
// Returns:
//   - bool: A boolean value indicating whether the user has completed KYC.
//   - error: An error if the operation fails.
func (ctx *TransactionContext) GetKYC(userId string) (bool, error) {
	// Set the function name and chaincode name for the cross-chaincode invocation.
	crossCCFunc := "KycExists"
	crossCCName := "kyc"

	// Call the GetChannelName function to obtain the channel name.
	channelName, err := ctx.GetChannelName()
	if err != nil {
		return false, fmt.Errorf("failed to get channel name: %s", err.Error())
	}

	// Set the parameters for the KycExists function.
	params := []string{crossCCFunc, userId}

	// Convert the parameters to byte arrays.
	queryArgs := make([][]byte, len(params))
	for i, arg := range params {
		queryArgs[i] = []byte(arg)
	}

	// Invoke the KycExists function on the kyc chaincode in the channel.
	response := ctx.GetStub().InvokeChaincode(crossCCName, queryArgs, channelName)

	// Check if the response status is not 200 OK.
	if response.Status != 200 {
		// Return an error with a descriptive message.
		return false, fmt.Errorf("failed to query kyc chaincode for user %s. Got status %d and error message: %s", userId, response.Status, response.Payload)
	}

	// Convert the response payload to a boolean and return it.
	return strconv.ParseBool(string(response.Payload))
}

// GetUserID retrieves the name of the minter from the CA certificate embedded in the client identity.
// It returns the user ID extracted from the client identity and an error if there was a failure in
// reading or extracting the user ID.
//
// Returns:
//   - string: The user ID extracted from the client identity.
//   - error: An error if there was a failure in reading or extracting the user ID.
func (ctx *TransactionContext) GetUserID() (string, error) {
	b64ID, err := ctx.GetClientIdentity().GetID()
	if err != nil {
		return "", fmt.Errorf("failed to read clientID: %v", err)
	}

	decodeID, err := base64.StdEncoding.DecodeString(b64ID)
	if err != nil {
		return "", fmt.Errorf("failed to base64 decode clientID: %v", err)
	}

	completeID := string(decodeID)
	userID := completeID[(strings.Index(completeID, "x509::CN=") + 9):strings.Index(completeID, ",")]
	return userID, nil
}

// GetState retrieves the value of the specified `key` from the ledger.
// It should be noted that GetState does not read data from the writeset,
// which contains data that has not been committed to the ledger yet.
// In other words, GetState only retrieves data that has been previously
// committed and is considered part of the current state of the ledger.
//
// If the `key` does not exist in the state database, the function returns
// (nil, nil) indicating the absence of the key in the ledger.
//
// Parameters:
//   - key: The key of the data to retrieve from the ledger.
//
// Returns:
//   - []byte: The value associated with the specified `key` in the ledger.
//   - error: An error if there was a failure in retrieving the data.
func (ctx *TransactionContext) GetState(key string) ([]byte, error) {
	return ctx.GetStub().GetState(key)
}

// SetEvent allows the chaincode to set an event on the response to the
// proposal to be included as part of a transaction. The event will be
// available within the transaction in the committed block regardless of the
// validity of the transaction.
//
// Only a single event can be included in a transaction, and must originate
// from the outer-most invoked chaincode in chaincode-to-chaincode scenarios.
// The marshaled ChaincodeEvent will be available in the transaction's ChaincodeAction.events field.
//
// Parameters:
//   - name (string): The name of the event to be set.
//   - payload ([]byte): The payload data associated with the event.
//
// Returns:
//   - error: An error if there was a failure in setting the event.
func (ctx *TransactionContext) SetEvent(name string, payload []byte) error {
	return ctx.GetStub().SetEvent(name, payload)
}

// GetTxID returns the transaction ID of the transaction proposal. The transaction ID is
// unique per transaction and per client. It can be used to uniquely identify and track a specific
// transaction within the blockchain network.
//
// Returns:
//   - string: The transaction ID of the transaction proposal.
func (ctx *TransactionContext) GetTxID() string {
	return ctx.GetStub().GetTxID()
}

// GetChannelID returns the channel the proposal is sent to for chaincode to process.
// This would be the channel_id of the transaction proposal
//
// Returns:
//   - string: The channel ID of the transaction proposal.
func (ctx *TransactionContext) GetChannelID() string {
	return ctx.GetStub().GetChannelID()
}

// GetStateByPartialCompositeKey queries the state in the ledger based on a given partial composite key.
// This function returns an iterator which can be used to iterate over all composite keys whose prefix matches
// the given partial composite key. However, if the number of matching composite keys is greater than the totalQueryLimit
// (defined in core.yaml), this iterator cannot be used to fetch all matching keys, and the results will be limited by the totalQueryLimit.
//
// The `objectType` and attributes are expected to have only valid UTF-8 strings and should not contain U+0000 (nil byte)
// and U+10FFFF (biggest and unallocated code point). See the related functions SplitCompositeKey and CreateCompositeKey
// for working with composite keys.
//
// Call Close() on the returned StateQueryIteratorInterface object when done.
//
// Please note that the query is re-executed during the validation phase to ensure the result set has not changed
// since transaction endorsement (phantom reads detected). This function should be used only for a partial composite key.
// For a full composite key, an iterator with an empty response would be returned.
//
// Parameters:
//   - objectType: The object type portion of the partial composite key.
//   - keys: The attributes that make up the partial composite key.
//
// Returns:
//   - StateQueryIteratorInterface: An iterator that can be used to iterate over the composite keys matching the partial composite key.
//   - error: An error if there was a failure in retrieving the composite keys by partial composite key.
func (ctx *TransactionContext) GetStateByPartialCompositeKey(objectType string, keys []string) (StateQueryIteratorInterface, error) {
	return ctx.GetStub().GetStateByPartialCompositeKey(objectType, keys)
}

// GetStateByRange returns a range iterator over a set of keys in the ledger.
// The iterator can be used to iterate over all keys between the startKey (inclusive) and endKey (exclusive).
// However, if the number of keys between startKey and endKey is greater than the totalQueryLimit (defined in core.yaml),
// this iterator cannot be used to fetch all keys, and the results will be capped by the totalQueryLimit.
// The keys are returned by the iterator in lexical order. Note that startKey and endKey can be an empty string,
// which implies an unbounded range query on the start or end.
//
// Call Close() on the returned StateQueryIteratorInterface object when done.
//
// Please note that the query is re-executed during the validation phase to ensure the result set has not changed
// since transaction endorsement (phantom reads detected).
//
// Parameters:
//   - startKey: The start key (inclusive) of the range.
//   - endKey: The end key (exclusive) of the range.
//
// Returns:
//   - StateQueryIteratorInterface: An iterator that can be used to iterate over the keys in the specified range.
//   - error: An error if there was a failure in retrieving the keys by range.
func (ctx *TransactionContext) GetStateByRange(startKey string, endKey string) (StateQueryIteratorInterface, error) {
	return ctx.GetStub().GetStateByRange(startKey, endKey)
}

// GetQueryResult performs a "rich" query against a state database that supports rich query,
// such as CouchDB. The query string is provided in the native syntax of the underlying state database.
// An iterator is returned, which can be used to iterate over all keys in the query result set.
//
// The query is NOT re-executed during the validation phase, and phantom reads are not detected.
// This means that other committed transactions may have added, updated, or removed keys that
// impact the result set, but this would not be detected at validation/commit time.
// Applications that are susceptible to this should avoid using GetQueryResult as part of transactions
// that update the ledger, and should limit its use to read-only chaincode operations.
//
// The iterator may not be able to fetch all keys in the query result set if the number of keys exceeds
// the totalQueryLimit defined in the core.yaml configuration file. The results will be limited by the totalQueryLimit.
//
// Parameters:
//   - query: The query string in the native syntax of the underlying state database.
//
// Returns:
//   - StateQueryIteratorInterface: An iterator that can be used to iterate over all keys in the query result set.
//   - error: An error if there was a failure in performing the query.
func (ctx *TransactionContext) GetQueryResult(query string) (StateQueryIteratorInterface, error) {
	return ctx.GetStub().GetQueryResult(query)
}

// GetHistoryForKey returns a history of key values across time.
// For each historic key update, the historic value and associated
// transaction ID and timestamp are returned. The timestamp is the
// timestamp provided by the client in the proposal header.
//
// GetHistoryForKey requires the peer configuration core.ledger.history.enableHistoryDatabase to be true.
// The query is NOT re-executed during the validation phase, and phantom reads are not detected.
// This means that other committed transactions may have updated the key concurrently, impacting the result set,
// but this would not be detected at validation/commit time.
// Applications that are susceptible to this should avoid using GetHistoryForKey as part of transactions
// that update the ledger, and should limit its use to read-only chaincode operations.
//
// Starting in Fabric v2.0, the GetHistoryForKey chaincode API will return results from newest to oldest
// in terms of ordered transaction height (block height and transaction height within a block).
// This allows applications to efficiently iterate through the top results to understand recent changes to a key.
//
// Parameters:
//   - key: The key for which to retrieve the history.
//
// Returns:
//   - HistoryQueryIteratorInterface: An iterator that can be used to iterate over the history of key values.
//   - error: An error if there was a failure in retrieving the history.
func (ctx *TransactionContext) GetHistoryForKey(key string) (HistoryQueryIteratorInterface, error) {
	return ctx.GetStub().GetHistoryForKey(key)
}

// CreateCompositeKey combines the given `objectType` and `attributes` to form a composite key.
// The `objectType` and `attributes` should be valid UTF-8 strings and must not contain the U+0000 (nil byte)
// or U+10FFFF (largest and unallocated code point) characters.
//
// Parameters:
//   - objectType (string): The type of the object for which the composite key is being created.
//   - attributes ([]string): The attributes used to form the composite key.
//
// Returns:
//   - string: The composite key formed by combining the `objectType` and `attributes`.
//   - error: An error if there was a failure in creating the composite key.
func (ctx *TransactionContext) CreateCompositeKey(objectType string, attributes []string) (string, error) {
	return ctx.GetStub().CreateCompositeKey(objectType, attributes)
}

// SplitCompositeKey splits the specified key into attributes on which the
// composite key was formed. Composite keys found during range queries
// or partial composite key queries can therefore be split into their
// composite parts.
// Parameters:
//   - compositeKey (string): The composite key which is to be splited.
//
// Returns:
//   - string: The composite key formed by combining the `objectType` and `attributes`.
//   - []string: list of individual keys after successful split of composite key.
//   - error: An error if there was a failure in split the composite key.
func (ctx *TransactionContext) SplitCompositeKey(compositeKey string) (string, []string, error) {
	return ctx.GetStub().SplitCompositeKey(compositeKey)
}

// GetTxTimestamp returns the timestamp when the transaction was created. The timestamp
// is extracted from the transaction's ChannelHeader, which ensures that it indicates the
// client's timestamp and has the same value across all endorsers.
//
// Returns:
//   - *timestamppb.Timestamp: The timestamp of the transaction.
//   - error: An error if there was a failure in retrieving the timestamp.
func (ctx *TransactionContext) GetTxTimestamp() (*timestamppb.Timestamp, error) {
	return ctx.GetStub().GetTxTimestamp()
}

// GetFunctionAndParameters returns the function name and parameters extracted from the transaction proposal.
// The first argument in the transaction proposal is considered as the function name, while the rest of the arguments
// are treated as parameters and returned as a string array.
//
// Note: Only use GetFunctionAndParameters if the client passes arguments intended to be used as strings.
//
// Returns:
//   - string: The function name extracted from the transaction proposal.
//   - []string: The parameters extracted from the transaction proposal as a string array.
func (ctx *TransactionContext) GetFunctionAndParameters() (string, []string) {
	return ctx.GetStub().GetFunctionAndParameters()
}
