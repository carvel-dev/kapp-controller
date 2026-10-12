# Credential Redaction in App and PackageRepository Status

## Overview

kapp-controller automatically redacts credential-bearing URLs from the status fields of `App` and `PackageRepository` resources. This prevents sensitive information from being exposed via `kubectl get ... -o yaml` or any tooling that reads Kubernetes resource status.

## What Gets Redacted

When subprocess tools (`vendir`, `ytt`, `kbld`, `kapp`) output URLs containing credentials in the following status fields, the credentials are automatically masked:

- **Fetch status:** `Fetch.Stdout`, `Fetch.Stderr`, `Fetch.Error`
- **Template status:** `Template.Stdout`, `Template.Stderr`, `Template.Error`
- **Deploy status:** `Deploy.Stdout`, `Deploy.Stderr`, `Deploy.Error`
- **Inspect status:** `Inspect.Stdout`, `Inspect.Stderr`, `Inspect.Error`
- **Conditions:** `Conditions[].Message`
- **Status messages:** `FriendlyDescription`, `UsefulErrorMessage`

## Redaction Pattern

The redaction covers URLs with userinfo (credentials) in the form:
```
scheme://user:password@host/path
ssh://git@github.com/repo.git
https://token@registry.example.com
```

These become:
```
scheme://[REDACTED]@host/path
ssh://[REDACTED]@github.com/repo.git
https://[REDACTED]@registry.example.com
```

### Supported URL Schemes

- `https://`, `http://`
- `ssh://` (Git over SSH)
- `ftp://`
- Custom schemes matching pattern `[a-zA-Z][a-zA-Z0-9+.-]*://`

### Edge Cases

- **Passwords with `@` symbol:** The regex matches greedily up to the last `@` before the host, so passwords like `pass@word` are fully masked: `user:pass@word@host.com` → `[REDACTED]@host.com`
- **Multiple URLs in output:** All credential-bearing URLs in the output are redacted
- **Non-URL text:** Plain text like email addresses (`user@example.com`) without a URL scheme are not affected

## Examples

### GOPROXY with credentials
```
GOPROXY=https://myuser:mytoken@proxy.example.com
↓
GOPROXY=https://[REDACTED]@proxy.example.com
```

### Git clone over HTTPS
```
git clone https://oauth2:token@github.com/myorg/myrepo.git
↓
git clone https://[REDACTED]@github.com/myorg/myrepo.git
```

### Git clone over SSH
```
Cloning ssh://git@github.com/carvel-dev/kapp-controller.git
↓
Cloning ssh://[REDACTED]@github.com/carvel-dev/kapp-controller.git
```

### Docker registry with credentials
```
Authenticating to https://user:password@registry.example.com
↓
Authenticating to https://[REDACTED]@registry.example.com
```

## Security Implications

### What Is Protected

- **Process inspection:** Credentials are not visible via `kubectl describe app` or similar commands
- **Logs and monitoring:** Status fields written to audit logs or monitoring systems will not contain credentials
- **API access:** Users with read-only access to App/PackageRepository resources cannot extract credentials

### What Is NOT Protected

This redaction protects against credential exposure in status fields only. It does **not**:
- Prevent credentials from being used by subprocess tools (which is necessary for their operation)
- Redact arbitrary secret formats (API keys, bare tokens outside URLs)
- Protect the kubeconfig token itself (protected separately via file permissions)
- Redact credentials in logs from the kapp-controller pod itself (outside status fields)

## Developer Notes

The redaction is implemented in `pkg/exec/cmd_run_result.go:RedactSecrets()` using a regex pattern that:
1. Matches URL schemes: `([a-zA-Z][a-zA-Z0-9+.-]*://)`
2. Matches userinfo component: `[^/\s]+@` (greedily up to the last `@` before host)
3. Replaces with: `${1}[REDACTED]@`

Redaction is applied at every point subprocess output is written to status:
- `pkg/app/app_reconcile.go` — App reconciliation paths
- `pkg/pkgrepository/app_reconcile.go` — PackageRepository reconciliation paths

The function is idempotent—calling `RedactSecrets()` multiple times on the same string produces the same result.
