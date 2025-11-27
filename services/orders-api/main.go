package main

import (
    "context"
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "os"
    "time"

    consulapi "github.com/hashicorp/consul/api"
    vaultapi "github.com/hashicorp/vault/api"
)

type ConfigResponse struct {
    Message   string `json:"message"`
    AppSecret string `json:"app_secret"`
}

func main() {
    port := getEnv("PORT", "8080")

    mux := http.NewServeMux()
    mux.HandleFunc("/health", healthHandler)
    mux.HandleFunc("/config", configHandler)

    addr := ":" + port
    log.Printf("orders-api listening on %s\n", addr)
    if err := http.ListenAndServe(addr, mux); err != nil {
        log.Fatalf("server failed: %v", err)
    }
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(http.StatusOK)
    _, _ = w.Write([]byte("ok"))
}

func configHandler(w http.ResponseWriter, r *http.Request) {
    // 1) Vault: read APP_SECRET from KV v2 at secret/orders-api
    appSecret, err := readVaultSecret()
    if err != nil {
        log.Printf("error reading Vault secret: %v", err)
        appSecret = "ERROR_FROM_VAULT"
    }
    if appSecret == "" {
        appSecret = "APP_SECRET not set"
    }

    // 2) Consul KV: read message
    message := "default message"
    if msg, err := readConsulKV("config/orders/message"); err == nil && msg != "" {
        message = msg
    }

    resp := ConfigResponse{
        Message:   message,
        AppSecret: fmt.Sprintf("len=%d", len(appSecret)),
    }

    w.Header().Set("Content-Type", "application/json")
    _ = json.NewEncoder(w).Encode(resp)
}

func readConsulKV(key string) (string, error) {
    cfg := consulapi.DefaultConfig()
    client, err := consulapi.NewClient(cfg)
    if err != nil {
        return "", err
    }

    kv := client.KV()
    pair, _, err := kv.Get(key, nil)
    if err != nil {
        return "", err
    }
    if pair == nil {
        return "", nil
    }
    return string(pair.Value), nil
}

func readVaultSecret() (string, error) {
    addr := getEnv("VAULT_ADDR", "http://127.0.0.1:8200")
    token := os.Getenv("VAULT_TOKEN")
    if token == "" {
        return "", fmt.Errorf("VAULT_TOKEN not set")
    }

    cfg := vaultapi.DefaultConfig()
    cfg.Address = addr

    client, err := vaultapi.NewClient(cfg)
    if err != nil {
        return "", fmt.Errorf("failed to create Vault client: %w", err)
    }
    client.SetToken(token)

    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    // KV v2 at mount "secret" -> logical path: secret/data/orders-api
    kv := client.KVv2("secret")
    secret, err := kv.Get(ctx, "orders-api")
    if err != nil {
        return "", fmt.Errorf("failed to read KV v2 secret: %w", err)
    }
    if secret == nil || secret.Data == nil {
        return "", fmt.Errorf("secret or data nil")
    }

    raw, ok := secret.Data["APP_SECRET"]
    if !ok {
        return "", fmt.Errorf("APP_SECRET key not found in Vault data")
    }

    s, ok := raw.(string)
    if !ok {
        return "", fmt.Errorf("APP_SECRET is not a string")
    }

    return s, nil
}

func getEnv(key, def string) string {
    v := os.Getenv(key)
    if v == "" {
        return def
    }
    return v
}
