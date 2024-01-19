// Package client implements the [REST API] exposed by [Paperless-ngx].
// Paperless-ngx is a document management system transforming physical
// documents into a searchable online archive.
//
// # Authentication
//
// Multiple [authentication schemes] are supported:
//
//   - [UsernamePasswordAuth]: HTTP basic authentication.
//   - [TokenAuth]: Paperless-ngx API authentication tokens.
//   - [GCPServiceAccountKeyAuth]: OpenID Connect (OIDC) using a Google Cloud
//     Platform service account.
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
