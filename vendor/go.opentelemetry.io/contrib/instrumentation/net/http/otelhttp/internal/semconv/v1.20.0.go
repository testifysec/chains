<<<<<<< HEAD
// Code created by gotmpl. DO NOT MODIFY.
// source: internal/shared/semconv/v120.0.go.tmpl

=======
>>>>>>> 70e0318b1 ([WIP] add archivista storage backend)
// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package semconv // import "go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp/internal/semconv"

import (
	"errors"
	"io"
	"net/http"
	"slices"
<<<<<<< HEAD
=======
	"strings"
>>>>>>> 70e0318b1 ([WIP] add archivista storage backend)

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp/internal/semconvutil"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/noop"
	semconv "go.opentelemetry.io/otel/semconv/v1.20.0"
)

<<<<<<< HEAD
type OldHTTPServer struct{}
=======
type oldHTTPServer struct{}
>>>>>>> 70e0318b1 ([WIP] add archivista storage backend)

// RequestTraceAttrs returns trace attributes for an HTTP request received by a
// server.
//
// The server must be the primary server name if it is known. For example this
// would be the ServerName directive
// (https://httpd.apache.org/docs/2.4/mod/core.html#servername) for an Apache
// server, and the server_name directive
// (http://nginx.org/en/docs/http/ngx_http_core_module.html#server_name) for an
// nginx server. More generically, the primary server name would be the host
// header value that matches the default virtual host of an HTTP server. It
// should include the host identifier and if a port is used to route to the
// server that port identifier should be included as an appropriate port
// suffix.
//
// If the primary server name is not known, server should be an empty string.
// The req Host will be used to determine the server instead.
<<<<<<< HEAD
func (o OldHTTPServer) RequestTraceAttrs(server string, req *http.Request) []attribute.KeyValue {
=======
func (o oldHTTPServer) RequestTraceAttrs(server string, req *http.Request) []attribute.KeyValue {
>>>>>>> 70e0318b1 ([WIP] add archivista storage backend)
	return semconvutil.HTTPServerRequest(server, req)
}

// ResponseTraceAttrs returns trace attributes for telemetry from an HTTP response.
//
// If any of the fields in the ResponseTelemetry are not set the attribute will be omitted.
<<<<<<< HEAD
func (o OldHTTPServer) ResponseTraceAttrs(resp ResponseTelemetry) []attribute.KeyValue {
=======
func (o oldHTTPServer) ResponseTraceAttrs(resp ResponseTelemetry) []attribute.KeyValue {
>>>>>>> 70e0318b1 ([WIP] add archivista storage backend)
	attributes := []attribute.KeyValue{}

	if resp.ReadBytes > 0 {
		attributes = append(attributes, semconv.HTTPRequestContentLength(int(resp.ReadBytes)))
	}
	if resp.ReadError != nil && !errors.Is(resp.ReadError, io.EOF) {
		// This is not in the semantic conventions, but is historically provided
		attributes = append(attributes, attribute.String("http.read_error", resp.ReadError.Error()))
	}
	if resp.WriteBytes > 0 {
		attributes = append(attributes, semconv.HTTPResponseContentLength(int(resp.WriteBytes)))
	}
	if resp.StatusCode > 0 {
		attributes = append(attributes, semconv.HTTPStatusCode(resp.StatusCode))
	}
	if resp.WriteError != nil && !errors.Is(resp.WriteError, io.EOF) {
		// This is not in the semantic conventions, but is historically provided
		attributes = append(attributes, attribute.String("http.write_error", resp.WriteError.Error()))
	}

	return attributes
}

// Route returns the attribute for the route.
<<<<<<< HEAD
func (o OldHTTPServer) Route(route string) attribute.KeyValue {
=======
func (o oldHTTPServer) Route(route string) attribute.KeyValue {
>>>>>>> 70e0318b1 ([WIP] add archivista storage backend)
	return semconv.HTTPRoute(route)
}

// HTTPStatusCode returns the attribute for the HTTP status code.
// This is a temporary function needed by metrics.  This will be removed when MetricsRequest is added.
func HTTPStatusCode(status int) attribute.KeyValue {
	return semconv.HTTPStatusCode(status)
}

// Server HTTP metrics.
const (
	serverRequestSize  = "http.server.request.size"  // Incoming request bytes total
	serverResponseSize = "http.server.response.size" // Incoming response bytes total
	serverDuration     = "http.server.duration"      // Incoming end to end duration, milliseconds
)

<<<<<<< HEAD
func (h OldHTTPServer) createMeasures(meter metric.Meter) (metric.Int64Counter, metric.Int64Counter, metric.Float64Histogram) {
=======
func (h oldHTTPServer) createMeasures(meter metric.Meter) (metric.Int64Counter, metric.Int64Counter, metric.Float64Histogram) {
>>>>>>> 70e0318b1 ([WIP] add archivista storage backend)
	if meter == nil {
		return noop.Int64Counter{}, noop.Int64Counter{}, noop.Float64Histogram{}
	}
	var err error
	requestBytesCounter, err := meter.Int64Counter(
		serverRequestSize,
		metric.WithUnit("By"),
		metric.WithDescription("Measures the size of HTTP request messages."),
	)
	handleErr(err)

	responseBytesCounter, err := meter.Int64Counter(
		serverResponseSize,
		metric.WithUnit("By"),
		metric.WithDescription("Measures the size of HTTP response messages."),
	)
	handleErr(err)

	serverLatencyMeasure, err := meter.Float64Histogram(
		serverDuration,
		metric.WithUnit("ms"),
		metric.WithDescription("Measures the duration of inbound HTTP requests."),
	)
	handleErr(err)

	return requestBytesCounter, responseBytesCounter, serverLatencyMeasure
}

<<<<<<< HEAD
func (o OldHTTPServer) MetricAttributes(server string, req *http.Request, statusCode int, additionalAttributes []attribute.KeyValue) []attribute.KeyValue {
=======
func (o oldHTTPServer) MetricAttributes(server string, req *http.Request, statusCode int, additionalAttributes []attribute.KeyValue) []attribute.KeyValue {
>>>>>>> 70e0318b1 ([WIP] add archivista storage backend)
	n := len(additionalAttributes) + 3
	var host string
	var p int
	if server == "" {
<<<<<<< HEAD
		host, p = SplitHostPort(req.Host)
	} else {
		// Prioritize the primary server name.
		host, p = SplitHostPort(server)
		if p < 0 {
			_, p = SplitHostPort(req.Host)
=======
		host, p = splitHostPort(req.Host)
	} else {
		// Prioritize the primary server name.
		host, p = splitHostPort(server)
		if p < 0 {
			_, p = splitHostPort(req.Host)
>>>>>>> 70e0318b1 ([WIP] add archivista storage backend)
		}
	}
	hostPort := requiredHTTPPort(req.TLS != nil, p)
	if hostPort > 0 {
		n++
	}
	protoName, protoVersion := netProtocol(req.Proto)
	if protoName != "" {
		n++
	}
	if protoVersion != "" {
		n++
	}

	if statusCode > 0 {
		n++
	}

	attributes := slices.Grow(additionalAttributes, n)
	attributes = append(attributes,
<<<<<<< HEAD
		semconv.HTTPMethod(standardizeHTTPMethod(req.Method)),
=======
		o.methodMetric(req.Method),
>>>>>>> 70e0318b1 ([WIP] add archivista storage backend)
		o.scheme(req.TLS != nil),
		semconv.NetHostName(host))

	if hostPort > 0 {
		attributes = append(attributes, semconv.NetHostPort(hostPort))
	}
	if protoName != "" {
		attributes = append(attributes, semconv.NetProtocolName(protoName))
	}
	if protoVersion != "" {
		attributes = append(attributes, semconv.NetProtocolVersion(protoVersion))
	}

	if statusCode > 0 {
		attributes = append(attributes, semconv.HTTPStatusCode(statusCode))
	}
	return attributes
}

<<<<<<< HEAD
func (o OldHTTPServer) scheme(https bool) attribute.KeyValue { // nolint:revive
=======
func (o oldHTTPServer) methodMetric(method string) attribute.KeyValue {
	method = strings.ToUpper(method)
	switch method {
	case http.MethodConnect, http.MethodDelete, http.MethodGet, http.MethodHead, http.MethodOptions, http.MethodPatch, http.MethodPost, http.MethodPut, http.MethodTrace:
	default:
		method = "_OTHER"
	}
	return semconv.HTTPMethod(method)
}

func (o oldHTTPServer) scheme(https bool) attribute.KeyValue { // nolint:revive
>>>>>>> 70e0318b1 ([WIP] add archivista storage backend)
	if https {
		return semconv.HTTPSchemeHTTPS
	}
	return semconv.HTTPSchemeHTTP
}

<<<<<<< HEAD
type OldHTTPClient struct{}

func (o OldHTTPClient) RequestTraceAttrs(req *http.Request) []attribute.KeyValue {
	return semconvutil.HTTPClientRequest(req)
}

func (o OldHTTPClient) ResponseTraceAttrs(resp *http.Response) []attribute.KeyValue {
	return semconvutil.HTTPClientResponse(resp)
}

func (o OldHTTPClient) MetricAttributes(req *http.Request, statusCode int, additionalAttributes []attribute.KeyValue) []attribute.KeyValue {
	/* The following semantic conventions are returned if present:
	http.method                     string
	http.status_code             int
	net.peer.name                   string
	net.peer.port                   int
	*/

	n := 2 // method, peer name.
	var h string
	if req.URL != nil {
		h = req.URL.Host
	}
	var requestHost string
	var requestPort int
	for _, hostport := range []string{h, req.Header.Get("Host")} {
		requestHost, requestPort = SplitHostPort(hostport)
		if requestHost != "" || requestPort > 0 {
			break
		}
	}

	port := requiredHTTPPort(req.URL != nil && req.URL.Scheme == "https", requestPort)
	if port > 0 {
		n++
	}

	if statusCode > 0 {
		n++
	}

	attributes := slices.Grow(additionalAttributes, n)
	attributes = append(attributes,
		semconv.HTTPMethod(standardizeHTTPMethod(req.Method)),
		semconv.NetPeerName(requestHost),
	)

	if port > 0 {
		attributes = append(attributes, semconv.NetPeerPort(port))
	}

	if statusCode > 0 {
		attributes = append(attributes, semconv.HTTPStatusCode(statusCode))
	}
	return attributes
}

// Client HTTP metrics.
const (
	clientRequestSize  = "http.client.request.size"  // Incoming request bytes total
	clientResponseSize = "http.client.response.size" // Incoming response bytes total
	clientDuration     = "http.client.duration"      // Incoming end to end duration, milliseconds
)

func (o OldHTTPClient) createMeasures(meter metric.Meter) (metric.Int64Counter, metric.Int64Counter, metric.Float64Histogram) {
	if meter == nil {
		return noop.Int64Counter{}, noop.Int64Counter{}, noop.Float64Histogram{}
	}
	requestBytesCounter, err := meter.Int64Counter(
		clientRequestSize,
		metric.WithUnit("By"),
		metric.WithDescription("Measures the size of HTTP request messages."),
	)
	handleErr(err)

	responseBytesCounter, err := meter.Int64Counter(
		clientResponseSize,
		metric.WithUnit("By"),
		metric.WithDescription("Measures the size of HTTP response messages."),
	)
	handleErr(err)

	latencyMeasure, err := meter.Float64Histogram(
		clientDuration,
		metric.WithUnit("ms"),
		metric.WithDescription("Measures the duration of outbound HTTP requests."),
	)
	handleErr(err)

	return requestBytesCounter, responseBytesCounter, latencyMeasure
}
=======
type oldHTTPClient struct{}

func (o oldHTTPClient) RequestTraceAttrs(req *http.Request) []attribute.KeyValue {
	return semconvutil.HTTPClientRequest(req)
}

func (o oldHTTPClient) ResponseTraceAttrs(resp *http.Response) []attribute.KeyValue {
	return semconvutil.HTTPClientResponse(resp)
}
>>>>>>> 70e0318b1 ([WIP] add archivista storage backend)
