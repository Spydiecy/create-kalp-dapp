package kalpsdk

import (
	//Standard Libs
	"fmt"

	//Custom Build Libs
	res "github.com/p2eengineering/kalp-sdk-public/response"
)

// PutStateWithKYC puts the specified `key` and `value` into the transaction's
// writeset as a data-write proposal, only if the user has completed KYC.
// If the user has not completed KYC, an error is returned.
// The data is not immediately written to the ledger, but instead, it becomes part of
// the transaction proposal and will be committed if the transaction is validated successfully.
//
// Parameters:
//   - key: The key under which the data will be stored in the ledger.
//   - value: The data to be stored in the ledger as a byte array.
//
// Returns:
//   - error: An error if the operation fails or if the user has not completed KYC.
func (ctx *TransactionContext) PutStateWithKYC(key string, value []byte) error {
	// Get the user ID
	userID, err := ctx.GetUserID()
	if err != nil {
		return err
	}

	// Check if the user has completed KYC.
	kycCheck, err := ctx.GetKYC(userID)
	if err != nil {
		return fmt.Errorf("failed to perform KYC check for user %s. Error: %v", userID, err)
	}
	// Return false if the user has not completed KYC.
	if !kycCheck {
		return fmt.Errorf("user %s has not completed KYC", userID)
	}

	// Put the state into the transaction's writeset.
	err = ctx.GetStub().PutState(key, value)
	if err != nil {
		return err
	}

	return nil
}

// PutStateWithoutKYC puts the specified `key` and `value` into the transaction's
// writeset as a data-write proposal without requiring KYC verification.
// The data is not immediately written to the ledger, but instead, it becomes part of
// the transaction proposal and will be committed if the transaction is validated successfully.
//
// This function does not enforce KYC restrictions, allowing any user to write data to
// the ledger without completing KYC. Use this function with caution as it may bypass
// security and compliance measures. It is recommended to use PutStateWithKYC instead,
// which enforces KYC restrictions and provides an additional layer of security.
//
// Parameters:
//   - key: The key under which the data will be stored in the ledger.
//   - value: The data to be stored in the ledger as a byte array.
//
// Returns:
//   - error: An error if the operation fails.
func (ctx *TransactionContext) PutStateWithoutKYC(key string, value []byte) error {
	return ctx.GetStub().PutState(key, value)
}

// InvokeChaincode locally calls the specified chaincode `Invoke` using the
// same transaction context. It allows one chaincode to invoke another chaincode
// within the same transaction.
//
// If the called chaincode is on the same channel as the calling chaincode, the
// called chaincode's read set and write set are added to the calling transaction.
//
// If the called chaincode is on a different channel, only the response from the
// called chaincode is returned to the calling chaincode. Any state changes made
// by the called chaincode will not affect the ledger. Essentially, the called
// chaincode on a different channel acts like a `Query`, and its read set and
// write set are not applied during the state validation checks in the subsequent
// commit phase. Only the calling chaincode's read set and write set are applied
// to the transaction.
//
// If the `channel` parameter is empty, it is assumed that the caller's channel is used.
//
// Parameters:
//   - chaincodeName: The name of the chaincode to invoke.
//   - args: The arguments to pass to the invoked chaincode.
//   - channel: The channel on which the chaincode is deployed. If empty, the caller's channel is assumed.
//
// Returns:
//   - res.Response: The response from the invoked chaincode.
func (ctx *TransactionContext) InvokeChaincode(chaincodeName string, args [][]byte, channel string) res.Response {
	return res.Response{Response: ctx.GetStub().InvokeChaincode(chaincodeName, args, channel)}
}

// PutKYC records the KYC information associated with a user.
// It invokes the "kyc" chaincode with the specified parameters to create a KYC record.
// This function should only be used by administrators to create KYC records. It invokes
// the "kyc" chaincode with the specified parameters using the "universalkyc" chaincode name.
//
// Parameters:
//   - id: The ID of the user.
//   - kycId: The ID of the KYC record.
//   - kycHash: The hash value representing the KYC information.
//
// Returns:
//   - error: An error if the operation fails.
func (ctx *TransactionContext) PutKYC(id string, kycId string, kycHash string) error {
	// Prepare the parameters for the chaincode invocation
	params := []string{"CreateKyc", id, kycId, kycHash}
	invokeArgs := make([][]byte, len(params))
	for i, arg := range params {
		invokeArgs[i] = []byte(arg)
	}

	channelName, err := ctx.GetChannelName()
	if err != nil {
		return fmt.Errorf("failed to get channel name: %s", err.Error())
	}

	// Invoke the "kyc" chaincode with the specified parameters using the "kyc" chaincode name
	response := ctx.GetStub().InvokeChaincode("kyc", invokeArgs, channelName)

	// Check the response status and return an error if it is not 200 (OK)
	if response.Status != 200 {
		return fmt.Errorf("failed to query kyc chaincode. Got error: %s", response.Payload)
	}

	return nil
}

// DelStateWithoutKYC records the specified `key` to be deleted in the writeset of
// the transaction proposal. The `key` and its value will be deleted from
// the ledger when the transaction is validated and successfully committed.
//
// This function does not require KYC verification, allowing any user
// to delete data from the ledger without KYC restrictions. Use this function
// with caution as it may bypass security and compliance measures.
// It is recommended to use DelStateWithKYC instead, which enforces KYC restrictions
// and provides an additional layer of security.
//
// Parameters:
//   - key: The key of the state to be deleted.
//
// Returns:
//   - error: An error if the deletion fails.
func (ctx *TransactionContext) DelStateWithoutKYC(key string) error {
	return ctx.GetStub().DelState(key)
}

// DelStateWithKYC records the specified `key` to be deleted in the writeset of
// the transaction proposal. The `key` and its value will be deleted from
// the ledger when the transaction is validated and successfully committed.
//
// It requires the user to have completed the KYC process before
// deleting the state. This ensures that only authorized users can delete
// data from the ledger, providing an additional layer of security and compliance.
//
// Parameters:
//   - key: The key of the state to be deleted.
//
// Returns:
//   - error: An error if the deletion fails or if the user has not completed KYC.
func (ctx *TransactionContext) DelStateWithKYC(key string) error {
	// Get the user ID
	userID, err := ctx.GetUserID()
	if err != nil {
		return err
	}

	// Check if the user has completed KYC.
	kycCheck, err := ctx.GetKYC(userID)
	if err != nil {
		return fmt.Errorf("failed to perform KYC check for user %s. Error: %v", userID, err)
	}

	// Return an error if the user has not completed KYC.
	if !kycCheck {
		return fmt.Errorf("user %s has not completed KYC", userID)
	}

	// Delete the state from the world state.
	err = ctx.GetStub().DelState(key)
	if err != nil {
		return err
	}

	return nil
}
