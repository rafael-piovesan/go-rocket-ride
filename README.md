# go-rocket-ride

This is a toy project based on [rocket-rides-atomic](https://github.com/brandur/rocket-rides-atomic) repo and, of course, on the original [Stripe's Rocket Rides](https://github.com/stripe/stripe-connect-rocketrides) demo as well. It aims to replicate the implementation of idempotency keys using Golang and Clean Architecture. Please refer to Brandur's amazing [article](https://brandur.org/idempotency-keys) about this topic for full details.

Quoting [Brandur's own words](https://github.com/brandur/rocket-rides-atomic#rocket-rides-atomic-) about the project:

>The work done API service is separated into atomic phases, and as the name suggests, all the work done during the phase is guaranteed to be atomic. Midway through each API request a call is made out to Stripe's API which can't be rolled back if it fails, so if it does we rely on clients re-issuing the API request with the same Idempotency-Key header until its results are definitive. After any request is considered to be complete, the results are stored on the idempotency key relation and returned for any future requests that use the same key.

## Architecture & code organization

```sh
.
├── adapters          # external data sources (e.g., 3rd party APIs, databases, etc.)
│   └── datastore     # Postgres repository
├── api               # application ports
│   └── http          # HTTP transport layer
├── cmd               # application commands
│   └── api           # 'main.go' for running the API server
├── db                # database related files
│   ├── fixtures      # fixtures used in integration tests and local development
│   └── migrations    # db migrations
├── entity            # application entities (including their specific enum types)
├── mocks             # interface mocks for unit testing
├── pkg               # 3rd party lib wrappers
│   ├── migrate       # help with db migrations during integration tests
│   ├── testcontainer # create db containers used in integration tests
│   ├── testfixtures  # load db fixtures needed for integration tests
│   └── tools         # keeps track of dev deps
├── usecase           # application use cases
├── config.go         # handles config via env vars and .env files
└── rocketride.go     # interface definitions
```

## Setup

Requirements:
1. Make a copy of the `app.env.sample` file and name it `app.env`, then use it to set the env vars as needed
1. A working instance of Postgres (for convenience, there's a `docker-compose.yaml` included to help with this step)
1. Docker is also needed for running the integration tests, since they rely on [testcontainers](https://github.com/testcontainers/testcontainers-go)
1. This project makes use of `testfixtures` CLI to facilitate loading db fixtures, please take a look at how to install it [here](https://github.com/go-testfixtures/testfixtures#cli)
1. Finally, run the following commands:

```sh
# download and install both the project and dev dependencies
make deps

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
-H 'http-authorization: local.user@email.com' \
-d '{ "origin_lat": 0.0, "origin_lon": 0.0, "target_lat": 0.0, "target_lon": 0.0 }'
```

## Development & testing

Use the provided `Makefile` to help you with dev & testing tasks:

```sh
help         # Show help
deps         # Install dependencies
migrate      # Run db migrations (expects $DSN env var, so, run it with: 'DSN=<postgres dsn> make fixtures')
fixtures     # Load db fixtures (expects $DSN env var, so, run it with: 'DSN=<postgres dsn> make fixtures')
lint         # Run linter
format       # Format source code
mock         # Generate interfaces mocks
integration  # Run integration tests
unit         # Run unit tests
server       # Run API server locally
```