// Copyright (c) 2020 SDSLabs
// Use of this source code is governed by an MIT license
// details of which can be found in the LICENSE file.

// Package icmp implements the ICMP checker and prober.
//
// The ICMP prober sends an ICMP ECHO request and receives a reply from the
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
// 	 "ICMP"   "", "PING", "ECHO"   Sends and receives ICMP ECHO
//
// Output:
//
// 	   Type           Value             Description
// 	----------- ----------------- ------------------------
// 	 "TIMEOUT"   <not validated>   Success is not-timeout
//
// Target:
//
// 	   Type                    Value                          Description
// 	----------- ----------------------------------- --------------------------------
// 	 "ADDRESS"   <valid address (can be hostname)>   Address to send the request to
//
// Payload is not required and hence not validated.
//
package icmp
