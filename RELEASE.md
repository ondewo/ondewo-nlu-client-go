# Release History

*****************

## Release ONDEWO NLU Go Client 7.3.1

### New Features

* `client.NewChannel(cfg client.Config, opts ...grpc.DialOption)` opens the gRPC connection for
  plaintext, TLS or mutual TLS, following the TLS contract of the python SDKs
  (`ondewo-client-utils` 4.1.1), so one set of certificates works with every ONDEWO client.
  `client.Config` carries `Host`, `Port`, `Insecure`, an optional `Logger` and the PEM **content**
  (never a file path) of the CA (`GrpcCert`, empty = system roots) and of an optional client
  identity (`GrpcClientCert` / `GrpcClientKey`).
* Refused with an error before gRPC sees them: half a client identity, a certificate and key that
  do not form a pair, a `GrpcCert` that holds no PEM certificate (typically a path), and
  `Insecure: true` combined with a client identity. No error message contains a PEM, a key or the
  whole `Config`; they name the field and `host:port`.
* A plaintext connection logs a warning naming `host:port` through `log/slog` (`Config.Logger`, or
  `slog.Default()`); the package never configures logging. Every `fmt` verb and `slog` rendering of
  a `Config` redacts the client key.
* Bare IPv6 literal hosts are bracketed (`::1` -> `[::1]:<port>`); bracketed hosts and hosts with a
  scheme are used as they are. PEMs with CRLF line endings work.
* Connection defaults of the python SDKs that grpc-go exposes: keepalive pings every 30 s with a
  20 s timeout, only while an RPC is active; 5 s maximum reconnect backoff; 2^31-1 byte maximum
  message size in both directions. Every `grpc.DialOption` passed in is applied after them and wins.
  Documented gaps: grpc-go has no `http2.max_pings_without_data` and only one keepalive timeout, and
  no per-method retry policy is configured (only gRPC's transparent retries apply).

### Improvements

* `tests/tls_test.go` runs real TLS and mutual-TLS handshakes against an in-process server with a
  PKI generated per run; the 100% coverage gate now spans `auth/` and `client/`.
* `tests/release_notes_test.go` pins the `RELEASE.md` slice the GitHub release body is built from:
  the Makefile's perl range, the spelling of every release heading, the `*****` separator closing
  every section, and non-empty notes for the current version.
* README: new section "TLS, mutual TLS and certificates" (modes, connection defaults, a test PKI with
  openssl, TLS security notes, troubleshooting).
*****************

## Release ONDEWO NLU Go Client 7.3.0

### Improvements

* Tracking API Version [7.3.0](https://github.com/ondewo/ondewo-nlu-api/releases/tag/7.3.0) ( [Documentation](https://ondewo.github.io/ondewo-nlu-api/) )

*****************

## Release ONDEWO NLU Go Client 7.2.0

### Improvements

* Tracking API Version [7.2.0](https://github.com/ondewo/ondewo-nlu-api/releases/tag/7.2.0) ( [Documentation](https://ondewo.github.io/ondewo-nlu-api/) )

*****************

## Release ONDEWO NLU Go Client 7.1.0

### Bug Fixes

* The module is declared as `github.com/ondewo/ondewo-nlu-client-go/v7`. From major version 2 on a
  Go module path has to carry its major version as a `/vN` suffix, so the previous, un-suffixed
  path made a `v7.1.0` tag unresolvable — `go get` reported `invalid version: module contains a
  go.mod file, so major version must be compatible`. The path is baked into every generated import
  by `protoc-gen-go`, so the stubs were regenerated with the corrected module path rather than
  patched in `go.mod` alone.

### Improvements

* `make check_go_module_path` fails the build when the `/vN` suffix of the module path stops
  agreeing with `ONDEWO_NLU_VERSION`, in CI and before every release. That mismatch compiles and
  tests green — it only surfaces once the tag is pushed and a consumer cannot resolve it.
* `make publish_dry_run` rehearses the whole publication without a single credential: it packs
  `git archive HEAD` into the module zip `proxy.golang.org` would serve, publishes it through a
  throwaway `file://` module proxy and resolves it from an outside consumer module under its real
  version. It runs on every push in CI and once more inside `make release`, right before the tags
  are created.
* `.github/workflows/release.yml` builds, verifies and publishes on a `v*` tag push, and refuses to
  start when the release credentials are absent instead of publishing an empty release.
* `make check_release_notes` refuses a release whose version has no entry in this file — `gh release
  create -n ""` publishes an empty release without complaining.

*****************

## Release ONDEWO NLU Go Client 0.1.0

### New Features

* Initial release of the ONDEWO NLU (Natural Language Understanding) gRPC client for Go. The module
  ships the stubs generated from the [ONDEWO NLU API](https://github.com/ondewo/ondewo-nlu-api)
  by version 5.15.1 of the
  [ONDEWO Proto Compiler](https://github.com/ondewo/ondewo-proto-compiler): one `*.pb.go` of
  messages and one `*_grpc.pb.go` of service stubs per `.proto`, below `api/ondewo/nlu/`,
  compiled against the `google.golang.org/protobuf` and `google.golang.org/grpc` runtimes pinned by
  the compiler image.
* `make build` reproduces the whole client from the two submodules — proto compiler image, stub
  generation and `go build` — and `make check_build` asserts that every `.proto` of the API
  produced a stub.

### Improvements

* The release targets tag each release twice on the same commit: with the ONDEWO release number
  (`0.1.0`) that the rest of the fleet uses, and with the `v`-prefixed spelling (`v0.1.0`) that is
  the only tag shape the Go module resolver accepts. `make publish_go_module` then warms
  `proxy.golang.org` so the new version is immediately installable with `go get`.

*****************
