package ridge_test

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"testing"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/fujiwara/ridge"
)

func TestDecodeLogStream(t *testing.T) {
	b, err := ioutil.ReadFile("test/logstream.json")
	if err != nil {
		t.Error(err)
		return
	}
	ls, err := ridge.DecodeLogStream(json.RawMessage(b))
	if err != nil {
		t.Error(err)
		return
	}
	if ls.Owner != "999999999999" {
		t.Errorf("owner %s", ls.Owner)
	}
	if ls.MessageType != "DATA_MESSAGE" {
		t.Errorf("messageType %s", ls.MessageType)
	}
	if ls.LogGroup != "/aws/lambda/ridge-test_main" {
		t.Errorf("logGroup %s", ls.LogGroup)
	}
	if ls.LogStream != "2017/01/16/[$LATEST]d2eb3fdfb4814aefbce6a55c3261be60" {
		t.Errorf("logStream %s", ls.LogStream)
	}
	if len(ls.SubscriptionFilters) < 1 || ls.SubscriptionFilters[0] != "LambdaStream_ridge-test_log" {
		t.Errorf("subscriptionFilters %v", ls.SubscriptionFilters)
	}
	if len(ls.LogEvents) != 2 {
		t.Errorf("len(logEvents) != 3 %v", ls.LogEvents)
	}
}

func ExampleDecodeLogStream() {
	lambda.Start(func(event json.RawMessage) (interface{}, error) {
		logStream, err := ridge.DecodeLogStream(event)
		if err != nil {
			return nil, err
		}
		for _, e := range logStream.LogEvents {
			//
			log.Println(e.Message)
		}
		return "", nil
	})
}
