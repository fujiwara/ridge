package ridge

import "net"

func (r *Ridge) SetStreamingResponse() {
	r.setStreamingResponse()
}

func (r *Ridge) Listen() (net.Listener, error) {
	return r.listen()
}
