/*
 *
 * Copyright 2020 gRPC authors.
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

package xdsclient

import (
<<<<<<< HEAD
=======
	"fmt"
>>>>>>> 70e0318b1 ([WIP] add archivista storage backend)
	"sync/atomic"
	"time"

	"google.golang.org/grpc/internal/backoff"
	"google.golang.org/grpc/internal/grpcsync"
	"google.golang.org/grpc/internal/xds/bootstrap"
)

<<<<<<< HEAD
const defaultWatchExpiryTimeout = 15 * time.Second
=======
const (
	defaultWatchExpiryTimeout         = 15 * time.Second
	defaultIdleAuthorityDeleteTimeout = 5 * time.Minute
)
>>>>>>> 70e0318b1 ([WIP] add archivista storage backend)

var (
	// The following functions are no-ops in the actual code, but can be
	// overridden in tests to give them visibility into certain events.
	xdsClientImplCreateHook = func(string) {}
	xdsClientImplCloseHook  = func(string) {}

	defaultStreamBackoffFunc = backoff.DefaultExponential.Backoff
)

func clientRefCountedClose(name string) {
	clientsMu.Lock()
<<<<<<< HEAD
	client, ok := clients[name]
	if !ok {
		logger.Errorf("Attempt to close a non-existent xDS client with name %s", name)
		clientsMu.Unlock()
		return
	}
	if client.decrRef() != 0 {
		clientsMu.Unlock()
		return
	}
	delete(clients, name)
	clientsMu.Unlock()

	// This attempts to close the transport to the management server and could
	// theoretically call back into the xdsclient package again and deadlock.
	// Hence, this needs to be called without holding the lock.
	client.clientImpl.close()
	xdsClientImplCloseHook(name)
=======
	defer clientsMu.Unlock()

	client, ok := clients[name]
	if !ok {
		logger.Errorf("Attempt to close a non-existent xDS client with name %s", name)
		return
	}
	if client.decrRef() != 0 {
		return
	}
	client.clientImpl.close()
	xdsClientImplCloseHook(name)
	delete(clients, name)
>>>>>>> 70e0318b1 ([WIP] add archivista storage backend)

}

// newRefCounted creates a new reference counted xDS client implementation for
// name, if one does not exist already. If an xDS client for the given name
// exists, it gets a reference to it and returns it.
<<<<<<< HEAD
func newRefCounted(name string, config *bootstrap.Config, watchExpiryTimeout time.Duration, streamBackoff func(int) time.Duration) (XDSClient, func(), error) {
=======
func newRefCounted(name string, watchExpiryTimeout, idleAuthorityTimeout time.Duration, streamBackoff func(int) time.Duration) (XDSClient, func(), error) {
>>>>>>> 70e0318b1 ([WIP] add archivista storage backend)
	clientsMu.Lock()
	defer clientsMu.Unlock()

	if c := clients[name]; c != nil {
		c.incrRef()
		return c, grpcsync.OnceFunc(func() { clientRefCountedClose(name) }), nil
	}

	// Create the new client implementation.
<<<<<<< HEAD
	c, err := newClientImpl(config, watchExpiryTimeout, streamBackoff)
=======
	config, err := bootstrap.GetConfiguration()
	if err != nil {
		return nil, nil, fmt.Errorf("xds: failed to get xDS bootstrap config: %v", err)
	}
	c, err := newClientImpl(config, watchExpiryTimeout, idleAuthorityTimeout, streamBackoff)
>>>>>>> 70e0318b1 ([WIP] add archivista storage backend)
	if err != nil {
		return nil, nil, err
	}
	c.logger.Infof("Created client with name %q and bootstrap configuration:\n %s", name, config)
	client := &clientRefCounted{clientImpl: c, refCount: 1}
	clients[name] = client
	xdsClientImplCreateHook(name)

	logger.Infof("xDS node ID: %s", config.Node().GetId())
	return client, grpcsync.OnceFunc(func() { clientRefCountedClose(name) }), nil
}

// clientRefCounted is ref-counted, and to be shared by the xds resolver and
// balancer implementations, across multiple ClientConns and Servers.
type clientRefCounted struct {
	*clientImpl

	refCount int32 // accessed atomically
}

func (c *clientRefCounted) incrRef() int32 {
	return atomic.AddInt32(&c.refCount, 1)
}

func (c *clientRefCounted) decrRef() int32 {
	return atomic.AddInt32(&c.refCount, -1)
}
