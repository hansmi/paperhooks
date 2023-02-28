// Package client implements the the [REST API] exposed by [Paperless-ngx].
// Paperless-ngx is a document management system transforming physical
// documents into a searchable online archive.
//
// # Authentication
//
// Paperless-ngx supports multiple [authentication schemes].
// [UsernamePasswordAuth] implements HTTP basic authentication and [TokenAuth],
// well, tokens.
//
// # Pagination
//
// APIs returning lists of items support pagination (e.g.
// [Client.ListDocuments]). The [ListOptions] struct embedded in the
// API-specific options supports specifying the page to request. Pagination
// tokens are available via the [Response] struct.
//
// [REST API]: https://docs.paperless-ngx.com/api/
// [Paperless-ngx]: https://docs.paperless-ngx.com/
// [authentication schemes]: https://docs.paperless-ngx.com/api/#authorization
package client
