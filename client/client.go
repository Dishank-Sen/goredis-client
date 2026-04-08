package client

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"net"
	"strings"
	"time"
)

type Config struct{
	addr string
}

type Client struct{
	cfg Config
	conn net.Conn
	reader *bufio.Reader
}

func NewClient(cfg Config) *Client{
	if strings.TrimSpace(cfg.addr) == ""{
		panic("redis: new client nil config")
	}
	conn, err := net.Dial("tcp", cfg.addr)
	if err != nil{
		panic(err)
	}
	return &Client{
		cfg: cfg,
		conn: conn,
		reader: bufio.NewReader(conn),
	}
}

func (c *Client) Close(){
	c.conn.Close()
}

func (c *Client) readRESP() (interface{}, error) {
	b, err := c.reader.Peek(1)
	if err != nil {
		return nil, err
	}

	switch b[0] {

	case '+': // Simple String
		line, err := c.reader.ReadString('\n')
		if err != nil {
			return nil, err
		}
		return strings.TrimSpace(line[1:]), nil

	case '-': // Error
		line, err := c.reader.ReadString('\n')
		if err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("%s", strings.TrimSpace(line[1:]))

	case '$': // Bulk String
		line, err := c.reader.ReadString('\n')
		if err != nil {
			return nil, err
		}

		var length int
		fmt.Sscanf(line, "$%d\r\n", &length)

		if length == -1 {
			return nil, nil // nil value
		}

		buf := make([]byte, length)
		_, err = io.ReadFull(c.reader, buf)
		if err != nil {
			return nil, err
		}

		// consume \r\n
		_, err = c.reader.ReadString('\n')
		if err != nil {
			return nil, err
		}

		return string(buf), nil

	default:
		return nil, fmt.Errorf("unknown RESP type: %q", b[0])
	}
}

func (c *Client) Set(ctx context.Context, key, value string) error {
	cmd := "SET"
	cmdStr := fmt.Sprintf("*3\r\n$%d\r\n%s\r\n$%d\r\n%s\r\n$%d\r\n%s\r\n",
		len(cmd), cmd, len(key), key, len(value), value)

	if err := c.writeWithContext(ctx, c.conn, []byte(cmdStr)); err != nil {
		return err
	}

	resp, err := c.readRESP()
	if err != nil {
		return err
	}

	// Expect "OK"
	if str, ok := resp.(string); ok && str == "OK" {
		return nil
	}

	return fmt.Errorf("unexpected response: %v", resp)
}

func (c *Client) Get(ctx context.Context, key string) (string, error) {
	cmd := "GET"
	cmdStr := fmt.Sprintf("*2\r\n$%d\r\n%s\r\n$%d\r\n%s\r\n",
		len(cmd), cmd, len(key), key)

	if err := c.writeWithContext(ctx, c.conn, []byte(cmdStr)); err != nil {
		return "", err
	}

	resp, err := c.readRESP()
	if err != nil {
		return "", err
	}

	if resp == nil {
		return "", fmt.Errorf("key not found")
	}

	if str, ok := resp.(string); ok {
		return str, nil
	}

	return "", fmt.Errorf("unexpected response type: %T", resp)
}

func (c *Client) writeWithContext(ctx context.Context, conn net.Conn, data []byte) error {
	errCh := make(chan error, 1)

	go func() {
		_, err := conn.Write(data)
		errCh <- err
	}()

	select {
	case <-ctx.Done():
		// cancel effect → force unblock
		conn.SetWriteDeadline(time.Now()) // immediate timeout
		return ctx.Err()

	case err := <-errCh:
		return err
	}
}