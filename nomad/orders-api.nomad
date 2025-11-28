job "orders-api" {
  datacenters = ["dc1"]
  type        = "service"

  group "orders" {
    count = 1

    network {
      port "http" {
        to = 8080
      }
    }

    # Consul service registration + health check
    service {
      name = "orders-api"
      port = "http"

      check {
        name     = "http health"
        type     = "http"
        path     = "/health"
        interval = "5s"
        timeout  = "2s"
      }
    }

    task "orders-api" {
      driver = "docker"

      config {
        image = "orders-api:dev"
        ports = ["http"]
      }

      env {
        PORT        = "8080"
        VAULT_ADDR  = "http://127.0.0.1:8200"
        VAULT_TOKEN = "root" # DEV ONLY - don't do this in prod
	CONSUL_HTTP_ADDR = "http://host.docker.internal:8500"
      }

      resources {
        cpu    = 200
        memory = 128
      }
    }
  }
}
