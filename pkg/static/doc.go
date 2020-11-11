// Package static includes the blob for generated files that lie in `static`
// directory of the project and other utilities to access them.
package static

// TODO(vrongmeal): When running in debug mode the resources should not be
// built again and again. This will consume more time. Rather just build the
// binary once and use file system in that case. When writing wrappers for
// using the resources from built content, make sure to include option to
// use file system to serve files.
