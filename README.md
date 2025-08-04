![Baton Logo](./baton-logo.png)

# `baton-miro` [![Go Reference](https://pkg.go.dev/badge/github.com/conductorone/baton-miro.svg)](https://pkg.go.dev/github.com/conductorone/baton-miro) ![main ci](https://github.com/conductorone/baton-miro/actions/workflows/main.yaml/badge.svg)

`baton-miro` is a connector for Baton built using the [Baton SDK](https://github.com/conductorone/baton-sdk).

Check out [Baton](https://github.com/conductorone/baton) to learn more the project in general.

## Connector Capabilities

1. **Resources synced**:

   - Users
   - Teams
   - Roles
   - Licenses

2. **Account provisioning**

   - Create Users

3. **Entitlement provisioning**

   - Assign User To Team
   - Unassign User To Team
   - Grant User To Role
   - Revoke User To Role

## Required permissions

- `identity:read`
- `team:read`
- `team:write` (required for team provisioning)
- `organizations:read`
- `organizations:team:read`
- `organizations:team:write` (required for team provisioning)

**Note:** For user creation (account provisioning), ensure your Miro app has SCIM API access configured.

# Getting Started

## brew

```
brew install conductorone/baton/baton conductorone/baton/baton-miro
baton-miro
baton resources
```

## docker

```
docker run --rm -v $(pwd):/out -e BATON_MIRO_ACCESS_TOKEN=token ghcr.io/conductorone/baton-miro:latest -f "/out/sync.c1z"
docker run --rm -v $(pwd):/out ghcr.io/conductorone/baton:latest -f "/out/sync.c1z" resources
```

## source

```
go install github.com/conductorone/baton/cmd/baton@main
go install github.com/conductorone/baton-miro/cmd/baton-miro@main

baton-miro

baton resources
```

# Data Model

`baton-miro` will pull down information about the following resources:

- Users
- Teams
- Roles
- Licenses

# Contributing, Support and Issues

We started Baton because we were tired of taking screenshots and manually
building spreadsheets. We welcome contributions, and ideas, no matter how
small&mdash;our goal is to make identity and permissions sprawl less painful for
everyone. If you have questions, problems, or ideas: Please open a GitHub Issue!

See [CONTRIBUTING.md](https://github.com/ConductorOne/baton/blob/main/CONTRIBUTING.md) for more details.

# `baton-miro` Command Line Usage

```
baton-miro

Usage:
  baton-miro [flags]
  baton-miro [command]

Available Commands:
  capabilities       Get connector capabilities
  completion         Generate the autocompletion script for the specified shell
  help               Help about any command

Flags:
      --client-id string           The client ID used to authenticate with ConductorOne ($BATON_CLIENT_ID)
      --client-secret string       The client secret used to authenticate with ConductorOne ($BATON_CLIENT_SECRET)
  -f, --file string                The path to the c1z file to sync with ($BATON_FILE) (default "sync.c1z")
  -h, --help                       help for baton-miro
      --log-format string          The output format for logs: json, console ($BATON_LOG_FORMAT) (default "json")
      --log-level string           The log level: debug, info, warn, error ($BATON_LOG_LEVEL) (default "info")
      --miro-access-token       string   Miro Access Token
      --miro-scim-access-token  string   Miro SCIM Access Token
  -p, --provisioning               This must be set in order for provisioning actions to be enabled. ($BATON_PROVISIONING)
  -v, --version                    version for baton-miro

Use "baton-miro [command] --help" for more information about a command.
```
