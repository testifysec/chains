// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package metric // import "go.opentelemetry.io/otel/sdk/metric"

import (
<<<<<<< HEAD
	"runtime"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/metric/exemplar"
	"go.opentelemetry.io/otel/sdk/metric/internal/aggregate"
)

// ExemplarReservoirProviderSelector selects the
// [exemplar.ReservoirProvider] to use
// based on the [Aggregation] of the metric.
type ExemplarReservoirProviderSelector func(Aggregation) exemplar.ReservoirProvider

// reservoirFunc returns the appropriately configured exemplar reservoir
// creation func based on the passed InstrumentKind and filter configuration.
func reservoirFunc[N int64 | float64](provider exemplar.ReservoirProvider, filter exemplar.Filter) func(attribute.Set) aggregate.FilteredExemplarReservoir[N] {
	return func(attrs attribute.Set) aggregate.FilteredExemplarReservoir[N] {
		return aggregate.NewFilteredExemplarReservoir[N](filter, provider(attrs))
	}
}

// DefaultExemplarReservoirProviderSelector returns the default
// [exemplar.ReservoirProvider] for the
// provided [Aggregation].
//
// For explicit bucket histograms with more than 1 bucket, it uses the
// [exemplar.HistogramReservoirProvider].
// For exponential histograms, it uses the
// [exemplar.FixedSizeReservoirProvider]
// with a size of min(20, max_buckets).
// For all other aggregations, it uses the
// [exemplar.FixedSizeReservoirProvider]
// with a size equal to the number of CPUs.
//
// Exemplar default reservoirs MAY change in a minor version bump. No
// guarantees are made on the shape or statistical properties of returned
// exemplars.
func DefaultExemplarReservoirProviderSelector(agg Aggregation) exemplar.ReservoirProvider {
=======
	"os"
	"runtime"
	"slices"

	"go.opentelemetry.io/otel/sdk/metric/internal/exemplar"
	"go.opentelemetry.io/otel/sdk/metric/internal/x"
)

// reservoirFunc returns the appropriately configured exemplar reservoir
// creation func based on the passed InstrumentKind and user defined
// environment variables.
//
// Note: This will only return non-nil values when the experimental exemplar
// feature is enabled and the OTEL_METRICS_EXEMPLAR_FILTER environment variable
// is not set to always_off.
func reservoirFunc[N int64 | float64](agg Aggregation) func() exemplar.FilteredReservoir[N] {
	if !x.Exemplars.Enabled() {
		return nil
	}
	// https://github.com/open-telemetry/opentelemetry-specification/blob/d4b241f451674e8f611bb589477680341006ad2b/specification/configuration/sdk-environment-variables.md#exemplar
	const filterEnvKey = "OTEL_METRICS_EXEMPLAR_FILTER"

	var filter exemplar.Filter

	switch os.Getenv(filterEnvKey) {
	case "always_on":
		filter = exemplar.AlwaysOnFilter
	case "always_off":
		return exemplar.Drop
	case "trace_based":
		fallthrough
	default:
		filter = exemplar.SampledFilter
	}

>>>>>>> 70e0318b1 ([WIP] add archivista storage backend)
	// https://github.com/open-telemetry/opentelemetry-specification/blob/d4b241f451674e8f611bb589477680341006ad2b/specification/metrics/sdk.md#exemplar-defaults
	// Explicit bucket histogram aggregation with more than 1 bucket will
	// use AlignedHistogramBucketExemplarReservoir.
	a, ok := agg.(AggregationExplicitBucketHistogram)
	if ok && len(a.Boundaries) > 0 {
<<<<<<< HEAD
		return exemplar.HistogramReservoirProvider(a.Boundaries)
=======
		cp := slices.Clone(a.Boundaries)
		return func() exemplar.FilteredReservoir[N] {
			bounds := cp
			return exemplar.NewFilteredReservoir[N](filter, exemplar.Histogram(bounds))
		}
>>>>>>> 70e0318b1 ([WIP] add archivista storage backend)
	}

	var n int
	if a, ok := agg.(AggregationBase2ExponentialHistogram); ok {
		// Base2 Exponential Histogram Aggregation SHOULD use a
		// SimpleFixedSizeExemplarReservoir with a reservoir equal to the
		// smaller of the maximum number of buckets configured on the
		// aggregation or twenty (e.g. min(20, max_buckets)).
		n = int(a.MaxSize)
		if n > 20 {
			n = 20
		}
	} else {
		// https://github.com/open-telemetry/opentelemetry-specification/blob/e94af89e3d0c01de30127a0f423e912f6cda7bed/specification/metrics/sdk.md#simplefixedsizeexemplarreservoir
		//   This Exemplar reservoir MAY take a configuration parameter for
		//   the size of the reservoir. If no size configuration is
		//   provided, the default size MAY be the number of possible
		//   concurrent threads (e.g. number of CPUs) to help reduce
		//   contention. Otherwise, a default size of 1 SHOULD be used.
		n = runtime.NumCPU()
		if n < 1 {
			// Should never be the case, but be defensive.
			n = 1
		}
	}

<<<<<<< HEAD
	return exemplar.FixedSizeReservoirProvider(n)
=======
	return func() exemplar.FilteredReservoir[N] {
		return exemplar.NewFilteredReservoir[N](filter, exemplar.FixedSize(n))
	}
>>>>>>> 70e0318b1 ([WIP] add archivista storage backend)
}
