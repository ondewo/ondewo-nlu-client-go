# Release History

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
  by version 5.15.0 of the
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
