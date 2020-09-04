---
title: Directory structure
parent: Development
nav_order: 4

---

# Directory structure

```
.
├── cmd
│   └── pinger
├── docs
├── pkg
├── target
└── vendor
```

## Root

Root directory contains all the build related files or the meta-data for the
project/repository. It contains files varying from project's `Dockerfile` to
linter config `.golangci.yml`.

## Command line

`cmd/` contains the commands that can be executed using the binary. It is
the building block package for the Pinger CLI.

`cmd/pinger/` contains the main (or entry-point) for the CLI.

## Documentation

`docs/` contains the documentation for the project. Documentation is
generated using the [Jekyll](https://jekyllrb.com/) framework. To know more
see [how to build documentation](./documentation.html).

## Target

`target/` contains any of the builds, including libraries and binaries for
the project.

## Vendor

`vendor/` contains the source code of dependencies of the project. See
[vendoring](./making-changes.html#vendoring) for more details.
