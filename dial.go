package reuse

import (
	"errors"
	"fmt"
	"net"
	"os"
	"strconv"
	"syscall"

	sockaddrnet "github.com/jbenet/go-sockaddr/net"
)

type Dialer struct {
	LocalAddr net.Addr
}

func (d *Dialer) Dial(network, address string) (c net.Conn, err error) {

	addr, err := ResolveAddr(network, address)
	if err != nil {
		return nil, fmt.Errorf("ResolveAddr failed: %s", err)
	}

	fd, err := Socket(addr)
	if err != nil {
		return nil, err
	}

	if d.LocalAddr != nil {
		localSockaddr := sockaddrnet.NetAddrToSockaddr(d.LocalAddr)
		if localSockaddr == nil {
			err = errors.New("sockaddr conversion failed (local)")
			return
		}

		if err = syscall.Bind(fd, localSockaddr); err != nil {
			fmt.Println("bind failed")
			return nil, err
		}
	}

	// Set backlog size to the maximum
	if err = syscall.Connect(fd, *addr.Sockaddr); err != nil {
		fmt.Println("connect failed")
		return nil, err
	}

	// File Name get be nil
	file := os.NewFile(uintptr(fd), filePrefix+strconv.Itoa(os.Getpid()))
	if c, err = net.FileConn(file); err != nil {
		return nil, err
	}

	if err = file.Close(); err != nil {
		return nil, err
	}

	return c, err
}
