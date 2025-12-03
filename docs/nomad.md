# Nomad (Workload Orchestration)

Nomad is a flexible, lightweight workload orchestrator and scheduler. Its main job is to ensure that applications, or "jobs," are running and healthy across a cluster of machines.

- Workload Orchestration: Nomad accepts job specifications (defined declaratively in HCL) and decides where and how to run them across the cluster (the Nomad Servers manage, and Nomad Clients run the jobs)

- Flexibility (Diverse Workloads): Unlike some orchestrators primarily focused on containers (like Kubernetes), Nomad can natively run a wide variety of workload types using different task drivers, including:
  - Containers (Docker, Podman)
  - Non-containerized applications (plain binaries, Java JARs, executables)
  - Virtual Machines (using QEMU)
  - Batch Jobs and System Services

- Simplicity and Scalability: Nomad runs as a single, self-contained binary on each node and is known for its architectural simplicity and ability to scale to thousands of nodes efficiently with minimal operational overhead.

- Integration with Vault: Nomad has deep, native integration with Vault. A Nomad job definition can reference secrets stored in Vault, and Nomad will securely retrieve the necessary dynamic credentials and inject them into the running container or application environment at runtime.

In summary, Vault is the security heart that provides secrets, and Nomad is the deployment brain that schedules the applications that need those secrets