package ridge

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"log"
	"mime"
	"net/http"
	"os"
	"strings"

	"github.com/apex/go-apex"
	"github.com/aws/aws-lambda-go/lambda"
)

// TextMimeTypes is a list of identified as text.
var TextMimeTypes = []string{"image/svg+xml", "application/json", "application/xml"}

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
