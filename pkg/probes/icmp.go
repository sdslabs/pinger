package probes

import (
	"errors"
	"math"
	"math/rand"
	"net"
	"syscall"
	"time"

	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
	"golang.org/x/net/ipv6"
)

// NewICMPProbe returns a `*ICMPProber` for given address and timeout.
// This prober can be used to check for ICMP requests.
func NewICMPProbe(address string, timeout time.Duration) (*ICMPProber, error) {
	pr := &ICMPProber{}

	if err := pr.SetAddress(address); err != nil {
		return nil, err
	}

	pr.SetTimeout(timeout)

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	pr.id = r.Intn(math.MaxInt16)

	return pr, nil
}

func isIPv4(ip net.IP) bool {
	return len(ip.To4()) == net.IPv4len
}

func isIPv6(ip net.IP) bool {
	return len(ip) == net.IPv6len
}

func timeToBytes(t time.Time) []byte {
	nsec := t.UnixNano()
	b := make([]byte, 8)
	for i := uint8(0); i < 8; i++ {
		b[i] = byte((nsec >> ((7 - i) * 8)) & 0xff)
	}
	return b
}

// ICMPProber is the prober for ICMP Echo checks.
type ICMPProber struct {
	timeout time.Duration
	address string
	ipv4    bool
	ipAddr  *net.IPAddr
	id      int
}

// SetAddress is a setter for address which is to be probed by the prober.
func (pr *ICMPProber) SetAddress(address string) error {
	ipAddr, err := net.ResolveIPAddr("ip", address)
	if err != nil {
		return err
	}
	var isipv4 bool
	if isIPv4(ipAddr.IP) {
		isipv4 = true
	} else if isIPv6(ipAddr.IP) {
		isipv4 = false
	} else {
		return errors.New("invalid address")
	}
	pr.address = address
	pr.ipv4 = isipv4
	pr.ipAddr = ipAddr
	return nil
}

// SetTimeout is a setter for timeout of prober.
func (pr *ICMPProber) SetTimeout(timeout time.Duration) {
	pr.timeout = timeout
}

// GetIPAddress returns the IP address of the prober address.
func (pr *ICMPProber) GetIPAddress() *net.IPAddr {
	return pr.ipAddr
}

// Probe method is used to execute the prober.
// It sends an ICMP echo request to given address
// and returns when receives a reply.
func (pr *ICMPProber) Probe() (*ICMPProbeResult, error) {
	var conn *icmp.PacketConn
	var err error
	if pr.ipv4 {
		conn, err = icmp.ListenPacket("udp4", "0.0.0.0")
		if err != nil {
			return nil, err
		}
		if err = conn.IPv4PacketConn().SetControlMessage(ipv4.FlagTTL, true); err != nil {
			return nil, err
		}
	} else {
		conn, err = icmp.ListenPacket("udp6", "::")
		if err != nil {
			return nil, err
		}
		if err = conn.IPv6PacketConn().SetControlMessage(ipv6.FlagHopLimit, true); err != nil {
			return nil, err
		}
	}
	defer conn.Close() //nolint:errcheck

	startTime := time.Now()

	timeoutResult := &ICMPProbeResult{
		Timeout: true,
		StartTime: startTime,
		Duration: pr.timeout,
	}

	if err := conn.SetDeadline(startTime.Add(pr.timeout)); err != nil {
		return nil, err
	}

	if err := pr.send(conn); err != nil {
		if errIsTimeout(err) {
			return timeoutResult, nil
		}

		return nil, err
	}

 	nb, ttl, err := pr.receive(conn);
	if err != nil {
		if errIsTimeout(err) {
			return timeoutResult, nil
		}

		return nil, err
	}

	return &ICMPProbeResult{
		Timeout: false,
		StartTime: startTime,
		Duration: time.Since(startTime),
		NumBytes: nb,
		TTL: ttl,
	}, nil
}

func (pr *ICMPProber) send(conn *icmp.PacketConn) error {
	var icmpType icmp.Type
	if pr.ipv4 {
		icmpType = ipv4.ICMPTypeEcho
	} else {
		icmpType = ipv6.ICMPTypeEchoRequest
	}

	var destination net.Addr = &net.UDPAddr{
		IP:   pr.ipAddr.IP,
		Zone: pr.ipAddr.Zone,
	}

	data := timeToBytes(time.Now()) // size of 8 bytes

	body := &icmp.Echo{
		ID:   pr.id,
		Seq:  0,
		Data: data,
	}

	message := &icmp.Message{
		Type: icmpType,
		Code: 0,
		Body: body,
	}

	msgBytes, err := message.Marshal(nil)
	if err != nil {
		return err
	}

	for {
		if _, err := conn.WriteTo(msgBytes, destination); err != nil {
			if neterr, ok := err.(*net.OpError); ok {
				if neterr.Err == syscall.ENOBUFS {
					continue
				}
			}
		}
		break
	}

	return nil
}

// Returns num bytes and time to live.
func (pr *ICMPProber) receive(conn *icmp.PacketConn) (n, t int, e error) {
	bytes := make([]byte, 512)
	var numBytes, ttl int
	var err error
	if pr.ipv4 {
		var cm *ipv4.ControlMessage
		numBytes, cm, _, err = conn.IPv4PacketConn().ReadFrom(bytes)
		if cm != nil {
			ttl = cm.TTL
		}
	} else {
		var cm *ipv6.ControlMessage
		numBytes, cm, _, err = conn.IPv6PacketConn().ReadFrom(bytes)
		if cm != nil {
			ttl = cm.HopLimit
		}
	}
	return numBytes, ttl, err
}

// ICMPProbeResult is the result of ICMP check probe.
type ICMPProbeResult struct {
	Timeout   bool
	StartTime time.Time
	Duration  time.Duration
	NumBytes  int
	TTL       int
}
