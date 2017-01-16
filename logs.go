package ridge

import (
	"bytes"
	"compress/gzip"
	"encoding/json"

	"github.com/pkg/errors"
)

type Message struct {
	Awslogs struct {
		Data []byte `json:"data"`
	} `json:"awslogs"`
}

type LogStream struct {
	MessageType         string     `json:"messageType"`
	Owner               string     `json:"owner"`
	LogGroup            string     `json:"logGroup"`
	LogStream           string     `json:"logStream"`
	SubscriptionFilters []string   `json:"subscriptionFilters"`
	LogEvents           []LogEvent `json:"logEvents"`
}

type LogEvent struct {
	ID        string `json:"id"`
	Timestamp int64  `json:"timestamp"`
	Message   string `json:"message"`
}

// DecodeLogStream decodes CloudwatchLogs stream passed to a lambda function.
func DecodeLogStream(event json.RawMessage) (*LogStream, error) {
	msg := Message{}
	err := json.Unmarshal(event, &msg)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode event")
	}
	gz, err := gzip.NewReader(bytes.NewReader(msg.Awslogs.Data))
	if err != nil {
		return nil, errors.Wrap(err, "cloud not create gzip reader")
	}
	dec := json.NewDecoder(gz)
	ls := LogStream{}
	err = dec.Decode(&ls)
	if err != nil {
		return nil, errors.Wrap(err, "cloud not decode log stream")
	}
	return &ls, nil
}
