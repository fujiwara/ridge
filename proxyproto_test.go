package ridge_test

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
	"testing"

	"github.com/fujiwara/ridge"
)

func TestProxyProtocol(t *testing.T) {
	r := ridge.New("127.0.0.1:0", "/", http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		host, _, _ := net.SplitHostPort(req.RemoteAddr)
		fmt.Fprint(w, host)
	}))
	r.ProxyProtocol = true
	listener, err := r.Listen()
	if err != nil {
		t.Fatal(err)
	}
	srv := &http.Server{Handler: r.Mux}
	go srv.Serve(listener)
	t.Cleanup(func() { srv.Close() })
	addr := listener.Addr().String()

	tests := []struct {
		name   string
		header string
		want   string
	}{
		{"without PROXY header", "", "127.0.0.1"},
		{"with PROXY v1 header", "PROXY TCP4 192.0.2.1 127.0.0.1 12345 80\r\n", "192.0.2.1"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			conn, err := net.Dial("tcp", addr)
			if err != nil {
				t.Fatal(err)
			}
			defer conn.Close()
			if _, err := io.WriteString(conn, tt.header+"GET / HTTP/1.0\r\nHost: localhost\r\n\r\n"); err != nil {
				t.Fatal(err)
			}
			resp, err := http.ReadResponse(bufio.NewReader(conn), nil)
			if err != nil {
				t.Fatal(err)
			}
			defer resp.Body.Close()
			body, _ := io.ReadAll(resp.Body)
			if got := strings.TrimSpace(string(body)); got != tt.want {
				t.Errorf("remote host = %q, want %q", got, tt.want)
			}
		})
	}
}
