package ipc

import (
	"fmt"
	"github.com/meowalien/go-meowalien-lib/errs"
	"net"
	"os"
)

func NewIPCListener(addr string) (lst net.Listener, err error) {
	fmt.Println("addr: ", addr)
	_ = os.Remove(addr)
	uaddr, err := net.ResolveUnixAddr("unix", addr)
	if err != nil {
		err = errs.New(err)
		return
	}
	lst, err = net.ListenUnix("unix", uaddr)
	if err != nil {
		err = errs.New(err)
		return
	}
	return
}

func NewIPCClient(addr string) (conn net.Conn, err error) {
	unixAddress, err := net.ResolveUnixAddr("unix", addr)
	if err != nil {
		err = errs.New(err)
		return
	}
	conn, err = net.DialUnix("unix", nil, unixAddress)
	if err != nil {
		err = errs.New(err)
		return
	}
	return
}
