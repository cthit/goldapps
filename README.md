# GoldApps

A system for syncing LDAP with gsuite and json files written in Go

### Features

Producers
- LDAP
- JSON
- Gamma (1.0)
- Auth (Gamma 2.0)

Consumers
- GApps
- JSON

## Dummy setup

Create the following files:

- `config.toml` - Copy from example.config.toml
- `gapps.json` - Containing `{}`
- `additions.json` - Containing `{}`
- `gamma.json` - Containing `{}`

## Setup

- Copy example.config.toml to config.toml and edit
- Grabb gapps.json and place in working directory
  - go to [Google developer console](https://console.developers.google.com)
  - go to credentials
  - create new service account för this app
  - use the downloaded file

## Usage

Read setup first

### Docker image

`$WAIT` specifies for how long the application should wait before running. This can bes jused in conjunction with `restart: always` to make the bridge run at regular intervals. If you don't desire any waiting effect you can simply set the entrypoint to `./goldapps`.

For some reason `entrypoint` has to be specified in the compose file or the docker run command.

The command should be your flags for the `goldapps` command

### Command `goldapps`

The following flags are available:

- `-y`: No interacting from the user required.
- `-i`: Ask the user about everything...
- `-dry`: Makes sure the program does not change anything.
- `-from someString`: Set the group source to `ldap`, `gapps` or `*.json`. In case of `gapps` config value `gapps.provider` will be used.
- `-to someString`: Set the group consumer to 'gapps' or '\*.json'. In case of `gapps` config value `gapps.consumer` will be used.
- `-users`: Only collect and sync users
- `-groups`: Only collect and sync groups
- `-additions *.json`: file with additions

Notice that flags should be combined on the form `goldapps -a -b` and **NOT** on the form `goldapps -ab`.
