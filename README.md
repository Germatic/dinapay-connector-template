# Dinapay Provider Connector Template

Starting point for an independently deployable provider connector implementing
the contract in `Germatic/dinapay-contracts/openapi/connector-v1.json`.

The template keeps provider code behind `core.ProviderAdapter`. A provider team
must not import or modify Dinapay V2 or routing code.

## Start a connector

1. Create a repository from this template and change the Go module path.
2. Replace `internal/provider/example` with the real provider adapter.
3. Implement a durable PostgreSQL `core.Store`; the included memory store is
   for tests and local scaffolding only.
4. Configure provider connections and credentials from a secret manager.
5. Declare only capabilities that are actually implemented.
6. Implement and verify the provider webhook under `/webhooks/{connectionId}`.
7. Add a reconciler that queries non-terminal operations and publishes the
   same normalized events as webhooks.
8. Run the acceptance matrix in `dinapay-contracts/docs/connector-implementation-guide.md`.

`RecoverPayment`, `RecoverRefund` and `RecoverPayout` are mandatory safety hooks for
the operations a connector declares. They query the
provider using the stable canonical `operationId` after an ambiguous timeout;
they must not create a new provider operation.

A connector may be pay-in-only, payout-only, or support both. Its capability
manifest is authoritative: declare only implemented operations and return
`unsupported` for operations the provider does not support. Adding a provider
must not require changes to Dinapay V2 or the router.

Provider-specific code should normally be limited to:

```text
internal/provider/<provider>/
├── adapter.go
├── client.go
├── mapping.go
├── webhook.go
└── *_test.go
```

## Generic components included

- Connector V1 HTTP surface and service-token authentication.
- Canonical payment, refund, payout, binding and event models.
- Provider adapter, store and event-publisher ports.
- Idempotency reservation and replay behavior.
- Webhook body limits and inbox deduplication hook.
- HTTP publisher for normalized provider events.
- Health/readiness endpoints and bounded HTTP server timeouts.
- Local memory adapter and an idempotency test.
- Docker image and CI test workflow.

## Required production work

This repository intentionally cannot process a provider payment as cloned. The
example adapter returns `unsupported`. Before deployment, replace the memory
store with durable storage, add migrations, implement provider credentials,
webhook verification, reconciliation, structured telemetry and readiness
checks. Never acknowledge an unverified provider webhook.

## Local verification

```sh
make test
SERVICE_TOKEN=local-secret PROVIDER_NAME=my-provider go run ./cmd/connector
```

Then query:

```sh
curl http://localhost:8092/health
curl http://localhost:8092/v1/capabilities \
  -H 'Authorization: Bearer local-secret'
```

## Environment

- `SERVICE_TOKEN`: required internal bearer token.
- `DINAPAY_V2_URL`: orchestrator base URL; defaults to `http://localhost:8112`.
- `PROVIDER_NAME`: placeholder capability name; replace with provider config.
- `PORT`: defaults to `8092`.

See the contracts repository for versioning, idempotency, event, tracing and
financial recovery requirements.
