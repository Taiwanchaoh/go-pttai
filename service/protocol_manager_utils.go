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
	"encoding/json"

	"github.com/ailabstw/go-pttai/log"
)

func PrestartPM(pm ProtocolManager) error {
	// 2. register entity
	ptt := pm.Ptt()
	err := ptt.RegisterEntity(pm.Entity(), false, false)
	if err != nil {
		return err
	}

	// 3. pre-start
	err = pm.Prestart()
	if err != nil {
		return err
	}

	return nil
}

/*
StartPM starts the pm
	1. go PMSync (require sync first to receive new-peer-ch)
	2. pm.Start
	3. go PMCreateOpKeyLoop
*/
func StartPM(pm ProtocolManager) error {
	log.Info("StartPM: start", "entity", pm.Entity().GetID())

	// 1. PMSync
	syncWG := pm.SyncWG()
	syncWG.Add(1)
	go func() {
		defer syncWG.Done()

		PMSync(pm)
	}()

	// 2. pm.Start
	err := pm.Start()
	if err != nil {
		return err
	}

	// 3. op-key
	syncWG.Add(1)
	go func() {
		defer syncWG.Done()
		pm.CreateOpKeyLoop()
	}()

	return nil
}

func StopPM(pm ProtocolManager) error {
	log.Info("Stop PM: to stop", "entity", pm.Entity().Name())

	err := pm.Stop()
	if err != nil {
		log.Warn("Stop PM: unable to stop", "entity", pm.Entity().Name(), "e", err)
		return err
	}

	log.Info("Stop PM: done", "entity", pm.Entity().Name())

	return nil
}

/*
Send Data to Peers using op-key
*/
func (pm *BaseProtocolManager) SendDataToPeers(op OpType, data interface{}, peerList []*PttPeer) error {

	dataBytes, err := json.Marshal(data)
	if err != nil {
		log.Error("SendDataToPeers: unable to marshal data", "e", err, "entity", pm.Entity().GetID())
		return err
	}

	opKeyInfo, err := pm.GetOldestOpKey(false)
	log.Debug("SendDataToPeers: after get opKey", "opKey", opKeyInfo.Hash, "entity", pm.Entity().GetID(), "e", err)

	if err != nil {
		return err
	}

	ptt := pm.Ptt()
	encData, err := ptt.EncryptData(op, dataBytes, opKeyInfo)
	if err != nil {
		return err
	}

	pttData, err := ptt.MarshalData(CodeTypeOp, opKeyInfo.Hash, encData)
	if err != nil {
		return err
	}

	okCount := 0
	for _, peer := range peerList {
		pttData.Node = peer.GetID()[:]
		err := peer.SendData(pttData)
		if err == nil {
			okCount++
		} else {
			log.Warn("PMSendDataToPeers: unable to SendData", "peer", peer, "e", err)
		}
	}
	if okCount == 0 {
		return ErrNotSent
	}

	return nil
}

func (pm *BaseProtocolManager) SendDataToPeer(op OpType, data interface{}, peer *PttPeer) error {
	return pm.SendDataToPeerWithCode(CodeTypeOp, op, data, peer)
}

/*
Send Data to Peers using op-key
*/
func (pm *BaseProtocolManager) SendDataToPeerWithCode(code CodeType, op OpType, data interface{}, peer *PttPeer) error {

	dataBytes, err := json.Marshal(data)
	if err != nil {
		log.Error("SendDataToPeerWithCode: unable to marshal data", "e", err, "entity", pm.Entity().GetID())
		return err
	}

	opKeyInfo, err := pm.GetOldestOpKey(false)
	log.Debug("SendDataToPeerWithCode: after get opKey", "opKey", opKeyInfo.Hash, "entity", pm.Entity().GetID(), "e", err)
	if err != nil {
		return err
	}

	ptt := pm.Ptt()
	encData, err := ptt.EncryptData(op, dataBytes, opKeyInfo)
	if err != nil {
		return err
	}

	pttData, err := ptt.MarshalData(code, opKeyInfo.Hash, encData)
	if err != nil {
		return err
	}

	pttData.Node = peer.GetID()[:]

	err = peer.SendData(pttData)
	log.Debug("SendDataToPeerWithCode: after SendData", "code", code, "peer", peer.PeerType, "e", err)
	if err != nil {
		return ErrNotSent
	}

	return nil
}
