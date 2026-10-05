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
9. Replace every placeholder in `contracts/` and `deploy/`. CI rejects
   placeholders automatically after the Go module stops being the template.

`RecoverPayment`, `RecoverRefund` and `RecoverPayout` are mandatory safety hooks for
the operations a connector declares. They query the
provider using the stable canonical `operationId` after an ambiguous timeout;
they must not create a new provider operation.

A connector may be pay-in-only, payout-only, or support both. Its capability
manifest is authoritative: declare only implemented operations and return
`unsupported` for operations the provider does not support. Adding a provider
must not require changes to Dinapay V2 or the router.

### Optional deterministic sandbox simulation

The template exposes the Connector Contract V1 simulation extension without
making it mandatory for real provider adapters. A connector that supports it:

1. implements `core.SimulationAdapter`;
2. advertises supported scenarios on the exact capability through
   `capability.simulation` and gives that capability a stable `capabilityId`;
3. declares the versioned behavior in `contracts/simulation-profile.json`;
4. persists idempotency using the simulation store operations; and
5. publishes an ordinary normalized provider event after accepting a command.

`POST /v1/payments/{providerPaymentId}/simulate` is registered only when
`DINARIA_ENVIRONMENT=sandbox`. Production returns 404 even when the adapter
implements simulation. The adapter must still reject undeclared scenarios and
delays greater than its advertised maximum.

The same provider connector handles real and simulated execution; there is no
generic simulator connector. Control Plane selects the mode on the sandbox
provider connection, and Routing keeps selecting the normal provider. A shared
simulation kit may provide canonical QR, bank-transfer, redirect, cash and card
generators, while this connector owns its profile and provider-specific rules.
Provider failures must never cause an automatic fallback to simulation.

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
- Optional, sandbox-only deterministic simulation extension.
- Provider adapter, store and event-publisher ports.
- Idempotency reservation and replay behavior.
- Webhook body limits and inbox deduplication hook.
- HTTP publisher for normalized provider events.
- Health/readiness endpoints and bounded HTTP server timeouts.
- Canonical `/version` build identity and low-cost Prometheus `/metrics`.
- Declarative capability and canonical-error manifests.
- Reproducible build metadata and release-component generation.
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
make verify test build
SERVICE_TOKEN=local-secret PROVIDER_NAME=my-provider go run ./cmd/connector
```

Then query:

```sh
curl http://localhost:8092/health
curl http://localhost:8092/version
curl http://localhost:8092/metrics
curl http://localhost:8092/v1/capabilities \
  -H 'Authorization: Bearer local-secret'
```

## Environment

- `SERVICE_TOKEN`: required internal bearer token.
- `DINARIA_ENVIRONMENT`: runtime environment; the simulation route exists only
  when its value is exactly `sandbox`.
- `DINAPAY_V2_URL`: orchestrator base URL; defaults to `http://localhost:8112`.
- `PROVIDER_NAME`: placeholder capability name; replace with provider config.
- `PORT`: defaults to `8092`.

See the contracts repository for versioning, idempotency, event, tracing and
financial recovery requirements.

## Release metadata

The build injects service, repository, version, commit, timestamp and
environment into `/version`. Generate the fragment consumed by Dinaria's
atomic release pipeline with:

```sh
make release-component \
  PROVIDER=my-provider \
  REPOSITORY=github.com/Germatic/dinapay-connector-my-provider \
  PROCESS_NAME=dinapay-connector-my-provider-v2 \
  PORT=8092
```

The generated SHA-256 identifies the exact artifact approved in sandbox and
later promoted to production. Do not edit that digest manually.
