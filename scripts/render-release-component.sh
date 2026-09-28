#!/usr/bin/env bash
set -Eeuo pipefail

required=(PROVIDER REPOSITORY ARTIFACT PROCESS_NAME PORT)
for name in "${required[@]}"; do
  [[ -n "${!name:-}" ]] || { printf 'ERROR: %s is required\n' "$name" >&2; exit 1; }
done

command -v jq >/dev/null || { printf 'ERROR: jq is required\n' >&2; exit 1; }
[[ -f "$ARTIFACT" ]] || { printf 'ERROR: artifact not found: %s\n' "$ARTIFACT" >&2; exit 1; }

if command -v sha256sum >/dev/null; then
  sha=$(sha256sum "$ARTIFACT" | awk '{print $1}')
else
  sha=$(shasum -a 256 "$ARTIFACT" | awk '{print $1}')
fi

commit=${COMMIT:-$(git rev-parse HEAD)}
contract_version=${CONTRACT_VERSION:-connector-v1}
process_manager=${PROCESS_MANAGER:-pm2}
health_status=${HEALTH_STATUS:-200}

jq -n \
  --arg service "dinapay-connector-$PROVIDER" \
  --arg repository "$REPOSITORY" \
  --arg contractVersion "$contract_version" \
  --arg commit "$commit" \
  --arg sha256 "$sha" \
  --arg artifactName "$(basename "$ARTIFACT")" \
  --arg processManager "$process_manager" \
  --arg processName "$PROCESS_NAME" \
  --arg healthUrl "http://127.0.0.1:$PORT/health" \
  --argjson healthStatus "$health_status" \
  '{service:$service,repository:$repository,contractVersion:$contractVersion,commit:$commit,sha256:$sha256,artifactName:$artifactName,processManager:$processManager,processName:$processName,healthUrl:$healthUrl,healthStatus:$healthStatus}'
