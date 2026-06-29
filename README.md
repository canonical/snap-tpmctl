# snap-tpmctl

[actions-image]: https://github.com/canonical/snap-tpmctl/actions/workflows/qa.yaml/badge.svg
[actions-url]: https://github.com/canonical/snap-tpmctl/actions?query=workflow%3AQA

[license-image]: https://img.shields.io/badge/License-GPL3.0-blue.svg

[codecov-image]: https://codecov.io/gh/canonical/snap-tpmctl/graph/badge.svg
[codecov-url]: https://codecov.io/gh/canonical/snap-tpmctl

[snap-badge-image]: https://snapcraft.io/snap-tpmctl/badge.svg
[snap-url]: https://snapcraft.io/snap-tpmctl

[![snap-tpmctl][snap-badge-image]][snap-url]
[![Code quality][actions-image]][actions-url]
[![License][license-image]](COPYING)
[![Code coverage][codecov-image]][codecov-url]


`snap-tpmctl` is a command-line tool to manage TPM-backed Full Disk Encryption (FDE) on Ubuntu systems.

It provides a single operational interface for:

- checking TPM/FDE status
- creating and validating recovery keys
- managing passphrase and PIN authentication methods
- listing LUKS keyslots metadata
- unlocking and mounting encrypted volumes

## Why this project is useful

Managing TPM and FDE operations often involves multiple low-level tools and repetitive, error-prone steps.

`snap-tpmctl` centralizes these workflows into clear, auditable commands designed for administrators, support engineers, and automation scenarios.

## Quick usage examples

Check current FDE status:

```bash
snap-tpmctl status
```

Create a recovery key:

```bash
sudo snap-tpmctl create-recovery-key my-recovery-key
```

Add PIN authentication:

```bash
sudo snap-tpmctl add-pin
```

List configured recovery keys:

```bash
snap-tpmctl list-recovery-keys
```

Unlock and mount an encrypted volume:

```bash
sudo snap-tpmctl mount-volume /dev/nvme0n1p4 /media/my-vol
```

## Contributing

Contributions are welcome. Please read [`CONTRIBUTING.md`](./CONTRIBUTING.md) for more info.