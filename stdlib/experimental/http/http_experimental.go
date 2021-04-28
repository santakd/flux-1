package http

import (
	"bytes"
	"context"
	"crypto/tls"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/influxdata/flux"
	"github.com/influxdata/flux/codes"
	fhttp "github.com/influxdata/flux/dependencies/http"
	"github.com/influxdata/flux/internal/errors"
	"github.com/influxdata/flux/runtime"
	"github.com/influxdata/flux/semantic"
	"github.com/influxdata/flux/values"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
)

// http get mirrors the http post originally completed for alerts & notifications
var request = values.NewFunction(
	"_request",
	runtime.MustLookupBuiltinType("experimental/http", "_request"),
	func(ctx context.Context, args values.Object) (values.Value, error) {
		// Get and validate URL
		uV, ok := args.Get("url")
		if !ok {
			return nil, errors.New(codes.Invalid, "missing \"url\" parameter")
		}
		u, err := url.Parse(uV.Str())
		if err != nil {
			return nil, err
		}
		deps := flux.GetDependencies(ctx)
		validator, err := deps.URLValidator()
		if err != nil {
			return nil, err
		}
		if err := validator.Validate(u); err != nil {
			return nil, errors.New(codes.Invalid, "no such host")
		}

		methodV, ok := args.Get("method")
		if !ok {
			return nil, errors.New(codes.Invalid, "missing \"method\" parameter")
		}
		if methodV.Type().Nature() != semantic.String {
			return nil, errors.Newf(codes.Invalid, "parameter \"method\" is not of type string: %v", methodV.Type())
		}
		method := methodV.Str()
		switch method {
		case "GET", "POST", "DELETE":
		default:
			return nil, errors.Newf(codes.Invalid, "invalid HTTP method %q", method)
		}

		configV, ok := args.Get("config")
		if !ok {
			return nil, errors.New(codes.Invalid, "missing \"config\" parameter")
		}
		if configV.Type().Nature() != semantic.Object {
			return nil, errors.Newf(codes.Invalid, "parameter \"config\" is not of type record: %v", configV.Type())
		}
		config := configV.Object()

		var body io.Reader
		bodyV, ok := args.Get("body")
		if ok {
			if bodyV.Type().Nature() != semantic.Bytes {
				return nil, errors.Newf(codes.Invalid, "parameter \"body\" is not of type bytes: %v", bodyV.Type())
			}
			body = bytes.NewReader(bodyV.Bytes())
		}

		// Construct HTTP request
		req, err := http.NewRequestWithContext(ctx, method, uV.Str(), body)
		if err != nil {
			return nil, err
		}

		// Add headers to request
		headersV, ok := args.Get("headers")
		if ok && !headersV.IsNull() {
			if headersV.Type().Nature() != semantic.Dictionary {
				return nil, errors.Newf(codes.Invalid, "parameter \"headers\" is not of type [string:string] : %v", headersV.Type())
			}
			var rangeErr error
			headersV.Dict().Range(func(k values.Value, v values.Value) {
				if k.Type().Nature() == semantic.String &&
					v.Type().Nature() == semantic.String {
					req.Header.Set(k.Str(), v.Str())
				} else {
					rangeErr = errors.Newf(codes.Invalid, "header key and values must be a string: %q", k)
				}
			})
			if rangeErr != nil {
				return nil, rangeErr
			}
		}

		// Get Client and configure it
		dc, err := deps.HTTPClient()
		if err != nil {
			return nil, errors.Wrap(err, codes.Aborted, "missing client in http.request")
		}

		timeoutV, ok := config.Get("timeout")
		if !ok {
			return nil, errors.New(codes.Invalid, "config is missing \"timeout\" property")
		}
		timeout := timeoutV.Duration()
		if timeout.IsMixed() {
			return nil, errors.New(codes.Invalid, "config timeout must not be a mixed duration")
		}
		dc, err = fhttp.WithTimeout(dc, timeout.Duration())
		if err != nil {
			return nil, err
		}

		verifyTLSV, ok := config.Get("verifyTLS")
		if !ok {
			return nil, errors.New(codes.Invalid, "config is missing \"verifyTLS\" property")
		}
		if !verifyTLSV.Bool() {
			dc, err = fhttp.WithTLSConfig(dc, &tls.Config{
				InsecureSkipVerify: true,
			})
			if err != nil {
				return nil, err
			}
		}

		// Do request, using local anonymous functions to facilitate timing the request
		statusCode, responseBody, headers, err := func(req *http.Request) (int, []byte, values.Dictionary, error) {
			s, cctx := opentracing.StartSpanFromContext(req.Context(), "http._request")
			s.SetTag("url", req.URL.String())
			defer s.Finish()

			req = req.WithContext(cctx)
			response, err := dc.Do(req)
			if err != nil {
				// Alias the DNS lookup error so as not to disclose the
				// DNS server address. This error is private in the net/http
				// package, so string matching is used.
				if strings.HasSuffix(err.Error(), "no such host") {
					return 0, nil, nil, errors.New(codes.Invalid, "no such host")
				}
				return 0, nil, nil, err
			}
			body, err := ioutil.ReadAll(response.Body)
			_ = response.Body.Close()
			if err != nil {
				return 0, nil, nil, err
			}
			s.LogFields(
				log.Int("statusCode", response.StatusCode),
				log.Int("responseSize", len(body)),
			)
			headers, err := headerToDict(response.Header)
			if err != nil {
				return 0, nil, nil, err
			}
			return response.StatusCode, body, headers, nil
		}(req)
		if err != nil {
			return nil, err
		}

		return values.NewObjectWithValues(map[string]values.Value{
			"statusCode": values.NewInt(int64(statusCode)),
			"headers":    headers,
			"body":       values.NewBytes(responseBody)}), nil

	},
	true, // get has side-effects
)

// headerToDict constructs a values.Dictionary from a map of header keys and values.
func headerToDict(header http.Header) (values.Dictionary, error) {
	builder := values.NewDictBuilder(semantic.NewDictType(semantic.BasicString, semantic.BasicString))
	for name, thevalues := range header {
		for _, onevalue := range thevalues {
			if err := builder.Insert(values.NewString(name), values.NewString(onevalue)); err != nil {
				return nil, err
			}
		}
	}
	return builder.Dict(), nil
}

func init() {
	runtime.RegisterPackageValue("experimental/http", "_request", request)

}
