package main

import (
	"fmt"
	"io"
	"net"
	"os"
	"time"
)

type TelnetClient interface {
	Connect() error
	io.Closer
	Send() error
	Receive() error
}

func NewTelnetClient(address string, timeout time.Duration, in io.ReadCloser, out io.Writer) TelnetClient {
	return &client{
		address: address,
		timeout: timeout,
		in:      in,
		out:     out,
	}
}

type client struct {
	address string
	timeout time.Duration
	in      io.ReadCloser
	out     io.Writer
	conn    net.Conn
}

func (c *client) Connect() error {
	connection, err := net.DialTimeout("tcp", c.address, c.timeout)
	if err != nil {
		return fmt.Errorf("connection failed: %w", err)
	}
	fmt.Fprintf(os.Stderr, "...Connected to %s\n", c.address)
	c.conn = connection

	return nil
}

func (c *client) Close() error {
	if c.conn != nil {
		if err := c.conn.Close(); err != nil {
			return fmt.Errorf("connection closure failed: %w", err)
		}
	}
	return nil
}

func (c *client) Send() error {
	if _, err := io.Copy(c.conn, c.in); err != nil {
		return fmt.Errorf("sending failed: %w", err)
	}
	fmt.Fprintln(os.Stderr, "...EOF")
	return nil
}

func (c *client) Receive() error {
	if _, err := io.Copy(c.out, c.conn); err != nil {
		return fmt.Errorf("receiving failed: %w", err)
	}
	fmt.Fprintln(os.Stderr, "...Connection was closed by peer")
	return nil
}
