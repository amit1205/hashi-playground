package main

import (
    "fmt"
    "io"
    "log"
    "net/http"
    "os"

    consulapi "github.com/hashicorp/consul/api"
)

type OrdersConfig struct {
    Message   string `json:"message"`
    AppSecret string `json:"app_secret"`
}

func main() {
    port := getEnv("PORT", "8081")

    mux := http.NewServeMux()
    mux.HandleFunc("/health", healthHandler)
    mux.HandleFunc("/call-orders", callOrdersHandler)

    addr := ":" + port
    log.Printf("reporting-worker listening on %s\n", addr)
    if err := http.ListenAndServe(addr, mux); err != nil {
        log.Fatalf("server failed: %v", err)
    }
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(http.StatusOK)
    _, _ = w.Write([]byte("ok"))
}

func callOrdersHandler(w http.ResponseWriter, r *http.Request) {
    // Discover orders-api via Consul
    addr, err := discoverService("orders-api")
    if err != nil {
        log.Printf("error discovering orders-api: %v", err)
        http.Error(w, "failed to discover orders-api", http.StatusBadGateway)
        return
    }

    url := fmt.Sprintf("http://%s/config", addr)
    resp, err := http.Get(url)
    if err != nil {
        log.Printf("error calling orders-api: %v", err)
        http.Error(w, "failed to call orders-api", http.StatusBadGateway)
        return
    }
    defer resp.Body.Close()

    body, _ := io.ReadAll(resp.Body)

    w.Header().Set("Content-Type", "application/json")
    _, _ = w.Write([]byte(fmt.Sprintf(`{"orders_api_url": %q, "response": %s}`, url, body)))
}

func discoverService(name string) (string, error) {
    cfg := consulapi.DefaultConfig()
    client, err := consulapi.NewClient(cfg)
    if err != nil {
        return "", err
    }

    // Simple: pick any passing instance
    services, _, err := client.Health().Service(name, "", true, nil)
    if err != nil {
        return "", err
    }
    if len(services) == 0 {
        return "", fmt.Errorf("no healthy instances found for %s", name)
    }

    svc := services[0]
    return fmt.Sprintf("%s:%d", svc.Service.Address, svc.Service.Port), nil
}

func getEnv(key, def string) string {
    v := os.Getenv(key)
    if v == "" {
        return def
    }
    return v
}
