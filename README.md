# replify

**replify** is a Go library designed to simplify and standardize API response wrapping for RESTful services. It leverages the Decorator Pattern to dynamically add error handling, metadata, pagination, and other response features in a clean and human-readable format. With Wrapify, you can ensure consistent and extensible API responses with minimal boilerplate. Perfect for building robust, maintainable REST APIs in Go!

### Requirements

- Go version 1.23 or higher

### Installation

To install, you can use the following commands based on your preference:

- For a specific version:

  ```bash
  go get github.com/sivaosorg/replify@v0.0.1
  ```

- For the latest version:
  ```bash
  go get -u github.com/sivaosorg/replify@latest
  ```

### Getting started

#### Getting replify

With [Go's module support](https://go.dev/wiki/Modules#how-to-use-modules), `go [build|run|test]` automatically fetches the necessary dependencies when you add the import in your code:

```go
import "github.com/sivaosorg/replify"
```

#### Usage

> An example of the wrapper-based standardized API response

```json
{
  "data": "response body here", // The primary data payload of the response.
  "debug": {
    "___abc": "trace sessions_id: 4919e84fc26881e9fe790f5d07465db4",
    "refer": 1234
  }, // Debugging information (useful for development).
  "message": "How are you? I'm good", // A message providing additional context about the response.
  "meta": {
    "api_version": "v0.0.1", // API version used for the request.
    "custom_fields": {
      "fields": "userID: 103"
    }, // Additional custom metadata fields.
    "locale": "en_US", // Locale used for the request, e.g., "en-US".
    "request_id": "80eafc6a1655ec5a06595d155f1e6951", // Unique identifier for the request, useful for debugging.
    "requested_time": "2024-12-14T20:24:23.983839+07:00" // Timestamp when the request was made.
  }, // Metadata about the API response.
  "pagination": {
    "is_last": true, // Indicates whether this is the last page.
    "page": 1000, // Current page number.
    "per_page": 2, // Number of items per page.
    "total_items": 120, // Total number of items available.
    "total_pages": 34 // Total number of pages.
  }, // Pagination details, if applicable.
  "path": "/api/v1/users", // Request path for which the response is generated.
  "status_code": 200, // HTTP status code for the response.
  "total": 1 // Total number of items (used in non-paginated responses).
}
```

> Structure of the wrapper-based standardized API response

```go
// R represents a wrapper around the main `wrapper` struct.
// It is used as a high-level abstraction to provide a simplified interface for handling API responses.
// The `R` type allows for easier manipulation of the wrapped data, metadata,
// and other response components, while maintaining the flexibility of the underlying `wrapper` structure.
type R struct {
    *wrapper
}
// Available checks whether the `wrapper` instance is non-nil.
func (w *wrapper) Available() bool
// Body retrieves the body data associated with the `wrapper` instance.
func (w *wrapper) Body() interface{}
// Cause traverses the error chain and returns the underlying cause of the error
// associated with the `wrapper` instance.
func (w *wrapper) Cause() error
// Debugging retrieves the debugging information from the `wrapper` instance.
func (w *wrapper) Debugging() map[string]interface{}
// Error retrieves the error associated with the `wrapper` instance.
func (w *wrapper) Error() string
// Header retrieves the `header` associated with the `wrapper` instance.
func (w *wrapper) Header() *header
// IsBodyPresent checks whether the body data is present in the `wrapper` instance.
func (w *wrapper) IsBodyPresent() bool
// IsClientError checks whether the HTTP status code indicates a client error.
func (w *wrapper) IsClientError() bool
// IsDebuggingKeyPresent checks whether a specific key exists in the `debug` information.
func (w *wrapper) IsDebuggingKeyPresent(key string) bool
// IsDebuggingPresent checks whether debugging information is present in the `wrapper` instance.
func (w *wrapper) IsDebuggingPresent() bool
// IsError checks whether there is an error present in the `wrapper` instance.
// This function returns `true` if the `wrapper` contains an error, which can be any of the following:
//   - An error present in the `errors` field.
//   - A client error (4xx status code) or a server error (5xx status code).
func (w *wrapper) IsError() bool
// IsErrorPresent checks whether an error is present in the `wrapper` instance.
func (w *wrapper) IsErrorPresent() bool
// IsHeaderPresent checks whether header information is present in the `wrapper` instance.
func (w *wrapper) IsHeaderPresent() bool
// IsLastPage checks whether the current page is the last page of results.
func (w *wrapper) IsLastPage() bool
// IsMetaPresent checks whether metadata information is present in the `wrapper` instance.
func (w *wrapper) IsMetaPresent() bool
// IsPagingPresent checks whether pagination information is present in the `wrapper` instance.
func (w *wrapper) IsPagingPresent() bool
// IsRedirection checks whether the HTTP status code indicates a redirection response.
func (w *wrapper) IsRedirection() bool
// IsServerError checks whether the HTTP status code indicates a server error.
func (w *wrapper) IsServerError() bool
// IsStatusCodePresent checks whether a valid status code is present in the `wrapper` instance.
func (w *wrapper) IsStatusCodePresent() bool
// IsSuccess checks whether the HTTP status code indicates a successful response.
func (w *wrapper) IsSuccess() bool
// IsTotalPresent checks whether the total number of items is present in the `wrapper` instance.
func (w *wrapper) IsTotalPresent() bool
// Json serializes the `wrapper` instance into a compact JSON string.
func (w *wrapper) Json() string
// JsonPretty serializes the `wrapper` instance into a prettified JSON string.
func (w *wrapper) JsonPretty() string
// Message retrieves the message associated with the `wrapper` instance.
func (w *wrapper) Message() string
// Meta retrieves the `meta` information from the `wrapper` instance.
func (w *wrapper) Meta() *meta
// OnKeyDebugging retrieves the value of a specific debugging key from the `wrapper` instance.
func (w *wrapper) OnKeyDebugging(key string) interface{}
// Pagination retrieves the `pagination` instance associated with the `wrapper`.
func (w *wrapper) Pagination() *pagination
// R represents a wrapper around the main `wrapper` struct. It is used as a high-level
// abstraction to provide a simplified interface for handling API responses.
// The `R` type allows for easier manipulation of the wrapped data, metadata, and other
// response components, while maintaining the flexibility of the underlying `wrapper` structure.
func (w *wrapper) Reply() R
// Respond generates a map representation of the `wrapper` instance.
func (w *wrapper) Respond() map[string]interface{}
// StatusCode retrieves the HTTP status code associated with the `wrapper` instance.
func (w *wrapper) StatusCode() int
// StatusText returns a human-readable string representation of the HTTP status.
func (w *wrapper) StatusText() string
// Total retrieves the total number of items associated with the `wrapper` instance.
func (w *wrapper) Total() int
// WithApiVersion sets the API version in the `meta` field of the `wrapper` instance.
func (w *wrapper) WithApiVersion(v string) *wrapper
// WithApiVersionf sets the API version in the `meta` field of the `wrapper` instance using a formatted string.
func (w *wrapper) WithApiVersionf(format string, args ...interface{}) *wrapper
// WithBody sets the body data for the `wrapper` instance.
func (w *wrapper) WithBody(v interface{}) *wrapper
// WithCustomFieldKV sets a specific custom field key-value pair in the `meta` field of the `wrapper` instance.
func (w *wrapper) WithCustomFieldKV(key string, value interface{}) *wrapper
// WithCustomFieldKVf sets a specific custom field key-value pair in the `meta` field of the `wrapper` instance
// using a formatted value.
func (w *wrapper) WithCustomFieldKVf(key string, format string, args ...interface{}) *wrapper
// WithCustomFields sets the custom fields in the `meta` field of the `wrapper` instance.
func (w *wrapper) WithCustomFields(values map[string]interface{}) *wrapper
// WithDebugging sets the debugging information for the `wrapper` instance.
func (w *wrapper) WithDebugging(v map[string]interface{}) *wrapper
// WithDebuggingKV adds a key-value pair to the debugging information in the `wrapper` instance.
func (w *wrapper) WithDebuggingKV(key string, value interface{}) *wrapper
// WithDebuggingKVf adds a formatted key-value pair to the debugging information in the `wrapper` instance.
func (w *wrapper) WithDebuggingKVf(key string, format string, args ...interface{}) *wrapper
// WithErrMessage adds a plain contextual message to an existing error and sets it for the `wrapper` instance.
func (w *wrapper) WithErrMessage(err error, message string) *wrapper
// WithErrMessagef adds a formatted contextual message to an existing error and sets it for the `wrapper` instance.
func (w *wrapper) WithErrMessagef(err error, format string, args ...interface{}) *wrapper
// WithErrorAck sets an error with a stack trace for the `wrapper` instance.
func (w *wrapper) WithErrorAck(err error) *wrapper
// WithErrWrap wraps an existing error with an additional message and sets it for the `wrapper` instance.
func (w *wrapper) WithErrWrap(err error, message string) *wrapper
// WithErrWrapf wraps an existing error with a formatted message and sets it for the `wrapper` instance.
func (w *wrapper) WithErrWrapf(err error, format string, args ...interface{}) *wrapper
// WithError sets an error for the `wrapper` instance using a plain error message.
func (w *wrapper) WithError(message string) *wrapper
// WithErrorf sets a formatted error for the `wrapper` instance.
func (w *wrapper) WithErrorf(format string, args ...interface{}) *wrapper
// WithHeader sets the header for the `wrapper` instance.
func (w *wrapper) WithHeader(v *header) *wrapper
// WithIsLast sets whether the current page is the last one in the wrapper's pagination.
func (w *wrapper) WithIsLast(v bool) *wrapper
// WithLocale sets the locale in the `meta` field of the `wrapper` instance.
func (w *wrapper) WithLocale(v string) *wrapper
// WithMessage sets a message for the `wrapper` instance.
func (w *wrapper) WithMessage(message string) *wrapper
// WithMessagef sets a formatted message for the `wrapper` instance.
func (w *wrapper) WithMessagef(message string, args ...interface{}) *wrapper
// WithMeta sets the metadata for the `wrapper` instance.
func (w *wrapper) WithMeta(v *meta) *wrapper
// WithPage sets the current page number in the wrapper's pagination.
func (w *wrapper) WithPage(v int) *wrapper
// WithPagination sets the pagination information for the `wrapper` instance.
func (w *wrapper) WithPagination(v *pagination) *wrapper
// WithPath sets the request path for the `wrapper` instance.
func (w *wrapper) WithPath(v string) *wrapper
// WithPathf sets a formatted request path for the `wrapper` instance.
func (w *wrapper) WithPathf(v string, args ...interface{}) *wrapper
// WithPerPage sets the number of items per page in the wrapper's pagination.
func (w *wrapper) WithPerPage(v int) *wrapper
// WithRequestID sets the request ID in the `meta` field of the `wrapper` instance.
func (w *wrapper) WithRequestID(v string) *wrapper
// WithRequestIDf sets the request ID in the `meta` field of the `wrapper` instance using a formatted string.
func (w *wrapper) WithRequestIDf(format string, args ...interface{}) *wrapper
// WithRequestedTime sets the requested time in the `meta` field of the `wrapper` instance.
func (w *wrapper) WithRequestedTime(v time.Time) *wrapper
// WithStatusCode sets the HTTP status code for the `wrapper` instance.
func (w *wrapper) WithStatusCode(code int) *wrapper
// WithTotal sets the total number of items for the `wrapper` instance.
func (w *wrapper) WithTotal(total int) *wrapper
// WithTotalItems sets the total number of items in the wrapper's pagination.
func (w *wrapper) WithTotalItems(v int) *wrapper
// WithTotalPages sets the total number of pages in the wrapper's pagination.
func (w *wrapper) WithTotalPages(v int) *wrapper
```

> Parse JSON string to wrapper-based standardized API response

```go
package main

import (
	"fmt"
	"log"
	"time"

	"github.com/sivaosorg/replify"
)

func main() {
	jsonStr := `{
		"data": "response body here",
		"debug": {
		  "___abc": "trace sessions_id: 4919e84fc26881e9fe790f5d07465db4",
		  "refer": 1234
		},
		"message": "How do you do? I'm good",
		"meta": {
		  "api_version": "v0.0.1",
		  "custom_fields": {
			"fields": "userID: 103"
		  },
		  "locale": "en_US",
		  "request_id": "80eafc6a1655ec5a06595d155f1e6951",
		  "requested_time": "2024-12-14T20:24:23.983839+07:00"
		},
		"pagination": {
		  "is_last": true,
		  "page": 1000,
		  "per_page": 2,
		  "total_items": 120,
		  "total_pages": 34
		},
		"path": "/api/v1/users",
		"status_code": 200,
		"total": 1
	  }`
	t := time.Now()
	w, err := replify.UnwrapJSON(jsonStr)
	diff := time.Since(t)
	if err != nil {
		log.Fatalf("Error parsing JSON: %v", err)
	}
	fmt.Printf("Exe time: %+v\n", diff.String())
	fmt.Printf("%+v\n", w.OnKeyDebugging("___abc"))
	fmt.Printf("%+v\n", w.JsonPretty())
}
```

> Example of creating a wrapper-based standardized API response

```go
package main

import (
	"fmt"

	"github.com/sivaosorg/replify"
	"github.com/sivaosorg/replify/randn"
)

func main() {
	p := replify.Pages().
		WithIsLast(true).
		WithPage(1).
		WithTotalItems(120).
		WithPage(1000).
		WithTotalPages(34).
		WithPerPage(2)
	w := replify.New().
		WithStatusCode(200).
		WithTotal(1).
		WithMessagef("How are you? %v", "I'm good").
		WithDebuggingKV("refer", 1234).
		WithDebuggingKVf("___abc", "trace sessions_id: %v", randn.CryptoID()).
		WithBody("response body here").
		WithPath("/api/v1/users").
		WithCustomFieldKVf("fields", "userID: %v", 103).
		WithPagination(p)
	if !w.Available() {
		return
	}
	fmt.Println(w.Json())
	fmt.Println(w.StatusCode())
	fmt.Println(w.StatusText())
	fmt.Println(w.Message())
	fmt.Println(w.Body())
	fmt.Println(w.IsSuccess())
	fmt.Println(w.Respond())
	fmt.Println(w.Meta().IsCustomFieldPresent())
	fmt.Println(w.Meta().IsApiVersionPresent())
	fmt.Println(w.Meta().IsRequestIDPresent())
	fmt.Println(w.Meta().IsRequestedTimePresent())
}
```

#### HTTP Status Codes for APIs

| **Scenario**                       | **HTTP Status Codes**                                     | **Example**                                                                         |
| ---------------------------------- | --------------------------------------------------------- | ----------------------------------------------------------------------------------- |
| **Successful Resource Retrieval**  | 200 OK, 304 Not Modified                                  | `GET /users/123` - Returns user data or indicates cached content is valid.          |
| **Resource Creation**              | 201 Created                                               | `POST /users` - Creates a new user, with a location header for the resource URL.    |
| **Asynchronous Processing**        | 202 Accepted                                              | `POST /large-file` - File upload starts, and the client polls for result.           |
| **Validation Errors**              | 400 Bad Request                                           | `POST /users` - Missing `name` field or invalid input format.                       |
| **Authentication Issues**          | 401 Unauthorized, 403 Forbidden                           | Accessing a secured endpoint without valid credentials or insufficient permissions. |
| **Rate Limiting**                  | 429 Too Many Requests                                     | Exceeded allowed API requests within a time window.                                 |
| **Missing Resource**               | 404 Not Found                                             | `GET /users/999` - User ID not found.                                               |
| **Server Failures**                | 500 Internal Server Error, 503 Service Unavailable        | Database failure, unhandled exception, or server in maintenance mode.               |
| **Version Conflicts**              | 409 Conflict                                              | `PUT /resource/123` with an outdated version, causing a conflict.                   |
| **Partial Responses**              | 206 Partial Content                                       | Video streaming or fetching paginated results.                                      |
| **Redirecting Resources**          | 302 Found, 307 Temporary Redirect, 308 Permanent Redirect | URL redirection during resource movement or temporary relocation.                   |
| **Client Data Too Large**          | 413 Payload Too Large, 414 URI Too Long                   | Request body or URL exceeds server-defined limits.                                  |
| **Unsupported Content Type**       | 415 Unsupported Media Type                                | Client uploads a file format not accepted by the server.                            |
| **Preconditions in Requests**      | 428 Precondition Required                                 | Conditional requests missing headers like `If-Match`.                               |
| **Unavailable Legal Restrictions** | 451 Unavailable For Legal Reasons                         | Server refuses access due to legal constraints (e.g., copyright violations).        |
| **Teapot Joke (RFC 2324)**         | 418 I'm a Teapot                                          | Easter egg for servers implementing the Hyper Text Coffee Pot Control Protocol.     |

#### Detailed Use of HTTP Status Codes

| **HTTP Status Code**                    | **Category**  | **Description**                                                                   | **Example**                                            |
| --------------------------------------- | ------------- | --------------------------------------------------------------------------------- | ------------------------------------------------------ |
| **100 Continue**                        | Informational | Request headers received; client can proceed with the request body.               | Client sends a large payload after server's readiness. |
| **200 OK**                              | Successful    | General success for GET, POST, PUT, or DELETE requests.                           | Returns requested resource or success confirmation.    |
| **201 Created**                         | Successful    | Resource successfully created.                                                    | `POST /users` creates a new user.                      |
| **202 Accepted**                        | Successful    | Request accepted but processing is asynchronous.                                  | Long-running jobs or file uploads.                     |
| **204 No Content**                      | Successful    | Request successful, but no response body is returned.                             | `DELETE /users/123` successfully deletes a user.       |
| **301 Moved Permanently**               | Redirection   | Resource has moved permanently to a new location.                                 | Redirect from `http://` to `https://`.                 |
| **304 Not Modified**                    | Redirection   | Cached content is still valid.                                                    | `GET` returns no new data if resource is unchanged.    |
| **400 Bad Request**                     | Client Error  | Malformed or invalid request from the client.                                     | Missing required fields or invalid input.              |
| **401 Unauthorized**                    | Client Error  | Authentication is required but missing or invalid.                                | Accessing protected routes without valid credentials.  |
| **403 Forbidden**                       | Client Error  | Access is denied due to insufficient permissions.                                 | Non-admin user attempting to access admin resources.   |
| **404 Not Found**                       | Client Error  | Requested resource could not be found.                                            | `GET /users/999` where the user doesn't exist.         |
| **409 Conflict**                        | Client Error  | Conflict with the current resource state.                                         | Versioning conflicts during updates.                   |
| **429 Too Many Requests**               | Client Error  | Too many requests sent by the client in a time window.                            | Exceeding API rate limits for a free-tier user.        |
| **500 Internal Server Error**           | Server Error  | Generic error due to unexpected server conditions.                                | Unhandled exception or database failure.               |
| **503 Service Unavailable**             | Server Error  | Server is temporarily unavailable (e.g., maintenance or overload).                | Server in maintenance mode.                            |
| **504 Gateway Timeout**                 | Server Error  | Upstream server timeout occurred.                                                 | API service failed while calling a dependent service.  |
| **205 Reset Content**                   | Successful    | Client should reset the view or form after the request.                           | Clearing form data after a successful submission.      |
| **206 Partial Content**                 | Successful    | Partial data returned, typically for range requests.                              | Video or file streaming.                               |
| **302 Found**                           | Redirection   | Resource temporarily located at a different URI.                                  | Temporary URL redirection during maintenance.          |
| **307 Temporary Redirect**              | Redirection   | POST request redirect maintaining HTTP method and body.                           | Redirecting after login to the original page.          |
| **308 Permanent Redirect**              | Redirection   | Resource has permanently moved; all requests must use the new location.           | Updating bookmarks for a relocated API endpoint.       |
| **413 Payload Too Large**               | Client Error  | Request entity exceeds the server’s capacity limits.                              | Uploading a file exceeding the server's maximum size.  |
| **414 URI Too Long**                    | Client Error  | Request URI is too long for the server to process.                                | Query parameters exceed the allowed URL length.        |
| **415 Unsupported Media Type**          | Client Error  | Media format of the request is not supported by the server.                       | Uploading a `.exe` file when only `.png` is allowed.   |
| **416 Range Not Satisfiable**           | Client Error  | Client requested an invalid range of data.                                        | Requesting byte range beyond the end of a file.        |
| **422 Unprocessable Entity**            | Client Error  | Server understands the request but is unable to process the content.              | Semantic validation errors for a JSON payload.         |
| **451 Unavailable For Legal Reasons**   | Client Error  | Resource cannot be provided due to legal restrictions.                            | Region-based content restrictions.                     |
| **500 Internal Server Error**           | Server Error  | Catch-all for unexpected server errors.                                           | Database unavailable or unhandled exceptions.          |
| **501 Not Implemented**                 | Server Error  | Server lacks the functionality to fulfill the request.                            | Unsupported HTTP method (e.g., `TRACE`).               |
| **502 Bad Gateway**                     | Server Error  | Server acting as a gateway received an invalid response from the upstream server. | Reverse proxy errors.                                  |
| **504 Gateway Timeout**                 | Server Error  | Timeout occurred while waiting for an upstream server.                            | API dependency fails to respond in time.               |
| **511 Network Authentication Required** | Server Error  | Network requires authentication to gain access.                                   | Captive portals in public Wi-Fi networks.              |

#### Test Case Scenarios

| **Test Case Scenario**                                                             | **Expected HTTP Status Code(s)**               | **Description**                                                                                                             |
| ---------------------------------------------------------------------------------- | ---------------------------------------------- | --------------------------------------------------------------------------------------------------------------------------- |
| **User successfully logs in**                                                      | 200 OK                                         | API response includes user details and authentication token.                                                                |
| **User account created successfully**                                              | 201 Created                                    | Resource (user account) is created, and the location of the new resource is returned.                                       |
| **File uploaded successfully but no additional data**                              | 204 No Content                                 | File upload completed successfully without a response body.                                                                 |
| **User submits invalid login credentials**                                         | 401 Unauthorized                               | Login fails due to incorrect username or password.                                                                          |
| **Accessing a protected resource without login**                                   | 403 Forbidden                                  | User is not authorized to view the resource, even with authentication.                                                      |
| **User attempts to access a non-existent endpoint**                                | 404 Not Found                                  | Requested API endpoint or resource does not exist.                                                                          |
| **Validation error on a user registration form**                                   | 422 Unprocessable Entity                       | Form data does not pass validation rules (e.g., password too short).                                                        |
| **Request body too large (e.g., large file upload)**                               | 413 Payload Too Large                          | File upload or JSON body exceeds server limits.                                                                             |
| **Incorrect file format uploaded (e.g., .exe file)**                               | 415 Unsupported Media Type                     | Server does not support the uploaded file type.                                                                             |
| **Pagination requested beyond total pages**                                        | 416 Range Not Satisfiable                      | Client requests data outside the valid range of pages or records.                                                           |
| **Server fails to connect to the database**                                        | 500 Internal Server Error                      | Unhandled error when the database connection cannot be established.                                                         |
| **Feature not yet implemented in the API**                                         | 501 Not Implemented                            | HTTP method or endpoint is recognized but not implemented.                                                                  |
| **API request timeout due to heavy backend processing**                            | 504 Gateway Timeout                            | Upstream server takes too long to process the request.                                                                      |
| **Resource moved to a new location**                                               | 301 Moved Permanently, 308 Permanent Redirect  | URL of the resource has changed permanently; clients should update their links.                                             |
| **Resource temporarily relocated during maintenance**                              | 302 Found, 307 Temporary Redirect              | Temporary redirect to another URL while the primary resource is unavailable.                                                |
| **Rate-limiting: Too many API requests**                                           | 429 Too Many Requests                          | Client exceeds the allowed number of API requests in a given time window.                                                   |
| **Precondition headers fail validation**                                           | 412 Precondition Failed                        | Conditional request fails because resource has been modified.                                                               |
| **JSON body missing required fields**                                              | 400 Bad Request                                | Client sends a malformed or incomplete request body.                                                                        |
| **Legal restriction prevents content access**                                      | 451 Unavailable For Legal Reasons              | Server is legally restricted from providing access to the resource.                                                         |
| **Client loses connection before the request completes**                           | 499 Client Closed Request                      | Client aborts the connection before receiving a response (commonly logged, not directly sent).                              |
| **Authentication required by proxy server**                                        | 407 Proxy Authentication Required              | Proxy server requires authentication to forward the request.                                                                |
| **File partially downloaded for video streaming**                                  | 206 Partial Content                            | Only a specific byte range of the resource is delivered for streaming or download resumption.                               |
| **Form submission resets after success**                                           | 205 Reset Content                              | Server advises client to clear the form view post submission.                                                               |
| **HTTP protocol not supported by server**                                          | 505 HTTP Version Not Supported                 | Server rejects the request because it uses an unsupported HTTP protocol version.                                            |
| **Conflict during resource creation (e.g., duplicate)**                            | 409 Conflict                                   | Duplicate data (e.g., username already taken) prevents successful resource creation.                                        |
| **Request missing mandatory headers**                                              | 428 Precondition Required                      | Server expects precondition headers but they are not provided in the request.                                               |
| **Server unavailable for scheduled maintenance**                                   | 503 Service Unavailable                        | API temporarily down due to maintenance.                                                                                    |
| **Captive portal blocking request**                                                | 511 Network Authentication Required            | Network requires user to authenticate via a captive portal (e.g., public Wi-Fi).                                            |
| **User attempts to delete a resource they don’t own**                              | 403 Forbidden                                  | User lacks the necessary permissions to delete the resource.                                                                |
| **Successful logout action**                                                       | 204 No Content                                 | Logout succeeds, and no further information is needed in the response body.                                                 |
| **Failed payment during checkout**                                                 | 402 Payment Required                           | Payment processing fails due to insufficient funds or invalid card details.                                                 |
| **Retryable error while interacting with an API**                                  | 503 Service Unavailable, 429 Too Many Requests | Server advises client to retry the request later due to high load or maintenance.                                           |
| **Resource has been permanently removed**                                          | 410 Gone                                       | The requested resource is no longer available and has been intentionally removed.                                           |
| **Attempting to use deprecated API functionality**                                 | 426 Upgrade Required                           | Server indicates the client must switch to a newer API version or feature.                                                  |
| **Server encounters infinite loop in processing**                                  | 508 Loop Detected                              | Recursive dependency or circular reference causes server to fail.                                                           |
| **Client request rejected for failing security policy**                            | 403 Forbidden                                  | Request rejected due to IP block, WAF policy, or lack of role-based access.                                                 |
| **Data synchronization between microservices fails**                               | 424 Failed Dependency                          | Dependent service failure causes the current request to fail.                                                               |
| **Cache expired and client retries fetching resource**                             | 304 Not Modified                               | Client reuses cached data because the resource hasn’t been modified.                                                        |
| **Conditional GET succeeds with ETag header validation**                           | 200 OK                                         | ETag or Last-Modified headers validate that the resource is still current.                                                  |
| **Debugging API errors via enhanced logs**                                         | 500 Internal Server Error                      | Logs include additional debugging information in non-production environments.                                               |
| **Database deadlock during bulk update**                                           | 503 Service Unavailable                        | Server temporarily pauses due to contention in backend services (e.g., database).                                           |
| **Sending an OPTIONS request for CORS preflight**                                  | 204 No Content                                 | Preflight request checks permissions and provides headers without returning a body.                                         |
| **User provides outdated resource version for update**                             | 409 Conflict                                   | Request fails because the resource version has changed (optimistic locking conflict).                                       |
| **Authorization header missing for private API**                                   | 401 Unauthorized                               | Server denies access to private endpoints due to missing or invalid authorization credentials.                              |
| **Request triggers alert for potential bot activity**                              | 429 Too Many Requests                          | Server throttles requests due to rate limits exceeded by suspected bot.                                                     |
| **Returning localized error messages**                                             | 422 Unprocessable Entity                       | Localized error message (e.g., user-facing validation errors) is included in the response.                                  |
| **API downtime monitored by a health check service**                               | 503 Service Unavailable                        | Health check detects outage, and server signals unavailability for monitoring tools.                                        |
| **Proxy service fails to resolve the target server**                               | 502 Bad Gateway                                | Reverse proxy server cannot reach the upstream server.                                                                      |
| **User requests unsupported language or locale**                                   | 406 Not Acceptable                             | Server cannot fulfill request due to missing translation or unsupported Accept-Language header.                             |
| **Request contains missing parameters**                                            | 400 Bad Request                                | Request body or query parameters are incomplete or improperly formatted.                                                    |
| **Rate-limited error with Retry-After header**                                     | 429 Too Many Requests                          | Response includes a Retry-After header to indicate when the client can send another request.                                |
| **Simultaneous update results in inconsistent state**                              | 409 Conflict                                   | Two users updating the same resource simultaneously cause a conflict.                                                       |
| **Long-running process accepted but not completed yet**                            | 202 Accepted                                   | Request acknowledged and being processed asynchronously; no result yet.                                                     |
| **Streaming file in chunks for large downloads**                                   | 206 Partial Content                            | Response delivers file in chunks using Content-Range headers.                                                               |
| **Server limits resource creation per account**                                    | 429 Too Many Requests                          | Too many POST requests for creating new resources are detected from the same account.                                       |
| **Request contains invalid authentication token**                                  | 401 Unauthorized                               | JWT token expired, malformed, or signed with an invalid key.                                                                |
| **Custom error page for 404 not found**                                            | 404 Not Found                                  | Server renders a user-friendly 404 error page for missing resources.                                                        |
| **JSON parsing error during request processing**                                   | 400 Bad Request                                | Malformed JSON request body results in a bad request error.                                                                 |
| **File uploaded to a temporary storage location**                                  | 201 Created                                    | Temporary resource is created, and location is returned in response headers.                                                |
| **Upstream server returns a malformed response**                                   | 502 Bad Gateway                                | Reverse proxy reports an invalid or unexpected response from the upstream server.                                           |
| **Deprecation warning in response headers**                                        | 200 OK                                         | Warning header advises clients that the endpoint or feature will be deprecated in future.                                   |
| **Custom quota limits for premium users**                                          | 429 Too Many Requests                          | Rate-limiting based on user’s subscription tier is enforced.                                                                |
| **User provides invalid credentials**                                              | 401 Unauthorized                               | Authentication fails due to incorrect or missing credentials in the request.                                                |
| **User requests a resource that no longer exists**                                 | 410 Gone                                       | Requested resource was previously available but has been permanently deleted.                                               |
| **File upload is too large**                                                       | 413 Payload Too Large                          | User attempts to upload a file that exceeds the allowed size.                                                               |
| **Access is denied due to insufficient permissions**                               | 403 Forbidden                                  | User lacks permissions to access the requested resource, even though they are authenticated.                                |
| **User tries to access a resource they are not authorized for**                    | 403 Forbidden                                  | The user is authenticated, but their role does not permit access to the resource.                                           |
| **API version mismatch**                                                           | 426 Upgrade Required                           | The client needs to upgrade to a supported API version.                                                                     |
| **API successfully processes a request**                                           | 200 OK                                         | Request is processed successfully, and the server returns the result.                                                       |
| **Resource creation is successful**                                                | 201 Created                                    | A new resource is created (e.g., new user, new item) and the response includes a location header.                           |
| **User tries to update a resource with invalid data**                              | 400 Bad Request                                | The user has sent invalid data or missing required parameters in the request body.                                          |
| **Invalid input format or type in request body**                                   | 400 Bad Request                                | Malformed JSON, XML, or any unsupported data format received in the body.                                                   |
| **Client exceeds rate limit for API calls**                                        | 429 Too Many Requests                          | Too many requests sent in a given time frame, exceeding the rate limit.                                                     |
| **Resource is temporarily unavailable due to maintenance**                         | 503 Service Unavailable                        | Service is down due to scheduled or unscheduled maintenance.                                                                |
| **Requested resource does not exist**                                              | 404 Not Found                                  | The server cannot find the resource requested by the client.                                                                |
| **User is logged out or session expired**                                          | 401 Unauthorized                               | Session has expired, requiring the user to authenticate again.                                                              |
| **API request is valid but cannot be processed immediately**                       | 202 Accepted                                   | Request is accepted, but processing is deferred (asynchronous).                                                             |
| **Invalid request headers**                                                        | 400 Bad Request                                | Client sends improper headers, such as missing required headers or unsupported content types.                               |
| **File type not supported during upload**                                          | 415 Unsupported Media Type                     | The server refuses to process the file because its type is not supported.                                                   |
| **Client attempts to access a restricted IP address**                              | 403 Forbidden                                  | Server blocks access from a specific IP address for security reasons.                                                       |
| **Required query parameters are missing**                                          | 400 Bad Request                                | A required query parameter is missing from the request, making it invalid.                                                  |
| **Client sends an invalid or expired authentication token**                        | 401 Unauthorized                               | Authentication fails due to the expiration or invalidity of the provided authentication token.                              |
| **The request is processed but some resources fail**                               | 206 Partial Content                            | The request was partially fulfilled; some resources were returned, but others failed.                                       |
| **Client requested a resource that has been moved**                                | 301 Moved Permanently                          | The requested resource has been permanently moved to a new URL.                                                             |
| **Client requests an old resource version**                                        | 409 Conflict                                   | A conflict occurs, typically when trying to update a resource that has been modified by others.                             |
| **Resource modification fails due to business rules**                              | 422 Unprocessable Entity                       | The resource update fails due to invalid data (e.g., violating constraints or business rules).                              |
| **Server error during processing due to unexpected issue**                         | 500 Internal Server Error                      | Generic error indicating that the server encountered an unexpected condition.                                               |
| **Proxy server encounters an error while communicating with the backend**          | 502 Bad Gateway                                | The reverse proxy server is unable to forward the request to the backend service.                                           |
| **Request body exceeds maximum allowed size**                                      | 413 Payload Too Large                          | The client sends a request with a body larger than what the server is willing or able to process.                           |
| **User attempts to upload a file without authorization**                           | 403 Forbidden                                  | User is authenticated, but lacks permission to upload files to the server.                                                  |
| **API request violates business logic (e.g., attempting to submit an empty form)** | 400 Bad Request                                | Request violates API business logic or validation rules (e.g., submitting empty or invalid fields).                         |
| **Server does not support the requested HTTP method**                              | 405 Method Not Allowed                         | The client attempts to use an unsupported HTTP method (e.g., using PUT on a read-only endpoint).                            |
| **Client attempts to submit a request with too many parameters**                   | 400 Bad Request                                | Request contains too many parameters or exceeds the API’s parameter limits.                                                 |
| **Service is temporarily unavailable**                                             | 503 Service Unavailable                        | The service is temporarily unavailable due to reasons like server overload, maintenance, etc.                               |
| **Invalid request format (e.g., XML instead of JSON)**                             | 415 Unsupported Media Type                     | Server responds with an error when an unsupported content type (like XML) is sent in the request.                           |
| **File not found in requested location**                                           | 404 Not Found                                  | Client requests a file that does not exist on the server.                                                                   |
| **Server experiences a timeout when processing request**                           | 504 Gateway Timeout                            | The server was unable to get a timely response from an upstream server.                                                     |
| **Incorrect authentication method provided**                                       | 400 Bad Request                                | Authentication method provided by the client is not supported by the server (e.g., Basic Auth instead of OAuth).            |
| **Request sent in unsupported language**                                           | 406 Not Acceptable                             | The server cannot provide a response in the requested language.                                                             |
| **Client sends invalid API version**                                               | 426 Upgrade Required                           | API version used in the request is outdated and no longer supported by the server.                                          |
| **The API is receiving an unexpected large number of requests**                    | 429 Too Many Requests                          | The client exceeded the rate limit defined for the API endpoint.                                                            |
| **Client requests a resource that is locked**                                      | 423 Locked                                     | The requested resource is locked and cannot be modified.                                                                    |
| **Request failed due to a custom business rule violation**                         | 422 Unprocessable Entity                       | The request contains valid input, but a business rule prevents it from being processed.                                     |
| **Client tries to access a deprecated endpoint**                                   | 410 Gone                                       | The endpoint has been deprecated and is no longer available.                                                                |
| **The server is unable to process the request due to a backend error**             | 502 Bad Gateway                                | Backend or upstream server issues prevent the processing of the request.                                                    |
| **The client’s request is incomplete (missing necessary fields)**                  | 400 Bad Request                                | Request cannot be processed because it lacks necessary fields or contains invalid values.                                   |
| **User reaches the maximum limit of API calls per day**                            | 429 Too Many Requests                          | The user has exceeded their daily API quota limit and must wait until the next period.                                      |
| **Client tries to access a resource that requires a higher subscription tier**     | 403 Forbidden                                  | Client is attempting to access a resource that is restricted based on their subscription tier.                              |
| **Client sends a request with a malformed header**                                 | 400 Bad Request                                | A malformed or missing required header causes the request to be rejected.                                                   |
| **Server successfully processes a POST request to create a new entity**            | 201 Created                                    | The POST request creates a new entity (such as a user or item), and the server returns the URI of the newly created entity. |
| **Client submits a valid request for batch processing**                            | 207 Multi-Status                               | A batch request is processed with mixed success, returning the status of each request in the batch.                         |
| **Request for an unsupported API version**                                         | 426 Upgrade Required                           | The requested API version is unsupported, and the client must upgrade.                                                      |

### Contributing

To contribute to project, follow these steps:

1. Clone the repository:

   ```bash
   git clone --depth 1 https://github.com/sivaosorg/replify.git
   ```

2. Navigate to the project directory:

   ```bash
   cd replify
   ```

3. Prepare the project environment:
   ```bash
   go mod tidy
   ```
