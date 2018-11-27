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

package me

import (
	"encoding/json"

	"github.com/ailabstw/go-pttai/common/types"
	"github.com/ailabstw/go-pttai/content"
	pkgservice "github.com/ailabstw/go-pttai/service"
)

type InternalSyncBoardAck struct {
	LogID *types.PttID `json:"l"`

	BoardData *pkgservice.ApproveJoinEntity `json:"B"`
}

func (pm *ProtocolManager) InternalSyncBoard(
	oplog *pkgservice.BaseOplog,
	peer *pkgservice.PttPeer,
) error {

	syncID := &pkgservice.SyncID{ID: oplog.ObjID, LogID: oplog.ID}

	return pm.SendDataToPeer(InternalSyncBoardMsg, syncID, peer)
}

func (pm *ProtocolManager) HandleInternalSyncBoard(
	dataBytes []byte,
	peer *pkgservice.PttPeer,
) error {

	syncID := &pkgservice.SyncID{}
	err := json.Unmarshal(dataBytes, syncID)
	if err != nil {
		return err
	}

	contentSPM := pm.Entity().Service().(*Backend).contentBackend.SPM()
	board := contentSPM.Entity(syncID.ID)
	if board == nil {
		return types.ErrInvalidID
	}
	boardPM := board.PM()

	myID := pm.Ptt().GetMyEntity().GetID()
	joinEntity := &pkgservice.JoinEntity{ID: myID}
	_, theApproveJoinEntity, err := boardPM.ApproveJoin(joinEntity, nil, peer)
	if err != nil {
		return err
	}

	approveJoinEntity, ok := theApproveJoinEntity.(*pkgservice.ApproveJoinEntity)
	if !ok {
		return pkgservice.ErrInvalidData
	}

	ackData := &InternalSyncBoardAck{LogID: syncID.LogID, BoardData: approveJoinEntity}

	pm.SendDataToPeer(InternalSyncBoardAckMsg, ackData, peer)

	return nil
}

func (pm *ProtocolManager) HandleInternalSyncBoardAck(
	dataBytes []byte,
	peer *pkgservice.PttPeer,

) error {

	// unmarshal data
	theBoardData := content.NewEmptyApproveJoinBoard()

	data := &InternalSyncBoardAck{BoardData: theBoardData}
	err := json.Unmarshal(dataBytes, data)
	if err != nil {
		return err
	}

	// oplog
	oplog := &pkgservice.BaseOplog{ID: data.LogID}
	err = oplog.Lock()
	if err != nil {
		return err
	}
	defer oplog.Unlock()

	pm.SetMeDB(oplog)
	err = oplog.Get(data.LogID, true)
	if oplog.IsSync {
		return nil
	}

	// contentSPM
	contentSPM := pm.Entity().Service().(*Backend).contentBackend.SPM().(*content.ServiceProtocolManager)

	err = contentSPM.Lock(oplog.ObjID)
	if err != nil {
		return err
	}
	defer contentSPM.Unlock(oplog.ObjID)

	board := contentSPM.Entity(oplog.ObjID)
	if board == nil {
		return pm.handleInternalSyncBoardAckNew(contentSPM, theBoardData, oplog, peer)
	}

	return types.ErrNotImplemented
}

func (pm *ProtocolManager) handleInternalSyncBoardAckNew(
	spm *content.ServiceProtocolManager,
	data *pkgservice.ApproveJoinEntity,
	oplog *pkgservice.BaseOplog,
	peer *pkgservice.PttPeer,
) error {

	_, err := spm.CreateJoinEntity(data, peer, oplog, true, true, true, true)
	if err != nil {
		return err
	}

	oplog.IsSync = true
	oplog.Save(true)

	return nil
}
