---
title: Making changes
parent: Development
nav_order: 3

---

# Making changes to the codebase

These are a few guidelines that need to be followed in-order to make changes
in the main repository.

## Vendoring

We vendor all the dependencies of our project. All the dependencies are
listed in `go.mod` file and are maintained in the `vendor/` directory.

When updating or adding or deleting a dependency, run the following to
keep the vendor up-to-date:

```sh
$ make vendor
```

## Protobufs

To install protobuf compiler, and other development tools, use:

```sh
$ make install
```

Run the following to compile protobufs:

```sh
$ make proto
```

When making changes to protobufs, make sure that you commit the compiled
proto files as-well.

## Dockerfile

When updating the `Dockerfile`, test building the image using:

```sh
$ make docker TAG="name:version" # default TAG="pinger:dev"

# Try running the container, should respond with pong
$ docker run --rm name:version
```

## Check for linting errors

We use [golangci-lint](https://golangci-lint.run/). Make sure your changes
pass the lint tests. To install golangci-lint, or any other development
dependencies, use:

```sh
$ make install
```

To check for linting errors, use:

```sh
$ make lint
```

Some errors can be fixed automatically. To do so, use:

```sh
$ make format
```

So, before committing any change, make sure there are no linting errors in
the code.

## Commits

Following are the guidelines for committing changes to the repository:

1. Break only logical changes into multiple commits. Commits such as "fix
   typo" or "address review commits" should be squashed into the one logical
   commit.

1. Each commit should individually pass tests and lint check. No separate
   commit should be made to fix these.

1. We don't have a defined commit message style for our codebase but the
   general idea is that the commit should include a heading, a body (if it's
   required) and reference to any issue that it might resolve. A good commit
   message looks something like this:

   ```
   Short commit heading with gist of changes.

   Body of commit trying to explain the old behaviour and how this commit
   changes it for the better.

   Resolves: #123

   Signed-off-by: Contributer <example@contributor.com>
   ```
