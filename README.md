# gobi

Gobi is a declarative build tool for go.

## Targets

### Lint

Lint is a golang target that uses gometalinter to run lint. It runs gofmt to
check for golang format compliance. The default configuration runs gofmt, govet
and golint on all golang sources in the project ignoring the vendor directory.
