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

Download a package for your platform from the [Releases](https://github.com/CrowdStrike/falcon-cli/releases) page:

| Platform | Package |
|---|---|
| Debian, Ubuntu | `falcon-cli_<version>-1_<arch>.deb` (`amd64`, `arm64`, `s390x`, `ppc64el`) |
| RHEL, Fedora | `falcon-cli-<version>-1.<arch>.rpm` (`x86_64`, `aarch64`, `s390x`, `ppc64le`) |
| macOS | `falcon-cli-<version>-macos-universal.pkg` (Intel and Apple silicon) |
| Windows | `falcon-cli-<version>-windows-<arch>.msi` (`x86_64`, `arm64`) |

Each release also has `.tar.gz` and `.zip` archives and plain binaries. `checksums.txt` lists the SHA-256 of the archives, binaries, and Linux packages; the macOS and Windows installers each have a `.sha256` file.

### Linux

```bash
sudo apt install ./falcon-cli_<version>-1_amd64.deb    # Debian, Ubuntu
sudo dnf install ./falcon-cli-<version>-1.x86_64.rpm   # RHEL, Fedora
```

To remove it, run `sudo apt remove falcon-cli` or `sudo dnf remove falcon-cli`.

### macOS

The pkg installs `falcon` to `/usr/local/falcon-cli` and links it into `/usr/local/bin`.

The pkg is not signed yet, so macOS blocks it if you downloaded it with a browser. To install it anyway:

1. Open the pkg, then click **Done** in the warning.
2. Open **System Settings** > **Privacy & Security**, scroll down to **Security**, and click **Open Anyway** next to the message about the pkg.
3. Confirm, and enter your password when asked.

Files downloaded with `curl` are not quarantined, so you can also install from a terminal without the warning:

```bash
curl -LO https://github.com/CrowdStrike/falcon-cli/releases/download/v<version>/falcon-cli-<version>-macos-universal.pkg
sudo installer -pkg falcon-cli-<version>-macos-universal.pkg -target /
```

macOS has no uninstaller for packages. To remove the CLI:

```bash
sudo rm /usr/local/bin/falcon
sudo rm -rf /usr/local/falcon-cli
sudo pkgutil --forget com.crowdstrike.falcon-cli
```

### Windows

The MSI installs `falcon.exe` to `C:\Program Files\CrowdStrike\Falcon CLI` and adds that folder to the system `PATH`. Open a new terminal afterward so it picks up the new `PATH`.

The MSI is not signed yet, so Microsoft Defender SmartScreen may warn you when you open it. Click **More info**, then **Run anyway**.

To remove it, uninstall **CrowdStrike Falcon CLI** from **Settings** > **Apps** > **Installed apps**.

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
