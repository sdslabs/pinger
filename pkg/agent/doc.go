// Copyright (c) 2020 SDSLabs
// Use of this source code is governed by an MIT license
// details of which can be found in the LICENSE file.

// Package agent runs a manager that runs various checks inside controllers.
//
// An agent can run in two modes, standalone mode and along with the central
// server. An agent technically runs independently, but to run with the
// central server it exposes a GRPC API which takes requests to push new
// checks and remove existing checks from the manager.
package agent
