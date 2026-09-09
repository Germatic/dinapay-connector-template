package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/Germatic/dinapay-connector-template/internal/adapters/memory"
	"github.com/Germatic/dinapay-connector-template/internal/adapters/publisher"
	"github.com/Germatic/dinapay-connector-template/internal/app"
	"github.com/Germatic/dinapay-connector-template/internal/provider/example"
	"github.com/Germatic/dinapay-connector-template/internal/transport/httpapi"
)

func main() {
	port := env("PORT", "8092")
	token := os.Getenv("SERVICE_TOKEN")
	orchestrator := env("DINAPAY_V2_URL", "http://localhost:8112")
	providerName := env("PROVIDER_NAME", "replace-me")
	if token == "" {
		log.Fatal("SERVICE_TOKEN is required")
	}
	provider := example.Adapter{Name: providerName}
	store := memory.New()
	events := publisher.New(orchestrator, token)
	service := app.New(provider, store, events)
	server := &http.Server{Addr: ":" + port, Handler: httpapi.New(service, token), ReadHeaderTimeout: 5 * time.Second, ReadTimeout: 15 * time.Second, WriteTimeout: 30 * time.Second, IdleTimeout: 60 * time.Second}
	log.Printf("connector template listening on %s", server.Addr)
	log.Fatal(server.ListenAndServe())
}
func env(name, fallback string) string {
	if value := os.Getenv(name); value != "" {
		return value
	}
	return fallback
}
