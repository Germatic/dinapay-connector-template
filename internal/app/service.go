package app

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"errors"

	"github.com/Germatic/dinapay-connector-template/internal/core"
)

type Service struct {
	provider  core.ProviderAdapter
	store     core.Store
	publisher core.EventPublisher
}

func New(provider core.ProviderAdapter, store core.Store, publisher core.EventPublisher) *Service {
	return &Service{provider: provider, store: store, publisher: publisher}
}

func hash(v any) []byte { raw, _ := json.Marshal(v); sum := sha256.Sum256(raw); return sum[:] }

func (s *Service) Capabilities(ctx context.Context) core.Capabilities {
	return s.provider.Capabilities(ctx)
}

func (s *Service) CreatePayment(ctx context.Context, key string, cmd core.CreatePaymentCommand) (core.ProviderPayment, bool, error) {
	if key == "" || cmd.OperationID == "" || cmd.TransactionID == "" || cmd.ProviderConnectionID == "" || cmd.Amount == "" || cmd.Currency == "" {
		return core.ProviderPayment{}, false, core.ErrRejected
	}
	existing, err := s.store.ReservePayment(ctx, cmd.ProviderConnectionID, key, hash(cmd))
	if errors.Is(err, core.ErrInProgress) {
		result, recoverErr := s.provider.RecoverPayment(ctx, cmd)
		if recoverErr != nil {
			return core.ProviderPayment{}, false, recoverErr
		}
		if recoverErr = s.store.CompletePayment(ctx, key, result); recoverErr != nil {
			return core.ProviderPayment{}, false, recoverErr
		}
		return result, true, nil
	}
	if err != nil {
		return core.ProviderPayment{}, false, err
	}
	if existing != nil {
		return *existing, true, nil
	}
	result, err := s.provider.CreatePayment(ctx, cmd)
	if err != nil {
		// A timeout is ambiguous: keep the reservation so the same operation cannot be duplicated.
		if !errors.Is(err, core.ErrUnavailable) {
			_ = s.store.FailPayment(context.WithoutCancel(ctx), cmd.ProviderConnectionID, key)
		}
		return core.ProviderPayment{}, false, err
	}
	if err = s.store.CompletePayment(ctx, key, result); err != nil {
		return core.ProviderPayment{}, false, err
	}
	return result, false, nil
}

func (s *Service) GetPayment(ctx context.Context, connection, id string) (core.ProviderPayment, error) {
	result, err := s.provider.GetPayment(ctx, connection, id)
	if errors.Is(err, core.ErrUnsupported) {
		return s.store.FindPayment(ctx, connection, id)
	}
	return result, err
}

func (s *Service) CancelPayment(ctx context.Context, id string, cmd core.CancelPaymentCommand) (core.ProviderPayment, error) {
	return s.provider.CancelPayment(ctx, cmd, id)
}

func (s *Service) CreateRefund(ctx context.Context, providerPaymentID, key string, cmd core.CreateRefundCommand) (core.ProviderRefund, bool, error) {
	if key == "" || cmd.OperationID == "" || cmd.RefundID == "" || cmd.TransactionID == "" || cmd.ProviderConnectionID == "" || cmd.Amount == "" || cmd.Currency == "" {
		return core.ProviderRefund{}, false, core.ErrRejected
	}
	existing, err := s.store.ReserveRefund(ctx, cmd.ProviderConnectionID, key, hash(cmd))
	if errors.Is(err, core.ErrInProgress) {
		result, recoverErr := s.provider.RecoverRefund(ctx, providerPaymentID, cmd)
		if recoverErr != nil {
			return core.ProviderRefund{}, false, recoverErr
		}
		if recoverErr = s.store.CompleteRefund(ctx, key, result); recoverErr != nil {
			return core.ProviderRefund{}, false, recoverErr
		}
		return result, true, nil
	}
	if err != nil {
		return core.ProviderRefund{}, false, err
	}
	if existing != nil {
		return *existing, true, nil
	}
	result, err := s.provider.CreateRefund(ctx, providerPaymentID, cmd)
	if err != nil {
		if !errors.Is(err, core.ErrUnavailable) {
			_ = s.store.FailRefund(context.WithoutCancel(ctx), cmd.ProviderConnectionID, key)
		}
		return core.ProviderRefund{}, false, err
	}
	if err = s.store.CompleteRefund(ctx, key, result); err != nil {
		return core.ProviderRefund{}, false, err
	}
	return result, false, nil
}

func (s *Service) GetRefund(ctx context.Context, connection, id string) (core.ProviderRefund, error) {
	result, err := s.provider.GetRefund(ctx, connection, id)
	if errors.Is(err, core.ErrUnsupported) {
		return s.store.FindRefund(ctx, connection, id)
	}
	return result, err
}

func (s *Service) HandleWebhook(ctx context.Context, raw core.RawWebhook) error {
	events, err := s.provider.ParseWebhook(ctx, raw)
	if err != nil {
		return err
	}
	for _, event := range events {
		fresh, err := s.store.RecordEvent(ctx, event, raw.Body)
		if err != nil {
			return err
		}
		if fresh && s.publisher != nil {
			if err = s.publisher.Publish(ctx, event); err != nil {
				return err
			}
		}
	}
	return nil
}
