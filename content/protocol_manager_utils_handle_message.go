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

package content

import (
	"github.com/ailabstw/go-pttai/log"
	pkgservice "github.com/ailabstw/go-pttai/service"
)

func (pm *ProtocolManager) HandleMessage(op pkgservice.OpType, dataBytes []byte, peer *pkgservice.PttPeer) error {

	log.Debug("content.HandleMessage: start", "op", op, "SyncCreateTitleMsg", SyncCreateTitleMsg)

	var err error
	switch op {
	// friend oplog
	case SyncBoardOplogMsg:
		err = pm.HandleSyncBoardOplog(dataBytes, peer)
	case SyncBoardOplogAckMsg:
		err = pm.HandleSyncBoardOplogAck(dataBytes, peer)
	case SyncBoardOplogNewOplogsMsg:
		err = pm.HandleSyncNewBoardOplog(dataBytes, peer)
	case SyncBoardOplogNewOplogsAckMsg:
		err = pm.HandleSyncNewBoardOplogAck(dataBytes, peer)
	case SyncPendingBoardOplogMsg:
		err = pm.HandleSyncPendingBoardOplog(dataBytes, peer)
	case SyncPendingBoardOplogAckMsg:
		err = pm.HandleSyncPendingBoardOplogAck(dataBytes, peer)

	case AddBoardOplogMsg:
		err = pm.HandleAddBoardOplog(dataBytes, peer)
	case AddBoardOplogsMsg:
		err = pm.HandleAddBoardOplogs(dataBytes, peer)
	case AddPendingBoardOplogMsg:
		err = pm.HandleAddPendingBoardOplog(dataBytes, peer)
	case AddPendingBoardOplogsMsg:
		err = pm.HandleAddPendingBoardOplogs(dataBytes, peer)

	// title
	case SyncCreateTitleMsg:
		err = pm.HandleSyncCreateTitle(dataBytes, peer, SyncCreateTitleAckMsg)
	case SyncCreateTitleAckMsg:
		err = pm.HandleSyncCreateTitleAck(dataBytes, peer)
	case SyncUpdateTitleMsg:
		err = pm.HandleSyncUpdateTitle(dataBytes, peer, SyncUpdateTitleAckMsg)
	case SyncUpdateTitleAckMsg:
		err = pm.HandleSyncUpdateTitleAck(dataBytes, peer)

	// article
	case SyncCreateArticleMsg:
		err = pm.HandleSyncCreateArticle(dataBytes, peer, SyncCreateArticleAckMsg)
	case SyncCreateArticleAckMsg:
		err = pm.HandleSyncCreateArticleAck(dataBytes, peer)
	case SyncCreateArticleBlockMsg:
		err = pm.HandleSyncArticleBlock(dataBytes, peer, SyncCreateArticleBlockAckMsg)
	case SyncCreateArticleBlockAckMsg:
		err = pm.HandleSyncCreateArticleBlockAck(dataBytes, peer)

	case SyncUpdateArticleMsg:
		err = pm.HandleSyncUpdateArticle(dataBytes, peer, SyncUpdateArticleAckMsg)
	case SyncUpdateArticleAckMsg:
		err = pm.HandleSyncUpdateArticleAck(dataBytes, peer)
	case SyncUpdateArticleBlockMsg:
		err = pm.HandleSyncArticleBlock(dataBytes, peer, SyncUpdateArticleBlockAckMsg)
	case SyncUpdateArticleBlockAckMsg:
		err = pm.HandleSyncUpdateArticleBlockAck(dataBytes, peer)

	// comment
	case SyncCreateCommentMsg:
		err = pm.HandleSyncCreateComment(dataBytes, peer, SyncCreateCommentAckMsg)
	case SyncCreateCommentAckMsg:
		err = pm.HandleSyncCreateCommentAck(dataBytes, peer)
	case SyncCreateCommentBlockMsg:
		err = pm.HandleSyncCommentBlock(dataBytes, peer, SyncCreateCommentBlockAckMsg)
	case SyncCreateCommentBlockAckMsg:
		err = pm.HandleSyncCreateCommentBlockAck(dataBytes, peer)

	// media
	case SyncCreateMediaMsg:
		err = pm.HandleSyncCreateMedia(dataBytes, peer, SyncCreateMediaAckMsg)
	case SyncCreateMediaAckMsg:
		err = pm.HandleSyncCreateMediaAck(dataBytes, peer, pm.SetBoardDB, pm.broadcastBoardOplogCore)
	case SyncCreateMediaBlockMsg:
		err = pm.HandleSyncMediaBlock(dataBytes, peer, SyncCreateMediaBlockAckMsg)
	case SyncCreateMediaBlockAckMsg:
		err = pm.HandleSyncCreateMediaBlockAck(dataBytes, peer, pm.SetBoardDB, pm.broadcastBoardOplogCore)

	default:
		log.Error("invalid op", "op", op, "SyncCreateTitleMsg", SyncCreateTitleMsg)
		err = pkgservice.ErrInvalidMsgCode
	}

	return err
}
