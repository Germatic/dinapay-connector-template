package core

import (
	"context"
	"errors"
)

var (
	ErrNotFound    = errors.New("not found")
	ErrConflict    = errors.New("idempotency conflict")
	ErrInProgress  = errors.New("operation in progress")
	ErrUnsupported = errors.New("operation unsupported")
	ErrRejected    = errors.New("provider rejected operation")
	ErrUnavailable = errors.New("provider unavailable")
)

// ProviderAdapter is the only provider-specific implementation required by the template.
type ProviderAdapter interface {
	Capabilities(context.Context) Capabilities
	CreatePayment(context.Context, CreatePaymentCommand) (ProviderPayment, error)
	RecoverPayment(context.Context, CreatePaymentCommand) (ProviderPayment, error)
	GetPayment(context.Context, string, string) (ProviderPayment, error)
	CancelPayment(context.Context, CancelPaymentCommand, string) (ProviderPayment, error)
	CreateRefund(context.Context, string, CreateRefundCommand) (ProviderRefund, error)
	RecoverRefund(context.Context, string, CreateRefundCommand) (ProviderRefund, error)
	GetRefund(context.Context, string, string) (ProviderRefund, error)
	CreatePayout(context.Context, CreatePayoutCommand) (ProviderPayout, error)
	RecoverPayout(context.Context, CreatePayoutCommand) (ProviderPayout, error)
	GetPayout(context.Context, string, string) (ProviderPayout, error)
	CancelPayout(context.Context, CancelPayoutCommand, string) (ProviderPayout, error)
	ParseWebhook(context.Context, RawWebhook) ([]ProviderEvent, error)
}

type Store interface {
	ReservePayment(context.Context, string, string, []byte) (*ProviderPayment, error)
	CompletePayment(context.Context, string, ProviderPayment) error
	FailPayment(context.Context, string, string) error
	FindPayment(context.Context, string, string) (ProviderPayment, error)
	ReserveRefund(context.Context, string, string, []byte) (*ProviderRefund, error)
	CompleteRefund(context.Context, string, ProviderRefund) error
	FailRefund(context.Context, string, string) error
	FindRefund(context.Context, string, string) (ProviderRefund, error)
	ReservePayout(context.Context, string, string, []byte) (*ProviderPayout, error)
	CompletePayout(context.Context, string, ProviderPayout) error
	FailPayout(context.Context, string, string) error
	FindPayout(context.Context, string, string) (ProviderPayout, error)
	RecordEvent(context.Context, ProviderEvent, []byte) (bool, error)
	CompleteEvent(context.Context, string) error
}

type EventPublisher interface {
	Publish(context.Context, ProviderEvent) error
}
