# Submitting a Pull Request

These are a few guidelines that need to be followed in-order to make
changes in the main repository.

## Pre-commit

There are a few things to take care of before committing your changes.

### Vendoring

All the dependencies in our repository are vendored. All the dependencies
are listed in `go.mod` file and are maintained in the `vendor` directory.
When updating/adding/deleting a dependency, run the following to keep the
`vendor` up-to-date:

```sh
$ make vendor
```

This will cleanup unused dependencies from `go.mod` as well as add new
dependencies to `vendor`.

### Lint

We use [golangci-lint](https://golangci-lint.run/) to lint our code. Each
commit should pass all the lint checks. To ensure that it does, use the
`make lint` command. Many times, these errors might be related to
formatting of the code. These are usually auto-fixable. `make fmt` fixes
all the errors that might be resolved automatically.

```sh
# To check for linting errors
$ make lint
# To fix auto-fixable errors
$ make format
```

### Other checks

Apart from the fact that the code should build successfully, we need to
take care of following:

- Ensure that you have compiled any protobuf that may have been updated.
- The Docker image should be built successfully as well.

## Commits

Break only logical changes into multiple commits. Commits such as "fix 
typo" or "address review commits" should be squashed into the one
logical commit. Each commit should individually pass tests and lint check
No separate commit should be made to fix these.

### Commit Messages

We don't have a defined commit message style for our codebase, but the
general idea is that the commit should include a heading, a body (if it's
required) and reference to any issue that it might resolve. A good commit
message looks something like this:

```
scope: Short commit heading with gist of changes.

Body of commit trying to explain the old behaviour and how this commit
changes it for the better.

Resolves: #123

Signed-off-by: Contributer <example@contributor.com>
```

## Pull Request

Pull requests follows same guidelines as mentioned for the commit messages.
The title should give a clear idea of what it is about, followed by a
descriptive body and should mention what issue it resolves (if any).
