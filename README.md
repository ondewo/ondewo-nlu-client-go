<div align="center">
  <table>
    <tr>
      <td>
        <a href="https://ondewo.com/">
            <img width="400px" src="https://raw.githubusercontent.com/ondewo/ondewo-logos/master/ondewo_we_automate_your_phone_calls.png"/>
        </a>
      </td>
    </tr>
    <tr>
       <td align="center">
          <a href="https://www.linkedin.com/company/ondewo"><img width="40px" src="https://cdn-icons-png.flaticon.com/512/3536/3536505.png"></a>
          <a href="https://www.facebook.com/ondewo"><img width="40px" src="https://cdn-icons-png.flaticon.com/512/733/733547.png"></a>
          <a href="https://twitter.com/ondewo"><img width="40px" src="https://cdn-icons-png.flaticon.com/512/733/733579.png"></a>
          <a href="https://www.instagram.com/ondewo.ai/"><img width="40px" src="https://cdn-icons-png.flaticon.com/512/174/174855.png"></a>
       </td>
    </tr>
  </table>
  <h1 align="center">
    ONDEWO NLU Client Go
  </h1>
</div>

## Overview

`ondewo-nlu-client-go` is the Go gRPC client of the ONDEWO NLU (Natural Language Understanding) API. It is a
compiled version of the [ONDEWO NLU API](https://github.com/ondewo/ondewo-nlu-api),
generated with the [ONDEWO PROTO COMPILER](https://github.com/ondewo/ondewo-proto-compiler).

ONDEWO APIs use [Protocol Buffers](https://github.com/protocolbuffers/protobuf) version 3 (proto3)
as their Interface Definition Language (IDL) to define the API interface and the structure of the
payload messages. The same interface definition is used for the gRPC version of the API in all
languages, so the service and message names below are the ones documented for the API itself.

Everything under `api/` is **generated**. It is nevertheless committed to this repository, because a
Go module is resolved straight from its version control tree — there is no build step between
`go get` and the consumer's compiler. Hand-written code therefore lives *outside* `api/`.

## Installation

The client is a plain Go module, published to the public module proxy
([`proxy.golang.org`](https://proxy.golang.org)) straight from this repository's release tags — there
is no registry account, no token and no `go install` step in between. Add it to your module with:

```shell
go get github.com/ondewo/ondewo-nlu-client-go/v7@latest   ## or @v7.2.0 to pin an exact release
```

Then import the package of the service you need:

```go
import nlupb "github.com/ondewo/ondewo-nlu-client-go/v7/api/ondewo/nlu"
```

> **The `/v7` is part of the name, not a version selector.** From major version 2 on, a Go module
> path carries its major version as a `/vN` suffix (see
> [the module reference](https://go.dev/ref/mod#major-version-suffixes)). Dropping it is not a
> shorter spelling of the same module — `go get github.com/ondewo/ondewo-nlu-client-go` names a
> *different*, unpublished module and fails with `invalid version: module contains a go.mod file, so
> major version must be compatible`. The suffix moves with the major version of the ONDEWO NLU API,
> so an `8.x` release will be imported as `/v8`, and a program can depend on both at once.

Releases are tagged twice on the same commit: with the ONDEWO release number (`7.2.0`), which is
what the [GitHub releases page](https://github.com/ondewo/ondewo-nlu-client-go/releases) lists and
what the rest of the ONDEWO client fleet uses, and with the `v`-prefixed spelling (`v7.2.0`), which
is the only tag shape Go tooling recognises as a module version. Use the `v`-prefixed one in
`go get`, `go.mod` and anywhere else a version is written.

Nothing needs to be configured for a private proxy or a credential: the module is public, so the
default `GOPROXY=https://proxy.golang.org,direct` and `GOSUMDB=sum.golang.org` resolve and verify it
as they do any other dependency. `make TEST` prints the exact module path and release tag of the
checked-out version.

To work on the client itself:

```shell
git clone https://github.com/ondewo/ondewo-nlu-client-go.git   ## Clone the repository
cd ondewo-nlu-client-go                                        ## Change into the repo directory
make setup_developer_environment_locally              ## Check out submodules, install pre-commit hooks
```

Building and releasing need only `make`, `git`, `docker` and `perl` on the host. The stubs are
generated in the proto compiler image, and every step that needs the Go toolchain or the `gh` CLI runs
in the utils image built from [`Dockerfile.utils`](Dockerfile.utils) (`ondewo-nlu-client-utils-go:<version>`),
which mounts the repository and runs as the invoking user, so nothing it writes is root-owned. The Go
targets themselves (`go_build`, `test`, `vet`, `publish_dry_run`, …) call `go` directly: that is how
they run inside the image, how CI runs its checks, and with a local Go toolchain they work without
docker as well.

## Usage

```go
package main

import (
    "context"
    "crypto/tls"
    "log"
    "os"
    "time"

    "google.golang.org/grpc"
    "google.golang.org/grpc/credentials"

    nlupb "github.com/ondewo/ondewo-nlu-client-go/v7/api/ondewo/nlu"
    "github.com/ondewo/ondewo-nlu-client-go/v7/auth"
)

func main() {
    // auth.WithBearerToken sends the Keycloak access token as `authorization: Bearer <token>` on
    // every call of this connection, exactly as the other ONDEWO clients do. It refuses to attach
    // itself to a plaintext connection, so it is paired with transport credentials here.
    conn, err := grpc.NewClient(
        "grpc-nlu.ondewo.com:443",
        grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{})),
        auth.WithBearerToken(os.Getenv("ONDEWO_NLU_ACCESS_TOKEN")),
    )
    if err != nil {
        log.Fatalf("could not connect: %v", err)
    }
    defer conn.Close()

    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()

    // Every service of the API has a generated New<Service>Client constructor. Browse
    // api/ondewo/nlu/ for the 17 services this product exposes.
    client := nlupb.NewAgentsClient(conn)

    response, err := client.ListAgents(ctx, &nlupb.ListAgentsRequest{})
    if err != nil {
        log.Fatalf("rpc failed: %v", err)
    }
    log.Printf("agents: %v", response.GetAgentsWithOwners())
}
```

## Repository structure

```
.
├── api                                    <----- GENERATED - do not edit, `make generate_ondewo_protos` rewrites it
│   └── ondewo
│       └── nlu
│           ├── *.pb.go                    <----- messages (protoc-gen-go)
│           └── *_grpc.pb.go               <----- service stubs (protoc-gen-go-grpc)
├── auth                                   <----- HAND WRITTEN - the `authorization: Bearer` credential
├── tests                                  <----- HAND WRITTEN - the go test suite (see Testing below)
├── ondewo-nlu-api                             <----- submodule @ https://github.com/ondewo/ondewo-nlu-api
├── ondewo-proto-compiler                  <----- submodule @ https://github.com/ondewo/ondewo-proto-compiler
├── .github
│   └── workflows
│       └── ci.yml                          <----- tests and lint on every push / PR - no secrets, publishes nothing
├── CONTRIBUTING.md
├── Dockerfile.utils                       <----- the go toolchain + gh CLI image every go/gh step of `make build` and `make release` runs in
├── go.mod                                 <----- module manifest, written by the compiler on the first run
├── go.sum
├── LICENSE
├── Makefile
├── README.md
└── RELEASE.md
```

## Regenerating the stubs

The generated code is produced by the `ondewo-go-proto-compiler` docker image, which is built from
the `ondewo-proto-compiler` submodule. The fixed image tag is the only contract between the two
repositories.

```shell
make update_submodules                      ## git submodule update --init --recursive
make checkout_defined_submodule_versions    ## check out the pins from the Makefile Variables chapter
make update_go_version                      ## derive go.mod's module path, the self-imports and README's versions from the Makefile
make build_compiler                         ## build ondewo-go-proto-compiler:latest from the submodule
make generate_ondewo_protos                 ## generate api/ from ondewo-nlu-api/ondewo
make go_build_via_docker                    ## build the utils image and compile the module in it
```

`make build` runs the whole chain in that order; `make check_build` then asserts that every `.proto`
produced a `*.pb.go`.

A few properties of the generation worth knowing:

* The image copies the mounted input volume into a temporary directory inside the container and
  compiles there, so the protos of `ondewo-nlu-api` are never modified.
* `api/` is **wiped** before every run, so a proto that was renamed or deleted upstream leaves no
  orphaned stub behind. Never put hand-written code below it.
* The Go import path is baked into every generated file by `protoc-gen-go`, so the module path is
  passed to the image as the third positional argument — it cannot be fixed up afterwards.
* `go.mod` and `go.sum` are written by the image **only when this repository has neither**. Once
  they are committed they are yours to maintain; the image prints the module requirements of the
  generated stubs at the end of every run so a drift is visible.
* Generation needs no network: every module the stubs are compiled against was pre-downloaded when
  the image was built.
* The compiler container runs as the invoking user (`--user`), so the generated files are never
  root-owned and no `sudo` is needed afterwards.

## Testing

```shell
make check_stubs            ## assert the generated stubs are committed
make test_via_docker        ## gofmt, go vet and test_coverage below, in the utils image (no local go needed)
make test                   ## go test over every package
make test_coverage          ## the same suite under -race, plus the hand-written coverage gate
make test_coverage_generated ## report (never gate) how much of api/ the suite exercises
```

The suite lives in `tests/` — never below `api/`, which is wiped on every regeneration — and needs
neither a network nor a running ONDEWO server. gRPC connections are made over an in-memory
`bufconn` listener, so a client stub, a server stub and a real HTTP/2 connection are exercised
in-process.

What it asserts about the **generated** code:

* a message with a scalar, a `map` of a nested message type and two well-known `Timestamp`s
  survives `proto.Marshal` → `proto.Unmarshal` unchanged, and a truncated payload is rejected;
* a proto3 `optional` scalar keeps its explicit presence — an explicitly set `0` is transmitted and
  arrives as a non-nil pointer, while an unset field stays `nil`. This is the distinction the
  angular target of the same compiler once lost, which made `false`/`0`/`""` unsendable;
* the enum zero value is the `*_UNSPECIFIED` member, with its generated name/value maps;
* the `grpc.ServiceDesc` of every service (from `protoc-gen-go-grpc`) lists exactly the RPCs the
  proto descriptor of that service (from `protoc-gen-go`) does — the two plugins run separately and
  each half compiles on its own, so a disagreement is otherwise invisible;
* every generated `New<Service>Client` binds to a connection, and every generated **unary** stub is
  actually called over the wire and has to come back as `codes.Unimplemented` — 328 of them for
  this product, which proves each one marshals its request and builds a method name the transport
  accepts;
* an RPC answered by a fake server round-trips its response, and one the server leaves to the
  generated `Unimplemented*Server` base type reports `codes.Unimplemented`;
* all 19 compiled `.proto` files are registered in the global descriptor registry as proto3.

**Coverage.** The threshold (`COVERAGE_THRESHOLD` in the `Makefile`, currently **100%**) is
enforced over the hand-written packages only — `auth/` — because everything below `api/` is machine
output: gating on it would measure how much of protoc's output a test happens to walk. The stubs
are still exercised for real, as listed above; `make test_coverage_generated` prints their figure
(**16.4%** of generated statements at the time of writing) for the record. `make test_coverage` also
fails if it ends up measuring no hand-written function at all, so a deleted package cannot turn the
gate into a green no-op.

`.github/workflows/ci.yml` runs these targets except `test`, which `test_coverage` supersedes, and
`test_via_docker`, because it installs go with `actions/setup-go` instead. It runs them on
`ubuntu-latest` against the go directive of `go.mod` and the toolchain the compiler image generates
with. It does **not** build the compiler image or check out the submodules: it builds and tests the
committed stubs, which is what a consumer of the module gets. CI is a test and lint gate only — it
uses no secret, and it never tags, packages or publishes anything.

## Release

A release runs entirely on the release host, driven by the `Makefile` — CI tags, packages and
publishes nothing, and nothing has to be approved or clicked by hand. See `make help` for the full
list of targets.

```shell
make ondewo_release                         ## credentials from the devops-accounts repo, then `make release`
```

The host needs `make`, `git` (with SSH access to GitHub and to Bitbucket, where
`ondewo-devops-accounts` lives), `docker` and `perl`; every step that needs `go` or `gh`
runs in the utils image. `make release` runs, in this order:

1. **Credentials** — `check_gh_credentials` fails when `GITHUB_GH_TOKEN` is empty or still the
   placeholder, and `validate_release_credentials_via_docker_image` then asks GitHub with one
   read-only call (`gh api repos/ondewo/ondewo-nlu-client-go`) whether the token is valid and may
   push to this repository. A revoked token stops the release here (HTTP 401 `Bad credentials`),
   before anything is pushed.
2. **Checks, build and tests** — `check_release_notes`, `build` (stubs in the compiler image,
   `go build` in the utils image), `check_build`, `check_go_module_path` and `test_via_docker`.
3. **Commit** — the generated stubs and the version-bearing files are committed as
   `PREPARING FOR RELEASE <version>`. A failing commit stops the release instead of letting the
   tags name the previous commit.
4. **Publication rehearsal** — `publish_dry_run_via_docker` on that commit (see below), the last
   step before anything is pushed.
5. **Push** — `master`, the release branch and **two** tags for the same commit: the ONDEWO release
   tag (`7.2.0`) and the `v`-prefixed tag Go tooling requires (`v7.2.0`). The `v` tag *is* the
   published module.
6. **Module proxy** — `publish_go_module_via_docker` asks `proxy.golang.org` for the new version,
   with a short retry. A failure only prints a warning: the module is released once its tag is
   pushed, and the proxy fetches the tag on the first `go get` anyway.
7. **GitHub release, last** — created from the matching `RELEASE.md` entry, so a GitHub release
   exists only for a version whose every earlier step succeeded.

`make ondewo_release` runs `make update_go_version` first (and `make build` runs it again), so a
release for which only `ONDEWO_NLU_VERSION` in the `Makefile` was changed, as ondewo-nlu-api's
`release_all_clients` does, still commits a `go.mod`, self-imports and README that agree with the tag.

### How publishing works, and what it costs

Nothing is uploaded. `proxy.golang.org` clones the tag on the first request for it and serves a zip
of exactly what git has under that tag, so **the tag *is* the published artifact** and there is no
registry account to own, no namespace to claim and no publishing credential to rotate. The whole
correctness question is therefore about the tag:

* it must be spelled `v<semver>` — `7.2.0` alone is not a Go module version;
* the module path in `go.mod` must end in `/vN` matching the tag's major version, and that path is
  baked by `protoc-gen-go` into every generated import, so it cannot be patched after generation —
  `make generate_ondewo_protos` passes it to the compiler image as the third positional argument;
* everything a consumer compiles has to be *committed*, because the proxy never sees a working tree.

Each of those is a gate rather than a convention:

```shell
make check_go_module_path   ## go.mod and every self-import carry the /vN that ONDEWO_NLU_VERSION implies
make check_release_notes    ## RELEASE.md has an entry for this version (`gh release create -n ""` would not complain)
make publish_dry_run        ## rehearse the whole publication, without credentials
```

`make publish_dry_run` is the interesting one. It packs `git archive HEAD` into the module zip the
proxy would serve, publishes it through a throwaway `file://` module proxy, and then resolves and
compiles it from a consumer module outside this tree under the real release version — so a module
path that disagrees with the tag, a self-import missing its suffix, or generated code that never
reached the commit all fail here instead of at a stranger's `go get`. It needs no secret of any
kind, and the network only for the module's dependencies, as any `go get` does. `make release` runs
it on the release commit, before anything is pushed; run it yourself (or
`make publish_dry_run_via_docker`) after committing a change to the stubs or to `go.mod`.

### Credentials

The release uses exactly one credential, and it buys the **GitHub release**, not the module:

| Variable | Where it comes from | What it is for |
| --- | --- | --- |
| `GITHUB_GH_TOKEN` | `ondewo-devops-accounts/account_github.env` | the validity check before the first push, and `gh release create` — the GitHub release page and its notes |

It lives only in `ondewo-devops-accounts`. `make ondewo_release` clones that repository, reads only
the `GITHUB_GH_TOKEN=` line of `account_github.env` and passes it into `make release`
(`clone_devops_accounts` + `run_release_with_devops`), which hands it to the utils container through
the environment (`docker run -e GITHUB_GH_TOKEN`), never on the `docker` command line. The clone is
gitignored and removed after a successful release. The working default in the `Makefile` is the
placeholder `ENTER_YOUR_TOKEN_HERE`, and a real token is never committed. No GitHub repository or
organisation secret is used anywhere: CI only tests, and needs none.

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md).

## License

Apache License 2.0 — see [LICENSE](LICENSE).
