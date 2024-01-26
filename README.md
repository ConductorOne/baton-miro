![Baton Logo](./docs/images/baton-logo.png)

# `baton-miro` [![Go Reference](https://pkg.go.dev/badge/github.com/conductorone/baton-miro.svg)](https://pkg.go.dev/github.com/conductorone/baton-miro) ![main ci](https://github.com/conductorone/baton-miro/actions/workflows/main.yaml/badge.svg)

`baton-miro` is a connector for Baton built using the [Baton SDK](https://github.com/conductorone/baton-sdk). It works with Miro API.

Check out [Baton](https://github.com/conductorone/baton) to learn more about the project in general.

# Prerequisites

Connector requires bearer access token that is used throughout the communication with API. To obtain this token, you have to create one in Miro. More in information about how to generate token [here](https://developers.miro.com/docs/try-out-the-rest-api-in-less-than-3-minutes)). 

After you have obtained access token, you can use it with connector. You can do this by setting `BATON_MIRO_ACCESS_TOKEN` or by passing `--miro-access-token`.

## Required permissions

- identity:read
- team:read
- team:write (could be just read if provisioning is not used)
- organizations:read
- organizations:team:read
- organizations:team:write (could be just read if provisioning is not used)

# Getting Started

## brew

```
brew install conductorone/baton/baton conductorone/baton/baton-miro

BATON_MIRO_ACCESS_TOKEN=token baton-miro
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

BATON_MIRO_ACCESS_TOKEN=token baton-miro
baton resources
```

# Data Model

`baton-miro` will fetch information about the following Baton resources:

- Users
- Teams
- Licenses
- Roles

# Contributing, Support and Issues

We started Baton because we were tired of taking screenshots and manually building spreadsheets. We welcome contributions, and ideas, no matter how small -- our goal is to make identity and permissions sprawl less painful for everyone. If you have questions, problems, or ideas: Please open a Github Issue!

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
      --miro-access-token string   Miro Access Token
  -p, --provisioning               This must be set in order for provisioning actions to be enabled. ($BATON_PROVISIONING)
  -v, --version                    version for baton-miro

Use "baton-miro [command] --help" for more information about a command.
```