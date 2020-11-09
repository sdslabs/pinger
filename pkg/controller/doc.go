// Package controller implements controller and manager.
//
// Controllers can be used to run a specific task repatitively at regular
// intervals of time. A controller stores the resulting stats of runs with
// itself until pulled from it.
//
// A manager is used to orchestrate multiple controllers running
// concurrently.
package controller
