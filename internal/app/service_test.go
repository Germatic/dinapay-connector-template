package app

import (
	"context"
	"testing"
	"time"

	"github.com/Germatic/dinapay-connector-template/internal/adapters/memory"
	"github.com/Germatic/dinapay-connector-template/internal/core"
)

type providerStub struct {
	creates   int
	createErr error
	recovered core.ProviderPayment
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
func (providerStub) CreateRefund(context.Context, string, core.CreateRefundCommand) (core.ProviderRefund, error) {
	return core.ProviderRefund{}, core.ErrUnsupported
}
func (providerStub) RecoverRefund(context.Context, string, core.CreateRefundCommand) (core.ProviderRefund, error) {
	return core.ProviderRefund{}, core.ErrUnsupported
}
func (providerStub) GetRefund(context.Context, string, string) (core.ProviderRefund, error) {
	return core.ProviderRefund{}, core.ErrUnsupported
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
