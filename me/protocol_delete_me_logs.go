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
	"github.com/ailabstw/go-pttai/common/types"
	pkgservice "github.com/ailabstw/go-pttai/service"
)

func (pm *ProtocolManager) handleDeleteMeLog(oplog *pkgservice.BaseOplog, info *ProcessMeInfo) ([]*pkgservice.BaseOplog, error) {

	opData := &MeOpDeleteMe{}

	toBroadcastLogs, err := pm.HandleDeleteEntityLog(
		oplog, info, opData, types.StatusDeleted,
		pm.SetMeDB, pm.postdeleteDeleteMe, pm.updateDeleteMeInfo)
	if err != nil {
		return nil, err
	}

	return toBroadcastLogs, nil
}

func (pm *ProtocolManager) handlePendingDeleteMeLog(oplog *pkgservice.BaseOplog, info *ProcessMeInfo) (types.Bool, []*pkgservice.BaseOplog, error) {

	opData := &MeOpDeleteMe{}
	return pm.HandlePendingDeleteEntityLog(
		oplog, info,
		types.StatusInternalDeleted, types.StatusPendingDeleted,
		MeOpTypeDeleteMe, opData,
		pm.SetMeDB, pm.setPendingDeleteMeSyncInfo, pm.updateDeleteMeInfo)
}

func (pm *ProtocolManager) setNewestDeleteMeLog(oplog *pkgservice.BaseOplog) (types.Bool, error) {

	return pm.SetNewestDeleteEntityLog(oplog)
}

func (pm *ProtocolManager) handleFailedDeleteMeLog(oplog *pkgservice.BaseOplog) error {

	return pm.HandleFailedDeleteEntityLog(oplog)
}
