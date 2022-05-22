# go-rocket-ride
![Lint and Tests](https://github.com/rafael-piovesan/go-rocket-ride/actions/workflows/lint-tests.yml/badge.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/rafael-piovesan/go-rocket-ride)](https://goreportcard.com/report/github.com/rafael-piovesan/go-rocket-ride)

## Related Articles
Read more about this project's motivations and reasonings:
* [Go (Golang): Clean Architecture & Repositories vs Transactions](https://medium.com/@rubens.piovesan/go-golang-clean-architecture-repositories-vs-transactions-9b3b7c953463)
* [Go (Golang): Testing tools & tips to step up your game](https://medium.com/@rubens.piovesan/go-golang-testing-tools-tips-to-step-up-your-game-4ed165a5b3b5)

## Description
This is a toy project based on [rocket-rides-atomic](https://github.com/brandur/rocket-rides-atomic) repo and, of course, on the original [Stripe's Rocket Rides](https://github.com/stripe/stripe-connect-rocketrides) demo as well. It aims to replicate the implementation of idempotency keys using Golang and Clean Architecture. Please refer to Brandur's amazing [article](https://brandur.org/idempotency-keys) about this topic for full details.

Quoting [Brandur's own words](https://github.com/brandur/rocket-rides-atomic#rocket-rides-atomic-) about the project:

>The work done API service is separated into atomic phases, and as the name suggests, all the work done during the phase is guaranteed to be atomic. Midway through each API request a call is made out to Stripe's API which can't be rolled back if it fails, so if it does we rely on clients re-issuing the API request with the same Idempotency-Key header until its results are definitive. After any request is considered to be complete, the results are stored on the idempotency key relation and returned for any future requests that use the same key.


## Architecture & code organization

```sh
.
├── api               # application ports
│   └── http          # HTTP transport layer
├── cmd               # application commands
│   └── api           # 'main.go' for running the API server
├── datastore         # app data stores (e.g., PostgreSQL, MySQL, etc.)
│   ├── bun           # Postgres data access based on Bun ORM 
│   └── sqlc          # Postgres data access based on Sqlc
├── db                # database related files
│   ├── fixtures      # fixtures used in integration tests and local development
│   ├── migrations    # db migrations
│   └── queries       # Sqlc db queries
├── entity            # application entities (including their specific enum types)
├── mocks             # interface mocks for unit testing
├── pkg               # 3rd party lib wrappers
│   ├── migrate       # help with db migrations during integration tests
│   ├── stripemock    # set Stripe's API SDK Backend to use stripe-mock
│   ├── testcontainer # create db containers used in integration tests
│   ├── testfixtures  # load db fixtures needed for integration tests
│   └── tools         # keep track of dev deps
├── usecase           # application use cases
├── config.go         # handle config via env vars and .env files
└── rocketride.go     # interface definitions
```

## Setup

Requirements:
1. Make a copy of the `app.env.sample` file and name it `app.env`, then use it to set the env vars as needed
1. A working instance of Postgres (for convenience, there's a `docker-compose.yaml` included to help with this step)
1. Stripe's [stripe-mock](https://github.com/stripe/stripe-mock) (also provided with the `docker-compose.yaml`)
1. Docker is also needed for running the integration tests, since they rely on [testcontainers](https://github.com/testcontainers/testcontainers-go)
1. This project makes use of `testfixtures` CLI to facilitate loading db fixtures, please take a look at how to install it [here](https://github.com/go-testfixtures/testfixtures#cli)
1. Instead of a `Makefile`, this project uses `Taskfile`, please check its installation procedure [here](https://taskfile.dev/#/installation)
1. Finally, run the following commands:

```sh
# download and install both the project and dev dependencies
task deps

# start the dependencies (postgres and stripe-mock)
docker-compose up -d

# run db migrations, remember to export the $DSN env var before running it
DSN=postgresql://postgres:postgres@localhost:5432/rides?sslmode=disable make migrate

# load db fixtures, remember to export the $DSN env var before running it
DSN=postgresql://postgres:postgres@localhost:5432/rides?sslmode=disable make fixtures

# start the API server
make server
```
Once the server is up running, send requests to it:
```sh
curl -i -w '\n' -X POST http://localhost:8080/ \
-H 'content-type: application/json' \
-H 'idempotency-key: key123' \
-H 'authorization: local.user@email.com' \
-d '{ "origin_lat": 0.0, "origin_lon": 0.0, "target_lat": 0.0, "target_lon": 0.0 }'
```

## Development & testing

Use the provided `Makefile` to help you with dev & testing tasks:

```log
task: Available tasks for this project:
* api:                  Run API server locally
* db:fixtures:          Load DB fixtures (expects $DSN env var to be set)
* db:migrate-drop:      Drop local DB (expects $DSN env var to be set)
* db:migrate-up:        Up DB migrations (expects $DSN env var to be set)
* db:sqlc:              Generate sqlc files
* deps:                 Install dependencies
* format:               Format source code
* lint:                 Run linter
* test:integration:     Run integration tests
* test:mock:            Generate interfaces mocks
* test:unit:            Run unit tests
```

## IDE
If you're using VS Code, make sure to check out this article with some tips on how to setup your IDE: [Setting up VS Code for Golang](https://medium.com/@rubens.piovesan/setting-up-vs-code-for-golang-2021-4cb6ebdd557c).

And here's a sample `settings.json` file with some suggestions:

```json
{
  "go.testFlags": [
    "-failfast",
    "-v"
  ],
  "go.toolsManagement.autoUpdate": true,
  "go.useLanguageServer": true,
  "go.lintFlags": [
    "--build-tags=integration,unit"
  ],
  "go.lintOnSave": "package",
  "go.lintTool": "golangci-lint",
  "gopls": {
    "build.buildFlags": [
      "-tags=integration,unit"
    ],
    "ui.semanticTokens": true
  }
}
```