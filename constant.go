package replify

// Standard HTTP headers related to content negotiation and encoding.
const (
	// Accept specifies the media types that are acceptable for the response.
	// 	Example: "application/json, text/html"
	HeaderAccept = "Accept"

	// AcceptCharset specifies the character sets that are acceptable.
	// 	Example: "utf-8, iso-8859-1"
	HeaderAcceptCharset = "Accept-Charset"

	// AcceptEncoding specifies the content encodings that are acceptable.
	//	Example: "gzip, deflate, br"
	HeaderAcceptEncoding = "Accept-Encoding"

	// AcceptLanguage specifies the acceptable languages for the response.
	// 	Example: "en-US, en;q=0.9, fr;q=0.8"
	HeaderAcceptLanguage = "Accept-Language"

	// Authorization contains the credentials for authenticating the client with the server.
	// 	Example: "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6..."
	HeaderAuthorization = "Authorization"

	// CacheControl specifies directives for caching mechanisms in both requests and responses.
	// 	Example: "no-cache, no-store, must-revalidate"
	HeaderCacheControl = "Cache-Control"

	// ContentDisposition specifies if the content should be displayed inline or treated as an attachment.
	// 	Example: "attachment; filename=\"document.pdf\""
	HeaderContentDisposition = "Content-Disposition"

	// ContentEncoding specifies the encoding transformations that have been applied to the body of the response.
	// 	Example: "gzip"
	HeaderContentEncoding = "Content-Encoding"

	// ContentLength specifies the size of the response body in octets.
	// 	Example: "1024"
	HeaderContentLength = "Content-Length"

	// ContentType specifies the media type of the resource.
	// 	Example: "application/json; charset=utf-8"
	HeaderContentType = "Content-Type"

	// Cookie contains stored HTTP cookies sent to the server by the client.
	// 	Example: "sessionId=abc123; userId=456"
	HeaderCookie = "Cookie"

	// Host specifies the domain name of the server (for virtual hosting) and the TCP port number.
	// 	Example: "www.example.com:8080"
	HeaderHost = "Host"

	// Origin specifies the origin of the cross-origin request or preflight request.
	// 	Example: "https://www.example.com"
	HeaderOrigin = "Origin"

	// Referer contains the address of the previous web page from which a link to the currently requested page was followed.
	// 	Example: "https://www.example.com/page1.html"
	HeaderReferer = "Referer"

	// UserAgent contains information about the user agent (browser or client) making the request.
	// 	Example: "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36"
	HeaderUserAgent = "User-Agent"

	// IfMatch makes the request conditional on the target resource having the same entity tag as the one provided.
	// 	Example: "\"686897696a7c876b7e\""
	HeaderIfMatch = "If-Match"

	// IfNoneMatch makes the request conditional on the target resource not having the same entity tag as the one provided.
	// 	Example: "\"686897696a7c876b7e\""
	HeaderIfNoneMatch = "If-None-Match"

	// ETag provides the entity tag for the resource.
	// 	Example: "\"33a64df551425fcc55e4d42a148795d9f25f89d4\""
	HeaderETag = "ETag"

	// LastModified specifies the last modified date of the resource.
	// 	Example: "Wed, 21 Oct 2015 07:28:00 GMT"
	HeaderLastModified = "Last-Modified"

	// Location specifies the URL to redirect a client to.
	// 	Example: "https://www.example.com/new-location"
	HeaderLocation = "Location"

	// Pragma specifies implementation-specific directives that might affect caching.
	// 	Example: "no-cache"
	HeaderPragma = "Pragma"

	// RetryAfter specifies the time after which the client should retry the request after receiving a 503 Service Unavailable status code.
	// 	Example: "120" or "Fri, 07 Nov 2014 23:59:59 GMT"
	HeaderRetryAfter = "Retry-After"

	// Server contains information about the software used by the origin server to handle the request.
	// 	Example: "Apache/2.4.41 (Ubuntu)"
	HeaderServer = "Server"

	// WWWAuthenticate indicates that the client must authenticate to access the requested resource.
	// 	Example: "Basic realm=\"Access to staging site\""
	HeaderWWWAuthenticate = "WWW-Authenticate"

	// Date specifies the date and time at which the message was sent.
	// 	Example: "Tue, 15 Nov 1994 08:12:31 GMT"
	HeaderDate = "Date"

	// Expires specifies the date/time after which the response is considered stale.
	// 	Example: "Thu, 01 Dec 1994 16:00:00 GMT"
	HeaderExpires = "Expires"

	// Age specifies the age of the response in seconds.
	// 	Example: "3600"
	HeaderAge = "Age"

	// Connection specifies control options for the current connection (e.g., keep-alive or close).
	// 	Example: "keep-alive"
	HeaderConnection = "Connection"

	// ContentLanguage specifies the language of the content.
	// 	Example: "en-US"
	HeaderContentLanguage = "Content-Language"

	// Forwarded contains information about intermediate proxies or gateways that have forwarded the request.
	// 	Example: "for=192.0.2.60;proto=http;by=203.0.113.43"
	HeaderForwarded = "Forwarded"

	// IfModifiedSince makes the request conditional on the target resource being modified since the specified date.
	// 	Example: "Wed, 21 Oct 2015 07:28:00 GMT"
	HeaderIfModifiedSince = "If-Modified-Since"

	// Upgrade requests the server to switch to a different protocol.
	// 	Example: "websocket"
	HeaderUpgrade = "Upgrade"

	// Via provides information about intermediate protocols and recipients between the user agent and the server.
	// 	Example: "1.1 proxy1.example.com, 1.0 proxy2.example.org"
	HeaderVia = "Via"

	// Warning carries additional information about the status or transformation of a message.
	// 	Example: "110 anderson/1.3.37 \"Response is stale\""
	HeaderWarning = "Warning"

	// XForwardedFor contains the originating IP address of a client connecting to a web server through an HTTP proxy or load balancer.
	// 	Example: "203.0.113.195, 70.41.3.18, 150.172.238.178"
	HeaderXForwardedFor = "X-Forwarded-For"

	// XForwardedHost contains the original host requested by the client in the Host HTTP request header.
	// 	Example: "example.com"
	HeaderXForwardedHost = "X-Forwarded-Host"

	// XForwardedProto specifies the protocol (HTTP or HTTPS) used by the client.
	// 	Example: "https"
	HeaderXForwardedProto = "X-Forwarded-Proto"

	// XRequestedWith identifies the type of request being made (e.g., Ajax requests).
	// 	Example: "XMLHttpRequest"
	HeaderXRequestedWith = "X-Requested-With"

	// XFrameOptions specifies whether the browser should be allowed to render the page in a <frame>, <iframe>, <object>, <embed>, or <applet>.
	// 	Example: "DENY" or "SAMEORIGIN"
	HeaderXFrameOptions = "X-Frame-Options"

	// XXSSProtection controls browser's built-in XSS (Cross-Site Scripting) filter.
	// 	Example: "1; mode=block"
	HeaderXXSSProtection = "X-XSS-Protection"

	// XContentTypeOpts prevents browsers from interpreting files as a different MIME type than what is specified.
	// 	Example: "nosniff"
	HeaderXContentTypeOpts = "X-Content-Type-Options"

	// ContentSecurity specifies security policy for web applications, helping to prevent certain types of attacks.
	// 	Example: "default-src 'self'; script-src 'self' 'unsafe-inline'"
	HeaderContentSecurity = "Content-Security-Policy"

	// StrictTransport enforces the use of HTTPS for the website to reduce security risks.
	// 	Example: "max-age=31536000; includeSubDomains"
	HeaderStrictTransport = "Strict-Transport-Security"

	// PublicKeyPins specifies public key pins to prevent man-in-the-middle attacks.
	// 	Example: "pin-sha256=\"base64+primary==\"; pin-sha256=\"base64+backup==\"; max-age=5184000"
	HeaderPublicKeyPins = "Public-Key-Pins"

	// ExpectCT allows websites to specify a Certificate Transparency policy.
	// 	Example: "max-age=86400, enforce"
	HeaderExpectCT = "Expect-CT"

	// AccessControlAllowOrigin specifies which domains are allowed to access the resources.
	// 	Example: "*" or "https://example.com"
	HeaderAccessControlAllowOrigin = "Access-Control-Allow-Origin"

	// AccessControlAllowMethods specifies which HTTP methods are allowed when accessing the resource.
	// 	Example: "GET, POST, PUT, DELETE"
	HeaderAccessControlAllowMethods = "Access-Control-Allow-Methods"

	// AccessControlAllowHeaders specifies which HTTP headers can be used during the actual request.
	// 	Example: "Content-Type, Authorization"
	HeaderAccessControlAllowHeaders = "Access-Control-Allow-Headers"

	// AccessControlMaxAge specifies how long the results of a preflight request can be cached.
	// 	Example: "86400"
	HeaderAccessControlMaxAge = "Access-Control-Max-Age"

	// AccessControlExposeHeaders specifies which headers can be exposed as part of the response.
	// 	Example: "Content-Length, X-JSON"
	HeaderAccessControlExposeHeaders = "Access-Control-Expose-Headers"

	// AccessControlRequestMethod indicates which HTTP method will be used during the actual request.
	// 	Example: "POST"
	HeaderAccessControlRequestMethod = "Access-Control-Request-Method"

	// AccessControlRequestHeaders specifies which headers can be sent with the actual request.
	// 	Example: "Content-Type, X-Custom-Header"
	HeaderAccessControlRequestHeaders = "Access-Control-Request-Headers"

	// AcceptPatch specifies which patch document formats are acceptable in the response.
	// 	Example: "application/json-patch+json"
	HeaderAcceptPatch = "Accept-Patch"

	// DeltaBase specifies the URI of the delta information.
	// 	Example: "\"abc123\""
	HeaderDeltaBase = "Delta-Base"

	// IfUnmodifiedSince makes the request conditional on the resource not being modified since the specified date.
	// 	Example: "Wed, 21 Oct 2015 07:28:00 GMT"
	HeaderIfUnmodifiedSince = "If-Unmodified-Since"

	// AcceptRanges specifies the range of the resource that the client is requesting.
	// 	Example: "bytes"
	HeaderAcceptRanges = "Accept-Ranges"

	// ContentRange specifies the range of the resource being sent in the response.
	// 	Example: "bytes 200-1000/5000"
	HeaderContentRange = "Content-Range"

	// Allow specifies the allowed methods for a resource.
	// 	Example: "GET, HEAD, PUT"
	HeaderAllow = "Allow"

	// AccessControlAllowCredentials indicates whether the response to the request can expose credentials.
	// 	Example: "true"
	HeaderAccessControlAllowCredentials = "Access-Control-Allow-Credentials"

	// XCSRFToken is used to prevent Cross-Site Request Forgery (CSRF) attacks.
	// 	Example: "i8XNjC4b8KVok4uw5RftR38Wgp2BF"
	HeaderXCSRFToken = "X-CSRF-Token"

	// XRealIP contains the real IP address of the client, often used in proxies or load balancers.
	// 	Example: "203.0.113.195"
	HeaderXRealIP = "X-Real-IP"

	// ContentSecurityPolicy specifies content security policies to prevent certain attacks.
	// 	Example: "default-src 'self'; img-src *; media-src media1.com media2.com"
	HeaderContentSecurityPolicy = "Content-Security-Policy"

	// ReferrerPolicy controls how much information about the referring page is sent.
	// 	Example: "no-referrer-when-downgrade"
	HeaderReferrerPolicy = "Referrer-Policy"

	// ExpectCt specifies a Certificate Transparency policy for the web server.
	// 	Example: "max-age=86400, enforce, report-uri=\"https://example.com/report\""
	HeaderExpectCt = "Expect-CT"

	// StrictTransportSecurity enforces HTTPS to reduce the chance of security breaches.
	// 	Example: "max-age=63072000; includeSubDomains; preload"
	HeaderStrictTransportSecurity = "Strict-Transport-Security"

	// UpgradeInsecureRequests requests the browser to upgrade any insecure requests to secure HTTPS requests.
	// 	Example: "1"
	HeaderUpgradeInsecureRequests = "Upgrade-Insecure-Requests"
)

// Media Type constants define commonly used MIME types for different content types in HTTP requests and responses.
const (
	// ApplicationJSON specifies that the content is JSON-formatted data.
	//  Example: "application/json"
	MediaTypeApplicationJSON = "application/json"

	// ApplicationXML specifies that the content is XML-formatted data.
	//  Example: "application/xml"
	MediaTypeApplicationXML = "application/xml"

	// ApplicationForm specifies that the content is URL-encoded form data.
	//  Example: "application/x-www-form-urlencoded"
	MediaTypeApplicationForm = "application/x-www-form-urlencoded"

	// ApplicationOctetStream specifies that the content is binary data (not interpreted by the browser).
	//  Example: "application/octet-stream"
	MediaTypeApplicationOctetStream = "application/octet-stream"

	// TextPlain specifies that the content is plain text.
	//  Example: "text/plain"
	MediaTypeTextPlain = "text/plain"

	// TextHTML specifies that the content is HTML-formatted data.
	//  Example: "text/html"
	MediaTypeTextHTML = "text/html"

	// ImageJPEG specifies that the content is a JPEG image.
	//  Example: "image/jpeg"
	MediaTypeImageJPEG = "image/jpeg"

	// ImagePNG specifies that the content is a PNG image.
	//  Example: "image/png"
	MediaTypeImagePNG = "image/png"

	// ImageGIF specifies that the content is a GIF image.
	//  Example: "image/gif"
	MediaTypeImageGIF = "image/gif"

	// AudioMP3 specifies that the content is an MP3 audio file.
	//  Example: "audio/mpeg"
	MediaTypeAudioMP3 = "audio/mpeg"

	// AudioWAV specifies that the content is a WAV audio file.
	//  Example: "audio/wav"
	MediaTypeAudioWAV = "audio/wav"

	// VideoMP4 specifies that the content is an MP4 video file.
	//  Example: "video/mp4"
	MediaTypeVideoMP4 = "video/mp4"

	// VideoAVI specifies that the content is an AVI video file.
	//  Example: "video/x-msvideo"
	MediaTypeVideoAVI = "video/x-msvideo"

	// ApplicationPDF specifies that the content is a PDF file.
	//  Example: "application/pdf"
	MediaTypeApplicationPDF = "application/pdf"

	// ApplicationMSWord specifies that the content is a Microsoft Word document (.doc).
	//  Example: "application/msword"
	MediaTypeApplicationMSWord = "application/msword"

	// ApplicationMSPowerPoint specifies that the content is a Microsoft PowerPoint presentation (.ppt).
	//  Example: "application/vnd.ms-powerpoint"
	MediaTypeApplicationMSPowerPoint = "application/vnd.ms-powerpoint"

	// ApplicationExcel specifies that the content is a Microsoft Excel spreadsheet (.xls).
	//  Example: "application/vnd.ms-excel"
	MediaTypeApplicationExcel = "application/vnd.ms-excel"

	// ApplicationZip specifies that the content is a ZIP archive.
	//  Example: "application/zip"
	MediaTypeApplicationZip = "application/zip"

	// ApplicationGzip specifies that the content is a GZIP-compressed file.
	//  Example: "application/gzip"
	MediaTypeApplicationGzip = "application/gzip"

	// MultipartFormData specifies that the content is a multipart form, typically used for file uploads.
	//  Example: "multipart/form-data; boundary=----WebKitFormBoundary7MA4YWxkTrZu0gW"
	MediaTypeMultipartFormData = "multipart/form-data"

	// ImageBMP specifies that the content is a BMP image.
	//  Example: "image/bmp"
	MediaTypeImageBMP = "image/bmp"

	// ImageTIFF specifies that the content is a TIFF image.
	//  Example: "image/tiff"
	MediaTypeImageTIFF = "image/tiff"

	// TextCSS specifies that the content is CSS (Cascading Style Sheets).
	//  Example: "text/css"
	MediaTypeTextCSS = "text/css"

	// TextJavaScript specifies that the content is JavaScript code.
	//  Example: "text/javascript"
	MediaTypeTextJavaScript = "text/javascript"

	// ApplicationJSONLD specifies that the content is a JSON-LD (JSON for Linked Data) document.
	//  Example: "application/ld+json"
	MediaTypeApplicationJSONLD = "application/ld+json"

	// ApplicationRDFXML specifies that the content is in RDF (Resource Description Framework) XML format.
	//  Example: "application/rdf+xml"
	MediaTypeApplicationRDFXML = "application/rdf+xml"

	// ApplicationGeoJSON specifies that the content is a GeoJSON (geospatial data) document.
	//  Example: "application/geo+json"
	MediaTypeApplicationGeoJSON = "application/geo+json"

	// ApplicationMsgpack specifies that the content is in MessagePack format (binary JSON).
	//  Example: "application/msgpack"
	MediaTypeApplicationMsgpack = "application/msgpack"

	// ApplicationOgg specifies that the content is an Ogg multimedia container format.
	//  Example: "application/ogg"
	MediaTypeApplicationOgg = "application/ogg"

	// ApplicationGraphQL specifies that the content is in GraphQL format.
	//  Example: "application/graphql"
	MediaTypeApplicationGraphQL = "application/graphql"

	// ApplicationProtobuf specifies that the content is in Protocol Buffers format (binary serialization).
	//  Example: "application/protobuf"
	MediaTypeApplicationProtobuf = "application/protobuf"

	// ImageWebP specifies that the content is a WebP image.
	//  Example: "image/webp"
	MediaTypeImageWebP = "image/webp"

	// FontWOFF specifies that the content is a WOFF (Web Open Font Format) font.
	//  Example: "font/woff"
	MediaTypeFontWOFF = "font/woff"

	// FontWOFF2 specifies that the content is a WOFF2 (Web Open Font Format 2) font.
	//  Example: "font/woff2"
	MediaTypeFontWOFF2 = "font/woff2"

	// AudioFLAC specifies that the content is a FLAC audio file (Free Lossless Audio Codec).
	//  Example: "audio/flac"
	MediaTypeAudioFLAC = "audio/flac"

	// VideoWebM specifies that the content is a WebM video file.
	//  Example: "video/webm"
	MediaTypeVideoWebM = "video/webm"

	// ApplicationDart specifies that the content is a Dart programming language file.
	//  Example: "application/dart"
	MediaTypeApplicationDart = "application/dart"

	// ApplicationXLSX specifies that the content is an Excel file in XLSX format.
	//  Example: "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"
	MediaTypeApplicationXLSX = "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"

	// ApplicationPPTX specifies that the content is a PowerPoint file in PPTX format.
	//  Example: "application/vnd.openxmlformats-officedocument.presentationml.presentation"
	MediaTypeApplicationPPTX = "application/vnd.openxmlformats-officedocument.presentationml.presentation"

	// ApplicationGRPC specifies that the content is in gRPC format (a high-performance RPC framework).
	//  Example: "application/grpc"
	MediaTypeApplicationGRPC = "application/grpc"
)

var (
	// 1xx Informational responses
	// Continue indicates that the initial part of the request has been received and has not yet been rejected by the server.
	Continue = Header().WithCode(100).WithText("Continue").WithType("Informational")
	// SwitchingProtocols indicates that the server will switch protocols as requested by the client.
	SwitchingProtocols = Header().WithCode(101).WithText("Switching Protocols").WithType("Informational")
	// Processing indicates that the server has received and is processing the request but no response is available yet.
	Processing = Header().WithCode(102).WithText("Processing").WithType("Informational")

	// 2xx Successful responses
	// OK indicates that the request has succeeded.
	OK = Header().WithCode(200).WithText("OK").WithType("Successful")
	// Created indicates that the request has been fulfilled and has resulted in a new resource being created.
	Created = Header().WithCode(201).WithText("Created").WithType("Successful")
	// Accepted indicates that the request has been accepted for processing, but the processing has not been completed.
	Accepted = Header().WithCode(202).WithText("Accepted").WithType("Successful")
	// NonAuthoritativeInformation indicates that the request was successful, but the enclosed metadata may be from a different source.
	NonAuthoritativeInformation = Header().WithCode(203).WithText("Non-Authoritative Information").WithType("Successful")
	// NoContent indicates that the server successfully processed the request, but is not returning any content.
	NoContent = Header().WithCode(204).WithText("No Content").WithType("Successful")
	// ResetContent indicates that the server successfully processed the request and requests the client to reset the document view.
	ResetContent = Header().WithCode(205).WithText("Reset Content").WithType("Successful")
	// PartialContent indicates that the server is delivering only part of the resource due to a range request.
	PartialContent = Header().WithCode(206).WithText("Partial Content").WithType("Successful")
	// MultiStatus provides status for multiple independent operations.
	MultiStatus = Header().WithCode(207).WithText("Multi-Status").WithType("Successful")
	// AlreadyReported indicates that the members of a DAV binding have already been enumerated in a previous reply.
	AlreadyReported = Header().WithCode(208).WithText("Already Reported").WithType("Successful")
	// IMUsed indicates that the server has fulfilled a GET request for the resource and the response is a representation of the result.
	IMUsed = Header().WithCode(226).WithText("IM Used").WithType("Successful")

	// 3xx Redirection responses
	// MultipleChoices indicates multiple options for the resource are available.
	MultipleChoices = Header().WithCode(300).WithText("Multiple Choices").WithType("Redirection")
	// MovedPermanently indicates that the resource has been permanently moved to a new URI.
	MovedPermanently = Header().WithCode(301).WithText("Moved Permanently").WithType("Redirection")
	// Found indicates that the resource has been temporarily moved to a different URI.
	Found = Header().WithCode(302).WithText("Found").WithType("Redirection")
	// SeeOther indicates that the response to the request can be found under another URI.
	SeeOther = Header().WithCode(303).WithText("See Other").WithType("Redirection")
	// NotModified indicates that the resource has not been modified since the last request.
	NotModified = Header().WithCode(304).WithText("Not Modified").WithType("Redirection")
	// UseProxy indicates that the requested resource must be accessed through the proxy given by the location field.
	UseProxy = Header().WithCode(305).WithText("Use Proxy").WithType("Redirection")
	// Reserved is a deprecated status code reserved for future use.
	Reserved = Header().WithCode(306).WithText("Reserved").WithType("Redirection")
	// TemporaryRedirect indicates that the resource has been temporarily moved to a different URI and will return to the original URI later.
	TemporaryRedirect = Header().WithCode(307).WithText("Temporary Redirect").WithType("Redirection")
	// PermanentRedirect indicates that the resource has been permanently moved to a new URI and future requests should use this URI.
	PermanentRedirect = Header().WithCode(308).WithText("Permanent Redirect").WithType("Redirection")

	// 4xx Client error responses
	// BadRequest indicates that the server could not understand the request due to invalid syntax.
	BadRequest = Header().WithCode(400).WithText("Bad Request").WithType("Client Error")
	// Unauthorized indicates that the client must authenticate itself to get the requested response.
	Unauthorized = Header().WithCode(401).WithText("Unauthorized").WithType("Client Error")
	// PaymentRequired is reserved for future use, indicating payment is required to access the resource.
	PaymentRequired = Header().WithCode(402).WithText("Payment Required").WithType("Client Error")
	// Forbidden indicates that the server understands the request but refuses to authorize it.
	Forbidden = Header().WithCode(403).WithText("Forbidden").WithType("Client Error")
	// NotFound indicates that the server can't find the requested resource.
	NotFound = Header().WithCode(404).WithText("Not Found").WithType("Client Error")
	// MethodNotAllowed indicates that the server knows the request method but the target resource doesn't support this method.
	MethodNotAllowed = Header().WithCode(405).WithText("Method Not Allowed").WithType("Client Error")
	// NotAcceptable indicates that the server cannot produce a response matching the list of acceptable values defined in the request's headers.
	NotAcceptable = Header().WithCode(406).WithText("Not Acceptable").WithType("Client Error")
	// ProxyAuthenticationRequired indicates that the client must first authenticate itself with the proxy.
	ProxyAuthenticationRequired = Header().WithCode(407).WithText("Proxy Authentication Required").WithType("Client Error")
	// RequestTimeout indicates that the server timed out waiting for the request.
	RequestTimeout = Header().WithCode(408).WithText("Request Timeout").WithType("Client Error")
	// Conflict indicates that the request conflicts with the current state of the server.
	Conflict = Header().WithCode(409).WithText("Conflict").WithType("Client Error")
	// Gone indicates that the requested resource is no longer available and will not be available again.
	Gone = Header().WithCode(410).WithText("Gone").WithType("Client Error")
	// LengthRequired indicates that the server requires the request to be sent with a Content-Length header.
	LengthRequired = Header().WithCode(411).WithText("Length Required").WithType("Client Error")
	// PreconditionFailed indicates that the server does not meet one of the preconditions set by the client.
	PreconditionFailed = Header().WithCode(412).WithText("Precondition Failed").WithType("Client Error")
	// RequestEntityTooLarge indicates that the request entity is larger than what the server is willing or able to process.
	RequestEntityTooLarge = Header().WithCode(413).WithText("Request Entity Too Large").WithType("Client Error")
	// RequestURITooLong indicates that the URI provided was too long for the server to process.
	RequestURITooLong = Header().WithCode(414).WithText("Request-URI Too Long").WithType("Client Error")
	// UnsupportedMediaType indicates that the media format of the requested data is not supported by the server.
	UnsupportedMediaType = Header().WithCode(415).WithText("Unsupported Media Type").WithType("Client Error")
	// RequestedRangeNotSatisfiable indicates that the range specified by the Range header cannot be satisfied.
	RequestedRangeNotSatisfiable = Header().WithCode(416).WithText("Requested Range Not Satisfiable").WithType("Client Error")
	// ExpectationFailed indicates that the server cannot meet the requirements of the Expect request-header field.
	ExpectationFailed = Header().WithCode(417).WithText("Expectation Failed").WithType("Client Error")
	// ImATeapot is a humorous response code indicating that the server is a teapot and refuses to brew coffee.
	ImATeapot = Header().WithCode(418).WithText("I’m a teapot").WithType("Client Error")
	// EnhanceYourCalm is a non-standard response code used to ask the client to reduce its request rate.
	EnhanceYourCalm = Header().WithCode(420).WithText("Enhance Your Calm").WithType("Client Error")
	// UnprocessableEntity indicates that the request was well-formed but could not be followed due to semantic errors.
	UnprocessableEntity = Header().WithCode(422).WithText("Unprocessable Entity").WithType("Client Error")
	// Locked indicates that the resource being accessed is locked.
	Locked = Header().WithCode(423).WithText("Locked").WithType("Client Error")
	// FailedDependency indicates that the request failed due to failure of a previous request.
	FailedDependency = Header().WithCode(424).WithText("Failed Dependency").WithType("Client Error")
	// UnorderedCollection is a non-standard response code indicating an unordered collection.
	UnorderedCollection = Header().WithCode(425).WithText("Unordered Collection").WithType("Client Error")
	// UpgradeRequired indicates that the client should switch to a different protocol.
	UpgradeRequired = Header().WithCode(426).WithText("Upgrade Required").WithType("Client Error")
	// PreconditionRequired indicates that the origin server requires the request to be conditional.
	PreconditionRequired = Header().WithCode(428).WithText("Precondition Required").WithType("Client Error")
	// TooManyRequests indicates that the user has sent too many requests in a given time.
	TooManyRequests = Header().WithCode(429).WithText("Too Many Requests").WithType("Client Error")
	// RequestHeaderFieldsTooLarge indicates that one or more header fields in the request are too large.
	RequestHeaderFieldsTooLarge = Header().WithCode(431).WithText("Request Header Fields Too Large").WithType("Client Error")
	// NoResponse is a non-standard code indicating that the server has no response to provide.
	NoResponse = Header().WithCode(444).WithText("No Response").WithType("Client Error")
	// RetryWith is a non-standard code indicating that the client should retry with different parameters.
	RetryWith = Header().WithCode(449).WithText("Retry With").WithType("Client Error")
	// BlockedByWindowsParentalControls is a non-standard code indicating that the request was blocked by parental controls.
	BlockedByWindowsParentalControls = Header().WithCode(450).WithText("Blocked by Windows Parental Controls").WithType("Client Error")
	// UnavailableForLegalReasons indicates that the server is denying access to the resource for legal reasons.
	UnavailableForLegalReasons = Header().WithCode(451).WithText("Unavailable For Legal Reasons").WithType("Client Error")
	// ClientClosedRequest is a non-standard code indicating that the client closed the connection before the server's response.
	ClientClosedRequest = Header().WithCode(499).WithText("Client Closed Request").WithType("Client Error")

	// 5xx Server error responses
	// InternalServerError indicates that the server encountered an unexpected condition that prevented it from fulfilling the request.
	InternalServerError = Header().WithCode(500).WithText("Internal Server Error").WithType("Server Error")
	// NotImplemented indicates that the server does not support the functionality required to fulfill the request.
	NotImplemented = Header().WithCode(501).WithText("Not Implemented").WithType("Server Error")
	// BadGateway indicates that the server received an invalid response from an upstream server.
	BadGateway = Header().WithCode(502).WithText("Bad Gateway").WithType("Server Error")
	// ServiceUnavailable indicates that the server is currently unavailable (e.g., overloaded or under maintenance).
	ServiceUnavailable = Header().WithCode(503).WithText("Service Unavailable").WithType("Server Error")
	// GatewayTimeout indicates that the server did not receive a timely response from an upstream server.
	GatewayTimeout = Header().WithCode(504).WithText("Gateway Timeout").WithType("Server Error")
	// HTTPVersionNotSupported indicates that the server does not support the HTTP protocol version used in the request.
	HTTPVersionNotSupported = Header().WithCode(505).WithText("HTTP Version Not Supported").WithType("Server Error")
	// VariantAlsoNegotiates indicates an internal server configuration error leading to circular references.
	VariantAlsoNegotiates = Header().WithCode(506).WithText("Variant Also Negotiates").WithType("Server Error")
	// InsufficientStorage indicates that the server is unable to store the representation needed to complete the request.
	InsufficientStorage = Header().WithCode(507).WithText("Insufficient Storage").WithType("Server Error")
	// LoopDetected indicates that the server detected an infinite loop while processing the request.
	LoopDetected = Header().WithCode(508).WithText("Loop Detected").WithType("Server Error")
	// BandwidthLimitExceeded is a non-standard code indicating that the server's bandwidth limit has been exceeded.
	BandwidthLimitExceeded = Header().WithCode(509).WithText("Bandwidth Limit Exceeded").WithType("Server Error")
	// NotExtended indicates that further extensions to the request are required for the server to fulfill it.
	NotExtended = Header().WithCode(510).WithText("Not Extended").WithType("Server Error")
	// NetworkAuthenticationRequired indicates that the client needs to authenticate to gain network access.
	NetworkAuthenticationRequired = Header().WithCode(511).WithText("Network Authentication Required").WithType("Server Error")
	// NetworkReadTimeoutError is a non-standard code indicating a network read timeout error.
	NetworkReadTimeoutError = Header().WithCode(598).WithText("Network Read Timeout Error").WithType("Server Error")
	// NetworkConnectTimeoutError is a non-standard code indicating a network connection timeout error.
	NetworkConnectTimeoutError = Header().WithCode(599).WithText("Network Connect Timeout Error").WithType("Server Error")
)

// StreamingStrategy defines the strategy used for streaming data.
// It determines how data is sent or received in a streaming manner.
const (
	// Direct streaming without buffering
	// Each piece of data is sent immediately as it becomes available.
	STRATEGY_DIRECT StreamingStrategy = "direct"

	// Buffered streaming with internal buffer
	// Data is collected in a buffer and sent in larger chunks to optimize performance.
	STRATEGY_BUFFERED StreamingStrategy = "buffered"

	// Chunked streaming with explicit chunk handling
	// Data is divided into chunks of a specified size and sent sequentially.
	STRATEGY_CHUNKED StreamingStrategy = "chunked"
)

// CompressionType defines the type of compression applied to data.
// It specifies the algorithm used to compress or decompress data.
const (
	// No compression applied
	CompressNone CompressionType = "none"

	// GZIP compression algorithm
	CompressGzip CompressionType = "gzip"

	// Deflate compression algorithm
	CompressDeflate CompressionType = "deflate"

	// Flate compression algorithm
	CompressFlate CompressionType = "flate"
)

// Common constants used across the package.
const (
	// ErrUnknown is a constant string used to represent an unknown or unspecified value in the context of XC (cross-cutting) concerns.
	// It is typically used as a placeholder when the actual value is not available or not applicable.
	ErrUnknown string = "unknown"

	// defaultChunkSize defines the maximum number of bytes in each chunk.
	// defaultChunkSize is used to limit the size of data chunks when processing large responses or requests.
	defaultChunkSize int = 1024
)

// Locale defines the language and regional settings for content localization.
// It specifies the language and country/region code.
const (
	// English (United States)
	//  Example: en_US
	LocaleEnUS Locale = "en_US"

	// English (United Kingdom)
	//  Example: en_GB
	LocaleEnGB Locale = "en_GB"

	// English (Australia)
	//  Example: en_AU
	LocaleEnAU Locale = "en_AU"

	// English (Canada)
	//  Example: en_CA
	LocaleEnCA Locale = "en_CA"

	// Vietnamese (Vietnam)
	//  Example: vi_VN
	LocaleViVN Locale = "vi_VN"

	// French (France)
	//  Example: fr_FR
	LocaleFrFR Locale = "fr_FR"

	// French (Canada)
	//  Example: fr_CA
	LocaleFrCA Locale = "fr_CA"

	// German (Germany)
	//  Example: de_DE
	LocaleDeDE Locale = "de_DE"

	// Spanish (Spain)
	//  Example: es_ES
	LocaleEsES Locale = "es_ES"

	// Spanish (Mexico)
	//  Example: es_MX
	LocaleEsMX Locale = "es_MX"

	// Portuguese (Brazil)
	//  Example: pt_BR
	LocalePtBR Locale = "pt_BR"

	// Portuguese (Portugal)
	//  Example: pt_PT
	LocalePtPT Locale = "pt_PT"

	// Italian (Italy)
	//  Example: it_IT
	LocaleItIT Locale = "it_IT"

	// Dutch (Netherlands)
	//  Example: nl_NL
	LocaleNlNL Locale = "nl_NL"

	// Russian (Russia)
	//  Example: ru_RU
	LocaleRuRU Locale = "ru_RU"

	// Japanese (Japan)
	//  Example: ja_JP
	LocaleJaJP Locale = "ja_JP"

	// Korean (South Korea)
	//  Example: ko_KR
	LocaleKoKR Locale = "ko_KR"

	// Chinese (Simplified, China)
	//  Example: zh_CN
	LocaleZhCN Locale = "zh_CN"

	// Chinese (Traditional, Taiwan)
	//  Example: zh_TW
	LocaleZhTW Locale = "zh_TW"

	// Thai (Thailand)
	//  Example: th_TH
	LocaleThTH Locale = "th_TH"

	// Indonesian (Indonesia)
	//  Example: id_ID
	LocaleIdID Locale = "id_ID"

	// Malay (Malaysia)
	//  Example: ms_MY
	LocaleMsMY Locale = "ms_MY"

	// Hindi (India)
	//  Example: hi_IN
	LocaleHiIN Locale = "hi_IN"

	// Arabic (Saudi Arabia)
	//  Example: ar_SA
	LocaleArSA Locale = "ar_SA"

	// Turkish (Turkey)
	//  Example: tr_TR
	LocaleTrTR Locale = "tr_TR"

	// Polish (Poland)
	//  Example: pl_PL
	LocalePlPL Locale = "pl_PL"

	// Swedish (Sweden)
	//  Example: sv_SE
	LocaleSvSE Locale = "sv_SE"

	// Danish (Denmark)
	//  Example: da_DK
	LocaleDaDK Locale = "da_DK"

	// Norwegian (Norway)
	//  Example: nb_NO
	LocaleNbNO Locale = "nb_NO"

	// Finnish (Finland)
	//  Example: fi_FI
	LocaleFiFI Locale = "fi_FI"
)
