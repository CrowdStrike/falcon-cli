# Command Reference

## Synopsis

```
falcon <command> <subcommand> [flags]
```

## Commands

| Command | Description |
|---------|-------------|
| `falcon auth` | Authenticate with the CrowdStrike Falcon API |
| `falcon auth config` | Configure API credentials |
| `falcon fcs` | Manage Falcon Cloud Security resources |
| `falcon version` | Display the CLI version |

## Falcon Cloud Security (`fcs`) subcommands

| Subcommand | Description |
|------------|-------------|
| `fcs assets list` | List cloud assets |
| `fcs risks list` | List cloud security risks |
| `fcs iom list` | List Indicators of Misconfiguration |
| `fcs compliance frameworks` | List compliance frameworks |
| `fcs compliance rules` | List compliance rules |
| `fcs groups list\|create\|update\|delete` | Manage cloud groups |
| `fcs suppression list\|create\|delete` | Manage suppression rules |
| `fcs kubernetes clusters\|containers\|images` | View Kubernetes resources |
| `fcs vulnerabilities list` | List vulnerabilities |
| `fcs iac list` | List IaC scan results |
| `fcs doctor` | Check API client permissions |

## Global Flags

```
      --config string          Path to config file (default "~/.falcon/config.yaml")
      --profile string         Named credential profile to use
      --cloud string           Falcon cloud region (auto-discovered by default)
      --cid string             CrowdStrike Customer ID
      --client-id string       OAuth2 Client ID (overrides config)
      --client-secret string   OAuth2 Client Secret (overrides config)
      --member-cid string      Member CID for MSSP/flight-control
      --base-url string        Override the API base URL
      --verbose                Enable verbose output
```

## Authentication

The CLI uses CrowdStrike API credentials (OAuth2 client ID and secret). Configure them interactively:

```bash
falcon auth config
```

Or pass them directly:

```bash
falcon auth config --client-id <YOUR_CLIENT_ID> --client-secret <YOUR_CLIENT_SECRET>
```

Credentials are stored in `~/.falcon/config.yaml`. You can also use environment variables:

```bash
export FALCON_CLIENT_ID=<YOUR_CLIENT_ID>
export FALCON_CLIENT_SECRET=<YOUR_CLIENT_SECRET>
```

## Examples

```bash
# Configure credentials interactively
falcon auth config

# Configure credentials directly
falcon auth config --client-id <ID> --client-secret <SECRET>

# List cloud assets
falcon fcs assets list

# View cloud security risks
falcon fcs risks list

# List Kubernetes clusters
falcon fcs kubernetes clusters

# Check compliance frameworks
falcon fcs compliance frameworks

# Check API client permissions
falcon fcs doctor
```
