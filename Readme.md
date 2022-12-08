# Cake Store Restful

## Tool dependencies 
- [Mockgen](https://github.com/golang/mock) to generate fake implementation of the interface
- [Docker](https://docs.docker.com/) and [Docker Compose](https://docs.docker.com/compose/) for dcokerize the app
- [Migrate](https://github.com/golang-migrate/migrate) for migration tool

## How to run
all command is mostly wrapped on `Makefile` for simplicity

- run database
    `make docker-compose-up-local`

- run the migration file
    `migration-up`

- run the app
    `go run main.go`

## other make command
- `migration-up`: up migration
- `migration-down`: delete last migration
- `coverage-test`: run coverage test
- `opan-api`: run swagger
- `docker-compose-up-local`: start container
- `docker-compose-down-local`: stop container

## Etc
- import request collection on path `/api/request-collection.json`