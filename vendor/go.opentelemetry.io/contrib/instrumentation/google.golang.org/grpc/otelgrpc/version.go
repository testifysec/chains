// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package otelgrpc // import "go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"

// Version is the current release version of the gRPC instrumentation.
func Version() string {
<<<<<<< HEAD
	return "0.57.0"
=======
	return "0.54.0"
>>>>>>> 70e0318b1 ([WIP] add archivista storage backend)
	// This string is updated by the pre_release.sh script during release
}

// SemVersion is the semantic version to be supplied to tracer/meter creation.
//
// Deprecated: Use [Version] instead.
func SemVersion() string {
	return Version()
}
