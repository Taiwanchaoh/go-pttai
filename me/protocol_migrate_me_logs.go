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

func (pm *ProtocolManager) handleMigrateMeLog(oplog *pkgservice.BaseOplog, info *ProcessMeInfo) ([]*pkgservice.BaseOplog, error) {

	opData := &MeOpMigrateMe{}

	toBroadcastLogs, err := pm.HandleDeleteEntityLog(
		oplog, info, opData, types.StatusMigrated,
		pm.SetMeDB, pm.postdeleteMigrateMe, pm.updateDeleteMeInfo)
	if err != nil {
		return nil, err
	}

	return toBroadcastLogs, nil
}

func (pm *ProtocolManager) handlePendingMigrateMeLog(oplog *pkgservice.BaseOplog, info *ProcessMeInfo) ([]*pkgservice.BaseOplog, error) {

	opData := &MeOpMigrateMe{}
	return pm.HandlePendingDeleteEntityLog(
		oplog, info,
		types.StatusInternalMigrate, types.StatusPendingMigrate,
		MeOpTypeMigrateMe, opData,
		pm.SetMeDB, pm.setPendingDeleteMeSyncInfo, pm.updateDeleteMeInfo)
}

func (pm *ProtocolManager) setNewestMigrateMeLog(oplog *pkgservice.BaseOplog) (types.Bool, error) {

	return pm.SetNewestDeleteEntityLog(oplog)
}

func (pm *ProtocolManager) handleFailedMigrateMeLog(oplog *pkgservice.BaseOplog) error {

	return pm.HandleFailedDeleteEntityLog(oplog)

}

func (pm *ProtocolManager) updateDeleteMeInfo(oplog *pkgservice.BaseOplog, theInfo pkgservice.ProcessInfo) error {

	entityID := pm.Entity().GetID()

	info, ok := theInfo.(*ProcessMeInfo)
	if !ok {
		return pkgservice.ErrInvalidData
	}

	info.DeleteMeInfo[*entityID] = oplog

	return nil
}