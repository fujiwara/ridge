package ridge

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"mime"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/apex/go-apex"
	"github.com/aws/aws-lambda-go/lambda"
)

// TextMimeTypes is a list of identified as text.
var TextMimeTypes = []string{"image/svg+xml", "application/json", "application/xml"}

// Request represents an HTTP request received by an API Gateway proxy integrations.
type Request struct {
	Body                            string              `json:"body"`
	Headers                         map[string]string   `json:"headers"`
	MultiValueHeaders               http.Header         `json:"multiValueHeaders"`
	HTTPMethod                      string              `json:"httpMethod"`
	Path                            string              `json:"path"`
	PathParameters                  map[string]string   `json:"pathParameters"`
	QueryStringParameters           map[string]string   `json:"queryStringParameters"`
	MultiValueQueryStringParameters map[string][]string `json:"multiValueQueryStringParameters"`
	Resource                        string              `json:"resource"`
	StageVariables                  map[string]string   `json:"stageVariables"`
	RequestContext                  RequestContext      `json:"requestContext"`
	IsBase64Encoded                 bool                `json:"isBase64Encoded"`
}

// NewRequest creates *net/http.Request from a Request.
func NewRequest(event json.RawMessage) (*http.Request, error) {
	var r Request
	if err := json.Unmarshal(event, &r); err != nil {
		return nil, err
	}
	return r.httpRequest()
}

func (r Request) httpRequest() (*http.Request, error) {
	header := make(http.Header)
	if len(r.MultiValueHeaders) > 0 {
		for key, values := range r.MultiValueHeaders {
			for _, value := range values {
				header.Add(key, value)
			}
		}
	} else {
		for key, value := range r.Headers {
			header.Add(key, value)
		}
	}
	host := header.Get("Host")
	header.Del("Host")
	v := make(url.Values)
	if len(r.MultiValueQueryStringParameters) > 0 {
		for key, values := range r.MultiValueQueryStringParameters {
			for _, value := range values {
				v.Add(key, value)
			}
		}
	} else {
		for key, value := range r.QueryStringParameters {
			v.Add(key, value)
		}
	}
	uri := r.Path
	if len(r.QueryStringParameters) > 0 {
		uri = uri + "?" + v.Encode()
	}
	u, _ := url.Parse(uri)
	var contentLength int64
	var b io.Reader
	if r.IsBase64Encoded {
		raw := make([]byte, len(r.Body))
		n, err := base64.StdEncoding.Decode(raw, []byte(r.Body))
		if err != nil {
			return nil, err
		}
		contentLength = int64(n)
		b = bytes.NewReader(raw[0:n])
	} else {
		contentLength = int64(len(r.Body))
		b = strings.NewReader(r.Body)
	}
	req := http.Request{
		Method:        r.HTTPMethod,
		Proto:         "HTTP/1.1",
		ProtoMajor:    1,
		ProtoMinor:    1,
		Header:        header,
		ContentLength: contentLength,
		Body:          ioutil.NopCloser(b),
		RemoteAddr:    r.RequestContext.Identity["sourceIp"],
		Host:          host,
		RequestURI:    uri,
		URL:           u,
	}
	return &req, nil
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
	StatusCode        int               `json:"statusCode"`
	Headers           map[string]string `json:"headers"`
	MultiValueHeaders http.Header       `json:"multiValueHeaders"`
	Body              string            `json:"body"`
	IsBase64Encoded   bool              `json:"isBase64Encoded"`
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

// ResponeWriter represents a response writer implements http.ResponseWriter.
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

	h := make(map[string]string, len(w.header))
	for key := range w.header {
		v := w.header.Get(key)
		if isBinary(key, v) {
			isBase64Encoded = true
			body = base64.StdEncoding.EncodeToString(w.Bytes())
		}
		h[key] = v
	}
	return Response{
		StatusCode:        w.statusCode,
		Headers:           h,
		MultiValueHeaders: w.header,
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

// Run runs http handler on Apex(nodejs runtime), go runtime, or net/http's server.
// If it is running on Apex (APEX_FUNCTION_NAME environment variable defined), call apex.HandleFunc().
// Otherwise start net/http server using prefix and address.
func Run(address, prefix string, mux http.Handler) {
	if env := os.Getenv("AWS_EXECUTION_ENV"); env != "" {
		handler := func(event json.RawMessage) (interface{}, error) {
			r, err := NewRequest(event)
			if err != nil {
				log.Println(err)
				return nil, err
			}
			w := NewResponseWriter()
			mux.ServeHTTP(w, r)
			return w.Response(), nil
		}
		if strings.HasPrefix(env, "AWS_Lambda_nodejs") && os.Getenv("APEX_FUNCTION_NAME") != "" {
			// Apex (node runtime)
			apex.HandleFunc(
				func(event json.RawMessage, ctx *apex.Context) (interface{}, error) {
					// redirect stdout to stderr in Apex functions
					stdout := os.Stdout
					os.Stdout = os.Stderr
					defer func() {
						os.Stdout = stdout
					}()
					return handler(event)
				},
			)
		} else if strings.HasPrefix(env, "AWS_Lambda_go") {
			// native Go runtime
			lambda.Start(handler)
		} else {
			log.Printf("Environment %s is not supported", env)
		}
	} else {
		m := http.NewServeMux()
		switch {
		case prefix == "/", prefix == "":
			m.Handle("/", mux)
		case !strings.HasSuffix(prefix, "/"):
			m.Handle(prefix+"/", http.StripPrefix(prefix, mux))
		default:
			m.Handle(prefix, http.StripPrefix(strings.TrimSuffix(prefix, "/"), mux))
		}
		log.Println("starting up with local httpd", address)
		log.Fatal(http.ListenAndServe(address, m))
	}
}
