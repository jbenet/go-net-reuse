package reuse

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"syscall"
)

// Listen returns net.Listener created from a file discriptor for a socket
// with SO_REUSEPORT and SO_REUSEADDR option set.
func Listen(network, address string) (l net.Listener, err error) {

	addr, err := ResolveAddr(network, address)
	if err != nil {
		return nil, fmt.Errorf("ResolveAddr failed: %s", err)
	}

	fd, err := Socket(addr)
	if err != nil {
		return nil, err
	}

	if err = syscall.Bind(fd, *addr.Sockaddr); err != nil {
		return nil, err
	}

	// Set backlog size to the maximum
	if err = syscall.Listen(fd, syscall.SOMAXCONN); err != nil {
		return nil, err
	}

	// File Name get be nil
	file := os.NewFile(uintptr(fd), filePrefix+strconv.Itoa(os.Getpid()))
	if l, err = net.FileListener(file); err != nil {
		return nil, err
	}

	if err = file.Close(); err != nil {
		return nil, err
	}

	return l, err
}
