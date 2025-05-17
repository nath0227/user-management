package logger

// HTTPPayload is the complete payload that can be interpreted by
// Stackdriver as a HTTP request.
type HTTPPayload struct {
	// The request method. Examples: "GET", "HEAD", "PUT", "POST".
	RequestMethod string `json:"requestMethod"`

	// The scheme (http, https), the host name, the path and the query portion of
	// the URL that was requested.
	//
	// Example: "http://example.com/some/info?color=red".
	RequestURL string `json:"requestUrl"`

	// The size of the HTTP request message in bytes, including the request
	// headers and the request body.
	RequestSize string `json:"requestSize"`

	// The response code indicating the status of response.
	//
	// Examples: 200, 404.
	Status int `json:"status"`

	// The size of the HTTP response message sent back to the client, in bytes,
	// including the response headers and the response body.
	ResponseSize string `json:"responseSize"`

	// The user agent sent by the client.
	//
	// Example: "Mozilla/4.0 (compatible; MSIE 6.0; Windows 98; Q312461; .NET CLR 1.0.3705)".
	UserAgent string `json:"userAgent"`

	// The IP address (IPv4 or IPv6) of the client that issued the HTTP request.
	//
	// Examples: "192.168.1.1", "FE80::0202:B3FF:FE1E:8329".
	RemoteIP string `json:"remoteIp"`

	// The IP address (IPv4 or IPv6) of the origin server that the request was
	// sent to.
	ServerIP string `json:"serverIp"`

	// The referrer URL of the request, as defined in HTTP/1.1 Header Field
	// Definitions.
	Referer string `json:"referer"`

	// The request processing latency on the server, from the time the request was
	// received until the response was sent.
	//
	// A duration in seconds with up to nine fractional digits, terminated by 's'.
	//
	// Example: "3.5s".
	Latency string `json:"latency"`

	// Whether or not a cache lookup was attempted.
	CacheLookup bool `json:"cacheLookup"`

	// Whether or not an entity was served from cache (with or without
	// validation).
	CacheHit bool `json:"cacheHit"`

	// Whether or not the response was validated with the origin server before
	// being served from cache. This field is only meaningful if cacheHit is True.
	CacheValidatedWithOriginServer bool `json:"cacheValidatedWithOriginServer"`

	// The number of HTTP response bytes inserted into cache. Set only when a
	// cache fill was attempted.
	CacheFillBytes string `json:"cacheFillBytes"`

	// Protocol used for the request.
	//
	// Examples: "HTTP/1.1", "HTTP/2", "websocket"
	Protocol string `json:"protocol"`
}
