package store

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"meme3/global"
	"meme3/service/model"

	"reflect"
	"time"

	"github.com/qiniu/qmgo"
	"github.com/qiniu/qmgo/middleware"
	"github.com/qiniu/qmgo/operator"
	"github.com/qiniu/qmgo/options"
	uuid "github.com/satori/go.uuid"
	"github.com/twmb/murmur3"
	"go.mongodb.org/mongo-driver/bson"
	moptions "go.mongodb.org/mongo-driver/mongo/options"
)

func NewId() (id string) {
	id, _ = newId()
	return
}

func newId() (id string, err error) {
	// uid, err := uuid.NewV1()
	// if err != nil {
	// 	return "", fmt.Errorf("generate uuid failed:%v", err)
	// }
	uid := uuid.NewV1()
	hash := murmur3.Sum32([]byte(uid.String()))
	id = fmt.Sprintf("%v", hash)
	return
}

func DBInit() error {
	err := InitMongoDB(&Mongo{
		URI:           global.Config.DBUrl,
		DBName:        global.Config.DBName,
		UserName:      global.Config.DBUser,
		Password:      global.Config.DBPassword,
		AuthMechanism: "SCRAM-SHA-1",
		AuthDBName:    "admin",
	})
	if err == nil {
		createOpt := options.CreateCollectionOptions{CreateCollectionOptions: Coptions}
		GetDB().CreateCollection(context.Background(), model.C_User, createOpt)
		GetDB().CreateCollection(context.Background(), model.C_Token, createOpt)
		GetDB().CreateCollection(context.Background(), model.C_Trade, createOpt)
		GetDB().CreateCollection(context.Background(), model.C_Referral, createOpt)
		GetDB().CreateCollection(context.Background(), model.C_Token_holder, createOpt)
		GetDB().CreateCollection(context.Background(), model.C_Listing, createOpt)
		GetDB().Collection(model.C_User).CreateIndexes(context.Background(), []options.IndexModel{{Key: []string{"-createdAt"}}})
		GetDB().Collection(model.C_Trade).CreateIndexes(context.Background(), []options.IndexModel{{Key: []string{"-createdAt"}}})
		GetDB().Collection(model.C_Referral).CreateIndexes(context.Background(), []options.IndexModel{{Key: []string{"-createdAt"}}})
		GetDB().Collection(model.C_Listing).CreateIndexes(context.Background(), []options.IndexModel{{Key: []string{"-createdAt"}}})
		GetDB().Collection(model.C_Token).CreateIndexes(context.Background(), []options.IndexModel{{Key: []string{"-createdAt"}}, {Key: []string{"-marketcap"}}, {Key: []string{"-volume"}}, {Key: []string{"-progress"}}})
		GetDB().Collection(model.C_Token_holder).CreateIndexes(context.Background(), []options.IndexModel{{Key: []string{"-createdAt"}}, {Key: []string{"-percent"}}})
	} else {
		DisableMongoDB = true
		slog.Warn("Init mongodb failed", "error", err.Error())
	}

	// test logic
	// UpdateValues()

	return nil
}

func DBCreate[T any](collName string, model, whereModel any) (retErr error) {
	if DisableMongoDB {
		return nil
	}
	coll := GetDB().Collection(collName)
	if whereModel != nil {
		if getModel, err := DBGet[T](collName, whereModel); err == nil && getModel != nil {
			return errors.New("duplicate document creation")
		}
	}

	rmodel := reflect.ValueOf(model)
	if rmodel.Kind() == reflect.Ptr {
		rmodel.Elem().FieldByName("CreatedAt").Set(reflect.ValueOf(time.Now().UnixMilli()))
		rmodel.Elem().FieldByName("UpdatedAt").Set(reflect.ValueOf(time.Now().UnixMilli()))
	}
	_, retErr = coll.InsertOne(context.Background(), model)
	global.Debug("DBCreate handle successed", "collName", collName, "model", model, "error", retErr)
	return
}

// upsert接口需要外部初始化model.CreatedAt
func DBSet(collName string, model, whereModel any) (retErr error) {
	if DisableMongoDB {
		return nil
	}
	coll := GetDB().Collection(collName)

	rmodel := reflect.ValueOf(model)
	if rmodel.Kind() == reflect.Ptr {
		rmodel.Elem().FieldByName("UpdatedAt").Set(reflect.ValueOf(time.Now().UnixMilli()))
	}
	if whereModel != nil {
		_, retErr = coll.UpdateAll(context.Background(), whereModel, bson.M{"$set": model}, options.UpdateOptions{UpdateOptions: moptions.Update().SetUpsert(true)})
	} else {
		_, retErr = coll.UpdateAll(context.Background(), bson.M{}, bson.M{"$set": model}, options.UpdateOptions{UpdateOptions: moptions.Update().SetUpsert(true)})
	}
	global.Debug("DBSet handle successed", "collName", collName, "model", model, "error", retErr)
	return
}

func DBGet[T any](collName string, whereModel any) (model *T, retErr error) {
	if DisableMongoDB {
		return new(T), nil
	}
	model = new(T)
	coll := GetDB().Collection(collName)
	if whereModel != nil {
		retErr = coll.Find(context.Background(), whereModel).One(model)
	} else {
		retErr = coll.Find(context.Background(), bson.M{}).One(model)
	}
	global.Debug("DBGet handle successed", "collName", collName, "whereModel", whereModel, "model", model, "error", retErr)
	return
}

func DBUpdate[T any](collName string, model, whereModel any) (retErr error) {
	if DisableMongoDB {
		return nil
	}
	// handle for updated_at
	rmodel := reflect.ValueOf(model)
	if rmodel.Kind() == reflect.Ptr {
		rmodel.Elem().FieldByName("UpdatedAt").Set(reflect.ValueOf(time.Now().UnixMilli()))
	}

	coll := GetDB().Collection(collName)
	if whereModel != nil {
		retErr = coll.UpdateOne(context.Background(), whereModel, bson.M{"$set": model})
	} else {
		retErr = coll.UpdateOne(context.Background(), bson.M{}, bson.M{"$set": model})
	}
	global.Debug("DBUpdate handle successed", "collName", collName, "whereModel", whereModel, "model", model, "error", retErr)
	return
}

// 优先从数据库获取数据，如果失败，请求缓存
func DBList[T any](collName string, whereModel any, bCache bool) (models []*T, retErr error) {
	defer func() {
		if retErr == nil && bCache {
			if crets := CacheGet(CacheListPageKey(collName, int64(len(models)), 0, 0), false, nil); crets != nil {
				models = crets.([]*T)
			}
		}
	}()
	if DisableMongoDB {
		return make([]*T, 0), nil
	}
	models = make([]*T, 0)
	coll := GetDB().Collection(collName)
	if whereModel != nil {
		retErr = coll.Find(context.Background(), whereModel).All(&models)
	} else {
		retErr = coll.Find(context.Background(), bson.M{}).All(&models)
	}
	if retErr != nil && retErr.Error() == "mongo: no documents in result" {
		return models, nil
	}
	// global.Debug("DBGetList handle successed", "collName", collName, "whereModel", whereModel, "models", models, "error", retErr)
	return
}

// 优先从缓存获取数据，如果失败，请求数据库
func DBListPage[T any](collName string, whereModel any, order, orderField string, offset, limit int, bCache bool) (models []*T, total int64, retErr error) {
	defer func() {
		if retErr == nil && bCache {
			if err := CacheSet(CacheListPageKey(collName, total, limit, offset), models, true); err != nil {
				slog.Error("CacheSet failed after DBListPage", "collName", collName, "whereModel", whereModel, "error", err.Error())
			}
		}
	}()
	total = DBCount(collName, whereModel)
	if bCache {
		if crets := CacheGet(CacheListPageKey(collName, total, limit, offset), false, nil); crets != nil {
			return crets.([]*T), total, nil
		}
	}
	if DisableMongoDB {
		return make([]*T, 0), 0, nil
	}
	coll := GetDB().Collection(collName)

	var collQuery qmgo.QueryI
	if whereModel != nil {
		collQuery = coll.Find(context.Background(), whereModel, options.FindOptions{
			QueryHook: new(T),
		})
	} else {
		collQuery = coll.Find(context.Background(), bson.M{}, options.FindOptions{
			QueryHook: new(T),
		})
	}

	if offset > 0 {
		collQuery = collQuery.Skip(int64(offset-1) * int64(limit))
	}
	if limit > 0 {
		collQuery = collQuery.Limit(int64(limit))
	}
	if order != "" {
		if order == "desc" {
			order = "-" + orderField
		} else {
			order = orderField
		}
		collQuery = collQuery.Sort(order)
	}
	models = make([]*T, 0)
	retErr = collQuery.All(&models)
	if retErr != nil && retErr.Error() == "mongo: no documents in result" {
		return models, 0, nil
	}

	// global.Debug("DBGetPageList handle successed", "whereModel", whereModel, "models", models, "retErr", retErr)
	return
}

func DBCount(collName string, whereModel any) (count int64) {
	if DisableMongoDB {
		return 0
	}
	var handleErr error
	coll := GetDB().Collection(collName)
	if whereModel != nil {
		count, handleErr = coll.Find(context.Background(), whereModel).Count()
	} else {
		count, handleErr = coll.Find(context.Background(), bson.M{}).Count()
	}
	if handleErr != nil {
		count = 0
		slog.Error("DBCount handle successed", "error", handleErr)
	}
	return
}

func DBDelete(collName string, whereModel any) (retErr error) {
	if DisableMongoDB {
		return nil
	}
	coll := GetDB().Collection(collName)
	if whereModel != nil {
		retErr = coll.Remove(context.Background(), whereModel)
	} else {
		return errors.New("invalid condition to delete model")
	}
	global.Debug("DBDelete handle successed", "whereModel", whereModel, "retErr", retErr)
	return
}

func DBAggregate[T any](collName string, pipeline any) ([]*T, error) {
	if DisableMongoDB {
		return make([]*T, 0), nil
	}

	ret := make([]*T, 0)
	if err := GetDB().Collection(collName).Aggregate(context.Background(), pipeline).All(&ret); err != nil {
		return nil, err
	}
	return ret, nil
}

func UpdateValues() {

	// logs, _ := DBList[model.TX](model.C_TX, bson.M{}, false)
	// for _, log := range logs {
	// 	DBSet(model.C_TX, log, bson.M{"_id": log.GID})
	// }

	// _, err := GetDB().Collection(model.C_Game).UpdateAll(context.Background(), bson.M{"status": "LOSE"}, bson.M{"$set": bson.M{"playresult": decimal.NewFromFloat(0)}})
	// slog.Info("update------", "error", err)
}

func RegistCallback(callback func(ctx context.Context, vmodel interface{}, opType operator.OpType, opts ...interface{}) error) {
	middleware.Register(callback)
}
