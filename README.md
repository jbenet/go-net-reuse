# go-net-reuse

This is a simple package to get around the problem of reusing addresses.
The go `net` package (to my knowledge) does not allow setting socket options.
This is particularly problematic when attempting to do TCP NAT holepunching,
which requires a process to both Listen and Dial on the same TCP port.
This package makes this possible for me. It is a pretty narrow use case, but
perhaps this package can grow to be more general over time.

## Notice

- Inspired by https://github.com/kavu/go_reuseport
- Currently only supports TCP.
- This could all be obviated by being able to set socket options on
  `net.Dialer` and `net.Listener`.

## Install

```
go get -u github.com/jbenet/go-net-reuse
```

## Usage

### Critical part

```Go
import reuse "github.com/jbenet/go-net-reuse"

l, _ := reuse.Listen("tcp", ":1111")
incoming, _ := l.Accept()

d := reuse.Dialer{LocalAddr: ":1111"}
outgoing, _ := d.Dial("tcp", ":2222")
```

### Full example

- [echonc](echonc/echonc.go) echo net-cat src

Try it out. In one terminal:

```sh
cd echonc
go build
./echonc :1111 :2222
```

In another:

```sh
./echonc :2222 :1111
```
