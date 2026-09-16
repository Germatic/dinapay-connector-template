package example

import contract "github.com/Germatic/dinapay-contracts/go/connectorcontract/failures"

type nativeFailure struct {
	Code     string
	Message  string
	Field    string
	Category string
}

// mapPayoutFailure is the provider-owned translation boundary. Replace the
// example codes with documented native provider failures. Never create a
// public Dinapay code in this connector.
func mapPayoutFailure(input nativeFailure) (contract.Failure, contract.ProviderFailure) {
	provider := contract.ProviderFailure{
		Code:    input.Code,
		Message: input.Message,
		Details: map[string]any{"field": input.Field, "category": input.Category},
	}
	switch input.Code {
	case "EXAMPLE_DESTINATION_REJECTED":
		return contract.NewPayout(contract.PayoutDestinationRejected), provider
	default:
		return contract.NewPayout(contract.PayoutUnknownError), provider
	}
}

func mapRefundFailure(input nativeFailure) (contract.Failure, contract.ProviderFailure) {
	provider := contract.ProviderFailure{
		Code:    input.Code,
		Message: input.Message,
		Details: map[string]any{"field": input.Field, "category": input.Category},
	}
	switch input.Code {
	case "EXAMPLE_REFUND_REJECTED":
		return contract.NewRefund(contract.RefundRejected), provider
	default:
		return contract.NewRefund(contract.RefundUnknownError), provider
	}
}
