# Consul 

Consul is a distributed, highly available service networking solution. It is the core tool that provides the centralized "source of truth" for your services.

*Feature*		*Description*
Service Discovery       Services register themselves with Consul's Service Catalog. Other services can then query 
                        Consul (via DNS or an HTTP API) to find the network location and address of a healthy 
                        instance of another service.

Health Checking         Consul runs integrated health checks against registered services. If a service instance is
                        unhealthy, Consul automatically stops routing traffic to it, improving application resilience
.

Key/Value (KV) Store    A lightweight, distributed store that allows you to store dynamic application configuration,
                        feature flags, or metadata.

Service Mesh            Consul can act as the control plane for a service mesh, managing sidecar proxies to enable 
                        secure, identity-based communication (with automatic TLS encryption) and traffic management 
                        between services.
