# Contributing to snap-tpmctl

Thank you for considering a contribution.

This project welcomes bug reports, feature proposals, code improvements, and documentation updates. The goal of these guidelines is to keep collaboration efficient, respectful, and reviewable.

## Code of conduct

By participating in this project, you agree to follow the Ubuntu Code of Conduct:

- https://ubuntu.com/community/ethos/code-of-conduct

## Getting started

### Report issues

Use GitHub Issues to report bugs or request improvements:

- https://github.com/canonical/snap-tpmctl/issues

### Security issues

Do not report security-sensitive issues publicly.

Use private reporting channels:

- https://github.com/canonical/snap-tpmctl/security/advisories

## Development setup

### Requirements

- Go 1.26+
- Linux environment with TPM/FDE tooling available for integration scenarios
- `golangci-lint-v2` for lint checks

### Build

```bash
go build ./cmd/tpmctl
```

### Run

```bash
go run ./cmd/tpmctl -- status
```

### Test

```bash
go test ./...
```

### Lint

```bash
go fmt ./...
golangci-lint-v2 run ./...
```
## Snap packaging

Build snap package locally:

```bash
snapcraft pack
```

Clean snap artifacts:

```bash
snapcraft clean && rm -f *.snap
```

The snap recipe is in `snap/snapcraft.yaml`.

## Contributor License Agreement

For Canonical projects, contributors may be required to sign the Canonical Contributor License Agreement:

- https://ubuntu.com/legal/contributors

## Need help?

For questions and discussion, use GitHub Issues or Ubuntu community channels:

- https://discourse.ubuntu.com/
