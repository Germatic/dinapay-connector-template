package httpapi

import (
	"crypto/subtle"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"

	"github.com/Germatic/dinapay-connector-template/internal/app"
	"github.com/Germatic/dinapay-connector-template/internal/core"
)

type Server struct {
	service      *app.Service
	serviceToken string
}

func New(service *app.Service, token string) http.Handler {
	s := &Server{service: service, serviceToken: token}
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, _ *http.Request) { write(w, 200, map[string]string{"status": "up"}) })
	mux.HandleFunc("GET /ready", func(w http.ResponseWriter, _ *http.Request) { write(w, 200, map[string]string{"status": "ready"}) })
	mux.HandleFunc("GET /v1/capabilities", s.auth(s.capabilities))
	mux.HandleFunc("POST /v1/payments", s.auth(s.createPayment))
	mux.HandleFunc("GET /v1/payments/{providerPaymentId}", s.auth(s.getPayment))
	mux.HandleFunc("POST /v1/payments/{providerPaymentId}/cancel", s.auth(s.cancelPayment))
	mux.HandleFunc("POST /v1/payments/{providerPaymentId}/refunds", s.auth(s.createRefund))
	mux.HandleFunc("GET /v1/refunds/{providerRefundId}", s.auth(s.getRefund))
	mux.HandleFunc("POST /v1/payouts", s.auth(s.createPayout))
	mux.HandleFunc("GET /v1/payouts/{providerPayoutId}", s.auth(s.getPayout))
	mux.HandleFunc("POST /v1/payouts/{providerPayoutId}/cancel", s.auth(s.cancelPayout))
	mux.HandleFunc("POST /webhooks/{connectionId}", s.webhook)
	return mux
}
func (s *Server) auth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
		if s.serviceToken == "" || subtle.ConstantTimeCompare([]byte(token), []byte(s.serviceToken)) != 1 {
			problem(w, 401, "unauthorized", "invalid service credentials")
			return
		}
		next(w, r)
	}
}
func (s *Server) capabilities(w http.ResponseWriter, r *http.Request) {
	write(w, 200, s.service.Capabilities(r.Context()))
}
func (s *Server) createPayment(w http.ResponseWriter, r *http.Request) {
	var cmd core.CreatePaymentCommand
	if err := decode(w, r, &cmd); err != nil {
		problem(w, 400, "invalid_request", err.Error())
		return
	}
	result, replayed, err := s.service.CreatePayment(r.Context(), r.Header.Get("Idempotency-Key"), cmd)
	if err != nil {
		mapError(w, err)
		return
	}
	if replayed {
		w.Header().Set("Idempotent-Replayed", "true")
		write(w, 200, result)
		return
	}
	write(w, 201, result)
}
func (s *Server) getPayment(w http.ResponseWriter, r *http.Request) {
	result, err := s.service.GetPayment(r.Context(), r.Header.Get("Provider-Connection-Id"), r.PathValue("providerPaymentId"))
	if err != nil {
		mapError(w, err)
		return
	}
	write(w, 200, result)
}
func (s *Server) cancelPayment(w http.ResponseWriter, r *http.Request) {
	var cmd core.CancelPaymentCommand
	if err := decode(w, r, &cmd); err != nil {
		problem(w, 400, "invalid_request", err.Error())
		return
	}
	result, err := s.service.CancelPayment(r.Context(), r.PathValue("providerPaymentId"), cmd)
	if err != nil {
		mapError(w, err)
		return
	}
	write(w, 200, result)
}
func (s *Server) createRefund(w http.ResponseWriter, r *http.Request) {
	var cmd core.CreateRefundCommand
	if err := decode(w, r, &cmd); err != nil {
		problem(w, 400, "invalid_request", err.Error())
		return
	}
	result, replayed, err := s.service.CreateRefund(r.Context(), r.PathValue("providerPaymentId"), r.Header.Get("Idempotency-Key"), cmd)
	if err != nil {
		mapError(w, err)
		return
	}
	if replayed {
		w.Header().Set("Idempotent-Replayed", "true")
		write(w, 200, result)
		return
	}
	write(w, 201, result)
}
func (s *Server) getRefund(w http.ResponseWriter, r *http.Request) {
	result, err := s.service.GetRefund(r.Context(), r.Header.Get("Provider-Connection-Id"), r.PathValue("providerRefundId"))
	if err != nil {
		mapError(w, err)
		return
	}
	write(w, 200, result)
}
func (s *Server) createPayout(w http.ResponseWriter, r *http.Request) {
	var cmd core.CreatePayoutCommand
	if err := decode(w, r, &cmd); err != nil {
		problem(w, 400, "invalid_request", err.Error())
		return
	}
	result, replayed, err := s.service.CreatePayout(r.Context(), r.Header.Get("Idempotency-Key"), cmd)
	if err != nil {
		mapError(w, err)
		return
	}
	if replayed {
		w.Header().Set("Idempotent-Replayed", "true")
		write(w, 200, result)
		return
	}
	write(w, 201, result)
}
func (s *Server) getPayout(w http.ResponseWriter, r *http.Request) {
	result, err := s.service.GetPayout(r.Context(), r.Header.Get("Provider-Connection-Id"), r.PathValue("providerPayoutId"))
	if err != nil {
		mapError(w, err)
		return
	}
	write(w, 200, result)
}
func (s *Server) cancelPayout(w http.ResponseWriter, r *http.Request) {
	var cmd core.CancelPayoutCommand
	if err := decode(w, r, &cmd); err != nil {
		problem(w, 400, "invalid_request", err.Error())
		return
	}
	result, err := s.service.CancelPayout(r.Context(), r.PathValue("providerPayoutId"), cmd)
	if err != nil {
		mapError(w, err)
		return
	}
	write(w, 200, result)
}
func (s *Server) webhook(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(http.MaxBytesReader(w, r.Body, 1<<20))
	if err != nil {
		problem(w, 400, "invalid_request", err.Error())
		return
	}
	headers := map[string][]string(r.Header.Clone())
	err = s.service.HandleWebhook(r.Context(), core.RawWebhook{ConnectionID: r.PathValue("connectionId"), Headers: headers, Body: body})
	if err != nil {
		mapError(w, err)
		return
	}
	write(w, 200, map[string]string{"status": "accepted"})
}
func decode(w http.ResponseWriter, r *http.Request, target any) error {
	d := json.NewDecoder(http.MaxBytesReader(w, r.Body, 1<<20))
	d.DisallowUnknownFields()
	return d.Decode(target)
}
func write(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}
func problem(w http.ResponseWriter, status int, code, message string) {
	write(w, status, map[string]any{"error": map[string]string{"code": code, "message": message}})
}
func mapError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, core.ErrNotFound):
		problem(w, 404, "not_found", err.Error())
	case errors.Is(err, core.ErrConflict):
		problem(w, 409, "idempotency_conflict", err.Error())
	case errors.Is(err, core.ErrUnsupported):
		problem(w, 422, "unsupported", err.Error())
	case errors.Is(err, core.ErrUnavailable):
		problem(w, 503, "provider_unavailable", err.Error())
	default:
		problem(w, 422, "provider_rejected", err.Error())
	}
}
