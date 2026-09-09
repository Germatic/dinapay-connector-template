package publisher

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/Germatic/dinapay-connector-template/internal/core"
)

type HTTP struct {
	baseURL, token string
	client         *http.Client
}

func New(baseURL, token string) *HTTP {
	return &HTTP{baseURL: strings.TrimRight(baseURL, "/"), token: token, client: &http.Client{Timeout: 10 * time.Second}}
}

func (p *HTTP) Publish(ctx context.Context, event core.ProviderEvent) error {
	body, err := json.Marshal(event)
	if err != nil {
		return err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, p.baseURL+"/internal/v1/provider-events", bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+p.token)
	if trace := traceFromContext(ctx); trace != "" {
		req.Header.Set("traceparent", trace)
	}
	resp, err := p.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode/100 != 2 {
		return fmt.Errorf("publish provider event: HTTP %d", resp.StatusCode)
	}
	return nil
}

// Replace this hook with the project's tracing library while preserving W3C propagation.
func traceFromContext(context.Context) string { return "" }
