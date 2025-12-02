package client

import (
	"crypto/rand"
	"net"
	"time"

	"golang.org/x/crypto/blake2b"
)

type Salamander struct {
	Connection net.PacketConn
	Key        []byte
}

func (s Salamander) ReadFrom(p []byte) (n int, addr net.Addr, err error) {
	packet := make([]byte, len(p))

	n, addr, err = s.Connection.ReadFrom(packet)

	if err != nil {
		return
	}

	hash := blake2b.Sum256(append(s.Key, packet[:8]...))

	for i, c := range packet[8:] {
		p[i] = packet[c] ^ hash[i%32]
	}

	return
}

func (s Salamander) WriteTo(p []byte, addr net.Addr) (n int, err error) {
	packet := make([]byte, len(p)+8)

	_, err = rand.Read(packet[:8])

	if err != nil {
		return 0, err
	}

	hash := blake2b.Sum256(append(s.Key, packet[:8]...))

	for i, c := range packet[8:] {
		packet[i] = p[c] ^ hash[i%32]
	}

	return s.Connection.WriteTo(packet, addr)
}

func (s Salamander) Close() error {
	return s.Connection.Close()
}

func (s Salamander) LocalAddr() net.Addr {
	return s.Connection.LocalAddr()
}

func (s Salamander) SetDeadline(t time.Time) error {
	return s.Connection.SetDeadline(t)
}

func (s Salamander) SetReadDeadline(t time.Time) error {
	return s.Connection.SetReadDeadline(t)
}

func (s Salamander) SetWriteDeadline(t time.Time) error {
	return s.Connection.SetWriteDeadline(t)
}
