package store

import (
	"context"
	"fmt"
	"reflect"

	"go.mongodb.org/mongo-driver/bson/bsoncodec"
	"go.mongodb.org/mongo-driver/bson/bsonrw"
	"go.mongodb.org/mongo-driver/bson/bsontype"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/event"

	"github.com/qiniu/qmgo"
	"github.com/qiniu/qmgo/options"
	"github.com/shopspring/decimal"
	"go.mongodb.org/mongo-driver/bson"
	moptions "go.mongodb.org/mongo-driver/mongo/options"
)

var (
	mdb            *qmgo.Database
	mClient        *qmgo.Client
	DisableMongoDB bool
)
var Coptions = &moptions.CreateCollectionOptions{}

type Mongo struct {
	URI           string
	DBName        string
	UserName      string
	AuthDBName    string
	Password      string
	AuthMechanism string
}

func InitMongoDB(mConf *Mongo) error {
	reg := bson.NewRegistryBuilder()
	reg.RegisterCodec(reflect.TypeOf(decimal.Decimal{}), Decimal{})
	opt := moptions.Client().SetRegistry(reg.Build())

	//打印command
	// startedCommands := make(map[int64]bson.Raw)
	cmdMonitor := &event.CommandMonitor{
		Started: func(_ context.Context, evt *event.CommandStartedEvent) {
			// startedCommands[evt.RequestID] = evt.Command
		},
		Succeeded: func(_ context.Context, evt *event.CommandSucceededEvent) {
			// slog.Debug("InitMongoDB successed", "cmd", startedCommands[evt.RequestID].String())
		},
		Failed: func(_ context.Context, evt *event.CommandFailedEvent) {
			// slog.Debug("InitMongoDB failed", "cmd", startedCommands[evt.RequestID].String())
		},
	}
	opt = opt.SetMonitor(cmdMonitor)
	var credential *qmgo.Credential
	if mConf.Password != "" {
		credential = &qmgo.Credential{
			AuthMechanism: mConf.AuthMechanism,
			AuthSource:    mConf.AuthDBName,
			Username:      mConf.UserName,
			Password:      mConf.Password,
		}
		if credential.AuthMechanism == "GSSAPI" {
			credential.PasswordSet = true
		}
	}
	client, err := qmgo.NewClient(
		context.Background(),
		&qmgo.Config{Uri: mConf.URI, Auth: credential},
		options.ClientOptions{ClientOptions: opt},
	)
	if err != nil {
		return err
	}
	mClient = client
	mdb = client.Database(mConf.DBName)

	// init
	Coptions.Collation = &moptions.Collation{
		Locale:          "en_US",
		CaseLevel:       false,
		CaseFirst:       "off",
		Strength:        2,
		NumericOrdering: false,
		Alternate:       "non-ignorable",
		MaxVariable:     "punct",
		Backwards:       false,
		Normalization:   false,
	}
	return nil
}

func CloseMongoDB() {
	mClient.Close(context.Background())
}

func GetDB() *qmgo.Database {
	return mdb
}

func GetClient() *qmgo.Client {
	return mClient
}

// overwrite decimal
type Decimal decimal.Decimal

func (d Decimal) DecodeValue(dc bsoncodec.DecodeContext, vr bsonrw.ValueReader, val reflect.Value) error {
	decimalType := reflect.TypeOf(decimal.Decimal{})
	if !val.IsValid() || !val.CanSet() || val.Type() != decimalType {
		return bsoncodec.ValueDecoderError{
			Name:     "decimalDecodeValue",
			Types:    []reflect.Type{decimalType},
			Received: val,
		}
	}

	var value decimal.Decimal
	switch vr.Type() {
	case bsontype.Decimal128:
		dec, err := vr.ReadDecimal128()
		if err != nil {
			return err
		}
		value, err = decimal.NewFromString(dec.String())
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("received invalid BSON type to decode into decimal.Decimal: %s", vr.Type())
	}

	val.Set(reflect.ValueOf(value))
	return nil
}

func (d Decimal) EncodeValue(ec bsoncodec.EncodeContext, vw bsonrw.ValueWriter, val reflect.Value) error {
	decimalType := reflect.TypeOf(decimal.Decimal{})
	if !val.IsValid() || val.Type() != decimalType {
		return bsoncodec.ValueEncoderError{
			Name:     "decimalEncodeValue",
			Types:    []reflect.Type{decimalType},
			Received: val,
		}
	}

	dec := val.Interface().(decimal.Decimal)
	dec128, err := primitive.ParseDecimal128(dec.String())
	if err != nil {
		return err
	}

	return vw.WriteDecimal128(dec128)
}
