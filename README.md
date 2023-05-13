# Client and hook toolkit for Paperless-ngx

[![Latest release](https://img.shields.io/github/v/release/hansmi/paperhooks)][releases]
[![CI workflow](https://github.com/hansmi/paperhooks/actions/workflows/ci.yaml/badge.svg)](https://github.com/hansmi/paperhooks/actions/workflows/ci.yaml)
[![Go reference](https://pkg.go.dev/badge/github.com/hansmi/paperhooks.svg)](https://pkg.go.dev/github.com/hansmi/paperhooks)

Paperhooks is a toolkit for [writing consumption hooks][paperless-hooks] for
Paperless-ngx written using the Go programming language. A
[REST API][paperless-api] client is part of the toolkit
([`pkg/client`](./pkg/client/)).

[Paperless-ngx][paperless] is a document management system transforming
physical documents into a searchable online archive.

## Run integration tests

[Integration tests](https://en.wikipedia.org/wiki/Integration_testing) execute
operations against a real Paperless-ngx server running in a Docker container.
The wrapper script enables _destructive_ tests and should not be run against
a production instance.

Commands:

```shell
env --chdir contrib/integration-env docker-compose up

contrib/run-integration
```

[paperless-api]: https://docs.paperless-ngx.com/api/
[paperless-hooks]: https://docs.paperless-ngx.com/advanced_usage/#consume-hooks
[paperless]: https://docs.paperless-ngx.com/
[releases]: https://github.com/hansmi/paperhooks/releases/latest

<!-- vim: set sw=2 sts=2 et : -->
