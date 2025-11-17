#!/usr/bin/env bash
set -euo pipefail

# Requires VAULT_ADDR and VAULT_TOKEN set for a privileged token.
vault kv put secret/orders-api APP_SECRET="super-secret-value"
