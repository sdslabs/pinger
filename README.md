# status

> An open source implmentation of status page.

A work in progress.

## Setup

1. Install the following:
    1. [GoLang](https://golang.org/doc/install)
    1. [PostgreSQL](https://www.postgresql.org/download/)
    1. [TimescaleDB](https://docs.timescale.com/latest/getting-started/installation)

1. Clone the repository and `cd` to it.

1. Run `make build` to build the project. Builds the project in `$GOPATH/bin`.

## Contributing

1. Run `make install` to install development tools.

1. Run `make lint` before pushing code and resolve all the errors.

1. Many errors can be auto-fixed. Just run `make format` to fix them.
