package icmp

import (
	"context"
	"fmt"
	"math"
	"math/rand"
	"net"
	"syscall"
	"time"

	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
	"golang.org/x/net/ipv6"

	"github.com/sdslabs/pinger/pkg/checker"
)

// Prober sends an ICMP ECHO request to the given address.
type Prober struct {
	addr    string
	timeout time.Duration
}

// NewProber creates a prober to send an ICMP ECHO ping to the provided
// address.
func NewProber(addr string, timeout time.Duration) (*Prober, error) {
	if timeout <= 0 {
		return nil, fmt.Errorf("timeout should be > 0")
	}

	return &Prober{
		addr:    addr,
		timeout: timeout,
	}, nil
}

// Probe sends an ICMP ECHO request to the given address and returns the
// result.
func (p *Prober) Probe(ctx context.Context) (*ProbeResult, error) {
	startTime := time.Now()

	timeoutResult := &ProbeResult{
		StartTime: startTime,
		Duration:  p.timeout,
		Timeout:   true,
	}

	baseCtx := ctx
	if ctx == nil {
		baseCtx = context.Background()
	}

	probeCtx, cancel := context.WithTimeout(baseCtx, p.timeout)
	defer cancel()

	addr, isIPv4, err := p.resolveAddress(probeCtx)
	if err != nil {
		if checker.ErrIsTimeout(err) {
			return timeoutResult, nil
		}

		return nil, err
	}

	var conn *icmp.PacketConn
	if isIPv4 {
		conn, err = icmp.ListenPacket("udp4", "0.0.0.0")
		if err != nil {
			return nil, err
		}
		err = conn.IPv4PacketConn().SetControlMessage(ipv4.FlagTTL, true)
		if err != nil {
			return nil, err
		}
	} else {
		conn, err = icmp.ListenPacket("udp6", "::")
		if err != nil {
			return nil, err
		}
		err = conn.IPv6PacketConn().SetControlMessage(ipv6.FlagHopLimit, true)
		if err != nil {
			return nil, err
		}
	}
	defer conn.Close() // nolint:errcheck

	nSend, err := p.sendMessage(probeCtx, conn, addr, isIPv4, startTime)
	if err != nil {
		if checker.ErrIsTimeout(err) {
			return timeoutResult, nil
		}

		return nil, err
	}

	nRecv, ttl, err := p.receiveMessage(probeCtx, conn, isIPv4, startTime)
	if err != nil {
		if checker.ErrIsTimeout(err) {
			return timeoutResult, nil
		}

		return nil, err
	}

	return &ProbeResult{
		Timeout:          false,
		StartTime:        startTime,
		Duration:         time.Since(startTime),
		NumBytesSent:     nSend,
		NumBytesReceived: nRecv,
		TTL:              ttl,
	}, nil
}

// resolveAddress resolves the prober IP address.
func (p *Prober) resolveAddress(ctx context.Context) (addr *net.IPAddr, isIPv4 bool, err error) {
	// we resolve the ip address in a go routine with the time deadline matching
	// our probe startTime + timeout. This is because the lookup doesn't accept
	// context and we need to handle timeout manually in this case. DNS Lookup
	// is going to end eventually and in reasonable time, this is just to handle
	// the case when it exceeds our timeout.
	type ipaddr struct {
		addr *net.IPAddr
		err  error
	}
	addrChan := make(chan ipaddr)
	go func(address string, stream chan<- ipaddr) {
		a, e := net.ResolveIPAddr("ip", address)
		stream <- ipaddr{addr: a, err: e}
	}(p.addr, addrChan)

	select {
	case <-ctx.Done():
		return nil, false, ctx.Err()

	case ipa := <-addrChan:
		err = ipa.err
		if err != nil {
			return nil, false, err
		}

		addr = ipa.addr
		isIPv4 = len(addr.IP.To4()) == net.IPv4len

		// check whether the ip address is correct, i.e, either of ipv4 or ipv6.
		if !isIPv4 && len(addr.IP) != net.IPv6len {
			err = fmt.Errorf("invalid IP address: %s", addr.String())
		}

		return
	}
}

// sendMessage sends an echo message to IP address.
func (p *Prober) sendMessage(
	ctx context.Context,
	conn *icmp.PacketConn,
	addr *net.IPAddr,
	isIPv4 bool,
	startTime time.Time,
) (n int, err error) {
	var msgType icmp.Type
	if isIPv4 {
		msgType = ipv4.ICMPTypeEcho
	} else {
		msgType = ipv6.ICMPTypeEchoRequest
	}

	destination := &net.UDPAddr{
		IP:   addr.IP,
		Zone: addr.Zone,
	}

	data := timeToBytes(startTime)
	msg := &icmp.Message{
		Type: msgType,
		Code: 0,
		Body: &icmp.Echo{
			ID: rand.New( // nolint:gosec
				rand.NewSource(startTime.UnixNano()),
			).Intn(math.MaxInt16),
			Seq:  0,
			Data: data,
		},
	}

	bin, err := msg.Marshal(nil)
	if err != nil {
		return 0, err
	}

	err = conn.SetWriteDeadline(startTime.Add(p.timeout))
	if err != nil {
		return 0, err
	}

	errChan := make(chan error)
	go send(conn, bin, destination, errChan)

	select {
	case <-ctx.Done():
		// no need to exit the routing since there is a deadline on write
		// already set and context canceled either due to timeout or closing of
		// agent which does eventually exit terminate the thread.
		return 0, ctx.Err()
	case err = <-errChan:
		return len(data), err
	}
}

// receiveMessage receives the reply from the address.
func (p *Prober) receiveMessage(
	ctx context.Context,
	conn *icmp.PacketConn,
	isIPv4 bool,
	startTime time.Time,
) (n, ttl int, err error) {
	err = conn.SetReadDeadline(startTime.Add(p.timeout))
	if err != nil {
		return 0, 0, err
	}

	packetChan := make(chan packet)
	go receive(conn, isIPv4, packetChan)

	select {
	case <-ctx.Done():
		// no need to exit the routing since there is a deadline on read
		// already set and context canceled either due to timeout or closing of
		// agent which does eventually exit terminate the thread.
		return 0, 0, ctx.Err()
	case pack := <-packetChan:
		n = pack.n
		ttl = pack.ttl
		err = pack.err
		return
	}
}

// ProbeResult is the result of an ICMP probe.
type ProbeResult struct {
	Timeout          bool
	StartTime        time.Time
	Duration         time.Duration
	NumBytesSent     int
	NumBytesReceived int
	TTL              int
}

// send sends the ICMP ECHO request to the address.
func send(conn *icmp.PacketConn, message []byte, destination net.Addr, stream chan<- error) {
	for {
		_, err := conn.WriteTo(message, destination)
		if err != nil {
			if netErr, ok := err.(*net.OpError); ok {
				if netErr.Err == syscall.ENOBUFS {
					continue
				}
			}

			stream <- err
			return
		}

		stream <- nil
		return
	}
}

// packet represents the ICMP ECHO REPLY.
type packet struct {
	n   int
	ttl int
	err error
}

// receive accepts and reads the message and passes it on to the receiver
// stream.
func receive(conn *icmp.PacketConn, isIPv4 bool, stream chan<- packet) {
	bytes := make([]byte, 512)
	var n, ttl int
	var err error

	if isIPv4 {
		var cm *ipv4.ControlMessage
		n, cm, _, err = conn.IPv4PacketConn().ReadFrom(bytes)
		if cm != nil {
			ttl = cm.TTL
		}
	} else {
		var cm *ipv6.ControlMessage
		n, cm, _, err = conn.IPv6PacketConn().ReadFrom(bytes)
		if cm != nil {
			ttl = cm.HopLimit
		}
	}

	stream <- packet{
		n:   n,
		ttl: ttl,
		err: err,
	}
}

// timeToBytes converts the given time to bytes.
func timeToBytes(t time.Time) []byte {
	nsec := t.UnixNano()
	b := make([]byte, 8)
	for i := uint8(0); i < 8; i++ {
		b[i] = byte((nsec >> ((7 - i) * 8)) & 0xff)
	}
	return b
}
