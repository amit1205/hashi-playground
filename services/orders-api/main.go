package main

import (
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "os"

    consulapi "github.com/hashicorp/consul/api"
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
    appSecret := os.Getenv("APP_SECRET") // injected via Vault/Template
    if appSecret == "" {
        appSecret = "APP_SECRET not set"
    }

    // Read message from Consul KV
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
    // In-cluster, CONSUL_HTTP_ADDR can be set by Nomad/Env if not default.
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

func getEnv(key, def string) string {
    v := os.Getenv(key)
    if v == "" {
        return def
    }
    return v
}
