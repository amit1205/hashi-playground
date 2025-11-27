job "orders-api-debug" {
  datacenters = ["dc1"]
  type        = "service"

  group "orders" {
    count = 1

    network {
      port "http" {
        to = 8080
      }
    }

    task "orders-api" {
      driver = "docker"

      config {
        image = "orders-api:dev"
        ports = ["http"]
      }

      env {
        PORT = "8080"
      }

      resources {
        cpu    = 100
        memory = 64
      }
    }
  }
}

