// Package ws implements the WS checker and prober.
//
// The WS prober sends a TCP ECHO request and receives a reply from the
// address the request was sent to.
//
// Valid check format is described as following:
//
// Interval and Timeout should be greater than 0.
//
// Input:
//
// 	  Type          Value                  Description
// 	-------- -------------------- ----------------------------
// 	 "WS"     "", "PING", "ECHO"   Sends and receives WS ECHO
//
// Output:
//
//       Type                   Value                              Description
// 	-------------- -------------------------------- ------------------------------------------
// 	 "TIMEOUT"      <not validated>                  Success is not-timeout
// 	 "STATUSCODE"   <valid status code>              Response status code
// 	 "BODY"         <not validated>                  Response body should match this
// 	 "HEADER"       <valid header of format "K=V">   Header with key 'K' should have value 'V'
// 	 "MESSAGE"      <non empty messages              Messages are split by "\n---\n".
// 	                separated by "\n---\n">          For "hello" and "world" as two messages,
// 	                                                 the output value should be
// 	                                                  hello
// 	                                                  ---
// 	                                                  world
// 	                                                 All messages should be equal in order
//
// Target:
//
// 	  Type          Value                Description
// 	-------- ------------------- ----------------------------
// 	 "URL"    <valid WS(S) URL>   URL to send the request to
//
// Payload:
//
// 	   Type               Value                        Description
// 	----------- -------------------------- -----------------------------------
// 	 "MESSAGE"   <non empty string>         Sent in order as in array
// 	 "HEADER"    <"K=V" formatted header>   Header with key='K' and value='V'
//
package ws
