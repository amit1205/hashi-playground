# Vault (Secrets Management)

Vault is a dedicated tool for securely storing, managing, and controlling access to secrets and sensitive data. In the context of applications, a "secret" is anything that grants access to sensitive systems.

- Centralized Secrets Store: It provides a central, highly available, and encrypted location for sensitive credentials, such as API keys, database passwords, SSH keys, encryption keys, and TLS certificates

- Identity-Based Access: Vault enforces strict access controls. Applications, users, or machines must first authenticate their identity using various methods (e.g., cloud provider roles, Kubernetes service accounts) before Vault grants them a token to access specific secrets based on defined policies.

- Dynamic Secrets: Vault can generate on-demand, short-lived credentials (e.g., temporary database usernames/passwords) for applications. This is a major security advantage because the secret exists only for the duration it's needed, reducing the risk of static credentials being leaked.

- Encryption as a Service (Transit Engine): It can perform cryptographic functions (encryption/decryption) on data without storing the encryption key itself, allowing applications to securely handle data without implementing complex cryptography.

- Audit Logging: Vault maintains a detailed audit log of every request and response, providing a clear record of who accessed which secret and when.