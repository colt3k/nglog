package ng

import (
	"fmt"
	"net"
	"strings"
)

//******************* Socket APPENDER ********************

type SocketAppender struct {
	*OutAppender
	host    string
	port    string
	contype string
}

func NewSocketAppender(filter, server, port string) (*SocketAppender, error) {
	if len(server) == 0 {
		return nil, fmt.Errorf("server required")
	}
	if len(port) == 0 {
		return nil, fmt.Errorf("port required")
	}

	oa := newOutAppender(filter, "")
	t := new(SocketAppender)
	t.OutAppender = oa
	t.host = server
	t.port = port
	t.contype = "tcp"

	return t, nil
}
func (f *SocketAppender) Name() string {
	if len(f.name) > 0 {
		return f.name
	}
	return fmt.Sprintf("%T", f)
}
func (f *SocketAppender) DisableColor() bool {
	return f.disableColor
}
func (f *SocketAppender) Applicable(msg string) bool {
	if f.filter == "*" {
		return true
	}
	if strings.Index(msg, f.filter) > -1 {
		return true
	}
	return false
}

func (f *SocketAppender) Process(msg []byte) {
	// Send via Socket

	address := f.host + ":" + f.port

	tcpAdr, err := net.ResolveTCPAddr(f.contype, address)
	if err != nil {
		fmt.Printf("failed to resolve address: %s on %s\n%+v", address, f.contype, err)
	}
	con, err := net.DialTCP(f.contype, nil, tcpAdr)
	if err != nil {
		fmt.Printf("dial failed on: %s\n%+v", address, err)
	}

	if con != nil {
		_, err = con.Write([]byte(strings.TrimSpace(string(msg))))
		if err != nil {
			fmt.Printf("write to server %s failed\n%+v", address, err)
			if con != nil {
				err = con.Close()
				if err != nil {
					fmt.Printf("issue closing %+v", err)
				}
			}
		}
		if con != nil {
			err = con.Close()
			if err != nil {
				fmt.Printf("issue closing %+v", err)
			}
		}
	}
}
