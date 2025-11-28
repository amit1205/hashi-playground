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
# Diagram
                             ┌──────────────────────────┐
                             │        User (curl)       │
                             │   GET /call-orders       │
                             └─────────────┬────────────┘
                                           │
                                           ▼
                      ┌─────────────────────────────────────────┐
                      │        reporting-worker (Nomad)         │
                      │-----------------------------------------│
                      │ • Discovers orders-api via Consul       │
                      │ • Calls http://host.docker.internal:PORT│
                      │   /config                               │
                      └─────────────┬───────────────────────────┘
                                    │ discover + call
                                    │
           ┌────────────────────────────────────────────────────────────────┐
           │                            Consul                              │
           │----------------------------------------------------------------│
           │  Service Catalog: orders-api, reporting-worker                 │
           │  KV: config/orders/message                                     │
           └───────────┬────────────────────────────────────────────────────┘
                       │ healthy instance: host.docker.internal:20088
                       ▼
              ┌─────────────────────────────────┐
              │       orders-api (Nomad)        │
              │---------------------------------│
              │ • Reads Consul KV (message)     │
              │ • Reads Vault KV (APP_SECRET)   │
              │ • /config returns JSON          │
              └──────────┬──────────────────────┘
                         │
       ┌─────────────────┼────────────────────────────────────┐
       │                 │                                    │
       ▼                 ▼                                    ▼
┌────────────────┐  ┌───────────────────┐           ┌──────────────────────┐
│ Consul KV      │  │ Vault KV v2       │           │   Final JSON returned│
│ config/orders/*│  │ secret/orders-api │           │ to reporting-worker  │
└────────────────┘  └───────────────────┘           └──────────────────────┘


---

# Prerequisites

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
```

# End to End Testing Flow

## 1. Start Consul in Dev Mode
In one terminal:
```
consul agent -dev -node=dev-consul -bind=127.0.0.1
```
This:
- Starts a single-node Consul server
- Listens on `127.0.0.1:8500`
- Enables a dev keyring and in-memory storage
You can open the Consul UI:
- URL: `http://127.0.0.1:8500/ui`

## 2. Start Vault in Dev Mode
In another terminal:
```
vault server -dev -dev-root-token-id=root
```
Notes:
- Vault listens on `http://127.0.0.1:8200`
- Dev mode auto-unseals and uses in-memory storage
- The root token is root (hard-coded with -dev-root-token-id).

Set environment variables in that terminal:

```
export VAULT_ADDR="http://127.0.0.1:8200"
export VAULT_TOKEN="root"
```
Confirm:
```
vault status
```
