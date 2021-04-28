package http

_emptyBody = bytes(v: "")

// DefaultConfig is the global default for all http requests using the http package.
// Changing this config will affect all other packages using the http package.
// To change the config for a single request pass a new config directly into the corresponding function.
option defaultConfig = {
    // Timeout on the request, if the timeout is zero no timeout is applied
    timeout: 0s,
    // VerifyTLS if false TLS verification will not be performed. This is insecure.
    verifyTLS: true,
}

// Internal method used to perform the actual request
builtin _request : (
    url: string,
    method: string,
    config: {A with timeout: duration, verifyTLS: bool},
    ?headers: [string:string],
    ?body: bytes,
) => {
    statusCode: int,
    body: bytes,
    headers: [string:string],
}

// Make an HTTP request using the provided config
request = (url, method, headers=[:], config=defaultConfig, body=_emptyBody) => _request(
    url: url,
    method: method,
    headers: headers,
    config: config,
    body: body,
)

// Post makes a POST HTTP request
post = (url, headers=[:], config=defaultConfig, body=_emptyBody) => request(
    url: url,
    method: "POST",
    headers: headers,
    config: config,
    body: body,
)

// Get makes a GET HTTP request
get = (url, headers=[:], config=defaultConfig, body=_emptyBody) => request(
    url: url,
    method: "GET",
    headers: headers,
    config: config,
    body: body,
)
