package kalpsdk

import (
	//Standard Libs
	"fmt"

	//Third party Libs
	"golang.org/x/exp/slices"
)

// ValidateCreateTokenTransaction checks if the contract has been initialized, if the operator is authorized
// to create the token, and if the token with the given ID and document type is already minted. Returns an error
// if any of the checks fail, or nil if the transaction is valid.
func (ctx *TransactionContext) ValidateCreateTokenTransaction(id string, docType string, account []string) error {
	// // Check if contract has been initialized.
	// initialized, err := kapsutils.CheckInitialized(ctx)
	// if err != nil {
	// 	return fmt.Errorf("failed to check if contract is already initialized: %v", err)
	// }
	// if !initialized {
	// 	return fmt.Errorf("contract options need to be set before calling any function, call Initialize() to initialize contract")
	// }

	// Check if operator is authorized to create token.
	operator, err := ctx.GetUserID()
	if err != nil {
		return fmt.Errorf("failed to get client id: %v", err)
	}
	if !slices.Contains(account, operator) {
		return fmt.Errorf("only the asset owner is allowed to initiate create transaction")
	}

	// Check if token is already minted.
	minted, err := IsMinted(ctx, id, docType)
	if err != nil {
		return fmt.Errorf("failed to check if token is already minted: %v", err)
	}
	if minted {
		return fmt.Errorf("the token with ID '%v' is already minted", id)
	}

	return nil
}

// IsMinted checks whether a token with the specified ID and document type is already minted or not.
// Returns true if minted, false otherwise.
func IsMinted(sdk *TransactionContext, id string, docType string) (bool, error) {
	queryString := fmt.Sprintf(`{"selector": {"id": "%s", "docType": "%s"}}`, id, docType)

	resultsIterator, err := sdk.GetStub().GetQueryResult(queryString)
	if err != nil {
		return false, fmt.Errorf("failed to get query result from the world state: %v", err)
	}

	return resultsIterator.HasNext(), nil
}
