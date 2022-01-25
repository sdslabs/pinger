# Build Process

We won't be running `make all` everytime. Since the command builds all the
components, it would take longer than usual when changing a couple lines in
one part of the code.

Currently we need to build three things:

1. Binary
1. Docker Image
1. Documentation

Let's take a look at both of them individually.

## Binary

The binary packs all the static files, including the CSS and JavaScript files,
required for the frontend. Once we have our resources ready, we can build the
executable using the `make build` command. We have `VERSION` flag for setting
the version of the binary.

```sh
# Building the binary
$ make build VERSION=1.0.1
```

The binary can be built in development mode which enables serving static files
from the file-system itself, instead of embedding them. It's works better for
development workflow since changes in static files don't require to rebuild the
Pinger executable.

```sh
$ make build DEV=on
```

> **Note:** This does pose the constraint that any Pinger command has to
> be executed in the repository's root directory. This shouldn't be an
> issue while developing.

### Protobuf

When making changes to a `.proto` file, we need to compile it into the
equivalent Go code. These generated files (`*.pb.go`) need to be committed.

```sh
$ make proto
```

## Docker Image

`Dockerfile` for the image is present in the root directory. To build the
image, use `make docker` command. This generates an image with tag –
`pinger:dev`. To change the tag, we can use the `TAG` option for make
command.

```sh
$ make docker TAG="pinger:v1.2.3"
```

> **Note:** In case of docker image, version is extracted from the tag.
>For example: in the aforementioned case, version of the binary will 
>set to `v1.2.3`.

## Documentation

Documentation is not (yet) packaged into the binary. We currently host it
using [Github Pages](https://pages.github.com/). We use
[mdBook](https://rust-lang.github.io/mdBook/) as the documentation
framework. Make sure you have that installed in your `$PATH`.

You can build the documentation using `make docs` command. This generates
the documentation in `docs/book` directory. While developing you might
require to watch the `docs` directory. We can again use the `DEV=on`
option in that case.

```sh
# To build the documentation
$ make docs
# Builds the documentation, watches for changes and serves on :3000
$ make docs DEV=on
```
