# Security Policy

## Supported Versions

| Version | Supported          |
| ------- | ------------------ |
| 0.1.x   | :white_check_mark: |

## Reporting a Vulnerability

If you discover a security vulnerability in NextTrace Exporter, please report it by emailing the maintainers or creating a private security advisory on GitHub.

**Please do not open public issues for security vulnerabilities.**

### What to Include

- Description of the vulnerability
- Steps to reproduce
- Potential impact
- Suggested fix (if any)

### Response Time

We aim to respond to security reports within 48 hours and will work to release a patch as soon as possible.

## Security Considerations

### Running with Privileges

NextTrace Exporter requires either:
- Root privileges (not recommended for production)
- `CAP_NET_RAW` capability (recommended)

**Recommended setup:**
```bash
sudo setcap cap_net_raw+ep /usr/local/bin/nexttrace_exporter
sudo setcap cap_net_raw+ep /usr/local/bin/nexttrace
```

### Network Exposure

- Bind to localhost (`127.0.0.1:9101`) if metrics are only needed locally
- Use firewall rules to restrict access to the metrics endpoint
- Consider using TLS reverse proxy (nginx, caddy) for public exposure
- Implement authentication at the reverse proxy level if needed

### Configuration Security

- Protect configuration files with appropriate file permissions:
  ```bash
  sudo chmod 600 /etc/nexttrace_exporter/config.yml
  ```
- Avoid storing sensitive information in configuration files
- Review target lists to prevent unintended network probing

### Docker Security

When running in Docker:
- Use `--cap-add=NET_RAW` instead of `--privileged`
- Run as non-root user when possible (note: may require additional setup)
- Use read-only volumes for configuration:
  ```bash
  -v $(pwd)/config.yml:/etc/nexttrace_exporter/config.yml:ro
  ```
- Keep the Docker image updated

### Input Validation

The exporter validates:
- Configuration file format and values
- Target hostnames/IPs
- Command-line arguments

However, be cautious with:
- User-supplied target lists
- Custom nexttrace arguments that might be exploitable

### Dependencies

- Regularly update Go dependencies
- Monitor security advisories for dependencies
- Use `go mod tidy` to remove unused dependencies

### Logging

- Logs may contain target hostnames/IPs
- Secure log files with appropriate permissions
- Consider log rotation and retention policies
- Avoid logging sensitive information

## Best Practices

1. **Principle of Least Privilege**: Run with minimal required permissions
2. **Network Segmentation**: Deploy in appropriate network zones
3. **Access Control**: Restrict who can access metrics and configuration
4. **Monitoring**: Monitor for unusual execution patterns or errors
5. **Updates**: Keep nexttrace_exporter and nexttrace updated
6. **Configuration Review**: Regularly audit target lists and settings

## Known Limitations

- Requires elevated privileges for raw socket access
- Targets are probed at configured intervals (consider privacy implications)
- No built-in authentication or encryption (use reverse proxy)

## Disclosure Policy

We follow responsible disclosure practices:
1. Report received and acknowledged
2. Vulnerability verified and assessed
3. Fix developed and tested
4. Security advisory published
5. Patch released
6. Public disclosure (typically 90 days after fix)

## Questions?

For security questions or concerns, please open a GitHub discussion or contact the maintainers.
