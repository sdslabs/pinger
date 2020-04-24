// Package request contains the structs that can be used to bind request JSON.
package request

// HTTPUserUpdate is the request used to update the user.
// This only contains the content to be updated.
type HTTPUserUpdate struct {
	Name string `json:"name" example:"Abc Xyz"`
}
