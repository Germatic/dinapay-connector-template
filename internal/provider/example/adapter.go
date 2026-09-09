package example

import (
	"context"
	"fmt"

	"github.com/Germatic/dinapay-connector-template/internal/core"
)

// Adapter is deliberately non-functional. Replace this package with the provider integration.
type Adapter struct{ Name string }

func (a Adapter) Capabilities(context.Context) core.Capabilities {
	return core.Capabilities{Provider: a.Name, ContractVersion: "1", Capabilities: []core.Capability{{
		Operation: "payment", Countries: []string{"XX"}, Currencies: []string{"XXX"},
		PaymentMethods: []string{"bank_transfer"}, Rails: []string{"replace_me"},
	}}}
}
func (Adapter) CreatePayment(context.Context, core.CreatePaymentCommand) (core.ProviderPayment, error) {
	return core.ProviderPayment{}, fmt.Errorf("%w: implement provider CreatePayment", core.ErrUnsupported)
}
func (Adapter) RecoverPayment(context.Context, core.CreatePaymentCommand) (core.ProviderPayment, error) {
	return core.ProviderPayment{}, core.ErrUnsupported
}
func (Adapter) GetPayment(context.Context, string, string) (core.ProviderPayment, error) {
	return core.ProviderPayment{}, core.ErrUnsupported
}
func (Adapter) CancelPayment(context.Context, core.CancelPaymentCommand, string) (core.ProviderPayment, error) {
	return core.ProviderPayment{}, core.ErrUnsupported
}
func (Adapter) CreateRefund(context.Context, string, core.CreateRefundCommand) (core.ProviderRefund, error) {
	return core.ProviderRefund{}, core.ErrUnsupported
}
func (Adapter) RecoverRefund(context.Context, string, core.CreateRefundCommand) (core.ProviderRefund, error) {
	return core.ProviderRefund{}, core.ErrUnsupported
}
func (Adapter) GetRefund(context.Context, string, string) (core.ProviderRefund, error) {
	return core.ProviderRefund{}, core.ErrUnsupported
}
func (Adapter) ParseWebhook(context.Context, core.RawWebhook) ([]core.ProviderEvent, error) {
	return nil, core.ErrUnsupported
}
