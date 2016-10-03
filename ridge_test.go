package ridge_test

import (
	"encoding/json"
	"net/url"
	"testing"

	"github.com/fujiwara/ridge"
)

var getEvent = `{
  "resource": "/{proxy+}",
  "path": "/path/to/example",
  "httpMethod": "GET",
  "headers": {
    "Accept": "*/*",
    "CloudFront-Forwarded-Proto": "https",
    "CloudFront-Is-Desktop-Viewer": "true",
    "CloudFront-Is-Mobile-Viewer": "false",
    "CloudFront-Is-SmartTV-Viewer": "false",
    "CloudFront-Is-Tablet-Viewer": "false",
    "CloudFront-Viewer-Country": "JP",
    "Host": "abcdefg.execute-api.ap-northeast-1.example.com",
    "User-Agent": "curl/7.43.0",
    "Via": "1.1 a3fed41c60e2fab219a274640e58ebe5.cloudfront.net (CloudFront)",
    "X-Amz-Cf-Id": "pxfb9Y0Mr-g0QNTExIE5_IAwk4yUy1ObJ7k7xsuJv_IMNsnYjQ0kvA==",
    "X-Forwarded-For": "203.0.113.1, 54.239.196.62",
    "X-Forwarded-Port": "443",
    "X-Forwarded-Proto": "https"
  },
  "queryStringParameters": {
    "foo": "bar baz"
  },
  "pathParameters": {
    "proxy": "path/to/example"
  },
  "stageVariables": null,
  "requestContext": {
    "accountId": "123456789123",
    "resourceId": "yotsyh",
    "stage": "prod",
    "requestId": "817175b9-890f-11e6-960e-4f321627a748",
    "identity": {
      "cognitoIdentityPoolId": null,
      "accountId": null,
      "cognitoIdentityId": null,
      "caller": null,
      "apiKey": null,
      "sourceIp": "203.0.113.1",
      "cognitoAuthenticationType": null,
      "cognitoAuthenticationProvider": null,
      "userArn": null,
      "userAgent": "curl/7.43.0",
      "user": null
    },
    "resourcePath": "/{proxy+}",
    "httpMethod": "GET",
    "apiId": "n1d78val4e"
  },
  "body": null
}
`

var postEvent = `{
  "resource": "/{proxy+}",
  "path": "/path/to/example",
  "httpMethod": "POST",
  "headers": {
    "Accept": "*/*",
    "CloudFront-Forwarded-Proto": "https",
    "CloudFront-Is-Desktop-Viewer": "true",
    "CloudFront-Is-Mobile-Viewer": "false",
    "CloudFront-Is-SmartTV-Viewer": "false",
    "CloudFront-Is-Tablet-Viewer": "false",
    "CloudFront-Viewer-Country": "JP",
    "Content-Type": "application/x-www-form-urlencoded",
    "Host": "abcdefg.execute-api.ap-northeast-1.example.com",
    "User-Agent": "curl/7.43.0",
    "Via": "1.1 736a82fbf158fe646f468bd5664ef95c.cloudfront.net (CloudFront)",
    "X-Amz-Cf-Id": "R9xVKMTNQUBmddfXFPdox98chlQAzzGB6mw7hxZa5aBLeEQZfRdeKw==",
    "X-Forwarded-For": "203.0.113.1, 54.239.196.51",
    "X-Forwarded-Port": "443",
    "X-Forwarded-Proto": "https"
  },
  "queryStringParameters": null,
  "pathParameters": {
    "proxy": "path/to/example"
  },
  "stageVariables": null,
  "requestContext": {
    "accountId": "123456789123",
    "resourceId": "yotsyh",
    "stage": "prod",
    "requestId": "8eed9b4f-890f-11e6-9f3c-1584342606cd",
    "identity": {
      "cognitoIdentityPoolId": null,
      "accountId": null,
      "cognitoIdentityId": null,
      "caller": null,
      "apiKey": null,
      "sourceIp": "203.0.113.1",
      "cognitoAuthenticationType": null,
      "cognitoAuthenticationProvider": null,
      "userArn": null,
      "userAgent": "curl/7.43.0",
      "user": null
    },
    "resourcePath": "/{proxy+}",
    "httpMethod": "POST",
    "apiId": "n1d78val4e"
  },
  "body": "foo=bar%20baz"
}
`

func TestGetRequest(t *testing.T) {
	r, err := ridge.NewRequest(json.RawMessage(getEvent))
	if err != nil {
		t.Fatalf("failed to decode getEvent: %s", err)
	}
	if r.Host != "abcdefg.execute-api.ap-northeast-1.example.com" {
		t.Errorf("Host: %s is not expected", r.Host)
	}
	if r.Method != "GET" {
		t.Errorf("Method: %s is not expected", r.Method)
	}
	u, _ := url.Parse("/path/to/example?foo=bar+baz")
	if r.URL.String() != u.String() {
		t.Errorf("URL: %s is not expected", r.URL)
	}
	if v := r.FormValue("foo"); v != "bar baz" {
		t.Errorf("FormValue(foo): %s is not expected", v)
	}
	if v := r.Header.Get("CloudFront-Viewer-Country"); v != "JP" {
		t.Errorf("Header[CloudFront-Viewer-Country]: %s is not expected", v)
	}
	if v := r.Header.Get("Via"); v != "1.1 a3fed41c60e2fab219a274640e58ebe5.cloudfront.net (CloudFront)" {
		t.Errorf("Header[Via]: %s is not expected", v)
	}
	if r.RemoteAddr != "203.0.113.1" {
		t.Errorf("RemoteAddr: %s is not expected", r.RemoteAddr)
	}
}

func TestPostRequest(t *testing.T) {
	r, err := ridge.NewRequest(json.RawMessage(postEvent))
	if err != nil {
		t.Fatalf("failed to decode postEvent: %s", err)
	}
	r.ParseForm()

	if r.Host != "abcdefg.execute-api.ap-northeast-1.example.com" {
		t.Errorf("Host: %s is not expected", r.Host)
	}
	if r.Method != "POST" {
		t.Errorf("Method: %s is not expected", r.Method)
	}
	u, _ := url.Parse("/path/to/example")
	if r.URL.String() != u.String() {
		t.Errorf("URL: %s is not expected", r.URL)
	}
	if v := r.PostFormValue("foo"); v != "bar baz" {
		t.Errorf("PostFormValue(foo): %s is not expected", v)
	}
	if v := r.Header.Get("CloudFront-Viewer-Country"); v != "JP" {
		t.Errorf("Header[CloudFront-Viewer-Country]: %s is not expected", v)
	}
	if v := r.Header.Get("Via"); v != "1.1 736a82fbf158fe646f468bd5664ef95c.cloudfront.net (CloudFront)" {
		t.Errorf("Header[Via]: %s is not expected", v)
	}
	if r.RemoteAddr != "203.0.113.1" {
		t.Errorf("RemoteAddr: %s is not expected", r.RemoteAddr)
	}
	if r.ContentLength != 13 {
		t.Errorf("Content-Length: %d is not expected", r.ContentLength)
	}
}
