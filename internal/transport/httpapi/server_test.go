package httpapi

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Germatic/dinapay-connector-template/internal/adapters/memory"
	"github.com/Germatic/dinapay-connector-template/internal/app"
	"github.com/Germatic/dinapay-connector-template/internal/buildinfo"
	"github.com/Germatic/dinapay-connector-template/internal/provider/example"
)

func TestOperationalEndpoints(t *testing.T) {
	handler := New(app.New(example.Adapter{Name: "example"}, memory.New(), nil), "secret")

	t.Run("version exposes canonical identity", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		handler.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/version", nil))
		if recorder.Code != http.StatusOK {
			t.Fatalf("status = %d", recorder.Code)
		}
		var got buildinfo.Info
		if err := json.Unmarshal(recorder.Body.Bytes(), &got); err != nil {
			t.Fatal(err)
		}
		if got.Service == "" || got.Commit == "" || got.ContractVersion == "" {
			t.Fatalf("incomplete build identity: %+v", got)
		}
	})

	t.Run("health includes build identity", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		handler.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/health", nil))
		if recorder.Code != http.StatusOK || !strings.Contains(recorder.Body.String(), `"build"`) {
			t.Fatalf("status=%d body=%s", recorder.Code, recorder.Body.String())
		}
	})

	t.Run("metrics expose build and request counters", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		handler.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/metrics", nil))
		body := recorder.Body.String()
		if recorder.Code != http.StatusOK || !strings.Contains(body, "dinapay_connector_build_info") || !strings.Contains(body, "dinapay_connector_http_requests_total") {
			t.Fatalf("status=%d body=%s", recorder.Code, body)
		}
	})
}
