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
//	Type               Value                          Description
//	----------- ------------------------- ------------------------------------------
//	 "TIMEOUT"   <not validated>           Success in not-timeout
//	 "MESSAGE"   <non empty messages       Messages are split by "\n---\n".
//	             separated by "\n---\n">   For "hello" and "world" as two messages,
//	                                       the output value should be
//	                                         hello
//	                                         ---
//	                                         world
//	                                       All messages should be equal in order
//
// Target:
//
//	   Type             Value                          Description
//	----------- -------------------------- --------------------------------
//	 "ADDRESS"   <valid HOST:PORT address>   Address to send the request to
//
// Payload:
//
//	   Type            Value                 Description
//	----------- -------------------- ---------------------------
//	 "MESSAGE"   <non empty string>   Sent in order as in array
//
package tcp
