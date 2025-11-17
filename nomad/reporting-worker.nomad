job "reporting-worker" {
  datacenters = ["dc1"]
  type        = "service"

  group "reporting" {
    count = 1

    network {
      port "http" {
        to = 8081
      }
    }

    service {
      name = "reporting-worker"
      port = "http"

      check {
        name     = "http health"
        type     = "http"
        path     = "/health"
        interval = "5s"
        timeout  = "2s"
      }
    }

    # Reporting worker doesn't need Vault for now, but you could
    # add a vault {} block similar to orders-api if needed.

    task "reporting-worker" {
      driver = "docker"

      config {
        image = "your-registry/reporting-worker:latest"
        ports = ["http"]
      }

      env {
        PORT = "8081"
        # Optionally, override CONSUL_HTTP_ADDR if needed.
        # CONSUL_HTTP_ADDR = "http://consul.service.consul:8500"
      }

      resources {
        cpu    = 200
        memory = 128
      }
    }
  }
}
