package fishfinger

import (
	"bufio"
	"fmt"
	"net"
	"time"
)

// -----------------------------------------------------------------------------

// SocketBackoff is an implementation of a `backoffFunc`. It
// retries to connect to the tcp port 9090 of `service` until the remote
// socket sends a message back.
//
// NOTE: This function does not return until the socket is reached or an `error`
// occurs. FishFinger project offers a default implementation of a
// `backoffFunc` that will work for the most simple of use cases, although I
// strongly recommend you provide your own, safer implementation while doing
// real work.
func SocketBackoff(c *Compose, service string, port string) error {
	var (
		msg  string
		conn net.Conn
	)

	addr, err := c.Port(service, port)
	if err != nil {
		return err
	}
	for ; msg != "ready\n"; time.Sleep(5 * time.Second) {
		if conn, err = net.Dial("tcp", addr); err == nil {
			fmt.Fprintf(conn, "ready\n")
			msg, _ = bufio.NewReader(conn).ReadString('\n')
			conn.Close()
		}
		fmt.Printf("Retry connection in 5s.\n")
	}
	return nil
}
