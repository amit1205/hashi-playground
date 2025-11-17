# Sample project using Hashi Stack

This project is a minimal “project” built on **Nomad**, **Consul**, and **Vault** with two containerized services that talk to each other:

- **orders-api**
  - HTTP API that:
    - Exposes `/health` and `/config`
    - Reads a secret `APP_SECRET` from **Vault** (injected by Nomad + Vault template)
    - Reads a config message from **Consul KV** at `config/orders/message`
- **reporting-worker**
  - HTTP API that:
    - Exposes `/health` and `/call-orders`
    - Uses **Consul service discovery** to locate `orders-api`
    - Calls `orders-api`’s `/config` endpoint and returns the response

This demonstrates key HashiStack use cases:
- **Nomad**: job scheduling, Docker driver, service registration, Vault integration via `vault` + `template` stanzas
- **Consul**: service discovery, health checks, KV config
- **Vault**: central secret storage + fine-grained policy, Nomad integration (apps do not talk to Vault directly)

---

## 1. Prerequisites

You’ll need the following installed locally:

- Docker (for building/running service images)
- Go 1.22+ (for local builds; optional if you only run Docker builds)
- Nomad (server + client in the same binary)
- Consul
- Vault

Check them:

```bash
nomad version
consul version
vault version
docker version

