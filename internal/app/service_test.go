package app

import (
	"context"
	"testing"
	"time"

	"github.com/Germatic/dinapay-connector-template/internal/adapters/memory"
	"github.com/Germatic/dinapay-connector-template/internal/core"
	"github.com/Germatic/dinapay-connector-template/internal/provider/example"
)

type providerStub struct {
	creates         int
	createErr       error
	recovered       core.ProviderPayment
	payoutCreates   int
	payoutCreateErr error
	payoutRecovered core.ProviderPayout
	simulations     int
}

func (providerStub) Capabilities(context.Context) core.Capabilities {
	return core.Capabilities{Provider: "test", ContractVersion: "1"}
}
func (p *providerStub) CreatePayment(_ context.Context, c core.CreatePaymentCommand) (core.ProviderPayment, error) {
	p.creates++
	if p.createErr != nil {
		return core.ProviderPayment{}, p.createErr
	}
	return core.ProviderPayment{TransactionID: c.TransactionID, Provider: "test", ProviderConnectionID: c.ProviderConnectionID, ProviderPaymentID: "provider-1", Status: "pending", ObservedAt: time.Now()}, nil
}
func (p *providerStub) RecoverPayment(_ context.Context, c core.CreatePaymentCommand) (core.ProviderPayment, error) {
	if p.recovered.ProviderPaymentID == "" {
		return core.ProviderPayment{}, core.ErrNotFound
	}
	return p.recovered, nil
}
func (providerStub) GetPayment(context.Context, string, string) (core.ProviderPayment, error) {
	return core.ProviderPayment{}, core.ErrUnsupported
}
func (providerStub) CancelPayment(context.Context, core.CancelPaymentCommand, string) (core.ProviderPayment, error) {
	return core.ProviderPayment{}, core.ErrUnsupported
}
func (p *providerStub) SimulatePayment(_ context.Context, _ string, c core.SimulatePaymentCommand) (core.SimulationAccepted, error) {
	p.simulations++
	return core.SimulationAccepted{SimulationID: "simulation-1", TransactionID: c.TransactionID, Scenario: c.Scenario, Status: "accepted", ScheduledAt: time.Now()}, nil
}
func (p *providerStub) RecoverSimulation(_ context.Context, _ string, c core.SimulatePaymentCommand) (core.SimulationAccepted, error) {
	return core.SimulationAccepted{SimulationID: "simulation-1", TransactionID: c.TransactionID, Scenario: c.Scenario, Status: "accepted", ScheduledAt: time.Now()}, nil
}
func (providerStub) CreateRefund(context.Context, string, core.CreateRefundCommand) (core.ProviderRefund, error) {
	return core.ProviderRefund{}, core.ErrUnsupported
}
func (providerStub) RecoverRefund(context.Context, string, core.CreateRefundCommand) (core.ProviderRefund, error) {
	return core.ProviderRefund{}, core.ErrUnsupported
}
func (providerStub) GetRefund(context.Context, string, string) (core.ProviderRefund, error) {
	return core.ProviderRefund{}, core.ErrUnsupported
}
func (p *providerStub) CreatePayout(_ context.Context, c core.CreatePayoutCommand) (core.ProviderPayout, error) {
	p.payoutCreates++
	if p.payoutCreateErr != nil {
		return core.ProviderPayout{}, p.payoutCreateErr
	}
	return core.ProviderPayout{PayoutID: c.PayoutID, Provider: "test", ProviderConnectionID: c.ProviderConnectionID, ProviderPayoutID: "provider-payout-1", Status: "processing", Source: c.Source, ObservedAt: time.Now()}, nil
}
func (p *providerStub) RecoverPayout(context.Context, core.CreatePayoutCommand) (core.ProviderPayout, error) {
	if p.payoutRecovered.ProviderPayoutID == "" {
		return core.ProviderPayout{}, core.ErrNotFound
	}
	return p.payoutRecovered, nil
}
func (providerStub) GetPayout(context.Context, string, string) (core.ProviderPayout, error) {
	return core.ProviderPayout{}, core.ErrUnsupported
}
func (providerStub) CancelPayout(context.Context, core.CancelPayoutCommand, string) (core.ProviderPayout, error) {
	return core.ProviderPayout{}, core.ErrUnsupported
}

func TestCreatePaymentRecoversAmbiguousTimeout(t *testing.T) {
	p := &providerStub{createErr: core.ErrUnavailable}
	s := New(p, memory.New(), nil)
	cmd := core.CreatePaymentCommand{OperationID: "op-timeout", TransactionID: "tx-timeout", ProviderConnectionID: "connection-1", Amount: "1.00", Currency: "USD"}
	if _, _, err := s.CreatePayment(context.Background(), "key-timeout", cmd); err != core.ErrUnavailable {
		t.Fatalf("first err=%v", err)
	}
	p.createErr = nil
	p.recovered = core.ProviderPayment{TransactionID: cmd.TransactionID, Provider: "test", ProviderConnectionID: cmd.ProviderConnectionID, ProviderPaymentID: "recovered-1", Status: "pending", ObservedAt: time.Now()}
	result, replayed, err := s.CreatePayment(context.Background(), "key-timeout", cmd)
	if err != nil || !replayed || result.ProviderPaymentID != "recovered-1" || p.creates != 1 {
		t.Fatalf("result=%#v replay=%v creates=%d err=%v", result, replayed, p.creates, err)
	}
}
func (providerStub) ParseWebhook(context.Context, core.RawWebhook) ([]core.ProviderEvent, error) {
	return nil, nil
}

func TestCreatePaymentIdempotency(t *testing.T) {
	p := &providerStub{}
	s := New(p, memory.New(), nil)
	cmd := core.CreatePaymentCommand{OperationID: "op-1", TransactionID: "tx-1", ProviderConnectionID: "connection-1", Amount: "1.00", Currency: "USD"}
	first, replay, err := s.CreatePayment(context.Background(), "key-1", cmd)
	if err != nil || replay {
		t.Fatalf("first: replay=%v err=%v", replay, err)
	}
	second, replay, err := s.CreatePayment(context.Background(), "key-1", cmd)
	if err != nil || !replay || first.ProviderPaymentID != second.ProviderPaymentID || p.creates != 1 {
		t.Fatalf("replay=%v creates=%d err=%v", replay, p.creates, err)
	}
	cmd.Amount = "2.00"
	_, _, err = s.CreatePayment(context.Background(), "key-1", cmd)
	if err != core.ErrConflict {
		t.Fatalf("conflict err=%v", err)
	}
}

func TestSimulatePaymentIdempotency(t *testing.T) {
	p := &providerStub{}
	s := New(p, memory.New(), nil)
	cmd := core.SimulatePaymentCommand{OperationID: "operation-1", TransactionID: "transaction-1", ProviderConnectionID: "connection-1", Scenario: "payment.confirmed"}
	first, replayed, err := s.SimulatePayment(context.Background(), "provider-payment-1", "simulation-key-1", cmd)
	if err != nil || replayed || first.Status != "accepted" {
		t.Fatalf("first=%#v replayed=%v err=%v", first, replayed, err)
	}
	second, replayed, err := s.SimulatePayment(context.Background(), "provider-payment-1", "simulation-key-1", cmd)
	if err != nil || !replayed || second.SimulationID != first.SimulationID || p.simulations != 1 {
		t.Fatalf("second=%#v replayed=%v calls=%d err=%v", second, replayed, p.simulations, err)
	}
	cmd.Scenario = "payment.expired"
	if _, _, err = s.SimulatePayment(context.Background(), "provider-payment-1", "simulation-key-1", cmd); err != core.ErrConflict {
		t.Fatalf("conflict err=%v", err)
	}
}

func TestSimulationIsOptional(t *testing.T) {
	s := New(example.Adapter{Name: "example"}, memory.New(), nil)
	cmd := core.SimulatePaymentCommand{OperationID: "operation-1", TransactionID: "transaction-1", ProviderConnectionID: "connection-1", Scenario: "payment.confirmed"}
	if _, _, err := s.SimulatePayment(context.Background(), "provider-payment-1", "simulation-key-1", cmd); err != core.ErrUnsupported {
		t.Fatalf("err=%v", err)
	}
}

func TestCreatePayoutIdempotency(t *testing.T) {
	p := &providerStub{}
	s := New(p, memory.New(), nil)
	cmd := payoutCommand("op-payout-1", "payout-1")
	first, replayed, err := s.CreatePayout(context.Background(), "payout-key-1", cmd)
	if err != nil || replayed {
		t.Fatalf("first: replay=%v err=%v", replayed, err)
	}
	second, replayed, err := s.CreatePayout(context.Background(), "payout-key-1", cmd)
	if err != nil || !replayed || first.ProviderPayoutID != second.ProviderPayoutID || p.payoutCreates != 1 {
		t.Fatalf("replay=%v creates=%d err=%v", replayed, p.payoutCreates, err)
	}
	cmd.Source.Amount = "11.00"
	if _, _, err = s.CreatePayout(context.Background(), "payout-key-1", cmd); err != core.ErrConflict {
		t.Fatalf("conflict err=%v", err)
	}
}

func TestCreatePayoutRecoversAmbiguousTimeout(t *testing.T) {
	p := &providerStub{payoutCreateErr: core.ErrUnavailable}
	s := New(p, memory.New(), nil)
	cmd := payoutCommand("op-payout-timeout", "payout-timeout")
	if _, _, err := s.CreatePayout(context.Background(), "payout-key-timeout", cmd); err != core.ErrUnavailable {
		t.Fatalf("first err=%v", err)
	}
	p.payoutCreateErr = nil
	p.payoutRecovered = core.ProviderPayout{PayoutID: cmd.PayoutID, Provider: "test", ProviderConnectionID: cmd.ProviderConnectionID, ProviderPayoutID: "recovered-payout-1", Status: "processing", Source: cmd.Source, ObservedAt: time.Now()}
	result, replayed, err := s.CreatePayout(context.Background(), "payout-key-timeout", cmd)
	if err != nil || !replayed || result.ProviderPayoutID != "recovered-payout-1" || p.payoutCreates != 1 {
		t.Fatalf("result=%#v replay=%v creates=%d err=%v", result, replayed, p.payoutCreates, err)
	}
}

func payoutCommand(operationID, payoutID string) core.CreatePayoutCommand {
	return core.CreatePayoutCommand{
		OperationID: operationID, PayoutID: payoutID, ProviderConnectionID: "connection-1",
		Source:      core.Money{Amount: "10.00", Currency: "USD"},
		Destination: core.PayoutDestination{Country: "VE", Currency: "VES", Beneficiary: map[string]any{"firstName": "Ana"}, Rail: map[string]any{"type": "bank_account"}},
	}
}
