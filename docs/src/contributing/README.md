# Contributing

> **Note:** This guide is for developers who want to make changes to the
> Pinger codebase. If you want to create your Pinger plugin, there's
> a completely different section dedicated for it.

## Prerequisites

These include some tools like the Go compiler and time series database for
storing metrics.

### Build dependencies

Before beginning make sure you have all the required tools installed:

- [Go](https://golang.org/doc/install) `v1.15`
- [Protobuf for Golang](https://developers.google.com/protocol-buffers/docs/gotutorial#compiling-your-protocol-buffers)
  - `protoc v3.14`
  - `protoc-gen-go v1.25`
- [Golang CI Lint](https://golangci-lint.run/usage/install/) `v1.32.2`
- [mdBook](https://rust-lang.github.io/mdBook/cli/index.html) `v0.4.4`

### Runtime Dependencies

Externally, Pinger only relies on databases.

- [PostgreSQL](https://www.postgresql.org/download/) `v12`
- [Timescale](https://docs.timescale.com/latest/getting-started/installation) `v2.0`

## Make Pinger

Before we can do anything, ewe need to fetch the source code present in
[this Github repository](https://github.com/sdslabs/pinger).

```sh
# If you use SSH
$ git clone git@github.com:sdslabs/pinger.git
# or with HTTP
$ git clone https://github.com/sdslabs/pinger.git
$ cd pinger
```

Pinger uses GNU Make. We can run a simple make command to build the Pinger
binary as well as the documentation.

> **Note:** In most of what we will learn in this section of documentation,
> everything is done using `make`. If you're ever stuck, run `make help`
> and hopefully you'll find something useful.

```sh
$ make all
```

The above command might take a while but at the end you should have an
executable called `pinger` and documentation built in the `docs/book`
directory.

We should take a deeper look into the build process and understand how
to use the `Makefile` for our development flow.
