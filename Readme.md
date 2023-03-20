Provider AMT
---

Provider AMT adds support for configuring Intel Active Management Tools during Kairos installation.

```yaml
#cloud-config
amt:
  dns_suffix_override: "" # Hostname override
  hostname: "" # Hostname override
  lms_address: "" # LMS address (default "localhost"). Can be used to change location of LMS for debugging.
  lms_port: "" # LMS port (default "16992")
  proxy_address: "" # Proxy address and port
  password: "" # AMT password
  profile: ""# Name of the profile to use
  server_address: "" # WebSocket address of server to activate against
  timeout: "" # timeout for activation
  extra: # extra arguments to pass to the activate command
    -foo: bar
```