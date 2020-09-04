# pinger

> This project is currently in active development

Pinger is an open-source implementation of an application that can create
status pages. Some similar applications are [Freshping](https://www.freshworks.com/website-monitoring/)
and [Apex Ping](https://apex.sh/ping/).

## Development

### Prerequisites

Before beginning make sure you have all the required tools installed.

- [Go](https://golang.org/doc/install) >= 1.4 (Preferred >= 1.5)
- [PostgreSQL](https://www.postgresql.org/download/) >= 11 (Preferred >= 12)
- [Timescale](https://docs.timescale.com/latest/getting-started/installation) >= 1.7
- [Docker](https://docs.docker.com/get-docker/)

### Setup

1. Clone the repository and change-directory to it.
   ```sh
   $ git clone git@github.com/sdslabs/pinger
   $ cd pinger
   ```

1. Install required development tools (linter, protoc etc.)
   ```sh
   $ make install
   ```

1. Build the project as `./target/pinger`.
   ```sh
   $ make build
   ```

1. Check for linting errors.
   ```sh
   $ make lint

   # Formats auto-fixable errors
   $ make format
   ```

**NOTE** Before committing changes, make sure the project builds and there are no linting errors.

***

Made by [SDSLabs](https://sdslabs.co)
