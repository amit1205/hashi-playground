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
        http.Error(w, fmt.Sprintf("failed to discover orders-api: %v", err), http.StatusBadGateway)
        return
    }

    url := fmt.Sprintf("http://%s/config", addr)
    log.Printf("calling orders-api at %s", url)

    resp, err := http.Get(url)
    if err != nil {
        log.Printf("error calling orders-api: %v", err)
        http.Error(w, fmt.Sprintf("failed to call orders-api: %v", err), http.StatusBadGateway)
        return
    }
    defer resp.Body.Close()

    body, _ := io.ReadAll(resp.Body)

    w.Header().Set("Content-Type", "application/json")
    _, _ = w.Write([]byte(fmt.Sprintf(
        `{"orders_api_url": %q, "status": %d, "response": %s}`,
        url,
        resp.StatusCode,
        body,
    )))
}


//Why this helps:
//Nomad/Consul often register the service with Service.Address = 127.0.0.1 (especially with your advertise http = 127.0.0.1).
//Inside Docker, 127.0.0.1 is the container, not the host.
//svc.Node.Address will be something like 10.0.0.107, which is reachable from the container.
//So now reporting-worker will call http://10.0.0.107:<mapped-port>/config instead of http://127.0.0.1:<mapped-port>/config


func discoverService(name string) (string, error) {
    cfg := consulapi.DefaultConfig()
    client, err := consulapi.NewClient(cfg)
    if err != nil {
        return "", err
    }

    // Only passing instances
    services, _, err := client.Health().Service(name, "", true, nil)
    if err != nil {
        return "", err
    }
    if len(services) == 0 {
        return "", fmt.Errorf("no healthy instances found for %s", name)
    }

    svc := services[0]

    // Prefer the service address, but if it's empty or loopback, fall back to the node address.
    //host := svc.Service.Address
    //if host == "" || host == "127.0.0.1" {
    //    host = svc.Node.Address
    //}
    host := "host.docker.internal"
    port := svc.Service.Port


    // Log what we’re going to use (optional but useful)
    log.Printf("discoverService(%q): using %s:%d (service=%q node=%q)",
        name, host, port, svc.Service.Address, svc.Node.Address)

    return fmt.Sprintf("%s:%d", host, port), nil
}


func getEnv(key, def string) string {
    v := os.Getenv(key)
    if v == "" {
        return def
    }
    return v
}
