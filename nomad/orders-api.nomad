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

    vault {
      policies = ["app"]
    }

    task "orders-api" {
      driver = "docker"

      config {
        image = "your-registry/orders-api:latest"
        ports = ["http"]
      }

      env {
        PORT = "8080"
      }

      template {
        destination = "secrets/env"
        env         = true

        data = <<EOH
APP_SECRET="{{ with secret "secret/data/orders-api" }}{{ .Data.data.APP_SECRET }}{{ end }}"
EOH
      }

      resources {
        cpu    = 200
        memory = 128
      }
    }
  }
}

