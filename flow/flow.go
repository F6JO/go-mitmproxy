package flow

import (
	"encoding/json"
	"net/http"
	"net/url"

	uuid "github.com/satori/go.uuid"
	_log "github.com/sirupsen/logrus"
)

var log = _log.WithField("at", "flow")

type Request struct {
	Method string
	URL    *url.URL
	Proto  string
	Header http.Header
	Body   []byte

	raw *http.Request
}

func (req *Request) MarshalJSON() ([]byte, error) {
	r := make(map[string]interface{})
	r["method"] = req.Method
	r["url"] = req.URL.String()
	r["proto"] = req.Proto
	r["header"] = req.Header
	r["body"] = req.Body
	return json.Marshal(r)
}

func NewRequest(req *http.Request) *Request {
	return &Request{
		Method: req.Method,
		URL:    req.URL,
		Proto:  req.Proto,
		Header: req.Header,
		raw:    req,
	}
}

func (r *Request) Raw() *http.Request {
	return r.raw
}

type Response struct {
	StatusCode int         `json:"statusCode"`
	Header     http.Header `json:"header"`
	Body       []byte      `json:"-"`

	decodedBody []byte
	decoded     bool // decoded reports whether the response was sent compressed but was decoded to decodedBody.
	decodedErr  error
}

type Flow struct {
	*Request
	*Response

	// https://docs.mitmproxy.org/stable/overview-features/#streaming
	// 如果为 true，则不缓冲 Request.Body 和 Response.Body，且不进入之后的 Addon.Request 和 Addon.Response
	Stream bool
	done   chan struct{}

	Id    uuid.UUID
	State map[string]interface{} // Can add value by addon
}

func (f *Flow) MarshalJSON() ([]byte, error) {
	j := make(map[string]interface{})
	j["id"] = f.Id
	j["request"] = f.Request
	j["response"] = f.Response
	return json.Marshal(j)
}

func NewFlow() *Flow {
	return &Flow{
		done:  make(chan struct{}),
		Id:    uuid.NewV4(),
		State: make(map[string]interface{}),
	}
}

func (f *Flow) Done() <-chan struct{} {
	return f.done
}

func (f *Flow) Finish() {
	close(f.done)
}
