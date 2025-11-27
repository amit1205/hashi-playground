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

    # Consul service registration + health check
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

    task "reporting-worker" {
      driver = "docker"

      config {
        image = "reporting-worker:dev"
        ports = ["http"]
      }

      env {
        PORT = "8081"
        # If needed, you can be explicit:
        # CONSUL_HTTP_ADDR = "http://127.0.0.1:8500"
      }

      resources {
        cpu    = 200
        memory = 128
      }
    }
  }
}
