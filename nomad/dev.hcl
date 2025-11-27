data_dir   = "/tmp/nomad-dev"
datacenter = "dc1"
name       = "dev-nomad"
bind_addr  = "0.0.0.0"

advertise {
  http = "127.0.0.1"
  rpc  = "127.0.0.1"
  serf = "127.0.0.1"
}

server {
  enabled          = true
  bootstrap_expect = 1
}

client {
  enabled = true
}

consul {
  address          = "127.0.0.1:8500"
  client_auto_join = false
  server_auto_join = false
}

vault {
  enabled = true
  address = "http://127.0.0.1:8200"
  token   = "root"
}

