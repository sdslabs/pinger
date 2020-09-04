---
title: Documentation
parent: Development
nav_order: 5

---

# Documentation

It's appericiated that you document features or release changes along-with
your code change itself.

## Building documentation

We use [Jekyll](https://jekyllrb.com/) to generate the documentation.

### Prerequisites

Make sure you have these installed before you can build docs.

- [Ruby](https://www.ruby-lang.org/) >= 2.6
- [Bundler](https://bundler.io/)

### Build

1. Install dependencies:
   ```sh
   $ make docs-install
   ```

1. Build documentation:
   ```sh
   $ make docs-build
   ```

<span class="label">Tip</span>
The above steps can be done together using `make docs`.

When developing it's helpful to build on file change.

```sh
# Build in watch mode
$ make docs-watch

# Serve documentation on :4000
$ make docs-serve
```
