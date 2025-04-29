# Changelog

## [v0.13.0](https://github.com/fujiwara/ridge/compare/v0.12.1...v0.13.0) - 2025-04-29
- implement for InvokeMode response_stream by @mashiike in https://github.com/fujiwara/ridge/pull/39
- adds tagpr by @fujiwara in https://github.com/fujiwara/ridge/pull/41
- Set ridge.StreamingResponse from environment variable RIDGE_STREAMING_RESPONSE by @fujiwara in https://github.com/fujiwara/ridge/pull/40
- Release for v0.13.0 by @github-actions in https://github.com/fujiwara/ridge/pull/42

## [v0.13.0](https://github.com/fujiwara/ridge/compare/v0.12.1...v0.13.0) - 2025-04-29
- implement for InvokeMode response_stream by @mashiike in https://github.com/fujiwara/ridge/pull/39
- adds tagpr by @fujiwara in https://github.com/fujiwara/ridge/pull/41
- Set ridge.StreamingResponse from environment variable RIDGE_STREAMING_RESPONSE by @fujiwara in https://github.com/fujiwara/ridge/pull/40

## [v0.12.1](https://github.com/fujiwara/ridge/compare/v0.12.0...v0.12.1) - 2024-11-18
- fix: mount at non-root path by @fujiwara in https://github.com/fujiwara/ridge/pull/38

## [v0.12.0](https://github.com/fujiwara/ridge/compare/v0.11.3...v0.12.0) - 2024-09-27
- set Lambda-Rutime headers into requests. by @fujiwara in https://github.com/fujiwara/ridge/pull/36

## [v0.11.3](https://github.com/fujiwara/ridge/compare/v0.11.2...v0.11.3) - 2024-07-05
- fix: OnLambdaRuntime() returns true even if as a extention. by @fujiwara in https://github.com/fujiwara/ridge/pull/35

## [v0.11.2](https://github.com/fujiwara/ridge/compare/v0.11.1...v0.11.2) - 2024-06-24
- ridge.Run runs net/http.Server if it is on lambda extensions. by @fujiwara in https://github.com/fujiwara/ridge/pull/32
- fix readme by @fujiwara in https://github.com/fujiwara/ridge/pull/33
- Bump actions/setup-go from 4 to 5 by @dependabot in https://github.com/fujiwara/ridge/pull/34

## [v0.11.1](https://github.com/fujiwara/ridge/compare/v0.11.0...v0.11.1) - 2024-06-20
- Goodby pkg/errors by @fujiwara in https://github.com/fujiwara/ridge/pull/31

## [v0.11.0](https://github.com/fujiwara/ridge/compare/v0.10.0...v0.11.0) - 2024-06-20
- Update examples by @ebi-yade in https://github.com/fujiwara/ridge/pull/29
- Add TermHandler by @fujiwara in https://github.com/fujiwara/ridge/pull/30

## [v0.10.0](https://github.com/fujiwara/ridge/compare/v0.9.2...v0.10.0) - 2024-06-10
- add x-amzn-requestid header from request context. by @fujiwara in https://github.com/fujiwara/ridge/pull/28

## [v0.9.2](https://github.com/fujiwara/ridge/compare/v0.9.1...v0.9.2) - 2024-06-07
- fix: duplicate headers by @fujiwara in https://github.com/fujiwara/ridge/pull/27

## [v0.9.1](https://github.com/fujiwara/ridge/compare/v0.9.0...v0.9.1) - 2024-06-07
- add ridge.Response.WriteTo() by @fujiwara in https://github.com/fujiwara/ridge/pull/26

## [v0.9.0](https://github.com/fujiwara/ridge/compare/v0.8.0...v0.9.0) - 2024-03-14
- Customizable Request Builder by @mashiike in https://github.com/fujiwara/ridge/pull/23
- not a http payload should return an error. by @fujiwara in https://github.com/fujiwara/ridge/pull/24

## [v0.8.0](https://github.com/fujiwara/ridge/compare/v0.7.0...v0.8.0) - 2024-03-01
- add ToRequestV1 and ToRequestV2 by @fujiwara in https://github.com/fujiwara/ridge/pull/22

## [v0.7.0](https://github.com/fujiwara/ridge/compare/v0.6.2...v0.7.0) - 2024-02-06
- Set "Cookies" into a response object. by @fujiwara in https://github.com/fujiwara/ridge/pull/20
- Deprecation io/ioutil by @fujiwara in https://github.com/fujiwara/ridge/pull/21

## [v0.6.2](https://github.com/fujiwara/ridge/compare/v0.6.1...v0.6.2) - 2023-10-17
- Bump github.com/pires/go-proxyproto from 0.6.0 to 0.6.1 by @dependabot in https://github.com/fujiwara/ridge/pull/16

## [v0.6.1](https://github.com/fujiwara/ridge/compare/v0.6.0...v0.6.1) - 2021-08-13
- update versions. by @fujiwara in https://github.com/fujiwara/ridge/pull/15

## [v0.6.0](https://github.com/fujiwara/ridge/compare/v0.5.0...v0.6.0) - 2020-11-13
- add RunWithContext by @fujiwara in https://github.com/fujiwara/ridge/pull/13

## [v0.5.0](https://github.com/fujiwara/ridge/compare/v0.4.1...v0.5.0) - 2020-10-28
- Drop apex support. by @fujiwara in https://github.com/fujiwara/ridge/pull/12

## [v0.4.1](https://github.com/fujiwara/ridge/compare/v0.4.0...v0.4.1) - 2020-10-28
- Enables to run on custom runtime by @fujiwara in https://github.com/fujiwara/ridge/pull/11

## [v0.4.0](https://github.com/fujiwara/ridge/compare/v0.3.0...v0.4.0) - 2020-06-19

## [v0.3.0](https://github.com/fujiwara/ridge/compare/v0.2.1...v0.3.0) - 2020-06-19
- Support to PROXY protocol by @fujiwara in https://github.com/fujiwara/ridge/pull/10

## [v0.2.1](https://github.com/fujiwara/ridge/compare/v0.2.0...v0.2.1) - 2020-04-21
- Set default content type by @fujiwara in https://github.com/fujiwara/ridge/pull/9

## [v0.2.0](https://github.com/fujiwara/ridge/compare/v0.1.0...v0.2.0) - 2020-03-18
- Support Lambda Payload v2 by @fujiwara in https://github.com/fujiwara/ridge/pull/8
