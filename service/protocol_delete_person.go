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
	"reflect"

	"github.com/ailabstw/go-pttai/common/types"
)

func (pm *BaseProtocolManager) DeletePerson(
	id *types.PttID,

	deleteOp OpType,
	origPerson Object,
	opData OpData,
	internalPendingStatus types.Status,
	pendingStatus types.Status,
	status types.Status,

	setLogDB func(oplog *BaseOplog),
	newOplog func(objID *types.PttID, op OpType, opData OpData) (Oplog, error),
	broadcastLog func(oplog *BaseOplog) error,
	postdelete func(id *types.PttID, oplog *BaseOplog, origPerson Object, opData OpData) error,
) error {

	myEntity := pm.Ptt().GetMyEntity()
	myID := myEntity.GetID()
	entity := pm.Entity()

	// validate
	if entity.GetStatus() != types.StatusAlive {
		return types.ErrInvalidStatus
	}
	if myEntity.GetStatus() != types.StatusAlive {
		return types.ErrInvalidStatus
	}

	if !pm.IsMaster(myID, false) && !reflect.DeepEqual(myID, id) {
		return types.ErrInvalidID
	}

	// lock orig-person
	origPerson.SetID(id)
	err := origPerson.Lock()
	if err != nil {
		return err
	}
	defer origPerson.Unlock()

	// get orig-person
	err = origPerson.GetByID(true)
	if err != nil {
		return err
	}

	// 3. check validity
	origStatus := origPerson.GetStatus()
	if origStatus >= types.StatusDeleted {
		return types.ErrAlreadyDeleted
	}

	// 4. oplog
	theOplog, err := newOplog(id, deleteOp, opData)
	if err != nil {
		return err
	}
	oplog := theOplog.GetBaseOplog()
	oplog.PreLogID = origPerson.GetLogID()

	err = pm.SignOplog(oplog)
	if err != nil {
		return err
	}

	// 5. update obj
	oplogStatus := types.StatusToDeleteStatus(oplog.ToStatus(), internalPendingStatus, pendingStatus, status)

	if oplogStatus >= types.StatusDeleted {
		err = pm.handleDeletePersonLogCore(oplog, origPerson, opData, oplogStatus, setLogDB, postdelete)
	} else {
		err = pm.handlePendingDeletePersonLogCore(oplog, origPerson, opData, internalPendingStatus, pendingStatus, setLogDB)
	}
	if err != nil {
		return err
	}

	// 6. oplog
	err = oplog.Save(true)
	if err != nil {
		return err
	}

	broadcastLog(oplog)

	return nil
}
