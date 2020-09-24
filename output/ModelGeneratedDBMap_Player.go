package main

import (
	"servers/protocol"
)

func init() {
	DBMsgStructMap[protocol.DBMsgId_PlayerCreate] = &protocol.PlayerDB{}
	DBMsgStructMap[protocol.DBMsgId_PlayerCreateMany] = &protocol.PlayerManyDB{}
	DBMsgStructMap[protocol.DBMsgId_PlayerDelete] = &protocol.ModelDeleteDB{}
	DBMsgStructMap[protocol.DBMsgId_PlayerDeleteMany] = &protocol.ModelDeleteManyDB{}
	DBMsgStructMap[protocol.DBMsgId_PlayerUpdate] = &protocol.PlayerDB{}
	DBMsgStructMap[protocol.DBMsgId_PlayerUpdateMany] = &protocol.PlayerManyDB{}
	DBMsgStructMap[protocol.DBMsgId_PlayerSelect] = &protocol.PlayerDB{}
	DBMsgStructMap[protocol.DBMsgId_PlayerSelectMany] = &protocol.PlayerManyDB{}
}
