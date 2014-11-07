package reuse

import (
	"errors"
	"net"
	"syscall"

	// sockaddr "github.com/jbenet/go-sockaddr"
	sockaddrnet "github.com/jbenet/go-sockaddr/net"
)

const (
	tcp4                  = 52 // "4"
	tcp6                  = 54 // "6"
	unsupportedProtoError = "Only tcp4 and tcp6 are supported"
	filePrefix            = "port."
)

// params converts a given net, addr string pair and turns it into a
// (*net.Addr, IPPROTO_ type, SOCK_ type, AFNET_ type)

type Addr struct {
	// these are the typical inputs
	Network string
	Address string

	Netaddr  net.Addr
	Sockaddr *syscall.Sockaddr
	IPPROTO  int
	SOCK     int
	AF       int
}

func Socket(addr Addr) (int, error) {
	fd, err := syscall.Socket(addr.AF, addr.SOCK, addr.IPPROTO)
	if err != nil {
		return fd, err
	}

	if err = syscall.SetsockoptInt(fd, syscall.SOL_SOCKET, syscall.SO_REUSEADDR, 1); err != nil {
		return fd, err
	}

	if err = syscall.SetsockoptInt(fd, syscall.SOL_SOCKET, syscall.SO_REUSEPORT, 1); err != nil {
		return fd, err
	}

	return fd, nil
}

func ResolveAddr(network, address string) (addr Addr, err error) {
	addr.Network = network
	addr.Address = address

	addr.Netaddr, err = ResolveNetAddr(network, address)
	if err != nil {
		return
	}

	*addr.Sockaddr = sockaddrnet.NetAddrToSockaddr(addr.Netaddr)
	if addr.Sockaddr == nil {
		err = errors.New("sockaddr conversion failed")
		return
	}

	addr.IPPROTO = sockaddrnet.NetAddrIPPROTO(addr.Netaddr)
	if addr.IPPROTO == -1 {
		err = errors.New("unknown IPPROTO type")
		return
	}

	addr.SOCK = sockaddrnet.NetAddrSOCK(addr.Netaddr)
	if addr.SOCK == 0 {
		err = errors.New("unknown SOCK type")
		return
	}

	addr.AF = sockaddrnet.NetAddrAF(addr.Netaddr)
	if addr.AF == 0 {
		err = errors.New("unknwon AF type")
		return
	}
	return
}

func ResolveNetAddr(network, address string) (net.Addr, error) {
	switch network {
	default:
		return nil, errors.New("unsupported network")
	case "ip", "ip4", "ip6":
		return net.ResolveIPAddr(network, address)
	case "tcp", "tcp4", "tcp6":
		return net.ResolveTCPAddr(network, address)
	case "udp", "udp4", "udp6":
		return net.ResolveUDPAddr(network, address)
	case "unix", "unixgram", "unixpacket":
		return net.ResolveUnixAddr(network, address)
	}
}
