//go:build vtprotobuf
// +build vtprotobuf

// Code generated by protoc-gen-go-vtproto. DO NOT EDIT.
// source: envoy/extensions/load_balancing_policies/client_side_weighted_round_robin/v3/client_side_weighted_round_robin.proto

package client_side_weighted_round_robinv3

import (
	protohelpers "github.com/planetscale/vtprotobuf/protohelpers"
	durationpb "github.com/planetscale/vtprotobuf/types/known/durationpb"
	wrapperspb "github.com/planetscale/vtprotobuf/types/known/wrapperspb"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

func (m *ClientSideWeightedRoundRobin) MarshalVTStrict() (dAtA []byte, err error) {
	if m == nil {
		return nil, nil
	}
	size := m.SizeVT()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBufferVTStrict(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *ClientSideWeightedRoundRobin) MarshalToVTStrict(dAtA []byte) (int, error) {
	size := m.SizeVT()
	return m.MarshalToSizedBufferVTStrict(dAtA[:size])
}

func (m *ClientSideWeightedRoundRobin) MarshalToSizedBufferVTStrict(dAtA []byte) (int, error) {
	if m == nil {
		return 0, nil
	}
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.unknownFields != nil {
		i -= len(m.unknownFields)
		copy(dAtA[i:], m.unknownFields)
	}
<<<<<<< HEAD
	if len(m.MetricNamesForComputingUtilization) > 0 {
		for iNdEx := len(m.MetricNamesForComputingUtilization) - 1; iNdEx >= 0; iNdEx-- {
			i -= len(m.MetricNamesForComputingUtilization[iNdEx])
			copy(dAtA[i:], m.MetricNamesForComputingUtilization[iNdEx])
			i = protohelpers.EncodeVarint(dAtA, i, uint64(len(m.MetricNamesForComputingUtilization[iNdEx])))
			i--
			dAtA[i] = 0x3a
		}
	}
=======
>>>>>>> 70e0318b1 ([WIP] add archivista storage backend)
	if m.ErrorUtilizationPenalty != nil {
		size, err := (*wrapperspb.FloatValue)(m.ErrorUtilizationPenalty).MarshalToSizedBufferVTStrict(dAtA[:i])
		if err != nil {
			return 0, err
		}
		i -= size
		i = protohelpers.EncodeVarint(dAtA, i, uint64(size))
		i--
		dAtA[i] = 0x32
	}
	if m.WeightUpdatePeriod != nil {
		size, err := (*durationpb.Duration)(m.WeightUpdatePeriod).MarshalToSizedBufferVTStrict(dAtA[:i])
		if err != nil {
			return 0, err
		}
		i -= size
		i = protohelpers.EncodeVarint(dAtA, i, uint64(size))
		i--
		dAtA[i] = 0x2a
	}
	if m.WeightExpirationPeriod != nil {
		size, err := (*durationpb.Duration)(m.WeightExpirationPeriod).MarshalToSizedBufferVTStrict(dAtA[:i])
		if err != nil {
			return 0, err
		}
		i -= size
		i = protohelpers.EncodeVarint(dAtA, i, uint64(size))
		i--
		dAtA[i] = 0x22
	}
	if m.BlackoutPeriod != nil {
		size, err := (*durationpb.Duration)(m.BlackoutPeriod).MarshalToSizedBufferVTStrict(dAtA[:i])
		if err != nil {
			return 0, err
		}
		i -= size
		i = protohelpers.EncodeVarint(dAtA, i, uint64(size))
		i--
		dAtA[i] = 0x1a
	}
	if m.OobReportingPeriod != nil {
		size, err := (*durationpb.Duration)(m.OobReportingPeriod).MarshalToSizedBufferVTStrict(dAtA[:i])
		if err != nil {
			return 0, err
		}
		i -= size
		i = protohelpers.EncodeVarint(dAtA, i, uint64(size))
		i--
		dAtA[i] = 0x12
	}
	if m.EnableOobLoadReport != nil {
		size, err := (*wrapperspb.BoolValue)(m.EnableOobLoadReport).MarshalToSizedBufferVTStrict(dAtA[:i])
		if err != nil {
			return 0, err
		}
		i -= size
		i = protohelpers.EncodeVarint(dAtA, i, uint64(size))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *ClientSideWeightedRoundRobin) SizeVT() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.EnableOobLoadReport != nil {
		l = (*wrapperspb.BoolValue)(m.EnableOobLoadReport).SizeVT()
		n += 1 + l + protohelpers.SizeOfVarint(uint64(l))
	}
	if m.OobReportingPeriod != nil {
		l = (*durationpb.Duration)(m.OobReportingPeriod).SizeVT()
		n += 1 + l + protohelpers.SizeOfVarint(uint64(l))
	}
	if m.BlackoutPeriod != nil {
		l = (*durationpb.Duration)(m.BlackoutPeriod).SizeVT()
		n += 1 + l + protohelpers.SizeOfVarint(uint64(l))
	}
	if m.WeightExpirationPeriod != nil {
		l = (*durationpb.Duration)(m.WeightExpirationPeriod).SizeVT()
		n += 1 + l + protohelpers.SizeOfVarint(uint64(l))
	}
	if m.WeightUpdatePeriod != nil {
		l = (*durationpb.Duration)(m.WeightUpdatePeriod).SizeVT()
		n += 1 + l + protohelpers.SizeOfVarint(uint64(l))
	}
	if m.ErrorUtilizationPenalty != nil {
		l = (*wrapperspb.FloatValue)(m.ErrorUtilizationPenalty).SizeVT()
		n += 1 + l + protohelpers.SizeOfVarint(uint64(l))
	}
<<<<<<< HEAD
	if len(m.MetricNamesForComputingUtilization) > 0 {
		for _, s := range m.MetricNamesForComputingUtilization {
			l = len(s)
			n += 1 + l + protohelpers.SizeOfVarint(uint64(l))
		}
	}
=======
>>>>>>> 70e0318b1 ([WIP] add archivista storage backend)
	n += len(m.unknownFields)
	return n
}
