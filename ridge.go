package ridge

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strings"
)

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

func NewRequest(event json.RawMessage) (*http.Request, error) {
	var r Request
	if err := json.Unmarshal(event, &r); err != nil {
		return nil, err
	}
	return r.HTTPRequest(), nil
}

type RequestBody struct {
	*strings.Reader
}

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
		Form:          formV,
		Host:          host,
		RequestURI:    uri,
		URL:           u,
	}
	return &req
}

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

type Response struct {
	StatusCode int               `json:"statusCode"`
	Headers    map[string]string `json:"headers"`
	Body       string            `json:"body"`
}

func NewResponseWriter() *ResponseWriter {
	return &ResponseWriter{
		statusCode: http.StatusOK,
		header:     make(http.Header),
		body:       []byte{},
	}
}

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
	for key, _ := range w.header {
		h[key] = w.header.Get(key)
	}
	return Response{
		StatusCode: w.statusCode,
		Headers:    h,
		Body:       string(w.body),
	}
}
