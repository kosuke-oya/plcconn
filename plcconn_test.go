package plcconn

import (
	"bytes"
	"errors"
	"net"
	"testing"
	"time"
)

var (
	ErrConnectionRefused = errors.New("connection refused")
	ErrTimeout           = errors.New("timeout")
)

func TestConnect(t *testing.T) {
	t.Parallel()
	// Create a mock TCP server
	// ip : 127.0.0.1, port : 任意
	mockServer, err := net.Listen("tcp", "127.0.0.1:")
	if err != nil {
		t.Fatalf("Failed to start mock server: %v", err)
	}
	port := mockServer.Addr().(*net.TCPAddr).Port
	defer mockServer.Close()

	// Start a goroutine to accept connections
	go func() {
		for {
			conn, err := mockServer.Accept()
			if err != nil {
				return
			}
			// Close the connection immediately
			conn.Close()
		}
	}()

	// time.Sleep(1 * time.Second)
	time.Sleep(1 * time.Second)

	// Create a TcpClient instance
	clientOk := NewPlcConn("localhost", port, 1)
	defer clientOk.Close()

	// Test connecting to the server
	err = clientOk.Connect()
	if err != nil {
		t.Errorf("Connect returned an error: %v", err)
	}

	// Test connecting to a non-existent server
	clientNg := NewPlcConn("localhost", port-1, 1)
	defer clientNg.Close()
	err = clientNg.Connect()
	if err == nil {
		t.Error("Connect did not return an error")
	}

}
func TestWrite(t *testing.T) {
	t.Parallel()
	// Create a mock TCP server
	mockServer, err := net.Listen("tcp", "127.0.0.1:")
	if err != nil {
		t.Fatalf("Failed to start mock server: %v", err)
	}
	port := mockServer.Addr().(*net.TCPAddr).Port
	defer mockServer.Close()

	// Start a goroutine to accept connections
	go func() {
		for {
			conn, err := mockServer.Accept()
			if err != nil {
				return
			}
			// Read the message from the connection
			buf := make([]byte, RESBUF_MAX_RLEN)
			_, err = conn.Read(buf)
			if err != nil {
				return
			}
			// Write a response back to the client
			_, err = conn.Write([]byte("response"))
			if err != nil {
				return
			}
			// Close the connection
			conn.Close()
		}
	}()

	// time.Sleep(1 * time.Second)
	time.Sleep(1 * time.Second)

	// Create a TcpClient instance
	client := NewPlcConn("localhost", port, 1)
	defer client.Close()

	// Test writing a message and receiving a response
	msg := []byte("hello")
	err = client.Connect()
	if err != nil {
		t.Fatalf("Connect returned an error: %v", err)
	}
	resp, err := client.Write(msg)
	if err != nil {
		t.Errorf("Write returned an error: %v", err)
	}
	// responseという先頭の8バイトを取得
	respMod := resp[:8]

	expectedResp := []byte("response")
	if !bytes.Equal(respMod, expectedResp) {
		t.Errorf("Write returned an unexpected response: got %s, want %s", resp, expectedResp)
	}

}
func TestIsConnected(t *testing.T) {
	t.Parallel()
	// Create a mock TCP server
	mockServer, err := net.Listen("tcp", "127.0.0.1:")
	if err != nil {
		t.Fatalf("Failed to start mock server: %v", err)
	}
	port := mockServer.Addr().(*net.TCPAddr).Port
	defer mockServer.Close()

	// Start a goroutine to accept connections
	go func() {
		for {
			conn, err := mockServer.Accept()
			if err != nil {
				return
			}
			// Read the message from the connection
			buf := make([]byte, RESBUF_MAX_RLEN)
			_, err = conn.Read(buf)
			if err != nil {
				return
			}
			// Write a response back to the client
			_, err = conn.Write([]byte("response"))
			if err != nil {
				return
			}
			// Close the connection
			conn.Close()
		}
	}()

	time.Sleep(1 * time.Second)

	client := NewPlcConn("localhost", port, 1)
	defer client.Close()

	// Test when the client is not connected
	if client.IsConnected() {
		t.Error("IsConnected returned true, expected false")
	}

	// Connect here
	err = client.Connect()
	if err != nil {
		t.Fatalf("Connect returned an error: %v", err)
	}
	// Check if the client is connected
	if !client.IsConnected() {
		t.Error("IsConnected returned false, expected true")
	}
}
func TestOpenWriteClose(t *testing.T) {
	t.Parallel()
	// Create a mock TCP server
	mockServer, err := net.Listen("tcp", "127.0.0.1:")
	if err != nil {
		t.Fatalf("Failed to start mock server: %v", err)
	}
	port := mockServer.Addr().(*net.TCPAddr).Port
	defer mockServer.Close()

	// Start a goroutine to accept connections
	go func() {
		for {
			conn, err := mockServer.Accept()
			if err != nil {
				return
			}
			// Read the message from the connection
			buf := make([]byte, RESBUF_MAX_RLEN)
			_, err = conn.Read(buf)
			if err != nil {
				return
			}
			// Write a response back to the client
			_, err = conn.Write([]byte("response"))
			if err != nil {
				return
			}
			// Close the connection
			conn.Close()
		}
	}()

	time.Sleep(1 * time.Second)

	client := NewPlcConn("localhost", port, 1)
	defer client.Close()

	// Test opening, writing, and closing the connection
	msg := []byte("hello")
	resp, err := client.OpenWriteClose(msg)
	if err != nil {
		t.Errorf("OpenWriteClose returned an error: %v", err)
	}
	// responseという先頭の8バイトを取得
	respMod := resp[:8]

	expectedResp := []byte("response")
	if !bytes.Equal(respMod, expectedResp) {
		t.Errorf("OpenWriteClose returned an unexpected response: got %s, want %s", resp, expectedResp)
	}
}
