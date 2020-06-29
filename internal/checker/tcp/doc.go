// Copyright (c) 2020 SDSLabs
// Use of this source code is governed by an MIT license
// details of which can be found in the LICENSE file.

// Package tcp implements the TCP checker and prober.
//
// The TCP prober sends a TCP ECHO request and receives a reply from the
// address the request was sent to.
//
// Valid check format is described as following:
//
// Interval and Timeout should be greater than 0.
//
// Input:
//
// 	  Type          Value                  Description
// 	-------- -------------------- ------------------------------
// 	 "TCP"   "", "PING", "ECHO"   Sends and receives TCP ECHO
//
// Output:
//
// 	        Type                 Value             Description
// 	----------------------- ----------------- ------------------------
// 	 "TIMEOUT"				 <not validated>   Success is not-timeout
//   "MESSAGE"               <validated>   	   Success is user defined messages
//
// Target:
//
// 	   Type             Value                          Description
// 	----------- -------------------------- --------------------------------
// 	 "ADDRESS"   <valid HOST:PORT address>   Address to send the request to
//
// Payload is not required and hence not validated.
//
// Output:
//
//    Type                    Value                               Description
// 	------------- ---------------------------------- ------------------------------------------
//   "MESSAGE"         <anything>
//
package tcp
