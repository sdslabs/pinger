// Package dns implements the DNS checker and prober.
//
// The DNS prober looks up the addresses for the host name and reports them
// to the checker which further checks if the output from the prober is
// desired.
//
// Valid check format is described as following:
//
// Interval and Timeout should be greater than 0.
//
// Input:
//
//	 Type          Value               Description
//	------- -------------------- ------------------------
//	 "DNS"   "", "PING", "ECHO"   Resolves the DNS names
//
// Output:
//
//	       Type               Value                       Description
//	----------------- -------------------- ----------------------------------------
//	 "TIMEOUT"         <not validated>      Success is not-timeout
//	 "ADDRESS", "IP"   <valid IP address>   One of the resolved IP should be this
//
// Target:
//
//	             Type                       Value               Description
//	------------------------------- ---------------------- ---------------------
//	 "HOST", "HOSTNAME", "DNSNAME"   <validated DNS name>   DNS name to resolve
//
// Payload is not required and hence not validated.
package dns
