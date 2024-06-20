package ridge

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
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
		return nil, fmt.Errorf("could not decode event: %w", err)
	}
	gz, err := gzip.NewReader(bytes.NewReader(msg.Awslogs.Data))
	if err != nil {
		return nil, fmt.Errorf("cloud not create gzip reader: %w", err)
	}
	dec := json.NewDecoder(gz)
	ls := LogStream{}
	err = dec.Decode(&ls)
	if err != nil {
		return nil, fmt.Errorf("cloud not decode log stream: %w", err)
	}
	return &ls, nil
}
