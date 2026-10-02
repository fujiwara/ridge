//go:build rie

// Integration tests with the AWS Lambda Runtime Interface Emulator (RIE).
// They build test/rie as a custom runtime bootstrap and run it in the
// official Lambda base image, so Docker is required.
//
//	go test -v -tags rie -run TestRIE .
//
// The image can be overridden by the RIDGE_RIE_IMAGE environment variable.
package ridge_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/fujiwara/ridge"
)

const (
	rieDefaultImage = "public.ecr.aws/lambda/provided:al2023"
	rieTermMessage  = "ridge-rie: TERM signal received" // must match test/rie
)

type rieEcho struct {
	Method     string              `json:"method"`
	RequestURI string              `json:"request_uri"`
	RemoteAddr string              `json:"remote_addr"`
	Header     http.Header         `json:"header"`
	Form       map[string][]string `json:"form"`
}

type rieContainer struct {
	id       string
	endpoint string
}

func startRIE(t *testing.T) *rieContainer {
	t.Helper()
	if _, err := exec.LookPath("docker"); err != nil {
		t.Skip("docker is not available")
	}
	image := os.Getenv("RIDGE_RIE_IMAGE")
	if image == "" {
		image = rieDefaultImage
	}

	dir := t.TempDir()
	bootstrap := filepath.Join(dir, "bootstrap")
	build := exec.Command("go", "build", "-o", bootstrap, "./test/rie")
	build.Env = append(os.Environ(), "CGO_ENABLED=0", "GOOS=linux", "GOARCH="+runtime.GOARCH)
	if out, err := build.CombinedOutput(); err != nil {
		t.Fatalf("failed to build bootstrap: %s\n%s", err, out)
	}

	out, err := exec.Command("docker", "run", "-d",
		"-p", "127.0.0.1::8080",
		"-v", bootstrap+":/var/runtime/bootstrap:ro",
		image, "bootstrap",
	).Output()
	if err != nil {
		t.Fatalf("failed to start container: %s", err)
	}
	c := &rieContainer{id: strings.TrimSpace(string(out))}
	t.Cleanup(func() {
		exec.Command("docker", "rm", "-f", c.id).Run()
	})

	out, err = exec.Command("docker", "port", c.id, "8080/tcp").Output()
	if err != nil {
		t.Fatalf("failed to get container port: %s", err)
	}
	addr := strings.TrimSpace(strings.SplitN(string(out), "\n", 2)[0])
	c.endpoint = "http://" + addr + "/2015-03-31/functions/function/invocations"

	// wait for RIE to accept connections
	deadline := time.Now().Add(30 * time.Second)
	for {
		resp, err := http.Get("http://" + addr)
		if err == nil {
			resp.Body.Close()
			break
		}
		if time.Now().After(deadline) {
			t.Fatalf("RIE did not start: %s", err)
		}
		time.Sleep(200 * time.Millisecond)
	}
	return c
}

func (c *rieContainer) invoke(t *testing.T, event []byte) *ridge.Response {
	t.Helper()
	resp, err := http.Post(c.endpoint, "application/json", bytes.NewReader(event))
	if err != nil {
		t.Fatalf("failed to invoke: %s", err)
	}
	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("unexpected invoke status %d: %s", resp.StatusCode, b)
	}
	var r ridge.Response
	if err := json.Unmarshal(b, &r); err != nil {
		t.Fatalf("failed to decode response: %s\n%s", err, b)
	}
	return &r
}

func (c *rieContainer) logs(t *testing.T) string {
	t.Helper()
	out, err := exec.Command("docker", "logs", c.id).CombinedOutput()
	if err != nil {
		t.Fatalf("failed to get logs: %s", err)
	}
	return string(out)
}

// loadEvent loads a fixture and overrides its path if path is not empty.
func loadEvent(t *testing.T, name, path string) []byte {
	t.Helper()
	b, err := os.ReadFile(filepath.Join("test", name))
	if err != nil {
		t.Fatal(err)
	}
	if path == "" {
		return b
	}
	var ev map[string]any
	if err := json.Unmarshal(b, &ev); err != nil {
		t.Fatal(err)
	}
	if ev["version"] == "2.0" {
		ev["rawPath"] = path
		ev["rawQueryString"] = ""
		delete(ev, "queryStringParameters")
	} else {
		ev["path"] = path
		delete(ev, "queryStringParameters")
		delete(ev, "multiValueQueryStringParameters")
	}
	b, err = json.Marshal(ev)
	if err != nil {
		t.Fatal(err)
	}
	return b
}

func TestRIE(t *testing.T) {
	c := startRIE(t)

	t.Run("echo", func(t *testing.T) {
		tests := []struct {
			fixture    string
			method     string
			requestURI string
			version    string
			form       map[string]string
		}{
			{"get-v1.json", "GET", "/path/to/example?foo=bar+baz&foo=boo+uoo", "1.0", map[string]string{"foo": "bar baz"}},
			{"get-rest.json", "GET", "/path/to/example?foo=bar+baz&foo=boo+uoo", "", map[string]string{"foo": "bar baz"}},
			{"get-v2.json", "GET", "/%E3%83%9E%E3%83%AB%E3%83%81%E3%83%90%E3%82%A4%E3%83%88?foo=1&foo=2&bar=3", "2.0", map[string]string{"foo": "1", "bar": "3"}},
			{"post-v2.json", "POST", "/", "2.0", map[string]string{"foo": "bar baz"}},
		}
		for _, tt := range tests {
			t.Run(tt.fixture, func(t *testing.T) {
				r := c.invoke(t, loadEvent(t, tt.fixture, ""))
				if r.StatusCode != http.StatusOK {
					t.Fatalf("unexpected status %d: %s", r.StatusCode, r.Body)
				}
				var echo rieEcho
				if err := json.Unmarshal([]byte(r.Body), &echo); err != nil {
					t.Fatalf("failed to decode body: %s\n%s", err, r.Body)
				}
				if echo.Method != tt.method {
					t.Errorf("method = %q, want %q", echo.Method, tt.method)
				}
				if echo.RequestURI != tt.requestURI {
					t.Errorf("request uri = %q, want %q", echo.RequestURI, tt.requestURI)
				}
				if echo.RemoteAddr != "203.0.113.1:0" {
					t.Errorf("remote addr = %q", echo.RemoteAddr)
				}
				if v := echo.Header.Get(ridge.PayloadVersionHeaderName); v != tt.version {
					t.Errorf("payload version = %q, want %q", v, tt.version)
				}
				for _, h := range []string{"Lambda-Runtime-Aws-Request-Id", "Lambda-Runtime-Invoked-Function-Arn"} {
					if echo.Header.Get(h) == "" {
						t.Errorf("%s header is empty", h)
					}
				}
				for k, v := range tt.form {
					if len(echo.Form[k]) == 0 || echo.Form[k][0] != v {
						t.Errorf("form %s = %v, want %q", k, echo.Form[k], v)
					}
				}
			})
		}
	})

	t.Run("cookie", func(t *testing.T) {
		want := []string{"foo=1", "bar=2"}
		for _, fixture := range []string{"get-v1.json", "get-rest.json", "get-v2.json"} {
			t.Run(fixture, func(t *testing.T) {
				r := c.invoke(t, loadEvent(t, fixture, "/cookie"))
				got := r.MultiValueHeaders["Set-Cookie"]
				if fixture == "get-rest.json" {
					// REST API has no version field, so cookies are returned only in multiValueHeaders.
					if len(r.Cookies) != 0 {
						t.Errorf("cookies = %v, want empty", r.Cookies)
					}
				} else {
					got = r.Cookies
				}
				if fmt.Sprint(got) != fmt.Sprint(want) {
					t.Errorf("cookies = %v, want %v", got, want)
				}
			})
		}
	})

	t.Run("binary", func(t *testing.T) {
		r := c.invoke(t, loadEvent(t, "get-v2.json", "/binary"))
		if !r.IsBase64Encoded {
			t.Error("isBase64Encoded = false, want true")
		}
		if r.Body != "AAEC/w==" {
			t.Errorf("body = %q", r.Body)
		}
	})

	t.Run("sigterm", func(t *testing.T) {
		if out, err := exec.Command("docker", "stop", c.id).CombinedOutput(); err != nil {
			t.Fatalf("failed to stop container: %s\n%s", err, out)
		}
		if logs := c.logs(t); !strings.Contains(logs, rieTermMessage) {
			t.Errorf("TermHandler was not called\n%s", logs)
		}
	})
}
