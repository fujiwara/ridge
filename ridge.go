package ridge

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"io"
	"log"
	"mime"
	"net"
	"net/http"
	"os"
	"strings"
	"sync"

	"github.com/aws/aws-lambda-go/lambda"
	proxyproto "github.com/pires/go-proxyproto"
)

// ProxyProtocol is a flag to support PROXY Protocol
var ProxyProtocol bool

// TextMimeTypes is a list of identified as text.
var TextMimeTypes = []string{"image/svg+xml", "application/json", "application/xml"}

// DefaultContentType is a default content-type when missing in response.
var DefaultContentType = "text/plain; charset=utf-8"

// Response represents a response for API Gateway proxy integration.
type Response struct {
	StatusCode        int               `json:"statusCode"`
	Headers           map[string]string `json:"headers"`
	MultiValueHeaders http.Header       `json:"multiValueHeaders"`
	Cookies           []string          `json:"cookies"`
	Body              string            `json:"body"`
	IsBase64Encoded   bool              `json:"isBase64Encoded"`
}

// WriteTo writes response to http.ResponseWriter.
func (r *Response) WriteTo(w http.ResponseWriter) (int64, error) {
	for k, vs := range r.MultiValueHeaders {
		for _, v := range vs {
			w.Header().Add(k, v)
		}
	}
	for _, c := range r.Cookies {
		w.Header().Add("Set-Cookie", c)
	}
	w.WriteHeader(r.StatusCode)
	if r.IsBase64Encoded {
		dec := base64.NewDecoder(base64.StdEncoding, strings.NewReader(r.Body))
		return io.Copy(w, dec)
	}
	n, err := io.WriteString(w, r.Body)
	return int64(n), err
}

// NewResponseWriter creates ResponseWriter
func NewResponseWriter() *ResponseWriter {
	w := &ResponseWriter{
		Buffer:     bytes.Buffer{},
		statusCode: http.StatusOK,
		header:     make(http.Header),
	}
	return w
}

// ResponseWriter represents a response writer implements http.ResponseWriter.
type ResponseWriter struct {
	bytes.Buffer
	header     http.Header
	statusCode int
}

func (w *ResponseWriter) Header() http.Header {
	return w.header
}

func (w *ResponseWriter) WriteHeader(code int) {
	w.statusCode = code
}

func (w *ResponseWriter) Response() Response {
	body := w.String()
	isBase64Encoded := false

	if t := w.header.Get("Content-Type"); t == "" {
		w.header.Set("Content-Type", DefaultContentType)
	}
	h := make(map[string]string, len(w.header))
	for key := range w.header {
		v := w.header.Get(key)
		if isBinary(key, v) {
			isBase64Encoded = true
		}
		h[key] = v
	}
	if isBase64Encoded {
		body = base64.StdEncoding.EncodeToString(w.Bytes())
	}

	return Response{
		StatusCode:        w.statusCode,
		Headers:           h,
		MultiValueHeaders: w.header,
		Cookies:           w.header.Values("Set-Cookie"),
		Body:              body,
		IsBase64Encoded:   isBase64Encoded,
	}
}

func isBinary(k, v string) bool {
	if k == "Content-Type" && !isTextMime(v) {
		return true
	}
	if k == "Content-Encoding" && v == "gzip" {
		return true
	}
	return false
}

func isTextMime(kind string) bool {
	mt, _, err := mime.ParseMediaType(kind)
	if err != nil {
		return false
	}

	if strings.HasPrefix(mt, "text/") {
		return true
	}

	isText := false
	for _, tmt := range TextMimeTypes {
		if mt == tmt {
			isText = true
			break
		}
	}
	return isText
}

// Run runs http handler on AWS Lambda runtime or net/http's server.
func Run(address, prefix string, mux http.Handler) {
	RunWithContext(context.Background(), address, prefix, mux)
}

// RunWithContext runs http handler on AWS Lambda runtime or net/http's server with context.
func RunWithContext(ctx context.Context, address, prefix string, mux http.Handler) {
	r := New(address, prefix, mux)
	r.RunWithContext(ctx)
}

// Ridge is a struct to run http handler on AWS Lambda runtime or net/http's server.
type Ridge struct {
	Address        string
	Prefix         string
	Mux            http.Handler
	RequestBuilder func(json.RawMessage) (*http.Request, error)
}

// New creates a new Ridge.
func New(address, prefix string, mux http.Handler) *Ridge {
	return &Ridge{
		Address:        address,
		Prefix:         prefix,
		Mux:            mux,
		RequestBuilder: NewRequest,
	}
}

// Run runs http handler on AWS Lambda runtime or net/http's server.
func (r *Ridge) Run() {
	r.RunWithContext(context.Background())
}

// RunWithContext runs http handler on AWS Lambda runtime or net/http's server with context.
func (r *Ridge) RunWithContext(ctx context.Context) {
	if strings.HasPrefix(os.Getenv("AWS_EXECUTION_ENV"), "AWS_Lambda") || os.Getenv("AWS_LAMBDA_RUNTIME_API") != "" {
		// go1.x or custom runtime(provided, provided.al2)
		handler := func(event json.RawMessage) (interface{}, error) {
			req, err := r.RequestBuilder(event)
			if err != nil {
				log.Println(err)
				return nil, err
			}
			w := NewResponseWriter()
			r.Mux.ServeHTTP(w, req)
			return w.Response(), nil
		}
		lambda.StartWithOptions(handler, lambda.WithContext(ctx))
	} else {
		m := http.NewServeMux()
		switch {
		case r.Prefix == "/", r.Prefix == "":
			m.Handle("/", r.Mux)
		case !strings.HasSuffix(r.Prefix, "/"):
			m.Handle(r.Prefix+"/", http.StripPrefix(r.Prefix, r.Mux))
		default:
			m.Handle(r.Prefix, http.StripPrefix(strings.TrimSuffix(r.Prefix, "/"), r.Mux))
		}
		log.Println("starting up with local httpd", r.Address)
		listener, err := net.Listen("tcp", r.Address)
		if err != nil {
			log.Fatalf("couldn't listen to %s: %s", r.Address, err.Error())
		}
		if ProxyProtocol {
			log.Println("enables to PROXY protocol")
			listener = &proxyproto.Listener{Listener: listener}
		}
		srv := http.Server{Handler: m}
		var wg sync.WaitGroup
		wg.Add(2)
		go func() {
			defer wg.Done()
			<-ctx.Done()
			log.Println("shutting down local httpd", r.Address)
			srv.Shutdown(ctx)
		}()
		if err := srv.Serve(listener); err != nil {
			if err != http.ErrServerClosed {
				log.Fatal(err)
			}
			wg.Done()
		}
		wg.Wait()
	}
}
