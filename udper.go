package udper

import (
	"net"
	"time"
)

var (
	addr     *string
	listener *net.UDPConn
	// Timeout specifies the connection timeout
	Timeout = time.Millisecond
)

type fn func()

// SetAddr sets the UDP port that will be listened on.
func SetAddr(a string) {
	addr = &a
}

// Start starts the connection
func Start() error {
	resAddr, err := net.ResolveUDPAddr("udp", *addr)
	if err != nil {
		return err
	}
	listener, err = net.ListenUDP("udp", resAddr)
	if err != nil {
		return err
	}

	return nil
}

// Stop ends the connection
func Stop() error {
	return listener.Close()
}

// Get returns the message after the execution
func Get(body fn) (string, error) {
	err := Start()

	if err != nil {
		return "", err
	}

	defer Stop()

	body()

	message := make([]byte, 1024*32)
	var bufLen int
	for {
		listener.SetReadDeadline(time.Now().Add(Timeout))
		n, _, _ := listener.ReadFrom(message[bufLen:])
		if n == 0 {
			break
		} else {
			bufLen += n
		}
	}

	return "" + string(message[0:bufLen]), nil
}
