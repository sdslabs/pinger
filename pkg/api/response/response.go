// Package response contains typed structs for various responses to api requests.
// These can be used comfortably to produce swagger api documentation.
package response

import "time"

// HTTPError is the response type for any error that occur for a HTTP request.
type HTTPError struct {
	Error string `json:"error" example:"JWT is expired"`
}

// HTTPInternalServerError returns an HTTPError with internal server error message.
// For internal server errors we don't want to present the errors to the client so
// the error message is constant. Instead, log the errors in these cases.
var HTTPInternalServerError = HTTPError{Error: "INTERNAL SERVER ERROR"}

// HTTPLogin is the response for login url request using OAuth service.
type HTTPLogin struct {
	LoginURL string `json:"login_url" example:"https://accounts.google.com/auth/..."`
}

// HTTPAuthorization is the response for login / register request.
type HTTPAuthorization struct {
	Token     string        `json:"token" example:"abcdefghijklmnopqrstuvwxyz1234567890"`
	ExpiresIn time.Duration `json:"expires_in" example:"3600"` // in seconds
	UserID    uint          `json:"user_id" example:"10"`
	UserEmail string        `json:"user_email" example:"myname@example.com"`
	UserName  string        `json:"user_name" example:"Go Pher"`
}

// HTTPRefreshToken is the response when client asks for refreshing the old token.
type HTTPRefreshToken struct {
	Token     string        `json:"token" example:"abcdefghijklmnopqrstuvwxyz1234567890"`
	ExpiresIn time.Duration `json:"expires_in" example:"3600"` // in seconds
}

// HTTPUserInfo is the response when client fetches or updates user.
type HTTPUserInfo struct {
	ID    uint   `json:"id" example:"123"`
	Email string `json:"Email" example:"abc@xyz.com"`
	Name  string `json:"Name" example:"Abc Xyz"`
}

// HTTPEmpty returns an empty response.
type HTTPEmpty struct{}
