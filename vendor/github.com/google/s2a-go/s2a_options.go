/*
 *
 * Copyright 2021 Google LLC
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     https://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

package s2a

import (
<<<<<<< HEAD
=======
	"context"
>>>>>>> 70e0318b1 ([WIP] add archivista storage backend)
	"crypto/tls"
	"errors"
	"sync"

	"github.com/google/s2a-go/fallback"
	"github.com/google/s2a-go/stream"
	"google.golang.org/grpc/credentials"

<<<<<<< HEAD
	s2av1pb "github.com/google/s2a-go/internal/proto/common_go_proto"
=======
	s2apbv1 "github.com/google/s2a-go/internal/proto/common_go_proto"
>>>>>>> 70e0318b1 ([WIP] add archivista storage backend)
	s2apb "github.com/google/s2a-go/internal/proto/v2/common_go_proto"
)

// Identity is the interface for S2A identities.
type Identity interface {
	// Name returns the name of the identity.
	Name() string
<<<<<<< HEAD
	Attributes() map[string]string
}

type UnspecifiedID struct {
	Attr map[string]string
}

func (u *UnspecifiedID) Name() string { return "" }

func (u *UnspecifiedID) Attributes() map[string]string {
	return u.Attr
=======
>>>>>>> 70e0318b1 ([WIP] add archivista storage backend)
}

type spiffeID struct {
	spiffeID string
}

func (s *spiffeID) Name() string { return s.spiffeID }

<<<<<<< HEAD
func (spiffeID) Attributes() map[string]string { return nil }

// NewSpiffeID creates a SPIFFE ID from id.
func NewSpiffeID(id string) Identity { return &spiffeID{spiffeID: id} }
=======
// NewSpiffeID creates a SPIFFE ID from id.
func NewSpiffeID(id string) Identity {
	return &spiffeID{spiffeID: id}
}
>>>>>>> 70e0318b1 ([WIP] add archivista storage backend)

type hostname struct {
	hostname string
}

func (h *hostname) Name() string { return h.hostname }

<<<<<<< HEAD
func (hostname) Attributes() map[string]string { return nil }

// NewHostname creates a hostname from name.
func NewHostname(name string) Identity { return &hostname{hostname: name} }
=======
// NewHostname creates a hostname from name.
func NewHostname(name string) Identity {
	return &hostname{hostname: name}
}
>>>>>>> 70e0318b1 ([WIP] add archivista storage backend)

type uid struct {
	uid string
}

func (h *uid) Name() string { return h.uid }

<<<<<<< HEAD
func (uid) Attributes() map[string]string { return nil }

// NewUID creates a UID from name.
func NewUID(name string) Identity { return &uid{uid: name} }
=======
// NewUID creates a UID from name.
func NewUID(name string) Identity {
	return &uid{uid: name}
}
>>>>>>> 70e0318b1 ([WIP] add archivista storage backend)

// VerificationModeType specifies the mode that S2A must use to verify the peer
// certificate chain.
type VerificationModeType int

// Three types of verification modes.
const (
	Unspecified VerificationModeType = iota
	Spiffe
	ConnectToGoogle
	ReservedCustomVerificationMode3
	ReservedCustomVerificationMode4
	ReservedCustomVerificationMode5
<<<<<<< HEAD
	ReservedCustomVerificationMode6
=======
>>>>>>> 70e0318b1 ([WIP] add archivista storage backend)
)

// ClientOptions contains the client-side options used to establish a secure
// channel using the S2A handshaker service.
type ClientOptions struct {
	// TargetIdentities contains a list of allowed server identities. One of the
	// target identities should match the peer identity in the handshake
	// result; otherwise, the handshake fails.
	TargetIdentities []Identity
	// LocalIdentity is the local identity of the client application. If none is
	// provided, then the S2A will choose the default identity, if one exists.
	LocalIdentity Identity
	// S2AAddress is the address of the S2A.
	S2AAddress string
	// Optional transport credentials.
	// If set, this will be used for the gRPC connection to the S2A server.
	TransportCreds credentials.TransportCredentials
	// EnsureProcessSessionTickets waits for all session tickets to be sent to
	// S2A before a process completes.
	//
	// This functionality is crucial for processes that complete very soon after
	// using S2A to establish a TLS connection, but it can be ignored for longer
	// lived processes.
	//
	// Usage example:
	//   func main() {
	//     var ensureProcessSessionTickets sync.WaitGroup
	//     clientOpts := &s2a.ClientOptions{
	//       EnsureProcessSessionTickets: &ensureProcessSessionTickets,
	//       // Set other members.
	//     }
	//     creds, _ := s2a.NewClientCreds(clientOpts)
	//     conn, _ := grpc.Dial(serverAddr, grpc.WithTransportCredentials(creds))
	//     defer conn.Close()
	//
	//     // Make RPC call.
	//
	//     // The process terminates right after the RPC call ends.
	//     // ensureProcessSessionTickets can be used to ensure resumption
	//     // tickets are fully processed. If the process is long-lived, using
	//     // ensureProcessSessionTickets is not necessary.
	//     ensureProcessSessionTickets.Wait()
	//   }
	EnsureProcessSessionTickets *sync.WaitGroup
	// If true, enables the use of legacy S2Av1.
	EnableLegacyMode bool
	// VerificationMode specifies the mode that S2A must use to verify the
	// peer certificate chain.
	VerificationMode VerificationModeType

	// Optional fallback after dialing with S2A fails.
	FallbackOpts *FallbackOptions

	// Generates an S2AStream interface for talking to the S2A server.
<<<<<<< HEAD
	getS2AStream stream.GetS2AStream
=======
	getS2AStream func(ctx context.Context, s2av2Address string) (stream.S2AStream, error)
>>>>>>> 70e0318b1 ([WIP] add archivista storage backend)

	// Serialized user specified policy for server authorization.
	serverAuthorizationPolicy []byte
}

// FallbackOptions prescribes the fallback logic that should be taken if the application fails to connect with S2A.
type FallbackOptions struct {
	// FallbackClientHandshakeFunc is used to specify fallback behavior when calling s2a.NewClientCreds().
	// It will be called by ClientHandshake function, after handshake with S2A fails.
	// s2a.NewClientCreds() ignores the other FallbackDialer field.
	FallbackClientHandshakeFunc fallback.ClientHandshake

	// FallbackDialer is used to specify fallback behavior when calling s2a.NewS2aDialTLSContextFunc().
	// It passes in a custom fallback dialer and server address to use after dialing with S2A fails.
	// s2a.NewS2aDialTLSContextFunc() ignores the other FallbackClientHandshakeFunc field.
	FallbackDialer *FallbackDialer
}

// FallbackDialer contains a fallback tls.Dialer and a server address to connect to.
type FallbackDialer struct {
	// Dialer specifies a fallback tls.Dialer.
	Dialer *tls.Dialer
	// ServerAddr is used by Dialer to establish fallback connection.
	ServerAddr string
}

// DefaultClientOptions returns the default client options.
func DefaultClientOptions(s2aAddress string) *ClientOptions {
	return &ClientOptions{
		S2AAddress:       s2aAddress,
		VerificationMode: ConnectToGoogle,
	}
}

// ServerOptions contains the server-side options used to establish a secure
// channel using the S2A handshaker service.
type ServerOptions struct {
	// LocalIdentities is the list of local identities that may be assumed by
	// the server. If no local identity is specified, then the S2A chooses a
	// default local identity, if one exists.
	LocalIdentities []Identity
	// S2AAddress is the address of the S2A.
	S2AAddress string
	// Optional transport credentials.
	// If set, this will be used for the gRPC connection to the S2A server.
	TransportCreds credentials.TransportCredentials
	// If true, enables the use of legacy S2Av1.
	EnableLegacyMode bool
	// VerificationMode specifies the mode that S2A must use to verify the
	// peer certificate chain.
	VerificationMode VerificationModeType

	// Generates an S2AStream interface for talking to the S2A server.
<<<<<<< HEAD
	getS2AStream stream.GetS2AStream
=======
	getS2AStream func(ctx context.Context, s2av2Address string) (stream.S2AStream, error)
>>>>>>> 70e0318b1 ([WIP] add archivista storage backend)
}

// DefaultServerOptions returns the default server options.
func DefaultServerOptions(s2aAddress string) *ServerOptions {
	return &ServerOptions{
		S2AAddress:       s2aAddress,
		VerificationMode: ConnectToGoogle,
	}
}

<<<<<<< HEAD
func toProtoIdentity(identity Identity) (*s2av1pb.Identity, error) {
=======
func toProtoIdentity(identity Identity) (*s2apbv1.Identity, error) {
>>>>>>> 70e0318b1 ([WIP] add archivista storage backend)
	if identity == nil {
		return nil, nil
	}
	switch id := identity.(type) {
	case *spiffeID:
<<<<<<< HEAD
		return &s2av1pb.Identity{
			IdentityOneof: &s2av1pb.Identity_SpiffeId{SpiffeId: id.Name()},
			Attributes:    id.Attributes(),
		}, nil
	case *hostname:
		return &s2av1pb.Identity{
			IdentityOneof: &s2av1pb.Identity_Hostname{Hostname: id.Name()},
			Attributes:    id.Attributes(),
		}, nil
	case *uid:
		return &s2av1pb.Identity{
			IdentityOneof: &s2av1pb.Identity_Uid{Uid: id.Name()},
			Attributes:    id.Attributes(),
		}, nil
	case *UnspecifiedID:
		return &s2av1pb.Identity{
			Attributes: id.Attributes(),
		}, nil
=======
		return &s2apbv1.Identity{IdentityOneof: &s2apbv1.Identity_SpiffeId{SpiffeId: id.Name()}}, nil
	case *hostname:
		return &s2apbv1.Identity{IdentityOneof: &s2apbv1.Identity_Hostname{Hostname: id.Name()}}, nil
	case *uid:
		return &s2apbv1.Identity{IdentityOneof: &s2apbv1.Identity_Uid{Uid: id.Name()}}, nil
>>>>>>> 70e0318b1 ([WIP] add archivista storage backend)
	default:
		return nil, errors.New("unrecognized identity type")
	}
}

func toV2ProtoIdentity(identity Identity) (*s2apb.Identity, error) {
	if identity == nil {
		return nil, nil
	}
	switch id := identity.(type) {
	case *spiffeID:
<<<<<<< HEAD
		return &s2apb.Identity{
			IdentityOneof: &s2apb.Identity_SpiffeId{SpiffeId: id.Name()},
			Attributes:    id.Attributes(),
		}, nil
	case *hostname:
		return &s2apb.Identity{
			IdentityOneof: &s2apb.Identity_Hostname{Hostname: id.Name()},
			Attributes:    id.Attributes(),
		}, nil
	case *uid:
		return &s2apb.Identity{
			IdentityOneof: &s2apb.Identity_Uid{Uid: id.Name()},
			Attributes:    id.Attributes(),
		}, nil
	case *UnspecifiedID:
		return &s2apb.Identity{
			Attributes: id.Attributes(),
		}, nil
=======
		return &s2apb.Identity{IdentityOneof: &s2apb.Identity_SpiffeId{SpiffeId: id.Name()}}, nil
	case *hostname:
		return &s2apb.Identity{IdentityOneof: &s2apb.Identity_Hostname{Hostname: id.Name()}}, nil
	case *uid:
		return &s2apb.Identity{IdentityOneof: &s2apb.Identity_Uid{Uid: id.Name()}}, nil
>>>>>>> 70e0318b1 ([WIP] add archivista storage backend)
	default:
		return nil, errors.New("unrecognized identity type")
	}
}
