package requester

import "github.com/Nigel2392/jsext/requester/encoders"

// Define methods and encodings type
type Methods string
type Encoding string

// Define standard error messages
var (
	ErrNoRequest  = "no request has been set"                  // Error message for when no request has been set
	ErrNoCallback = "no callback has been set"                 // Error message for when no callback has been set
	ErrNoEncoding = "no encoding has been set or is not valid" // Error message for when no encoding has been set
)

// Define request methods
const (
	GET     Methods = "GET"     // GET method
	POST    Methods = "POST"    // POST method
	PUT     Methods = "PUT"     // PUT method
	PATCH   Methods = "PATCH"   // PATCH method
	DELETE  Methods = "DELETE"  // DELETE method
	OPTIONS Methods = "OPTIONS" // OPTIONS method
	HEAD    Methods = "HEAD"    // HEAD method
	TRACE   Methods = "TRACE"   // TRACE method

)

// Define methods of encoding
const (
	FORM_URL_ENCODED Encoding = "application/x-www-form-urlencoded" // FORM_URL_ENCODED encoding
	MULTIPART_FORM   Encoding = "multipart/form-data"               // MULTIPART_FORM encoding
	JSON             Encoding = "json"                              // JSON encoding
	XML              Encoding = "xml"                               // XML encoding
)

type Client interface {
	// Initialize a GET request
	Get(url string) Client
	// Initialize a POST request
	Post(url string) Client
	// Initialize a PUT request
	Put(url string) Client
	// Initialize a PATCH request
	Patch(url string) Client
	// Initialize a DELETE request
	Delete(url string) Client
	// Initialize a OPTIONS request
	Options(url string) Client
	// Initialize a HEAD request
	Head(url string) Client
	// Add form data to the request
	WithData(formData map[string]interface{}, encoding Encoding, file ...encoders.File) Client
	// Add a header to the request
	WithHeader(header map[string]string) Client
	// Add a query to the request
	WithQuery(query map[string]string) Client
}
