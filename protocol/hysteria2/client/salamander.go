package client

import (
	"crypto/rand"
	"net"
	"time"

	"golang.org/x/crypto/blake2b"
)

type SalamanderConnection struct {
	Connection net.PacketConn
	Key        []byte
}

func (s SalamanderConnection) ReadFrom(p []byte) (n int, addr net.Addr, err error) {
	packet := make([]byte, len(p)+8)

	n, addr, err = s.Connection.ReadFrom(packet)

	if err != nil {
		return
	}

	hash := blake2b.Sum256(append(s.Key, packet[:8]...))

	for i, c := range packet[8:n] {
		p[i] = c ^ hash[i%32]
	}

	n -= 8

	return
}

func (s SalamanderConnection) WriteTo(p []byte, addr net.Addr) (n int, err error) {
	packet := make([]byte, len(p)+8)

	rand.Read(packet[:8])

	hash := blake2b.Sum256(append(s.Key, packet[:8]...))

	for i, c := range p {
		packet[i+8] = c ^ hash[i%32]
	}

	n, err = s.Connection.WriteTo(packet, addr)

	if err == nil {
		n -= 8
	}

	return
}

func (s SalamanderConnection) Close() error {
	return s.Connection.Close()
}

func (s SalamanderConnection) LocalAddr() net.Addr {
	return s.Connection.LocalAddr()
}

func (s SalamanderConnection) SetDeadline(t time.Time) error {
	return s.Connection.SetDeadline(t)
}

func (s SalamanderConnection) SetReadDeadline(t time.Time) error {
	return s.Connection.SetReadDeadline(t)
}

func (s SalamanderConnection) SetWriteDeadline(t time.Time) error {
	return s.Connection.SetWriteDeadline(t)
}
