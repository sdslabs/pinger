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

Each commit message consists of a header, a body and a footer. The header
has a special format that includes a scope and a subject:

```
<scope>: <subject>
<BLANK LINE>
<body>
<BLANK LINE>
<footer>
```

Any line of the commit message cannot be longer 76 characters! This allows
the message to be easier to read on github as well as in various git tools.

#### Scope

The scope could be anything specifying place/context of the commit change.
When making changes to Go code, this will most likely be the package name
where the changes are made. For documentation changes scope will be `docs`
and for UI changes scope will be `ui`.

> **Note:** Scope should not be capitalized. `docs` should not be spelled
> _`Docs`_.

#### Subject

The subject contains succinct description of the change:

- Use the imperative, present tense: "change" not "changed" nor "changes".
- Do capitalize first letter.
- Add a dot (.) at the end.

#### Body

The body should include the motivation for the change and contrast this with
previous behavior. Just as in the **subject**, try to use the imperative,
present tense: "change" not "changed" nor "changes", though, this rule is
not enforced for the body.

#### Footer

The footer should contain any information about **Breaking Changes** and is
also the place to reference issues that this commit **Closes**.

An example of a good commit message would be:

```
exporter: Update timescale exporter to use Gorm v2.

Previously timescale exporter used Gorm v2 which did not support batch
insert resulting in raw SQL queries. Gorm v2 now includes batch insert,
hence, the same is used to refactor the exporter queries.

Closes #123

Signed-off-by: Contributer <example@contributor.com>
```

## Pull Request

Pull requests follows same guidelines as mentioned for the commit messages.
The title should give a clear idea of what it is about, followed by a
descriptive body and should mention what issue it resolves (if any).
