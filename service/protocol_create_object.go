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
	"github.com/ailabstw/go-pttai/common/types"
	"github.com/ailabstw/go-pttai/log"
)

type CreateData interface{}

func (pm *BaseProtocolManager) CreateObject(
	data CreateData,
	createOp OpType,

	newObj func(data CreateData) (Object, OpData, error),
	newOplogWithTS func(objID *types.PttID, ts types.Timestamp, op OpType, opData OpData) (Oplog, error),
	increate func(obj Object, data CreateData, oplog *BaseOplog, opData OpData) error,
	broadcastLog func(oplog *BaseOplog) error,
	postcreate func(obj Object, oplog *BaseOplog) error,
) (Object, error) {

	entity := pm.Entity()

	// 1. validate
	if entity.GetStatus() != types.StatusAlive && entity.GetStatus() != types.StatusToBeSynced {
		return nil, types.ErrInvalidStatus
	}

	// 2. new-obj
	obj, opData, err := newObj(data)
	if err != nil {
		log.Warn("CreateObject: unable to newObj", "e", err)
		return nil, err
	}

	// 3. oplog
	theOplog, err := newOplogWithTS(obj.GetID(), obj.GetUpdateTS(), createOp, opData)
	if err != nil {
		log.Warn("CreateObject: unable to newOplogWithTS", "e", err)
		return nil, err
	}
	oplog := theOplog.GetBaseOplog()

	// 4. increate
	if increate != nil {
		err = increate(obj, data, oplog, opData)
		if err != nil {
			log.Warn("CreateObject: unable to increate", "e", err)
			return nil, err
		}
	}

	// 4.1. set is good
	obj.SetIsGood(true)
	blockInfo := obj.GetBlockInfo()
	if blockInfo != nil {
		blockInfo.SetIsAllGood()
	}
	obj.SetIsAllGood(true)

	// 5. sign oplog
	err = pm.SignOplog(oplog)
	if err != nil {
		log.Warn("CreateObject: unable to sign", "e", err)
		return nil, err
	}

	err = oplog.Verify()
	if err != nil {
		log.Error("CreateObject: unable to sign", "e", err)
		return nil, err
	}

	// 6. save object
	err = pm.saveNewObjectWithOplog(obj, oplog, false, false, postcreate)
	if err != nil {
		log.Warn("CreateObject: unable to saveNewObjectWithOplog", "e", err)
		return nil, err
	}

	// 7. oplog-save
	err = oplog.Save(false)
	if err != nil {
		log.Warn("CreateObject: unable to save oplog", "e", err)
		return nil, err
	}

	log.Debug("CreateObject: to broadcastLog", "obj", obj.GetID(), "oplog", oplog.ID)

	broadcastLog(oplog)

	return obj, nil
}

/*
SaveNewObjectWithOplog saves New Object with Oplog.
*/
func (pm *BaseProtocolManager) saveNewObjectWithOplog(
	obj Object,
	oplog *BaseOplog,

	isLocked bool,
	isForceNot bool,

	postcreate func(obj Object, oplog *BaseOplog) error,
) error {

	var err error

	if !isLocked {
		err = obj.Lock()
		if err != nil {
			return err
		}
		defer obj.Unlock()
	}

	// check is-synced
	isAllGood := obj.GetIsAllGood()
	if isAllGood {
		SetNewObjectWithOplog(obj, oplog)
		oplog.IsSync = true
	}
	log.Debug("saveNewObjectWithOplog: after isAllGood", "isAllGood", isAllGood, "obj", obj, "oplog", oplog.ID, "obj.status", obj.GetStatus(), "isForceNot", isForceNot)

	// save
	err = obj.Save(true)
	log.Debug("saveNewObjectWithOplog: after Save", "e", err)
	if err != nil {
		return err
	}

	// postcreate
	if isForceNot {
		return nil
	}

	if !isAllGood {
		return nil
	}

	if postcreate == nil {
		return nil
	}

	if oplog.ToStatus() != types.StatusAlive {
		return nil
	}

	return postcreate(obj, oplog)
}
