package memory

import (
	"bytes"
	"context"
	"sync"

	"github.com/Germatic/dinapay-connector-template/internal/core"
)

type operation[T any] struct {
	hash     []byte
	complete bool
	value    T
}

type Store struct {
	mu       sync.RWMutex
	payments map[string]operation[core.ProviderPayment]
	refunds  map[string]operation[core.ProviderRefund]
	events   map[string]struct{}
}

func New() *Store {
	return &Store{payments: map[string]operation[core.ProviderPayment]{}, refunds: map[string]operation[core.ProviderRefund]{}, events: map[string]struct{}{}}
}

func scope(connection, key string) string { return connection + ":" + key }

func (s *Store) ReservePayment(_ context.Context, connection, key string, hash []byte) (*core.ProviderPayment, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	k := scope(connection, key)
	if op, ok := s.payments[k]; ok {
		if !bytes.Equal(op.hash, hash) {
			return nil, core.ErrConflict
		}
		if !op.complete {
			return nil, core.ErrInProgress
		}
		v := op.value
		return &v, nil
	}
	s.payments[k] = operation[core.ProviderPayment]{hash: append([]byte(nil), hash...)}
	return nil, nil
}
func (s *Store) CompletePayment(_ context.Context, key string, value core.ProviderPayment) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	k := scope(value.ProviderConnectionID, key)
	op, ok := s.payments[k]
	if !ok {
		return core.ErrConflict
	}
	op.complete = true
	op.value = value
	s.payments[k] = op
	return nil
}
func (s *Store) FailPayment(_ context.Context, connection, key string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.payments, scope(connection, key))
	return nil
}
func (s *Store) FindPayment(_ context.Context, connection, id string) (core.ProviderPayment, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, op := range s.payments {
		if op.complete && op.value.ProviderConnectionID == connection && op.value.ProviderPaymentID == id {
			return op.value, nil
		}
	}
	return core.ProviderPayment{}, core.ErrNotFound
}
func (s *Store) ReserveRefund(_ context.Context, connection, key string, hash []byte) (*core.ProviderRefund, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	k := scope(connection, key)
	if op, ok := s.refunds[k]; ok {
		if !bytes.Equal(op.hash, hash) {
			return nil, core.ErrConflict
		}
		if !op.complete {
			return nil, core.ErrInProgress
		}
		v := op.value
		return &v, nil
	}
	s.refunds[k] = operation[core.ProviderRefund]{hash: append([]byte(nil), hash...)}
	return nil, nil
}
func (s *Store) CompleteRefund(_ context.Context, key string, value core.ProviderRefund) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	k := scope(value.ProviderConnectionID, key)
	op, ok := s.refunds[k]
	if !ok {
		return core.ErrConflict
	}
	op.complete = true
	op.value = value
	s.refunds[k] = op
	return nil
}
func (s *Store) FailRefund(_ context.Context, connection, key string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.refunds, scope(connection, key))
	return nil
}
func (s *Store) FindRefund(_ context.Context, connection, id string) (core.ProviderRefund, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, op := range s.refunds {
		if op.complete && op.value.ProviderConnectionID == connection && op.value.ProviderRefundID == id {
			return op.value, nil
		}
	}
	return core.ProviderRefund{}, core.ErrNotFound
}
func (s *Store) RecordEvent(_ context.Context, event core.ProviderEvent, _ []byte) (bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.events[event.EventID]; ok {
		return false, nil
	}
	s.events[event.EventID] = struct{}{}
	return true, nil
}
