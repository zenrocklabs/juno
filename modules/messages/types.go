package messages

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/zenrocklabs/juno/types"
)

// MessageAddressesParser represents a function that extracts all the
// involved addresses from a provided message (both accounts and validators)
type MessageAddressesParser = func(tx *types.Transaction) ([]string, error)

// CosmosMessageAddressesParser represents a MessageAddressesParser that parses a
// Chain message and returns all the involved addresses (both accounts and validators)
var CosmosMessageAddressesParser = DefaultMessagesParser

// DefaultMessagesParser represents the default messages parser that simply returns the list
// of all the signers of a message
func DefaultMessagesParser(tx *types.Transaction) ([]string, error) {
	return parseAddressesFromEvents(tx), nil
}

// function to remove duplicate values
func removeDuplicates(s []string) []string {
	bucket := make(map[string]bool)
	result := []string{}
	for _, str := range s {
		if _, ok := bucket[str]; !ok {
			bucket[str] = true
			result = append(result, str)
		}
	}
	return result
}

func parseAddressesFromEvents(tx *types.Transaction) []string {
	addresses := []string{}
	for _, event := range tx.Events {
		for _, attribute := range event.Attributes {
			// Try parsing the address as a validator address
			validatorAddress, _ := sdk.ValAddressFromBech32(attribute.Value)
			if validatorAddress != nil {
				addresses = append(addresses, validatorAddress.String())
			}

			// Try parsing the address as an account address
			accountAddress, err := sdk.AccAddressFromBech32(attribute.Value)
			if err != nil {
				// Skip if the address is not an account address
				continue
			}

			addresses = append(addresses, accountAddress.String())
		}
	}

	return removeDuplicates(addresses)
}
