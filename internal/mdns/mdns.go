package mdns

import (
	"fmt"
	"net"

	"github.com/pion/mdns/v2"
	"golang.org/x/net/ipv4"
	"golang.org/x/net/ipv6"
)

type MdnsServer struct {
	domainName string
	addr4      *net.UDPAddr
	addr6      *net.UDPAddr
	l4         *net.UDPConn
	l6         *net.UDPConn
}

func NewMDNsServer(domainName string) *MdnsServer {
	addr4, err := net.ResolveUDPAddr("udp4", mdns.DefaultAddressIPv4)
	if err != nil {
		panic(err)
	}

	addr6, err := net.ResolveUDPAddr("udp6", mdns.DefaultAddressIPv6)
	if err != nil {
		panic(err)
	}

	l4, err := net.ListenUDP("udp4", addr4)
	if err != nil {
		panic(err)
	}

	l6, err := net.ListenUDP("udp6", addr6)
	if err != nil {
		panic(err)
	}
	return &MdnsServer{
		domainName: domainName,
		addr4:      addr4,
		addr6:      addr6,
		l4:         l4,
		l6:         l6,
	}
}
func (m *MdnsServer) Serve(domain string) {
	_, err := mdns.NewServer(
		ipv4.NewPacketConn(m.l4),
		ipv6.NewPacketConn(m.l6),
		mdns.WithLocalNames(domain),
	)
	if err != nil {
		panic(err)
	}
	fmt.Println("Listening on " + domain)
}
