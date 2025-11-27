# Allow reading data for any secret under secret/ in dev
path "secret/data/*" {
  capabilities = ["read"]
}

# Allow reading metadata (KV v2 needs this sometimes)
path "secret/metadata/*" {
  capabilities = ["read"]
}
