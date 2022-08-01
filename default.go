package pool

import (
	"context"
	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
	"io"
	"net"
)

func WebsocketDialer(dialAddr string) DialFunc {
	return func(ctx context.Context) (net.Conn, error) {
		c, _, _, err := ws.Dial(ctx, dialAddr)
		return c, err
	}
}

func TcpDialer(addr string) DialFunc {
	return func(ctx context.Context) (net.Conn, error) {
		return net.Dial("tcp", addr)
	}
}

func WsWriter(data []byte) func(writer io.Writer) error {
	return func(w io.Writer) error {
		writer := wsutil.NewWriter(w, ws.StateClientSide, ws.OpText)
		if _, err := writer.Write(data); err != nil {
			return err
		}
		return writer.Flush()
	}
}

func WebsocketReadFunc(fn func(p []byte)) ReadFunc {
	return func(conn net.Conn) error {
		h, r, err := wsutil.NextReader(conn, ws.StateClientSide)
		if err != nil {
			return err
		}
		if h.OpCode.IsControl() {
			return wsutil.ControlFrameHandler(conn, ws.StateClientSide)(h, r)
		}

		data := make([]byte, h.Length)
		if _, err = r.Read(data); err != nil && err != io.EOF {
			return err
		}

		fn(data)
		return nil
	}
}

func TcpReadFunc(fn func(p []byte)) ReadFunc {
	return func(conn net.Conn) error {
		data := make([]byte, 1024)
		n, err := conn.Read(data)
		if err != nil && err != io.EOF {
			return err
		}
		fn(data[:n])
		return nil
	}
}
