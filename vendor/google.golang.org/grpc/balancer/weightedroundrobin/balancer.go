/*
 *
 * Copyright 2023 gRPC authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

package weightedroundrobin

import (
<<<<<<< HEAD
	"encoding/json"
	"fmt"
	rand "math/rand/v2"
=======
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
>>>>>>> 70e0318b1 ([WIP] add archivista storage backend)
	"sync"
	"sync/atomic"
	"time"
	"unsafe"

	"google.golang.org/grpc/balancer"
<<<<<<< HEAD
	"google.golang.org/grpc/balancer/endpointsharding"
	"google.golang.org/grpc/balancer/pickfirst/pickfirstleaf"
=======
	"google.golang.org/grpc/balancer/base"
>>>>>>> 70e0318b1 ([WIP] add archivista storage backend)
	"google.golang.org/grpc/balancer/weightedroundrobin/internal"
	"google.golang.org/grpc/balancer/weightedtarget"
	"google.golang.org/grpc/connectivity"
	estats "google.golang.org/grpc/experimental/stats"
	"google.golang.org/grpc/internal/grpclog"
<<<<<<< HEAD
	"google.golang.org/grpc/internal/grpcsync"
=======
>>>>>>> 70e0318b1 ([WIP] add archivista storage backend)
	iserviceconfig "google.golang.org/grpc/internal/serviceconfig"
	"google.golang.org/grpc/orca"
	"google.golang.org/grpc/resolver"
	"google.golang.org/grpc/serviceconfig"

	v3orcapb "github.com/cncf/xds/go/xds/data/orca/v3"
)

// Name is the name of the weighted round robin balancer.
const Name = "weighted_round_robin"

var (
	rrFallbackMetric = estats.RegisterInt64Count(estats.MetricDescriptor{
		Name:           "grpc.lb.wrr.rr_fallback",
		Description:    "EXPERIMENTAL. Number of scheduler updates in which there were not enough endpoints with valid weight, which caused the WRR policy to fall back to RR behavior.",
		Unit:           "update",
		Labels:         []string{"grpc.target"},
		OptionalLabels: []string{"grpc.lb.locality"},
		Default:        false,
	})

	endpointWeightNotYetUsableMetric = estats.RegisterInt64Count(estats.MetricDescriptor{
		Name:           "grpc.lb.wrr.endpoint_weight_not_yet_usable",
		Description:    "EXPERIMENTAL. Number of endpoints from each scheduler update that don't yet have usable weight information (i.e., either the load report has not yet been received, or it is within the blackout period).",
		Unit:           "endpoint",
		Labels:         []string{"grpc.target"},
		OptionalLabels: []string{"grpc.lb.locality"},
		Default:        false,
	})

	endpointWeightStaleMetric = estats.RegisterInt64Count(estats.MetricDescriptor{
		Name:           "grpc.lb.wrr.endpoint_weight_stale",
		Description:    "EXPERIMENTAL. Number of endpoints from each scheduler update whose latest weight is older than the expiration period.",
		Unit:           "endpoint",
		Labels:         []string{"grpc.target"},
		OptionalLabels: []string{"grpc.lb.locality"},
		Default:        false,
	})
	endpointWeightsMetric = estats.RegisterFloat64Histo(estats.MetricDescriptor{
		Name:           "grpc.lb.wrr.endpoint_weights",
		Description:    "EXPERIMENTAL. Weight of each endpoint, recorded on every scheduler update. Endpoints without usable weights will be recorded as weight 0.",
		Unit:           "endpoint",
		Labels:         []string{"grpc.target"},
		OptionalLabels: []string{"grpc.lb.locality"},
		Default:        false,
	})
)

<<<<<<< HEAD
// endpointSharding which specifies pick first children.
var endpointShardingLBConfig serviceconfig.LoadBalancingConfig

func init() {
	balancer.Register(bb{})
	var err error
	endpointShardingLBConfig, err = endpointsharding.ParseConfig(json.RawMessage(endpointsharding.PickFirstConfig))
	if err != nil {
		logger.Fatal(err)
	}
=======
func init() {
	balancer.Register(bb{})
>>>>>>> 70e0318b1 ([WIP] add archivista storage backend)
}

type bb struct{}

func (bb) Build(cc balancer.ClientConn, bOpts balancer.BuildOptions) balancer.Balancer {
	b := &wrrBalancer{
<<<<<<< HEAD
		ClientConn:       cc,
		target:           bOpts.Target.String(),
		metricsRecorder:  bOpts.MetricsRecorder,
		addressWeights:   resolver.NewAddressMap(),
		endpointToWeight: resolver.NewEndpointMap(),
		scToWeight:       make(map[balancer.SubConn]*endpointWeight),
	}

	b.child = endpointsharding.NewBalancer(b, bOpts)
=======
		cc:                cc,
		subConns:          resolver.NewAddressMap(),
		csEvltr:           &balancer.ConnectivityStateEvaluator{},
		scMap:             make(map[balancer.SubConn]*weightedSubConn),
		connectivityState: connectivity.Connecting,
		target:            bOpts.Target.String(),
		metricsRecorder:   bOpts.MetricsRecorder,
	}

>>>>>>> 70e0318b1 ([WIP] add archivista storage backend)
	b.logger = prefixLogger(b)
	b.logger.Infof("Created")
	return b
}

func (bb) ParseConfig(js json.RawMessage) (serviceconfig.LoadBalancingConfig, error) {
	lbCfg := &lbConfig{
		// Default values as documented in A58.
		OOBReportingPeriod:      iserviceconfig.Duration(10 * time.Second),
		BlackoutPeriod:          iserviceconfig.Duration(10 * time.Second),
		WeightExpirationPeriod:  iserviceconfig.Duration(3 * time.Minute),
		WeightUpdatePeriod:      iserviceconfig.Duration(time.Second),
		ErrorUtilizationPenalty: 1,
	}
	if err := json.Unmarshal(js, lbCfg); err != nil {
		return nil, fmt.Errorf("wrr: unable to unmarshal LB policy config: %s, error: %v", string(js), err)
	}

	if lbCfg.ErrorUtilizationPenalty < 0 {
		return nil, fmt.Errorf("wrr: errorUtilizationPenalty must be non-negative")
	}

	// For easier comparisons later, ensure the OOB reporting period is unset
	// (0s) when OOB reports are disabled.
	if !lbCfg.EnableOOBLoadReport {
		lbCfg.OOBReportingPeriod = 0
	}

	// Impose lower bound of 100ms on weightUpdatePeriod.
	if !internal.AllowAnyWeightUpdatePeriod && lbCfg.WeightUpdatePeriod < iserviceconfig.Duration(100*time.Millisecond) {
		lbCfg.WeightUpdatePeriod = iserviceconfig.Duration(100 * time.Millisecond)
	}

	return lbCfg, nil
}

func (bb) Name() string {
	return Name
}

<<<<<<< HEAD
// updateEndpointsLocked updates endpoint weight state based off new update, by
// starting and clearing any endpoint weights needed.
//
// Caller must hold b.mu.
func (b *wrrBalancer) updateEndpointsLocked(endpoints []resolver.Endpoint) {
	endpointSet := resolver.NewEndpointMap()
	addressSet := resolver.NewAddressMap()
	for _, endpoint := range endpoints {
		endpointSet.Set(endpoint, nil)
		for _, addr := range endpoint.Addresses {
			addressSet.Set(addr, nil)
		}
		var ew *endpointWeight
		if ewi, ok := b.endpointToWeight.Get(endpoint); ok {
			ew = ewi.(*endpointWeight)
		} else {
			ew = &endpointWeight{
				logger:            b.logger,
				connectivityState: connectivity.Connecting,
				// Initially, we set load reports to off, because they are not
				// running upon initial endpointWeight creation.
				cfg:             &lbConfig{EnableOOBLoadReport: false},
				metricsRecorder: b.metricsRecorder,
				target:          b.target,
				locality:        b.locality,
			}
			for _, addr := range endpoint.Addresses {
				b.addressWeights.Set(addr, ew)
			}
			b.endpointToWeight.Set(endpoint, ew)
		}
		ew.updateConfig(b.cfg)
	}

	for _, endpoint := range b.endpointToWeight.Keys() {
		if _, ok := endpointSet.Get(endpoint); ok {
			// Existing endpoint also in new endpoint list; skip.
			continue
		}
		b.endpointToWeight.Delete(endpoint)
		for _, addr := range endpoint.Addresses {
			if _, ok := addressSet.Get(addr); !ok { // old endpoints to be deleted can share addresses with new endpoints, so only delete if necessary
				b.addressWeights.Delete(addr)
			}
		}
		// SubConn map will get handled in updateSubConnState
		// when receives SHUTDOWN signal.
	}
}

// wrrBalancer implements the weighted round robin LB policy.
type wrrBalancer struct {
	// The following fields are set at initialization time and read only after that,
	// so they do not need to be protected by a mutex.
	child               balancer.Balancer
	balancer.ClientConn // Embed to intercept NewSubConn operation
	logger              *grpclog.PrefixLogger
	target              string
	metricsRecorder     estats.MetricsRecorder

	mu               sync.Mutex
	cfg              *lbConfig // active config
	locality         string
	stopPicker       *grpcsync.Event
	addressWeights   *resolver.AddressMap  // addr -> endpointWeight
	endpointToWeight *resolver.EndpointMap // endpoint -> endpointWeight
	scToWeight       map[balancer.SubConn]*endpointWeight
}

func (b *wrrBalancer) UpdateClientConnState(ccs balancer.ClientConnState) error {
	if b.logger.V(2) {
		b.logger.Infof("UpdateCCS: %v", ccs)
	}
=======
// wrrBalancer implements the weighted round robin LB policy.
type wrrBalancer struct {
	// The following fields are immutable.
	cc              balancer.ClientConn
	logger          *grpclog.PrefixLogger
	target          string
	metricsRecorder estats.MetricsRecorder

	// The following fields are only accessed on calls into the LB policy, and
	// do not need a mutex.
	cfg               *lbConfig            // active config
	subConns          *resolver.AddressMap // active weightedSubConns mapped by address
	scMap             map[balancer.SubConn]*weightedSubConn
	connectivityState connectivity.State // aggregate state
	csEvltr           *balancer.ConnectivityStateEvaluator
	resolverErr       error // the last error reported by the resolver; cleared on successful resolution
	connErr           error // the last connection error; cleared upon leaving TransientFailure
	stopPicker        func()
	locality          string
}

func (b *wrrBalancer) UpdateClientConnState(ccs balancer.ClientConnState) error {
	b.logger.Infof("UpdateCCS: %v", ccs)
	b.resolverErr = nil
>>>>>>> 70e0318b1 ([WIP] add archivista storage backend)
	cfg, ok := ccs.BalancerConfig.(*lbConfig)
	if !ok {
		return fmt.Errorf("wrr: received nil or illegal BalancerConfig (type %T): %v", ccs.BalancerConfig, ccs.BalancerConfig)
	}

<<<<<<< HEAD
	// Note: empty endpoints and duplicate addresses across endpoints won't
	// explicitly error but will have undefined behavior.
	b.mu.Lock()
	b.cfg = cfg
	b.locality = weightedtarget.LocalityFromResolverState(ccs.ResolverState)
	b.updateEndpointsLocked(ccs.ResolverState.Endpoints)
	b.mu.Unlock()

	// Make pickfirst children use health listeners for outlier detection to
	// work.
	ccs.ResolverState = pickfirstleaf.EnableHealthListener(ccs.ResolverState)
	// This causes child to update picker inline and will thus cause inline
	// picker update.
	return b.child.UpdateClientConnState(balancer.ClientConnState{
		BalancerConfig: endpointShardingLBConfig,
		ResolverState:  ccs.ResolverState,
	})
}

func (b *wrrBalancer) UpdateState(state balancer.State) {
	b.mu.Lock()
	defer b.mu.Unlock()

	if b.stopPicker != nil {
		b.stopPicker.Fire()
		b.stopPicker = nil
	}

	childStates := endpointsharding.ChildStatesFromPicker(state.Picker)

	var readyPickersWeight []pickerWeightedEndpoint

	for _, childState := range childStates {
		if childState.State.ConnectivityState == connectivity.Ready {
			ewv, ok := b.endpointToWeight.Get(childState.Endpoint)
			if !ok {
				// Should never happen, simply continue and ignore this endpoint
				// for READY pickers.
				continue
			}
			ew := ewv.(*endpointWeight)
			readyPickersWeight = append(readyPickersWeight, pickerWeightedEndpoint{
				picker:           childState.State.Picker,
				weightedEndpoint: ew,
			})
		}
	}
	// If no ready pickers are present, simply defer to the round robin picker
	// from endpoint sharding, which will round robin across the most relevant
	// pick first children in the highest precedence connectivity state.
	if len(readyPickersWeight) == 0 {
		b.ClientConn.UpdateState(balancer.State{
			ConnectivityState: state.ConnectivityState,
			Picker:            state.Picker,
		})
		return
	}

	p := &picker{
		v:               rand.Uint32(), // start the scheduler at a random point
		cfg:             b.cfg,
		weightedPickers: readyPickersWeight,
		metricsRecorder: b.metricsRecorder,
		locality:        b.locality,
		target:          b.target,
	}

	b.stopPicker = grpcsync.NewEvent()
	p.start(b.stopPicker)

	b.ClientConn.UpdateState(balancer.State{
		ConnectivityState: state.ConnectivityState,
		Picker:            p,
	})
}

type pickerWeightedEndpoint struct {
	picker           balancer.Picker
	weightedEndpoint *endpointWeight
}

func (b *wrrBalancer) NewSubConn(addrs []resolver.Address, opts balancer.NewSubConnOptions) (balancer.SubConn, error) {
	addr := addrs[0] // The new pick first policy for DualStack will only ever create a SubConn with one address.
	var sc balancer.SubConn

	oldListener := opts.StateListener
	opts.StateListener = func(state balancer.SubConnState) {
		b.updateSubConnState(sc, state)
		oldListener(state)
	}

	b.mu.Lock()
	defer b.mu.Unlock()
	ewi, ok := b.addressWeights.Get(addr)
	if !ok {
		// SubConn state updates can come in for a no longer relevant endpoint
		// weight (from the old system after a new config update is applied).
		return nil, fmt.Errorf("balancer is being closed; no new SubConns allowed")
	}
	sc, err := b.ClientConn.NewSubConn([]resolver.Address{addr}, opts)
	if err != nil {
		return nil, err
	}
	b.scToWeight[sc] = ewi.(*endpointWeight)
	return sc, nil
}

func (b *wrrBalancer) ResolverError(err error) {
	// Will cause inline picker update from endpoint sharding.
	b.child.ResolverError(err)
=======
	b.cfg = cfg
	b.locality = weightedtarget.LocalityFromResolverState(ccs.ResolverState)
	b.updateAddresses(ccs.ResolverState.Addresses)

	if len(ccs.ResolverState.Addresses) == 0 {
		b.ResolverError(errors.New("resolver produced zero addresses")) // will call regeneratePicker
		return balancer.ErrBadResolverState
	}

	b.regeneratePicker()

	return nil
}

func (b *wrrBalancer) updateAddresses(addrs []resolver.Address) {
	addrsSet := resolver.NewAddressMap()

	// Loop through new address list and create subconns for any new addresses.
	for _, addr := range addrs {
		if _, ok := addrsSet.Get(addr); ok {
			// Redundant address; skip.
			continue
		}
		addrsSet.Set(addr, nil)

		var wsc *weightedSubConn
		wsci, ok := b.subConns.Get(addr)
		if ok {
			wsc = wsci.(*weightedSubConn)
		} else {
			// addr is a new address (not existing in b.subConns).
			var sc balancer.SubConn
			sc, err := b.cc.NewSubConn([]resolver.Address{addr}, balancer.NewSubConnOptions{
				StateListener: func(state balancer.SubConnState) {
					b.updateSubConnState(sc, state)
				},
			})
			if err != nil {
				b.logger.Warningf("Failed to create new SubConn for address %v: %v", addr, err)
				continue
			}
			wsc = &weightedSubConn{
				SubConn:           sc,
				logger:            b.logger,
				connectivityState: connectivity.Idle,
				// Initially, we set load reports to off, because they are not
				// running upon initial weightedSubConn creation.
				cfg: &lbConfig{EnableOOBLoadReport: false},

				metricsRecorder: b.metricsRecorder,
				target:          b.target,
				locality:        b.locality,
			}
			b.subConns.Set(addr, wsc)
			b.scMap[sc] = wsc
			b.csEvltr.RecordTransition(connectivity.Shutdown, connectivity.Idle)
			sc.Connect()
		}
		// Update config for existing weightedSubConn or send update for first
		// time to new one.  Ensures an OOB listener is running if needed
		// (and stops the existing one if applicable).
		wsc.updateConfig(b.cfg)
	}

	// Loop through existing subconns and remove ones that are not in addrs.
	for _, addr := range b.subConns.Keys() {
		if _, ok := addrsSet.Get(addr); ok {
			// Existing address also in new address list; skip.
			continue
		}
		// addr was removed by resolver.  Remove.
		wsci, _ := b.subConns.Get(addr)
		wsc := wsci.(*weightedSubConn)
		wsc.SubConn.Shutdown()
		b.subConns.Delete(addr)
	}
}

func (b *wrrBalancer) ResolverError(err error) {
	b.resolverErr = err
	if b.subConns.Len() == 0 {
		b.connectivityState = connectivity.TransientFailure
	}
	if b.connectivityState != connectivity.TransientFailure {
		// No need to update the picker since no error is being returned.
		return
	}
	b.regeneratePicker()
>>>>>>> 70e0318b1 ([WIP] add archivista storage backend)
}

func (b *wrrBalancer) UpdateSubConnState(sc balancer.SubConn, state balancer.SubConnState) {
	b.logger.Errorf("UpdateSubConnState(%v, %+v) called unexpectedly", sc, state)
}

func (b *wrrBalancer) updateSubConnState(sc balancer.SubConn, state balancer.SubConnState) {
<<<<<<< HEAD
	b.mu.Lock()
	ew := b.scToWeight[sc]
	// updates from a no longer relevant SubConn update, nothing to do here but
	// forward state to state listener, which happens in wrapped listener. Will
	// eventually get cleared from scMap once receives Shutdown signal.
	if ew == nil {
		b.mu.Unlock()
		return
	}
	if state.ConnectivityState == connectivity.Shutdown {
		delete(b.scToWeight, sc)
	}
	b.mu.Unlock()

	// On the first READY SubConn/Transition for an endpoint, set pickedSC,
	// clear endpoint tracking weight state, and potentially start an OOB watch.
	if state.ConnectivityState == connectivity.Ready && ew.pickedSC == nil {
		ew.pickedSC = sc
		ew.mu.Lock()
		ew.nonEmptySince = time.Time{}
		ew.lastUpdated = time.Time{}
		cfg := ew.cfg
		ew.mu.Unlock()
		ew.updateORCAListener(cfg)
		return
	}

	// If the pickedSC (the one pick first uses for an endpoint) transitions out
	// of READY, stop OOB listener if needed and clear pickedSC so the next
	// created SubConn for the endpoint that goes READY will be chosen for
	// endpoint as the active SubConn.
	if state.ConnectivityState != connectivity.Ready && ew.pickedSC == sc {
		// The first SubConn that goes READY for an endpoint is what pick first
		// will pick. Only once that SubConn goes not ready will pick first
		// restart this cycle of creating SubConns and using the first READY
		// one. The lower level endpoint sharding will ping the Pick First once
		// this occurs to ExitIdle which will trigger a connection attempt.
		if ew.stopORCAListener != nil {
			ew.stopORCAListener()
		}
		ew.pickedSC = nil
=======
	wsc := b.scMap[sc]
	if wsc == nil {
		b.logger.Errorf("UpdateSubConnState called with an unknown SubConn: %p, %v", sc, state)
		return
	}
	if b.logger.V(2) {
		logger.Infof("UpdateSubConnState(%+v, %+v)", sc, state)
	}

	cs := state.ConnectivityState

	if cs == connectivity.TransientFailure {
		// Save error to be reported via picker.
		b.connErr = state.ConnectionError
	}

	if cs == connectivity.Shutdown {
		delete(b.scMap, sc)
		// The subconn was removed from b.subConns when the address was removed
		// in updateAddresses.
	}

	oldCS := wsc.updateConnectivityState(cs)
	b.connectivityState = b.csEvltr.RecordTransition(oldCS, cs)

	// Regenerate picker when one of the following happens:
	//  - this sc entered or left ready
	//  - the aggregated state of balancer is TransientFailure
	//    (may need to update error message)
	if (cs == connectivity.Ready) != (oldCS == connectivity.Ready) ||
		b.connectivityState == connectivity.TransientFailure {
		b.regeneratePicker()
>>>>>>> 70e0318b1 ([WIP] add archivista storage backend)
	}
}

// Close stops the balancer.  It cancels any ongoing scheduler updates and
// stops any ORCA listeners.
func (b *wrrBalancer) Close() {
<<<<<<< HEAD
	b.mu.Lock()
	if b.stopPicker != nil {
		b.stopPicker.Fire()
		b.stopPicker = nil
	}
	b.mu.Unlock()

	// Ensure any lingering OOB watchers are stopped.
	for _, ewv := range b.endpointToWeight.Values() {
		ew := ewv.(*endpointWeight)
		if ew.stopORCAListener != nil {
			ew.stopORCAListener()
		}
	}
}

func (b *wrrBalancer) ExitIdle() {
	if ei, ok := b.child.(balancer.ExitIdler); ok { // Should always be ok, as child is endpoint sharding.
		ei.ExitIdle()
	}
=======
	if b.stopPicker != nil {
		b.stopPicker()
		b.stopPicker = nil
	}
	for _, wsc := range b.scMap {
		// Ensure any lingering OOB watchers are stopped.
		wsc.updateConnectivityState(connectivity.Shutdown)
	}
}

// ExitIdle is ignored; we always connect to all backends.
func (b *wrrBalancer) ExitIdle() {}

func (b *wrrBalancer) readySubConns() []*weightedSubConn {
	var ret []*weightedSubConn
	for _, v := range b.subConns.Values() {
		wsc := v.(*weightedSubConn)
		if wsc.connectivityState == connectivity.Ready {
			ret = append(ret, wsc)
		}
	}
	return ret
}

// mergeErrors builds an error from the last connection error and the last
// resolver error.  Must only be called if b.connectivityState is
// TransientFailure.
func (b *wrrBalancer) mergeErrors() error {
	// connErr must always be non-nil unless there are no SubConns, in which
	// case resolverErr must be non-nil.
	if b.connErr == nil {
		return fmt.Errorf("last resolver error: %v", b.resolverErr)
	}
	if b.resolverErr == nil {
		return fmt.Errorf("last connection error: %v", b.connErr)
	}
	return fmt.Errorf("last connection error: %v; last resolver error: %v", b.connErr, b.resolverErr)
}

func (b *wrrBalancer) regeneratePicker() {
	if b.stopPicker != nil {
		b.stopPicker()
		b.stopPicker = nil
	}

	switch b.connectivityState {
	case connectivity.TransientFailure:
		b.cc.UpdateState(balancer.State{
			ConnectivityState: connectivity.TransientFailure,
			Picker:            base.NewErrPicker(b.mergeErrors()),
		})
		return
	case connectivity.Connecting, connectivity.Idle:
		// Idle could happen very briefly if all subconns are Idle and we've
		// asked them to connect but they haven't reported Connecting yet.
		// Report the same as Connecting since this is temporary.
		b.cc.UpdateState(balancer.State{
			ConnectivityState: connectivity.Connecting,
			Picker:            base.NewErrPicker(balancer.ErrNoSubConnAvailable),
		})
		return
	case connectivity.Ready:
		b.connErr = nil
	}

	p := &picker{
		v:               rand.Uint32(), // start the scheduler at a random point
		cfg:             b.cfg,
		subConns:        b.readySubConns(),
		metricsRecorder: b.metricsRecorder,
		locality:        b.locality,
		target:          b.target,
	}
	var ctx context.Context
	ctx, b.stopPicker = context.WithCancel(context.Background())
	p.start(ctx)
	b.cc.UpdateState(balancer.State{
		ConnectivityState: b.connectivityState,
		Picker:            p,
	})
>>>>>>> 70e0318b1 ([WIP] add archivista storage backend)
}

// picker is the WRR policy's picker.  It uses live-updating backend weights to
// update the scheduler periodically and ensure picks are routed proportional
// to those weights.
type picker struct {
<<<<<<< HEAD
	scheduler unsafe.Pointer // *scheduler; accessed atomically
	v         uint32         // incrementing value used by the scheduler; accessed atomically
	cfg       *lbConfig      // active config when picker created

	weightedPickers []pickerWeightedEndpoint // all READY pickers
=======
	scheduler unsafe.Pointer     // *scheduler; accessed atomically
	v         uint32             // incrementing value used by the scheduler; accessed atomically
	cfg       *lbConfig          // active config when picker created
	subConns  []*weightedSubConn // all READY subconns
>>>>>>> 70e0318b1 ([WIP] add archivista storage backend)

	// The following fields are immutable.
	target          string
	locality        string
	metricsRecorder estats.MetricsRecorder
}

<<<<<<< HEAD
func (p *picker) endpointWeights(recordMetrics bool) []float64 {
	wp := make([]float64, len(p.weightedPickers))
	now := internal.TimeNow()
	for i, wpi := range p.weightedPickers {
		wp[i] = wpi.weightedEndpoint.weight(now, time.Duration(p.cfg.WeightExpirationPeriod), time.Duration(p.cfg.BlackoutPeriod), recordMetrics)
	}
	return wp
}

func (p *picker) Pick(info balancer.PickInfo) (balancer.PickResult, error) {
	// Read the scheduler atomically.  All scheduler operations are threadsafe,
	// and if the scheduler is replaced during this usage, we want to use the
	// scheduler that was live when the pick started.
	sched := *(*scheduler)(atomic.LoadPointer(&p.scheduler))

	pickedPicker := p.weightedPickers[sched.nextIndex()]
	pr, err := pickedPicker.picker.Pick(info)
	if err != nil {
		logger.Errorf("ready picker returned error: %v", err)
		return balancer.PickResult{}, err
	}
	if !p.cfg.EnableOOBLoadReport {
		oldDone := pr.Done
		pr.Done = func(info balancer.DoneInfo) {
			if load, ok := info.ServerLoad.(*v3orcapb.OrcaLoadReport); ok && load != nil {
				pickedPicker.weightedEndpoint.OnLoadReport(load)
			}
			if oldDone != nil {
				oldDone(info)
			}
		}
	}
	return pr, nil
=======
func (p *picker) scWeights(recordMetrics bool) []float64 {
	ws := make([]float64, len(p.subConns))
	now := internal.TimeNow()
	for i, wsc := range p.subConns {
		ws[i] = wsc.weight(now, time.Duration(p.cfg.WeightExpirationPeriod), time.Duration(p.cfg.BlackoutPeriod), recordMetrics)
	}

	return ws
>>>>>>> 70e0318b1 ([WIP] add archivista storage backend)
}

func (p *picker) inc() uint32 {
	return atomic.AddUint32(&p.v, 1)
}

func (p *picker) regenerateScheduler() {
	s := p.newScheduler(true)
	atomic.StorePointer(&p.scheduler, unsafe.Pointer(&s))
}

<<<<<<< HEAD
func (p *picker) start(stopPicker *grpcsync.Event) {
	p.regenerateScheduler()
	if len(p.weightedPickers) == 1 {
=======
func (p *picker) start(ctx context.Context) {
	p.regenerateScheduler()
	if len(p.subConns) == 1 {
>>>>>>> 70e0318b1 ([WIP] add archivista storage backend)
		// No need to regenerate weights with only one backend.
		return
	}

	go func() {
		ticker := time.NewTicker(time.Duration(p.cfg.WeightUpdatePeriod))
		defer ticker.Stop()
		for {
			select {
<<<<<<< HEAD
			case <-stopPicker.Done():
=======
			case <-ctx.Done():
>>>>>>> 70e0318b1 ([WIP] add archivista storage backend)
				return
			case <-ticker.C:
				p.regenerateScheduler()
			}
		}
	}()
}

<<<<<<< HEAD
// endpointWeight is the weight for an endpoint. It tracks the SubConn that will
// be picked for the endpoint, and other parameters relevant to computing the
// effective weight. When needed, it also tracks connectivity state, listens for
// metrics updates by implementing the orca.OOBListener interface and manages
// that listener.
type endpointWeight struct {
=======
func (p *picker) Pick(balancer.PickInfo) (balancer.PickResult, error) {
	// Read the scheduler atomically.  All scheduler operations are threadsafe,
	// and if the scheduler is replaced during this usage, we want to use the
	// scheduler that was live when the pick started.
	sched := *(*scheduler)(atomic.LoadPointer(&p.scheduler))

	pickedSC := p.subConns[sched.nextIndex()]
	pr := balancer.PickResult{SubConn: pickedSC.SubConn}
	if !p.cfg.EnableOOBLoadReport {
		pr.Done = func(info balancer.DoneInfo) {
			if load, ok := info.ServerLoad.(*v3orcapb.OrcaLoadReport); ok && load != nil {
				pickedSC.OnLoadReport(load)
			}
		}
	}
	return pr, nil
}

// weightedSubConn is the wrapper of a subconn that holds the subconn and its
// weight (and other parameters relevant to computing the effective weight).
// When needed, it also tracks connectivity state, listens for metrics updates
// by implementing the orca.OOBListener interface and manages that listener.
type weightedSubConn struct {
>>>>>>> 70e0318b1 ([WIP] add archivista storage backend)
	// The following fields are immutable.
	balancer.SubConn
	logger          *grpclog.PrefixLogger
	target          string
	metricsRecorder estats.MetricsRecorder
	locality        string

	// The following fields are only accessed on calls into the LB policy, and
	// do not need a mutex.
	connectivityState connectivity.State
	stopORCAListener  func()
<<<<<<< HEAD
	// The first SubConn for the endpoint that goes READY when endpoint has no
	// READY SubConns yet, cleared on that sc disconnecting (i.e. going out of
	// READY). Represents what pick first will use as it's picked SubConn for
	// this endpoint.
	pickedSC balancer.SubConn
=======
>>>>>>> 70e0318b1 ([WIP] add archivista storage backend)

	// The following fields are accessed asynchronously and are protected by
	// mu.  Note that mu may not be held when calling into the stopORCAListener
	// or when registering a new listener, as those calls require the ORCA
	// producer mu which is held when calling the listener, and the listener
	// holds mu.
	mu            sync.Mutex
	weightVal     float64
	nonEmptySince time.Time
	lastUpdated   time.Time
	cfg           *lbConfig
}

<<<<<<< HEAD
func (w *endpointWeight) OnLoadReport(load *v3orcapb.OrcaLoadReport) {
	if w.logger.V(2) {
		w.logger.Infof("Received load report for subchannel %v: %v", w.SubConn, load)
	}
	// Update weights of this endpoint according to the reported load.
=======
func (w *weightedSubConn) OnLoadReport(load *v3orcapb.OrcaLoadReport) {
	if w.logger.V(2) {
		w.logger.Infof("Received load report for subchannel %v: %v", w.SubConn, load)
	}
	// Update weights of this subchannel according to the reported load
>>>>>>> 70e0318b1 ([WIP] add archivista storage backend)
	utilization := load.ApplicationUtilization
	if utilization == 0 {
		utilization = load.CpuUtilization
	}
	if utilization == 0 || load.RpsFractional == 0 {
		if w.logger.V(2) {
			w.logger.Infof("Ignoring empty load report for subchannel %v", w.SubConn)
		}
		return
	}

	w.mu.Lock()
	defer w.mu.Unlock()

	errorRate := load.Eps / load.RpsFractional
	w.weightVal = load.RpsFractional / (utilization + errorRate*w.cfg.ErrorUtilizationPenalty)
	if w.logger.V(2) {
		w.logger.Infof("New weight for subchannel %v: %v", w.SubConn, w.weightVal)
	}

	w.lastUpdated = internal.TimeNow()
	if w.nonEmptySince.Equal(time.Time{}) {
		w.nonEmptySince = w.lastUpdated
	}
}

// updateConfig updates the parameters of the WRR policy and
// stops/starts/restarts the ORCA OOB listener.
<<<<<<< HEAD
func (w *endpointWeight) updateConfig(cfg *lbConfig) {
=======
func (w *weightedSubConn) updateConfig(cfg *lbConfig) {
>>>>>>> 70e0318b1 ([WIP] add archivista storage backend)
	w.mu.Lock()
	oldCfg := w.cfg
	w.cfg = cfg
	w.mu.Unlock()

	if cfg.EnableOOBLoadReport == oldCfg.EnableOOBLoadReport &&
		cfg.OOBReportingPeriod == oldCfg.OOBReportingPeriod {
		// Load reporting wasn't enabled before or after, or load reporting was
		// enabled before and after, and had the same period.  (Note that with
		// load reporting disabled, OOBReportingPeriod is always 0.)
		return
	}
<<<<<<< HEAD
	// (Re)start the listener to use the new config's settings for OOB
	// reporting.
	w.updateORCAListener(cfg)
}

func (w *endpointWeight) updateORCAListener(cfg *lbConfig) {
=======
	if w.connectivityState == connectivity.Ready {
		// (Re)start the listener to use the new config's settings for OOB
		// reporting.
		w.updateORCAListener(cfg)
	}
}

func (w *weightedSubConn) updateORCAListener(cfg *lbConfig) {
>>>>>>> 70e0318b1 ([WIP] add archivista storage backend)
	if w.stopORCAListener != nil {
		w.stopORCAListener()
	}
	if !cfg.EnableOOBLoadReport {
		w.stopORCAListener = nil
		return
	}
<<<<<<< HEAD
	if w.pickedSC == nil { // No picked SC for this endpoint yet, nothing to listen on.
		return
	}
	if w.logger.V(2) {
		w.logger.Infof("Registering ORCA listener for %v with interval %v", w.pickedSC, cfg.OOBReportingPeriod)
	}
	opts := orca.OOBListenerOptions{ReportInterval: time.Duration(cfg.OOBReportingPeriod)}
	w.stopORCAListener = orca.RegisterOOBListener(w.pickedSC, w, opts)
}

// weight returns the current effective weight of the endpoint, taking into
=======
	if w.logger.V(2) {
		w.logger.Infof("Registering ORCA listener for %v with interval %v", w.SubConn, cfg.OOBReportingPeriod)
	}
	opts := orca.OOBListenerOptions{ReportInterval: time.Duration(cfg.OOBReportingPeriod)}
	w.stopORCAListener = orca.RegisterOOBListener(w.SubConn, w, opts)
}

func (w *weightedSubConn) updateConnectivityState(cs connectivity.State) connectivity.State {
	switch cs {
	case connectivity.Idle:
		// Always reconnect when idle.
		w.SubConn.Connect()
	case connectivity.Ready:
		// If we transition back to READY state, reset nonEmptySince so that we
		// apply the blackout period after we start receiving load data. Also
		// reset lastUpdated to trigger endpoint weight not yet usable in the
		// case endpoint gets asked what weight it is before receiving a new
		// load report. Note that we cannot guarantee that we will never receive
		// lingering callbacks for backend metric reports from the previous
		// connection after the new connection has been established, but they
		// should be masked by new backend metric reports from the new
		// connection by the time the blackout period ends.
		w.mu.Lock()
		w.nonEmptySince = time.Time{}
		w.lastUpdated = time.Time{}
		cfg := w.cfg
		w.mu.Unlock()
		w.updateORCAListener(cfg)
	}

	oldCS := w.connectivityState

	if oldCS == connectivity.TransientFailure &&
		(cs == connectivity.Connecting || cs == connectivity.Idle) {
		// Once a subconn enters TRANSIENT_FAILURE, ignore subsequent IDLE or
		// CONNECTING transitions to prevent the aggregated state from being
		// always CONNECTING when many backends exist but are all down.
		return oldCS
	}

	w.connectivityState = cs

	return oldCS
}

// weight returns the current effective weight of the subconn, taking into
>>>>>>> 70e0318b1 ([WIP] add archivista storage backend)
// account the parameters.  Returns 0 for blacked out or expired data, which
// will cause the backend weight to be treated as the mean of the weights of the
// other backends. If forScheduler is set to true, this function will emit
// metrics through the metrics registry.
<<<<<<< HEAD
func (w *endpointWeight) weight(now time.Time, weightExpirationPeriod, blackoutPeriod time.Duration, recordMetrics bool) (weight float64) {
=======
func (w *weightedSubConn) weight(now time.Time, weightExpirationPeriod, blackoutPeriod time.Duration, recordMetrics bool) (weight float64) {
>>>>>>> 70e0318b1 ([WIP] add archivista storage backend)
	w.mu.Lock()
	defer w.mu.Unlock()

	if recordMetrics {
		defer func() {
			endpointWeightsMetric.Record(w.metricsRecorder, weight, w.target, w.locality)
		}()
	}

<<<<<<< HEAD
	// The endpoint has not received a load report (i.e. just turned READY with
=======
	// The SubConn has not received a load report (i.e. just turned READY with
>>>>>>> 70e0318b1 ([WIP] add archivista storage backend)
	// no load report).
	if w.lastUpdated.Equal(time.Time{}) {
		endpointWeightNotYetUsableMetric.Record(w.metricsRecorder, 1, w.target, w.locality)
		return 0
	}

	// If the most recent update was longer ago than the expiration period,
	// reset nonEmptySince so that we apply the blackout period again if we
	// start getting data again in the future, and return 0.
	if now.Sub(w.lastUpdated) >= weightExpirationPeriod {
		if recordMetrics {
			endpointWeightStaleMetric.Record(w.metricsRecorder, 1, w.target, w.locality)
		}
		w.nonEmptySince = time.Time{}
		return 0
	}

	// If we don't have at least blackoutPeriod worth of data, return 0.
	if blackoutPeriod != 0 && (w.nonEmptySince.Equal(time.Time{}) || now.Sub(w.nonEmptySince) < blackoutPeriod) {
		if recordMetrics {
			endpointWeightNotYetUsableMetric.Record(w.metricsRecorder, 1, w.target, w.locality)
		}
		return 0
	}

	return w.weightVal
}
