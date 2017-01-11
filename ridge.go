package ridge

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/apex/go-apex"
)

// Request represents an HTTP request received by an API Gateway proxy integrations.
type Request struct {
	Body                  string            `json:"body"`
	Headers               map[string]string `json:"headers"`
	HTTPMethod            string            `json:"httpMethod"`
	Path                  string            `json:"path"`
	PathParameters        map[string]string `json:"pathParameters"`
	QueryStringParameters map[string]string `json:"queryStringParameters"`
	Resource              string            `json:"resource"`
	StageVariables        map[string]string `json:"stageVariables"`
	RequestContext        RequestContext    `json:"requestContext"`
}

// NewRequest creates *net/http.Request from a Request.
func NewRequest(event json.RawMessage) (*http.Request, error) {
	var r Request
	if err := json.Unmarshal(event, &r); err != nil {
		return nil, err
	}
	return r.HTTPRequest(), nil
}

// RequestBody represents HTTP Request body impliments io.ReadCloser.
type RequestBody struct {
	*strings.Reader
}

// Close closes requestBody.
func (b *RequestBody) Close() error {
	return nil
}

func (r Request) HTTPRequest() *http.Request {
	header := make(http.Header)
	for key, value := range r.Headers {
		header.Add(key, value)
	}
	host := header.Get("Host")
	header.Del("Host")
	formV := make(url.Values)
	for key, value := range r.QueryStringParameters {
		formV.Add(key, value)
	}
	uri := r.Path
	if len(r.QueryStringParameters) > 0 {
		uri = uri + "?" + formV.Encode()
	}
	u, _ := url.Parse(uri)
	req := http.Request{
		Method:        r.HTTPMethod,
		Proto:         "HTTP/1.1",
		ProtoMajor:    1,
		ProtoMinor:    1,
		Header:        header,
		ContentLength: int64(len(r.Body)),
		Body:          &RequestBody{strings.NewReader(r.Body)},
		RemoteAddr:    r.RequestContext.Identity["sourceIp"],
		Host:          host,
		RequestURI:    uri,
		URL:           u,
	}
	return &req
}

// RequestContext represents request contest object.
type RequestContext struct {
	AccountID    string            `json:"accountId"`
	ApiID        string            `json:"apiId"`
	HTTPMethod   string            `json:"httpMethod"`
	Identity     map[string]string `json:"identity"`
	RequestID    string            `json:"requestId"`
	ResourceID   string            `json:"resourceId"`
	ResourcePath string            `json:"resourcePath"`
	Stage        string            `json:"stage"`
}

// Response represents a response for API Gateway proxy integration.
type Response struct {
	StatusCode int               `json:"statusCode"`
	Headers    map[string]string `json:"headers"`
	Body       string            `json:"body"`
}

// NewResponseWriter creates ResponseWriter
func NewResponseWriter() *ResponseWriter {
	return &ResponseWriter{
		statusCode: http.StatusOK,
		header:     make(http.Header),
		body:       []byte{},
	}
}

// ResponeWriter represents a response writer implements http.ResponseWriter.
type ResponseWriter struct {
	header     http.Header
	body       []byte
	statusCode int
}

func (w *ResponseWriter) Header() http.Header {
	return w.header
}

func (w *ResponseWriter) WriteHeader(code int) {
	w.statusCode = code
}

func (w *ResponseWriter) Write(b []byte) (int, error) {
	w.body = append(w.body, b...)
	return len(b), nil
}

func (w *ResponseWriter) Response() Response {
	h := make(map[string]string, len(w.header))
	for key := range w.header {
		h[key] = w.header.Get(key)
	}
	return Response{
		StatusCode: w.statusCode,
		Headers:    h,
		Body:       string(w.body),
	}
}

// Run runs http handler on Apex or net/http's server.
// If it is running on Apex (APEX_FUNCTION_NAME environment variable defined), call apex.HandleFunc().
// Otherwise start net/http server using prefix and address.
func Run(address, prefix string, mux http.Handler) {
	if os.Getenv("APEX_FUNCTION_NAME") != "" {
		apex.HandleFunc(func(event json.RawMessage, ctx *apex.Context) (interface{}, error) {
			r, err := NewRequest(event)
			if err != nil {
				log.Println(err)
				return nil, err
			}
			w := NewResponseWriter()
			mux.ServeHTTP(w, r)
			return w.Response(), nil
		})
	} else {
		m := http.NewServeMux()
		m.Handle(prefix+"/", http.StripPrefix(prefix, mux))
		log.Println("starting up with local httpd", address)
		log.Fatal(http.ListenAndServe(address, m))
	}
}
