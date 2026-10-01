#!/usr/bin/env bash
set -Eeuo pipefail

root=$(git rev-parse --show-toplevel)
cd "$root"

command -v jq >/dev/null || { printf 'ERROR: jq is required\n' >&2; exit 1; }

jq -e '
  .schemaVersion == "1" and
  (.provider | type == "string" and length > 0) and
  (.contractVersion | type == "string" and length > 0) and
  (.capabilities | type == "array" and length > 0) and
  ([.capabilities[] |
    (.operation | IN("payment", "refund", "payout")) and
    ((.dataRequirements // []) | type == "array") and
    ([((.dataRequirements // [])[]) |
      (.path | type == "string" and contains(".")) and
      (.presence | IN("required", "conditional"))
    ] | all)
  ] | all)
' contracts/capabilities.json >/dev/null

jq -e '
  .schemaVersion == "1" and
  (.provider | type == "string" and length > 0) and
  (.mappings | type == "array") and
  ([.mappings[] | (.operation | IN("payment", "refund", "payout")) and (.providerCode | length > 0) and (.canonicalCode | length > 0)] | all)
' contracts/canonical-errors.json >/dev/null

jq -e '
  (.service | length > 0) and (.repository | length > 0) and
  (.contractVersion | length > 0) and (.commit | length > 0) and
  (.sha256 | length > 0) and (.artifactName | length > 0) and
  (.processManager | IN("pm2", "systemd", "none")) and
  (.healthUrl | length > 0)
' deploy/release-component.template.json >/dev/null

module=$(go list -m)
if [[ "$module" != github.com/Germatic/dinapay-connector-template ]]; then
  if grep -RniE --exclude='*.md' 'replace-me|replace_me|REPLACE_[A-Z_]+' \
    contracts deploy internal/provider README.md; then
    printf 'ERROR: generated connector still contains template placeholders\n' >&2
    exit 1
  fi
  if grep -Rn 'github.com/Germatic/dinapay-connector-template\|dinapay-connector-template' \
    go.mod internal/buildinfo Makefile Dockerfile; then
    printf 'ERROR: generated connector still contains template build identity\n' >&2
    exit 1
  fi
fi

unformatted=$(gofmt -l cmd internal)
[[ -z "$unformatted" ]] || { printf 'ERROR: gofmt required:\n%s\n' "$unformatted" >&2; exit 1; }

for endpoint in '/health' '/ready' '/version' '/metrics' '/v1/capabilities'; do
  grep -Rq "$endpoint" internal/transport/httpapi || {
    printf 'ERROR: required endpoint is missing: %s\n' "$endpoint" >&2
    exit 1
  }
done

printf 'Connector contract verification passed.\n'
