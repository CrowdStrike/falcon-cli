# CrowdStrike Falcon CLI

A command-line interface for working with the CrowdStrike Falcon platform.

## Description

The CrowdStrike Falcon CLI provides terminal access to Falcon platform APIs. It supports querying and managing cloud security posture (assets, risks, misconfigurations, compliance, Kubernetes resources, IaC scanning, vulnerabilities, and suppression rules), and authenticating across multiple Falcon clouds and tenants.

Credentials are stored in `~/.falcon/config.yaml` with named profile support. Environment variables prefixed with `FALCON_` are also supported.

## Falcon API Permissions

The CLI requires a CrowdStrike API client with appropriate scopes for the commands you intend to use. Each command group requires its own set of scopes.

> [!NOTE]
> Run `falcon fcs doctor` to check which scopes your API client currently has and which commands are available.

## Usage

See the full [command reference](docs/usage.md) for all commands, subcommands, and flags.

```bash
# Configure credentials interactively
falcon auth config

# List cloud assets
falcon fcs assets list

# View cloud security risks
falcon fcs risks list

# List Kubernetes clusters
falcon fcs kubernetes clusters
```

## Installation

Download pre-built binaries for your platform from the [Releases](https://github.com/CrowdStrike/falcon-cli/releases) page. Builds are available for Linux, macOS, and Windows on amd64 and arm64.

### Building from source

Requires Go 1.26+:

```bash
git clone https://github.com/CrowdStrike/falcon-cli.git
cd falcon-cli
make build
```

This produces a `falcon` binary in the project root.

## Contributing

We welcome contributions that improve the CLI. Please ensure that your contributions pass all CI checks (`make test && make lint`).

See [CONTRIBUTING.md](CONTRIBUTING.md) for details.

## Support

CrowdStrike Falcon CLI is a community-driven, open source project designed to streamline interactions with the CrowdStrike Falcon platform. While not a formal CrowdStrike product, CrowdStrike Falcon CLI is maintained by CrowdStrike and supported in partnership with the open source developer community.

For additional information, please refer to the [SUPPORT.md](SUPPORT.md) file.

## License

See [LICENSE](LICENSE)
