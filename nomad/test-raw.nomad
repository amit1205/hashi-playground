job "test-raw" {
  datacenters = ["dc1"]
  type        = "service"

  group "g1" {
    task "t1" {
      driver = "raw_exec"

      config {
        command = "sh"
        args    = ["-c", "sleep 3600"]
      }

      resources {
        cpu    = 50
        memory = 64
      }
    }
  }
}

