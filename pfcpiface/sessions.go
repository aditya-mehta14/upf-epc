// SPDX-License-Identifier: Apache-2.0
// Copyright 2020 Intel Corporation

package pfcpiface

import (
	"fmt"
	"sync"

	"github.com/omec-project/upf-epc/pfcpiface/metrics"
	log "github.com/sirupsen/logrus"
)

type notifyFlag struct {
	flag bool
	mux  sync.Mutex
}

type PacketForwardingRules struct {
	pdrs []pdr
	fars []far
	qers []qer
}

// PFCPSession implements one PFCP session.
type PFCPSession struct {
	localSEID        uint64
	remoteSEID       uint64
	notificationFlag notifyFlag
	metrics          *metrics.Session
	PacketForwardingRules
}

func (p PacketForwardingRules) String() string {
	return fmt.Sprintf("PDRs=%v, FARs=%v, QERs=%v", p.pdrs, p.fars, p.qers)
}

// NewPFCPSession allocates an session with ID.
func (pConn *PFCPConn) NewPFCPSession(rseid uint64) uint64 {
	lseid, err := pConn.GenerateFseid()
	if err != nil {
		log.Errorln("Failed to generate session Id", err)
		return 0
	}

	s := PFCPSession{
		localSEID:  lseid,
		remoteSEID: rseid,
		PacketForwardingRules: PacketForwardingRules{
			pdrs: make([]pdr, 0, MaxItems),
			fars: make([]far, 0, MaxItems),
			qers: make([]qer, 0, MaxItems),
		},
	}
	pConn.sessions[lseid] = &s

	// Metrics update
	s.metrics = metrics.NewSession(pConn.nodeID.remote)
	pConn.SaveSessions(s.metrics)

	return lseid
}

// RemoveSession removes session using lseid.
func (pConn *PFCPConn) RemoveSession(lseid uint64) {
	s, ok := pConn.sessions[lseid]
	if !ok {
		return
	}

	// Metrics update
	s.metrics.Delete()
	pConn.SaveSessions(s.metrics)

	delete(pConn.sessions, lseid)
}
