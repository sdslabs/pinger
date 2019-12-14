package probes

import (
	"errors"
	"math"
	"math/rand"
	"net"
	"sync"
	"syscall"
	"time"

	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
	"golang.org/x/net/ipv6"
)

// NewICMPProbe returns a `*ICMPProber` for given address and timeout.
// This prober can be used to check for ICMP requests.
func NewICMPProbe(address string, timeout time.Duration) (*ICMPProber, error) {
	pr := &ICMPProber{
		quit: make(chan bool),
	}
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
	quit    chan bool
	id      int
}

// SetAddress is a setter for address which is to be probed by the prober.
func (pr *ICMPProber) SetAddress(address string) error {
	ipAddr, err := net.ResolveIPAddr("ip", address)
	if err != nil {
		return err
	}
	var ipv4 bool
	if isIPv4(ipAddr.IP) {
		ipv4 = true
	} else if isIPv6(ipAddr.IP) {
		ipv4 = false
	} else {
		return errors.New("invalid address")
	}
	pr.address = address
	pr.ipv4 = ipv4
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
			close(pr.quit)
			return nil, err
		}
		conn.IPv4PacketConn().SetControlMessage(ipv4.FlagTTL, true)
	} else {
		conn, err = icmp.ListenPacket("udp6", "::")
		if err != nil {
			close(pr.quit)
			return nil, err
		}
		conn.IPv6PacketConn().SetControlMessage(ipv6.FlagHopLimit, true)
	}
	defer conn.Close()

	var wg sync.WaitGroup

	recv := make(chan *packet, 1)
	defer close(recv)

	wg.Add(1) // receive go routine
	go pr.receive(conn, recv, &wg)

	startTime := time.Now()

	// send request
	if err := pr.send(conn); err != nil {
		return nil, err
	}

	timer := time.NewTimer(pr.timeout)
	defer timer.Stop()

	for {
		select {
		case <-pr.quit:
			pr.waitAndClose(&wg)
			return nil, errors.New("error while receiving")
		case <-timer.C:
			pr.quit <- true
			pr.waitAndClose(&wg)
			return &ICMPProbeResult{
				Timeout:   true,
				StartTime: startTime,
				Duration:  pr.timeout,
			}, nil
		case r := <-recv:
			pr.quit <- true
			pr.waitAndClose(&wg)
			return &ICMPProbeResult{
				Timeout:   false,
				StartTime: startTime,
				Duration:  time.Since(startTime),
				NumBytes:  r.numBytes,
				TTL:       r.ttl,
			}, nil
		}
	}
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

func (pr *ICMPProber) receive(conn *icmp.PacketConn, recv chan<- *packet, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		select {
		case <-pr.quit:
			return
		default:
			bytes := make([]byte, 512)
			conn.SetReadDeadline(time.Now().Add(time.Millisecond * 100))
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
			if err != nil {
				if neterr, ok := err.(*net.OpError); ok {
					if neterr.Timeout() {
						continue
					} else {
						pr.quit <- true
						return
					}
				}
			}
			recv <- &packet{bytes: bytes, numBytes: numBytes, ttl: ttl}
		}
	}
}

func (pr *ICMPProber) waitAndClose(wg *sync.WaitGroup) {
	wg.Wait()
	close(pr.quit)
}

// ICMPProbeResult is the result of ICMP check probe.
type ICMPProbeResult struct {
	Timeout   bool
	StartTime time.Time
	Duration  time.Duration
	NumBytes  int
	TTL       int
}

type packet struct {
	bytes    []byte
	numBytes int
	ttl      int
}
