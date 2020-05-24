// Copyright (c) 2020 SDSLabs
// Use of this source code is governed by an MIT license
// details of which can be found in the LICENSE file.

// Package http implements the HTTP checker and prober.
//
// The HTTP prober creates a new request from the URL provided and sends
// the request with the specified method. This further is checked by the
// checker if the prober results match with the expected output and result
// is returned.
//
// Valid check format is described below:
//
// Input:
//
// 	  Type              Value              Description
// 	-------- --------------------------- ----------------
// 	 "HTTP"   "", "GET", "POST", "PUT",   Request method
// 	          "PATCH", "DELETE"
//
// Output:
//
// 	     Type                   Value                               Description
// 	-------------- -------------------------------- -------------------------------------------
// 	 "TIMEOUT"      <not validated>                  Success is not-timeout
// 	 "STATUSCODE"   <valid status code>              Response status code
// 	 "BODY"         <not validated>                  Response body should match this
// 	 "HEADER"       <valid header of format "K=V">   Header with key 'K' should have value 'V'
//
// Target:
//
// 	 Type           Value                 Description
// 	------- --------------------- ----------------------------
// 	 "URL"   <valid HTTP(S) URL>   URL to send the request to
//
// Payloads:
//
// 	    Type                    Value                               Description
// 	------------- ---------------------------------- ------------------------------------------
// 	 "HEADER"      <"K=V" formatted header>           Header with key='K' and value='V'
// 	 "PARAMETER"   <"K=V" formatted string where      This is passed on either as query (GET),
// 	               'V' is a valid JSON object, i.e,   or form-data/json body (depending on
// 	               string, bool, number or null>      value of "Content-Type" header)
//
package http
