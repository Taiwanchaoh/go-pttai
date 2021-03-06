// Copyright 2018 The go-pttai Authors
// This file is part of the go-pttai library.
//
// The go-pttai library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-pttai library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-pttai library. If not, see <http://www.gnu.org/licenses/>.

package service

import (
	"time"

	"github.com/ailabstw/go-pttai/common/types"
	"github.com/ailabstw/go-pttai/log"
)

func PMOplogMerkleTreeLoop(pm ProtocolManager, merkle *Merkle) error {
	ticker := time.NewTicker(merkle.GenerateSeconds)
	defer ticker.Stop()

	pmGenerateOplogMerkleTree(pm, merkle)

loop:
	for {
		select {
		case <-ticker.C:
			pmGenerateOplogMerkleTree(pm, merkle)
		case <-pm.QuitSync():
			log.Debug("PMOplogMerkleTreeLoop: QuitSync", "entity", pm.Entity().GetID(), "service", pm.Entity().Service().Name())
			break loop
		}
	}

	return nil
}

func pmGenerateOplogMerkleTree(pm ProtocolManager, merkle *Merkle) error {
	status := pm.Entity().GetStatus()
	if status != types.StatusAlive {
		return nil
	}

	log.Debug("pmGenerateOplogMerkleTree: start", "entity", pm.Entity().GetID())

	now, err := types.GetTimestamp()
	if err != nil {
		return err
	}

	isBusy := pmGenerateOplogMerkleTreeIsBusy(merkle, now)
	if isBusy {
		return ErrBusy
	}

	// set busy
	merkle.BusyGenerateTS = now
	defer func() {
		merkle.BusyGenerateTS = types.ZeroTimestamp
	}()

	// save-merkle-tree

	ts := merkle.LastGenerateTS
	for ; ts.IsLess(now); ts.Ts += 3600 {
		err := merkle.SaveMerkleTree(ts)
		if err != nil {
			break
		}
	}

	log.Debug("pmGenerateOplogMerkleTree: done", "entity", pm.Entity().GetID())

	return nil
}

func pmGenerateOplogMerkleTreeIsBusy(merkle *Merkle, now types.Timestamp) bool {
	if merkle.BusyGenerateTS.IsEqual(types.ZeroTimestamp) {
		return false
	}

	expireTimestamp := now
	expireTimestamp.Ts -= merkle.ExpireGenerateSeconds

	if merkle.BusyGenerateTS.IsLess(expireTimestamp) {
		log.Warn("GenerateOplogMerkleTree expired", "busy-ts", merkle.BusyGenerateTS, "expire-ts", expireTimestamp)
		merkle.BusyGenerateTS = types.ZeroTimestamp
		return false
	}

	log.Warn("GenerateOplogMerkleTree is-busy", "busy-ts", merkle.BusyGenerateTS)
	return true
}
