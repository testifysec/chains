// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package gcp // import "go.opentelemetry.io/contrib/detectors/gcp"

// Version is the current release version of the GCP resource detector.
func Version() string {
<<<<<<< HEAD
	return "1.32.0"
=======
	return "1.29.0"
>>>>>>> 70e0318b1 ([WIP] add archivista storage backend)
	// This string is updated by the pre_release.sh script during release
}

// SemVersion is the semantic version to be supplied to tracer/meter creation.
//
// Deprecated: Use [Version] instead.
func SemVersion() string {
	return Version()
}
