`baton-google-identity-platform` is a connector for Google Identity Platform built using the [Baton SDK](https://github.com/conductorone/baton-sdk). It communicates with the Google API to sync data about users.

Check out [Baton](https://github.com/conductorone/baton) to learn more the project in general.

# Getting Started

## Prerequisites

Service account key for your project. If you don't already have one follow the steps [here](https://cloud.google.com/identity-platform/docs/install-admin-sdk#create-service-account-console) to create a service account and a service account key. If you already have a service account, make sure you have all the permissions set and APIs enabled. Then download the service key which will be used in the connector.

## brew

```
brew install conductorone/baton/baton conductorone/baton/baton-google-identity-platform
baton-google-identity-platform
baton resources
```

## docker

```
docker run --rm -v $(pwd):/out -e BATON_CREDENTIALS_JSON_FILE_PATH=pathOfServiceKey ghcr.io/conductorone/baton-google-identity-platform:latest -f "/out/sync.c1z"
docker run --rm -v $(pwd):/out ghcr.io/conductorone/baton:latest -f "/out/sync.c1z" resources
```

## source

```
go install github.com/conductorone/baton/cmd/baton@main
go install github.com/conductorone/baton-google-identity-platform/cmd/baton-google-identity-platform@main

BATON_CREDENTIALS_JSON_FILE_PATH=pathOfServiceKey
baton resources
```

# Data Model

`baton-google-identity-platform` will pull down information about the following Google Identity Platform resources:

- Users

# Contributing, Support and Issues

We started Baton because we were tired of taking screenshots and manually building spreadsheets. We welcome contributions, and ideas, no matter how small -- our goal is to make identity and permissions sprawl less painful for everyone. If you have questions, problems, or ideas: Please open a Github Issue!

See [CONTRIBUTING.md](https://github.com/ConductorOne/baton/blob/main/CONTRIBUTING.md) for more details.

# `baton-google-identity-platform` Command Line Usage

```
baton-google-identity-platform

Usage:
  baton-google-identity-platform [flags]
  baton-google-identity-platform [command]

Available Commands:
  completion         Generate the autocompletion script for the specified shell
  help               Help about any command

Flags:
  --credentials-json-file-path string       JSON credentials file name for the Google identity platform account  ($BATON_CREDENTIALS_JSON_FILE_PATH)
  -f, --file string                         The path to the c1z file to sync with ($BATON_FILE) (default "sync.c1z")
  -h, --help                                help for baton-google-identity-platform
      --log-format string                   The output format for logs: json, console ($BATON_LOG_FORMAT) (default "json")
      --log-level string                    The log level: debug, info, warn, error ($BATON_LOG_LEVEL) (default "info")
  -v, --version                             version for baton-google-identity-platform

Use "baton-google-identity-platform [command] --help" for more information about a command.

```
