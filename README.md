# GoldApps 
A system for syncing LDAP with gsuite and json files written in Go

### Features
* Migrate groups (from -> to):
  * ldap -> gapps
  * ldap -> json
  * gapps -> json
  * gapps -> gapps
  * json -> gapps
  * (json -> json)

## Setup
* Copy example.config.toml to config.toml and edit
* Grabb gapps.json and place in working directory
    * go to [Google developer console](https://console.developers.google.com)
    * go to credentials
    * create new service account för this app
    * use the downloaded file

## Usage
[docker image](https://hub.docker.com/r/cthit/goldapps/)

### Command `goldapps`

The following flags are available:
* `-y`: No interacting from the user required.
* `-i`: Ask the user about everything...
* `-dry`: Makes sure the program does not change anything.
* `-from someString`: Set the group source to `ldap`, `gapps` or `*.json`. In case of `gapps` config value `gapps.provider` will be used.
* `-to someString`: Set the group consumer to 'gapps' or '*.json'. In case of `gapps` config value `gapps.consumer` will be used.
* `-users`: Only collect and sync users
* `-groups`: Only collect and sync groups


Notice that flags should be combined on the form `goldapps -a -b` and **NOT** on the form `goldapps -ab`.