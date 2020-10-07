---
title: Building from source
parent: Development
nav_order: 2
---

# Building from source

Pinger is written purely in Go but is dependent on a few other external
dependencies.

## Prerequisites

Before beginning make sure you have all the required tools installed.

- [Go](https://golang.org/doc/install) >= 1.13 (Preferred >= 1.15)
- [PostgreSQL](https://www.postgresql.org/download/) >= 11 (Preferred >= 12)
- [Timescale](https://docs.timescale.com/latest/getting-started/installation) >= 1.7
- [Docker](https://docs.docker.com/get-docker/) >= 17.05

## Setup

1. Clone the repository and change-directory to it.

   ```sh
   $ git clone git@github.com/sdslabs/pinger.git
   $ cd pinger
   ```

1. Build the project as `./target/pinger`.

   ```sh
   $ make build
   ```

1. Verify a successful build by pinging it.
   ```sh
   $ ./target/pinger ping
   INFO[0000] pong
   ```
