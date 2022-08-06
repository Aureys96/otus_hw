package main

import (
	"bytes"
	"io/ioutil"
	"net"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestTelnetClient(t *testing.T) {
	t.Run("basic", func(t *testing.T) {
		l, err := net.Listen("tcp", "127.0.0.1:")
		require.NoError(t, err)
		defer func() { require.NoError(t, l.Close()) }()

		var wg sync.WaitGroup
		wg.Add(2)

		go func() {
			defer wg.Done()

			in := &bytes.Buffer{}
			out := &bytes.Buffer{}

			timeout, err := time.ParseDuration("10s")
			require.NoError(t, err)

			client := NewTelnetClient(l.Addr().String(), timeout, ioutil.NopCloser(in), out)
			require.NoError(t, client.Connect())
			defer func() { require.NoError(t, client.Close()) }()

			in.WriteString("hello\n")
			err = client.Send()
			require.NoError(t, err)

			err = client.Receive()
			require.NoError(t, err)
			require.Equal(t, "world\n", out.String())
		}()

		go func() {
			defer wg.Done()

			conn, err := l.Accept()
			require.NoError(t, err)
			require.NotNil(t, conn)
			defer func() { require.NoError(t, conn.Close()) }()

			request := make([]byte, 1024)
			n, err := conn.Read(request)
			require.NoError(t, err)
			require.Equal(t, "hello\n", string(request)[:n])

			n, err = conn.Write([]byte("world\n"))
			require.NoError(t, err)
			require.NotEqual(t, 0, n)
		}()

		wg.Wait()
	})

	t.Run("error when try to connect to unavailable server", func(t *testing.T) {
		in := &bytes.Buffer{}
		out := &bytes.Buffer{}
		client := NewTelnetClient("127.0.0.1:19646", 10*time.Second, ioutil.NopCloser(in), out)
		require.EqualError(t, client.Connect(),
			"connection failed: dial tcp 127.0.0.1:19646: connectex: No connection could be made because the target machine actively refused it.")
	})
	t.Run("error when server is already stopped", func(t *testing.T) {
		l, err := net.Listen("tcp", "127.0.0.1:")
		require.NoError(t, err)

		in := &bytes.Buffer{}
		out := &bytes.Buffer{}
		client := NewTelnetClient(l.Addr().String(), 10*time.Second, ioutil.NopCloser(in), out)
		require.NoError(t, client.Connect())

		require.NoError(t, l.Close())

		in.WriteString("test\n")
		require.Error(t, client.Send())
	})
}
