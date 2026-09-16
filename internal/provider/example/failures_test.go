package example

import (
	"testing"

	contract "github.com/Germatic/dinapay-contracts/go/connectorcontract/failures"
)

func TestPayoutFailureMappingUsesContractCatalog(t *testing.T) {
	public, internal := mapPayoutFailure(nativeFailure{Code: "EXAMPLE_DESTINATION_REJECTED", Message: "native detail"})
	if public.Code != string(contract.PayoutDestinationRejected) || !contract.ValidPayout(public) {
		t.Fatalf("public=%#v", public)
	}
	if internal.Code != "EXAMPLE_DESTINATION_REJECTED" || internal.Message != "native detail" {
		t.Fatalf("internal=%#v", internal)
	}
}

func TestUnknownPayoutFailureFallsBack(t *testing.T) {
	public, internal := mapPayoutFailure(nativeFailure{Code: "NEW_NATIVE_CODE"})
	if public.Code != string(contract.PayoutUnknownError) || internal.Code != "NEW_NATIVE_CODE" {
		t.Fatalf("public=%#v internal=%#v", public, internal)
	}
}
