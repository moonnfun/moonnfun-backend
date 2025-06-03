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
)

func GetUser(address any) (*model.User, error) {
	return store.DBGet[model.User](model.C_User, bson.M{"address": address})
}

func UserLogin(address, signature, message string) (*model.User, error) {
	if err := verifyWalletSignature(address, signature, message); err != nil {
		return nil, err
	}

	var user *model.User
	if u, err := GetUser(address); err != nil {
		user = &model.User{
			Address: address,
		}
	} else {
		user = u
	}

	// session
	store.CacheSetByTime(user.Address, user, true, time.Duration(global.Config.WebSessionTimeout)*time.Second, func(val any) bool {
		go SaveUser(val.(*model.User))
		return true
	})
	slog.Info("user login successed", "user", user)
	return user, nil
}

func RemoveUser(userID any) {
	store.CacheGet(userID, true, func(v any) bool {
		go SaveUser(v.(*model.User))
		return true
	})
}

func SaveUser(user *model.User) error {
	if err := store.DBSet(model.C_User, user, bson.M{"_id": user.DBID}); err != nil {
		slog.Error("update user failed", "user", user, "error", err.Error())
		return err
	}
	slog.Info("update user successed", "user", user)
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
