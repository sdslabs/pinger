// Package oauth defines methods to setup oauth servers such as google and
// establishes status page authentication via JSON Web Token using the oauth.
//
// Example:
//     // Say `prov` is an OAuth Provider and `oauthRouter` is your gin router group.
//     if err := oauth.AddProvider(prov); err != nil {
//	       // handle error
//     }
//     if err := oauth.Setup(oauthRouter); err != nil {
//	       // handle error
//     }
//
// You can also use `oauth.Initialise` as a shorthand for above.
package oauth
