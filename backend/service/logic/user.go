package logic

import (
	"fmt"
	"log/slog"
	"meme3/global"
	"meme3/service/model"
	"meme3/service/store"
	"time"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetUser(address any, bCache bool) (*model.User, error) {
	// if bCache {
	// 	u := store.CacheGet(address, false, nil)
	// 	if u == nil {
	// 		return nil, fmt.Errorf("can not find user: %s", address)
	// 	}
	// 	return u.(*model.User), nil
	// }
	return store.DBGet[model.User](model.C_User, bson.M{"address": address})
}

func GetUserByRefer(refer string) (*model.User, error) {
	return store.DBGet[model.User](model.C_User, bson.M{"refferalid": refer})
}

func UserLogin(address, signature, message, refer string) (*model.User, error) {
	if err := verifyWalletSignature(address, signature, message); err != nil {
		return nil, err
	}

	bSave := false
	bUpdate := false
	var user *model.User
	if u, err := GetUser(address, false); err != nil {
		user = &model.User{
			Address:    address,
			RefferalID: store.NewId(),
		}
		bSave = true
		user.DBID = primitive.NewObjectID()
		user.CreatedAt = time.Now().UnixMilli()
		user.RefferalUrl = fmt.Sprintf("%s?refer=%s", global.Config.HostURL, user.RefferalID)
	} else {
		user = u
		if user.RefferalID == "" {
			bSave = true
			bUpdate = true
			user.RefferalID = store.NewId()
			user.RefferalUrl = fmt.Sprintf("%s?refer=%s", global.Config.HostURL, user.RefferalID)
		}
	}

	if bSave {
		SaveUser(refer, user, bUpdate)
	}

	// // session
	// store.CacheSetByTime(user.Address, user, true, time.Duration(global.Config.WebSessionTimeout)*time.Second, func(val any) bool {
	// 	// go SaveUserToDB(val.(*model.User), true)
	// 	return true
	// })

	slog.Info("user login successed", "user", user)
	return user, nil
}

func SaveUser(refer string, user *model.User, bUpdate bool) error {
	global.Debug("before save user", "refer", refer, "user", user)
	if refer != "" && user.FromID == "" && refer != user.RefferalID {
		fromUser, err := GetUserByRefer(refer)
		if err != nil {
			return err
		}
		if err := SaveRefer(fromUser, user); err != nil {
			return err
		}
		user.FromID = fromUser.RefferalID
	} else {
		slog.Warn("can not bind user again and again", "refer", refer)
	}
	return SaveUserToDB(user, bUpdate)
}

func RemoveUser(userID any) {
	// store.CacheGet(userID, true, func(v any) bool {
	// 	// go SaveUserToDB(v.(*model.User), true)
	// 	return true
	// })
}

func SaveUserToDB(user *model.User, bUpdate bool) error {
	dbID := user.DBID
	if bUpdate {
		user.DBID = primitive.ObjectID{}
	}
	if err := store.DBSet(model.C_User, user, bson.M{"_id": dbID}); err != nil {
		slog.Error("update user failed", "user", user, "error", err.Error())
		return err
	}
	slog.Info("update user successed", "user", user)
	return nil
}

func SaveRefer(fromUser, referUser *model.User) error {
	referral := &model.Refferal{
		Wallet:  fromUser.Address,
		Address: referUser.Address,
	}
	referral.DBID = primitive.NewObjectID()
	referral.CreatedAt = time.Now().UnixMilli()
	if err := store.DBSet(model.C_Referral, referral, bson.M{"wallet": fromUser.Address, "address": referUser.Address}); err != nil {
		slog.Error("save refer failed", "refer", referral, "fromUser", fromUser, "referUser", referUser, "error", err.Error())
		// return err
	} else {
		slog.Info("save refer successed", "user", referral)
	}
	return nil
}

func UpdateUser(address string) error {
	// if u, err := GetUser(address, true); err != nil {
	// 	return err
	// } else {
	// 	u.TotalTrading += 1
	// 	if err := store.DBSet(model.C_User, bson.M{"totaltrading": u.TotalTrading}, bson.M{"_id": u.DBID}); err != nil {
	// 		slog.Error("update user failed", "user", u, "error", err.Error())
	// 		return err
	// 	}
	// }
	return nil
}

func verifyWalletSignature(address, signatureHex, message string) error {
	signature, err := hexutil.Decode(signatureHex)
	if err != nil {
		return err
	}

	signature[crypto.RecoveryIDOffset] -= 27 // Transform yellow paper V from 27/28 to 0/1

	messageHash := accounts.TextHash([]byte(message))

	pubKey, err := crypto.SigToPub(messageHash, signature)
	if err != nil {
		return err
	}

	if common.HexToAddress(address) != crypto.PubkeyToAddress(*pubKey) {
		return fmt.Errorf("failed to verify signature")
	}
	return nil
}
