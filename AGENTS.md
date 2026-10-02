# AGENTS.md

Guidance for AI coding agents working in this repository.

## Project overview

ridge is a Go library that bridges AWS Lambda HTTP events (API Gateway REST API / HTTP API, ALB, Lambda Function URLs) to Go's `net/http`.
The same binary runs either as a Lambda handler or as a standalone `net/http` server, chosen at runtime from environment variables.

- Module: `github.com/fujiwara/ridge` (single package `ridge` at the repository root)
- Dependencies: `github.com/aws/aws-lambda-go`, `github.com/pires/go-proxyproto`, `github.com/google/go-cmp` (tests)

## Repository layout

- `ridge.go` — `Ridge` type, `Run` / `RunWithContext`, runtime detection (`OnLambdaRuntime`, `AsLambdaHandler`, `AsLambdaExtension`), `Response`, `ResponseWriter`, `StreamingResponseWriter`, path prefix mounting, SIGTERM handling, PROXY protocol
- `request.go` — `NewRequest` (default `RequestBuilder`), `RequestV1` / `RequestV2` payload types and conversion to `*http.Request`, `ToRequestV1` / `ToRequestV2`
- `logs.go` — `DecodeLogStream` for CloudWatch Logs subscription events
- `export_test.go` — exports unexported functions for tests in the `ridge_test` package
- `test/` — JSON event payload fixtures (`get-v1.json`, `get-v2.json`, `get-rest.json`, etc.) used by tests
- `example/` — sample application with its own `go.mod`, deployed with lambroll (not part of the main module's tests)

## Key behaviors to preserve

- **Runtime detection**: Lambda mode when `AWS_EXECUTION_ENV` starts with `AWS_Lambda` or `AWS_LAMBDA_RUNTIME_API` is set, **and** `_HANDLER` is set. Without `_HANDLER` (Lambda extension), it runs as a `net/http` server.
- **Payload version detection**: `NewRequest` reads the `version` field (`"2.0"` → `RequestV2`, `"1.0"` or empty → `RequestV1`). The global `PayloadVersion` overrides auto detection. The detected version is set to the `X-Lambda-Payload-Version` request header.
- **Version-specific responses**: `ResponseWriter.ResponseFor(version)` sets `Cookies` only when the version is non-empty (HTTP API). REST API (no `version` field) relies on `multiValueHeaders` for `Set-Cookie`.
- **Binary bodies**: responses are base64 encoded when `Content-Type` is not text (see `isTextMime` / `TextMimeTypes`) or `Content-Encoding: gzip`.
- **Streaming responses**: enabled via `Ridge.StreamingResponse` or `RIDGE_STREAMING_RESPONSE=1|true`; requires the Function URL `InvokeMode` to be `RESPONSE_STREAM`.
- **PROXY protocol**: when `ProxyProtocol` is enabled, the local server accepts connections both with and without a PROXY header (`proxyproto.USE`). go-proxyproto v0.15+ defaults to `REQUIRE`, so the policy is set explicitly.
- **RemoteAddr**: set to `sourceIp` with port `0` in the `host:port` format (`remoteAddr` in `request.go`); empty when the event has no source IP (e.g. ALB).
- Lambda mode adds `Lambda-Runtime-Aws-Request-Id` and `Lambda-Runtime-Invoked-Function-Arn` request headers.
- Keep the public API backward compatible (e.g. the aliases `Request = RequestV1` and `RequetContext` are intentionally kept).

## Development

```sh
go test -v ./...      # run all tests (same as `make test`)
go fmt ./...
go fix ./...
```

- Before committing, run `go fmt ./...` and `go fix ./...`.
- CI (`.github/workflows/go.yml`) runs `go test -v .` on multiple Go versions; keep the code compatible with the `go` directive in `go.mod`.
- Tests mostly live in the external `ridge_test` package. Add new event payloads as fixtures under `test/` and assert on the resulting `*http.Request` / `Response`.
- Add tests for any new or changed behavior, covering each payload format (v1, v2, REST API) when relevant.
- Update `README.md` when adding or changing user-facing features.

## Conventions

- Write code comments, commit messages, PR titles/descriptions, and documentation in English.
- Keep commit messages concise and focused on why the change is made.
- Do not commit directly to `main`; create a branch first.
- Releases are managed by tagpr (`.tagpr`, `vPrefix = true`); do not edit `CHANGELOG.md` by hand.
- In GitHub Actions workflows, pin `uses:` to commit SHAs (use `pinact run`) and pass `${{ ... }}` expressions to `run:` via `env:`.
