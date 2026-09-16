# Failure mapping

Every terminal provider failure must produce two separate values:

- `Failure`: provider-neutral and built from the versioned Dinapay catalog;
- `ProviderFailure`: sanitized native evidence for internal support only.

The template pins the contract module in `go.mod`. Mapping is local and has no
runtime dependency. A connector may upgrade the contract version and deploy
independently from every other connector.

Rules:

1. Use `NewPayment`, `NewRefund` or `NewPayout`; do not invent public code strings.
2. Do not copy provider identity, codes or messages into public `Failure`.
3. Translate useful native fields into canonical Dinapay field paths.
4. Map unknown native codes to `unknown_error` and preserve them internally.
5. Test every known mapping and the unknown fallback.
6. Propose a catalog change only when existing codes cannot describe a
   materially different client-visible outcome.
7. Do not expose internal retry or reconciliation policy publicly.

Refund observations follow the same boundary. A terminal `rejected` or
`failed` result must include a `Failure` validated by `ValidRefund` and may
include a `ProviderFailure` for internal evidence. Pending and confirmed
refunds must not carry a failure. The orchestrator persists both objects,
publishes only `Failure`, and delays public `failed` until any required balance
compensation has completed.
