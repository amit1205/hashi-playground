# Consul Template (The Dynamic Configuration Tool)

Consul Template is a standalone application that is a template renderer, notifier, and supervisor for data stored in Consul (and HashiCorp Vault). It allows you to automatically generate configuration files from the live data in Consul.

|Feature	    |	Description    |
|-------------------|------------------|
|Dynamic Configuration| It continuously watches for changes to the service catalog or key/value store in Consul.|
|Template Rendering| When data changes, it uses Go-based templates to re-render configuration files (like Nginx reverse proxy configs, application settings, etc.) on the local file system.|
|Process Supervision| After rendering a new file, it can optionally execute a command, such as sending a signal (like SIGHUP) to a service like Nginx or HAProxy to gracefully reload the new configuration without dropping traffic.|

Consul Template is used to make static configuration files behave dynamically by keeping them synchronized with the changing, real-time data from the Consul cluster.