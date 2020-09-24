package main

import (
	"context"
	"github.com/golang/protobuf/proto"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
	"servers/protocol"
	"wrappers/common"
	"wrappers/db_wrapper"
)

func init() {
	DBMsgProcessorMap[protocol.DBMsgId_PlayerCreate] = MsgProcPlayerCreate
	DBMsgProcessorMap[protocol.DBMsgId_PlayerCreateMany] = MsgProcPlayerCreateMany
	DBMsgProcessorMap[protocol.DBMsgId_PlayerDelete] = MsgProcPlayerDelete
	DBMsgProcessorMap[protocol.DBMsgId_PlayerDeleteMany] = MsgProcPlayerDeleteMany
	DBMsgProcessorMap[protocol.DBMsgId_PlayerUpdate] = MsgProcPlayerUpdate
	DBMsgProcessorMap[protocol.DBMsgId_PlayerUpdateMany] = MsgProcPlayerUpdateMany
	DBMsgProcessorMap[protocol.DBMsgId_PlayerSelect] = MsgProcPlayerSelect
	DBMsgProcessorMap[protocol.DBMsgId_PlayerSelectMany] = MsgProcPlayerSelectMany
}

func MsgProcPlayerCreate(srcAddr string, msgId uint16, userId string, message proto.Message) bool {
	zLog := common.ZLog
	msg := message.(*protocol.PlayerDB)
	collection := db_wrapper.MongoClient.Database(DBName).Collection("player")

	ctx := context.Background()

	for {
		bytes, err := proto.Marshal(msg.Info)
		if err != nil {
			zLog.Error("proto.Marshal error", zap.Error(err))
			break
		}

		infoBsonObject := bson.M{"Data": bytes}
		result, err := collection.InsertOne(ctx, infoBsonObject)
		if err != nil {
			zLog.Error("MsgProcPlayerCreate collection.InsertOne error", zap.Error(err))
			break
		} else {
			if msg.Info.Id == "" {
				msg.Info.Id = result.InsertedID.(primitive.ObjectID).Hex()
				zLog.Debug("MsgProcPlayerCreate", zap.String("InsertedID", msg.Info.Id))
			}
		}

		msg.SucceededCount = 1
		break
	}
	SendMsgBack(srcAddr, msgId, userId, msg)

	return true
}

func MsgProcPlayerCreateMany(srcAddr string, msgId uint16, userId string, message proto.Message) bool {
	zLog := common.ZLog
	msg := message.(*protocol.PlayerManyDB)
	collection := db_wrapper.MongoClient.Database(DBName).Collection("player")

	ctx := context.Background()

	for {
		inserts := bson.A{}

		marshalOk := true
		for _, info := range msg.Infos {
			bytes, err := proto.Marshal(info)
			if err != nil {
				zLog.Error("MsgProcPlayerCreateMany proto.Marshal error", zap.Error(err))
				marshalOk = false
				break
			}
			inserts = append(inserts, bson.M{"Data": bytes})
		}
		if !marshalOk {
			break
		}

		result, err := collection.InsertMany(ctx, inserts)
		if err != nil {
			zLog.Error("MsgProcPlayerCreateMany collection.InsertMany error", zap.Error(err))
			break
		}

		for idx, info := range msg.Infos {
			if idx >= len(result.InsertedIDs) {
				break
			}
			if info.Id == "" {
				info.Id = result.InsertedIDs[idx].(primitive.ObjectID).Hex()
				zLog.Debug("MsgProcPlayerCreateMany", zap.String("InsertedID", info.Id))
			}
		}

		msg.SucceededCount = int32(len(result.InsertedIDs))

		/*
			for _, info := range msg.Infos {
				bytes, err := proto.Marshal(info)
				if err != nil {
					zLog.Error("MsgProcPlayerCreateMany proto.Marshal error", zap.Error(err))
					return false
				}

				infoBsonObject := bson.M{"Data": bytes}
				result, err := collection.InsertOne(ctx, infoBsonObject)
				if err != nil {
					zLog.Error("MsgProcPlayerCreateMany collection.InsertOne error", zap.Error(err))
					return false
				} else {
					if info.Id == "" {
						info.Id = result.InsertedID.(primitive.ObjectID).Hex()
						zLog.Debug("MsgProcPlayerCreateMany", zap.String("InsertedID", info.Id))
					}
				}
			}
		*/
		break
	}

	SendMsgBack(srcAddr, msgId, userId, msg)

	return true
}

func MsgProcPlayerDelete(srcAddr string, msgId uint16, userId string, message proto.Message) bool {
	zLog := common.ZLog
	msg := message.(*protocol.ModelDeleteDB)
	collection := db_wrapper.MongoClient.Database(DBName).Collection("player")

	ctx := context.Background()

	for {
		Id, err := primitive.ObjectIDFromHex(msg.Id)
		if err != nil {
			zLog.Error("MsgProcPlayerDelete primitive.ObjectIDFromHex error", zap.String("msg.Id", msg.Id), zap.Error(err))
			break
		}
		filter := bson.M{"_id": Id}
		deleteResult, err := collection.DeleteOne(ctx, filter)
		if err != nil {
			zLog.Error("MsgProcPlayerDelete collection.DeleteOne error", zap.Error(err))
			break
		}

		msg.DeletedCount = deleteResult.DeletedCount

		break
	}

	SendMsgBack(srcAddr, msgId, userId, msg)

	return true
}

func MsgProcPlayerDeleteMany(srcAddr string, msgId uint16, userId string, message proto.Message) bool {
	zLog := common.ZLog
	msg := message.(*protocol.ModelDeleteManyDB)
	collection := db_wrapper.MongoClient.Database(DBName).Collection("player")

	ctx := context.Background()

	for {

		ids := bson.A{}

		for _, id := range msg.Ids {
			objId, err := primitive.ObjectIDFromHex(id)
			if err != nil {
				zLog.Error("MsgProcPlayerDeleteMany primitive.ObjectIDFromHex error", zap.String("id", id), zap.Error(err))
				break
			}
			ids = append(ids, objId)
		}

		filter := bson.M{"_id": bson.D{{"$in", ids}}}

		deleteResult, err := collection.DeleteMany(ctx, filter)
		if err != nil {
			zLog.Error("MsgProcPlayerDeleteMany collection.DeleteMany error", zap.Error(err))
			break
		}

		msg.DeletedCount = deleteResult.DeletedCount

		break
	}

	SendMsgBack(srcAddr, msgId, userId, msg)

	return true
}

func MsgProcPlayerUpdate(srcAddr string, msgId uint16, userId string, message proto.Message) bool {
	zLog := common.ZLog
	msg := message.(*protocol.PlayerDB)
	collection := db_wrapper.MongoClient.Database(DBName).Collection("player")

	ctx := context.Background()

	for {
		bytes, err := proto.Marshal(msg.Info)
		if err != nil {
			zLog.Error("proto.Marshal error", zap.Error(err))
			break
		}

		infoBsonObject := bson.M{"Data": bytes}

		Id := msg.Info.Id
		objId, err := primitive.ObjectIDFromHex(Id)
		if err != nil {
			zLog.Error("MsgProcPlayerUpdate primitive.ObjectIDFromHex failed", zap.String("Id", Id))
			break
		}
		filter := bson.M{"_id": objId}
		updateResult, err := collection.UpdateOne(ctx, filter, bson.M{"$set": infoBsonObject})
		if err != nil {
			zLog.Error("collection.UpdateOne err", zap.Error(err))
			break
		}

		msg.SucceededCount = int32(updateResult.ModifiedCount)
		break
	}
	SendMsgBack(srcAddr, msgId, userId, msg)

	return true
}

func MsgProcPlayerUpdateMany(srcAddr string, msgId uint16, userId string, message proto.Message) bool {
	zLog := common.ZLog
	msg := message.(*protocol.PlayerManyDB)
	collection := db_wrapper.MongoClient.Database(DBName).Collection("player")

	ctx := context.Background()

	for {
		msg.SucceededCount = 0
		for _, info := range msg.Infos {
			bytes, err := proto.Marshal(info)
			if err != nil {
				zLog.Error("MsgProcPlayerUpdateMany proto.Marshal error", zap.Error(err))
				continue
			}

			infoBsonObject := bson.M{"Data": bytes}

			Id := info.Id
			objId, err := primitive.ObjectIDFromHex(Id)
			if err != nil {
				zLog.Error("MsgProcPlayerUpdateMany primitive.ObjectIDFromHex failed", zap.String("Id", Id))
				continue
			}
			filter := bson.M{"_id": objId}
			updateResult, err := collection.UpdateOne(ctx, filter, bson.M{"$set": infoBsonObject})
			if err != nil {
				zLog.Error("collection.UpdateOne err", zap.Error(err))
				continue
			}

			msg.SucceededCount += int32(updateResult.ModifiedCount)
		}
		break
	}
	SendMsgBack(srcAddr, msgId, userId, msg)

	return true
}

func MsgProcPlayerSelect(srcAddr string, msgId uint16, userId string, message proto.Message) bool {
	zLog := common.ZLog
	msg := message.(*protocol.PlayerDB)
	collection := db_wrapper.MongoClient.Database(DBName).Collection("player")

	ctx := context.Background()

	for {
		var result bson.M

		var filter bson.M
		Id := msg.GetInfo().Id
		objId, err := primitive.ObjectIDFromHex(Id)
		if err != nil {
			zLog.Error("MsgProcPlayerSelect primitive.ObjectIDFromHex failed", zap.String("Id", Id))
			break
		}
		filter = bson.M{"_id": objId}

		fOne := collection.FindOne(ctx, filter, &options.FindOneOptions{})

		if fOne.Err() != nil {
			zLog.Warn("MsgProcPlayerSelect collection.FindOne fOne.Err()", zap.Error(fOne.Err()))
			break
		} else {
			err := fOne.Decode(&result)
			if err != nil {
				zLog.Error("MsgProcPlayerSelect fOne.Decode error", zap.Error(err))
				break
			}

			bytes := result["Data"].(primitive.Binary).Data
			o := protocol.Player{}
			err = proto.Unmarshal(bytes, &o)
			if err != nil {
				zLog.Error("proto.Unmarshal failed.", zap.Error(err))
				break
			}
			zLog.Debug("select", zap.String("info", o.String()))
			msg.Info = &o
			msg.Info.Id = Id
		}

		msg.SucceededCount = 1
		break
	}
	SendMsgBack(srcAddr, msgId, userId, msg)

	return true
}

func MsgProcPlayerSelectMany(srcAddr string, msgId uint16, userId string, message proto.Message) bool {
	zLog := common.ZLog
	msg := message.(*protocol.PlayerManyDB)
	collection := db_wrapper.MongoClient.Database(DBName).Collection("player")

	ctx := context.Background()

	for {
		ids := bson.A{}

		for _, info := range msg.Infos {
			objId, err := primitive.ObjectIDFromHex(info.Id)
			if err != nil {
				zLog.Error("MsgProcPlayerSelectMany primitive.ObjectIDFromHex error", zap.String("id", info.Id), zap.Error(err))
				break
			}
			ids = append(ids, objId)
		}

		filter := bson.M{"_id": bson.D{{"$in", ids}}}

		cursor, err := collection.Find(ctx, filter, &options.FindOptions{})
		if err != nil {
			zLog.Error("collection.Find failed", zap.Error(err))
			break
		}

		var results []bson.M
		if err = cursor.All(ctx, &results); err != nil {
			zLog.Error("cursor.All", zap.Error(err))
			break
		}
		for _, result := range results {
			bytes := result["Data"].(primitive.Binary).Data
			o := protocol.Player{}
			err = proto.Unmarshal(bytes, &o)
			if err != nil {
				zLog.Error("proto.Unmarshal failed.", zap.Error(err))
				break
			}
			zLog.Debug("selectmany", zap.String("info", o.String()))
			for _, _info := range msg.Infos {
				_objId := result["_id"].(primitive.ObjectID)
				_Id := _objId.Hex()
				if _info.Id == _Id {
					_info = &o
					_info.Id = _Id
					break
				}
			}
		}

		msg.SucceededCount = 1
		break
	}
	SendMsgBack(srcAddr, msgId, userId, msg)

	return true
}
