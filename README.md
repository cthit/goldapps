# GoldApps 
A system for syncing LDAP with gsuite written in GO

### Features
* TODO

## Setup
* Copy config.toml.example to config.toml and edit
* Grabb gapps.json and place in working directory
    * go to [Google developer console](https://console.developers.google.com)
    * go to credentials
    * create new service account för this app
    * use the downloaded file

## Production
Use the automated [docker image](https://hub.docker.com/r/cthit/goldapps/)

### Example compose file
```yml
version: '2.2'
services:
    ...
```

## Development

### Software requirements
* docker
* docker-compose

### Setup
Run the following command:
1. `docker-compose up`

### Local production environment
You can compile and build the production image with your local codebase using:
`docker-compose -f prod.docker-compose up --build`